# Lighting — Overview

> Stub — content pending.

The big picture for the lighting discipline. What the **Light** step does in the modtools launcher, where lighting state lives at runtime, the boundary between baked and dynamic.

## TODO

### The Light step
- What runs: GI bake, light grid sample placement, lightmap generation, probe baking
- Cost: minutes-to-hours depending on map size and quality
- Inputs from Radiant: lights, sun, lightclip volumes, surfaces tagged for lighting
- Outputs: BSP lightmaps, light grid, baked probes
- The "iteration loop" reality — full bakes are slow; preview bakes are common

### Layers of lighting in BO3
- **Baked direct/indirect**: surface lightmaps, fixed at compile
- **Light grid**: volumetric ambient samples for dynamic objects (props, characters)
- **Reflection probes**: cubemap captures for shiny surfaces
- **Sun + sky**: real-time directional + ambient-from-skybox
- **Volumetric (lit) fog**: traced through the world fog volume
- **Dynamic lights**: real-time lights script can toggle
- **Visionsets / LUTs**: post-process colour grading
- **Exposure / tonemapping**: per-region adjustment

### Runtime lighting state
- The `set_lighting_state(N)` mechanism — how a single call swaps SSI + light toggles + fog gating in coordination
- See [`asset-pipeline/05-ssi.md`](../asset-pipeline/05-ssi.md) for the SSI side and the late-joiner gotcha

### Sub-pages
- [`light-bake.md`](./light-bake.md) — the bake step in detail
- [`probes.md`](./probes.md) — light grid + reflection probes
- [`visionsets-and-luts.md`](./visionsets-and-luts.md) — post-processing
- [`skybox.md`](./skybox.md) — skybox model + animation
- [`weather.md`](./weather.md) — rain, lightning, weather state coupling

## Cross-references

- [`asset-pipeline/05-ssi.md`](../asset-pipeline/05-ssi.md) — SSI authoring and the lighting-state mechanism
- [`asset-pipeline/04-fog.md`](../asset-pipeline/04-fog.md) — fog banks gated by lighting state
- [`mapping/zoning-portals-umbra.md`](../mapping/zoning-portals-umbra.md) — portals constrain the bake
