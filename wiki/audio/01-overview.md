# Audio — Overview

> Stub — content pending. Audio is the section Reapy is shallowest on; expect this page to be reconstructed mostly from community sources rather than from lived experience. Use [`/corpus-challenge`](../CONTRIBUTING.md) once it exists.

The big picture for the audio subsystem in BO3: how a `.wav` file becomes a sound a player hears.

## TODO

### The data flow
- Author: WAV file (48 kHz, signed 16-bit PCM, stereo or mono — anything else fails silently per `KNOWLEDGE.md`)
- Declare: row in a sound-alias `.csv` mapping a name to the file + playback parameters
- Compile: alias compiler (Harmony) bakes the `.csv` into the engine-readable form
- Reference: GSC / CSC calls `PlaySoundAtPosition`, `PlayLoopedSoundAtPosition`, `PlaySoundOnTag`, `entity PlaySound(alias)` etc. by alias name

### The layered systems
- **Sound aliases** — the addressable unit. Per-alias config: bus, volume curve, 2D vs 3D, looping, rolloff distances. (See `sound-aliases.md`.)
- **Buses** — `BUS_FX`, `BUS_MUS`, `BUS_VOX`, `BUS_REVERB` (verify exact list). Player volume sliders mix buses.
- **Streamed vs loaded** — short SFX get loaded; music and ambience stream. See `formats-and-encoding.md`.
- **Reverb / ambient rooms** — region-based audio post (Ardivee's Acoustix tool). See `ambient-and-reverb.md`.
- **Music** — round music, hellround music, jingles, easter-egg music — multi-layer. See `music.md`.
- **VOX** — announcer + per-AI voice. See `vox.md`.
- **Splitscreen** — host-only emit deduplication; the rain-sounds gotcha. See `splitscreen-audio.md`.

### Server vs client
- Most sound calls run **client-side (CSC)** — sound is local to the receiver
- Server-side calls (`entity PlaySound(alias)` from GSC) replicate to all clients automatically
- The choice matters for splitscreen — see `splitscreen-audio.md`

### Adjacent: in-map video
- BO3 supports `.mkv` video on materials/screens; specific Handbrake version + encoding constraints (per `KNOWLEDGE.md`)
- Lives in `tools.md` for now

## Cross-references

- All sub-pages in this section
- [`asset-pipeline/01-overview.md`](../asset-pipeline/01-overview.md) — where sound assets fit in the build
- `KNOWLEDGE.md` — Reapy's hands-on notes on WAV constraints, MKV constraints
