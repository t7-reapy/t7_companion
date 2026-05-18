# Light Bake

> Stub — content pending.

What the **Light** step in the modtools launcher actually does, end to end. The longest single phase of a full build, and the one most likely to surface latent geometry / portal / volume issues.

## TODO

### What runs during the bake
- Direct lighting solve per surface
- GI bounce passes (count + quality settings)
- Light grid sample placement and bake
- Probe baking (reflection probes are baked here too — see [`probes.md`](./probes.md))
- Lightmap UV unwrap (or is that done at Compile? — verify)
- Output: lightmap textures, light grid data, probe cubemaps

### Inputs that drive the bake
- Light entities (point, spot, sun, area — verify exact types BO3 supports)
- `volume_lightclip` — bounds the bake (don't bake the void)
- Surface flags — which materials accept GI, which don't
- Sun + skybox — real or baked sky contribution

### Quality vs time
- Quality presets in the launcher (verify levels)
- Preview bakes — what shortcuts they take
- The "ship-quality" final bake — when to run it (typically once, very late)

### Reading the bake output
- Where the launcher's bake log lives
- Common warnings: light leak, unbounded bake, surface without UV, etc.
- VRAM impact of higher-quality lightmaps

### Common gotchas
- Light leak through brush slivers — pre-bake fix is geometry, not lighting
- Probes baked into walls — they sample inside-the-wall, then dynamic objects look wrong nearby
- GI bleed across portals — sometimes desired (continuity), sometimes not
- Sun direction changing → re-bake required
- Adding a `volume_lightclip` after a bake → must re-bake to take effect

## Cross-references

- [`01-overview.md`](./01-overview.md) — where the bake fits among lighting layers
- [`probes.md`](./probes.md) — probe baking is part of this step
- [`mapping/zoning-portals-umbra.md`](../mapping/zoning-portals-umbra.md) — portals constrain the GI solver
- Asset pipeline: [`01-overview.md`](../asset-pipeline/01-overview.md) — where Light fits in the four-step build
