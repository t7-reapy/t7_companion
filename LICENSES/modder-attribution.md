# Per-modder attribution for `source_workspace` content

The `source_workspace` source bundled in `t7kb.db` mirrors the maintainer's own BO3 modding workspace, publicly available at [t7-reapy/t7_workspace](https://github.com/t7-reapy/t7_workspace).

## Three content buckets

| Bucket | Path patterns | Posture |
|---|---|---|
| **Maintainer's own work** | `_reapy/`, top-level `KNOWLEDGE.md` / `README.md` | Original content authored by the maintainer (Reapy). Included under the same posture as any other content the maintainer chose to publish in their workspace repo. |
| **Maintainer's usermaps** | `usermaps/` | The maintainer's own custom maps — scripts, zone configs, and project structure are authored by the maintainer, but the maps include imports from community asset packs (HarryBo21 guns, Kingslayerkyle HUDs, MADGAZ texture packs, etc.). Attribution for the imported packs is encoded by directory under the bucket below + cross-referenced in the workspace `README.md` credits table. |
| **Edits to shipped Treyarch content** | workspace default | Same posture as `source_scripts` — Treyarch IP mirrored under the public-on-GitHub fair-use stance; see [`fair-use-notice.md`](fair-use-notice.md). |
| **Third-party modder content** | `_custom/<modder>/`, `model_export/<modder>/`, `source_data/<modder>/`, `sound_assets/<modder>/` | Community modder packs the maintainer installed. **Attribution encoded by directory name**, plus a cross-referenced credits table in the workspace's own `README.md` (mirrored at `source_workspace::workspace.README` in `t7kb.db`) listing **author + original download links** for every external pack. |

The workspace `README.md` (upstream: [t7-reapy/t7_workspace](https://github.com/t7-reapy/t7_workspace/blob/main/README.md)) is the **canonical credits table** — it lists every modder pack present, the author, and the original distribution URLs. That sheet travels with the data into `t7kb.db`.

## Community convention

The BO3 modding scene shares modder packs under an informal **"credit me if you use this"** convention — typically equivalent to CC-BY without share-alike. Each pack the maintainer installed was chosen specifically because the pack's author distributes it openly for community reuse (the workspace's credits table includes the original distribution URLs for verification).

## DMCA-ready, per-modder removal

If any specific modder objects to inclusion of their content, the affected subtree (`_custom/<modder>/`, `model_export/<modder>/`, etc.) can be flipped to pointer-only (path + metadata, no body) or removed entirely — without affecting other modders' content. The exact request flow is documented in [`../docs/add-remove-knowledge.md`](../docs/add-remove-knowledge.md). Per-modder removals are surgical: `doc_id`s and search behavior stay stable, and unrelated rows are untouched.

## Contact

To request removal or repositioning of any content under `source_workspace`, contact the maintainer [McReaper](https://github.com/McReaper) directly (private channel preferred for takedowns).
