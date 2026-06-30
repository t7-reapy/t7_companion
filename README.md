# t7kb

Agent-first search over a local **Black Ops 3 modding** knowledge base. A single
pure-Go binary serving hybrid retrieval (keyword + semantic) over a bundled
SQLite index (`t7kb.db`) — built to be driven by an AI agent over MCP, with a
small CLI for direct use. Everything runs locally and offline.

> **Status:** early. The retrieval core (BM25 + vector + RRF) works; the MCP
> server and packaged releases are in progress.

## Install

1. From the latest release, download your OS archive (`t7kb_windows_amd64.zip` or
   `t7kb_linux_amd64.tar.gz`) **and** `t7kb.db.zip`.
2. Extract the OS archive **at your Black Ops III root**.
3. Drop `t7kb.db.zip` into the same folder, next to the binary — it is unpacked
   automatically on first run (no manual decompress).
4. Point your agent at the MCP server: run `t7kb mcp` (stdio).

The OS archive carries the binary, the embedding model, and the docs; the database
ships as a separate `t7kb.db.zip` (it's large). Everything then runs locally and
offline.

## CLI

```
t7kb search "how do I add a custom perk"   # hybrid (keyword + semantic)
t7kb search --bm25 weaponfull attachment   # keyword only
t7kb get <doc_id>                          # print a document's full body
```

`--db PATH` overrides the database location (default: `$T7KB_DB`, then beside the
binary, then `./t7kb.db`).

## Build from source

Pure Go, no C compiler needed (`CGO_ENABLED=0`); a recent Go toolchain is fetched
automatically.

```
go build ./cmd/t7kb
```

## Licensing

The code here is **MIT** (see [`LICENSE`](LICENSE)). `t7kb.db` bundles knowledge
from the BO3 modding community; every row carries its `source` and `url` for
attribution. See [`NOTICE.md`](NOTICE.md).
