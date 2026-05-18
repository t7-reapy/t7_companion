# Skybox

> **Status: stub.** The skybox is the geometry the engine renders at the world horizon. Referenced from SSI's `skyboxmodel` field; can be static or animated.

## What this page will cover

- Skybox model anatomy: it's a regular `xmodel`, not a special asset type. Usually a large sphere or cube with sky textures applied.
- Sourcing skybox models: shipped Treyarch ones (`skybox_zm_castle`, `skybox_default_day`, etc.), community packs (Deadshot's extras, BO2 skyboxes by Carabella/sgt. mckonow, T9 skyboxes by Nastian), or porting via Saluki.
- Wiring a skybox to an SSI record (`skyboxmodel "..."`).
- Animated skyboxes — rotation scripts, scrolling cloud textures, morphing day/night cycles.
- Sky size considerations — `worldfogskysize` on the fog record sets the dome size that wraps the skybox.
- Common gotchas: skybox model not zoned (visible black void at horizon), wrong scale (clipping into world geometry), seams visible at horizon.

## Reference reading

- [Epic Rotating World Script — Rotating Skybox (YouTube)](https://www.youtube.com/watch?v=F_S_mVl2VLo) — Gmzorz's rotating skybox script, demonstrates the animated case.
- See SSI's `skyboxmodel` field in [`asset-pipeline/05-ssi.md`](../asset-pipeline/05-ssi.md).
