# Callbacks recipes — `endon` / `notify` / `callback::*` / custom function pointers

> Three parallel mechanisms for "run this when X happens": the **`callback::*` registry** (engine-named events), the **`endon` / `notify` idiom** (in-script signalling), and **custom function-pointer registries** (your own per-module event hooks). Used pervasively across `zm_test`. This page is the reference for all three.

## Mental model

Three layers of "wait for X" / "run when X":

1. **`callback::*` registry** — the engine fires named events at predetermined points in the gameplay lifecycle (player connect, spawn, damage, kill, etc.). You register a function to be called when each event fires. The registry is in `share/raw/scripts/shared/callbacks_shared.gsc` (GSC side) and `callbacks_shared.csc` (CSC side).
2. **`endon` / `notify` idiom** — your scripts declare *their own* named events on entities (`level`, `self`, any entity), then other threads either *wait* for them (`waittill`) or *terminate* on them (`endon`). This is how scripts coordinate within themselves — concurrent threads communicating via named signals on a shared entity.
3. **Custom function-pointer registries** — your module exposes an `add_<event>_callback(handler)` API for downstream callers to register hooks for events your module fires (e.g. "when the player enters the room of thanks"). Just an array of function pointers your module walks when the event happens.

Together they cover almost all "react to something happening" needs without you having to poll.

## Built-in `callback::*` callbacks

### GSC side (server)

From `share/raw/scripts/shared/callbacks_shared.gsc` — register-functions that take `( func, obj )` (your handler + an optional context object):

| Function                                | Fires when                                              |
| --------------------------------------- | -------------------------------------------------------- |
| `callback::on_finalize_initialization`  | Final init phase — after most level setup is done.       |
| `callback::on_start_gametype`           | At gametype start.                                      |
| `callback::on_connecting(func, obj)`    | Player begins to connect.                                |
| `callback::on_connect(func, obj)`       | Player has fully connected (handler runs with `self == player`). |
| `callback::on_disconnect(func, obj)`    | Player disconnects.                                     |
| `callback::on_spawned(func, obj)`       | Player spawns / respawns (handler runs with `self == player`). |
| `callback::on_loadout(func, obj)`       | Player loadout being assigned.                           |
| `callback::on_player_damage(func, obj)` | Player takes damage.                                    |
| `callback::on_player_killed(func, obj)` | Player dies.                                            |
| `callback::on_joined_team(func, obj)` / `on_joined_spectate` | Team / spectate transitions.        |
| `callback::on_ai_killed(func, obj)` / `on_actor_killed`      | AI / actor dies (handler runs with `self == ai_or_actor`). |
| `callback::on_vehicle_spawned` / `on_vehicle_killed`         | Vehicle lifecycle events.           |
| `callback::on_laststand(func, obj)`     | Player went into Last Stand (laststand_shared register). |

### CSC side (client)

From `share/raw/scripts/shared/callbacks_shared.csc`:

| Function                                       | Fires when                                       |
| ---------------------------------------------- | ------------------------------------------------ |
| `callback::on_localclient_connect(func, obj)`  | Local client connects. (Probably per-local-client in splitscreen — verify in `callbacks_shared.csc` if it ever matters.) |
| `callback::on_localclient_shutdown(func, obj)` | Local client shuts down.                         |
| `callback::on_localplayer_spawned(func, obj)`  | Local player spawn (CSC equivalent of `on_spawned`). |
| `callback::on_finalize_initialization`         | Final init.                                      |
| `callback::on_spawned`                         | Generic spawn (less specific than localplayer).  |
| `callback::on_shutdown`                        | World shutdown.                                  |
| `callback::on_start_gametype`                  | Gametype start.                                  |

> When in doubt about which to use, prefer the **`localplayer` / `localclient`** variants on CSC — they're per-local-player and play nicely with splitscreen.

## Real call sites in `zm_test`

The most-used by far across **all** `zm_test` scripts (including hellround, weather, room-of-thanks, weapon scripts): **`on_spawned`** (≈13 distinct registrations across modules) and **`on_localclient_connect`** (≈5 on the CSC side, mostly weather/hellround init paths).

From `usermaps/zm_test/scripts/zm/zm_test.gsc` (~lines 144–147):

```c
callback::on_connect(&disable_hitmarkers);
callback::on_connect(&notify_ui_for_nuke_powerup);
callback::on_connect(&apply_current_lighting_state);   // late-joiner lighting fix
callback::on_spawned(&on_player_spawned);
callback::on_laststand(&onlaststand);
```

From hellround scripts (`zm_hellround_environment.gsc:40`):

```c
callback::on_connect(&sync_hellround_environment);     // late-joiner hellround sync
```

From magmagat / minigun pickup scripts (note `on_spawned` fires every time the player spawns, including respawns — so these handlers run on every life, not just the first):

```c
callback::on_spawned(&magmagat_on_player_spawned);
callback::on_spawned(&give_hellround_minigun);
```

## Defining your own callbacks (custom event hooks)

Beyond the engine-built-ins, you can declare your own callback registry for module-internal events. The convention in `zm_test` is the `add_<event>_callback(handler)` pattern:

```c
// Caller (your map's main script)
zm_room_of_thanks::add_enter_room_of_thanks_callback(&set_lighting_state_clear);
zm_room_of_thanks::add_exit_room_of_thanks_callback(&set_lighting_state_normal);
```

The implementing module stores handler arrays internally and walks them when the event fires:

```c
#insert scripts\shared\shared.gsh;   // for MAKE_ARRAY / ARRAY_ADD

// In zm_room_of_thanks.gsc (sketch)

function init()
{
    MAKE_ARRAY(level.rot_enter_callbacks);   // once, up front
}

function add_enter_room_of_thanks_callback(handler)
{
    if (IsFunctionPtr(handler))
    {
        ARRAY_ADD(level.rot_enter_callbacks, handler);
    }
}

function private fire_enter_callbacks()
{
    foreach (handler in level.rot_enter_callbacks)
    {
        [[ handler ]]();   // function-pointer invocation syntax
    }
}
```

Notes on the GSC idioms:

- **Use the `shared.gsh` macros** for array boilerplate:
  - `MAKE_ARRAY(arr)` — initialise as `[]` if undefined, or wrap a single value as a one-element array if it's not already an array. Saves the `if (!IsArray(...)) arr = [];` boilerplate.
  - `ARRAY_ADD(arr, item)` — `MAKE_ARRAY` + push in one step (used above).
  - `DEFAULT(var, default)` — set `var = default` only if `var` is undefined.
- **Use `IsFunctionPtr(handler)` to validate function pointers** before storing or calling them. Calling `[[ undefined ]]()` will crash.
- **Use `IsArray(x)` rather than `isdefined(x)`** for arrays — more precise about what you actually expect.
- **Always brace single-statement `if` bodies** — bare-line ifs are valid GSC but harder to scan; the project style consistently uses curly braces.
- **`[[ handler ]]()` is the GSC syntax for calling a function pointer.** You can pass `self` context with `entity [[ handler ]]( args )` — `entity` becomes the handler's `self`.

## `endon` / `notify` — the foundational signalling pattern

This is how scripts coordinate within themselves, beyond engine-named events.

### The shape

```c
// One thread declares a name for an event the entity will eventually emit:
self endon("end_game");      // "this function call terminates if 'end_game' is notified on self"
level endon("hellround_off"); // same idea, on level

// Another thread waits for a notify before continuing:
self waittill("done_loading");
level waittill("hellround_on");

// To fire the event, any thread can:
self notify("done_loading");
level notify("hellround_on");
level notify("end_game");
```

Three operations, all on a target entity (`level` or `self` or any entity):

- **`endon(name)`** — *terminate the current function call* if `name` is notified on the target entity. Place at the top of a `while(1)` thread to make it self-clean-up on a known signal. (Per the GSC language reference — verify in `docs_modtools/GSC_Language.pdf` for the exact semantics; "terminate" applies to the current invocation, not the entire script module.)
- **`waittill(name)`** — *block here until* `name` is notified. Resumes when it fires.
- **`notify(name)`** — *emit the named signal*. Wakes every thread `waittill`-ing on it and terminates every thread `endon`-ing on it.

### Passing values through `notify` / `waittill`

`notify` can pass any number of values, and `waittill` reads them by name:

```c
// notifier:
self notify("damage_taken", attacker, damage_amount);

// listener (in a thread that already endon'd or otherwise survives):
self waittill("damage_taken", attacker, damage_amount);
PrintLn("Hit by " + attacker.name + " for " + damage_amount);
```

Multi-arg notify/waittill is common when the signal carries information the receiver needs (entity references, integer state, strings).

### Useful `waittill_*` helpers from `util_shared.gsc`

The shared util library provides higher-level wait patterns that hand-rolled `waittill` chains can't elegantly express:

| Helper                                                | What it does                                                       |
| ----------------------------------------------------- | ------------------------------------------------------------------ |
| `util::waittill_any(s1, s2, s3, s4, s5, s6)`          | Wait for any one of N notifies (up to 6).                          |
| `util::waittill_any_array(notifies)`                  | Same, but takes an array of names — more flexible.                 |
| `util::waittill_any_return(s1, ..., s7)`              | Variant that returns *which* notify fired.                         |
| `util::waittill_any_timeout(timeout, s1, s2, ...)`    | Wait for any of the names, **or** a timeout. Returns which fired (or `undefined` on timeout). |
| `util::waittill_notify_or_timeout(msg, timer)`        | Single notify with timeout fallback.                               |
| `util::waittill_either(msg1, msg2)`                   | Wait for one of two specific notifies.                             |
| `util::waittill_multiple(...)`                        | Wait for **all** of the listed notifies (AND semantics).           |
| `util::waittill_any_ents(ent1, s1, ent2, s2, ...)`    | Wait across multiple entities (ent + name pairs).                  |
| `util::waittillmatch(name, value)`                    | Engine builtin (no `util::` prefix). Wait for a notify whose first param matches `value` — used heavily in animation notetracks (`entity waittillmatch("_anim_notify_", "gib_annihilate")`). |

When a state machine needs to react to several distinct signals or has timeout semantics, reach for these instead of hand-rolling.

### Canonical loop pattern

Almost every long-running thread in `zm_test` opens with `endon` lines for the lifecycle events that should kill it:

```c
function watch_thing()
{
    level endon("end_game");      // map ends → cancel this function
    self endon("disconnect");      // player leaves → cancel
    self endon("bled_out");        // player dies → cancel

    while (true)
    {
        self waittill("trigger");
        // ... do thing
    }
}
```

Without those `endon` lines, the loop keeps ticking after fast_restart / disconnect / death and you get apparent lag and ghost behaviour. This was the **missing `endon("end_game")` gotcha** from [`01-overview.md`](./01-overview.md). (Exact "terminate the current invocation" semantics — verify in `docs_modtools/GSC_Language.pdf` if it ever matters at the edges.)

### Real example — PaP cost modulation under bonfire sale

From `usermaps/zm_test/scripts/zm/_zm_pack_a_punch.gsc` (Reapy's override of the shipped PaP script):

```c
level endon("Pack_A_Punch_off");

while (1)
{
    self.cost = 5000;
    self.aat_cost = 2500;
    level waittill("powerup bonfire sale");

    self.cost = 1000;
    self.aat_cost = 500;
    level waittill("bonfire_sale_off");
}
```

Two `level` notifies (`"powerup bonfire sale"`, `"bonfire_sale_off"`) drive the cost-modulation state machine, and `level endon("Pack_A_Punch_off")` cleans up the whole thread when PaP shuts down.

## `level` vs `self` scope — where does `self` come from?

Callback handlers and notify-driven threads receive a `self` context that depends on the registration / notification path:

| Registration / call site                            | `self` inside the handler                                      |
| --------------------------------------------------- | -------------------------------------------------------------- |
| `callback::on_connect(&handler)`                    | The connecting **player**.                                     |
| `callback::on_spawned(&handler)`                    | The spawning **player**.                                       |
| `callback::on_ai_killed(&handler)`                  | The killed **AI / actor**.                                     |
| `callback::on_localclient_connect(&handler)` (CSC)  | The local **client** entity.                                   |
| `callback::on_finalize_initialization(&handler)`    | `level` (no per-entity context).                               |
| Custom `[[ handler ]]()` with no entity prefix      | Whatever `self` was in the caller's context.                   |
| `entity [[ handler ]]()`                            | `entity` becomes the handler's `self`.                         |
| `waittill` after `level waittill(...)`              | Whatever `self` was at the `waittill` line.                    |

Practical rules:

- **`self`** in player-event callbacks (`on_connect`, `on_spawned`, etc.) is the player. Use `self` to read/write that player's stuff.
- **`level`** is the project-wide singleton. Anything global (game state, level flags, world arrays) lives on `level.something`.
- A function declared as `function private foo() // self == X` (with the comment as the convention `zm_test` uses to label `self`) helps readers know what to expect.

## Common gotchas

- **Loop lacks `endon("end_game")`** — keeps running after `fast_restart`, causes apparent lag, doubled spawns. Always include the lifecycle endons at the top of a `while(1)`.
- **Wrong scope on `notify`** — `level notify("foo")` and `self notify("foo")` are *independent*; a thread `level endon("foo")` will not terminate when `self notify("foo")` fires (and vice versa). The notify target entity must match the endon target entity.
- **`waittill` race** — if the notify fires *before* you reach the `waittill` line, you miss it. For one-shot signals where ordering matters, use a flag (`flag::wait_till(...)` from `flag_shared`) instead — flags persist their state. Note: flags are themselves *implemented* on top of `notify` / `waittill`; they just keep a backing boolean so a `wait_till` issued after the flag was set returns immediately instead of blocking forever.
- **Late-joiner missing `on_connect` events** — `callback::on_connect` only fires for the joining player; existing players don't re-fire it. If a system needs all current players synced when one joins, the joiner's handler must scan and apply state, not assume initial setup.
- **Custom callback array uninitialised** — `[[ handler ]]()` on `undefined` crashes. Use `IsArray(level.your_callbacks)` to gate iteration, and `IsFunctionPtr(handler)` before storing.

## Reference reading

- `share/raw/scripts/shared/callbacks_shared.gsc` / `.csc` — the registry implementation; grep here for the full list of register functions.
- `share/raw/scripts/shared/laststand_shared.gsc` — registers `callback::on_laststand`.
- `share/raw/scripts/shared/flag_shared.gsc` / `.gsh` — the `flag::*` API for stateful signals (use this when `notify` race conditions matter).
- [`api-reference.md`](./api-reference.md) — `endon`, `notify`, `waittill`, `waittillmatch`, etc. are core engine builtins; signatures live in the GSCode JSON.

---

## Open questions / TODO

- [ ] Cover `flag::*` patterns end-to-end as a separate page — they're the persistent counterpart to `notify` and crop up in many state machines.
