# Brush textures (tool textures, clip flags)

> **Status: stub.** Special textures applied to brush faces that change *engine behaviour* rather than appearance — caulk, nodraw, hint, skip, clip variants, etc. These are the "metadata" textures of the BSP.

## What this page will cover

- Why "tool textures" exist: faces don't always need to be drawn; some carry compile-time hints; some control collision selectively.
- Common tool textures: `caulk` (don't render this face), `nodraw` (similar with subtle differences), `hint` (visibility hint for portal generation), `skip` (no portals), various `caulk_shadow*` variants.
- Clip flags by texture name: `clip` (fully blocking), `clip_player`, `clip_ai`, `clip_vehicle`, `clip_full`, `clipshot` (block bullets), `missileclip`, `clip_physics`, `clip_out_of_bounds`.
- Special-purpose: `traverse`, `ladder`, `mantle`, `umbra_*`.
- The `_carver`, `_blend`, `*_decal` patterns and what they signal at compile time.
- Naming patterns and how `RadiantFilters.json` exposes them as filter categories.
- Common gotchas: applying clip on the wrong face, hint brushes generating useless portal cuts, carver brushes leaving fragments.

## Reference reading

- [Call of Duty 5: Tool Textures (zeroy wiki)](https://wiki.zeroy.com//index.php?title=Call_of_Duty_5:_Tool_Textures) — older but the concepts apply directly to T7. The closest thing to a single reference for what each tool texture does.
- `radiant/configs/RadiantFilters.json` in the install — confirms which textures the editor recognises as filterable categories.
