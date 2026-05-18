# Lua & LUI

> Stub — content pending.

The HUD/menu layer in BO3 is **LUI** (Lua-driven UI), running on its own VM separate from GSC and CSC. Big enough to be its own section: it has its own language (Lua), its own asset types, its own debugging story, and its own runtime.

Most map-side HUD work is *driving* LUI from script — sending values from GSC over a clientfield, and having a CSC- or LUI-side widget react.

## Pages (planned)

- `01-overview.md` — what LUI is, the Lua VM, how it relates to GSC/CSC, what's authored where.
- `02-bridges.md` — clientfield → CSC → LUI data flow; supported APIs to push state into the UI.
- `03-hud-elements.md` — score, ammo, perk icons, intro text, round banner — existing widgets and how to override per map.
- `04-custom-widgets.md` — adding a new HUD widget: Lua file location, asset registration, anchoring, render order.
- `05-loading-screen.md` — Steam Workshop preview image, loading screen, custom fonts, intro cinematic glue.
- `06-debugging-lui.md` — how to actually see Lua errors (no debugger, console output gated — document the workarounds).

## TODO (scope)

- LUI vs HUD elements vs `luiroot` / `hudelem`-era APIs — what's still supported in BO3
- The CSC ↔ LUI bridge (how Lua widgets read CSC state)
- Map-specific HUD element registration: where it lives, asset type involved (`luimenu`, `luiwidget`? — verify)
- Custom perk icons / shaders and their LUI binding
- Intro text and round-banner customization
- Score / point HUD overrides
- Loading screen and Steam Workshop preview image
- Custom font registration

## Research starting points

- YouTube corpus: `theme` includes `hud-lui`.
- Vanilla LUI source: `share/raw/ui/lui/` (verify path).
- See also: [community packs](../asset-pipeline/07-community-packs.md) for HarryBo21 perk icon / shader integration.
