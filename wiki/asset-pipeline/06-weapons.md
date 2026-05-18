# Weapons

> ~40 weapons across the `zm_test` roster, spanning WW2 (M1 Garand, MP40, Thompson) through modern (AS VAL, Peacekeeper, Vintorez) to wonder weapons (Ray Gun, Tesla, Thundergun, Shredder). **Most ported via SkyeLord's CoD-port packs**, with several exceptions:
> - **Shredder** (sourced separately, replaced an earlier CNG attempt that's still visible in git history under `b7ac4108` / `12ed0070` reverts).
> - **Several harrybo21 weapons** ported from older CoDs — actual list in `zm_test`'s zone files: **Blundergat** (`t8_shotgun_blundergat_zm`, MOTD), **Magmagat** (`t8_shotgun_magmagat_zm` + `_upgraded_zm`, MOTD), **Fire Staff** and **Lightning Staff** (Origins-derived, see `_hb21_zm_weap_staff_*` scripts), and the **Black Hole / Gersh device** (Cosmodrome, see `_hb21_zm_weap_black_hole_projectile`). Plus their second-upgrade variants for the ones Reapy extended.
> - **Blast-O-Matic** (OwenC137) — installed in commit `f78a683f` alongside harrybo21's weapons pack (which it depends on as a prerequisite).
>
> Reapy never authors a weapon record from scratch. This page reflects that reality — it's an "integrator's manual" for community-pack weapons, not a "how to design a gun" tutorial.

> Heads-up: weapon GDT records are huge (~1300 fields per record, often ~1300 lines of text) **because the engine supports a huge variety of options** — every fire mode, attachment slot, ADS curve, recoil pattern, and material binding is its own field. Most defaults are sensible. This page does **not** enumerate all of them — it gives you the mental model, the meaningful field clusters, and the integration touchpoints (the CSV, the camo table, the upgrade chain).

## Mental model

A "weapon" in BO3 is not one record — it's a **family**:

```
base record  (e.g. iw8_asval)             ← what you spawn with
upgrade      (e.g. iw8_asval_up)          ← Pack-a-Punch result
2nd upgrade  (e.g. iw8_asval_up_up)       ← Reapy's custom "double-PaP" convention
alt weapon   (e.g. iw8_ak47_launcher_zm)  ← alternative weapon (e.g. under-barrel launcher) the parent's altWeapon field points at
```

Each entry in the family is a **separate GDT record**, separately zoned, with its own ~1300 fields. The relationships between them are wired in two places:

1. **`share/raw/gamedata/weapons/zm/zm_levelcommon_weapons.csv`** — the master CSV mapping each `weapon_name` to its `upgrade_name` (and other meta — cost, ammo cost, mystery-box availability, wonder-weapon flag, `force_attachments`, etc.).
2. **Shared GSC scripts** consume that CSV — the upgrade behaviour itself lives in `share/raw/scripts/zm/_zm_pack_a_punch.gsc` (and `_zm_pack_a_punch_util.gsc`), not in the engine. The engine's role here is small (mostly CSC-side: rendering the PaP weapon list, etc.). When you need behaviour beyond what the shared script provides — e.g. Reapy's `_up_up` second-upgrade chain — you **override** the shared script in your usermap (`usermaps/zm_test/scripts/zm/_zm_pack_a_punch.gsc`).

## Record types

The two main GDT types declaring a weapon — APE picks which one based on the weapon's primary fire behaviour:

- **`bulletweapon`** (`bulletweapon.gdf`) — hitscan: pistols, ARs, SMGs, LMGs, snipers, shotguns. Damage resolved instantly along a ray with `bulletDamage` + `playerDamage` curves.
- **`projectileweapon`** (`projectileweapon.gdf`) — physical projectile: rockets, grenade launchers, the under-barrel launchers, some wonder weapons. Spawns an entity that flies, can detonate, can ricochet.

The wider weapon-family GDT types (per [`reference/asset-types.md`](../reference/asset-types.md) — Weapons & combat group):

- **`dualwieldweapon`** / **`dualwieldprojectileweapon`** — left-hand companion of an akimbo pair (the *right* hand keeps using the regular `bulletweapon` / `projectileweapon` record, just with the `_rdw_*` suffix; see naming conventions).
- **`grenadeweapon`** — frag/semtex/throwables.
- **`meleeweapon`** — bowie, knives, the WW2 PaP melee.
- **`gasweapon`** — chemical gas, e.g. cellbreaker tear-gas grenade.
- **`cybercomweapon`** — multiplayer cybercom abilities (mostly Treyarch-shipped).
- **`turretweapon`** — placeable turrets, sentry-style.
- **`attachment`**, **`attachmentunique`**, **`attachmentcosmeticvariant`** — attachment definitions and the per-attachment cosmetic variants the weapon record references via its `acv_*` slots.

## Naming conventions

You'll see these suffixes everywhere in `usermaps/zm_test/zone_source/zm_test_weapons.zpkg`:

| Suffix          | Meaning                                                                           | Where it's wired                                                              |
| --------------- | --------------------------------------------------------------------------------- | ------------------------------------------------------------------------------ |
| (none)          | Base weapon                                                                       | Native — referenced everywhere                                                 |
| `_up`           | Pack-a-Punched version                                                            | `zm_levelcommon_weapons.csv` `upgrade_name` column → consumed by `_zm_pack_a_punch.gsc` |
| `_up_up`        | **Reapy's custom second-PaP** (not a Treyarch convention)                         | Reapy's override of `_zm_pack_a_punch.gsc` — adds `second_upgrade_cost` flow  |
| `_zm`           | Zombies-mode-specific variant of an MP/SP weapon                                  | Convention only — referenced by name in the CSV/zone                           |
| `_launcher_zm`  | Alternative weapon (e.g. under-barrel launcher) as a separate switchable weapon   | The parent weapon's `altWeapon` field points at it                             |
| `_rdw_up`       | **Right-hand Dual Wield** upgrade — the right-hand half of a dual-wield pair. Uses the regular `bulletweapon` / `projectileweapon` GDT type, but configured to participate in the dual-wield pairing. | Convention; any dual-wield-capable weapon (not just pistols)               |
| `_ldw_up`       | **Left-hand Dual Wield** companion — the left-hand half of the same pair. Uses the dedicated `dualwieldweapon` / `dualwieldprojectileweapon` GDT type. The right-hand record's `DualWieldWeapon` field points at this `_ldw_*` record by name. | Convention; pairs with `_rdw_*`                                            |

> The `_up_up` chain is one of Reapy's signature features. Since Treyarch's PaP system bottoms out after one upgrade, doubling required custom GSC: detect the player has already PaP'd, swap from `weapon_up` to `weapon_up_up` on a second interaction. Worth documenting in `docs/scripting/` when we get there.

## The weapon GDT record (~1300 fields)

You will not edit most of these. You will read a few, and tweak a handful when something's wrong. Here's how to navigate.

### Field groups, by "how often does Reapy actually touch this?"

| Frequency        | Group                                  | Fields you might genuinely care about                                       |
| ---------------- | -------------------------------------- | ---------------------------------------------------------------------------- |
| **Often**        | Damage & ranges                        | `bulletDamage`, `minDamage`, `playerDamage`, `minPlayerDamage`, `maxDamageRange`, `minDamageRange`, `meleeDamage` |
| **Often**        | Ammo & reload                          | `clipSize`, `maxAmmo`, `startAmmo`, `reloadTime`, `reloadEmptyTime`, `reloadEmptyAddTime` |
| **Often**        | Recoil & sway                          | `recoilShotsViewKick*` family, `swayMaxAngle`, `swayShootingScalars` — Reapy tuned these heavily for the ported feel |
| **Sometimes**    | Fire behavior                          | `fireType` (full-auto / semi / burst), `fireTime`, `burstCount`, `cookTime`  |
| **Sometimes**    | Models                                 | `model` (viewmodel), `worldModel`, `hideTags`                                |
| **Sometimes**    | Sounds                                 | `fireSound`, `fireSoundPlayer`, `pickupSound`, `reloadSound` (and ~50 variants) |
| **Sometimes**    | FX                                     | `flashEffect`, `flashEffectPlayer`, `viewFlashEffect`, `tracer`              |
| **Rarely**       | ADS / movement                         | `adsTime`, `adsTransInTime`, `moveSpeedScale`                                |
| **Rarely**       | Attachment cosmetic variant slots      | `acv_*` (~50 fields, one per attachment type — links to `attachmentcosmeticvariant` records) |
| **Almost never** | Engine-internal tuning                 | Hundreds of fields with engine-default values                                |

### Read the GDT, don't edit by hand

For ported weapons, the pack author has already tuned these for the source game's feel. Don't waste time reviewing all 1300 fields — instead:

- If the weapon **shoots**, **reloads**, and **deals damage** in a vague approximation of what you expect, ship it.
- If something's *visibly* wrong (no flash, no sound, wrong damage curve, wrong attachment look), open the GDT in APE, find the relevant field group from the table above, tweak.
- If the weapon **doesn't show up at all**, the problem is in the CSV / zoning / dependencies — see Common Gotchas below.

## The level CSV (`<map>_weapons.csv`)

Every weapon you want available in your map needs a row in your map's weapons CSV. This is what the zombies game-mode scripts read to know about your roster.

> ⚠️ Two CSVs to know about, **only one of which you edit:**
>
> - `share/raw/gamedata/weapons/zm/zm_levelcommon_weapons.csv` — the **shipped** Treyarch reference roster. Read-only for you; useful for column-format reference and stock-weapon examples.
> - `usermaps/<your_map>/gamedata/weapons/zm/<your_map>_weapons.csv` — your **per-map override**. This is the one you actually edit. For `zm_test`: `usermaps/zm_test/gamedata/weapons/zm/zm_test_weapons.csv`.

Real header (same in both files):

```
weapon_name,upgrade_name,hint,cost,weaponVO,weaponVOresp,ammo_cost,create_vox,
obsolete_false,in_box,upgrade_in_box,is_limited,limit,upgrade_limit,obsolete2_false,
wallbuy_autospawn,class,is_aat_exempt,is_wonder_weapon,force_attachments
```

Real rows:

```
t9_1911,t9_1911_rdw_up,,300,pistol,,,,,TRUE,FALSE,FALSE,,,FALSE,TRUE,pistol,,,
iw8_asval,iw8_asval_up,,1350,rifle,,,,,TRUE,FALSE,FALSE,,,FALSE,TRUE,rifle,,,
iw8_vintorez,iw8_vintorez_up,,1600,sniper,,,,,TRUE,FALSE,FALSE,,,FALSE,TRUE,sniper,,,
pistol_standard,pistol_standard_upgraded,,50,pistol,,0,,,FALSE,FALSE,TRUE,0,,FALSE,FALSE,pistol,TRUE,,
```

**Key columns:**

- **`weapon_name`** / **`upgrade_name`** — the PaP wiring. Engine reads `upgrade_name` when the player Pack-a-Punches the active weapon.
- **`cost`** — wallbuy cost (not used if the weapon isn't on a wallbuy).
- **`ammo_cost`** — refill cost from the wallbuy.
- **`weaponVO`** — voice-line category (`pistol`, `rifle`, `sniper`, `shotgun`, `smg`, `lmg`, `launcher`, …); drives "out of ammo" / pickup VO.
- **`in_box`** — `TRUE` if the weapon can come out of the Mystery Box.
- **`upgrade_in_box`** — `TRUE` if the **upgraded** version can come out of the box (e.g. some Wonder Weapons).
- **`is_limited`** / **`limit`** — limited-quantity weapons (e.g. only N of them ever exist on the map).
- **`wallbuy_autospawn`** — auto-create wallbuy entities for this weapon.
- **`class`** — used by perks like Mule Kick to slot weapons by category.
- **`is_aat_exempt`** — exempt from Alternate Ammo Types (which AAT-style perks/effects can't apply to it).
- **`is_wonder_weapon`** — flagged as a wonder weapon (drives box appearance rate, point-cost differences, certain script paths).
- **`force_attachments`** — space-separated list of attachments the weapon must always have. Examples from the shipped CSV: `ar_famas` → `reddot grip`, `ar_garand` → `reflex rf`, `ar_peacekeeper` → `holo damage fmj`. **In `zm_test` this column is always empty** — Reapy attaches things via the weapon's GDT directly instead (e.g. commit `a7a9e04e` "stg44 scope change" is a GDT-only edit, no CSV touched). Use the GDT route if you want per-weapon attachment defaults you can tune in APE; use this column if you want a CSV-level override.

> ⚠️ **Forgetting to add a row to this CSV is the most common reason a "fully zoned" weapon doesn't show up in-game**, or shows up but can't be Pack-a-Punched. The CSV is the contract between your GDT records and the zombies scripts.

## Camos: `weaponcamotable` + `weaponcamo`

Reapy's deepest non-trivial weapons-side work. The system has two layers:

### Layer 1 — `weaponcamotable` (the registry)

A **table of camo tables**. Lists up to 10 sub-`weaponcamo` records by name. Real example from `source_data/skye_up_camo.gdt`:

```
"skye_up_camo" ( "weaponcamotable.gdf" )
{
    "configstringFileType" "WEAPONCAMOTABLE"
    "numCamoTables"        "5"
    "table_01_name"        "skye_up_camo_dlc1"
    "table_02_name"        "skye_up_camo_base76"
    "table_03_name"        "skye_up_camo_base121"
    "table_04_name"        "skye_up_camo_base128"
    "table_05_name"        "skye_up_camo_your_own_base136"
    "table_06_name"        ""
    ...  (table_07_name through table_10_name empty)
}
```

So this `weaponcamotable` is a *manifest* that tells the engine "here are 5 camos in this table, by name." You zone the table; the engine resolves the named sub-records.

### Layer 2 — `weaponcamo` (a single camo definition)

Each `weaponcamo` record describes how a single camo paints onto the weapon's individual material slots. The naming pattern in the GDT is `material1_<slot>_<property>`:

```
"skye_up_camo_base121" ( "weaponcamo.gdf" )
{
    "baseIndex" "121"
    "configstringFileType" "WEAPONCAMO"
    "material1_10_base_material_1" "skye_up_camo"
    "material1_10_base_material_2" "mtl_wpn_t7_loot_ar_m16_handguard_base"
    "material1_10_base_material_3" "mtl_wpn_t7_loot_ar_m16_sights_base"
    ...
    "material1_10_camo_mask_2"     "i_wpn_t7_loot_ar_m16_handguard_base_r"
    "material1_10_camo_mask_3"     "i_wpn_t7_loot_ar_m16_sights_base_r"
    ...
    "material1_10_material"        "mtl_wpn_t7_camo_dlc3_pap_base"
    "material1_10_scale_x"         "0.5"
    "material1_10_scale_y"         "1"
    "material1_10_useGlossMap"     "1"
    "material1_10_useNormalMap"    "1"

    "material1_11_..."  (next slot, same field shape)
    ...
}
```

Reading the prefix: `material1_<N>_*` defines slot N's substitution rules. For each slot you specify:

- **`base_material_1..10`** — the *original* materials this slot might use on the weapon (a single weapon material slot can resolve to multiple base materials depending on attachments / variants — that's why there are 10 slots).
- **`camo_mask_1..10`** — the mask images that drive how camo paint applies onto each base material (alpha-style masks).
- **`material`** — the *replacement* camo material to swap in.
- **`scale_x` / `scale_y` / `rotation` / `trans_x` / `trans_y`** — UV transform on the camo overlay.
- **`useGlossMap` / `useNormalMap` / `gloss_blend` / `normal_amount`** — whether the camo brings its own gloss/normal or inherits the underlying material's.
- **`detail_normal_*`** — detail-normal-map overlay on top.

`baseIndex` is the engine-side numeric ID this camo claims (e.g. `121`). Camos collide if two records share a baseIndex, so packs co-ordinate ranges (the `base76 / base121 / base128 / base136` naming reflects index allocations).

### Why this is non-trivial

A camo isn't "swap one material" — it's "swap N materials, each potentially with different masks and UV transforms, across all slot variants the weapon can express." Authoring a new camo from scratch means walking every weapon-material-slot in your roster and authoring a `material1_<slot>_*` block per slot. Pack authors do this work once for their weapon ports; you mostly just zone their `weaponcamotable` and reference it.

## Worked example: zoning a ported weapon

> Always **follow the README inside the pack zip first** — SkyeLord (and most pack authors) ship explicit step-by-step install instructions per weapon. The general shape below mirrors a typical SkyeLord README; deviate only when the README tells you to.

Real SkyeLord install pattern (Maddox RFB shown — applies to most ports):

**Step 1 — Drop the folders into the BO3 root.** Pack ships pre-laid-out folders (`model_export/skye_ports/...`, `texture_assets/skye_ports/...`, `source_data/...`, etc.). Drag them into `D:/.../Call of Duty Black Ops III/`. They merge into the existing search roots.

**Step 2 — Zone the weapon family** in `usermaps/zm_test/zone_source/zm_test_weapons.zpkg`:

```
weapon,t8_maddox_rfb
weapon,t8_maddox_rfb_up
```

Plus a `_up_up` line if you're doing the second-upgrade thing.

**Step 3 — Add a CSV row** to `share/raw/gamedata/weapons/zm/zm_levelcommon_weapons.csv` (or any per-level weapon table):

```
t8_maddox_rfb,t8_maddox_rfb_up,,1400,rifle,,,,,TRUE,FALSE,FALSE,,,FALSE,TRUE,rifle,,,
```

**Step 4 — Add sound aliases** to `share/raw/sound/aliases/user_aliases.csv` (or another sound alias CSV). The pack's README dumps the exact lines to paste — fire/foley/reload sounds for the base + PaP weapon. Example shape:

```
#BO4 Maddox RFB,,,...
wpn_t8_maddoxrfb_shot_plr,,,skye_ports\t8_maddoxrfb\fire\wpn_t8_maddoxrfb_shot.wav,,,UIN_MOD,,,,,BUS_FX,,,,,,85,85,,,,,,,,,,,,,,,,,,,2d,,,NONLOOPING,...
wpn_t8_maddoxrfb_pap_shot_plr,,,skye_ports\t8_maddoxrfb\fire\wpn_t8_maddoxrfb_pap_shot.wav,,,...
... bolt, mag in/out, foley etc.
```

**Step 5 (optional, for inspect)** — register the weapon for the inspect animation system. In your map's GSC, before `zm::usermap_main();`:

```c
inspectable::add_inspectable_weapon( GetWeapon("t8_maddox_rfb"), 5 );
inspectable::add_inspectable_weapon( GetWeapon("t8_maddox_rfb_up"), 5 );
```

**Step 6 — Place a wallbuy** in Radiant. SkyeLord packs ship prebuilt prefabs; open the prefab browser → `zm/skye_prefabs/`, drag the matching wallbuy prefab into your map, line it up against a wall.

**Step 7 — Link in the modtools launcher → Run.** Most weapon installs are link-only (no Compile/Light needed) since you didn't change BSP geometry. *Exception*: if you placed a wallbuy prefab (Step 6 above), that's a Radiant geometry change → you'll need a full Compile (or "Just Ents" if you only added entities) before Link.

### When you'll have to do more than the README says

- **ACV (`acv_*`) attachment cosmetics**: most pack-authored weapons pull their `attachmentcosmeticvariant` records in as dependencies; you only have to manually zone them if your `assetinfo/zm_test.csv` shows they didn't land in the fastfile. *Tweak Reapy actually had to do for some weapons.*
- **Double-PaP (`_up_up`)**: extend your CSV with a chained mapping and add the `_up_up` GDT record (typically copied from `_up` and re-tuned).
- **Camo support for the upgraded weapon**: if PaP applies a camo, ensure the weapon's materials are in your `weaponcamotable`'s `base_material_*` slots (see camo section + the t7wiki references below).

### Reference reading (pack, PaP, and camo authoring)

Weapon ports and integration:
- [SkyeLord's master hub on UGX-Mods](https://www.ugx-mods.com/forum/full-weapons/84/skyes-weapon-ports-to-bo3-master-hub/16874/) — the canonical place for SkyeLord weapon port READMEs and updates. Worth scraping when you want install fidelity.
- [Setting Up Mystery Box Weapons (Modme Wiki)](https://wiki.modme.co/wiki/black_ops_3/basics/Setting-up-mystery-box-weapons.html) — the box rotation side of weapon integration.
- [NevisX](https://github.com/Scobalula/NevisX) (Scobalula) — **live in-game weapon stat editor**. Watches your GDTs on disk for changes; when you save in APE (or any other editor), NevisX pushes the updated weapon definition into the running game process — **no Link, no relaunch**. Auto re-attaches across game instances so you can leave it running. Massive iteration-speed boost when tuning recoil / state timers / damage / ammo.
  > ⚠️ NevisX only syncs **simple float/number fields** (~99% of weapon settings according to the README). Model swaps, material swaps, sound aliases, FX bindings — anything non-numeric — still need a full Link cycle. Use it for *tuning*, not for *swapping*.
- [How To Make Custom Pack-a-Punch (YouTube)](https://www.youtube.com/watch?v=EBzzZAa3D90) — covers the PaP machine setup side that pairs with the GDT/CSV work documented above.

Camo authoring:
- [t7wiki — Change Pack-a-Punch camo](https://www.t7wiki.com/guides/change-pack-a-punch-camo)
- [t7wiki — Add camo table to custom weapons](https://www.t7wiki.com/guides/add-camo-table-to-custom-weapons)
- [t7wiki — Increase camo count](https://www.t7wiki.com/guides/increase-camo-count)

## Common gotchas

- **Map closes immediately on start with `zm_weapons.gsc:0` error.** Almost always means your map has a wallbuy entity referencing a weapon that's *not actually included* in the map (missing zone entry, missing CSV row, or wallbuy `script_string` typo). Check Radiant for orphaned wallbuy prefabs.
- **Weapon doesn't appear in the Mystery Box.** Missing CSV row, or `in_box=FALSE`. Or the Box rotation script picks from a curated list and you didn't add this weapon to it.
- **PaP returns a "broken" weapon (bad name).** Either the `upgrade_name` in the CSV is misspelled, or the `_up` GDT record exists but isn't zoned.
- **Camo is missing on the upgraded weapon only.** PaP applies a camo via the `weaponcamotable` you've configured for the level — verify your `_up` weapon's materials are listed in the camo's `base_material_*` slots, otherwise the camo can't bind to them.
- **Wallbuy crashes or silently breaks the map when Pack-a-Punching.** Hard wallbuy limit (~20 wall weapons) per the README. Reduce wallbuy count or move some weapons to box-only.
- **Weapon fires but no flash / no sound / no tracer.** Asset chain broken: `flashEffect` / `fireSound` / `tracer` reference assets that didn't get zoned. Check `assetinfo/zm_test.csv` for missing dependencies.

## Where to go next

- **Custom upgrade chains** (`_up_up` script-side wiring, hellround camo rewards on win) → `docs/scripting/` (forthcoming)
- **Community packs and dupe handling** → [`07-community-packs.md`](./07-community-packs.md)
- **Sounds for weapons** (sound aliases, `sharedweaponsounds`) → covered when we tackle audio (research-driven section)

---

## Open questions / TODO

- [ ] Build a small dedicated `attachmentcosmeticvariant` reference page enumerating which `acv_<slot>` field on a weapon corresponds to which attachment type, plus what the `_smg_standard` / `_ar_standard` / `_pistol_standard` naming convention encodes. Useful when troubleshooting "the suppressor is floating off the muzzle."
