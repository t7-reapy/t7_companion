# Zoning, Portals, and Umbra

> Stub — content pending.

Visibility culling: the systems that decide what the engine actually has to draw each frame. Three layered concepts in BO3:

1. **Map zones** — the script-time partition of the level (`zone_unlocked` flags, `zone_volume` brushes for spawning).
2. **Portals** — Radiant-time hint surfaces that split the BSP into rooms the renderer can cull against.
3. **Umbra** — the runtime occlusion-culling solver that uses the BSP + portal data to decide visibility per frame.

These are easy to confuse because the words overlap with neighbouring concepts (see asset-pipeline zone files, which are unrelated).

## TODO

### Map zones (script-time)
- The `zone_volume` entity and its targetname conventions
- `zone_unlocked` / `zone_active` flags and when zones flip
- Spawning constraints — zombies only spawn inside active zones
- The "no powerup outside active zone" gotcha that bit `zm_test` (commit reference once we re-grep)

### Portals (compile-time)
- Portal hint brushes — what surface counts as a portal
- The "one portal per doorway" guideline and when to break it
- How portal placement affects light bake / GI bleed

### Umbra (runtime)
- What umbra actually solves — BSP leaf visibility per frame
- The Compile step generates the umbra data; can fail silently
- Reading the compile output for umbra warnings
- The one unresolvable umbra issue in `zm_test` (per `PRESENTATION_PLAN.md`)
- How to debug "thing not rendering" issues — usually umbra over-culling

## Common gotchas

- Zone names that don't match between `.map` and `.gsc` — silent failure, zone never activates
- Portal brushes that don't fully seal a room — defeats the cull
- Light leaks across portals — usually means the portal is on a thin wall

## Cross-references

- Asset pipeline: [`02-zone-files.md`](../asset-pipeline/02-zone-files.md) — *unrelated* despite the name overlap (those are linker zones)
- [`brushes-patches.md`](./brushes-patches.md) — portal hint surfaces are brushes
- Lighting: [`light-bake.md`](../lighting/light-bake.md) — the bake follows portal boundaries
- Scripting: [`01-overview.md`](../scripting/01-overview.md) — `zone_unlocked` flag use
