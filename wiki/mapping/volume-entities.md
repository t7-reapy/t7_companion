# Volume Entities — Catalogue

> Stub — content pending.

BO3 places a lot of behavior on **volume** entities — invisible brush-textured boxes that mark regions where some engine subsystem should activate. This page is the canonical catalogue.

## TODO — Volumes to document

| Entity | What it does | Documented on page |
|---|---|---|
| `volume_sun` | Selects which `lightingstate` SSI slot is active in the region | Lighting / SSI |
| `volume_worldfog` | Defines a fog bank; gated by `LIGHTSTATE_N` flags | [`asset-pipeline/04-fog.md`](../asset-pipeline/04-fog.md) |
| `volume_litfog` | Lit (volumetric) fog variant | Lighting |
| `volume_lut` | LUT colour-grading region | [`lighting/visionsets-and-luts.md`](../lighting/visionsets-and-luts.md) |
| `volume_reflection` | Reflection probe sampling region | [`lighting/probes.md`](../lighting/probes.md) |
| `volume_lightclip` | Bounds the light bake | [`lighting/light-bake.md`](../lighting/light-bake.md) |
| `volume_outdoor` | Marks open-sky regions (rain/weather affect player) | [`lighting/weather.md`](../lighting/weather.md) |
| `volume_exposure` | Camera exposure adjustment region | Lighting |
| `volume_weathergrime` | Weather wear effect on materials | [`lighting/weather.md`](../lighting/weather.md) |
| `nav_volume` | AI navigation bounding | [`ai-pathing.md`](./ai-pathing.md) |
| `zone_volume` | Map zone partition | [`zoning-portals-umbra.md`](./zoning-portals-umbra.md) |
| `trigger_*` (various) | Interactivity triggers | Scripting recipes |

This page is the index — each entity gets a short "what it is, what brush-texture to use, KVPs that matter, where it's authoritatively documented" entry. Don't duplicate the deep coverage; link out.

## Conventions

- All volumes are **brushes** with a special tool-texture (usually `clip_volumes/<name>` or similar — verify the exact path).
- Most volumes are *additive*: overlap means both apply. A few are *exclusive*; flag those when documenting.
- Volume KVPs override the world default for the region they cover.

## Cross-references

- [`brush-textures.md`](./brush-textures.md) — the tool-texture side of the volume tag
- Each linked page above is the canonical home for its volume's behavior
