# GDT asset types — full registry

The complete list of **GDT (engine) asset types** that APE exposes for BO3, captured from the "New Asset" type-selector dropdown. These are the things you can author / edit as records in a `.gdt`.

> ⚠️ Note: this list is the *engine's* asset taxonomy, **not** a 1:1 list of valid `.zone` / `.zpkg` line types. Zone syntax has its own dispatch — many GDT records are pulled in transitively (no zone line needed) or via umbrella aliases (e.g. weapon records of type `bulletweapon` / `projectileweapon` are zoned with the umbrella `weapon,<name>` line). Some entries here may not have a corresponding zone-line type at all. Treat the list below as "what APE can edit," and the section further down on **non-APE zone types** as the orthogonal "what zone files actually accept."

Use this page as the source of truth when sanity-checking _"is `xyz` actually a valid GDT asset type?"_.

> Some types are shipped/Treyarch-only and won't matter for community modding (multiplayer body styles, training-sim ratings, accolades, etc.). Others are workhorses (`image`, `material`, `xmodel`, `weapon`, `fx`-related…). The list below is alphabetical; the **"Most-touched in this project"** section at the bottom calls out which Reapy actually messes with daily on `zm_test`.

## Full list (alphabetical) — GDT/engine types editable in APE

```
accolade
accoladelist
aiassassination
aifxtable
aimtable
ainames
aitype
attachment
attachmentcosmeticvariant
attachmentunique
attackables
beam
bonuszmdata
botsettings
bullet_penetration
bulletweapon
cgmediatable
character
charactercustomizationtable
characterweaponcustomsettings
codfuanims
collectible
collectiblelist
containers
customizationcolor
cybercomweapon
destructiblecharacterdef
destructibledef
destructiblepiece
doors
dualwieldprojectileweapon
dualwieldweapon
duprenderbundle
emblem
entityfximpacts
entityfxtable
entitysoundimpacts
flametable
fog
footsteptable
fxcharacterdef
gallery_image
gallery_imagelist
gamedifficulty
gasweapon
gibcharacterdef
glass
grenadeweapon
image
impactsfxtable
impactsoundstable
killcam
killstreak
laser
lensflare
light
lightdescription
locdmgtable
maptable
maptableentry
maptableloadingimages
material
medal
medalcase
medalcaseentry
medaltable
meleeweapon
mpbody
mpdialog
mpdialog_commander
mpdialog_player
mpdialog_scorestreak
mpdialog_taacom
objective
objectivelist
physconstraints
physpreset
player_character
playerbodystyle
playerbodytype
playerfxtable
playerhead
playerhelmetstyle
playersoundstable
postfxbundle
projectileweapon
ragdollsettings
rumble
sanim
scriptbundle
sentientevents
sharedweaponsounds
shellshock
sitrep
ssi
surfacefxtable
surfacesounddef
tagfx
teamcolorfx
tracer
trainingsimrating
trainingsimratinglist
turretanims
turretweapon
vehicle
vehiclecustomsettings
vehiclefxdef
vehicleriders
vehiclesounddef
weaponcamo
weaponcamotable
xanim
xcam
xmodel
xmodelalias
zbarrier
```

### Types used in zone files but **not** in APE's New Asset dropdown

These are real asset types you'll see referenced as `<type>,<name>` in `.zone` / `.zpkg` files, but APE doesn't expose them under "New Asset" because they originate from a different authoring path (raw files on disk, engine internals, or umbrella aliases).

The complete list as observed in `usermaps/zm_test/zone_source/`:

| Type                   | Source / origin                                                                                                                                                                                                           |
| ---------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `scriptparsetree`      | Compiled GSC/CSC. Declared by pointing at a raw `.gsc` / `.csc` file under `scripts/`.                                                                                                                                    |
| `stringtable`          | CSV string tables.                                                                                                                                                                                                        |
| `structuredtable`      | JSON-based structured tables, e.g. `gamedata/tables/common/objectives.json`.                                                                                                                                              |
| `localize`             | Localized string keys (per-language). Lives in language-specific zone fragments.                                                                                                                                          |
| `ttf`                  | TrueType fonts.                                                                                                                                                                                                           |
| `rawfile` / `lua_file` | Raw text / compiled Lua payloads.                                                                                                                                                                                         |
| `fx`                   | FX exporter assets (`.efx`), e.g. `fx,custom/AI/cellbreaker_death.efx`. Authored in the FX editor, not in APE.                                                                                                            |
| `techset`              | Engine "technique set" — the lower-level shader binding underneath a material, e.g. `mc/lit_emissive_scroll_advanced#e6142445`. The trailing `#hex` is a hash. You almost never hand-author these; the linker emits them. |
| `weapon`               | Umbrella zoning name. The _actual_ GDT record might be `bulletweapon`, `projectileweapon`, `meleeweapon`, etc., but you zone it with `weapon,<name>` and the linker dispatches by looking up the GDT entry's true type.   |
| `customizationtable`   | Character customization table. Probably the same surface as APE's `charactercustomizationtable` (TBC).                                                                                                                    |
| `col_map`, `gfx_map`   | Extracted views of a `.d3dbsp` (see `asset-pipeline/01-overview.md`).                                                                                                                                                     |
| `sound`                | Sound zoning. The zone line `sound,<name>` references a **`.szc` (Sound Zone Config)** file, *not* the per-alias CSVs. The aliases themselves are declared in `share/raw/sound/aliases/*.csv` but they're *referenced from* the `.szc`, which is what zoning actually pulls in. (To be expanded in [`docs/audio/`](../audio/).)                                                                                                                                                              |
| `include`              | Not a real asset type — the `.zpkg` include directive.                                                                                                                                                                    |

> ⚠️ The `weapon` zoning alias is a clear example: a zone line `weapon,iw8_ak47` resolves to a `bulletweapon` GDT record, while `weapon,iw8_ak47_launcher_zm` resolves to a `projectileweapon`. The linker dispatches by reading the GDT record's *actual* type. So zone-line type ≠ GDT type, and several GDT types may share one zone alias.

## Loose categorization

Rough buckets to navigate the list. Boundaries are fuzzy; some types straddle.

### Visuals & rendering

`image`, `material`, `xmodel`, `xmodelalias`, `light`, `lightdescription`, `lensflare`, `fog`, `ssi`, `postfxbundle`, `duprenderbundle`, `cgmediatable`, `tagfx`, `entityfxtable`, `entityfximpacts`, `surfacefxtable`, `playerfxtable`, `vehiclefxdef`, `fxcharacterdef`, `gibcharacterdef`, `teamcolorfx`, `aifxtable`, `gallery_image`, `gallery_imagelist`, `glass`, `tracer`, `beam`, `laser`, `flametable`, `xcam`, `sanim`, `xanim`, `scriptbundle` (animated/scripted scenes — pairs with `xanim`)

### Weapons & combat

`bulletweapon`, `projectileweapon`, `dualwieldweapon`, `dualwieldprojectileweapon`, `grenadeweapon`, `gasweapon`, `meleeweapon`, `cybercomweapon`, `turretweapon`, `attachment`, `attachmentunique`, `attachmentcosmeticvariant`, `weaponcamo`, `weaponcamotable`, `aimtable`, `bullet_penetration`, `impactsfxtable`, `impactsoundstable`, `entitysoundimpacts`, `shellshock`, `rumble`

### AI & character

`aitype`, `aiassassination`, `ainames`, `character`, `characterweaponcustomsettings`, `charactercustomizationtable`, `customizationcolor`, `playerbodystyle`, `playerbodytype`, `playerhead`, `playerhelmetstyle`, `player_character`, `mpbody`, `attackables`, `sentientevents`, `codfuanims`, `turretanims`, `botsettings`

### Audio

`sharedweaponsounds`, `playersoundstable`, `surfacesounddef`, `vehiclesounddef`

### Map / level / world

`zbarrier`, `doors`, `containers`, `destructibledef`, `destructiblecharacterdef`, `destructiblepiece`, `vehicle`, `vehiclecustomsettings`, `vehicleriders`, `maptable`, `maptableentry`, `maptableloadingimages`, `footsteptable`, `physpreset`, `physconstraints`, `ragdollsettings`

### Game-mode metadata (mostly MP / Treyarch-shipped)

`bonuszmdata`, `gamedifficulty`, `objective`, `objectivelist`, `killcam`, `killstreak`, `sitrep`, `medal`, `medalcase`, `medalcaseentry`, `medaltable`, `accolade`, `accoladelist`, `mpdialog`, `mpdialog_commander`, `mpdialog_player`, `mpdialog_scorestreak`, `mpdialog_taacom`, `emblem`, `collectible`, `collectiblelist`, `locdmgtable`, `trainingsimrating`, `trainingsimratinglist`

## Most-touched in this project (`zm_test`)

Reapy's hands-on depth, by his own ranking + cross-referenced against `git log` for completeness. Note: most assets came **pre-extracted via community packs** (Midgetblaster, harrybo21, SkyeLord ports, etc.); only a couple of items were Greyhound-extracted by Reapy himself. So "depth" here means _records edited / debugged in APE_ more than _raw extraction work_.

1. `image`, `material` — daily driver (covered in [`asset-pipeline/02-images-and-materials.md`](../asset-pipeline/02-images-and-materials.md))
2. `xmodel` — heavy use, mostly via porting (covered in [`asset-pipeline/03-models.md`](../asset-pipeline/03-models.md))
3. `fog` — built the per-zone / per-state fog system (covered in [`asset-pipeline/04-fog.md`](../asset-pipeline/04-fog.md))
4. `ssi` — sun & sky info, heavy iteration during weather/hellround work (covered in [`asset-pipeline/05-ssi.md`](../asset-pipeline/05-ssi.md))
5. `projectileweapon`, `bulletweapon`, `attachment`, `attachmentcosmeticvariant` — weapon record shapes + cosmetic-variant tweaks across the ~40-weapon roster (covered in [`asset-pipeline/06-weapons.md`](../asset-pipeline/06-weapons.md))
6. `weaponcamotable` + `weaponcamo` — non-trivial work for camos and hellround camo rewards
7. **`fx`** — extensive practical FX integration: hellround FX, weapon impact FX, parasite/zombie FX, decoration FX, FX count tuning. Belongs in a future FX page or in the upcoming weather/lighting docs.
8. **`xanim`** + **`scriptbundle`** — death animations (BO2 ports), fauna animations, FX-anim props, siege scripted scenes in hellround. Animation work is broader than the depth-list earlier suggested.
9. **`sound`** zoning + sound aliases — substantial *practical* work (fire sounds, steamfire, splitscreen sound bug fixes, environmental triggers, music management, easter-egg sound choreography). Belongs in [`docs/audio/`](../audio/).

> *Audio nuance*: the GDT-side audio types (`sharedweaponsounds`, `surfacesounddef`, `playersoundstable`, `impactsoundstable`, etc.) remain mostly **research territory** — Reapy hasn't gone deep into authoring those records. But practical sound *integration* (aliases, `.szc` zoning, in-game triggers) was extensive. Two different layers of "audio depth."

This ordering should drive what gets a deep-dive page next, not textbook order.
