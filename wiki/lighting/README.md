# Lighting

> **Status: stub.** Top-level discipline alongside asset-pipeline / scripting / mapping. Tightly coupled to mapping (and to the SSI / fog pages in asset-pipeline) but big enough to stand alone.

## What this section will cover

- The Light step in the modtools launcher: what it actually does (GI bake, light grid, probe placement, lightmap generation).
- Light entities: types, falloff, shadow casting, baked vs dynamic.
- The light grid (volumetric ambient sampling for dynamic objects).
- Reflection probes (`reflection_probe`, `volume_reflection`) — how the engine samples the world and re-uses it for shiny surfaces.
- The lighting-state mechanism end-to-end (`util::set_lighting_state`, `lightingstate_N` KVPs on lights, `LIGHTSTATE_N` flags on fog volumes, `ssiN` slots on `volume_sun`).
- Volumetric lighting / lit fog perf considerations.
- Sun cookies (`lightdescription` assets, the `cookie_sun_*` records).
- Visionsets, LUTs, colour grading (separate sub-page).
- Weather effects coupling (rain droplets, lightning strikes — separate sub-page).
- Skybox setup (model selection, animated/rotating skybox — separate sub-page).
- Common gotchas: light leaks, umbra issues, probe misplacement, baked vs dynamic mismatches.

## Sub-pages (planned)

- `01-overview.md` — the big picture (TODO)
- `light-bake.md` — Light-step internals (TODO)
- `probes.md` — light grid + reflection probes (TODO)
- `visionsets-and-luts.md` — post-FX colour grading
- `skybox.md` — skybox model + animation
- `weather.md` — rain droplets, thunder/lightning, weather state coupling

## Reference reading (general)

- [Lighting Tutorial (YouTube playlist)](https://www.youtube.com/watch?v=hPweAEu8zJY&list=PLYAF4YwatlpU47HOddWaEPJlQt3yjoNtY) — the canonical multi-part walkthrough; closest to a "lighting from scratch" series.
- [`docs_modtools/`](file:///D:/SteamLibrary/steamapps/common/Call%20of%20Duty%20Black%20Ops%20III/docs_modtools/) — Treyarch's own PDFs include lighting + probes references.

## Status

Reapy is hands-on with the **lighting-state mechanism** (sun/sky/fog coordination) and tuned several specifics for `zm_test` (storm penumbras, exposure ranges). Light bake / probe placement / GI is more of an "I made it work" area — research-driven content welcome.
