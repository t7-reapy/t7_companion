# Scripting

> **Status: stub.** Top-level discipline. The second-largest section after asset-pipeline. GSC (server-side) and CSC (client-side), the LUI/Lua HUD layer, and everything around them.

## What this section will cover

- GSC vs CSC: the server↔client split, what runs where, why it matters.
- The GSC language: syntax, threads, `endon`/`notify`, namespaces, `array::thread_all`, common idioms.
- Clientfields: the data sync mechanism, bit budgets, registration, `CF_CALLBACK_ZERO_ON_NEW_ENT` and other late-joiner flags.
- Script overrides: how `usermaps/zm_test/scripts/zm/*.gsc` shadows `share/raw/scripts/zm/*.gsc`, the `zm_patch.csv` gotcha.
- Engine quirks: silent failures, `%`-prefixed animations crashing CSC, wallbuy hard limit, the `_up_up` second-PaP override pattern.
- Tooling: GSCode, Zoroth, Cerberus, BO3 source explorer, GSC Injector hot reload.
- Debugging: dvars as poor-man's feature flags, Sphynx commands as the closest thing to a debugger, reading the t7-source decompile.
- LUI / Lua HUD: separate runtime, separate failure modes, T7Overcharged for advanced cases.
- Splitscreen as a distributed-systems problem.
- Custom features in `zm_test`: hellround state machine, weather state machine, the `_up_up` PaP chain, room-of-thanks elevator/board logic, easter-egg systems.

## Sub-pages (planned)

- `01-overview.md` — server↔client model, pipeline mental model
- `gsc-language.md` — language reference
- `clientfields.md` — sync mechanism, late-joiner fix, splitscreen
- `script-overrides.md` — overriding shipped scripts
- `engine-quirks.md` — silent failures, hard limits
- `tooling.md` — GSCode / Zoroth / Cerberus / Injector / dvars
- `api-reference.md` — pointers to the script API surface
- `lui-lua.md` — HUD scripting
- `debugging.md` — dvars, sphynx, t7-source

## Reference reading (general)

- [GSCode — VS Code extension (blakintosh)](https://marketplace.visualstudio.com/items?itemName=blakintosh.gscode) — syntax, region folding, go-to-definition.
- [`shiversoftdev/Black-Ops-3-Projects`](https://github.com/shiversoftdev/Black-Ops-3-Projects/tree/main) — released BO3 mod projects, useful as reference patterns.
- [Modme forum — Scripting](https://forum.modme.co/wiki/forums/21.html) — the long tail of community Q&A.
- [Abnormal202 Scripting Tutorials (UGX)](https://www.ugx-mods.com/forum/scripting/91/abnormal202-scripting-tutorials-master/16746/).

## Status

This is Reapy's second-largest area after mapping. Most of his lived experience post-March 2025 is here. Expect this section to be the densest with original signal once we start writing it.
