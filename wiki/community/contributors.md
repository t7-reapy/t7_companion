# Contributors registry

> **Status: stub.** Per-author breakdown of community contributors whose packs ship in `zm_test`. The asset side of [`docs/community/`](./README.md). Eventually a denser per-author table covering: what they ship, canonical download links, last-known-good URLs, common quirks / version drift, and the right way to credit.

## What this page will eventually contain

For each author:
- Pack(s) shipped, by name and asset category.
- Canonical download links (with both original source and self-host backup paths).
- Discord / website / GitHub presence.
- Quirks: install gotchas, common collision targets, deprecation notes.
- Credit string format the author prefers.

## Top-of-mind authors (loose grouping until the registry is fleshed out)

### Umbrella / foundational
- **Midgetblaster** — base v2.5 pack, one of the most-installed foundational packs.
- **harrybo21** — perks (Ultimate Perk Pack), gun pack, FX library, physics presets, BT assets, napalm/parasites/apothicon furies AI.
- **Scobalula** — tools (NevisX, Harmony, Cerberus, Greyhound, HydraX) + assets (blood splatter, MWR player models, on-screen rain drops).

### Weapon ports
- **SkyeLord** — the canonical CoD-port hub: BO2/BO4/CW/MW/MW2/MW3/IW/AW/Vanguard/WW2/Ghosts.
- **MADGAZ** — BO6 textures/materials, vehicles/foliage/gore/Liberty Falls assets, PaP camo, knives.
- **Pmr360** — knives.
- **Ma7** — IW Shredder.

### AI / character
- **NSZ / Spiki** — NSZ Brutus, empty bottle powerup.
- **Lednor (MrLednor)** — wolf heads.
- **Logical** — charred zombies, player models.
- **HarryBo21** — apothicon furies, parasites, napalm zombies (also under "Foundational").
- **ninjamanny829** — Origins zombies.

### Sound / VOX
- **Werelupus** — Cold War Zombie VOX.
- **SaintVertigo** — RE2 Remake zombie VOX, dynamic weapon movements.
- **westchief596** — Kortifex announcer pack.
- **Uk_ViiPeR** — MW2022 campaign sound assets.
- **WetEgg** — death animations, Ultimate Round Sounds.
- **Combat** — original perk jingle instrumentals.
- **Starzismik** — BO3 sound assets fixes (ZXG hosting).

### Maps / cinematic
- **KingslayerKyle** — character models (MW19 TF141, MWR USMC), WW2 PaP, MWR HUD, T7LuaRepo.
- **MoiCestTOM** — Doom assets, blood GFX.
- **Symbo** — vehicle scripter, Apothicon Glaive sword powerup.
- **Vertasea** — animated power switch, dog traversals, perk bump audio.
- **Topato** — MWR frag grenades.

### FX / VFX
- **coolyer** — BO1 numbers FX.
- **Mike Pence (modme Dick_Nixon)** — BO2 MOTD lightning FX.
- **Program115** — rain sounds.
- **VerK0** — thunder & lightning (with Program115).
- **JBird** — traversals.

### Distribution helpers / engine support
- **JariK / JariKCoding** — 1-byte game data persistence script (the save system); the working T7Overcharged fork.
- **JxstNoTex** — original T7Overcharged fork (often pointed at first; the JariKCoding fork is what actually works).
- **shiversoftdev** — `t7patch` (engine fixes), `t7-source` (Treyarch source dump), BO3Enhanced.
- **dest1yo / echo000** — Saluki (successor to Greyhound).
- **DTZxPorter** — porter-lib (underpins Saluki), Cast format, ModmeWiki original maintainer.
- **ATE47** — atian-cod-tools (`acts`) for on-disk `.ff` parsing.
- **shidouri** — original GDTDupePurgerPy that became Reapy's `dupe_fixer.py`.

### Tools (script + UI side)
- **blakintosh** — GSCode VS Code extension.
- **ardivee** — Acoustix (sound editor), Cyph3r (localization editor), LuiGUI (PSD→LUI), Zoroth.

### Honourable / unaffiliated
- **Keysia** — original drawings/illustrations used in `zm_test`.
- **xoxor4d** — Flashlight (with Program115).
- **Holofya** — Flashlight with UV (with xoxor4d).
- **HiMMA** — Furniture desk pack.
- **Surfypolecat4** — The Trinket Box.
- **Kimday** — WinRar license.
- **OwenC137** — Blast-O-Matic.
- **Robit** — BO4 wallbuy texture.
- **LG-RZ** — BO3 shaders.
- **Fearlessninja98** — Perk Poster Challenge.
- **Sphynx** — dev commands utility scripts.
- **XcDylan93** — weapon camo utility scripts (lives in repo).
- **lilrobot** — inspectable weapons.
- **Symbo, Pmr360, et al.** — see weapon-ports / map-cinematic groupings above.

## How this page should evolve

When filling in: prefer per-author *one-screen profiles* over flat lists. Each profile = pack name → canonical URL → backup path → install notes → known-bad-versions. Mirror the README dependency table where possible so this page doesn't drift from the lockfile.
