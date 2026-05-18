# BGCache caps — per-script `#precache` ceilings

> **BGCache** = "**Both Game** cache" — distinct from server-only and client-only caches. Each precache type has a hard cap on how many entries can be registered. Hit the cap → in-game error dialog `BG_Cache_GetIndexInternal - Exceeded 'N' items for type 'X'` and the map fails to load.

## What BGCache actually is (verbatim, from a community-shared explanation)

The first parameter to `#precache` is the *type* of precache list. Each type has its own pool with its own cap. The list of pool names + caps below is the canonical reference — you'll see the **`type`** name appear in error messages exactly as written.

> "string" is an alias to `locstring` which is used for HUD elements and certain other things. `triggerstring` is specifically for triggers. They're different BGCache pools, so the same literal string can end up registered in *both* if you precache it as `string` and then use it in a `triggerstring` context (like `SetHintString` on a trigger) — the engine re-registers it lazily at runtime under the new pool.
>
> The linker takes advantage of the fact that all clients on the same fastfile have the same BGCache layout — it pregenerates a list of cached assets at link time. That way the server can use **cache indexes** instead of alias strings when network-replicating fields, which is much smaller on the wire.
>
> Some BGCache pools allow **dynamic allocation at runtime** via the configstrings table that's sent to clients. `triggerstring` and `localizedstring` allow this — when you call `SetText(string)` or `SetHintString(string)` and the string isn't precached, the game registers it in the configstring table and records that string at the next available BGCache index.

So caching is partly compile-time (`#precache` directives → linker pre-fills the table) and partly runtime (dynamic allocation for the strings pools). The cap applies to the total: precached + runtime-allocated entries combined cannot exceed it.

## Cap table

| Type                       | Cap   | Notes                                                                    |
| -------------------------- | -----:| ------------------------------------------------------------------------ |
| `models`                   | 4 096 | World/static/script models.                                              |
| `weapons`                  |   512 | All weapon GDT records the script declares.                              |
| `rumbles`                  |   128 | Controller rumble presets.                                               |
| `shellshocks`              |    64 | Shell-shock effect presets.                                              |
| `xcams`                    |   256 | Cinematic / scripted cameras.                                            |
| `destructibles`            |    64 | Destructible definitions.                                                |
| `streamerHints`            |    64 | Streaming hint volumes.                                                  |
| `headicons`                |    15 | Above-player head icons.                                                 |
| `statusicons`              |     8 | HUD status icons.                                                        |
| `locationselectoricons`    |    15 | Mortar / killstreak target icons.                                        |
| `scriptmenus`              |    64 | Script-driven UI menus.                                                  |
| `materials`                |   512 | Material assets.                                                         |
| `localizedstrings`         | 2 048 | Localized strings (per language).                                        |
| `debugstrings`             |   512 | Debug-only strings.                                                      |
| `eventstrings`             |   256 | Event-trigger strings.                                                   |
| `triggerstrings`           |   250 | Trigger hint strings.                                                    |
| `objectivestrings`         |   256 | Objective text.                                                          |
| **`serverfx`**             |   256 | **GSC** `#precache("fx", ...)` declarations.                             |
| `luimenus`                 |    64 | LUI menu definitions.                                                    |
| `luimenudata`              |   128 | LUI menu data records.                                                   |
| **`clientfx`**             | 1 024 | **CSC** `#precache("client_fx", ...)` declarations.                      |
| `clienttagfxset`           |    64 | `client_tagfxset` declarations (multi-tag FX bundles).                   |

## When you exceed a cap

The error dialog format is consistent across types. Real screenshots from in-game:

```
BG_Cache_GetIndexInternal - Exceeded '256' items for type 'fx'
BG_Cache_GetIndexInternal - Exceeded '250' items for type 'triggerstring'
BG_Cache_GetIndexInternal - Exceeded '1024' items for type 'client_fx'
```

The `'N'` in the message is exactly the BGCache cap, and `'<type>'` is the pool name as listed in the table above. So the dialog tells you which pool is full and at what number.

> *Note*: this is a separate failure mode from the **"Could not find fx …"** console-log errors, which are usually a *zoning* issue (the FX wasn't included in the fastfile) — not a BGCache cap. Don't conflate them.

## Caps vs the engine asset-pool limits

These BGCache caps are **per-script `#precache` declaration limits** — how many entries the precache table can hold. They are different from the **engine asset-pool limits** that the linker reports in `usermaps/<map>/zone_source/<lang>/assetinfo/<map>_poolinfo.csv` and that `/listassetpool` queries at runtime. The poolinfo CSV is genuinely useful: it reports your map's *current usage* against the engine's *hard ceiling* per type, so you can see at a glance which pool is closest to filling. (Same shape as the [`runtime-asset-pools.md`](./runtime-asset-pools.md) reference, but scoped to your specific build.)

Example for FX (the most-discussed pool):

| Cap                                        | Number    | Source                                                |
| ------------------------------------------ | ---------:| ----------------------------------------------------- |
| BGCache `clientfx` precache cap            | **1 024** | This page (community-shared list)                     |
| BGCache `serverfx` precache cap            | **256**   | This page                                             |
| Engine `fx` pool cap (runtime)             | **2 000** | `/listassetpool` live (see [`runtime-asset-pools.md`](./runtime-asset-pools.md)) |
| Engine `fx` pool cap (linker poolinfo CSV) | **1 875** | `<map>_poolinfo.csv` — this number can drift from the live `listassetpool` value depending on engine version / patch state |

The BGCache caps are stricter and apply to *what you declare in scripts*. The engine pool is the runtime limit on *active asset instances*. Both real ceilings; the BGCache cap is the one most people see first because it fires the named error dialog.

## How to verify what's currently declared

Two paths:

1. **At link time** — `-printbgcache` and `-writebgcachetofile` exist as args (see [`modtools-cli-args.md`](./modtools-cli-args.md)). Caveat: these may need to be passed directly to `linker_modtools.exe` / `modtools.exe` rather than the launcher GUI — the launcher doesn't always forward custom args to the linker process. Behaviour not fully verified; experiment if you need it.
2. **At runtime** — open the in-game console (the `~` key, or `²` on AZERTY) and run:

   ```
   /con_channelshow *;logfile 1;listassetpool
   ```

   This first prints **the list of pool types with their indices** (one row per pool — name + index, no per-asset detail). It also enables full console channels and opens a logfile. Use this run to find the index of the pool you care about.

   Then to dump the **detailed contents of a specific pool** (or a range), pass the index:

   ```
   /listassetpool 38           # just the fx pool
   /listassetpool 0-102        # every pool
   ```

   Each detail row looks like `Total of 1499/2000 assets in fx pool, bytes 215856` — current/cap plus byte aggregate. (That example is **index 38 = `fx`**.) Full index → name → cap mapping in [`runtime-asset-pools.md`](./runtime-asset-pools.md).

   > 💡 Console output is often **not visible on screen** in BO3 — sometimes it scrolls off, sometimes the channel filter swallows it. The reliable place to read it is one of the engine logfiles at the BO3 folder root: `console_mp.log` or `console.log`. Tail or open these after running `listassetpool` to see the full dump.

> 💡 Once you know *which pool* is near-full, see [`runtime-asset-pools.md`](./runtime-asset-pools.md#tracing-down-a-near-cap-pool-to-a-specific-prefab) for the per-prefab tracing workflow.

## Cross-VM rule

The two script VMs (server / client) **don't share precaches**. Precaching an FX (or any other asset) on the GSC side does not make it available on the CSC side, and vice versa. If both sides need the asset, both sides must `#precache` it. This effectively *doubles* the precache budget you're working with for any cross-side asset, but each side counts separately.

## Related

- [`docs/scripting/03-fx.md`](../scripting/03-fx.md) — FX recipes including the precache discussion.
- [`docs/scripting/01-overview.md`](../scripting/01-overview.md) — overview-level coverage of the `#precache` directive and asset-precaching mental model.
- [`docs/asset-pipeline/01-overview.md`](../asset-pipeline/01-overview.md#generated-reports--diagnostic-files) — the linker reports including the `_poolinfo.csv` engine-pool view.
