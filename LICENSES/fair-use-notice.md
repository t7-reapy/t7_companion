# Fair-use notice for bundled third-party content

This notice covers content bundled inside `t7kb.db` from the following sources, none of which `t7_companion` claims to own or sublicense:

- **`source_scripts`** — decompiled BO3 GSC/CSC engine scripts, mirrored from the publicly-available [shiversoftdev/t7-source](https://github.com/shiversoftdev/t7-source) GitHub repository (Treyarch IP).
- **`source_dump`** — a broad BO3 game-data dump (script bundles, vision sets, gamedata tables, Radiant key definitions), mirrored from the publicly-available [ate47/bo3-source](https://github.com/ate47/bo3-source) GitHub repository (Treyarch IP, dumped with Atian CoD Tools).
- **`docs-bo3`** — Treyarch's official BO3 mod tools reference documentation (PDF/HTM), distributed openly by Treyarch as part of the BO3 mod tools install (`docs_modtools/`).
- **`video-youtube`** — per-segment transcripts (~30 sec chunks) of publicly available YouTube tutorials, each linked via `?t=<seconds>` deep-link URLs to the original video.

## Posture

These entries are included **as-is, with attribution**, under a fair-use mirror posture:

- The content is already broadly distributed by its respective rights holders (Treyarch via the free mod tools download; YouTube videos publicly visible with creator-uploaded captions; shiversoftdev's repository public on GitHub since 2021 without takedown).
- Inclusion here is for **research, education, preservation, and reference** purposes, with full attribution preserved on every entry.
- We do not relicense, claim ownership of, or grant any rights to these works that we do not possess.
- Every `t7kb.db` row carries an `url` (or local-install path) and a `source`-prefixed `doc_id` so users can verify provenance.

## DMCA-ready posture

If Treyarch / Activision, any YouTube creator, or any other rights-holder requests removal or change of inclusion, the affected rows in `t7kb.db` can be flipped to **pointer-only** (URL/path + structured metadata only, no body text) or removed entirely — without affecting unrelated rows. The exact request flow is documented in [`../docs/add-remove-knowledge.md`](../docs/add-remove-knowledge.md). `doc_id`s and search behavior stay stable; only the row's `body` column is affected. For the data model itself, see [`../docs/data-model.md`](../docs/data-model.md).

## Attribution lines on every entry

When `t7kb.db` is rendered by the consumer, entries from these sources should display:

- `source_scripts` rows: *"Decompiled BO3 script from [shiversoftdev/t7-source](https://github.com/shiversoftdev/t7-source) @ commit `<sha>` — Treyarch IP, included under fair-use reference posture."*
- `source_dump` rows: *"BO3 game data from [ate47/bo3-source](https://github.com/ate47/bo3-source) @ commit `<sha>` — Treyarch IP, included under fair-use reference posture."*
- `docs-bo3` rows: *"From Treyarch's official BO3 mod tools reference documentation (shipped openly with the mod tools)."*
- `video-youtube` rows: *"Transcript segment from `<creator>` — YouTube `<url>?t=<seconds>` — fair-use deep-link."*

## Contact

To request removal or repositioning of any content covered by this notice, contact the maintainer [McReaper](https://github.com/McReaper) directly (private channel preferred for takedowns).
