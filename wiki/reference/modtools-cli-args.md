# Modtools launcher CLI arguments

> **Provenance**: gathered from the **community in the newest BO3 patch dump (February 2026)**, published on Discord. Not official Treyarch documentation — semantics for many flags have to be inferred or experimented with. Use as a reference when scripting builds, automating compiles, or chasing obscure flags. The full list below is the *names of all flags the binaries appear to recognise*; descriptions in the "useful ones" table further down are best-effort and clearly marked as guesses where applicable.

## Sorted alphabetically

```
-auth
-baselang
-bulletReport
-cachepath
-cacheupload
-cleanup
-clearwindowsfilecache
-compilealltechsets
-compileshaders
-compress
-compression
-convert
-convertall
-convertalljobindex
-convertallnumjobs
-convertumbratome
-ddlPaths
-debugdeps
-debugweaponmerging
-debugweaponmergingverbose
-debugweaponvariant
-dedicated
-disable--sprites
-disable--translate
-disableumbrasndbs
-dontencryptpcffotd
-echo
-enable--translate
-filtercolor
-force--sprites
-forceconvert
-forcecrypt
-forceimagestreaming
-forcelod
-forcestable
-forceunstable
-fs_game
-gameDir
-gencrypt:
-gsc_emitwarnings
-gsc_optimizedvar
-gsc_profiler
-gsc_stripdebug
-gsc_testfunctions
-ignoreModelErrors
-includexmodelcollisiontree
-keepshaderdebuginfo
-linkshipped
-localized
-luacheapdedicatedcheck
-mapEntsFileName
-modelinfo
-modSource
-modZone
-mtumbraconvert
-noassertdlg
-noauth
-nocachedownload
-nohlsllinedirective
-nolinkerdb
-nonetwork:
-nosndbsshadercompile
-nosndconvert
-nosounds
-nosummary
-noumbraoptimizations
-noxpak
-pause
-printbgcache
-printledinfo
-shaderwarnings
-skipdeps
-skiploc
-skipprelink
-smppercentage
-sndbsshadercompilebatchsize
-spawnedchild
-stripdebug
-stripshaderdebuginfo
-summary
-textureComboStreamedImgReport
-threads--11
-umbra_disabletomefx
-umbra_disabletomelights
-unittest
-useemptybitfields
-validateai
-validateCoverNodes
-verbose
-version
-waitfordebugger
-writeArchiveTypeInfoArray
-writebgcachetofile
```

## Notes on the most useful ones

Most of these are internal / Treyarch-build-time flags. The handful that come up in modder workflows:

| Flag                     | What it does                                                                      |
| ------------------------ | --------------------------------------------------------------------------------- |
| `-printbgcache`          | Dumps BGCache contents at link time. Useful when chasing precache caps (see [`bgcache-caps.md`](./bgcache-caps.md)). *Verified existence via Discord discussion; behaviour with the launcher GUI is finicky — may need direct invocation on `linker_modtools.exe` / `modtools.exe`.* |
| `-bulletReport`          | Generates the `_bulletreport.csv` model-size diagnostic (verified — the file appears in `zone_source/<lang>/assetinfo/` after a Link). |
| `-writebgcachetofile`    | **guess:** writes BGCache contents to a file rather than the console — extrapolated from `-printbgcache`'s name. Verify before relying on it. |
| `-stripdebug`            | **guess:** strip debug info from the build (name-inferred).                       |
| `-gsc_stripdebug`        | **guess:** strip GSC-specific debug info (name-inferred).                         |
| `-gsc_emitwarnings`      | **guess:** tell the GSC compiler to emit warnings rather than swallowing them silently (name-inferred). |
| `-gsc_profiler`          | **guess:** enable the GSC profiler (name-inferred).                               |
| `-verbose`               | **guess:** verbose linker logging — common universal flag convention.             |
| `-summary`               | **guess:** print a build summary at the end (name-inferred).                      |
| `-force` / `-forceconvert` | **guess:** force re-conversion of assets even when caches say they're up-to-date (name-inferred). |
| `-skipdeps`              | **guess:** skip dependency resolution (dangerous; diagnostic builds only — name-inferred). |
| `-noassertdlg`           | **guess:** suppress assert dialogs (useful for batch / CI runs — name-inferred). |
| `-modSource`, `-modZone`, `-fs_game` | **guess:** mod-targeting flags for which mod to build / load.            |
| `-baselang`, `-localized`, `-skiploc` | **guess:** localisation-related flags for language-specific builds.    |

## Categories at a glance

- **Build control**: `-cleanup`, `-clearwindowsfilecache`, `-noxpak`, `-skipprelink`, `-summary`, `-verbose`, `-pause`, `-spawnedchild`
- **Asset conversion**: `-convert`, `-convertall`, `-forceconvert`, `-compileshaders`, `-compilealltechsets`, `-mtumbraconvert`
- **Compression / encryption**: `-compress`, `-compression`, `-forcecrypt`, `-gencrypt:`, `-dontencryptpcffotd`
- **Debug / diagnostics**: `-printbgcache`, `-writebgcachetofile`, `-bulletReport`, `-modelinfo`, `-debugdeps`, `-debugweaponmerging[verbose]`, `-debugweaponvariant`, `-validateai`, `-validateCoverNodes`, `-textureComboStreamedImgReport`, `-printledinfo`, `-keepshaderdebuginfo`, `-shaderwarnings`, `-noassertdlg`, `-summary`
- **GSC compiler**: `-gsc_emitwarnings`, `-gsc_optimizedvar`, `-gsc_profiler`, `-gsc_stripdebug`, `-gsc_testfunctions`, `-stripdebug`, `-stripshaderdebuginfo`, `-nohlsllinedirective`
- **Umbra (visibility)**: `-convertumbratome`, `-disableumbrasndbs`, `-mtumbraconvert`, `-noumbraoptimizations`, `-umbra_disabletomefx`, `-umbra_disabletomelights`
- **Networking / auth**: `-auth`, `-noauth`, `-dedicated`, `-nonetwork:`
- **Cache control**: `-cachepath`, `-cacheupload`, `-nocachedownload`, `-clearwindowsfilecache`
- **Misc**: `-echo`, `-filtercolor`, `-pause`, `-threads--11`, `-smppercentage`, `-luacheapdedicatedcheck`, `-waitfordebugger`, `-writeArchiveTypeInfoArray`, `-useemptybitfields`, `-unittest`

## Status

This list captures *the existence of* each flag — actual semantics for many of them are best discovered by experiment or by the modtools launcher's own help output. **Most modders never touch more than a handful** (`-printbgcache`, `-bulletReport`, `-verbose`, `-summary`). Treat the long list as "things that exist if you ever need them," not "knobs you should tune."
