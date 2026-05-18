# Images & Materials

> The most-used GDT asset types, and the place where most "why doesn't it look right?" debugging happens.

## Mental model

Two layered concepts. Don't conflate them.

- **`image`** = a *bitmap on disk* plus the engine's settings for how to compress, mipmap, sample, and stream it. One image record = one rendering-ready texture asset addressable by name.
- **`material`** = a *shader instance*. It picks a shader template (the `materialType`), binds 1+ images to its slots (`colorMap`, `normalMap`, `specMap`, etc.), and sets per-instance parameters (color tints, gloss, scrolling, alpha behavior, surface tags). One material record = one thing the renderer can apply to a brush face or a model surface.

The renderer never sees images directly. It sees materials, and follows them to their bound images. So:

```
Radiant brush face   ──►   material   ──►   image(s)   ──►   raw bitmap on disk
script SetMaterial   ──►   material   ──►   image(s)   ──►   raw bitmap on disk
```

When you "make a wall texture," you're making **at minimum a color image and one material**, and probably also a normal-map image (and maybe a spec/gloss image, an emissive image, etc.) bound into that one material.

## Where bitmaps come from

Realistic mix in this project (and most community BO3 mods):

| Origin                                                                       | Frequency in `zm_test` |
| ---------------------------------------------------------------------------- | ---------------------- |
| **Pre-extracted community packs** (Midgetblaster, harrybo21, SkyeLord ports…)| Dominant. Pack ships GDTs + raw assets; you install + occasionally tweak. |
| **Greyhound / Saluki extraction by you**                                     | Rare — Reapy ran it for ~2 missing assets. Most porting work was already done by pack authors. |
| **Authored** (Photoshop, NormalMap-Online, Substance…)                       | Small but non-zero — custom text textures, powerup originals, the **rain viewmodel droplets** (sprite flipbook, hand-authored frame-by-frame), the room-of-thanks `thx_` credit textures, and a few one-offs. |

> *Bitmaps aren't the only thing that gets ported.* The same install-or-extract loop applies to xmodels (.xmodel_export), animations (.xanim_export), sounds, and other raw asset types. This page is image-and-material focused, but the workflow generalizes — see `03-models.md` for the model-specific variant.

### Installing a community pack (the dominant path)

Most "porting" work in this project is really *consuming* someone else's port:

1. Drop the pack's raw assets under the appropriate BO3 search root (`texture_assets/_<author>/...`, `model_export/_<author>/...`, `_custom/...`). Match the directory shape the pack's GDT expects.
2. Drop the pack's GDT(s) into `source_data/` so APE indexes them.
3. **Check paths inside the GDT.** Often the pack assumes a particular search root; if your install layout differs, you'll need to either move files or edit `baseImage` / `baseModel` paths in the GDT records.
4. Reference / zone what you actually want (per the rules in `01-overview.md`).

### Doing your own Greyhound/Saluki extraction (the rare path)

When the asset isn't in any pack and you have to extract from a source game yourself:

1. Run the extractor against a running source-game session (or its on-disk `.xpak` for image data). Output lands in the **tool's** install dir, not in BO3.
2. Copy the raw assets into a BO3 search root.
3. Either adopt the extractor-emitted GDT (fixing paths) or write the GDT entries fresh in APE.
4. Reference / zone.

### Authored / generated bitmaps

For Photoshop, Substance, NormalMap-Online, etc.: you make the bitmap, drop it in `texture_assets/` (or `_custom/`), then create the GDT image record from scratch in APE. The engine doesn't care how it was made — only the GDT entry crosses into runtime.

Common pitfalls across all three flows: GDT path drift after install (linker says "could not find asset"), accidentally importing a GDT whose asset names collide with something already in your registry (silent shadowing — see the dupe-purger note in `01-overview.md`), or copying only some of the chained dependencies (a material with a missing image, an xmodel with a missing material).

## The `image` GDT record

Most fields are tuning knobs. The semantic-meaningful ones cluster into 5 groups.

### 1. The source

- **`baseImage`** — relative path to the raw bitmap (`.png`, `.tga`, `.tif`). Resolved against `bin/converter_gdt_dirs_0.txt` search roots. Example: `texture_assets\\_bo6\\anh_dirt_rocky_pebbles_01_dec_disp\\ximage_456a5dff2cd319a_c.png`.
- **`type`** — always `image`.
- **`imageType`** — `Texture` (almost always — flat 2D bitmap), `Cube` (cubemap, e.g. for reflection probes / skyboxes), or `Volume` (3D volume texture, e.g. for LUTs and effects). The dropdown in APE is labelled `Image Type`.
- **Map** (composite mode toggle in APE; not always exposed as a single GDT field) — switch between a single `image` (one bitmap one record) and a `composite` mode that stitches multiple source images into one output. Composite is how 4-channel packs are authored (e.g. R=AO, G=roughness, B=metallic, A=height).

### 2. Channel semantic — *the most error-prone field*

Two related fields tell the engine *what kind of data* the image carries:

- **`semantic`** — what role the bitmap plays in a shader. Shown as **`Image Usage`** in APE. Full dropdown:

  | Value           | Meaning                                                                     |
  | --------------- | --------------------------------------------------------------------------- |
  | `2d`            | 2D HUD/UI image, no lighting.                                               |
  | `diffuseMap`    | Base color / albedo (sRGB).                                                 |
  | `effectMap`     | FX-system bitmap (sprites, particles).                                      |
  | `normalMap`     | Tangent-space normals (linear).                                             |
  | `specularMask`  | Specular masking (where to apply spec).                                     |
  | `specularMap`   | Specular reflectance (linear).                                              |
  | `glossMap`      | Gloss / roughness (linear).                                                 |
  | `occlusionMap`  | Ambient occlusion bake.                                                     |
  | `revealMap`     | Alpha reveal / dissolve mask.                                               |
  | `normalBody`    | Body-specific normal map variant.                                           |
  | `multipleMask`  | Packed multi-channel mask.                                                  |
  | `thicknessMap`  | Subsurface / thickness data.                                                |
  | `camoMap`       | Camo masking (used by `weaponcamotable`).                                   |
  | `One Channel`   | Single-channel utility image.                                               |
  | `Two Channel`   | Two-channel utility image.                                                  |
  | `Emblem`        | Player emblem texture.                                                      |
  | `HDR`           | High dynamic range source (e.g. environment captures).                      |
  | `Eye Caustic`   | Eye-specific caustic pattern.                                               |
  | `Custom`        | Manually configured semantic — you take responsibility for the pipeline.    |
  | `LutTpage`      | Look-up table tile-page (LUT atlas, used by visionsets/colour grading).     |

- **`coreSemantic`** — encodes two things in one string:
  1. **Channel layout** — how many channels the bitmap carries (1ch, 3ch, 3ch+alpha, 4ch).
  2. **Colorspace** — sRGB (gamma-encoded, perceptual) vs linear (raw values for math).

  Real palette from this project's GDTs (`grep -rh '"coreSemantic"' source_data/ | sort | uniq -c`):

  | Value             | Approx. usage | Meaning / typical use                                                |
  | ----------------- | -------------:| -------------------------------------------------------------------- |
  | `Linear1ch`       | 8,084         | Single linear channel — gloss, AO, masks, packed scalar data.        |
  | `sRGB3chAlpha`    | 5,633         | 3 sRGB color channels + alpha — diffuse/albedo with cutout.          |
  | `Normal`          | 3,974         | Tangent-space normal map (engine handles the channel/colorspace pair internally). |
  | `Linear4ch`       | 2,402         | 4 linear channels — packed multi-data textures (e.g. R=AO, G=rough, B=metal, A=height). |
  | `sRGB3ch`         | 1,706         | 3 sRGB color channels, no alpha — opaque diffuse.                    |
  | `LutTpage`        | 110           | LUT tile-page atlas (visionsets / colour grading).                   |
  | `HDR`             | 80            | High dynamic range source.                                           |
  | `Custom`          | 15            | Manually configured — you take responsibility for the pipeline.      |
  | (empty)           | 365           | Unset — likely defaulted from `semantic`.                            |

  The compiler reads `coreSemantic` to (a) pick the right DDS pixel format on disk and (b) flag gamma handling for the sampler so the GPU does the right conversion at sample time.

  > 💡 In practice you rarely set this by hand — APE derives it from your `semantic` (Image Usage) choice. The advanced `Custom` value is the override hatch.

> ⚠️ **Get this wrong and you get visible bugs**: a normal map flagged as `sRGB` will be gamma-corrected on sample → lighting looks wrong, especially on glancing angles. A diffuse map flagged as linear will look washed out. The "this normal map looks weird" debugging path almost always lands here.

### 3. Compression & mips

- **`compressionMethod`** — effectively a binary toggle: **compressed** (DXT/BC block compression, default for almost everything) or **uncompressed** (raw pixels — only for the rare cases where you can't tolerate compression artifacts, like UI icons with hard edges or LUT atlases). The verbose strings you'll see in GDTs (`compressed high color`, etc.) are internal labels; in practice this is a one-bit decision.
- **`mipMode`** — mip filter (`Average`, …).
- **`mipBase`** — `1/1` is full-res; `1/2`, `1/4`, etc. start the chain at lower res to save memory.
- **`noMipMaps`** — `1` to disable mips entirely (UI/HUD textures usually).
- **`noPicMip`** — `1` to opt out of `r_picmip`-driven downscaling on low-end machines.

### 4. Sampling

- **`clampU`** / **`clampV`** — `0` = repeat (tiling), `1` = clamp at edge (decals, HUD elements).
- **`colorSRGB`** — overrides default colorspace handling. Usually leave at `0` and let `coreSemantic` drive it.

### 5. Streaming

- **`streamable`** — `1` lets the heavy mip levels live in the paired `.xpak` and stream on demand. Default `1` for most world textures.
- **`forceStreaming`** / **`himipStreaming`** — finer control for special cases.

> 💡 The vast majority of fields on an image record are defaulted and you never touch them. The fields that matter for "did I add this asset right?" are the 5 groups above.

## The `material` GDT record

Materials are *much* bigger records (~900 fields in BO3 — most defaulted to `0`). The fields that matter for understanding the asset are a small subset.

### Picking a shader: `materialType` and `materialCategory`

- **`materialType`** — names the shader template the material instantiates. This is the single most important field; it determines what shader code runs and which image slots are even meaningful.
- **`materialCategory`** — a coarser grouping (e.g. `Geometry Advanced`) that APE uses to filter which `materialType` values are offered.
- **`template`** — usually just `material.template`; the GDF file that defines the schema.

There are dozens of `materialType` values. There is no canonical short list — in practice, modders pick by:
1. Looking at the `materialType` of a similar existing asset and copying it.
2. Asking Discord or grepping community packs for working examples.

### Actual `materialType` palette in this project

A `grep -h '"materialType"' source_data/ | sort | uniq -c | sort -rn` across the project's GDTs reveals the *real* palette in use (top entries by frequency):

| `materialType`                          | Approx. usage | Reading the name                                                                |
| --------------------------------------- | -------------:| ------------------------------------------------------------------------------- |
| `lit_micro_tile_blend_advanced`         | 3,069         | Lit material, terrain-style micro-tile blending, full PBR slot set.             |
| `lit_plus_cc`                           | 2,900         | Lit, "plus" tier, with colour-correction (`_cc`) extras.                        |
| `lit_plus`                              | 2,768         | Lit, "plus" tier — common world props.                                          |
| `lit_decal`                             | 2,711         | Lit decal, baseline.                                                            |
| `lit_micro_tile_blend_plus`             | 2,681         | Lit terrain-blend, "plus" tier.                                                 |
| `lit_micro_tile_blend`                  | 2,323         | Lit terrain-blend, baseline.                                                    |
| `lit_weapon`                            | 2,203         | Viewmodel/weapon shader.                                                        |
| `lit_advanced_fullspec`                 | 2,062         | Lit "advanced" tier with full specular path. Common world hero materials.       |
| `lit_advanced_cc`                       | 1,804         | Lit "advanced" tier with colour-correction extras.                              |
| `lit_advanced`                          | 1,466         | Lit "advanced" tier baseline.                                                   |
| `lit`                                   | 1,366         | Plain lit, no extras.                                                           |
| `lit_decal_diffuse_spec_gloss`          | 805           | Decal carrying diffuse + spec + gloss channels.                                 |
| `lit_decal_diffuse`                     | 711           | Decal carrying just diffuse.                                                    |
| `lit_emissive`                          | 411           | Self-illuminating surfaces.                                                     |
| `lit_emissive_scroll_transparent`       | 374           | Scrolling emissive transparent (e.g. animated screens, scrolling text).         |
| `fxanim_vegetation`                     | 375           | FX-anim vegetation shader (separate family).                                    |

Naming-convention reading guide:

- `lit_*` — needs lighting (almost everything that's not HUD/effect).
- `_advanced_*` — full PBR-style slot set (color + normal + spec + gloss + AO).
- `_plus_*` — lighter tier than advanced, fewer slots.
- `_micro_tile_blend_*` — terrain-style tiling with blend masks.
- `_decal_*` — decal projection family.
- `_cc` — colour-correction extras (per-instance tinting).
- `_emissive` — self-illumination paths.
- `lit_weapon` — viewmodel-specific shader.
- `fxanim_*` — FX system animated materials, separate family.

> ⚠️ Picking the wrong `materialType` is one of the silent failure modes: the material may compile, the texture may bind, but the lighting will look wrong or some channels won't be sampled. The cargo-cult-from-existing-assets approach is genuinely the safest.

### Image slots

The slots a material exposes depend on its `materialType`. Common ones across world-geometry types:

| Slot field        | What goes in it                                                                |
| ----------------- | ------------------------------------------------------------------------------ |
| `colorMap`        | The diffuse/albedo image — base color of the surface.                          |
| `normalMap`       | The tangent-space normal-map image.                                            |
| `specMap`         | Specular reflectance image (or packed PBR data).                               |
| `glossMap`        | Gloss/roughness image.                                                         |
| `aoMap`           | Ambient occlusion image.                                                       |
| `emissiveMap`     | Self-illumination image.                                                       |
| `alphaMap` / `alphaRevealMap` | For dissolve/reveal effects.                                       |

The slot value is the **name** of an `image` GDT record, not a path. Example: `"colorMap" "i_t10_anh_dirt_rocky_pebbles_01_dec_disp_c"`.

### Render-side metadata

- **`sort`** — render queue / depth-sort hint, e.g. `<default>*`. Controls when in the frame this material draws (opaque vs transparent vs decal).
- **`surfaceType`** — physical/audio surface tag, e.g. `dirt`, `metal`, `wood`. Drives footstep sounds, bullet impact FX, decal projection. Matched against `surfaceFXTable` and `surfaceSoundTable` records elsewhere.
- **`usage`** — coarse hint, e.g. `terrain`, `weapon`, `viewmodel`.

### Collision/clip flags

A material can also carry collision behavior, used when applied to a brush:

- **`aiClip`** / **`aiSightClip`** — block AI / AI sight only.
- **`bulletClip`** / **`canShootClip`** — block bullets / be shot through.
- **`caulk`** — back-side caulk material (don't render this face).

These let you build invisible collision walls (`clip` brushes) by giving a brush a material with the right flag pattern instead of inventing a separate "clip" entity type. Classic idTech move.

### Shader constants: `cg00_x` / `cg00_y` / ...

Hundreds of fields named `cgNN_x/y/z/w`. These are **per-instance shader constants** packed into vec4 slots — the material's tunable knobs once a `materialType` is bound. Most are `0` for any given material because most shaders don't use most slots.

**Which `cg*` slots actually matter** is dictated by the **`.techsetdef`** files under `share/raw/techsetdefs_stable/`. A techsetdef binds shader-side parameters to specific `cg*` indices via syntax like:

```
x = <cg00_y>
y = <cg02_z>
```

So `cg02_z` only carries meaning if the techsetdef bound to your material's technique set references it. APE's per-`materialType` UI surfaces the relevant ones with friendly names so you usually don't touch the raw `cg*` indices yourself — but for a custom shader (or when debugging "why does this slot do nothing"), the techsetdef is where you check what's wired.

> 💡 Some of these slots are **scriptable** — i.e. driven by GSC at runtime — when the underlying HLSL shader references them as engine-driven inputs. Grep `share/raw/techsetdefs_stable/` for `<cgNN_X>` patterns to see real bindings.

## Worked example: a wall texture, end-to-end (the porting flow)

Most realistic flow in your project: you want a brick wall texture from BO6.

1. **Extract.** Greyhound (or Saluki) on a running BO6 session: tick the materials/images you want, Export. Output is a `.png` per image plus reference info that lets you reconstruct the material.
2. **Drop the bitmaps** into `texture_assets/_bo6/<asset_name>/`.
3. **Create image GDT entries** in APE — one per channel (`_c` color, `_n` normal, `_s` spec/gloss, etc.). For each, set `baseImage`, `semantic`, `coreSemantic`, `compressionMethod`. **The most common mistake**: leaving `coreSemantic` as the default sRGB on a normal map.
4. **Create the material GDT entry** in APE. Pick a `materialType` matching the shader the source game used (cargo-culting from a similar BO3 material if unsure). Bind your image names into `colorMap` / `normalMap` / etc. Set `surfaceType` so footsteps and impacts behave right.
5. **Apply** the material to a brush face in Radiant, or assign it to an `xmodel` slot in APE.
6. **Zone if needed.** If the material is on a brush in your `.map`, the BSP carries it. If you'll only swap to it from script, add `material,<name>` to your `.zone`.
7. **Link → Run.** Inspect in-game. If lighting looks off, your `coreSemantic` or `materialType` is the prime suspect.

For an authored texture, swap step 1 for "make the bitmap in Photoshop / NormalMap-Online / Substance" — the rest is the same.

## PBR authoring conventions (per Treyarch)

Treyarch ships `docs_modtools/PhysicallyBasedRendering.pptx` — their internal PBR primer. The technical model in BO3:

- **Full HDR lighting pipeline**, GGX specular + Oren-Nayar diffuse.
- **Max gloss `2^17`** (vs BO1's `2^13`) — finer-grained specular highlights.
- **Reflections** sampled from multiple **256×256 blended reflection probes** (vs BO1's single 128×128 per-object probe).
- **GI** baked as a **point cloud diffuse bake** (BO1 baked to lightmaps).
- **Both forward and deferred renderers coexist.** Deferred handles opaque geometry, GBuffer decals, and alpha-test (on/off) geometry. Forward handles transparency and "odd-ball" lighting (anisotropic, etc.). **Transparency is the single most expensive thing you can author** — try alpha-test first.
- **Decals** ride on top of the GBuffer and have a **Decal Layer** that must be set per intent (grime / dent / wet patch / etc.). They override parts of the underlying surface rather than blending atop it.

### Authoring rules — read these before painting

These are the rules Treyarch's slides explicitly call out. Get them wrong and your asset will look subtly off in ways that are hard to diagnose later.

- **Author everything in 16-bit.** Always. No 8-bit source files for the PBR pipeline.
- **Diffuse / albedo**: keep values in **40–238**. Below 40 = no headroom for shadows; above 238 = no headroom for specular highlights. **Exception**: pure metals get a *fully black* diffuse (their color comes from the spec channel).
- **Never paint AO or highlights into the diffuse map.** Those belong in `occlusionMap` and the lighting model. Pre-baked AO in diffuse breaks under dynamic lighting.
- **Specular**: most surfaces don't need a custom spec map — defaults are good. Author one only when you have **bare metal showing through** a non-metal (paint white = metal, black = non-metal, grayscale for transitions / dust / pitting around metal).
- **Gloss**: start with a fully white texture (max glossiness). Set the **gloss range in APE** per surface type — there are presets, plus `Custom` for high/low thresholds. Paint the *least* glossy areas black. Detail-normal-map roughness automatically modulates gloss too.
- **AO**: bake from a high-poly model. Don't darken everything — AO only affects *indirect* lighting (reflections + bounced diffuse). Over-baked AO makes models look like they have "creepy mustaches" growing on them under shadow.
- **Normal maps** are the highest-value map: highest resolution, most detail. Bake from a high-poly mesh. Use a **detail normal map** for any surface feature smaller than ~1/8 inch — those contribute to gloss variance and read as micro-roughness.
- **Substance Designer** is the recommended authoring environment. One graph emits diffuse / spec / normal / gloss in a coherent set; batch updates are cheap.

> 💡 The cleanest way to validate your authoring is to view the asset under multiple lighting conditions in-game (rotate the camera, change zone). Most PBR mistakes only become visible under specific angles or lighting probes.

## Common gotchas

- **Normal maps look wrong.** Almost always `coreSemantic` set to sRGB instead of linear. Fix the image record, re-Link.
- **Diffuse looks washed out.** Opposite case: linear sRGB. Same fix.
- **Material disappears** in some lighting conditions. Wrong `sort` value (transparent material drawn before opaques). Default usually fine; only worth touching when you know why.
- **Footsteps sound generic.** `surfaceType` not set or set to a value the surface tables don't recognize.
- **Texture is blurry on close inspection.** `mipBase` lowered, or `compressionMethod` too aggressive, or `noPicMip` not set on a HUD texture.
- **Texture missing in-game (purple checker / black).** The image GDT record exists but the bitmap path is wrong, or the bitmap isn't under any `converter_gdt_dirs_0.txt` root, or you forgot to zone the material when nothing else pulls it in transitively.
- **Duplicate material name** silently shadowed by another GDT (community pack collision). See [`07-community-packs.md`](./07-community-packs.md) for the dupe-purger script and the Midgetblaster-priority handling.

## Where to go next

Sequenced by where there's most lived signal in this project (not textbook order):

- **xmodels** (placement, LODs, material assignment, porting workflow) → [`03-models.md`](./03-models.md)
- **Fog** (per-zone / per-state — `fog` GDT type) → [`04-fog.md`](./04-fog.md)
- **SSI** (sun & sky info) → [`05-ssi.md`](./05-ssi.md)
- **Weapons** (`projectileweapon`, `bulletweapon`, attachments, `weaponcamotable`) → [`06-weapons.md`](./06-weapons.md)
- **Community packs & dupe handling** → [`07-community-packs.md`](./07-community-packs.md)

> Audio is *not* part of asset-pipeline in this doc — it's its own top-level section (see [`docs/audio/`](../audio/)). Sound aliases, zone configs, ambient rooms, mixing, VOX, and music sourcing are big enough together to warrant their own discipline.

For the full GDT asset-type taxonomy, see [`reference/asset-types.md`](../reference/asset-types.md).

---

## Reference reading

Texture authoring, shaders, and material-adjacent topics:

- [`LG-RZ/BlackOps3Shaders`](https://github.com/LG-RZ/BlackOps3Shaders) — community shader repository for custom material effects.
- [`Scobalula/Bo3MWStyleScope`](https://github.com/Scobalula/Bo3MWStyleScope) — example of a custom shader (MW-style zoom scope) showing how shader work plugs into the material pipeline.
- [Unshade](https://unshade.ie/) — shader-editing tool.
- [NormalMap-Online (cpetry)](https://cpetry.github.io/NormalMap-Online/) — quick browser-based normal-map generator from a height/colour input.
- [Custom Textured Flag tutorial (YouTube)](https://www.youtube.com/watch?v=wREUQgcWF5g) — small but useful walkthrough of custom 2D-image authoring through the GDT.

> Visionsets / LUT / colour-grading work lives in [`docs/lighting/visionsets-and-luts.md`](../lighting/visionsets-and-luts.md). Brush-side tool textures (caulk, nodraw, hint, skip, clip flags) are mapping concerns, not material-pipeline ones — see [`docs/mapping/brush-textures.md`](../mapping/brush-textures.md).

## Open questions / TODO

