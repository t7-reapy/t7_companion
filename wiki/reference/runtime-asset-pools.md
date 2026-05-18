# Runtime asset pools — `/listassetpool` index → name → cap

> Captured from a live BO3 session via the in-game console command `/listassetpool 0-102`. Each row gives the pool index (the integer you can pass to `listassetpool` to filter), the asset-type name, and the engine's runtime cap on simultaneously-loaded assets of that type.

> **Distinct from BGCache caps** ([`bgcache-caps.md`](./bgcache-caps.md)): BGCache is the per-script `#precache` directive ceiling. The runtime pool here is the engine's "maximum number of <type> assets active in one map" cap. Both are real ceilings; you typically trip BGCache first.

## How to read

- **Pool index** — pass to `/listassetpool <index>` to dump just that one pool. Or use `/listassetpool 0-102` for the whole table.
- **Cap** — engine's hard ceiling. The poolinfo CSV (`zone_source/<lang>/assetinfo/<map>_poolinfo.csv`) reports your map's current usage vs this cap.

## Full pool list (BO3, verified live)

| Idx | Type                              | Cap     |
| ---:| --------------------------------- | -------:|
|   0 | `physpreset`                      | 275     |
|   1 | `physconstraints`                 | 128     |
|   2 | `destructibledef`                 | 128     |
|   3 | `xanim`                           | 25 200  |
|   4 | `xmodel`                          | 11 264  |
|   5 | `xmodelmesh`                      | 35 840  |
|   6 | `material`                        | 22 528  |
|   7 | `computeshaderset`                | 256     |
|   8 | `techset`                         | 1 024   |
|   9 | `image`                           | 49 152  |
|  10 | `sound`                           | 32      |
|  11 | `sound_patch`                     | 16      |
|  12 | `col_map`                         | 2       |
|  13 | `com_map`                         | 2       |
|  14 | `game_map`                        | 2       |
|  15 | `map_ents`                        | 2       |
|  16 | `gfx_map`                         | 2       |
|  17 | `lightdef`                        | 32      |
|  18 | `lensflaredef`                    | 70      |
|  19 | `ui_map`                          | 0       |
|  20 | `font`                            | 16      |
|  21 | `fonticon`                        | 16      |
|  22 | `localize`                        | 25 600  |
|  23 | `weapon`                          | 1 536   |
|  24 | `weapondef`                       | 0       |
|  25 | `weaponvariant`                   | 0       |
|  26 | `weaponfull`                      | 0       |
|  27 | `cgmediatable`                    | 5       |
|  28 | `playersoundstable`               | 16      |
|  29 | `playerfxtable`                   | 16      |
|  30 | `sharedweaponsounds`              | 64      |
|  31 | `attachment`                      | 128     |
|  32 | `attachmentunique`                | 2 148   |
|  33 | `weaponcamo`                      | 512     |
|  34 | `customizationtable`              | 8       |
|  35 | `customizationtable_feimages`     | 8       |
|  36 | `customizationtablecolor`         | 1 024   |
|  37 | `snddriverglobals`                | 1       |
|  38 | `fx`                              | 2 000   |
|  39 | `tagfx`                           | 64      |
|  40 | `klf`                             | 70      |
|  41 | `impactsfxtable`                  | 256     |
|  42 | `impactsoundstable`               | 64      |
|  43 | `player_character`                | 8       |
|  44 | `aitype`                          | 96      |
|  45 | `character`                       | 150     |
|  46 | `xmodelalias`                     | 48      |
|  47 | `rawfile`                         | 5 000   |
|  48 | `stringtable`                     | 220     |
|  49 | `structuredtable`                 | 105     |
|  50 | `leaderboarddef`                  | 256     |
|  51 | `ddl`                             | 64      |
|  52 | `glasses`                         | 2       |
|  53 | `texturelist`                     | 8       |
|  54 | `scriptparsetree`                 | 1 150   |
|  55 | `keyvaluepairs`                   | 64      |
|  56 | `vehicle`                         | 64      |
|  57 | `addon_map_ents`                  | 1       |
|  58 | `tracer`                          | 100     |
|  59 | `slug`                            | 5       |
|  60 | `surfacefxtable`                  | 64      |
|  61 | `surfacesounddef`                 | 256     |
|  62 | `footsteptable`                   | 32      |
|  63 | `entityfximpacts`                 | 256     |
|  64 | `entitysoundimpacts`              | 256     |
|  65 | `zbarrier`                        | 16      |
|  66 | `vehiclefxdef`                    | 32      |
|  67 | `vehiclesounddef`                 | 32      |
|  68 | `typeinfo`                        | 0       |
|  69 | `scriptbundle`                    | 1 024   |
|  70 | `scriptbundlelist`                | 64      |
|  71 | `rumble`                          | 280     |
|  72 | `bulletpenetration`               | 1       |
|  73 | `locdmgtable`                     | 1       |
|  74 | `aimtable`                        | 12      |
|  75 | `animselectortable`               | 64      |
|  76 | `animmappingtable`                | 64      |
|  77 | `animstatemachine`                | 64      |
|  78 | `behaviortree`                    | 64      |
|  79 | `behaviorstatemachine`            | 128     |
|  80 | `ttf`                             | 48      |
|  81 | `sanim`                           | 1 024   |
|  82 | `lightdescription`                | 550     |
|  83 | `shellshock`                      | 64      |
|  84 | `xcam`                            | 532     |
|  85 | `bgcache`                         | 32      |
|  86 | `texturecombo`                    | 16      |
|  87 | `flametable`                      | 16      |
|  88 | `bitfield`                        | 52      |
|  89 | `attachmentcosmeticvariant`       | 640     |
|  90 | `maptable`                        | 25      |
|  91 | `maptableloadingimages`           | 25      |
|  92 | `medal`                           | 768     |
|  93 | `medaltable`                      | 32      |
|  94 | `objective`                       | 256     |
|  95 | `objectivelist`                   | 64      |
|  96 | `umbra_tome`                      | 0       |
|  97 | `navmesh`                         | 2       |
|  98 | `navvolume`                       | 2       |
|  99 | `binaryhtml`                      | 2 048   |
| 100 | `laser`                           | 50      |
| 101 | `beam`                            | 50      |
| 102 | `streamerhint`                    | 50      |

## Sample output format

`/listassetpool` responses look like:

```
Total of 1499/2000 assets in fx pool, bytes 215856
```

(`current / cap`, plus an aggregate byte count.)

> 💡 If the output isn't visible in-game, check the BO3 logfiles — see [`bgcache-caps.md`](./bgcache-caps.md#how-to-verify-whats-currently-declared) for the full caveat and the recommended log-channel command.

## Caps that look like 0

A few entries (`ui_map`, `weapondef`, `weaponvariant`, `weaponfull`, `typeinfo`, `umbra_tome`) report a cap of 0 in this engine build. Either deprecated, MP-only, or never wired up for zombies — don't try to use these.

## Tracing down a near-cap pool to a specific prefab

When `/listassetpool` (or the linker's `_poolinfo.csv`) shows you're getting close to a runtime cap, the next question is *"which prefab is bloating the count?"* The Radiant **Basic Map Stats** report (`map_stats`) breaks down per-prefab `Model Count` / `Fx Count` / `Entity Count` / `Light Count` — quickly identifies which prefab is responsible. Generation steps + CSV column reference: see [`asset-pipeline/01-overview.md`](../asset-pipeline/01-overview.md#generated-reports--diagnostic-files).

## Related

- [`bgcache-caps.md`](./bgcache-caps.md) — the *per-script `#precache`* caps (the ceiling you hit first when authoring scripts).
- [`docs/scripting/03-fx.md`](../scripting/03-fx.md) — FX-specific guidance citing both pool views.
- [`docs/asset-pipeline/01-overview.md`](../asset-pipeline/01-overview.md#generated-reports--diagnostic-files) — the linker emits a `_poolinfo.csv` per build that mirrors these caps with current usage.
