# Clientfields recipes — server↔client state propagation

> A **clientfield** is a typed, bit-budgeted, automatically-replicated value that the server (GSC) sets and every client (CSC) sees. The cleanest way to push gameplay state from server logic to per-client visual / audio / HUD reactions. Used pervasively in `zm_test` for hellround toggles, weather state, AI VFX, fog swaps. This page is the protocol reference + recipe collection.

## Mental model

Picture a **typed bit-pipe** from server to clients, attached to a target entity-type:

```
GSC (server)                              CSC (client, per local player)
──────────────                            ──────────────────────────────
clientfield::register(...)                clientfield::register(..., &handler, flags)
                                          ^ same name, version, bits, type
                                          ^ CSC also declares a handler

clientfield::set(name, value)   ───────►  on next propagation tick:
clientfield::increment(name)              handler(client, oldVal, newVal, ...) fires
                                          on every connected client (including late-joiners
                                          if the right flag is set)
```

Why this exists (instead of just calling a CSC function from GSC): **you can't call CSC from GSC**. The two scripting VMs don't share a stack, can't pass function pointers across. Clientfields are the engine's typed bit-channel for crossing the boundary safely — small, replicated, persistent across reconnects, integrated with the network layer.

The **engine itself** wires the propagation — your script just declares the field, sets it on the server, and writes a handler on the client.

## Registration — both sides must declare

A clientfield exists only if **both** sides register it with matching name / version / bits / type. If GSC registers with 2 bits but CSC registers with 1, you get a **runtime "clientfield mismatch"** error — debuggable but not caught at compile time. See the dedicated section [The "clientfield mismatch" runtime error](#the-clientfield-mismatch-runtime-error) below for the full diagnosis-and-fix playbook.

### GSC registration (5 args)

```c
clientfield::register(str_pool_name, str_name, n_version, n_bits, str_type);
```

No handler, no flags — server-side is just declaring the field exists and naming the bit budget it'll consume.

### CSC registration (8 args)

```c
clientfield::register(str_pool_name, str_name, n_version, n_bits, str_type,
                      &handler_function,
                      host_only_flag, callback_zero_on_new_ent_flag);
```

CSC adds three extra args:

- **`&handler`** — function pointer called when the field changes value on this client.
- **`host_only_flag`** — pass `CF_HOST_ONLY` (true) if only the splitscreen host should run the handler, or `!CF_HOST_ONLY` (false, the common case) for every local client.
- **`callback_zero_on_new_ent_flag`** — pass `CF_CALLBACK_ZERO_ON_NEW_ENT` (true) to **fire the handler on late-joiners and entity-spawn**. This is the late-joiner-fix flag — without it, mid-game-join clients miss the current state. With it, the engine fires the handler with `bNewEnt=true` and the current value as soon as the client sees the entity for the first time.

The two flag macros are defined as plain `true` in `share/raw/scripts/shared/version.gsh:102-103`:

```c
#define CF_HOST_ONLY                  true
#define CF_CALLBACK_ZERO_ON_NEW_ENT   true
```

Convention: write `CF_HOST_ONLY` when you want it on, `!CF_HOST_ONLY` when off. Keeps the call-sites readable about *what* you intended without having to remember which flag is `true` and which is `false`.

## Target types (the `str_pool_name`)

The first argument names which **entity type** the clientfield is attached to. Each type has its own **bit pool**, and registered fields consume bits from that pool. Free-bit numbers below are for a **starter ZM map** — you have fewer if shipped Treyarch scripts have already consumed some:

| Pool             | CSC `self`        | Free bits (starter ZM) | Used in zombies? | General usage                                                                                                            |
| ---------------- | ----------------- | ---------------------:| ---------------- | ------------------------------------------------------------------------------------------------------------------------ |
| `world`          | `level`           | **1 917**             | Yes              | Huge pool for global state — manipulate almost any client-side entity. Useful for LUA UI toggles too.                    |
| `actor`          | AI                | 59                    | Yes              | Per-AI CSC FX/sounds on basic zombies and other actors.                                                                  |
| `vehicle`        | vehicle entity    | 76                    | Sometimes        | GK drones, SOE wasps, Origins tank. CSC FX/sounds on vehicles.                                                           |
| `allplayers`     | player            | 107                   | Yes              | Fires on **every** player. Common for 3rd-person FX where you filter the local player from seeing it on themselves.      |
| `toplayer`       | player            | 79                    | Yes              | Fires on **a single** player. Per-local-player effects (XcDylan93's `dw_fired`, `lh_fired`).                              |
| `playercorpse`   | dead player body  | 112                   | No               | FX on dead-player bodies (mostly MP).                                                                                    |
| `clientuimodel`  | N/A               | 79                    | Yes              | Direct interface to the client's LUA UI Model. No CSC handler functions needed, but still requires CSC registration.     |
| `scriptmover`    | basic entity      | 84                    | Yes              | Client-side FX on a server-side entity, especially when it's moving.                                                     |
| `helicopter`     | helicopter entity | 36                    | No               | FX/sounds for helicopters and drones.                                                                                    |
| `plane`          | plane entity      | 63                    | No               | Plane assets (rare in BO3).                                                                                              |
| `missile`        | missile entity    | 56                    | Sometimes        | FX/sounds on missiles from launchers, bows, AI weapons.                                                                  |
| `zbarrier`       | zbarrier entity   | 44                    | Yes              | FX/sounds on zbarriers — Pack-A-Punch, GobbleGum machines, zombie barricades.                                            |
| `item`           | item entity       | 64                    | No               | Rare campaign-specific scenarios.                                                                                        |

Notes:

- The **free-bit count is per-pool** — fields on `world` don't compete with fields on `actor`. The "1 917 free in `world`" sounds generous, but stack enough wide (4–8 bit) world fields and you'll notice.
- **Running out of bits in a pool is a hard limit** — registration fails. Reduce widths (`GetMinBitCountForNum`), consolidate toggles into one wider field, or move the state out of clientfields entirely.

## Field types (the `str_type`)

The fifth arg names how to interpret the bits the engine replicates:

| Type        | Semantics                                                         | Frequency in shipped Treyarch source |
| ----------- | ----------------------------------------------------------------- | ------------------------------------:|
| `"int"`     | Plain integer in the bit-width you declared. Use for state values (0/1 toggle, 0/1/2/3 phase, etc.). | ~570 registrations (GSC + CSC combined) |
| `"counter"` | Each `increment(name)` bumps the value by 1; on natural rollover (when it overflows the declared bit-width) the wrap is invisible to the handler — what matters is that the CSC handler fires once per increment. Effectively a "fire-and-forget pulse." Use for one-shot events (FX bursts, sound triggers). | ~67 |
| `"float"`   | Quantised float in the **0.0 – 1.0 range** (it's a normalised fraction, not an arbitrary float). Used for UI / progress-bar / healthbar / fill-percentage values where a smooth gradient between 0 and 1 needs to replicate. Grep `t7-source` for `"float"` in `clientfield::register` calls to see the real shipped uses. | 2 (CSC only) |

**Choose `int` vs `counter` carefully**:

- **`int`** holds *state*. Setting it to the same value twice does nothing. Use for "is hellround on or off?" — the CSC handler reacts to the *transition*.
- **`counter`** is a *pulse*. Every increment fires the handler once on every client. Use for "play this FX again now" — even if it just played, increment again to play again.

Picking the wrong type leads to "the FX only plays the first time" (you used int when you needed counter) or "the state keeps re-firing my expensive setup" (you used counter when you needed int).

`GetMinBitCountForNum(N)` is a useful helper: returns the minimum bit-width to hold values 0..N. Used in shipped scripts: `clientfield::register("actor", FURY_DAMAGE_CLIENTFIELD, VERSION_DLC4, GetMinBitCountForNum(7), "counter")` for a 0..7 counter.

## Operations — set / increment

From GSC (or CSC, but typically GSC):

```c
// Set a value (most common — works on any-side, any-target):
target_entity clientfield::set(str_field_name, n_value);
level clientfield::set("hellround_on", 1);
self clientfield::set("perk_halo_color", 3);

// Increment a counter (fires the CSC handler once per call):
target_entity clientfield::increment(str_field_name);
target_entity clientfield::increment(str_field_name, by_n);   // bump by N at once
```

The first arg of `set` (or implicit `self` of the call) is the entity the clientfield is registered on. So `level clientfield::set("hellround_on", 1)` requires the field to have been registered on the `"world"` (or `"level"`) pool; `self clientfield::set(...)` requires the field to be on whatever pool matches `self`'s type.

## Handler signature (CSC side)

The handler runs on every client that matches `host_only_flag`, fired by the engine when the value changes (or on entity spawn if `CF_CALLBACK_ZERO_ON_NEW_ENT` is set). Standard signature:

```c
function on_my_field_change( n_local_client_num, _oldVal, n_new_val, b_new_ent,
                              _bInitialSnap, _fieldName, _bWasTimeJump )
{
    util::waitforclient( n_local_client_num );

    if ( b_new_ent )
    {
        // First time we're seeing this entity (e.g. late-joiner sync).
        // Apply current state without playing the "transition" FX.
    }
    else
    {
        // Value changed during normal gameplay — fire the transition reaction.
    }
}
```

Conventional argument names (from shipped Treyarch scripts):
- `n_local_client_num` — which local client this handler is firing for. Pass through to all `PlayFx` / `PlaySound` / etc. calls.
- `n_new_val` — the new value (the only one you usually care about).
- `b_new_ent` — true if this is a "first sight" snapshot vs a value change (the late-joiner discriminator).
- The other args (`_oldVal`, `_bInitialSnap`, `_fieldName`, `_bWasTimeJump`) are usually unused — prefix with `_` to signal that.

`util::waitforclient( n_local_client_num )` is the standard first call: blocks until the engine confirms the client is fully ready to render/sound, avoiding race conditions with the spawn process.

## Real worked example — hellround environment toggle

The cleanest end-to-end clientfield in `zm_test`. Files: `usermaps/zm_test/scripts/zm/hellround/zm_hellround_environment.{gsc,csc,gsh}`.

**GSH header** declares the constant:

```c
#define HRENV_TOGGLE_CLIENT_FIELD "hr_env_toggle"
```

**GSC side** registers with 1 bit `"int"`, and sets when hellround toggles:

```c
clientfield::register("world", HRENV_TOGGLE_CLIENT_FIELD, VERSION_SHIP, 1, "int");

// later, when toggling hellround on/off:
level clientfield::set(HRENV_TOGGLE_CLIENT_FIELD, b_enable);
```

**CSC side** registers with the matching shape + a handler + the **`CF_CALLBACK_ZERO_ON_NEW_ENT`** flag (so late-joiners get the right state):

```c
clientfield::register("world", HRENV_TOGGLE_CLIENT_FIELD, VERSION_SHIP, 1, "int",
                      &hellround_environment, !CF_HOST_ONLY, CF_CALLBACK_ZERO_ON_NEW_ENT);

function hellround_environment(n_client_num, _oldVal, n_new_val, b_new_ent, ...)
{
    util::waitforclient(n_client_num);

    if (!b_new_ent)
    {
        // Live toggle — play transition FX/sound.
        play_transition_fx(n_client_num);
        play_transition_sounds(n_client_num);
    }
    fog_update(IS_TRUE(n_new_val), b_new_ent);
    show_hellround_volumes(IS_TRUE(n_new_val));
    show_hellround_models(IS_TRUE(n_new_val));
    play_environment_sounds(n_client_num, IS_TRUE(n_new_val));
}
```

Why this works for late-joiners: the `CF_CALLBACK_ZERO_ON_NEW_ENT` flag tells the engine to fire `hellround_environment` *with the current value and `b_new_ent=true`* the moment a client first sees the world. The handler skips the transition FX (since the client wasn't there for the toggle) but still applies the state (fog, volumes, models, ambient sounds). Mid-game-join into hellround → the joining client correctly sees bloody fog, hellround models visible, ambient sound playing.

## Patterns

### Counter for one-shot FX events

Want to spawn an FX every time something happens, reliably across clients?

```c
// GSC
clientfield::register("scriptmover", "play_burst_fx", VERSION_SHIP, 1, "counter");
self clientfield::increment("play_burst_fx");   // every call fires the CSC handler once

// CSC
clientfield::register("scriptmover", "play_burst_fx", VERSION_SHIP, 1, "counter",
                      &on_play_burst_fx, !CF_HOST_ONLY, !CF_CALLBACK_ZERO_ON_NEW_ENT);
function on_play_burst_fx(client, _o, _n, _new, ...)
{
    util::waitforclient(client);
    PlayFxOnTag(client, level._effect["my_burst"], self, "tag_origin");
}
```

Note `!CF_CALLBACK_ZERO_ON_NEW_ENT` — counter pulses don't replay on late-join (you don't want the FX to fire when a player joins mid-game just because the counter happens to be at 7).

### Int for sticky state

Want to track and react to "what phase of hellround are we in" (0/1/2/3)?

```c
// GSC
clientfield::register("world", "hellround_phase", VERSION_SHIP, 2, "int");
level clientfield::set("hellround_phase", n_iteration);

// CSC
clientfield::register("world", "hellround_phase", VERSION_SHIP, 2, "int",
                      &on_phase_change, !CF_HOST_ONLY, CF_CALLBACK_ZERO_ON_NEW_ENT);
function on_phase_change(client, oldVal, newVal, bNewEnt, ...)
{
    util::waitforclient(client);
    apply_phase_visuals(client, newVal, bNewEnt);
}
```

Note `CF_CALLBACK_ZERO_ON_NEW_ENT` for state — late-joiner gets the current phase.

### Per-player state via `"toplayer"`

When the state is per-local-player (not world or AI), use the `"toplayer"` pool. Example from XcDylan93's `_zm_mod_fx.csc`:

```c
clientfield::register("toplayer", "dw_fired", VERSION_SHIP, 1, "counter",
                      &dw_fired, !CF_HOST_ONLY, !CF_CALLBACK_ZERO_ON_NEW_ENT);
clientfield::register("toplayer", "lh_fired", VERSION_SHIP, 1, "counter",
                      &lh_fired, !CF_HOST_ONLY, !CF_CALLBACK_ZERO_ON_NEW_ENT);
```

`"toplayer"` fields fire the handler on the *target* local client only — useful for per-player visual feedback (their own muzzle flash colour, their own kill streak FX) that other splitscreen players shouldn't see double.

## The "clientfield mismatch" runtime error

Your most likely first encounter with clientfields is this error. It fires at runtime when the GSC and CSC `clientfield::register` calls don't agree on the field's shape.

**Triggers**: any of `str_pool_name`, `str_name`, `n_version`, `n_bits`, or `str_type` differs between the two sides for the same name.

**Surface as**: console error message + (with developer mode enabled) an in-game popup. To get diagnostics:

- `+set developer 2 +set logfile 2` — generic developer-mode flags that surface most errors with stack-trace context. (Detail in [`01-overview.md`](./01-overview.md#5-silent-by-default-failures-until-you-turn-the-lights-on).)
- `+set com_clientfieldsdebug 1` — clientfield-specific debug dvar. Adds extra diagnostic logging when registrations are evaluated and when fields are set/incremented.

**Fix**: open both files (GSC and CSC) and diff the registration lines for the named field. Make sure all 5 of the shape args match exactly, including the `VERSION_*` constant.

## Common pitfalls

- **GSC and CSC registrations don't match.** See the dedicated section above — runtime "clientfield mismatch" error.
- **Forgot `CF_CALLBACK_ZERO_ON_NEW_ENT` on a state field.** Mid-game-join client doesn't get the current state, sees the default until the next change. Add the flag for `int` fields that represent persistent state. Skip it for `counter` pulses you don't want to retroactively re-fire.
- **Using `int` when you wanted `counter`** (or vice versa). Setting an `int` to the same value twice does nothing on the handler side. With `counter`, every `increment` fires the handler once — and the value wrap on overflow is invisible to the handler, so increments keep working past the bit-width cap. If your FX "only plays once," you probably want `counter`.
- **Bit budget exhausted in a target pool.** Registration fails. Reduce bit widths where possible (`GetMinBitCountForNum`), consolidate multiple toggles into one wider field with bit-flags, or move state to script variables that don't need replication. Per-pool free-bit table is in the Target types section above.
- **Forgot `util::waitforclient(n_client_num)` at the top of a CSC handler.** Race conditions with spawn FX, audio, anything that needs the client to be ready.
- **Setting from CSC.** Clientfields propagate *server → client*. Setting from CSC affects only that client's local state and doesn't propagate.
- **Using `self` ambiguously.** `self clientfield::set(...)` requires `self` to match the pool's target type. Inside a player function you can `self clientfield::set("perk_halo", 1)` if `"perk_halo"` was registered on `"allplayers"` or `"toplayer"`, but not if it was registered on `"world"`.

## Reference reading

- `share/raw/scripts/shared/clientfield_shared.gsc` / `.csc` — the protocol implementation; grep here when in doubt about behaviour.
- `share/raw/scripts/shared/clientfields.gsh` — header declarations.
- `share/raw/scripts/shared/version.gsh:102+` — `CF_*` flag macros.
- The "clientfield mismatch" error is documented inline above — see "The 'clientfield mismatch' runtime error" section.
- `usermaps/zm_test/scripts/zm/hellround/zm_hellround_environment.{gsc,csc,gsh}` — the worked example end-to-end.
- `share/raw/scripts/shared/ai/archetype_apothicon_fury.gsc` / `.csc` — extensive use of `actor`-pool clientfields for per-AI VFX (`FURY_DAMAGE_CLIENTFIELD`, `FURY_FURIOUS_MODE_CLIENTFIELD`, etc.).
- [`XcDylan932/ReDempTion`'s `_zm_mod_fx.csc`](https://github.com/XcDylan932/ReDempTion/blob/main/scripts/zm/_xcdylan93/_zm_mod_fx.csc) — `"toplayer"` pool examples for per-local-client effects.

---

## Open questions / TODO

- [x] **Bit-budget table per pool** — resolved inline (see Target types). Sourced from `KNOWLEDGE.md`.
