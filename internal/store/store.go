// Package store opens a t7kb.db and serves hybrid retrieval over it:
// FTS5 BM25 + cosine over precomputed embeddings, fused with RRF.
package store

import (
	"context"
	"database/sql"
	"encoding/binary"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"

	_ "modernc.org/sqlite" // pure-Go SQLite driver (no cgo); FTS5 compiled in
)

// RRFK is the Reciprocal Rank Fusion constant (matches query_kb.py).
const RRFK = 60

// Pool is the per-retriever candidate count fused before truncating to limit.
const Pool = 50

// Store is a read-only handle to a t7kb.db.
type Store struct {
	db *sql.DB
}

// Hit is one fused search result.
type Hit struct {
	DocID   string
	Title   string
	Source  string
	Score   float64 // fused RRF score (higher = better)
	Snippet string
}

// Doc is a full document body.
type Doc struct {
	DocID       string
	Source      string
	Title       string
	Summary     string
	Body        string
	URL         string
	ContentPath string
	Themes      string
	Metadata    string
	SourceType  string
	Reliability float64
}

// Open opens the database at path in read-only mode.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", "file:"+path+"?mode=ro")
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping %s: %w", path, err)
	}
	return &Store{db: db}, nil
}

// Close releases the handle.
func (s *Store) Close() error { return s.db.Close() }

var ftsToken = regexp.MustCompile(`[\p{L}\p{N}_]+`)

// ftsQuery turns free text into a safe FTS5 MATCH expression: each word
// becomes a quoted term, joined implicitly (AND). Empty if no usable tokens.
func ftsQuery(q string) string {
	toks := ftsToken.FindAllString(q, -1)
	if len(toks) == 0 {
		return ""
	}
	quoted := make([]string, len(toks))
	for i, t := range toks {
		quoted[i] = `"` + t + `"`
	}
	return strings.Join(quoted, " ")
}

// rankItem is one retriever's ranked candidate.
type rankItem struct {
	docID   string
	snippet string
}

// bm25Rank returns up to limit doc_ids ranked by FTS5 BM25, with a snippet.
func (s *Store) bm25Rank(ctx context.Context, query string, limit int) ([]rankItem, error) {
	match := ftsQuery(query)
	if match == "" {
		return nil, nil
	}
	const q = `
		SELECT d.doc_id, snippet(docs_fts, 3, '', '', ' … ', 12) AS snip
		FROM docs_fts
		JOIN documents d ON d.rowid = docs_fts.rowid
		WHERE docs_fts MATCH ?
		ORDER BY bm25(docs_fts)
		LIMIT ?`
	rows, err := s.db.QueryContext(ctx, q, match, limit)
	if err != nil {
		return nil, fmt.Errorf("bm25 search: %w", err)
	}
	defer rows.Close()
	var out []rankItem
	for rows.Next() {
		var it rankItem
		var snip sql.NullString
		if err := rows.Scan(&it.docID, &snip); err != nil {
			return nil, err
		}
		it.snippet = snip.String
		out = append(out, it)
	}
	return out, rows.Err()
}

func decodeVec(b []byte) []float32 {
	n := len(b) / 4
	v := make([]float32, n)
	for i := range v {
		v[i] = math.Float32frombits(binary.LittleEndian.Uint32(b[i*4:]))
	}
	return v
}

// vectorRank streams the embeddings table, scoring each chunk by dot product
// with qvec (cosine, since both are L2-normalized), keeping the best chunk per
// doc. Chunks whose dimension differs from qvec are skipped — so the tool
// degrades to BM25-only against a db embedded with a different model.
func (s *Store) vectorRank(ctx context.Context, qvec []float32, limit int) ([]rankItem, error) {
	if len(qvec) == 0 {
		return nil, nil
	}
	rows, err := s.db.QueryContext(ctx, `SELECT doc_id, chunk_text, embedding FROM embeddings`)
	if err != nil {
		// No embeddings table (or unreadable) → vector disabled, not fatal.
		return nil, nil
	}
	defer rows.Close()

	type best struct {
		score   float64
		snippet string
	}
	bestByDoc := map[string]best{}
	for rows.Next() {
		var docID string
		var chunkText sql.NullString
		var blob []byte
		if err := rows.Scan(&docID, &chunkText, &blob); err != nil {
			return nil, err
		}
		v := decodeVec(blob)
		if len(v) != len(qvec) {
			continue // dimension mismatch → skip
		}
		var dot float64
		for i := range qvec {
			dot += float64(qvec[i]) * float64(v[i])
		}
		if cur, ok := bestByDoc[docID]; !ok || dot > cur.score {
			bestByDoc[docID] = best{score: dot, snippet: snippetOf(chunkText.String)}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	ranked := make([]rankItem, 0, len(bestByDoc))
	type ds struct {
		doc string
		sc  float64
		sn  string
	}
	all := make([]ds, 0, len(bestByDoc))
	for d, b := range bestByDoc {
		all = append(all, ds{d, b.score, b.snippet})
	}
	sort.Slice(all, func(i, j int) bool { return all[i].sc > all[j].sc })
	if limit > 0 && len(all) > limit {
		all = all[:limit]
	}
	for _, a := range all {
		ranked = append(ranked, rankItem{docID: a.doc, snippet: a.sn})
	}
	return ranked, nil
}

func snippetOf(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	if len(s) > 160 {
		s = s[:160] + " …"
	}
	return s
}

// SearchHybrid fuses BM25 and vector rankings with RRF, applies reliability as
// a soft tiebreak, and returns the top `limit` hits. qvec may be nil/empty (or
// a different dimension than the db) — then it's BM25-only.
func (s *Store) SearchHybrid(ctx context.Context, query string, qvec []float32, limit int) ([]Hit, error) {
	bm25, err := s.bm25Rank(ctx, query, Pool)
	if err != nil {
		return nil, err
	}
	vec, err := s.vectorRank(ctx, qvec, Pool)
	if err != nil {
		return nil, err
	}

	// RRF over both rankings; remember the first snippet seen per doc.
	fused := map[string]float64{}
	snip := map[string]string{}
	add := func(items []rankItem) {
		for rank, it := range items {
			fused[it.docID] += 1.0 / float64(RRFK+rank+1)
			if snip[it.docID] == "" && it.snippet != "" {
				snip[it.docID] = it.snippet
			}
		}
	}
	add(bm25)
	add(vec)
	if len(fused) == 0 {
		return nil, nil
	}

	docIDs := make([]string, 0, len(fused))
	for d := range fused {
		docIDs = append(docIDs, d)
	}
	rel, title, source, err := s.meta(ctx, docIDs)
	if err != nil {
		return nil, err
	}
	sort.Slice(docIDs, func(i, j int) bool {
		di, dj := docIDs[i], docIDs[j]
		if fused[di] != fused[dj] {
			return fused[di] > fused[dj]
		}
		if rel[di] != rel[dj] {
			return rel[di] > rel[dj] // reliability tiebreak
		}
		return di < dj
	})
	if limit > 0 && len(docIDs) > limit {
		docIDs = docIDs[:limit]
	}

	hits := make([]Hit, 0, len(docIDs))
	for _, d := range docIDs {
		hits = append(hits, Hit{
			DocID: d, Title: title[d], Source: source[d],
			Score: fused[d], Snippet: snip[d],
		})
	}
	return hits, nil
}

// meta fetches title/source/reliability for a set of doc_ids.
func (s *Store) meta(ctx context.Context, docIDs []string) (rel map[string]float64, title, source map[string]string, err error) {
	rel = map[string]float64{}
	title = map[string]string{}
	source = map[string]string{}
	if len(docIDs) == 0 {
		return
	}
	ph := strings.TrimSuffix(strings.Repeat("?,", len(docIDs)), ",")
	args := make([]any, len(docIDs))
	for i, d := range docIDs {
		args[i] = d
	}
	rows, err := s.db.QueryContext(ctx,
		`SELECT doc_id, title, source, reliability FROM documents WHERE doc_id IN (`+ph+`)`, args...)
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var d, t, src string
		var r float64
		if err := rows.Scan(&d, &t, &src, &r); err != nil {
			return nil, nil, nil, err
		}
		title[d], source[d], rel[d] = t, src, r
	}
	return rel, title, source, rows.Err()
}

// Get returns the full document for a doc_id, or (nil, nil) if absent.
func (s *Store) Get(ctx context.Context, docID string) (*Doc, error) {
	const q = `
		SELECT doc_id, source, title, summary, body, url, content_path,
		       themes, metadata, source_type, reliability
		FROM documents WHERE doc_id = ?`
	var d Doc
	var summary, url, contentPath, themes, metadata, sourceType sql.NullString
	err := s.db.QueryRowContext(ctx, q, docID).Scan(
		&d.DocID, &d.Source, &d.Title, &summary, &d.Body, &url, &contentPath,
		&themes, &metadata, &sourceType, &d.Reliability,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get %s: %w", docID, err)
	}
	d.Summary = summary.String
	d.URL = url.String
	d.ContentPath = contentPath.String
	d.Themes = themes.String
	d.Metadata = metadata.String
	d.SourceType = sourceType.String
	return &d, nil
}
