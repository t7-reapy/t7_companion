# Per-modder attribution for `source_workspace` content

The `source_workspace` source bundled in `kb.db` mirrors the maintainer's
own BO3 modding workspace, publicly available at
[t7-reapy/t7_workspace](https://github.com/t7-reapy/t7_workspace).

## Three content buckets, three postures

### Maintainer's own work

Paths under `usermaps/`, `_reapy/`, or top-level `KNOWLEDGE.md` /
`README.md`. Original IP of the maintainer (Reapy). Included under
the same posture as any other content the maintainer chose to publish
in their workspace repo.

### Edits to shipped Treyarch content

Workspace files derived from Treyarch's shipped mod tools content
(custom modifications to engine GSC, zone configs, etc.). Same posture
as `source_scripts` — Treyarch IP mirrored under the public-on-GitHub
fair-use stance. See [`fair-use-notice.md`](fair-use-notice.md).

### Third-party modder content

Paths under `_custom/<modder>/`, `model_export/<modder>/`,
`source_data/<modder>/`, `sound_assets/<modder>/`. These are community
modder packs the maintainer installed in their workspace.

**Attribution is encoded by directory name**, and the workspace's own
committed `README.md` (mirrored at `source_workspace::workspace.README`
in `kb.db`) carries a cross-referenced credits table with **author +
original download links** for every external pack installed.

Known authors whose packs are present in the workspace include — but
are not limited to: **HarryBo21, Kingslayerkyle, Midgetblaster,
MoiCestTOM, MADGAZ, Symbo, Zeroy, Slick Willy, Scobalula, LG-RZ,
DTZxPorter (D3V Team)**, and others; see the workspace `README.md` row
in `kb.db` (or
[the upstream](https://github.com/t7-reapy/t7_workspace/blob/main/README.md))
for the full credits table.

## Community convention

The BO3 modding scene shares modder packs under an informal
**"credit me if you use this"** convention — typically equivalent to
CC-BY without share-alike. Each pack the maintainer installed was
chosen specifically because the pack's author distributes it openly
for community reuse (the workspace's credits table includes the
original distribution URLs for verification).

## DMCA-ready, per-modder removal

If any specific modder objects to inclusion of their content:

1. The affected `_custom/<modder>/`, `model_export/<modder>/` (etc.)
   subtree gets flipped to pointer-only (path + metadata, no body) in
   one commit, leaving other modders' content unaffected.
2. Full removal of a specific modder's subtree is equally one commit.
3. `doc_id`s and search indexing remain stable across these changes.

## Contact

To request removal or repositioning of any content under
`source_workspace`, open an issue against this repository or contact
the maintainer [McReaper](https://github.com/McReaper) directly.
Per-modder removals are surgical — they don't affect anyone else's
content.
