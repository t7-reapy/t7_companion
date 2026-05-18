# Mapping

> **Status: stub.** Top-level discipline. Radiant-side work: BSP geometry, brushes, patches, prefabs, zones, AI pathing, portals, umbra. The third-largest section after asset-pipeline and scripting.

## What this section will cover

- Radiant editor basics: brushes, patches, entities, prefabs.
- Brush construction: convex requirement, subtraction, carving.
- Patches: curved/planar 3D surfaces, bezier control points.
- Tool-textures (caulk, nodraw, hint, skip, the various `clip_*` flags) — what each does at compile time.
- Prefabs: how `_prefabs/zm/zm_test/house/*.map` instances assemble into the parent map.
- Zoning: portals, hint brushes, umbra targeting.
- AI navigation: path nodes, traversal nodes, cover nodes, navmesh.
- Volume entities catalogue (`volume_sun`, `volume_worldfog`, `volume_litfog`, `volume_lut`, `volume_reflection`, `volume_lightclip`, `volume_outdoor`, `volume_exposure`, `volume_weathergrime`, `nav_volume`, etc.).
- Wallbuy / barrier / spawner placement patterns.
- Compile reports: `.badnodes`, `traversals`, `map_stats` (the on-demand reports under **Misc → Generate Reports**).
- Getting Radiant KVPs visible to scripts (the registration-step gap Reapy mentioned).
- Common gotchas: brush slivers, light leaks, umbra issues, pathing dead ends, prefab origin drift.

## Sub-pages (planned)

- `01-overview.md` — Radiant + BSP fundamentals (TODO)
- `brushes-patches.md` — geometry primitives (TODO)
- `brush-textures.md` — tool textures + clip flags
- `prefabs.md` — modular map assembly (TODO)
- `zoning-portals-umbra.md` — visibility (TODO)
- `ai-pathing.md` — nodes + navmesh (TODO)
- `volume-entities.md` — catalogue (TODO)
- `radiant-kvps-to-scripts.md` — KVP exposure to GSC

## Reference reading (general)

- [BO3 Mod Tools playlist tuto (YouTube)](https://www.youtube.com/watch?v=tj7guP_ZeUI&list=PLnt_Nobu89HtwqIkEtRt4zj6RC_s-_f2q) — long-form Radiant walkthrough.
- [Modme Wiki — BO3 mapping pages](https://wiki.modme.co/wiki/Game-Support-_-Black-Ops-3.html).
- [Advanced Terrain Series (YouTube)](https://www.youtube.com/playlist?list=PLoSUR65Ja21rvlDtjcBZMu3KCMike0H7C) — patch-heavy terrain workflow.

## Status

Reapy spent ~6 months in this discipline (the pre-git phase of `Rainy Doom`). Lived-experience density is high but with the caveat that the first 6 months of work has no commit history. Lighting + AI pathing are likely the densest sub-areas.

**Expected to be a huge section** — the Radiant editor alone has enough surface area to warrant deep coverage. Possible content sources include transcribing community YouTube tutorials (Advanced Terrain Series, the BO3 Mod Tools playlist, etc.) into structured written form alongside Reapy's lived experience.
