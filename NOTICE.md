# NOTICE — Polyglot Licensing

`t7_companion` is maintained by [McReaper](https://github.com/McReaper)
under the [T7 Modding Area](https://github.com/t7-reapy) organization.

This repository bundles content from multiple sources, each under its own
license or use posture. The repository as a whole is an *aggregate* in the
GPLv3 sense — bundled works retain their original licenses, and that
inclusion does not relicense them.

| Component | Path | License / posture |
|---|---|---|
| **t7_companion code** (CLI / MCP server / scripts) | `src/`, top-level scripts | **MIT** — see [`LICENSE`](LICENSE) |
| **Synthesized wiki pages** | `wiki/` | **CC-BY-SA 4.0** — see [`wiki/LICENSE.md`](wiki/LICENSE.md) and [`LICENSES/CC-BY-SA-4.0.txt`](LICENSES/CC-BY-SA-4.0.txt) |
| **`gscode-api` entries** in `kb.db` | rows where `source = 'gscode-api'` | **GPLv3** (upstream [Blakintosh/gscode](https://github.com/Blakintosh/gscode)) — see [`LICENSES/GPL-3.0.txt`](LICENSES/GPL-3.0.txt) |
| **`source_scripts` entries** in `kb.db` (decompiled BO3 engine GSC/CSC) | rows where `source = 'source_scripts'` | **Treyarch IP, fair-use mirror** of public upstream [shiversoftdev/t7-source](https://github.com/shiversoftdev/t7-source) — see [`LICENSES/fair-use-notice.md`](LICENSES/fair-use-notice.md) |
| **`docs-bo3` entries** in `kb.db` (Treyarch mod tools reference docs) | rows where `source = 'docs-bo3'` | **Treyarch IP, shipped openly with the BO3 mod tools** — same fair-use posture, see [`LICENSES/fair-use-notice.md`](LICENSES/fair-use-notice.md) |
| **`source_workspace` entries** in `kb.db` (maintainer's modding workspace + community modder packs) | rows where `source = 'source_workspace'` | **Maintainer IP + per-modder attribution** under BO3 community "credit me if you use this" convention — see [`LICENSES/modder-attribution.md`](LICENSES/modder-attribution.md) |
| **`video-youtube` entries** in `kb.db` (transcript segments with `?t=N` deep links) | rows where `source = 'video-youtube'` | **Per-creator copyright, fair-use deep-linking** — see [`LICENSES/fair-use-notice.md`](LICENSES/fair-use-notice.md) |
| **`personal-notes` entries** in `kb.db` (maintainer's firsthand BO3 modding notes) | rows where `source = 'personal-notes'` | **CC-BY 4.0** (maintainer's choice) — attribution-only |
| **`tools-dtzxporter` entries** in `kb.db` (public GitHub READMEs) | rows where `source = 'tools-dtzxporter'` | **Per-repo license** (typically MIT) — preserved in each entry's metadata |

## How to use this DB legally

- **Code in `src/` is MIT** — fork it, integrate it into your own tools.
- **Wiki pages are CC-BY-SA 4.0** — republish freely with attribution and
  share-alike.
- **`kb.db` is an aggregate** — querying it is fine; redistributing it
  whole means redistributing all the bundled licenses. Each row carries
  its `source` so consumers can filter by license posture.
- **DMCA-ready**: if any rights-holder objects to as-is inclusion of
  their content, the affected entries can be flipped to pointer-only
  (URL/path + metadata, no body) in one commit. `doc_id`s stay stable.

## Attribution requirements when redistributing

Any redistribution of this repo, of `kb.db`, or of derivative wikis must
preserve:

1. This `NOTICE.md` file (or equivalent).
2. The full text of every bundled license (the `LICENSES/` directory).
3. Per-source attribution as already encoded in each kb.db row's
   `source`, `url`, and (where applicable) author metadata.
4. The full `wiki/` directory's CC-BY-SA share-alike requirement on
   downstream wiki derivatives.

## Reporting an issue

If you are a rights-holder of any bundled content and want it removed or
its disposition changed, open an issue against this repository or contact
the maintainer directly. Pointer-only conversion happens in one commit;
full removal is also straightforward — every kb.db row tracks its source.
