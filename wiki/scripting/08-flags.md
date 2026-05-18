# Recipe — Flags (`flag::*`)

> Stub — content pending.

The `flag::*` namespace (defined in `share/raw/scripts/shared/flag_shared.gsc` and CSC equivalent) is BO3's lightweight named-state system. It's the canonical way to gate a thread on "did event X happen yet" without polling.

## TODO

### Mental model
- A flag is a named boolean stored on an entity (usually `level`, sometimes `self` for per-player flags)
- Three states: not initialized → false → true (and back)
- Multiple threads can `wait_till` the same flag; all wake up when it flips

### The core API (verify against `flag_shared.gsc`)
- `level flag::init("flag_name")` — declare the flag (must come before any other op)
- `level flag::set("flag_name")` — flip true; wakes any `wait_till`'ers
- `level flag::clear("flag_name")` — flip false
- `level flag::set_clear("flag_name", value)` — set to a boolean
- `level flag::wait_till("flag_name")` — yield until the flag is true (returns immediately if already true)
- `level flag::wait_till_clear("flag_name")` — symmetric
- `level flag::wait_till_any(["a", "b"])` / `wait_till_all([...])` — multiplex
- `level flag::exists("flag_name")` — guard against double-init

### Real example
From `usermaps/zm_test/scripts/zm/room_of_thanks/zm_room_of_thanks_elevator.gsc`:

```cpp
function private elevator_think()
{
    self endon("kill_elevator_think");
    level flag::wait_till("initial_blackscreen_passed");
    // ... rest of elevator logic only runs after the load-screen blackscreen ends
}
```

Why this pattern: the elevator should only start its logic *after* the player's intro blackscreen finishes. Polling a global would work; a flag is cheaper and clearer.

### Common flags worth knowing
- `initial_blackscreen_passed` — load screen done
- `start_zombie_round_logic` — round system active
- `solo_game` / `splitscreen` (verify exact names) — game mode predicates
- TODO: scrape `share/raw/scripts/zm/_zm.gsc` and friends for the complete vocabulary the engine sets

### Per-player flags
- `player flag::init("my_flag")` — flag lives on the player entity
- Useful for "this player has unlocked the door" without a global side table
- Cleared on `player_disconnect` / `player_died` depending on how you wire it

## Common gotchas

- Reading a flag before `init` returns false silently — no error, just wrong behavior
- `wait_till` on a flag that's never set hangs the thread forever (use `endon` to guard)
- Flag names are strings — typos compile fine, fail silently at runtime
- Per-player flags vs `level` flags — easy to put it on the wrong scope

## Cross-references

- [`04-callbacks.md`](./04-callbacks.md) — flags are often set inside callbacks
- [`01-overview.md`](./01-overview.md) — `endon` interaction
- Treyarch source: `share/raw/scripts/shared/flag_shared.gsc` (canonical)
