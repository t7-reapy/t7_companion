# BO3 Modding — Technical Documentation

A working knowledge base built while shipping **Rainy Doom** (`zm_test`), my first Black Ops 3 zombie map. Captured before it fades.

> Audience: future me, future contributors, and engineers from other domains curious how a 2015 game engine actually thinks. Plain English, pragmatic, opinionated.

**Contributing / writing protocol:** see [`CONTRIBUTING.md`](./CONTRIBUTING.md) — trust ranking, per-page loop, how external community resources flow in.

## Reading order

The doc is organized by the *size of the topic* in this project, biggest first:

1. **[Asset pipeline](./asset-pipeline/)** — how raw files become things the engine can load. The biggest, most-touched-everywhere subsystem. *Reviewed and stable.*
2. **[Scripting](./scripting/)** — GSC/CSC, the server↔client model, and the joy of debugging without a debugger. *Overview + API reference + several recipes done; more pending.*
3. **[Mapping](./mapping/)** — Radiant, BSP geometry, brushes/patches/prefabs, zoning, AI pathing, portals/umbra. *Stubs.*
4. **[Lighting](./lighting/)** — light bake, light grid, probes, reflection probes, the Light step, volumetric lighting. Its own discipline; tightly coupled to mapping but big enough to stand alone. *Stubs.*
5. **[Audio](./audio/)** — sound aliases, zones, ambient rooms, music, VOX, splitscreen positional audio, royalty-free sourcing. *Stub — the section I'm shallowest on.*
6. **[Systems](./systems/)** — gameplay mechanics that thread through every other section: perks, Pack-a-Punch, mystery box, wonder weapons, traps, traversals, easter eggs, power-ups. *Stubs.*
7. **[Lua & LUI](./lua-lui/)** — the HUD/menu layer: Lua VM, CSC↔LUI bridges, custom widgets, loading screens, fonts. *Stubs.*

Then the supporting cast:

8. **[Community](./community/)** — wikis, Discords, contributors. *Stubs + contributors registry.*
9. **[Reference](./reference/)** — file formats, glossary, command cheatsheets. *Reviewed.*

## Conventions

- Anything I'm not 100% sure about is flagged with `> ⚠️ uncertain` and a hint about how to verify.
- Tool pages follow a fixed template: *what it is / problem it solves / pipeline fit / gotchas / canonical source*.
- File paths are relative to the BO3 install root (`D:/SteamLibrary/steamapps/common/Call of Duty Black Ops III/`).
- Concrete file references prefer real examples from this repo over hypothetical syntax.
