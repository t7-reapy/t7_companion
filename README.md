# t7_companion

:robot: Rich LLM-loadable collection of capabilities for **Black Ops 3
modding** — a synthesized wiki, a portable search index (`kb.db`), and
a consumer CLI / MCP server / Claude Code plugin. Try it, keep it
:chains:.

> **Status**: scaffolding. The wiki content and `kb.db` are produced
> upstream in [t7-reapy/t7_knowledge](https://github.com/t7-reapy/t7_knowledge)
> (the maintainer's ingestion pipeline over the BO3 community corpus).
> This repo will be the **public, shippable** artifact once synthesis
> matures. Today it just holds the licensing scaffold + wiki domain
> skeleton.

## Repo layout

```
t7_companion/
  src/                     CLI / MCP / plugin code (MIT) — not yet started
  wiki/                    Synthesized markdown wiki (CC-BY-SA 4.0)
    ai/                    Per-domain page directories (placeholder skeleton)
    asset-pipeline/
    audio/
    ...
    LICENSE.md             ← wiki content license attestation (CC-BY-SA 4.0)
  kb.db                    Shipped search index (built by upstream pipeline)
  LICENSE                  ← MIT for the code in this repo
  LICENSES/                ← bundled third-party license texts and notices
    GPL-3.0.txt              gscode-api upstream license
    CC-BY-SA-4.0.txt         wiki content license (full text)
    fair-use-notice.md       Treyarch IP / YouTube transcript posture
    modder-attribution.md    source_workspace per-modder attribution posture
  NOTICE.md                ← polyglot-licensing top-level summary
  README.md                ← this file
```

## Licensing at a glance

This repo bundles content from multiple sources, each under its own
license or use posture (a *polyglot* repository). Full breakdown in
[`NOTICE.md`](NOTICE.md), but the quick version:

| Component | License |
|---|---|
| CLI / tooling code (`src/`) | **MIT** ([`LICENSE`](LICENSE)) |
| Wiki pages (`wiki/`) | **CC-BY-SA 4.0** ([`wiki/LICENSE.md`](wiki/LICENSE.md)) |
| `gscode-api` rows in `kb.db` | **GPLv3** ([`LICENSES/GPL-3.0.txt`](LICENSES/GPL-3.0.txt)) |
| `source_scripts` / `docs-bo3` / `video-youtube` rows in `kb.db` | **fair-use mirror** ([`LICENSES/fair-use-notice.md`](LICENSES/fair-use-notice.md)) |
| `source_workspace` rows in `kb.db` | **per-modder attribution** ([`LICENSES/modder-attribution.md`](LICENSES/modder-attribution.md)) |
| `personal-notes` / `tools-dtzxporter` rows in `kb.db` | per-row license preserved in metadata |

**DMCA-ready**: if any rights-holder objects to as-is inclusion of their
content, the affected `kb.db` rows can be flipped to pointer-only
(URL/path + metadata, no body) in one commit. See
[`NOTICE.md`](NOTICE.md#how-to-use-this-db-legally) for the full posture.

## Upstream

- **Knowledge pipeline**: [t7-reapy/t7_knowledge](https://github.com/t7-reapy/t7_knowledge)
  — scraping, distillation, tagging, kb.db build.
- **Architecture plan**: [`.docs/architecture/t7_companion.md`](https://github.com/t7-reapy/t7_knowledge/blob/main/.docs/architecture/t7_companion.md)
  in t7_knowledge.

## Contributing / reporting issues

Wiki edits, license-posture concerns, attribution corrections, and
content-removal requests all welcome — open an issue or contact the
maintainer [McReaper](https://github.com/McReaper) directly.
