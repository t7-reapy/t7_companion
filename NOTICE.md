# NOTICE — Polyglot Licensing

`t7_companion` is maintained by [McReaper](https://github.com/McReaper) under the [T7 Modding Area](https://github.com/t7-reapy) organization.

This repository bundles content from multiple sources, each under its own license or use posture. The repository as a whole is an *aggregate* in the GPLv3 sense — bundled works retain their original licenses, and that inclusion does not relicense them.

| Component | Path | License / posture |
|---|---|---|
| **t7_companion code** | `src/`, top-level scripts | **MIT** — see [`LICENSE`](LICENSE) |
| **Synthesized wiki pages** | `wiki/` | **CC-BY-SA 4.0** — see [`wiki/LICENSE.md`](wiki/LICENSE.md) and [`LICENSES/CC-BY-SA-4.0.txt`](LICENSES/CC-BY-SA-4.0.txt) |
| **Per-domain skill files** | `skills/` | **CC-BY-SA 4.0** — same license as `wiki/` |
| **Contributor docs** | `docs/` | **CC-BY-SA 4.0** — same license as `wiki/` |
| **`gscode-api` rows** in `t7kb.db` | `source = 'gscode-api'` | **GPLv3** (upstream [Blakintosh/gscode](https://github.com/Blakintosh/gscode)) — see [`LICENSES/GPL-3.0.txt`](LICENSES/GPL-3.0.txt) |
| **`source_scripts` rows** in `t7kb.db` | `source = 'source_scripts'` | **Treyarch IP, fair-use mirror** of public upstream [shiversoftdev/t7-source](https://github.com/shiversoftdev/t7-source) — see [`LICENSES/fair-use-notice.md`](LICENSES/fair-use-notice.md) |
| **`docs-bo3` rows** in `t7kb.db` | `source = 'docs-bo3'` | **Treyarch IP, shipped openly with the BO3 mod tools** — same fair-use posture, see [`LICENSES/fair-use-notice.md`](LICENSES/fair-use-notice.md) |
| **`source_workspace` rows** in `t7kb.db` | `source = 'source_workspace'` | **Maintainer IP + per-modder attribution** under BO3 community convention — see [`LICENSES/modder-attribution.md`](LICENSES/modder-attribution.md) |
| **`video-youtube` rows** in `t7kb.db` | `source = 'video-youtube'` | **Per-creator copyright, fair-use deep-linking** — see [`LICENSES/fair-use-notice.md`](LICENSES/fair-use-notice.md) |
| **`personal-notes` rows** in `t7kb.db` | `source = 'personal-notes'` | **CC-BY 4.0** — the maintainer owns this content and licenses it for attribution-only redistribution |
| **`tools-dtzxporter` rows** in `t7kb.db` | `source = 'tools-dtzxporter'` | **Per-repo upstream license** — recorded in each row's `license` metadata field |

> [!NOTE]
> **Where does `t7kb.db` actually live?** It ships as a **GitHub Release asset** (downloaded by the consumer at install / update time), not committed verbatim in the repo — file size makes that impractical. The schema is documented at [`docs/data-model.md`](docs/data-model.md).

## How to use this DB

Most users will simply run the consumer alongside an LLM of their choice (Claude Code, OpenCode, any tool-use-capable client) — install the consumer, the DB downloads, and queries go via the search/get tools. The DB is queryable for any personal modding-research use; no extra ceremony needed.

The legal terms below only matter if you intend to **redistribute** the DB or build derivatives of it.

## Redistribution — what you need to preserve

Any redistribution of this repo, of `t7kb.db`, or of derivative wikis must preserve:

1. This `NOTICE.md` file (or equivalent).
2. The full text of every bundled license (the `LICENSES/` directory).
3. Per-source attribution as already encoded in each `t7kb.db` row's `source`, `url`, `license`, and (where applicable) author metadata.
4. The `wiki/` directory's CC-BY-SA share-alike requirement on any downstream wiki derivatives.

> [!IMPORTANT]
> The wiki, skills, and docs are CC-BY-SA 4.0 — if you remix or build upon them, your derivative must also be CC-BY-SA 4.0. This protects against closed-source forks of community-derived knowledge.

## Reporting an issue

Open an issue on this repository, or contact the maintainer [McReaper](https://github.com/McReaper) directly. The exact flow for content removal is documented at [`docs/add-remove-knowledge.md`](docs/add-remove-knowledge.md).
