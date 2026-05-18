# Splitscreen Audio

> Stub — content pending. **The one place in `audio/` where Reapy has lived experience** — splitscreen sound bugs are the dominant pain point of the discipline.

How sound emission interacts with splitscreen, and the patterns that survive both single-screen and 2-player local play.

## TODO

### The core problem
- BO3 supports 2-player splitscreen on the same machine
- A "client" in BO3 terms is a *local player view*, not a connected machine — splitscreen has 1 machine but 2 clients
- CSC code runs *per client*. A `PlaySound` call inside CSC fires for **each** local client
- Result: 2 players → sound emits twice → audio doubles up at the same emitter location
- The fix is host-only emission + a `GetLocalPlayers()` loop where appropriate

### The host-only check
- Pattern: gate a server-side emit behind `IsHost()` or equivalent (verify exact API)
- Variant: emit on the *server* (GSC), let the engine replicate to all clients
- When CSC emission is the right call: per-player feedback (kill confirms, perk-acquired stings)

### `GetLocalPlayers()` for per-player feedback
- Returns the array of local player entities (1 in single-screen, 2 in splitscreen)
- Use when the *intent* is per-player feedback that should fire once per local view
- Wrong use: looping `GetLocalPlayers()` for a *world* sound — that's the doubling bug

### Real examples from `zm_test`
- **Teddy bear plays twice in splitscreen** — fixed in commit `8cbb07db` (in the new history). Pattern: world sound was being emitted CSC-side per local client.
- **Perk sound + ambient sound in splitscreen** — commit `c4c31bfb`. Pattern: similar — host check missing for replicated ambience.
- **Client sounds during meteor event only triggered if under certain distance** — commit `5f32d243`. Distance gating in CSC interacted poorly with splitscreen's two camera positions.

### Rain-sounds gotcha
- Rain emitters are world sounds; without host gating, splitscreen doubled the entire rainfall
- Fix pattern documented in the source (when filling, point at the exact file)

## Common gotchas

- Defaulting to "emit in CSC" for everything → splitscreen always doubles
- Defaulting to "emit in GSC" for everything → no per-player audio (kill confirms feel wrong)
- `IsHost()` check that runs *only* on the host machine vs *only* for the host *player* — these are different concepts in splitscreen; verify which you mean
- Distance gating using a single player position when there are two cameras

## Cross-references

- [`sound-aliases.md`](./sound-aliases.md) — the API surface
- Scripting: [`01-overview.md`](../scripting/01-overview.md) — splitscreen sound dedup gotcha cross-reference
- `KNOWLEDGE.md` / commit history — concrete fixes
