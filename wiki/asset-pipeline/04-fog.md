# Fog

> Per-zone / per-state atmospheric fog. The most-iterated visual subsystem in `zm_test`: weather, hellround state changes, the bloody-environment shift. This page is heavier than the others in this section because the lived experience is deepest here.

## Mental model

Fog in BO3 is **declared per-record in GDT**, **wired to bank slots inside per-region brush volumes in Radiant**, and **swapped at runtime by bank index**. Name → bank → active.

```
fog GDT record (named)   ──►   volume_worldfog/volume_litfog entity slot N in Radiant   ──►   CSC: SetWorldFogActiveBank(1 << N)
```

This is a variation of the **lighting-state mechanism** — the broader system that also drives SSI swaps and per-state light toggles. **The canonical reference for that mechanism is [`05-ssi.md`](./05-ssi.md)**, which lays out the full activation path (`util::set_lighting_state(N)`, the unified state enum, propagation across SSI / lights / fog). Read that first if you haven't yet; the fog page below assumes you've internalised it and focuses on the fog-specific quirks.

Fog has two notable departures from the SSI baseline:

1. The **slot field naming** is `fsi`/`fsi_N` (fog state index) rather than `ssi`/`ssiN`.
2. Beyond the lighting-state path, fog also exposes a **direct CSC API** (`SetWorldFogActiveBank`, `SetLitFogBank`) that scripts can call to swap fog independently of the project-wide lighting state — useful for fine-grained transitions (`zm_test`'s hellround fog uses this for a 3-step lerp through a transition bank). The two paths can complement each other: lighting state still gates which banks are *available* via `LIGHTSTATE_N`; the direct API picks which available bank is *active*.

The runtime never sees fog *names*. It addresses fog by **bank index** (an integer 0..N). The mapping from name to bank lives on the *fog volume entity* covering the player's current region, baked into the BSP at Compile time. So adding a new fog state is two steps: (1) author the record in APE, (2) wire it to a bank slot on every fog volume that should expose that state. Then your script swaps banks.

There are **two parallel fog systems** that you usually wire (and swap) together:

- **World fog** — the cheap, traditional distance-based atmospheric haze. Per-pixel fog factor based on camera distance + height. Wired via **`volume_worldfog`** entities. Activated by `SetWorldFogActiveBank`.
- **Lit fog** — fog that interacts with light sources (sun, lamps, spots). It's the system whose `litfog "1"` toggle **enables volumetric lighting** (god rays, light shafts, lamps glowing through mist, sunbeams through windows), but it's not *only* volumetric lighting — it also handles broader fog↔light interactions (fog colour shifting under coloured lights, fog appearing darker in shadow, etc.). Wired via **`volume_litfog`** entities. Activated by `SetLitFogBank` (bank-as-single-index, smooth lerp between source and target). Significantly more expensive than world fog.

In practice you keep an *atmosphere* world-fog record (`_a` suffix, e.g. `house_fog_a`) on a `volume_worldfog`, and a paired *volumetric* lit-fog record (`_v` suffix, e.g. `house_fog_v`) on a `volume_litfog` covering the same area. Switching state means flipping both bank indices at once. The `_a` / `_v` suffix convention reads as **a**tmospheric vs **v**olumetric.

> 💡 Conceptual split worth internalising: *world fog* answers "how much haze is between the camera and that wall?"; *lit fog* answers "how does that haze look when light passes through it from this lamp / the sun / the shadow over there?" Two different physics, two different APIs, but you author them as a coordinated pair so the look stays consistent.

> 💡 The volume-based architecture is what lets `zm_test` use one fog palette inside the apartment and a different one in the courtyard. Each region of the map gets its own `volume_worldfog`+`volume_litfog` pair, each with its own per-region bank → record mapping. The bank *index* is a global runtime addressing scheme; what each index *means* depends on which volume the camera is in.

## The `fog` GDT record

~45 fields per record. Authored in APE. Five logical groups.

### 1. World fog (the broad atmospheric haze)

The basic distance-fog the engine has had since forever. Computes a per-pixel fog factor based on distance from camera.

- **`worldfog`** — `1` = enable world fog for this record.
- **`fogcolor`** — RGB (+ alpha unused) of the fog tint.
- **`fogintensity`** — overall opacity multiplier.
- **`fogopacity`** — alpha curve modulator.
- **`halfdist`** — distance (in inches) at which fog reaches 50% blend. Lower = thicker fog up close.
- **`halfheight`** — vertical falloff (lets fog be thicker at low altitudes).
- **`basedist`**, **`baseheight`** — starting offsets.
- **`distribution`** — falloff curve shape (0 = sharper, 1 = softer).
- **`densityscaler`** — global density multiplier.
- **`worldfogskysize`** — size of the fog dome wrapping the skybox (e.g. `65000`).

### 2. Sun fog (light scattered through fog from the sun direction)

Adds the "sunset glow", "god ray", or sun-disc-through-mist effect.

- **`sunfog`** — `1` = enable.
- **`sunfogcolor`** — RGB of the scattered sun colour.
- **`sunfogintensity`**, **`sunfogopacity`** — strength.
- **`sunfoginner`** / **`sunfogouter`** — angular falloff (degrees from sun direction). Inner = full strength, outer = zero.
- **`sunpitchoffset`** / **`sunyawoffset`** — offset the perceived sun direction without moving the actual sun.
- **`suncoloroverride`** + **`suncoloroverridenenabled`** — override the sun colour for fog calculation only (so fog can be coloured differently from the actual lighting).
- **`sunintensityscale`** — multiply the sun fog contribution.

### 3. Atmosphere fog (volumetric height-fog, the thick stuff)

The richer, modern fog — supports sun shafts, height layers, denser haze.

- **`atmospherefog`** — `1` = enable.
- **`atmospherefogcolor`** / **`atmospherefogdensity`** — base colour and density.
- **`atmospherehazecolor`** / **`atmospherehazedensity`** / **`atmospherehazespread`** — separate **haze** layer on top of base atmosphere fog.
- **`atmospherehazebasedist`** / **`atmospherehazefadedist`** — haze distance bounds.
- **`atmosphereextinctionstrength`** — how much light gets eaten by the atmosphere along its path.
- **`atmosphereinscatterstrength`** — how much ambient light scatters into the camera ray (the "blue sky" effect).
- **`atmospherepbramount`** — PBR contribution of the atmosphere.
- **`atmospheresunenabled`** — couple atmosphere fog to the sun direction.
- **`skyhalfheightoffset`** — vertical offset for the sky-half-height calculation.

### 4. Lit fog (fog × lights, including volumetric lighting)

Turns the fog into a **medium that interacts with light sources** — the sun, omni lights, spot lights all influence how this fog looks. Setting `litfog "1"` is what turns on volumetric lighting (god rays, light shafts, lamps glowing through mist, sunbeams through windows), but lit fog also covers other fog↔light interactions: fog colour shifting under coloured lights, fog darkening in shadow, fog inheriting probe-baked light. So volumetric lighting is the most visually striking outcome but not the only thing this group enables.

- **`litfog`** — `1` = enable lit-fog evaluation for this record (which includes volumetric lighting). **Significant performance cost** — render-time scales with the lit volume size and number of contributing lights.
- **`maxlitsunfogdistance`** / **`maxlitomnispotfogdistance`** — clamp distances beyond which the volumetric march stops sampling sun / omni+spot lights. Lower = cheaper but kills shafts at distance.
- **`probebakelitfogdensityscaler`** / **`probebakeworldfogdensityscaler`** — how much lit / world fog density contributes to the baked light probes (so probes "know" they're inside fog).
- **`probecontributionscaler`** — how much probes contribute back into fog colouring (the inverse direction — so fog inherits the surrounding light environment).

> Light grids, light probes, baked GI, and the Light step generally are their own discipline — see `docs/lighting/` (top-level section) for how the lighting subsystem works. This page only covers the *fog-side* knobs that control how fog interacts with that system.

### 5. Misc / shared

- **`albedo`** — fog "surface" albedo for indirect colour bounce.
- **`extcolor`** — extinction colour (the colour the world *loses* through the fog).
- **`type`** — always `fog`.

## How fog gets activated: name → bank → runtime

Two-stage wiring:

### Stage 1 — bank assignment via fog volume entities (compile-time)

Inside Radiant, you paint a brush volume covering the region you want a fog palette to apply to, and give it the classname **`volume_worldfog`** (for world/atmosphere fog) or **`volume_litfog`** (for lit/volumetric fog). Each entity carries its own bank-to-record mapping in its key/value pairs.

Real example from `map_source/zm/zm_ai_test.map`:

```
{
  "classname"      "volume_worldfog"
  "fsi"            "default"
  "fsi_1"          "house_fog_a"
  "fsi_2"          "default"
  "fsi_3"          "house_hell_fog_a_transition"
  "fsi_4"          "house_hell_fog_a"
  "BANK_1"         "1"
  "BANK_3"         "1"
  "BANK_4"         "1"
  "LIGHTSTATE_1"   "1"
  "LIGHTSTATE_2"   "1"
  "LIGHTSTATE_3"   "1"
  "LIGHTSTATE_4"   "1"
  "fogtime"        "0.5"
  "ENABLE_SUN_FOG" "1"
  "AUTO_PRIORITY"  "1"
  "DISABLE_FOG"    "0"
  // brush data follows
}
```

And the paired `volume_litfog` covers the same region with `_v` (volumetric) records:

```
{
  "classname"        "volume_litfog"
  "fsi"              "default"
  "fsi_1"            "house_fog_v"
  "fsi_2"            "default"
  "fsi_3"            "house_hell_fog_v_transition"
  "fsi_4"            "house_hell_fog_v"
  "BANK_1"           "1"
  "BANK_3"           "1"
  "BANK_4"           "1"
  "LIGHTSTATE_1..4"  "1"
  "fogtime"          "0.5"
  "ENABLE_LIGHTS"    "1"
  "ENABLE_SUN"       "1"
  "ambientColor"     "0 0 0"
  "ambientIntensity" "1"
  "AUTO_PRIORITY"    "1"
  // brush data follows
}
```

**Field reference:**

- **`fsi`**, **`fsi_1`** … **`fsi_4`** — fog **state index** slots. The bare `fsi` is bank **0**; `fsi_N` is bank **N** (1–4). Each value is the name of a `fog` GDT record, or `"default"` for an empty/passthrough slot. **5 slots total per volume**, so you get bank indices 0–4.
- **`BANK_1`** … **`BANK_4`** — boolean flags marking which numbered banks are "in use" on this volume. (Bank 0 is implicit.)
- **`LIGHTSTATE_1`** … **`LIGHTSTATE_4`** — gate each bank by the map's current **lighting state**. Working hypothesis (per Reapy): a fog bank only activates if its `LIGHTSTATE_N` flag matches the currently-active lighting state, letting you keep different fog palettes for "lights on" vs "lights off" vs intermediate states. Cross-confirmed by `radiant/configs/RadiantFilters.json`, which exposes `lightingstate1..4` as filterable KVPs (so they're a real, project-wide Radiant concept). Lighting states themselves are a mapping topic — see `docs/lighting/` (forthcoming).
- **`fogtime`** — default transition lerp time when this volume's fog activates.
- **`AUTO_PRIORITY`** — flag related to volume nesting / overlap resolution. Niche use case — most maps don't have multiple *active* fog volumes overlapping the same point, so this rarely matters. If you do end up with overlap, expect to figure out the resolution rule by experimentation.
- **`DISABLE_FOG`** — kill switch for the whole volume.
- **`ENABLE_SUN_FOG`** (worldfog only) — couple to the sun direction.
- **`ENABLE_LIGHTS`** / **`ENABLE_SUN`** / **`ambientColor`** / **`ambientIntensity`** (litfog only) — lit-fog-specific controls.

> 💡 The `fsi`/`fsi_N` ↔ bank-index mapping resolves an earlier ambiguity: `HRENV_FOG_INDEX_NORMAL = 0` corresponds to the unnumbered `fsi` field on the volume; `_TRANSITION = 2` to `fsi_2`; `_BLOODY = 3` to `fsi_3`. Index `1` is "deliberately skipped" in `zm_test` because that volume reserves `fsi_1` for some other purpose (or just leaves it as a spare slot to keep numbering aligned).

> ℹ️ **Related but separate Radiant volume classes.** Full list of `volume_*` entity classes Radiant's entity browser exposes (verified from the editor's "Volume" entity-type submenu): `attenuation`, `exposure`, `fpstool`, `lightclip`, `litfog`, `outdoor`, `performance`, `reflection`, `sun`, `vista`, `weathergrime`, `worldfog`, `worldfogmodifier`, plus `Sun_extension`. Notable ones for this page: **`worldfogmodifier`** (locally tweaks worldfog without replacing it), **`reflection`** (reflection-probe volumes), **`sun`** / **`Sun_extension`** (sun setup, covered in [`05-ssi.md`](./05-ssi.md)). The rest are mapping/lighting concerns and live in [`docs/mapping/`](../mapping/) and [`docs/lighting/`](../lighting/) rather than here.

**Multiple volumes for multiple regions.** Critically, you can have many `volume_worldfog`/`volume_litfog` pairs in one map — the apartment, the courtyard, the underwater section, etc. Each carries its own bank → record mapping, so the same bank index can mean *different* fog records in different regions of the map. The engine applies whichever volume the camera is currently inside. **Avoid overlapping multiple *active* fog volumes** — it's a niche pattern most maps never need, and the resolution rules are undocumented.

### Stage 2 — bank swap from CSC at runtime

Engine APIs (CSC-side):

```gsc
// In a CSC script
SetWorldFogActiveBank(client_number, fog_bank_bitmask);
SetLitFogBank(client_number, from_bank_index, to_bank_index, transition_time);
```

- `SetWorldFogActiveBank` takes a **bitmask**, but in practice you must activate exactly one bank — `fog_bank = 1 << index`. Composing multiple (`(1<<0) | (1<<3)` for layered fog) **does not work**: confirmed by experiment, the engine breaks visually rather than blending. Treat this API as "single bank only, expressed as a bitmask."
- `SetLitFogBank` takes a **from-bank** and **to-bank** index plus a transition time, so the lit fog *lerps* smoothly between the two banks over `transition_time` seconds. Pass `-1` as `from_bank_index` to mean "from whatever's currently active."
- Both are **per-player** in splitscreen. Iterate `GetLocalPlayers()` and pass each `client_number` separately.

### The transition pattern

In `zm_test`, the hellround swap uses a 3-step transition rather than a hard cut:

```
1. Active bank → TRANSITION bank (lerp over RADIANT_TIME)
2. Wait TRANSITION_TIME real seconds
3. TRANSITION bank → target bank (lerp over RADIANT_TIME)
```

This avoids a snappy "bloody fog SLAMS in" feel — you get a gradual desaturation/fade through a transitional palette, then settle into the new state. Real values from `usermaps/zm_test/scripts/zm/hellround/zm_hellround_environment.gsh`:

```c
#define HRENV_FOG_INDEX_NORMAL      0
#define HRENV_FOG_INDEX_TRANSITION  2
#define HRENV_FOG_INDEX_BLOODY      3
#define HRENV_FOG_RADIANT_TIME    0.50  // lerp time per bank swap
#define HRENV_FOG_TRANSITION_TIME 2.00  // hold on transition state before final swap
```

And the actual swap function (`zm_hellround_environment.csc`):

```c
function fog_update(b_hellfog, skip_transition)
{
    fog_index = (b_hellfog ? HRENV_FOG_INDEX_BLOODY : HRENV_FOG_INDEX_NORMAL);
    if (!skip_transition)
    {
        set_fog_index(HRENV_FOG_INDEX_TRANSITION, HRENV_FOG_RADIANT_TIME);
        waitrealtime(HRENV_FOG_TRANSITION_TIME);
    }
    set_fog_index(fog_index, HRENV_FOG_RADIANT_TIME);
}

function private set_fog_index(index, transition_time)
{
    fog_bank     = 1 << index;
    lit_fog_bank = index;
    foreach (player in GetLocalPlayers())
    {
        client_number = player GetLocalClientNumber();
        SetWorldFogActiveBank(client_number, fog_bank);
        SetLitFogBank(client_number, -1, lit_fog_bank, transition_time);
    }
}
```

Note the **`skip_transition` flag** — used when initialising fog for a freshly-joined client (you want the correct fog *now*, not a 2-second fade from whatever the engine defaulted to).

### Triggered via clientfield

Fog state is a piece of *world state*, not per-player state, so it propagates from server to client through a clientfield registered on the `"world"` entity:

```c
clientfield::register("world", HRENV_TOGGLE_CLIENT_FIELD,
                      VERSION_SHIP, 1, "int",
                      &hellround_environment,
                      !CF_HOST_ONLY,
                      CF_CALLBACK_ZERO_ON_NEW_ENT);  // fires for late-joiners too
```

The `CF_CALLBACK_ZERO_ON_NEW_ENT` flag is the **fix for splitscreen / late-join fog desync** — without it, a player joining mid-hellround would see normal fog because the clientfield wouldn't fire on their client. (See `scripting/clientfields.md` when written.)

## Worked example: adding a new fog state

Say you want a "smoke" fog state for a new event.

1. **Author 2 fog records** in APE in your project GDT (e.g. `source_data/custom/house.gdt`):
   - `house_smoke_fog_a` — world+atmosphere variant (sets `worldfog`, `atmospherefog`, picks colour and densities).
   - `house_smoke_fog_v` — lit variant (sets `litfog`, picks colour and density that interacts with lights).
2. **Wire to a bank slot** on the relevant `volume_worldfog` and `volume_litfog` — e.g. bank 4. (Skip already-used 0/2/3.)
3. **Trigger from script** by setting the clientfield to a new value the CSC handler maps to bank 4 (extend `fog_update` or write a parallel helper).
4. **Compile (All) → Light → Link → Run.** Test in-game; iterate on the GDT until the look is right (only Link needed for GDT-only changes).

## Splitscreen considerations

Two players can be in different fog environments only if your gameplay design allows — but the **engine fog state is per-client**, not per-world. The `foreach (player in GetLocalPlayers())` loop is mandatory. Skip it and only player 1 gets the fog change.

The README note *"rain sounds in splitscreen should only be played once if both players are in same environment"* applies the same way to fog: if both local players are in the same zone, you may want to only emit the *transition FX/sound* once (host check), but the fog bank change itself must run for both clients.

## Common gotchas

- **Fog "glitches thicker" after a `fast_restart`.** Documented in your README known-bugs list. Likely cause: a CSC handler doesn't reset to bank 0 on level restart, leaving the previous run's bank active until the next state change. Workaround / fix: have the init path explicitly call `set_fog_index(NORMAL, 0)` with a zero-second transition once the world loads.
- **Late-joining client sees default/wrong fog.** Missing `CF_CALLBACK_ZERO_ON_NEW_ENT` on the registering clientfield. Add the flag.
- **One player in splitscreen has wrong fog.** `GetLocalPlayers()` loop missing — only one client got the API call.
- **Fog color looks "off" at sunset only.** Your `sunfog*` values are interacting badly with `atmospherefog` — try toggling `sunfog 0` to confirm which subsystem is to blame, then re-tune.
- **Lit fog (volumetric lighting) kills framerate.** `litfog "1"` is expensive — drop `maxlitomnispotfogdistance` and `maxlitsunfogdistance` to clamp the sampling distance, and check whether the affected zone really needs god-rays. Sometimes turning lit fog *off* in a region and leaning on world fog only is the right call for performance.
- **Probes look wrong inside fog.** `probebakelitfogdensityscaler` / `probebakeworldfogdensityscaler` were set without re-running Light. Fog→probe interaction is baked at Light time, not Link time.
- **Trying to layer two fog banks `(1<<0) | (1<<3)` to blend them.** Doesn't work — engine breaks visually. The bitmask API shape is misleading; in practice it's "exactly one bank at a time." If you want to *cross-fade* between fogs, use the transition-bank pattern; if you want *layered* fog look, author it into a single record (use both world fog and atmosphere fog within one record's fields).

## Where to go next

- **Visionsets** (post-FX colour grading, often swapped at the same time as fog) → covered later, possibly as a separate page

## Reference reading

Fog-adjacent material that intersects this page:

- [Thunder & Lightning tutorial (YouTube)](https://www.youtube.com/watch?v=D1HjKRy1RzA) — pairs naturally with fog state changes during storm sequences.
- [LightningStrikeEffects (zeroy wiki)](https://wiki.zeroy.com/index.php?title=Call_of_duty_bo3:_LightningStrikeEffects) — engine-side lightning FX reference, useful when wiring storm states to fog swaps.

> Sibling topics with their own pages: lighting & lit-fog details in [`docs/lighting/`](../lighting/), visionsets/LUT colour-grading in [`docs/lighting/visionsets-and-luts.md`](../lighting/visionsets-and-luts.md), weather effects (rain drops, thunder system, sky animation) in [`docs/lighting/weather.md`](../lighting/weather.md).

