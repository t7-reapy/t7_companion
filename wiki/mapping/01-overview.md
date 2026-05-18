# Mapping — Overview

> Stub — content pending.

The big-picture page for the mapping discipline. What Radiant is, how a `.map` file becomes BSP geometry, where the boundaries with asset-pipeline / scripting / lighting fall.

## TODO

- The Radiant editor in 60 seconds: 2D/3D viewports, brush vs entity, the entity browser
- Files that make up a map: `.map` (Radiant source), `_prefabs/`, the compile artefacts
- The Compile step in the modtools launcher — what it actually does (BSP solve, lighting prep, etc.)
- BSP fundamentals: convex hulls, splitting, leaves, portals (link to `zoning-portals-umbra.md`)
- The Radiant ↔ scripting handshake: `targetname`, `script_noteworthy`, `script_string`, `script_int`, `script_float` and how scripts read them via `GetEnt` / `GetEntArray` / `getentarray("targetname")` (link to `radiant-kvps-to-scripts.md`)
- Where to start when opening a fresh `.map` (template? `zm_template`?)
- How to share / iterate on a map across machines (Radiant doesn't merge cleanly)
- Common pitfalls in the very first hours

## Cross-references

- [`brushes-patches.md`](./brushes-patches.md) — the geometry primitives
- [`brush-textures.md`](./brush-textures.md) — caulk / nodraw / clip / hint
- [`zoning-portals-umbra.md`](./zoning-portals-umbra.md) — visibility
- [`ai-pathing.md`](./ai-pathing.md) — path nodes + traversals
- [`volume-entities.md`](./volume-entities.md) — the volume catalogue
- [`prefabs.md`](./prefabs.md) — modular assembly
- Asset pipeline: [`01-overview.md`](../asset-pipeline/01-overview.md) for the build flow Compile → Light → Link → Run
