# Community packs & dupe handling

> The "third-party libraries" of BO3 modding. ~50+ packs are installed in `zm_test`, ranging from giant umbrella collections (Midgetblaster v2.5) to single-feature drops (Vertasea's animated power switch). This page is about *living with them at scale*: install conventions, the dupe problem they create, the purger tooling that survives in the repo, and the ecology of who-ships-what.

## Mental model

Treat community packs as **transitive dependencies in an unmanaged package system**:

- No central registry. No semver. No changelog discipline (you read the commits, if there are any).
- Distribution is fragmented across Mega / Google Drive / MediaFire / iCloud / Discord CDN attachments / GitHub releases / DevRaw / individual creator websites — links rot constantly. The README's dependency table exists *as a lockfile* because of this.
- Installing a pack adds entries to your **global GDT asset registry** and drops raw files under your search roots. There's no dependency *resolution* — you do that work yourself when something collides.

The spectrum (covered briefly in `01-overview.md`, restated here for context):

| Pack shape          | Example                              | What it ships                                                                  |
| ------------------- | ------------------------------------ | ------------------------------------------------------------------------------ |
| **Umbrella**        | Midgetblaster v2.5, harrybo21 gun pack | Thousands of GDT records + raw assets. You install once, cherry-pick by zoning. |
| **Targeted feature** | Vertasea animated power switch, NSZ Brutus | One feature: a `.zpkg` + the GDTs it needs. You `include` the zpkg from your zone. |
| **Hybrid**          | SkyeLord weapon ports, MADGAZ packs   | Some pre-bundled `.zpkg`s + a lot of loose GDTs for cherry-picking.            |

## Install conventions on disk

Convention that emerged in `zm_test` (and that's broadly compatible with how packs ship):

```
_custom/<author>/                 # raw assets (xmodels, textures, sounds) from the pack
source_data/_<author>/            # installed-pack GDTs grouped by author
source_data/custom/               # YOUR own GDTs and tweaked ones (house.gdt, etc.)
source_data/<pack-name>.gdt       # single-file packs that don't warrant a folder
```

Real examples from `zm_test`:

- `_custom/kingslayer_kyle/`, `_custom/wetegg/`, `_custom/_coolyer/`, `_custom/_moicesttom/` — raw assets per author.
- `source_data/_midgetblaster/`, `source_data/_charred_zombies.gdt`, `source_data/skye_*.gdt` — installed pack GDTs.
- `source_data/custom/house.gdt` — Reapy's own fog/SSI/lighting records.

> The `_custom/` folder is **not** a default linker search root — it's listed explicitly in `bin/converter_gdt_dirs_0.txt` (per `01-overview.md`). Without that line, none of these raw assets resolve. Worth checking when something's "missing" after a fresh install.

## Install flow

For a typical pack drop:

1. **Read the pack's README first.** Skye / harrybo21 / madgaz README files name exact zone lines, CSV rows, and (sometimes) sound aliases to paste. Follow them — your time-saving will lose to surprise behaviour you didn't expect otherwise.
2. **Drop the pack's folders into the BO3 root.** Drag `model_export/`, `texture_assets/`, `_custom/`, `source_data/` etc. — they merge with the existing roots.
3. **Drop the pack's GDTs into `source_data/`** (or wherever the pack expects). APE indexes anything in there.
4. **Fix paths inside the GDT** if the pack assumes a different layout from yours. Most pack-supplied GDTs reference assets via paths like `_custom/<author>/...` or `texture_assets/<author>/...` — check they resolve under your `converter_gdt_dirs_0.txt` roots.
5. **Add zone entries** for the assets you actually want (per the pack README). Either `include,<pack_zpkg>` or individual `<type>,<name>` lines.
6. **Add CSV rows** for any weapons / per-level data the pack wants you to register (`zm_levelcommon_weapons.csv`, sound alias CSVs, etc.).
7. **Link → Run.** Check the linker output for the inevitable dupes (next section).

## The dupe problem

The single biggest *recurring* friction with community packs is **asset-name collision**. Two packs declare an `image`/`material`/`xmodel`/whatever asset under the same name → the linker silently picks one (usually the last GDT loaded), the other is shadowed, and you find out only when the wrong one shows up in-game.

You see it in the linker output as lines like:

```
ERROR: Duplicate 'image' asset 'mtl_t7_world_lit_concrete_d' found in source_data/_midgetblaster/t7_props.gdt:12345
```

These accumulate by the hundreds when you stack umbrella packs. You can't ignore them because *which one wins* is non-deterministic relative to your install order.

## Dupe-purger tooling in the repo

Two scripts ship in the repo to deal with this:

### `dupe_fixer.py`

Reapy's modified version of **Shidouri's [GDTDupePurgerPy](https://github.com/shidouri/GDTDupePurgerPy)**. Walks all GDTs in `source_data/`, finds duplicate asset names, and removes the redundant blocks from chosen GDTs. Modifications: (a) Midgetblaster GDTs are preferred deletion targets, so collisions favour keeping the more-specific pack's copy, (b) optional flags for printing / backup / verbose output.

**Invocation:**

```bash
python dupe_fixer.py [flags]
```

**Inputs the script reads from its working directory** (the BO3 install root, where the script lives):

| File             | Role                                                                                       |
| ---------------- | ------------------------------------------------------------------------------------------ |
| `dupe_error.txt` | The linker error log. Auto-created on first run if missing — paste your linker output here. |
| `midget.gdtdef`  | List of Midgetblaster GDT relative paths to prefer-purge.                                 |
| `stock.gdtdef`   | List of stock GDTs that must NOT be touched (defensive).                                  |

**Flags** (any of these spellings, with `+`/`-`/`/` prefix — all stripped before matching):

- Quiet output: `noshow`, `quiet`, `shh`, `no_print`, `no_log`
- Skip backup (dangerous): `nobak`, `developer_no_backup_use_wisely`
- Verbose: `verbose`, `v`, `logs`, `log`

Default behaviour: print on, backup on, non-verbose. Backups go alongside the modified GDT.

**Workflow:**

1. Run a Link in modtools.
2. Capture the linker output to `./dupe_error.txt` (the script reads only this file).
3. Make sure `midget.gdtdef` + `stock.gdtdef` are present (they live in repo).
4. Run `python dupe_fixer.py` (add `-verbose` if you want the per-asset log).
5. Re-Link; collisions should drop. Repeat if new ones surface.

Why Midgetblaster-priority: Midgetblaster v2.5 is the largest umbrella pack — it ships *thousands* of generic asset records. When a more-specific pack (SkyeLord port, harrybo21 perk variant, etc.) collides with Midgetblaster, the safer "deletion target" is almost always the Midgetblaster copy, because the specific pack was authored for a particular use case while Midgetblaster's record is a generic catch-all.

### `midgetblaster_duplicates/duplicates.py`

Companion script for the targeted case: scrape an existing `duplicates.txt` linker-error log, grep only the `_midgetblaster`-path lines, and surgically remove just those asset blocks from Midgetblaster GDTs. Use this when you know you only want to deduplicate against Midgetblaster's contributions specifically (e.g. after adding a new pack and seeing the new collisions in the error output).

**Invocation:**

```bash
cd midgetblaster_duplicates
python duplicates.py
```

Reads `./duplicates.txt` (drop the linker output here), edits Midgetblaster GDTs in-place. Pattern it greps:

```
ERROR: Duplicate '<TYPE>' asset '<NAME>' found in <PATH containing _midgetblaster>:
```

There's also a sibling **`duplicate_parent.py`** in the same folder for the related `GDT ParseError: Parent Entity does not exist` errors (different failure mode, same Midgetblaster-priority approach).

### When to run

- **First install of a pack**: run after the first failing Link. Typically generates new collisions.
- **After tweaking your own GDTs**: only if you edited GDTs *manually* (text-edited the file). APE-driven edits don't introduce duplicates because APE manages asset uniqueness for you — only hand-edits can drift.
- **After updating the modtools or re-installing assets**: defensively.

The workflow is fully manual — run the script, re-Link, repeat until the collision count stabilises (or until the remaining ones are intentional).

## Community-pack ecology

The major authors you'll repeatedly install from. Roughly grouped by what they ship.

### Umbrella / foundational
- **Midgetblaster (v2.5)** — the giant base pack. Many other packs assume Midgetblaster is installed (it provides shared dependencies).
- **harrybo21** — perks (`Ultimate Perk Pack`), gun pack, FX library, physics presets, BT assets, napalm/parasites/apothicon furies AI. Foundational for many ZM features.
- **Scobalula** — utility-end of the ecosystem: tools (NevisX, Harmony, Cerberus, Greyhound, HydraX), assets (blood splatter, MWR player models, Bo3OnScreenRainDrops, [Bo3Mutators](https://github.com/Scobalula/Bo3Mutators) — a Lua/GSC mutator framework).

### Weapon-port specialists
- **SkyeLord** — the canonical CoD-port hub: BO2/BO4/CW/MW/MW2/MW3/IW/AW/Vanguard/WW2/Ghosts weapons. Almost the entire `zm_test` weapon roster.
- **MADGAZ** — BO6 texture / material packs, BO6 vehicle/foliage/gore/Liberty Falls assets, PaP camo, knife packs.
- **Pmr360** — knives.

### AI / character / specialised entities
- **NSZ / Spiki** — NSZ Brutus, empty bottle powerup.
- **Lednor (MrLednor)** — wolf heads, dog/wolf assets.
- **Logical** — charred zombies, player models.

### Sound / VOX
- **Werelupus** — Cold War Zombie VOX.
- **SaintVertigo** — RE2 Remake zombie VOX (also ships dynamic weapon movements outside the VOX bucket).
- **westchief596** — Kortifex announcer pack.
- **Uk_ViiPeR** — MW2022 campaign sound assets.
- **WetEgg** — death animations, ultimate round sounds.

### Maps / cinematic
- **KingslayerKyle** — character models, MWR USMC, WW2 PaP, MWR HUD, T7LuaRepo.
- **MoiCestTOM** — Doom assets, blood GFX.
- **Symbo** — vehicle scripter, Apothicon Glaive sword powerup.
- **Vertasea** — animated power switch, dog traversals, perk bump audio.

### Distribution helpers
- **JariK** — 1-byte game-data persistence script (the save system).
- **shiversoftdev** — `t7patch` (engine fixes), `t7-source` (the Treyarch source dump).
- **JxstNoTex** / **JariKCoding** — T7Overcharged forks (the working fork is JariKCoding).
- **dest1yo / echo000** — Saluki successor to Greyhound.

(See [`docs/community/contributors.md`](../community/contributors.md) for the running registry.)

## Distribution chaos & link rot

Where packs actually live on the internet, by frequency:

- **Mega** — most common. Often the *only* link, frequently goes down.
- **Google Drive** — Skye's hub uses iCloud + Google Drive heavily; iCloud links sometimes 404 for non-Apple users.
- **MediaFire** — older WaW-era convention, still used.
- **Discord CDN attachments** — fragile (links rot when posts are deleted, attachments expire). Vertasea's assets are mostly here.
- **GitHub releases** — the most reliable. Scobalula's tools all live here.
- **Personal websites** — KingslayerKyle, MrLednor, Airyz host on their own sites.
- **DevRaw + DevRaw Database Legacy spreadsheet** — community-curated index trying to track what's where.

> **Mirror what you use.** Reapy's working pattern: the README dependency table points at the *external* source so other people can find it from the same place he originally did. The HDD copy of `bo3 tools - sources/` is a personal backup for when those external links die — not an alternate distribution channel for others.

## Common gotchas

- **Linker says "could not find asset" after a fresh pack install.** Either the pack's GDT references a path under a search root not in `converter_gdt_dirs_0.txt` (often `_custom/`), or you copied the GDT but not the raw assets it references.
- **A texture/model/weapon looks wrong after a pack install.** Almost always silent shadowing — another pack already shipped that asset name and yours got deduplicated to the wrong winner. Run the dupe-purger; check the linker output for the explicit conflict lines.
- **Removing a pack leaves orphan references.** GDT records can outlive the assets they reference. Audit `source_data/` for entries pointing at deleted folders before you re-Link.
- **Pack ships its own modified shared script.** Some packs replace `share/raw/scripts/zm/_zm_*.gsc` files. This collides with *your* override directory. Treat shared-script replacements as breaking changes — copy them into your `usermaps/zm_test/scripts/zm/` directory and merge by hand.
- **Pack's README is wrong / out of date.** Common with old packs whose author moved on. Cross-reference the modme thread or the UGX hub for corrections; ask in Discord.

## Where to go next

- **Asset taxonomy**: full GDT type list in [`docs/reference/asset-types.md`](../reference/asset-types.md).
- **Author registry**: [`docs/community/contributors.md`](../community/contributors.md) — per-author breakdown of what ships, where to download, common version drift.
- **Discord servers** / **wikis**: [`docs/community/`](../community/) for the index of where to ask when a pack misbehaves.

## Reference reading

- [Sphynx's Mega-Release Thread (Modme Forums)](https://forum.modme.co/wiki/threads/3031.html) — craftables, dev commands, scripts; one of the densest single-author release threads.
- [NSZ_POWERUPS_MEGATHREAD (Modme Forums)](https://forum.modme.co/wiki/threads/2831.html).
- [Mike's repertoire of assets (Modme Forums)](https://forum.modme.co/wiki/threads/2402.html).
- [DOWNLOADS — MRLEDNOR](https://mrlednor.wixsite.com/mrlednor/downloads).
- [KingslayerKyle's Assets](https://sites.google.com/view/kingslayerkyle).
- [DEVRAW Database Legacy spreadsheet](https://docs.google.com/spreadsheets/d/10aQLnuZUgvduFS4zgPNOTBlFbD-pBfpy--Gm9ilIqRg/edit?pli=1&gid=1848896341#gid=1848896341) — the closest thing to a community-curated asset index.
- [Sharp's BO3 Modtools Assets forum](https://sg4yforums.com/Assets/BO3-Modtools-Assets/).
- [shidouri/GDTDupePurgerPy](https://github.com/shidouri/GDTDupePurgerPy) — the upstream of the in-repo `dupe_fixer.py`.

---

## Open questions / TODO

- [ ] Build out [`docs/community/contributors.md`](../community/contributors.md) into per-author *one-screen profiles* (pack list, canonical URL, backup path, install notes, known-bad versions). The author registry currently exists as a stub with the loose grouping; needs density.
- [ ] Write up the "shared script replacement" merge workflow (when a pack ships a modified `_zm_*.gsc`) — belongs in [`docs/scripting/`](../scripting/), to be done when the scripting section is written.
