# t7_companion

:robot: Rich LLM-loadable collection of capabilities for **Black Ops 3 modding** — a synthesized wiki, per-domain skills, a portable search index (`t7kb.db`), and a consumer CLI / MCP server / Claude Code plugin. Try it, keep it :chains:.

> [!IMPORTANT]
> **Status**: scaffolding. Today this repo holds the licensing structure; the wiki content, `t7kb.db`, and CLI/MCP code land progressively as the project matures.

## Repo layout

```
t7_companion/
├── src/                          CLI / MCP / plugin code (MIT) — not yet started
├── wiki/                         Synthesized markdown wiki (CC-BY-SA 4.0)
│   └── LICENSE.md                  wiki content license attestation
├── skills/                       Per-domain LLM-loadable skill files (CC-BY-SA 4.0) — not yet started
├── docs/                         Contributor reference
│   ├── data-model.md               t7kb.db schema (what each row looks like)
│   └── add-remove-knowledge.md     Adding / removing content (DMCA-ready)
├── LICENSES/                     Bundled third-party license texts + notices
│   ├── GPL-3.0.txt                 gscode-api upstream license
│   ├── CC-BY-SA-4.0.txt            wiki / skills / docs content license
│   ├── fair-use-notice.md          Treyarch / YouTube fair-use posture
│   └── modder-attribution.md       source_workspace per-modder attribution posture
├── LICENSE                       MIT for the code in this repo
├── NOTICE.md                     polyglot-licensing top-level summary (canonical reference)
└── README.md                     this file
```

> [!NOTE]
> `t7kb.db` itself is **not committed** to this repo — it ships as a **GitHub Release asset** downloaded by the consumer at install / update time. The repo holds the licensing scaffold, the wiki sources, and the consumer code; the search index is built upstream and published per release.

## Licensing — quick summary

This is a *polyglot* repo. Each component has its own license:

- **Code** (`src/`) → **MIT**
- **Wiki / skills / docs prose** (`wiki/`, `skills/`, `docs/`) → **CC-BY-SA 4.0**
- **`t7kb.db`** → an aggregate; each row keeps its own posture (GPLv3 / fair-use / per-modder / CC-BY / etc.)

For the full per-source breakdown, see **[`NOTICE.md`](NOTICE.md)** — that's the canonical reference. The `LICENSES/` directory holds the bundled third-party license texts and our fair-use notices.

> [!CAUTION]
> **DMCA-ready**: if any rights-holder objects to as-is inclusion of their content, the affected `t7kb.db` rows can be flipped to **pointer-only** (URL/path + metadata, no body) without affecting unrelated rows. See [`docs/add-remove-knowledge.md`](docs/add-remove-knowledge.md) for the flow.

## Contributing / reporting issues

A full `CONTRIBUTING.md` will land alongside the first wiki content. In the meantime: wiki edits, license-posture concerns, attribution corrections, and content-removal requests are all welcome — open an issue or contact the maintainer [McReaper](https://github.com/McReaper) directly.

> [!TIP]
> Want to fix a wiki page? Open a PR against `wiki/*.md` — pages are CC-BY-SA 4.0 and contributions are welcome under the same license. `t7kb.db` rows are built deterministically from the upstream source-set and aren't directly editable here (see [`docs/add-remove-knowledge.md`](docs/add-remove-knowledge.md) for how source-data changes happen).
