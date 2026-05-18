# Asset Pipeline — Overview

> "How does a PNG on my hard drive end up as a wall texture inside a running map?"

This page lays out the universal flow that *every* asset (image, model, sound, FX, weapon, even a script) goes through. Subsequent pages drill into specific asset types and tools — this is the mental model.

## TL;DR

```
raw file  →  GDT entry  →  zone manifest  →  Linker  →  fastfile  →  engine loads
(bitmap,     (declares      (says "include    (packs        (.ff +       (in-game)
 .xmodel,     it as an       this in zm_test")  everything)   .xpak)
 .wav...)     asset, with
              params)
```

Three artifacts, three steps, two outputs. Everything else is a variation on this theme.

## The mental model

BO3 is a **declarative** asset system on top of a **link-time** packer:

1. **You declare assets in GDT files.** A GDT (Game Data Table) is a pseudo-JSON file holding records like *"this `image` asset is named `tile_concrete_d`, points at `texture_assets/tile_concrete.tga`, uses BC1 compression, has these mipmap settings."* The GDT is opened and edited in **APE** (`bin/asseteditor_modtools.exe`), but it's just a text file underneath — you can grep it.

2. **You list assets in a zone manifest.** A `.zone` (or `.zpkg` sub-package) file is a flat 2-column CSV: `<asset_type>,<asset_name>`. It's the equivalent of a linker script: *"to build `zm_test.ff`, include these assets."*

3. **The Linker packs everything into a fastfile.** `bin/linker_modtools.exe` reads the zone manifest, resolves each asset by looking up its name in your GDTs, pulls the raw file off disk, applies the engine-side processing (DDS conversion, mesh chunking, audio resampling…), and writes one `.ff` (manifest + small data) plus one `.xpak` (bulk binary blobs) per language to `usermaps/zm_test/zone/`.

4. **The engine loads the fastfile.** When the map starts, the engine mmap's `.ff`, follows references into `.xpak` for the heavy blobs, and the asset is now addressable by its name from script.

The crucial insight: **the asset name in the GDT is the only handle the engine ever sees.** Filenames on disk are linker-time only. From script you say `LoadFX("hellround/meteor")`, not `LoadFX("F:/foo/meteor.efx")`.

## Where raw files actually live

There is no single "assets folder." The Linker resolves asset paths against a configured set of search roots, listed in **`bin/converter_gdt_dirs_0.txt`**:

```
_custom
aitype
archetypes
model_export
source_data
temp_assets
texture_assets
xanim_export
art_assets
usermaps
mods
```

Anything you reference from a GDT is searched in these roots, in order. Practically:

- **`source_data/`** — weapon GDTs and other "source-of-truth" project files for major systems.
- **`texture_assets/`** — bitmap source files (TGA/TIFF/PNG) for materials.
- **`model_export/`** — XMODEL_EXPORT files (the ASCII mesh format the linker consumes), often grouped by author/pack.
- **`_custom/`** — local stash for assets you authored or reorganized; **not** searched by default — it's there because the user added `_custom` to `converter_gdt_dirs_0.txt` explicitly.
- **`xanim_export/`, `aitype/`, `archetypes/`, `art_assets/`, `temp_assets/`, `mods/`, `usermaps/`** — specialized roots for animation exports, AI configs, etc.

> 💡 If the linker complains *"could not find asset X"*, ask in this order: (1) is the asset zoned (or pulled in transitively by something that is)? (2) which root did you drop the raw file in, and is that root listed in `converter_gdt_dirs_0.txt`? (3) does a GDT actually declare `X`?

## Zone files: the manifest layer

Every fastfile is built from one `.zone` file plus any `.zpkg` it `include`s. From `usermaps/zm_test/zone_source/zm_test.zone`:

```
>class,zm_mod_level
>group,modtools

// Custom LUA HUD
ttf,fonts/default.ttf
include,h1_hud
scriptparsetree,scripts/zm/_typewriter.gsc

// Custom perks
include,zm_test_perks
scriptparsetree,scripts/zm/_zm_perk_light_fix.gsc

xmodel,skybox_default_day
material,luts_t7_default

// BSP
col_map,maps/zm/zm_test.d3dbsp
gfx_map,maps/zm/zm_test.d3dbsp

// Audio
sound,zm_test
```

Two types of lines that matter:

- **`<asset_type>,<asset_name>`** — pull this asset (resolved via your GDTs) into the fastfile. Common types: `image`, `material`, `xmodel`, `xanim`, `sound`, `fx`, `weapon`, `scriptparsetree`, `scriptbundle`, `stringtable`, `ttf`, `col_map`, `gfx_map`, `rawfile`, `lua_file`. (Full alphabetical taxonomy: see [`reference/asset-types.md`](../reference/asset-types.md).)
- **`include,<name>`** — pull in the contents of `<name>.zpkg` (sub-manifest). Pure organizational tool — there's no namespacing or scoping, it's literally textual inclusion. Use it to keep the root `.zone` readable.

The `>class,...` and `>group,modtools` lines at the top tell the linker *what kind* of fastfile to build. **Verified** by grepping `share/raw/` + `usermaps/` for `>class,*`:

| `>class,...`     | What gets built                                                                          |
| ---------------- | ---------------------------------------------------------------------------------------- |
| `core`           | Core engine fastfile (shipped game code — usually only Treyarch builds these).           |
| `mp_level`       | A multiplayer map fastfile.                                                              |
| `mp_mod_level`   | A multiplayer **mod** map fastfile (community MP map).                                   |
| `zm_common`      | Shared zombies code/assets fastfile that every zombies level depends on.                 |
| `zm_mod_level`   | A zombies **mod** map fastfile (community ZM map — what `zm_test` is).                   |

That's the complete set in this install.

### When can you skip zoning?

You only have to add an asset to a zone manifest if nothing else already pulls it in. Skip cases:

- **It's used by your `.map` file.** Anything placed in Radiant gets auto-included via the BSP (`col_map`/`gfx_map`) entries.
- **It's a dependency of an already-zoned asset.** A material that's referenced by a zoned `xmodel` comes along for free; an image referenced by a zoned material comes along for free; etc. The linker walks the dependency graph.
- **It's already pulled in by a `.zpkg` you `include`.**
- **A `.zpkg` you include itself includes another `.zpkg`.** Inclusion is recursive — only `.zpkg` files participate in include chains, never another `.zone`.
- **An `aitype` GDT record names a `.zpkg` via its `csvInclude` field.** `aitype` is (as far as I know) the only asset type that carries this — every record has a `csvInclude` slot that names a `.zpkg` the linker pulls in whenever this aitype is zoned. So zoning a single `aitype,zombie_default` transitively brings in `zombie.zpkg` (with all its xmodels, xanims, sounds, fx) for free. See `source_data/_midgetblaster/t7_aitype.gdt` — most records there set `"csvInclude" "<something>.zpkg"`.

> ⚠️ There are likely other implicit-inclusion cases beyond these. Add as discovered.

### How community packs actually fit in

There's a spectrum:

- **Targeted feature packs** (a single weapon, a perk variant, one wonder weapon) usually ship a `.zpkg` *and* the GDTs the zpkg references. To install: drop both, then `include` their zpkg from your `.zone`. The zpkg *is* the public surface of the pack.
- **Umbrella asset packs** (Midgetblaster v2.5, harrybo21's gun pack, the SkyeLord weapon ports) ship mostly **GDTs** plus raw assets, expecting you to pick what you want and zone it yourself. They rarely include a single "include this zpkg and you're done" surface — the surface area is too big.
- **Hybrid** (most common): the pack ships a few zpkgs grouping the most-likely-used bundles, *plus* loose GDTs for cherry-picking. Midgetblaster's `t7_aitype.gdt` is a good example of the GDT-side: each `aitype` record points at a `.zpkg` via `csvInclude`, so zoning one aitype is enough to pull in everything it needs.

So the dependency direction is *not* one-shaped. Sometimes you `include` a pack-supplied zpkg; sometimes you write your own zpkg that lists pack-supplied GDT assets; sometimes the GDT itself transitively pulls a zpkg in. The common thread: **GDTs declare what assets *exist*; zpkgs declare what bundles *belong together*.**

## The Linker / Modtools Launcher ritual

Day-to-day work happens in **`bin/modlauncher.exe`** (the GUI). The launcher also shows a list of installed maps and mods with a checkbox per item — meant for **batching** several compile/light/link cycles together. In practice the batch mode is unreliable on heavy work: a long Link step can block the launcher's next planned action (e.g. Run, or the next item in the batch), so for big maps you tend to drive each step manually. Per-item, the launcher exposes 4 checkboxes, executed top-to-bottom when you hit Run:

| Step    | What it does                                                                          | When you need it                                                                              |
| ------- | ------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| Compile | Compiles the `.map` file from Radiant into a `.d3dbsp` (geometry, collision, portals, navmesh). Two modes: **All** (full rebuild — both brushes and entities) and **Just Ents** (only entity changes, much faster). | After geometry / brush / entity changes in Radiant. |
| Light   | Bakes lighting into the BSP — light grid, probes, lightmaps.                          | After moving lights, changing materials' lighting, or after a full Compile.                   |
| Link    | Runs `linker_modtools.exe` to pack everything into `.ff` + `.xpak` files in `zone/`.  | After **any** asset change — script, sound, image, model, weapon, GDT edit.                   |
| Run     | Launches `BlackOps3.exe` with the right command-line to load your fastfile.           | When you're ready to test.                                                                    |

**The Link step is the magic one.** It's the fast iteration loop: edit a `.gsc`, edit a sound alias, swap a texture — re-Link, re-Run, you're testing in seconds. You only re-Compile when geometry changed, only re-Light when lighting math changed. *Caveat for big maps*: Link is "fast" relative to Compile/Light, but on a map with thousands of zoned assets and many scripts, Link itself can take **several minutes**.

> *Brushes vs entities*: in Radiant, **brushes** are the static convex geometry that builds the world (walls, floors, optimized into the BSP). **Entities** are placed objects with properties and behavior — lights, spawn points, script_models, triggers, etc. The "Just Ents" Compile mode skips the slow BSP rebuild because entity placements live in a separate part of the `.map` file. (Brushes also include **patches** — curved surfaces — covered in the mapping doc later.)

### Which language gets built

The Link step localizes per language. The set of languages it builds is configured in the modtools launcher: **Edit → Options → Build Language**, expand the dropdown and pick either *All* or a specific language (English, French, etc.). That choice determines which of the per-language outputs you'll find in `zone/` — `en_zm_test.ff`, `fr_zm_test.ff`, etc. The base `zm_test.ff` (language-neutral) is always built.

## Worked example: image → in-game wall texture

Let's trace one image, end to end.

1. **Author the bitmap.** You make `tile_concrete_d.tga` in Photoshop. Drop it into `texture_assets/myproject/`.
2. **Declare it as an `image` GDT entry.** Open APE, create a new asset of type `image` named `tile_concrete_d`. Point its `baseImage` field at `texture_assets/myproject/tile_concrete_d.tga`. Set compression (e.g. BC1 for color, BC5 for normal), mipmap behavior, color space. APE saves this into a `.gdt` file under `source_data/` (or wherever your GDT lives) as a JSON-ish record.
3. **Declare a `material` GDT entry that uses it.** Materials are the engine's *shader instances* — they bind a shader template (e.g. `mtl_t7_world_lit`) to a set of image slots (color, normal, spec/gloss, etc.) plus parameters (roughness, tiling, emissive). The material is what the renderer actually sees; the image is just one of its ingredients. Name it e.g. `mc_tile_concrete`.
4. **Apply the material** to a brush face in Radiant, *or* assign it to an XMODEL slot in APE. This is where the dependency graph forms: BSP → material → image.
5. **Zone it (or don't).**
   - If the material is on a brush in your `.map`, the BSP pulls in the material, which pulls in the image. **Zero zone edits needed.**
   - If the material is referenced from script (e.g. swapping a model's material at runtime), add `material,mc_tile_concrete` to your `.zone` so the linker doesn't garbage-collect it.
6. **Hit Link** in modlauncher. The linker:
   - Resolves `mc_tile_concrete` against your GDTs → finds the record.
   - Walks its image references → finds `tile_concrete_d`.
   - Resolves `tile_concrete_d` → finds the TGA on disk.
   - Converts TGA → DDS with the compression you specified, builds mipmaps.
   - Writes the small material record into `zm_test.ff`, the heavy DDS mip chain into `zm_test.xpak`.
7. **Run.** The engine streams it on demand when the camera looks at it.

The asset name `mc_tile_concrete` is now your handle to it from anywhere — script, other materials, whatever. The TGA path is forgotten the moment the linker finishes.

## Engine container formats: `.d3dbsp`, `.ff`, `.xpak`

> Engine internals here lean on the IW/idTech lineage (Quake → Quake III → CoD) plus reading the source of community asset extractors.

The reference tools to know:

- **[Greyhound](https://github.com/Scobalula/Greyhound/) (Scobalula, C++, open source)** — extracts *models, images, animations, sounds* from a running BO3 process and from on-disk `.xpak` files. The XPak parsing code (`src/WraithXCOD/WraithXCOD/XPAKSupport.cpp`, `XPAKCache.cpp`, `DBGameFiles.h`) is the most readable open-source spec for the BO3 streaming container.
- **[HydraX](https://github.com/Scobalula/HydraX/) (Scobalula, C#, open source)** — extracts the *specialized* asset types Greyhound doesn't cover: AI files (`aitype`, `behaviortree`, `behaviorstatemachine`, `animstatemachine`), weapon defs (`weapon`, `attachment`, `attachmentunique`, `attachmentcosmeticvariant`, `weaponcamo`, `tracer`, `flametable`), tables (`footsteptable`, `surfacefxtable`, `playersoundstable`, etc.), `scriptbundle`, `stringtable`, `fonticon`, `localize`, `ttf`, `xcam`, `physpreset`, and more. **Also live-process** — runs against a running BO3 with the relevant level loaded.
- **[Saluki](https://github.com/echo000/saluki-releases) (echo000 / dest1yo, Rust, binary-only releases)** — Greyhound's spiritual successor, supports CoD1 → BO7. Built on dtzxporter's `porter-lib`. Source isn't published, so for *format* questions you read Greyhound; for *modern game support* you run Saluki.

### On-disk `.ff` parsing: ATE47's `acts`

The widely-used open-source on-disk `.ff` reader is **[atian-cod-tools](https://github.com/ate47/atian-cod-tools/) (ATE47, C++)**. The `acts` binary ships a `fastfile` subcommand that opens a `.ff` directly and dumps recognized asset pools:

```pwsh
# Dump GSC scripts from a BO3 fastfile
acts fastfile -r gsc "...\Call of Duty Black Ops III\zone\zm_zod_patch.ff"
# Then decompile the dumped GSC bytecode
acts gscd output_ff\spt -f treyarch -g -o output
```

**BO3 coverage in acts is partial**: only `rawfile` and `scriptparsetree` pools. That's enough to crack open compiled GSC scripts and raw text files from a stock fastfile — extremely useful for understanding how Treyarch's own scripts work. Models, images, weapons, sounds on BO3 still have to come from live-process tools.

Newer titles (BO4 → BO7) have much broader on-disk coverage in acts — the format is better-understood / more-invested-in for newer games.

#### What we learn from acts about the BO3 `.ff` format

- **Compression: zlib.** ACT's linker exposes a `compression` setting with options `none`, `zlib`, `zlib_hc`, `lz4`, `lz4_hc`, and several Oodle variants for BO4+. BO3 fastfiles compress with zlib. (The Oodle family applies to `.xpak` on BO3 and to `.ff` on BO4+.)
- **Patch files exist.** Treyarch ships `.fd` and `.fp` files alongside some `.ff` files as patches/deltas. acts takes a `-p` option to apply them when reading. Whether you'll ever produce them as a modder is a separate question (probably no).
- **The acts linker uses a `.zone` syntax close to the BO3 modtools' own** (per its docs). Different config keys (`>game`, `>name`, `>compression`, etc. instead of `>class`, `>group`), but the *asset line* format `type,name` is identical. Confirms our description of the modtools `.zone` format is on the right track.

### Other reference tools

- **`shiversoftdev/t7-source`** — Treyarch source dump, closest thing to engine ground truth. The actual fastfile loader code lives in there. *Origin*: this is **not** an `acts`-style on-disk extraction — it's a leak / source-code dump of Treyarch's actual development scripts and engine fragments. Useful precisely *because* it goes deeper than what the public extractors can reach.
- **Older-CoD ZoneTool/ZoneBuilder forks** — family-resemblance only (T6 `.ff` ≠ T7 `.ff`, but the lineage informs).

### BSP and `.d3dbsp`

**BSP** stands for *Binary Space Partitioning*. It's a spatial data structure invented for Doom (in 2D) and perfected for Quake (in 3D), where the world is recursively split along chosen planes into a tree whose leaves are convex regions of empty space. Every CoD engine since *Call of Duty 1* descends from idTech 3 (Quake III) and inherits this idea.

A BSP gives the engine two things almost for free:
1. **Visibility**: from any camera position, walking the tree tells you which leaf you're in, and a precomputed PVS (Potentially Visible Set) — built during the Compile/Light steps — tells you which other leaves can possibly be seen from there. The renderer skips the rest.
2. **Collision**: ray and point queries against the static world resolve to a handful of node tests, not millions of triangle tests.

The Radiant Compile step turns your `.map` (text source: brushes, patches, entities) into `usermaps/zm_test/maps/zm/zm_test.d3dbsp`. CoD's `.d3dbsp` extends the classic Quake BSP with: rendered surface chunks with material assignments, lightmap UVs and lightmap chunk references, portal data for visibility, and pathnode/AI navigation data. It is **the** compiled world artifact; everything else about the world is metadata layered on top.

### `col_map` vs `gfx_map`

Two of the most-confusing zone lines:

```
col_map,maps/zm/zm_test.d3dbsp
gfx_map,maps/zm/zm_test.d3dbsp
```

Same source file, two different *views* extracted from it:

- **`col_map`** tells the linker to extract the **collision** representation from the `.d3dbsp` and pack it as a `clipMap_t` engine asset. Used by physics / AI / projectile systems. No materials, no rendering, just the geometry + surface flags needed to answer "did this ray hit something solid?"
- **`gfx_map`** tells the linker to extract the **rendered** representation — surface chunks with material bindings, lightmap data, portals, the PVS — and pack it as a `GfxWorld` engine asset. This is what the renderer actually walks every frame.

The split exists because the engine subsystems consume them separately: a dedicated server only needs the `clipMap_t`; a render-only preview only needs the `GfxWorld`. They live as two distinct entries in the fastfile even though both originate from one Radiant compile.

### Fastfile (`.ff`)

CoD's proprietary asset bundle format. Genealogy: idTech used `.pak` (Quake) then `.pk3` (Quake III, just a renamed zip) — flat archives of arbitrary files. CoD broke from that and went *typed and binary*: a `.ff` is a **zlib-compressed** stream of *engine-ready* asset records, indexed by name and type, structured to be near-mmap-able. Loading a fastfile is closer to deserializing a memory dump than to extracting an archive — most assets need only a pointer fix-up before the engine can use them. That's where the "fast" in fastfile comes from: minimal CPU work between disk and gameplay.

A `.ff` holds the manifest plus the small/critical chunks: asset headers, compiled GSC/CSC bytecode, localized strings, materials, sound aliases, small images, the `clipMap_t` and `GfxWorld` records, etc. Streaming-friendly bulk data (image mips, audio waveforms) is *not* embedded — it lives in the paired `.xpak` and is referenced by hash.

Two-tier verification status:

- ✅ **zlib compression** — confirmed by ATE47's `acts` linker pool, which exposes BO3-compatible compression as `zlib` only (Oodle family is BO4+).
- ✅ **`rawfile` and `scriptparsetree` pools exist** with extractable on-disk layout — confirmed by `acts fastfile -r gsc *.ff` actually working on shipped BO3 fastfiles.
- ⚠️ **Internal record/header layout, pool indexing, pointer fixup table** — informed-but-unverified. The full spec lives in `shiversoftdev/t7-source` if it ever matters.

### Streaming pak (`.xpak`)

Newer companion file (BO3 era). Holds the **bulk binary blobs** that would blow memory if all loaded at once: the heavy mip levels of textures, raw audio waveform data, video. The engine streams these on demand — when the camera approaches a wall, the higher-res mips of its texture get pulled in from `.xpak`; when they're not needed, they're evicted.

The `.xpak` always pairs 1:1 with a `.ff` of the same basename. The `.ff` holds the *handle* to the streaming chunk; the `.xpak` holds the *bytes*.

#### Verified format details (from Greyhound)

The on-disk layout, confirmed by Greyhound's `BO3XPakHeader` struct in `DBGameFiles.h` and the parsing logic in `XPAKSupport.cpp`:

```
[ Header ]                — 64 bytes
[ Hash table ]            — HashCount * 24 bytes (Key, Offset, Size per entry)
[ Index table ]           — IndexCount entries, each: Hash + a key:value text properties blob
[ Data segment ]          — DataSize bytes, Oodle-compressed in blocks
```

- **Magic**: `0x4950414b` (the four bytes `K A P I` on disk, i.e. "KAPI" — likely a stylized "IPAK"). *How this was identified*: hex inspection of the first 4 bytes of any `.xpak` file (`xxd zm_test.xpak | head -1`) reveals the magic; Greyhound's parser hardcodes the same value as the validity check (`if (Header.Magic == 0x4950414b)` in `XPAKSupport.cpp`).
- **Index entries** are *self-describing text*: each holds free-form `name:`, `type:`, `width:`, `height:`, `format:`, `size0:`, `hashN:` lines. This is why an `.xpak` can carry many asset categories without a typed schema — it just hands back bytes plus a small textual description.
- **Compression** is Oodle (`oo2core_6_win64.dll`, the BO3-era Oodle build). Each compressed segment is a sequence of `BO3XPakDataHeader { Count, Commands[Count] }` blocks, where each Command's high byte is a flag and low 24 bits is the block size.
- **Versioning**: Version `0xB` is BO4-era (Greyhound loads a separate `bo4_ximage.wni` name index for it). Version `0xD` is MW2019/MW4 with extra padding bytes (Greyhound special-cases skipping 288 bytes between header halves).

This matters in practice because: image mip data, sound waveforms, and video chunks all live here, addressed by 64-bit hash. The asset record in the `.ff` carries that hash; resolving it means seeking into the matching `.xpak`, decompressing the block, and getting raw bytes back.

### Putting it together: what you find in `zone/`

```
zm_test.ff         en_zm_test.ff         fr_zm_test.ff
zm_test.xpak       en_zm_test.xpak       fr_zm_test.xpak
```

- **`zm_test.ff` + `zm_test.xpak`** — the language-neutral payload (all the assets that don't change between languages).
- **`en_zm_test.ff` + `en_zm_test.xpak`** — English-specific payload (localized strings, English VO).
- **`fr_zm_test.ff` + `fr_zm_test.xpak`** — French-specific payload.

The engine loads the language-neutral fastfile plus the active language's fastfile. Other languages stay on disk — only the player's chosen language ships extra weight.

## Generated reports & diagnostic files

After a Compile/Light/Link run, several side-channel files appear under `usermaps/zm_test/zone_source/`, organised in language subfolders: **`english/`** and **`french/`** for per-language data, plus **`all/`** for language-neutral data. The `<lang>` placeholder in the table below resolves to whichever of those subfolders contains the file. They're not consumed by the engine — they're *for you*, to debug and audit what the toolchain just did.

| File                                                  | Source step | What it tells you                                                                                                                          |
| ----------------------------------------------------- | ----------- | ------------------------------------------------------------------------------------------------------------------------------------------ |
| `<lang>/assetlist/zm_test.csv`                        | Linker      | The **resolved** flat list of every asset packed into that language's fastfile, after walking all `include`s and dependency graphs.        |
| `<lang>/assetinfo/zm_test.csv`                        | Linker      | Per-asset breakdown: index, type, name, **resident bytes**, **streamed bytes**, and a `parentStack` showing the dependency chain that brought it in. Gold for *"why is this asset in my fastfile?"* and *"what's eating my memory?"* |
| `<lang>/assetinfo/zm_test_xmodel.csv`                 | Linker      | Per-model breakdown of every xmodel in the fastfile, with **per-LOD vert/tri counts** (`verts0..7`, `tris0..7`), screen-space radius, ref count. Verified header: `name,refs,radius,onePixelDist,lodCount,verts0,tris0,verts1,tris1,...`. The go-to for "which models are killing performance?" |
| `<lang>/assetinfo/zm_test_bulletreport.csv`           | Linker      | Despite the name, this is a **small-model size report**: lists xmodels under area thresholds (header: `xmodel,avgArea (min 50.0),triCount (max 0),dimX,dimY,dimZ,volume,refs`). Likely a culling / bullet-collision sanity check — flags models too small to render or that may need bullet-pen tuning. *Exact purpose isn't documented; read the file to understand what your map flagged.* |
| `<lang>/assetinfo/zm_test_poolinfo.csv`               | Linker      | Pool budget tracker. Header: `type,limit,total`. One row per asset type showing how many of that type the engine pool can hold vs how many your fastfile uses. Spot-on for "am I about to hit the engine cap?" |
| `<lang>/assetinfo/zm_test.deps`                       | Linker      | Dependency configuration: `ignore` and `ignore_missing_shipped` directives for shipped fastfiles (`core_*`, `zm_common`, etc.) the linker references but doesn't have to resolve from your project. |
| `all/assetinfo/zm_test.badnodes`                      | Compile     | Pathnode/AI-nav diagnostic — bad or unreachable nodes flagged during BSP compile. Useful for tracking down zombie pathing bugs.            |
| `all/scriptgdb/scripts/.../*.gsc.gdb`, `*.csc.gdb`, `*.gsh.gdb` | Linker | Per-script debug info ("game DB"), produced when GSC/CSC/GSH files are compiled to bytecode. Used by tooling, not user-edited.              |
| `<your-chosen-path>/traversals`                       | Radiant (on demand) | Output of **Misc → Generate Reports → Print Traversal Coordinates**. Flat CSV listing every traversal node placed in the map (prefab name, local origin/angles, top-level origin). Useful for auditing AI traversal coverage and importing/relocating them. You pick where it gets written. |
| `<your-chosen-path>/map_stats`                        | Radiant (on demand) | Output of **Misc → Generate Reports → Basic Map Stats**. Per-prefab counts: times included, brush count, entity count, prefab count, model count, fx count, light count. Plus a `Total` row. Use it to spot prefabs that bloated unexpectedly. You pick where it gets written. |

> 💡 If a Link succeeds but the asset doesn't behave as expected in-game, the `<lang>/assetinfo/zm_test.csv` file (e.g. `usermaps/zm_test/zone_source/english/assetinfo/zm_test.csv`) is usually the first place to look. The `parentStack` column answers *"who pulled this in?"* in one line.

**Concrete example** of a parentStack row:

```
21,material,mc/mtl_gib_chunk_set,22,0,|p7_gib_chunk_bone_03|xmodel|zone_source/loc/zm_test.zone|csv
```

Read it bottom-up: the material `mc/mtl_gib_chunk_set` is in your fastfile (22 resident bytes, 0 streamed) because the xmodel `p7_gib_chunk_bone_03` references it, and that xmodel was zoned via `zone_source/loc/zm_test.zone`'s csv chain. So if you saw "why is this gore-material in my map?", parentStack tells you: it came along for free with that gibbing xmodel.

Radiant itself also produces compile reports (BSP build logs, leak files, lighting reports) — see [`docs/mapping/`](../mapping/) for those.

## Cross-cutting concerns to know about

- **Asset names are global.** All asset names live in one flat namespace per type (no folders, no scoping). If two GDTs declare an `image` named `wood_planks_d`, one will silently win. This is a real problem with community packs — see [`07-community-packs.md`](./07-community-packs.md) for the dupe-purger tooling and Midgetblaster-priority handling.
- **Many community packs ship their own GDTs** (often in `_custom/` sub-folders or scattered across `model_export/`). Installing a pack is essentially "adding entries to the global asset registry." This is why install order matters and why duplicate-resolution matters.
- **Script assets are special.** GSC/CSC files compile to bytecode at link time, but their `.gsh.gdb` / `.gsc.gdb` artifacts (in `usermaps/zm_test/zone_source/all/scriptgdb/`) are intermediate debug info — useful for tooling, not user-edited.
- **Scripts can be hot-reloaded** in some setups via the GSC Injector — to be covered in [`docs/scripting/`](../scripting/) (forthcoming).

## Where to go next

- **Adding a new image / material in detail** → [`02-images-and-materials.md`](./02-images-and-materials.md)
- **Adding a custom xmodel** → [`03-models.md`](./03-models.md)
- **Fog (`fog` GDT type — per-zone / per-state)** → [`04-fog.md`](./04-fog.md)
- **SSI (sun & sky info, lighting-state mechanism)** → [`05-ssi.md`](./05-ssi.md)
- **Adding a weapon** (the deep end) → [`06-weapons.md`](./06-weapons.md)
- **Community packs and dupe handling** → [`07-community-packs.md`](./07-community-packs.md)
- **Sounds & audio** → [`docs/audio/`](../audio/) (separate top-level section)

---

## Reference reading

Big-picture / megathread material covering the modtools end-to-end:

- [BO3 Mod Tools — Megathread (IceGrenade wiki)](https://github.com/IceGrenade/bo3/wiki/Black-Ops-3-Mod-Tools-Get-Started) — the canonical "where do I even start" landing page.
- [Black Ops 3 Mod Tools Super Guide (r/CODZombies)](https://www.reddit.com/r/CODZombies/comments/58nbvq/black_ops_3_mod_tools_super_guide/) — older Reddit megaguide; still a good orientation.
- [Modme Wiki — Game Modding BO3](https://wiki.modme.co/wiki/Game-Support-_-Black-Ops-3.html) — the established community wiki.
- [T7Wiki](https://www.t7wiki.com/en/home) — best-maintained reference site, error list, dev commands.
- [ZGC Tutorials (icegrenade.co.uk)](https://icegrenade.co.uk/tutorials) — curated tutorial index.
- [`docs_modtools/`](file:///D:/SteamLibrary/steamapps/common/Call%20of%20Duty%20Black%20Ops%20III/docs_modtools/) — Treyarch's own PDFs (PBR, lighting, sound, FX, GSC reference).
- [BO3 Mod Tools playlist tuto (YouTube)](https://www.youtube.com/watch?v=tj7guP_ZeUI&list=PLnt_Nobu89HtwqIkEtRt4zj6RC_s-_f2q) — long video walkthrough series.
- [Modding tools (DEVRAW)](https://www.devraw.net/resources) and [DTZxPorter Tools](https://dtzxporter.com/tools#utilities-view) — tooling indexes.
- [`shiversoftdev/t7-source`](https://github.com/shiversoftdev/t7-source) — the Treyarch source dump (closest to engine ground truth).

## Open questions / TODO

- [ ] **BSP / `d3dbsp` internal layout** — partially resolved. The closest open-source reference is **[Husky](https://github.com/Scobalula/Husky/) (Scobalula, C#)** — a BSP extractor for CoD 1 → Vanguard. Husky reads geometry (vertices, faces, UVs, material names) from a *running* game's BSP and writes OBJ/MTL — same live-process pattern as Greyhound/HydraX. So Husky covers the *visible-geometry* slice of the d3dbsp container. The full on-disk container format (collision geometry, portals, PVS, pathnodes) is still not fully unpacked in any public tool. The realistic deeper reference remains `t7-source` if you ever need it.
- [ ] Patches (curved Radiant surfaces) — covered in [`docs/mapping/brushes-patches.md`](../mapping/brushes-patches.md) (planned) — see [`docs/mapping/`](../mapping/) for the section overview.
