# Weather effects

> **Status: stub.** Rain, thunder, lightning, wind — the visual layer on top of fog/SSI state changes that sells "this is a stormy environment." Reapy iterated heavily on this for `zm_test`.

## What this page will cover

- Rain: 3D rain volumes (`volume_weathergrime`, particle FX), camera-screen rain droplets, drainpipe overflow, rain audio.
- Thunder: lightning flash via SSI swap (the `house_thunder` lighting state), `LightningStrikeEffects` engine system, audio sync.
- Wind: wind FX exploders, wind sound layering, animated foliage response.
- The state machine that ties weather phases together with the hellround state.
- Splitscreen considerations (audio one-shot vs per-player) — same patterns as fog.
- Coupling to fog states (`volume_worldfog` + `volume_litfog` swaps during storm peak).
- Performance: weather is a busy zone for FX exploders — culling distance matters.

## Reference reading

- [Thunder & Lightning tutorial (YouTube)](https://www.youtube.com/watch?v=D1HjKRy1RzA) — practical setup, BO3-specific.
- [LightningStrikeEffects (zeroy wiki)](https://wiki.zeroy.com/index.php?title=Call_of_duty_bo3:_LightningStrikeEffects) — engine-side reference for the lightning system.
- [`Scobalula/Bo3OnScreenRainDrops`](https://github.com/Scobalula/Bo3OnScreenRainDrops) — installed in `zm_test`; renders rain droplets directly on the player's camera.

## Status

This is one of Reapy's most-iterated systems. Lots of original signal here when we get to writing it: rain pipe authoring, the wind exploder loop, lightning state coupling, splitscreen sound deduplication.
