# Probes — Light Grid & Reflection Probes

> Stub — content pending.

The two probe systems in BO3:

1. **Light grid** — a volumetric grid of ambient samples used for *dynamic objects* (props, characters, anything not part of the baked BSP).
2. **Reflection probes** — cubemap captures of the scene used for *shiny surfaces* (glossy materials, screens, water).

They're both authored at Radiant time, baked during the Light step, sampled at runtime.

## TODO

### Light grid
- The grid is automatic — placed by the bake based on map bounds and density settings (verify exactly how density is controlled)
- Per-cell ambient sampling: direction + colour
- Dynamic objects sample the cell they're in
- Visualization: dvar to show the grid (verify exact name)
- Common gotcha: sample inside a wall → dynamic prop near the wall picks up wrong colour

### Reflection probes
- `reflection_probe` entity placement: in the centre of the volume the probe should serve
- `volume_reflection` entity defines the *region* a given probe applies to
- Multiple probes per region — engine picks based on distance? — verify
- Cubemap resolution / quality — perf vs visual
- Re-bake required after moving probes or volumes
- Common gotcha: probe placed too close to a wall → reflections pick up the wall geometry

### Manual probe placement
- The "manual probe placement issues" mentioned in `PRESENTATION_PLAN.md` — surface the specifics when filling this page
- The "single probe orbe" experiment in `zm_test`'s last zone — what worked, what didn't (commit history)

## Common gotchas

- Probe inside-wall syndrome — the most common visual bug
- Reflection probes don't capture dynamic objects — only the baked scene
- Probe density too high → bake time explodes, runtime cost rises, visual gain is small past a point
- Light grid sampling on a fast-moving object can pop between cells — soften with grid density, not script

## Cross-references

- [`01-overview.md`](./01-overview.md) — where probes fit among lighting layers
- [`light-bake.md`](./light-bake.md) — probes are baked during this step
- [`mapping/volume-entities.md`](../mapping/volume-entities.md) — `volume_reflection` and others
