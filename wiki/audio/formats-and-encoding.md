# Audio Formats & Encoding

> Stub — content pending. Mostly grounded in `KNOWLEDGE.md`'s constraints — fill out with concrete tool commands when the page is reviewed.

What BO3 will actually load. Get this wrong and the engine fails silently or crashes — the worst combination.

## TODO

### WAV requirements (from `KNOWLEDGE.md`)
- **48 kHz sample rate** — anything else fails
- **Signed 16-bit PCM** — no float, no 24-bit, no compressed PCM
- **Mono or stereo** — surround channels not supported
- **Anything else: silent fail or crash** (per Reapy's investigation)

### Streamed vs loaded
- **Loaded**: short SFX (footstep, gunshot, perk drink). Loaded into memory at level start; instant playback; counts toward asset budget.
- **Streamed**: long-form audio (music, ambience). Read from disk on play; cheap memory; small streaming overhead per active stream.
- The choice is per-alias (verify the CSV column name)
- Rule of thumb: streaming above ~10 seconds, loaded below

### MKV in-map video (related but not audio)
- BO3's in-map video format is `.mkv`
- Specific Handbrake version + encoding settings required (per `KNOWLEDGE.md` — get the exact version)
- Used for: room-of-thanks credit videos, intro cinematics, screens
- Lives in this section as adjacent to audio in the build pipeline

### Conversion workflow
- Source → 48 kHz / 16-bit / PCM via [tool TBD: Audacity? ffmpeg? Acoustix?]
- ffmpeg one-liner: `ffmpeg -i in.<any> -ar 48000 -sample_fmt s16 -ac 2 out.wav` (verify Reapy's exact command)
- Batch conversion script for whole sound packs (TODO: write one and link it)

### Royalty-free music sourcing
- Reapy uses **Pixabay** for djent/metal tracks (per `PRESENTATION_PLAN.md`)
- Specific artists: Void Construct, Dark Matter, Melodic Metal, Deathcore, Sinister; AlexGrohl, AudioDollar, BrightestAvenue
- One Pixabay-sourced track triggered a false-positive YouTube Content ID claim (commit history: replaced for safety)
- Workflow: vet on YouTube Content ID before committing to a track

## Common gotchas

- 44.1 kHz file (most music tracks default) → silent fail in-game; convert to 48 kHz first
- 32-bit float WAV (some DAWs default) → fail; convert to 16-bit PCM
- MKV with the wrong Handbrake version → won't play; engine doesn't say why

## Cross-references

- [`sound-aliases.md`](./sound-aliases.md) — the alias points at the WAV
- [`tools.md`](./tools.md) — the conversion tools
- `KNOWLEDGE.md` — Reapy's investigation notes on format constraints
