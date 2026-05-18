# Audio

> **Status: stub.** Top-level discipline. *Mostly research-driven for now* — Reapy hasn't gone deep here, doesn't know what `sharedweaponsounds` / `surfacesounddef` etc. actually are. We'll fill this in primarily by reading docs + community wikis rather than from lived experience.

## What this section will cover

- Sound alias system: `.csv` alias files, the alias columns (name, file path, channel/bus, volume curve, 2D vs 3D, looping vs non-looping, rolloff distances).
- Audio formats: WAV, 48 kHz, signed 16-bit PCM, stereo or mono — anything else silently fails or crashes (per `KNOWLEDGE.md`).
- Streamed vs loaded sounds: when to use which, performance trade-offs.
- Bus routing: `BUS_FX`, `BUS_MUS`, `BUS_VOX`, etc.
- Sound zones / reverb: `volume_litfog` is *not* sound zones; sound zones use a different volume entity. Documentation TODO.
- Ambient rooms (Ardivee's tool).
- Music management: round music, hellround music, perk jingles, easter-egg music.
- VOX (voice-over) systems: announcer, zombie VOX, kortifex announcer pack.
- Splitscreen positional audio: `GetLocalPlayers()` patterns, host-only emit deduplication (the rain-sounds example from `zm_test`'s known bugs).
- Royalty-free sourcing: Pixabay djent tracks worked except for false-positive Content ID claims.
- Tools: Acoustix, Sound Studio Extended, Harmony, WAV converters.
- The MKV in-map video format quirks (specific Handbrake version requirement) — also lives here adjacent to audio.

## Sub-pages (planned)

- `01-overview.md`
- `sound-aliases.md`
- `formats-and-encoding.md`
- `ambient-and-reverb.md`
- `music.md`
- `vox.md`
- `splitscreen-audio.md`
- `tools.md`

## Reference reading

- [Acoustix (Sound Editor) — ardivee wiki](https://wiki.ardivee.com/article/acoustix/) — the BO3-integrated sound editor.
- [`Scobalula/Harmony`](https://github.com/Scobalula/Harmony) — sound alias compiler tool.
- [`tovaru/Black-Ops-II-Sound-Studio-Extended`](https://github.com/tovaru/Black-Ops-II-Sound-Studio-Extended) — sound alias editor.
- [Cyph3r — localization editor (ardivee wiki)](https://wiki.ardivee.com/article/cyph3r-localization-editor/) — adjacent (localized strings often pair with localized VO).
- [Adding Environment / Ambient Sounds (YouTube)](https://www.youtube.com/watch?v=EKEyCStNKTw)
- [Ambient Rooms / Reverb (YouTube)](https://www.youtube.com/watch?v=OgxJKjQZmxw)
- [Ambient Rooms (Modme Wiki)](https://wiki.codmods.com/docs/bo3/ambient_rooms)

## Status

Audio is a **gap, not a depth** for Reapy. When we eventually write it, expect content to be reconstructed from docs / community wikis / the script source rather than from project commits. Splitscreen sound bugs (per the README known-issues list) are the one place lived experience exists.
