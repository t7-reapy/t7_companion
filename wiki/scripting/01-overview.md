# Scripting — Overview

> ~128 GSC/CSC files in `usermaps/zm_test/scripts/`, organised across hellround / weather / room-of-thanks / harrybo21 weapons / utility modules. Reapy's second-deepest area after mapping. About a year of lived experience post-March 2025, captured in the commit history.

> **The biggest "learned the hard way" lesson** isn't the language itself — GSC's syntax is small and approachable. It's everything *around* the language: API discovery, asset caching, how scripts interact with engine systems (weapons, FX, fog, lighting, **triggers / unitriggers, animations, scriptbundles, sounds**, and other typed engine structs), and the script-override mechanism. Most of this section is dedicated to those non-obvious things rather than language reference.

## Mental model

BO3 ships **two parallel scripting languages** with the same syntax but different runtimes:

- **GSC** (`.gsc`) — runs **server-side**. Game logic, gameplay state, AI behaviour, score, round flow, weapon awards, save data. The authoritative game state lives here.
- **CSC** (`.csc`) — runs **client-side**, per local player. Visual effects, audio triggers, HUD-adjacent logic that needs to react instantly without a server round-trip. Splitscreen runs N CSC instances simultaneously, one per local client.

Plus a third layer:

- **GSH** (`.gsh`) — header / preprocessor includes. `#define` macros, constants, struct layouts. `#insert`-ed by both GSC and CSC files.

```
.gsc  →  GSC bytecode  →  zoned via scriptparsetree  →  runs on server
.csc  →  CSC bytecode  →  zoned via scriptparsetree  →  runs on each client
.gsh  →  text-included via #insert  (no separate compile, no zone line)
```

The split matters because **state on one side isn't visible to the other** without an explicit channel. Crossing the server↔client boundary requires **clientfields** (the typed bit-channel propagation system) — covered in `clientfields.md` *(planned)*. If you're trying to make a CSC react to something the GSC knows, you go through a clientfield.

## File-header anatomy

Every script starts with a familiar pattern:

```c
#using scripts\shared\callbacks_shared;     // imports another module's namespace
#using scripts\zm\_zm_utility;
#using scripts\zm\_hb21_zm_weap_staff_fire; // community pack module

#insert scripts\shared\shared.gsh;          // text-include macros / constants
#insert scripts\shared\version.gsh;
#insert scripts\zm\_zm_utility.gsh;

#namespace zm_test;                         // this script's namespace

REGISTER_SYSTEM("zm_test", &init, undefined) // optional: register a startup hook
```

- **`#using`** imports another script's *namespace* — its public functions become callable as `<namespace>::function(...)`.
- **`#insert`** is text-include, like C `#include`. Used for `.gsh` macro/constant headers, NOT for code modules.
- **`#namespace`** declares this file's namespace (one per file).
- **`REGISTER_SYSTEM`** (a macro from `system_shared.gsh`) registers an init hook that runs when the script system bootstraps. **`REGISTER_SYSTEM_EX`** is the variant for systems that also need a `main()` function to autoexec after init.

> ⚠️ **`#using` is not zone-checked at compile time.** The linker doesn't verify that every `#using`-ed module's assets actually got zoned, so you can have a using that points at something half-installed and the map *compiles* fine. It blows up at **runtime** instead — the BO3 linker console fills with `Could not find fx "..."` / `Could not find tagfx "..."` lines followed by `Fatal Error: Error linking script 'scripts/zm/...'`. **Moral**: keep `#using` declarations and your zone files in sync.

## The script-override mechanism — *the most useful thing to know*

Any shipped script under `share/raw/scripts/...` can be **shadowed** by placing a file with the same path under `usermaps/<your_map>/scripts/...`. Your version replaces Treyarch's at link time for *your map only*.

Example from `zm_test`:

- `share/raw/scripts/zm/_zm_pack_a_punch.gsc` — Treyarch's PaP machine logic
- `usermaps/zm_test/scripts/zm/_zm_pack_a_punch.gsc` — Reapy's override (adds `_up_up` second-upgrade chain via `self.second_upgrade_cost = 11500`)

The override pattern is how `_up_up`, custom hellround logic, weather state, room-of-thanks elevator, and most of `zm_test`'s features are wired without forking Treyarch's whole script tree.

> ⚠️ **The `zm_patch.csv` gotcha**: some shipped scripts are pre-bundled in the `zm_patch` fastfile (`share/raw/scripts/zm/zm_patch.csv` lists them). To override one of *those*, you may need to first comment out its reference in the patch CSV, otherwise the linker pulls in both versions and the wrong one wins. Classic pitfall the BO3 modding community knows well.

## Why scripting feels opaque (even after you understand the syntax)

The language is small. The hard part is everything around it. Five recurring sources of "what is this even doing" frustration:

### 1. API discovery — "what function do I call?"

There are ~3 000 built-in script API functions across GSC + CSC. **Modern reality: you do get IntelliSense** thanks to community VS Code extensions — [GSCode](https://marketplace.visualstudio.com/items?itemName=blakintosh.gscode) is the workhorse, providing hover docs, signatures, and go-to-definition out of the box. Install it first; it changes the discovery experience entirely.

When IntelliSense isn't enough or you want to browse / cite / verify:

- **`share/raw/scripts/`** — grep this directory to see how Treyarch *actually uses* a function in real code. The most useful single source of truth — start here.
- **[`api-reference.md`](./api-reference.md)** — the GSCode JSON ground truth, with direct links and schema reference.
- **`shiversoftdev/t7-source`** — broader source dump when `share/raw/` is silent.
- **Discord** — when nothing else helps, ask in MT / ZGC / T7 servers.

The "right" function for a given task is often non-obvious and only learned by seeing it used in context. Even with IntelliSense, the *which-function* question routinely needs a `share/raw/scripts/` grep to answer.

### 2. Asset precaching — "why did my model lag the first time it appeared?"

Most engine-system assets need to be **precached** so the engine has them ready when script first asks for them. Skip precaching and the asset usually still loads — but on first use the engine has to fetch it synchronously, which manifests as a hitch / lag spike at the worst possible moment (the first zombie spawn, the first pickup, the first FX trigger).

Precaching in BO3 happens via the **`#precache(<type>, <name>)` script directive** at the top of a `.gsc` / `.csc` file — NOT runtime function calls. Real types in use across `share/raw/scripts/` and `usermaps/zm_test/scripts/`:

| `#precache` type   | What it precaches                                                       |
| ------------------ | ----------------------------------------------------------------------- |
| `fx`, `client_fx`  | FX assets. `client_fx` is the CSC-side variant.                         |
| `model`, `xmodel`, `client_model` | Static / world models, viewmodels.                       |
| `client_tagfxset`  | Tag FX sets (FX bound to a model tag).                                  |
| `material`         | Material assets.                                                        |
| `script_bundle`    | Scriptbundles (animated scripted scenes).                               |
| `lui_menu`, `lui_menu_data`, `menu` | LUI menus and supporting data.                         |
| `statusicon`       | Status icons (HUD-side glyphs).                                         |
| `string`, `eventstring`, `triggerstring` | Localized / runtime strings.                      |
| `objective`        | Objective definitions.                                                  |
| `vehicle`          | Vehicle assets.                                                         |
| `locationselector` | Location-selector asset (mortar / killstreak targeting).                |

Real examples from this project:

```c
// from the fauna animation script (zm_animated_fauna)
#precache( "script_bundle", "p7_fxanim_cp_lotus_atrium_ravens_bundle" );

// from share/raw/scripts/shared/ai/archetype_apothicon_fury.csc
// (harrybo21's apothicon furies AI)
#precache( "client_fx", FURY_DAMAGE_EFFECT );
```

The convention for FX you'll trigger from script:

```c
// In init or main:
level._effect["my_fx_name"] = "path/to/fx.efx";

// Then later:
PlayFX( level._effect["my_fx_name"], origin );
```

The `level._effect[...]` dictionary is **populated separately on GSC and CSC sides** by convention — script writers assign to the same key on both sides where each side needs it. Real usage in shipped script: `share/raw/scripts/zm/_zm_weapons.csc:298` reads `level._effect["870mcs_zm_fx"]` on the client side. Used **221 times** across `zm_test`'s scripts.

If an asset doesn't appear (or causes a hitch on first use) and the linker found it (per `<lang>/assetinfo/zm_test.csv`), the next suspect is "is it actually `#precache`d in the script that uses it?"

### 3. How weapons actually work from script

The weapon-record GDT is half the story; the other half is the script-side `level.zombie_weapons[...]` registry, the upgrade chain in `_zm_weapons.gsc`, the wallbuy spawning logic, the box rotation. To add a custom weapon you don't just zone its GDT — you also `#using` the right script modules and may need to extend the registry. See [`asset-pipeline/06-weapons.md`](../asset-pipeline/06-weapons.md) for the GDT half; the script half deserves its own page.

### 4. Engine-system interaction patterns

The engine exposes many typed systems through specific function families that all have their own conventions. Each system has its own setup ritual:

- **Fog**: `SetWorldFogActiveBank(client, 1 << index)` + `SetLitFogBank(...)` (CSC, per-player loop).
- **Lighting state**: `util::set_lighting_state(N)` on level or self ([`asset-pipeline/05-ssi.md`](../asset-pipeline/05-ssi.md)).
- **Clientfields**: register on a target type → set on server → handle in CSC callback. Specific flag soup (`CF_HOST_ONLY`, `CF_CALLBACK_ZERO_ON_NEW_ENT`, defined in `share/raw/scripts/shared/version.gsh:102+`).
- **Exploders**: declare in Radiant, trigger via `Exploder("name")` in script.
- **Brush triggers** placed in Radiant (`trigger_use`, `trigger_radius`, `trigger_damage`, etc.) — entities that fire script callbacks when a player enters/uses them.
- **Unitriggers** — *script-created* triggers (no Radiant brush). Spin one up via `set_up_unitrigger(width, length, height, LOOKAT, HINTSTRING, TRIGGERED_FUNC, UPDATE_PROMPT)` in `_load.gsc`. Useful for craftables, dynamic wall-buys, anything where you want a use-prompt that doesn't exist as a brush.
- **Animations & scriptbundles** — playing canned animations on entities (doors, scripted scenes, fxanim props) via `scriptbundle_scene` lookups, `play_anim_*` helpers, scene tags.
- **Sounds** — alias-driven via `PlayLoopSound` / `PlaySound` / `PlaySoundAtPosition` etc.; aliases compiled into `.sabl` (loaded) and `.sabs` (streamed) banks, with load-vs-stream chosen per-sound in CSV files referenced from a `.szc` zone config (zoned with `sound,<name>`). The **GSC and CSC sides have different sound APIs** — different functions available, different signatures. Always check [`api-reference.md`](./api-reference.md) before assuming a function exists or behaves the same on both sides.
- **FX** — `level._effect[name]` cache + `playfx`/`playfxoncamera`/`playfxontag`.

Each system is small, but the conventions are different and need to be learned individually. Each will get a dedicated reference page in [`docs/scripting/`](./).

### 5. Silent-by-default failures (until you turn the lights on)

GSC has no try/catch in the modern sense. A wrong type passed to a built-in often **returns `undefined` and continues**, leading to cascading silent breakages. **By default** there's nothing on screen — the script just stops doing what you expected.

But it's not actually black-box once you enable developer mode. Launch with **`+set developer 2`** and:

- **Intrusive in-game popups** appear when an exception fires.
- A **stack trace prints to the expanded console** (`Shift+~` on QWERTY, `Shift+²` on AZERTY).
- For shipped Treyarch scripts you usually get a useful source-line reference. **For map-mod scripts, source-line attribution is often missing** — you get the call stack but not the line, so you still do guess-work to pinpoint the bug. Less painful than fully silent, more painful than a real debugger.

Combined with `+set logfile 2 +set scr_mod_enable_devblock 1`, `AssertMsg` calls also start firing as actual fail-fasts (otherwise they're no-ops in release builds).

The standard debug pattern:

- **`PrintLn` / `IPrintLn`** (and the **Bold** variants `PrintLnBold` / `IPrintLnBold`) — quick logging. `PrintLn*` to console, `IPrintLn*` to all players' on-screen feed. The Bold versions are stylistically preferred in `zm_test` (they stand out more), but the non-Bold ones work too — pick by the visual prominence you want.
- **`AssertMsg`** — explicit fail-fast for impossible states. Only fires under the developer-mode launch flags noted above.
- **dvars as feature flags** — `GetDvarInt("debug_hellround")` etc. for runtime toggles without a re-Link.
- **Sphynx dev commands** — the closest thing to a debugger; full cheat sheet in-game.

Covered in detail in `debugging.md` *(planned)*.

## Tooling at a high level

| Tool                                         | Role                                                                |
| -------------------------------------------- | ------------------------------------------------------------------- |
| **GSCode** (blakintosh, VS Code extension)   | The best GSC/CSC editor extension; still actively maintained, supports cross-script reference / go-to-definition / region folding. Install this first; check the GitHub README for setup details. |
| **[`shiversoftdev/t7-source`](https://github.com/shiversoftdev/t7-source)** | The Treyarch source dump — already in cleartext. **Reach for this first** when you want to read shipped scripts; no decompile needed. |
| **[Cerberus](https://github.com/Scobalula/Cerberus-Repo)** (Scobalula, decompiler) | GSC/CSC bytecode → readable source. **Repo is archived** as of 2026, but still works for the cases t7-source doesn't cover (community-shipped compiled scripts, third-party `.ff` dumps via `acts`). |
| **`shiversoftdev/t7-compiler`** (often called "the GSC Injector") | [Repo + setup guide on GitHub](https://github.com/shiversoftdev/t7-compiler). Can hot-reload compiled GSC bytecode into a running BO3 process — but in practice it's **buggy** and prone to user-error crashes. The pragmatic alternative everyone settles on (per Discord consensus): **keep the game open, re-Link in the modtools launcher, then `map_restart` or `fast_restart` from the in-game console** — no game close needed. That loop ends up faster than fighting t7-compiler edge cases. Setup steps if you do try it: see the repo's README. |
| **`acts gscd`** (ATE47)                      | Decompile compiled GSC dumped from a `.ff` (see [`asset-pipeline/01-overview.md`](../asset-pipeline/01-overview.md)). |
| **Sphynx dev commands**                      | Runtime cheat-sheet of in-game console commands (spawn zombies, skip rounds, give perks, etc.). |

Detail in `tooling.md` *(planned)* and [`api-reference.md`](./api-reference.md).

## Common pitfalls

This section is grounded in real `zm_test` commit history — every entry below has been hit and fixed in this project at least once.

- **`%`-prefixed animations in CSC silently won't play.** Not a hard crash, just no playback and no error log. **Almost impossible to debug** — full detective mode, eliminating possibilities one at a time, before realising the prefix character was the issue. The animation work carries scars from this kind of debugging — see commits `77a87a04` ("Added fauna animations to the map **after many attempt and trials**") and `2ea45dc7` ("fix(tweak): siege anim in hellround now works"). Use the unprefixed animation name in CSC.
- **Forgot to register a clientfield's `CF_CALLBACK_ZERO_ON_NEW_ENT` flag** → late-joiners / splitscreen partners don't get the state. Repeatedly hit on splitscreen polish work — `e380d5ae` ("Fixed 2 bugs in hellround environment scripts for splitscreen") is one of several. See SSI page Common Gotchas for the working pattern.
- **Wallbuy hard limit (~20 wall weapons)** before the map silently breaks when Pack-a-Punching. See [`asset-pipeline/06-weapons.md`](../asset-pipeline/06-weapons.md).
- **`endon("end_game")` missing from a long-running thread** → loops keep ticking after a `fast_restart`, causing apparent lag and double-spawned entities. See `4d55ef3f` ("Forbid duplicated loops of spawners") for an instance of the duplicated-loop class of this bug.
- **Calling `set_lighting_state` on the wrong scope** (`level` vs `self`). Use `level` for shared environment changes; `self` is rare. (See SSI page.)
- **Splitscreen sound deduplication** — emitting an FX/sound once per player when both are in the same area double-plays it. Hit several times in `zm_test`: `8cbb07db` ("teddy bear doesn't play twice in splitscreen"), `c4c31bfb` ("Fixed perk sound + ambient sound in splitscreen"). Pattern: host-only check + `GetLocalPlayers()` loop.
- **Splitscreen LUI bugs** — typewriter-style overlays draw twice or in the wrong viewport in splitscreen. See `87075c22` ("fixed type writing in splitscreen").

## Where to go next

Sequenced by depth-of-need:

- **GSC language reference** → `gsc-language.md` *(TODO)* — syntax, threads, `endon`/`notify`, namespaces, `array::thread_all`, common idioms. Primary source: [`docs_modtools/GSC_Language.pdf`](file:///D:/SteamLibrary/steamapps/common/Call%20of%20Duty%20Black%20Ops%20III/docs_modtools/GSC_Language.pdf). Function-API content cross-references [`api-reference.md`](./api-reference.md), which scrapes the GSCode JSON for the up-to-date truth (the local `bo3_scriptapifunctions.htm` is known to be outdated and unreliable).
- **Clientfields** (the server↔client sync mechanism, late-joiner flag soup) → `clientfields.md` *(TODO)*.
- **Script overrides** (override patterns, the `zm_patch.csv` workflow) → `script-overrides.md` *(TODO)*.
- **Engine quirks** (silent failures, hard limits, the `%`-CSC trap, wallbuy cap) → `engine-quirks.md` *(TODO)*.
- **Tooling** (GSCode / Cerberus / t7-compiler setups) → `tooling.md` *(TODO)*.
- **API reference cross-walk** → [`api-reference.md`](./api-reference.md).
- **Debugging** (dvars, Sphynx, t7-source reading) → `debugging.md` *(TODO)*.
- **LUI / Lua HUD** → `lui-lua.md` *(TODO)*. *Honest note*: this layer was largely AI-assisted in `zm_test` rather than authored from scratch — coverage will be lighter here.
- **Custom upgrade chains** (`_up_up` script-side wiring, hellround camo rewards on win) → woven into `script-overrides.md` and the engine-quirks page.

## Reference reading

- [GSCode (blakintosh, VS Code extension)](https://marketplace.visualstudio.com/items?itemName=blakintosh.gscode) — install first.
- [Modme Forums — Scripting](https://forum.modme.co/wiki/forums/21.html) — the long tail of community Q&A.
- [Abnormal202 Scripting Tutorials (UGX)](https://www.ugx-mods.com/forum/scripting/91/abnormal202-scripting-tutorials-master/16746/) — solid intro tutorials.
- [BO3 Source Code Explorer (zeroy)](https://bo3explorer.zeroy.com/) — the API search engine.
- [`shiversoftdev/Black-Ops-3-Projects`](https://github.com/shiversoftdev/Black-Ops-3-Projects/tree/main) — released BO3 mod projects, useful as reference patterns.
- [`XcDylan932/ReDempTion`](https://github.com/XcDylan932/ReDempTion/tree/main) — Dylan's full mod source code (a **mod** with custom UI, not a map). Covers Custom Pause Menu, Wonder Weapon FX, Character/Zombie Models, Camos, Zombie Eyes, Powerup FX, AAT FX, Custom HUD, Deadwire Color. Strong reference repo when you need patterns for FX / character work / gameplay tweaks / QoL features.
- [`docs_modtools/GSC_Language.pdf`](file:///D:/SteamLibrary/steamapps/common/Call%20of%20Duty%20Black%20Ops%20III/docs_modtools/GSC_Language.pdf) — Treyarch's language reference.

---

## Open questions / TODO

- [ ] Document the full clientfield flag set on its dedicated page. Source: `share/raw/scripts/shared/version.gsh:102+` defines the flag macros (`CF_HOST_ONLY`, `CF_CALLBACK_ZERO_ON_NEW_ENT`, etc.).
