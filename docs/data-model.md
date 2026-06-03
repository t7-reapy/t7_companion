# Data model — `t7kb.db`

`t7kb.db` is a single-file SQLite database that bundles every searchable entry shipped with `t7_companion`. It is built upstream from the BO3 community corpus and published as a **GitHub Release asset** alongside each `t7_companion` release; the consumer downloads it at install or update time.

This page describes its schema, the per-row fields each consumer can rely on, and how license posture is recorded per row.

## High-level shape

Two tables:

- **`documents`** — one row per searchable entry. Holds the body text, metadata, and license info.
- **`docs_fts`** — a [SQLite FTS5](https://www.sqlite.org/fts5.html) virtual table over `documents` in external-content mode, providing BM25 ranking over `title + summary + body`. No duplicate storage; FTS5 reads bodies directly from `documents`.

> [!NOTE]
> **Hybrid retrieval is the design target.** The schema reserves a nullable `embedding BLOB` column on `documents` so a vector layer (`sqlite-vec` or equivalent) can be added additively without re-designing — combined with BM25 via [Reciprocal Rank Fusion](https://learn.microsoft.com/en-us/azure/search/hybrid-search-ranking) for synonym-tolerant retrieval on top of exact-term matching. BM25 ships first; the vector layer comes in a later release once empirical recall warrants it.

## `documents` schema

```sql
CREATE TABLE documents (
    doc_id        TEXT PRIMARY KEY,   -- '<source>::<source_local_id>'
    source        TEXT NOT NULL,      -- e.g. 'gscode-api', 'source_scripts', 'wiki_page'
    title         TEXT NOT NULL,
    summary       TEXT,               -- short paraphrased blurb (~1 line)
    body          TEXT NOT NULL,      -- full content; '' when row is pointer-only
    url           TEXT,               -- canonical upstream URL when public
    content_path  TEXT,               -- local-install path when applicable
    themes        TEXT,               -- JSON array of taxonomy themes
    metadata      TEXT,               -- JSON array of free-form tags
    source_type   TEXT,               -- 'api', 'script', 'wiki_page', 'transcript_segment', ...
    reliability   REAL NOT NULL,      -- [0.0, 1.0] — soft tiebreaker for ranking
    license       TEXT,               -- per-row license (see "License posture" below)
    upstream      TEXT,               -- JSON: { "repo": "...", "commit": "...", "author": "..." }
    embedding     BLOB                -- reserved for the future vector layer
);

CREATE VIRTUAL TABLE docs_fts USING fts5(
    doc_id UNINDEXED, title, summary, body,
    content='documents', content_rowid='rowid',
    tokenize='unicode61 remove_diacritics 2'
);

CREATE INDEX idx_documents_source      ON documents(source);
CREATE INDEX idx_documents_reliability ON documents(reliability);
CREATE INDEX idx_documents_source_type ON documents(source_type);
CREATE INDEX idx_documents_license     ON documents(license);
```

## `doc_id` format

Every row's id is `<source>::<source_local_id>`. The `source` prefix prevents collisions across sources; the `source_local_id` is whatever scheme the source's own ingest pipeline assigned.

Examples:

- `gscode-api::api.gsc.setclientfield` — a GSC built-in function reference
- `source_scripts::t7src.gsc.scripts.zm._zm_perks` — a decompiled engine script file
- `source_workspace::workspace.usermaps.zm_test.scripts.zm.zm_cellbreaker.csc` — a per-usermap script in the maintainer's workspace (attribution lives in the workspace `README.md` credits table)
- `wiki_page::scripting.clientfields` — a synthesized wiki page
- `video-youtube::SomeTutorialName.seg.04` — a transcript segment with a deep-link URL
- `docs-bo3::bo3-docs-gsc-language` — a Treyarch reference document

> [!IMPORTANT]
> Sources that are **synthesized into wiki pages** (Discord, forums, community wikis) do **not** appear as their own rows in `t7kb.db`. They flow into the wiki via the synthesis pipeline; the resulting wiki page rows cite them in their `## Sources` block, but the original raw entries themselves are not shipped.

Consumers never construct doc_ids by hand — every `search` result hands them the `doc_id` to pass to `get`.

## License posture per row

The `license` field is a short SPDX-like string capturing the row's authoritative license or fair-use posture:

| `license` value | Meaning | Bundled doc |
|---|---|---|
| `MIT` / `BSD-3-Clause` / `Apache-2.0` / etc. | Standard SPDX identifier — upstream's own license | upstream repo |
| `GPL-3.0` | gscode-api entries | [`../LICENSES/GPL-3.0.txt`](../LICENSES/GPL-3.0.txt) |
| `CC-BY-SA-4.0` | Synthesized wiki / skill / docs rows | [`../LICENSES/CC-BY-SA-4.0.txt`](../LICENSES/CC-BY-SA-4.0.txt) |
| `CC-BY-4.0` | `personal-notes` rows (maintainer owns the IP and chose attribution-only redistribution) | maintainer-set |
| `treyarch-fair-use` | `source_scripts`, `source_dump`, `docs-bo3` rows | [`../LICENSES/fair-use-notice.md`](../LICENSES/fair-use-notice.md) |
| `youtube-fair-use` | `video-youtube` transcript-segment rows | [`../LICENSES/fair-use-notice.md`](../LICENSES/fair-use-notice.md) |
| `modder-attribution` | `source_workspace` third-party modder rows | [`../LICENSES/modder-attribution.md`](../LICENSES/modder-attribution.md) |
| `maintainer-owned` | `source_workspace` paths under `_reapy/` and top-level `KNOWLEDGE.md` / `README.md` | maintainer-set |
| `maintainer-owned-with-community-imports` | `source_workspace` paths under `usermaps/` — the maintainer's own map scripts/configs but with imported community asset packs referenced inside; attribution for the imported packs lives in the workspace `README.md` credits table | maintainer-set |

Consumers can:

- **Filter** results by license posture (e.g. *"only show me MIT/GPL — I want to copy code freely"*).
- **Display** the license shorthand alongside each result so users see the posture at a glance.
- **Cite** the license file via [`../NOTICE.md`](../NOTICE.md) when redistributing.

## Pointer-only rows

When a row has been flipped to pointer-only (e.g. via a removal request — see [`add-remove-knowledge.md`](add-remove-knowledge.md)), `body` becomes an empty string and the row keeps everything else (title, url, content_path, metadata). Search still finds the row via its title/summary tokens; the `body` column just doesn't contribute matches or readable content.

Consumers should treat pointer-only rows as **"here's where to find this — go read it upstream"**: surface `title`, `url`, and `summary`, and mark the result clearly as a pointer.

## Build provenance

Each `t7kb.db` release ships with a sibling `t7kb.manifest.json` describing:

- DB build timestamp + content hash
- Source-set version (which scrape commits each source was built from)
- Per-source row count and total size
- Schema version (so consumers can detect breaking changes)

The manifest is small (~1 KB) and committed to the release alongside the DB file.
