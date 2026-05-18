# SSI — Sun & Sky Info

> SSI is the GDT type that defines the **sun + skybox + exposure setup** for a region of the world. In `zm_test`, swapping SSI is what makes the world *feel* different in normal play vs hellround vs thunder vs the room-of-thanks: different sun direction, different shadow softness, different dynamic range, different skybox.

This page also unlocks the **lighting-state mechanism** that we hand-waved on the fog page — SSI is the cleanest place to learn it because lighting state is *the* canonical activation path for SSI.

## Mental model

An SSI record describes the *global lighting environment*: where the sun is, how its shadows behave, what skybox you see, how exposure clamps the dynamic range, what cookie pattern the sun casts.

```
ssi GDT record (named)   ──►   volume_sun entity slot ssiN in Radiant   ──►   GSC: util::set_lighting_state(N)   ──►   active
```

This is the **canonical example of the lighting-state mechanism** in BO3 — the same pattern is reused (with subtle variations) by lights and fog volumes, but SSI is the cleanest place to learn it because every meaningful piece of the system is right here in plain sight.

The shape:

- **The volume**: a brush entity called **`volume_sun`** carved over the region the SSI applies to.
- **The slots**: **`ssi`, `ssi1`, `ssi2`, `ssi3`, `ssi4`** — the bare `ssi` is the default fallback; `ssi1`..`ssi4` are the four user-addressable lighting states (4 SSI max, hard cap). Each numbered slot is paired with a **`ssiN_runtime_override`** holding the SSI actually applied at runtime (see below).
- **The activation API**: **`util::set_lighting_state(N)`** — server-side GSC. Not a low-level `SetSomething(...)` call.

That single `set_lighting_state(N)` call simultaneously drives three subsystems:

1. **SSI swap** on every `volume_sun` (the slot we're documenting here).
2. **Light toggles** for any light entity whose `lightingstate_N` KVP is `1`.
3. **Per-entity-type lighting-state gates** across the rest of the map. The same `LIGHTSTATE_N` / `lightingstate_N` flag pattern shows up on many entity types, each subscribing to the same enum:
   - **Fog volumes** (`volume_worldfog` / `volume_litfog`) — `LIGHTSTATE_N = 1` enables that bank for the state. Combined with the script's direct `SetWorldFogActiveBank` calls — the two paths can complement each other (see [`04-fog.md`](./04-fog.md)).
   - **`misc_model`** entities — `lightingstate4 = 0` skips the model in the GI bake for state 4 (model has no baked lighting in that state).
   - **FX entities** — `lightingstate4 = 0` prevents the FX from playing during state 4. *Quirk*: the engine only re-evaluates whether to play on a *state switch*; default is "play," so an FX whose `lightingstate1 = 0` will still play on map start until the state is toggled at least once.

So SSI, lights, and fog all subscribe to **one project-wide enum**, the lighting state. Setting it once flips the whole visual mood of the world atomically.

## The `ssi` GDT record

~30 fields per record. Authored in APE. Five logical groups.

### 1. Sun direction & shadows

- **`pitch`** — sun elevation in degrees (`60.0` ≈ mid-morning / mid-afternoon; `0` = horizon; `90` = directly overhead).
- **`yaw`** — sun azimuth in degrees (`300.0` = roughly NW direction).
- **`enablesun`** — toggles the **sun lensflare** (the visible sun disc / glare in the sky). `0` for overcast/storm/indoor states where you want to suppress the lensflare. The sun's directional lighting itself is always part of the SSI's pitch/yaw/colour pipeline — there's no "kill the sun lighting" switch, since the lighting environment as a whole is the sun.
- **`dynamicShadow`** — `1` to enable real-time shadows from the sun on dynamic objects. Static-only world geo gets baked shadows regardless. *Aesthetic trade-off*: dynamic shadows look noticeably worse than baked (lower resolution, often noisy edges), but they add a real layer of realism by tying moving objects (zombies, swung doors, dropped items) into the lighting environment instead of letting them float without ground contact.
- **`penumbra_inches`** — softness of the shadow edge in world inches. Smaller = sharper. Useful tell from `zm_test`: `house = 8.0` (soft daytime), `house_thunder = 2.5` (sharp lightning-light edges).

### 2. Exposure & dynamic range

These control the auto-exposure / tonemapper. Beyond the obvious dynamic-range work, the exposure knobs are also the practical handle for "wet map" looks — pushing the camera response so reflective surfaces, water sheen, and rain-coated materials read as luminous against the rest of the scene.

- **`ev`** — current exposure value override (`0` = use auto).
- **`evmin`** / **`evmax`** — auto-exposure clamp range, in stops.
- **`evcmp`** — exposure compensation (additive offset).
- **`stops`** — dynamic range width. `7.0` for normal scenes; **`12.0` for `house_thunder`** — a clever choice: widening DR during thunder lets lightning flashes register their full intensity instead of clipping.
- **`bounceCount`** — global illumination bounce count for the bake (`4` is typical).
- **`spec_comp`** — specular compensation multiplier.

### 3. Sun colour

- **`colorSRGB`** — RGBA sun light colour. In `zm_test` the tint is `0.516 0.685 0.991 1` — heavily blue-shifted, reading more like cool **moonlight** than direct sunlight (the map's main vibe is overcast / dim / unsettled).

### 4. Sun cookie (gobo / shadow-projection pattern)

A "cookie" is a texture projected through the sun light to break up its perfect-disc shadow into something irregular — leaves rustling, clouds drifting, etc. Hugely useful for outdoor atmosphere.

- **`sunCookieLightDefName`** — link to a `lightdescription` asset that holds the cookie pattern. Reapy's records use `cookie_sun_mp_apartments`.
- **`sunCookieIntensity`** — `1.0` normal, `0.0` disables the cookie entirely (`house_thunder` does this — no leaf-shadow effect during a storm).
- **`sunCookieAngle`** — rotation in degrees.
- **`sunCookieScale`** — scale factor.
- **`sunCookieRotation`** — additional rotation knob (separate from angle).
- **`sunCookieOffsetX`** / **`sunCookieOffsetY`** — cookie position offset.
- **`sunCookieScrollX`** / **`sunCookieScrollY`** — animated scroll rate (gives the "drifting clouds" feel).
- **`sunVolumetricCookie`** — `1` to also evaluate the cookie inside the volumetric lighting pass, so god-rays carry the cookie pattern.

### 5. Skybox & lens flare

- **`skyboxmodel`** — the `xmodel` asset name of the skybox (e.g. `skybox_zm_castle`). Skyboxes are 3D models, not panoramic textures — the engine renders them at the world horizon.
- **`lensFlare`** — link to a `lensflare` asset (or empty for none).
- **`lensFlarePitchOffset`** / **`lensFlareYawOffset`** — offset the flare from the actual sun direction.

## The `volume_sun` entity

The Radiant counterpart to the SSI registry. Real example from `map_source/zm/zm_test.map`:

```
{
  "classname"          "volume_sun"
  "ssi"                "default_day"
  "ssi1"               "house"
  "ssi1_runtime_override"  "house_override"
  "ssi2"               "house_power_on"
  "ssi2_runtime_override"  "house_power_on_override"
  "ssi3"               "house_thunder"
  "ssi3_runtime_override"  "house_thunder_override"
  "ssi4"               "house_hell"
  "ssi4_runtime_override"  "house_hell_override"
  "shadowSplitDistance"    "1500"
  "global_fill_color"      "0 0 0"
  "global_fill_intensity"  "1"
  "grid_density"           "32"
  "respectLightLods"       "1"
  "shadowBiasScale"        "1"
  "shadowVistaDetail"      "1"
  "streamLighting"         "1"
  // brush data follows
}
```

**Field reference:**

- **`ssi`** — bank 0 (default). Usually a generic shipped record like `default_day`.
- **`ssi1`** … **`ssi4`** — banks 1–4. Each is the name of an `ssi` GDT record.
- **`ssiN_runtime_override`** — **the SSI that's actually applied at runtime**. The base `ssiN` slot, in contrast, is the SSI **the GI bake uses** to compute static lighting / probes / lightmaps. So the pair is: `ssiN` = bake-time lighting reference, `ssiN_runtime_override` = what the player actually sees in-game. Source: confirmed via a Discord message in the modding community. This inverts the obvious naming — "override" sounds like an exception you reach for, but in practice it *is* the live lighting setup, while the base is what the offline tools consult.
- **`shadowSplitDistance`** — distance (inches) at which cascade-shadow splits occur. Higher = softer transitions, more VRAM.
- **`global_fill_color`** / **`global_fill_intensity`** — uniform ambient fill applied alongside the sun.
- **`grid_density`** — light-grid sample spacing for this volume.
- **`respectLightLods`** — `1` to honour LOD culling on lights inside this volume.
- **`shadowBiasScale`**, **`shadowVistaDetail`** — shadow-map tuning knobs.
- **`streamLighting`** — `1` enables on-demand streaming of lighting data rather than preloading it all. *Exactly what gets streamed isn't fully pinned down* — likely a mix of light grid samples + reflection probe data + per-light shadow buffers, since those are the largest lighting payloads. To verify if it ever matters, look at `usermaps/zm_test/zone_source/english/assetinfo/zm_test.csv` — its `resident` vs `streamed` byte columns show what got loaded vs streamed for lighting-tagged assets.

The brush carved as the volume defines the *region* this SSI registry applies to. Like fog volumes, you can have multiple `volume_sun` entities in one map, each carrying its own SSI palette, applied wherever the camera is.

## Activation: the lighting state mechanism

The activation path is intentionally *unified* across SSI / lights / fog:

```c
// GSC, server-side
util::set_lighting_state(N);   // N is 0, 1, 2, 3, or 4
```

That single call simultaneously:

1. **Swaps the active SSI** on every `volume_sun` to slot `ssiN`.
2. **Toggles every light entity** whose `lightingstate_N` KVP is `1`.
3. **Activates the matching bank** on fog volumes whose `LIGHTSTATE_N` flag is set (combined with the script's direct `SetWorldFogActiveBank` calls — the two paths can complement each other; see [`04-fog.md`](./04-fog.md)).

### Real call sites in `zm_test`

From `usermaps/zm_test/scripts/zm/zm_test.gsc`:

```gsc
util::set_lighting_state(level.power_on_lightstate);   // power switch flip
```

From `zm_test.gsc` room-of-thanks callbacks:

```gsc
function private set_lighting_state_clear()
{
    util::set_lighting_state(2);   // entering the elevator/clear state
}
```

From `zm_hellround_environment.gsc`:

```gsc
self util::set_lighting_state(lightstate);    // per-player apply (splitscreen-safe)
level util::set_lighting_state(lightstate);   // world-level apply
```

From `zm_weather_thunder.gsc`:

```gsc
level util::set_lighting_state(thunder_lightstate);
// ... lightning flash duration ...
level util::set_lighting_state(self.lightstate_missing);  // restore previous
```

Note the pattern: `level` set vs `self` (player) set — `set_lighting_state` can be called on a player entity for per-player visual states, or on `level` for global effects. Both forms exist in the project.

## Reapy's setup in `zm_test`

Eight SSI records in `source_data/custom/house.gdt`, used as four (state, override) pairs:

| Lighting state | Base SSI            | Runtime override                | Purpose                                     |
| --------------:| ------------------- | ------------------------------- | ------------------------------------------- |
| 0 (`ssi`)      | `default_day`       | —                               | Engine default.                             |
| 1 (`ssi1`)     | `house`             | `house_override`                | Baseline daytime "lights off" state.        |
| 2 (`ssi2`)     | `house_power_on`    | `house_power_on_override`       | After power switch; lights-on baseline.     |
| 3 (`ssi3`)     | `house_thunder`     | `house_thunder_override`        | Storm / thunder peak — wider DR, sharp shadows, no cookie. |
| 4 (`ssi4`)     | `house_hell`        | `house_hell_override`           | Hellround state.                            |

> 🛠️ *Authoring note*: the `_override` records here were created by **duplicating** the base records in APE. In hindsight that's confusing — APE actually supports an *override* mechanism for child records that would have made the relationship explicit and avoided drift between the base and override field values. If you're starting fresh, prefer creating each `_override` via APE's override mechanism rather than as a copy.

Cross-reference the corresponding lighting-state semantics on fog volumes, lights, and `volume_sun` — the index maps consistently across all three subsystems.

### The clever bit: thunder dynamics

`house_thunder` doesn't just darken the sky — it *re-tunes the entire camera*:

- `stops 12.0` (vs `7.0` baseline) — wider dynamic range so a lightning flash registers as a real flash instead of clipping to white.
- `penumbra_inches 2.5` (vs `8.0` baseline) — sharper shadow edges, characteristic of harsh storm light cutting through narrow cloud breaks.
- `sunCookieIntensity 0.0` (vs `1.0` baseline) — no leaf-shadow cookie pattern; storm light is too diffuse for those.

Worth noting because it shows SSI is not just "where is the sun" — it's the entire camera tuning of the scene, swappable as a unit.

## Worked example: adding a new SSI state

> ⚠️ **Hard cap: 4 user-addressable lighting states.** `volume_sun` only exposes `ssi1`..`ssi4` (plus the unnumbered `ssi` default). There is no `ssi5`. **Workarounds** if you need more: (a) repurpose one of the existing four, or (b) build a *separate `volume_sun` in a different region of the map* with its own bank palette, then teleport the player into that region for the duration of the new state — gives you a clean fresh 4 slots without touching the others.

Say you want a "fog-of-war" lighting state, repurposing slot 4 (we'll move hellround somewhere else, or re-frame state 4 as "war" since hellround already cycles through transitions).

1. **Author the SSI records** in APE in `source_data/custom/house.gdt`:
   - `house_war` (base, used by the GI bake).
   - `house_war_override` (the actual runtime SSI). Author this via APE's override mechanism on top of the base when possible, rather than duplicating.
2. **Wire to slot 4** on the relevant `volume_sun` entity: set `ssi4 = "house_war"` and `ssi4_runtime_override = "house_war_override"`.
3. **Set the state from script**: `util::set_lighting_state(4);`.
4. **Coordinate with fog and lights**: any `volume_worldfog`/`volume_litfog` for the same region needs `LIGHTSTATE_4 = 1` and `fsi_4` populated; any lights you want toggled per state need `lightingstate_4 = 1` on their entity KVPs.
5. **Compile (All) → Light → Link → Run.** Lighting changes need a Light step because GI is baked from the *base* `ssiN` (not the override).

## Common gotchas

- **SSI swap looks wrong / muddy.** GI bake hasn't been re-run after changing the base SSI. Re-run Light, not just Link. Light bake captures the ambient/bounced light from the *baseline* SSI; runtime overrides blend on top but can't fix a wrong-bake.
- **In-game sun and lighting look black / completely unbaked even after a successful Light step.** Stale baked lighting data wasn't refreshed by the Link step. Lighting bake artifacts live in **both** files: the `GfxWorld` asset record (lightmap UVs, light grid samples, probe references) sits in the `.ff`, and the heavy mip data of lightmap textures sits in the paired `.xpak`. **Delete both** `usermaps/zm_test/zone/zm_test.ff` and `usermaps/zm_test/zone/zm_test.xpak` (plus their `en_*` / `fr_*` variants if needed) and re-run Link to regenerate the lot from the fresh bake.
- **Lighting state changes but nothing happens visually.** No light or volume on the map has the matching `lightingstate_N` / `LIGHTSTATE_N` flag, and the `volume_sun` slot for that state is empty or points at `default_day`. The mechanism succeeded; nothing was wired to it.
- **Per-player state doesn't sync to splitscreen partner.** `set_lighting_state` was called on `level` instead of `self` (or vice versa). For shared environment changes use `level` — this is the common case. The `self` form is rare; reach for it only when you actually need *per-player* visual states (e.g. one player in a special-vision mode while the other isn't).
- **Thunder flash doesn't visually punch.** Auto-exposure clamps the camera's brightness range via `evmin` / `evmax`. If `evmax` is too low, the lightning flash hits the cap and reads as just "kind of bright" — same as a normal scene. Widen `evmax` (and `stops`) in the thunder SSI to give the flash headroom. `zm_test`'s `house_thunder` does this right: `stops 12.0` vs daytime `7.0`.
- **Skybox doesn't change between states.** Most projects reuse one `skyboxmodel` across all SSI records and modulate *colour/exposure* across states instead. **Exception worth knowing**: `zm_test`'s `house_hell` and `house_hell_override` swap to `skybox_hellround` (vs `skybox_zm_castle` for the other states), so the *sky itself* visibly changes during hellround. If you want this kind of state-driven skybox swap, set a different `skyboxmodel` per SSI record.
- **Setting a lighting state on player join doesn't take.** Same late-joiner pattern as fog — make sure your initialisation path applies the current lighting state on client connect, or use a clientfield with `CF_CALLBACK_ZERO_ON_NEW_ENT` to broadcast it. The function `util::set_lighting_state()` called with no argument re-applies `level.lighting_state` to `self`, which is the natural shape for an `on_connect` callback.

## Where to go next

- **Lighting** (light bake, light grid, probes, reflection probes, the Light step, `lightdescription`, sun cookies) → `docs/lighting/` (forthcoming top-level section).
- **Mapping** (`volume_sun` and other volume entities, lighting-state script wiring details) → `docs/mapping/` (forthcoming).
- **Weapons** → [`06-weapons.md`](./06-weapons.md).
- **FX** (`fx` / `.efx` files — animated/particle/visual effects, used heavily in `zm_test` for hellround visuals, weather, weapon impacts, ambient props) → planned `08-fx.md`. *We've been treating FX implicitly throughout the SSI/fog/weapons pages but it deserves its own asset-pipeline page.*

---

## Reference reading

Lighting / sun / sky / colour-grading topics that intersect SSI:

- [Lighting Tutorial (YouTube playlist)](https://www.youtube.com/watch?v=hPweAEu8zJY&list=PLYAF4YwatlpU47HOddWaEPJlQt3yjoNtY) — the canonical multi-part walkthrough of the BO3 lighting system that SSI sits inside.
- [Visionsets and Overlays Tutorial (YouTube)](https://www.youtube.com/watch?v=D9Uq0zo-NrU) — visionsets / LUT colour-grading; swapped at the same time as SSI states for full-mood transitions.
- [Change Look of Map / LuT Edit (YouTube)](https://www.youtube.com/watch?v=p6iXtx50bxI) — practical companion to the visionset/LUT swap pattern.
- [LightningStrikeEffects (zeroy wiki)](https://wiki.zeroy.com/index.php?title=Call_of_duty_bo3:_LightningStrikeEffects) — relevant for the thunder lightingstate, where lightning flashes interact with the wider exposure range.
- [Epic Rotating World Script — Rotating Skybox (YouTube)](https://www.youtube.com/watch?v=F_S_mVl2VLo) — animated skybox model, ties into the SSI's `skyboxmodel` field. There's also a **dvar that changes the sky rotation in-game** for live tuning without a re-Link (look it up in the Sphynx dev-command set or `dvar list`).

## Open questions / TODO

- [ ] **`shadowSplitDistance` × sun pitch** — best guess: steeper pitch (sun more overhead) means shadows are shorter, so longer split distance is *less* needed for cascade quality. So a high-pitch SSI can usually afford a smaller `shadowSplitDistance` for a perf win. Not formally confirmed.
