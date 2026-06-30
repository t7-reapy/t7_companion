# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

`t7kb` is the shipped, **agent-first query tool** for a Black Ops 3 modding knowledge base: a single **pure-Go, no-cgo** binary that serves hybrid retrieval (FTS5 BM25 + vector cosine, fused with RRF) over a local SQLite index, `t7kb.db`. The primary surface is an **MCP server** for AI agents; it also has a small CLI. Everything runs locally and offline.

This repo is the **consumer**. It does NOT build `t7kb.db` — that is produced by the sibling **t7_knowledge** repo (see "Two-repo split").

## Repo layout

The repo is three things in one: the **Go tool**, a **Claude Code plugin**, and **install/docs**. Map of what lives where (so nothing reads as accidental duplication):

- `cmd/` + `internal/` — the Go tool (`t7kb`): CLI + MCP server. The only compiled artifact.
- `install/install.sh` + `install.ps1` — the **single source of truth for installation**. Everything else that mentions install just invokes these.
- `.claude-plugin/marketplace.json` — the **marketplace** catalog (what makes `/plugin marketplace add t7-reapy/t7_companion` work). Lists the plugin below.
- `plugin/.claude-plugin/plugin.json` — the **plugin** manifest itself. *Two different files on purpose:* marketplace = the catalog, plugin = the thing in it. Not a duplicate.
- `plugin/skills/*/SKILL.md` — the plugin's skills. `setup` is the **agent's** install action (it runs `install/install.sh` + `claude mcp add`) — i.e. the agent-facing mirror of the README's install, not a separate install method.
- `templates/AGENTS.md` — the vendor-neutral primer a user drops in their *map project*.
- `README.md` — the lean **human-facing** front door: install + a one-line connect pointer + CLI.
- `docs/clients.md` — per-client MCP config (Codex/OpenCode/Cursor/Copilot/Claude) + the AGENTS.md editor table. The detail the README points at; lives only here.
- `CLAUDE.md` (this file) / `NOTICE.md` / `docs/data-model.md` — contributor + licensing + schema docs.

Install is intentionally described for two audiences — the README (human) and the `setup` skill (agent) — but both call the same `install/` scripts, so there is one real source.

## Commands

Targets **Go 1.25** (`go.mod`). The machine's default `go` may be older; `GOTOOLCHAIN=auto` (the default) fetches the right toolchain on first build.

- Build: `go build -o t7kb.exe ./cmd/t7kb`
- Cross-compile (what CI ships): `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ./cmd/t7kb`
- Vet: `go vet ./...`
- Test: `go test ./...`
- Single test: `go test ./internal/store -run TestHybridSearch -v`

Any hybrid search (and `TestHybridSearch`) loads the embedding model via go-sentex. By default the model is looked up in a `models/` dir **beside the binary** (the release bundle); for dev, set `HF_HOME` to your global HuggingFace cache or it downloads ~87 MB on first use. `search --bm25` skips the embedder.

Run surfaces:
- `t7kb mcp` — stdio MCP server (primary; exposes `search` + `get` tools)
- `t7kb search "<query>"` / `t7kb get <doc_id>` — one-shot CLI (`--scores` shows RRF + reliability)
- `t7kb` (no args) — interactive browse: query → numbered hits → type a number to read one
- `--version` reports the build version (injected via ldflags; "dev" locally)
- DB resolution (`resolveDB`): `--db` > `$T7KB_DB` > beside the binary > `./t7kb.db`. It returns the intended path even if the file is absent — `ensureDB` (`internal/cli/db.go`) then unpacks a sibling `t7kb.db.zip` there on first run.

## Architecture

**Two-repo split (read this first).** The corpus and the DB build live upstream in the sibling `t7_knowledge` repo; the chain is `sources → build_kb.py → kb.db → make_t7kb.py → t7kb.db`. Do **not** add corpus ingestion or batch-embedding logic here — it belongs upstream. The ~3.5 GB `t7kb.db` exceeds GitHub's 2 GiB asset limit, so it ships **compressed as a separate `t7kb.db.zip` release asset** (the OS archives don't contain it); the binary auto-unpacks it beside itself on first run (see `ensureDB`).

**The embedding-parity invariant (the #1 gotcha).** Vector search only works if the query is embedded with the **same model** that produced the stored chunk vectors. Both sides use **all-MiniLM-L6-v2 (384-d, L2-normalized)** — the corpus via Python fastembed upstream, the query via **go-sentex** here — verified to produce identical vectors. If you change the model here you MUST rebuild the corpus upstream (`corpus.py` `EMBEDDING_MODEL` / `EMBEDDING_DIM`). Watch out: `store.vectorRank` silently **skips** chunks whose dimension ≠ the query's, so a model/dim mismatch degrades to BM25-only with no error. The corpus is also chunked for MiniLM's ~256-token window, so the model choice is load-bearing on both chunking and vectors.

**Pure-Go, no cgo — a hard constraint.** It is what makes the single static binary and clean cross-compile work (there is no C compiler on the build host). SQLite is `modernc.org/sqlite` (not mattn); embeddings are `go-sentex`/gomlx (not onnxruntime). Don't introduce cgo deps. (Aside: the corpus build embeds in Python because the pure-Go gomlx backend is ~3 chunks/s — fine for one query embed, far too slow for the ~350K-chunk corpus build.)

**Retrieval pipeline.** `internal/store` owns it: `bm25Rank` (FTS5 `bm25()`), `vectorRank` (streams the `embeddings` table; cosine = dot since normalized; best chunk per doc), `rrfFuse` (k=60), `sortByFusedScore` (reliability tiebreak). `internal/embed` wraps go-sentex (one query embedding per search) and defaults `HF_HOME` to the beside-binary `models/`. `internal/cli` builds the cobra tree.

**MCP shape.** The `search` tool returns ranked `doc_id` / `title` / `source` / `reliability` / snippet — deliberately NOT the RRF/vector internals (noise to an agent; reliability is the one ranking signal it gets). RRF scores are CLI-only.

The DB schema (`documents`, `docs_fts`, `embeddings`) is in `docs/data-model.md`; per-row `source` + `url` carry attribution (`NOTICE.md`).

## Releasing

Pushing a `v*` tag runs GoReleaser (`.goreleaser.yaml` + `.github/workflows/release.yml`): it cross-builds the binaries, runs a before-hook that downloads the embedding model into `models/`, and bundles binary + model + docs into per-platform archives. The version is injected via `-ldflags` into `internal/cli.version`.

- The auto-changelog is **disabled** on purpose — GoReleaser's git changelog leaks full SHAs + author emails. Release notes are the curated `header` / `footer` in `.goreleaser.yaml`; edit those, not a changelog config.
- `t7kb.db.zip` is **not** built by CI — attach it to the release manually (`gh release upload <tag> t7kb.db.zip`), since it's large and built upstream.
- Validate config changes with `goreleaser check`; dry-run with `goreleaser release --snapshot --clean`.

## Conventions

- Conventional commits (`feat:`, `fix:`, `refactor:`, `chore:`, scopes like `feat(cli):`) — they drive the version bump intent even though the changelog body is curated.
- `*.db`, `*.db.zip`, and `models/` are build/ship artifacts — gitignored.
- Never hard-wrap markdown at 80 columns (or any fixed width). One line per paragraph / list item; let the editor soft-wrap. The maintainer is allergic to fixed-width reflow.
