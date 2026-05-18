# Sound Aliases

> Stub — content pending.

The addressable unit of sound in BO3. An alias is a named entry in a `.csv` file that points at a `.wav` and configures playback parameters (bus, volume curve, 2D vs 3D, rolloff, looping).

## TODO

### The alias `.csv` format
- Where alias CSVs live (verify exact path under `share/raw/sound/aliases/` or similar)
- Column reference: name, file path, channel/bus, volume curve, 2D vs 3D, looping vs non-looping, rolloff distances, randomization fields
- Authoring tool: Acoustix (ardivee) — opens and edits these CSVs in a UI
- Compiler: Harmony (Scobalula) — `.csv` → engine-readable

### Calling an alias from script
- `entity PlaySound("alias_name")` (server-side, replicated)
- `entity PlaySoundOnTag("alias_name", "tag_name")` (server-side, attached to a model bone)
- `PlaySoundAtPosition("alias_name", origin)` (server-side at a world point)
- `PlayLoopedSound*` family for looping
- `entity StopSound* / StopAllSounds` for cleanup
- CSC equivalents (`localClientNum` parameter)

### Real example
`zm_room_of_thanks_elevator.gsc` uses dedicated *sound entities* that ride along with moving geometry:
```cpp
self.snd_ent_left_door PlaySoundOnTag(door_sound, "tag_origin");
self.snd_ent_grid_door PlaySoundOnTag(ELEVATOR_SOUND_GRID, "tag_origin");
```
Why a separate entity per sound: positional audio tracks the `script_model`'s position; if you parent the sound to the *door brush*, the listener hears it at the door's centroid. Often you want a different anchor — hence the dedicated `tag_origin` script_model that you `MoveTo` together with the visual.

### Bus routing
- `BUS_FX` / `BUS_MUS` / `BUS_VOX` / `BUS_REVERB` (verify list)
- Player slider mapping
- Per-bus volume control from script (verify API)

### Rolloff and 2D/3D
- 2D = always at full volume regardless of position (HUD beeps, announcer, music)
- 3D = positional, falls off with distance per the rolloff curve

## Common gotchas

- WAV not 48 kHz / 16-bit PCM → silent failure
- Forgetting to recompile aliases after editing the CSV → script calls return silently
- Calling `PlaySound` from CSC and forgetting it'll only emit on that local client → splitscreen mismatch
- Looping sound without a `Stop*` partner → audio leak, persists across rounds

## Cross-references

- [`01-overview.md`](./01-overview.md) — where aliases sit in the data flow
- [`formats-and-encoding.md`](./formats-and-encoding.md) — the WAV constraints
- [`splitscreen-audio.md`](./splitscreen-audio.md) — the local-vs-replicated emission gotcha
- [`tools.md`](./tools.md) — Acoustix + Harmony
