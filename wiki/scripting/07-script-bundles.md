# Recipe — Script Bundles

> Stub — content pending.

Script bundles (`script_bundle` asset type) are a packaging format Treyarch uses to ship *grouped* gameplay data — fxanim sequences, ai data, weapon variants, scriptable definitions — as a single addressable asset. From a custom-map perspective they show up most often as **fxanim bundles** (animated environment props) and as **scriptbundle-driven AI archetypes**.

## TODO

### Mental model
- A scriptbundle is a precachable, addressable bundle of data declared in GDT
- The actual *contents* of a bundle vary by `bundleType` — fxanim, ai, weapon variant, etc.
- Scripts reference a bundle by name; the engine resolves payload at load time
- Bundles are precached with `#precache("script_bundle", "<name>");`

### fxanim bundles (the common case)
- Content: an authored animated sequence (mesh + anims + maybe FX hookups) intended to play on a `script_model` in the map
- Real example: `usermaps/zm_test/scripts/zm/zm_animated_fauna.gsc` precaches four bundles
  ```cpp
  #precache("script_bundle", "p7_fxanim_cp_lotus_atrium_ravens_bundle");
  #precache("script_bundle", "p7_fxanim_mp_apartments_rat_comd_p1_bundle");
  ...
  ```
- The bundle ships the actual xanims; the script_model uses them via `UseAnimTree(#animtree)` + `AnimScripted(...)` (see [`06-animations.md`](./06-animations.md))
- These specific bundles are **stock Treyarch content from other levels** repurposed — that's a viable strategy when you don't have animator time

### Scriptbundle-driven AI / collectors
- Reapy uses scriptbundles for the **cerberus collectors** in `zm_test`
- TODO: pull the specific files and document the pattern (probably in `usermaps/zm_test/scripts/zm/` somewhere with `_cerberus` or `_collector` in the name)
- Compare/contrast with the fxanim flavor — what changes, what stays the same

### Authoring your own bundle
- Create a GDT entry with `bundleType` set appropriately
- Populate the bundle's content fields (varies by type)
- Add `script_bundle,<name>` to the zone file
- Precache from script
- Reference by name from runtime APIs

## Common gotchas

- Forgetting `#precache("script_bundle", "<name>");` → asset missing at runtime, silent failure or runtime error
- Forgetting the zone-file line → linker doesn't include it; same outcome
- Repurposing a stock bundle whose internal anim names you don't know — the bundle exposes specific anim aliases; you have to know them to call `AnimScripted`

## Cross-references

- [`06-animations.md`](./06-animations.md) — the playback APIs for the anims a bundle ships
- Asset pipeline: [`02-zone-files.md`](../asset-pipeline/02-zone-files.md) — zone-file syntax for `script_bundle,<name>`
- Reference: [`asset-types.md`](../reference/asset-types.md) — the `script_bundle` asset row
