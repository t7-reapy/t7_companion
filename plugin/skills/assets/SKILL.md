---
name: bo3-assets
description: How to get a model, material/texture, or animation into Black Ops 3 — extracting from other CoD titles (Wraith, Greyhound, Cordycep), the weapon/character porting pipeline (Maya/Blender export → APE compile → materials → anims), rigging custom models, and common export/GDT pitfalls. Use for porting, custom modeling, texturing, or animation work, as distinct from GSC/CSC scripting.
---

# Assets: models, materials, porting, animation

Sourcing here is mixed: raw Discord threads run low reliability (~0.25), but a real chunk of this domain is backed by UGX/ModMe/T7-wiki writeups (~0.70) and a few schema/source-verified references (~0.90 — e.g. collmaps, script bundles). Check the `reliability` score per hit rather than assuming this whole domain is low-confidence. Look up exact APE fields, GDT syntax, and material settings in **t7kb** (`search` then `get`); this skill is the pipeline and the gotchas around it.

## Extraction tools aren't interchangeable

- **Wraith** — the classic ripper for BO3/BO2/BO1/MW/MW2/MW3; exports `.MA` (Maya scene) + `.XMODEL_EXPORT` directly. Good default for older titles.
- **Greyhound** — exports across most CoD games, including live-loading a running game to pull whatever's in memory (e.g. the factory zombie rig, a stock character). Also used to automate GDT/asset creation, not just export.
- **Cordycep** — needed for newer Treyarch/IW titles (Cold War, MW2019+, Vanguard, Warzone) that Wraith/Greyhound can't reach directly, since it works around the anti-cheat/newer packing rather than reading files live.
- **Legion** — a *different* tool needed specifically for **Apex Legends**, which packs assets as `.rpak` archives Cordycep doesn't handle. Don't assume Cordycep covers every "newer" title — check which archive format the source game actually uses.
- **Spiki's tools** — automate GDT/APE asset creation specifically for Treyarch titles (BO3/BO4); pairs well with Greyhound/Cordycep exports to skip a lot of manual APE data entry.
- **Kronos** or **`export2bin.exe`** (ships in the mod tools' `bin/` folder) — either converts `.xmodel_export`/`.xanim_export` to the `.xmodel_bin`/`.xanim_bin` APE actually wants; Kronos additionally converts ripped textures to `.TIFF` (the only image format BO3 supports).

## The weapon porting pipeline (also the template for character/prop ports)

Order matters here — skipping ahead (e.g. exporting before attaching joints) is the usual cause of a silently broken port:

1. **Obtain** the source assets (Wraith/Greyhound/Cordycep/Legion — pick per source game, see above). Rip **every** piece — scope, magazine, and body are separate models.
2. **Prepare in Maya.** Open the `.ma` via `File > Open`, never drag-and-drop — dragging merges it into the scene and namespaces everything, causing problems later.
3. **Attach imported joints** to the main model's root (`j_gun`/`tag_weapon`).
4. **Export** via **Call of Duty Tools → Export XModel**, selecting the full hierarchy (**`Select > Hierarchy`**, not just the root/`tag_origin` — clicking only the root joint looks like it selected everything but doesn't, and is the single most common "export fails silently" cause).
5. **Convert** with Kronos or `export2bin.exe` to `.xmodel_bin`.
6. **Compile in APE**: new GDT (save it under `Black Ops 3\source_data` — a GDT saved elsewhere silently won't reappear in APE/Radiant next session), `xmodel` asset, type **animated**, `BulletCollisionLOD` = LOD0, submodels parented to the main model's root tag.

**ADS export tags depend on the source game** — a common silent-fail point: from a Treyarch-source rig, export only `tag_view`+`tag_torso` for ADS; from an IW-source rig, export only `tag_view`+`tag_ads`. Grabbing the wrong pair for the source you ripped from is a frequent cause of broken ADS.

**Naming convention** (Treyarch's own): `<game_id>_<usage>_<class>_<name>_<type>`, e.g. `t6_attach_mag_dsr50_view` — keep ported assets consistent with this so your GDT stays navigable.

**Materials**: BO3 is PBR — diffuse/albedo, ambient occlusion, normal, specular, gloss. From an older CoD, the color map (`_c`) maps to diffuse and the environment map (`_e`/`_env`) maps to gloss. Weapon materials: `Material Category` = Geometry, `Material Type` = `lit_weapon`. Two easy-to-hit pitfalls: leaving **Surface Type** at `<error>`/`error` (use `<none>` instead), and non-power-of-2 image dimensions (BO3 requires power-of-2 textures). Set one material's fields, then duplicate + rename for the rest rather than re-entering settings each time — same trick works for xanim assets.

**Animations**: convert with the community `conversion_rig.ma` (root bone becomes a child of `t7:tag_weapon_right`, rename joints with `renameRig.mel`, import the anim, strip the rig with `removeNamespace.mel`, rename the root back to `tag_weapon`). **Cold War's animation filenames are dehashed** — pull weapon anims from Modern Warfare instead, it works more reliably. Compile in APE as an `xanim` asset — the settings differ by animation kind: **viewmodel anims** use Use Bones unchecked / `Type` = `relative`; **everything else (world/character anims)** uses Use Bones checked / `Type` = `delta` — mixing these up is a common cause of a ported anim looking right on the weapon but broken on the world model, or vice versa. Check **Looping** for idle/sprint-loop/slide-loop/swim anims regardless of type.

Finish: `bulletweapon` asset in APE, add to the map's zone file (`weapon,<name>`), drag it in from Radiant's entity browser.

## Custom/ported characters and zombies

Same shape as weapons, plus rigging: extract the **factory zombie/character rig** (Greyhound), import your mesh in Maya/Blender, bind joints to the mesh, and weight-paint. **Split the head from the body and rig it to the default BO3 head separately** if you need jaw movement (a community `binder.mel` script automates the head-to-body attachment). An alternative some modders use: import a stock BO3 character alongside yours, copy its weight paint, then swap the armature onto your model and delete the stock one — works as long as the joint hierarchy matches.

For wiring a rigged zombie into custom AI behavior (archetypes, spawners) rather than the modeling/rigging itself, see **bo3-zombies-ai**.

Swapping the **stock Zombies player model** is a separate, simpler task from custom-character porting: swap to another stock crew member via `zm_usermap.gsc`'s character-index override, or build a fully custom player model via `customizationtable`/`playerbodytype`/`playerbodystyle` duplication in APE.

## Collision: ported/custom models have zero collision by default

Unlike stock assets, a ported or custom xmodel has **no collision unless you give it one** — a bare model will let players/zombies walk straight through it. Fix via the xmodel's `CollisionMap` field pointing at a `.map` under `share/raw/collmaps/`, textured with `clip_physics`; the collmap must be **brush-only** (patches don't produce valid collision). Scaled models need `scaleCollMap` set to match, and `use_collmap`/`no_collmap`/`use_misc_models_collmaps` control whether it actually gets baked in. This is easy to miss because the model looks and behaves fine right up until something needs to path around or stand on it.

## Script bundles (the asset-pipeline side of data-driven content)

Scene/cinematic, vehicle, killstreak, and collectible data lives in **script bundles** — GDT-authored, data-driven assets read at runtime via `struct::get_script_bundle`, `get_script_bundle_list`, and `get_script_bundle_instances` (verified against shipped `struct.gsc`). Reach for a script bundle instead of hardcoding a table in GSC when the data is really asset content (e.g. a set of placeable collectible variants) — it keeps the data in APE/GDT where the rest of the pipeline expects it.

## Common pitfalls

- **Blender's current Blender-COD plugin (the GitHub one) can break UV export.** If ported textures look wrong/shifted after export, community consensus is to fall back to a legacy release rather than debug the current one.
- **Duplicate GDT asset errors** (`Duplicate 'material' asset '<name>' found in ...gdt:<line>`) mean the same asset name exists in two GDTs (yours and a shared/stock one) — delete your duplicate entry, don't rename around it; it's a naming collision, not a corruption.
- **Ragdoll behavior for a custom model** goes through `RagdollSettings` — a dragged-in stock ragdoll setup silently keeps stock proportions/behavior unless you edit it for your model.

## Don't invent

Community export/rigging workflows here have real version-specific gotchas (a plugin build, a Maya version, a specific rig) that change over time — verify the current tool/plugin version against what the corpus and raw install actually show before asserting a fix still applies, rather than assuming yesterday's Discord answer is timeless.
