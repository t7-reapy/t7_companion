# Visionsets & LUTs

> **Status: stub.** Post-FX colour grading and tonemapping. Frequently swapped *together with SSI* per lighting state — a hellround LUT, a thunder LUT, a normal-day LUT.

## What this page will cover

- What a visionset is — `.vision` files in `share/raw/vision/`, the engine-side post-FX stack.
- Visionset slots on `volume_lut` and similar entities (per-region LUT regions).
- LUT (look-up table) authoring: the `LutTpage` image semantic, how to bake a colour-grade in Photoshop / DaVinci / Substance and import it.
- Coupling visionsets to SSI swaps so a state change rotates the colour grade alongside the lighting.
- Examples of state-driven LUT swaps in `zm_test` (hellround, thunder, default day).
- Performance: visionsets are mostly free; LUTs are tiny.
- Common gotchas: wrong colour space on the LUT image, blending two visionsets producing a muddy intermediate.

## Reference reading

- [Visionsets and Overlays Tutorial (zeroy / Slick Willy, YouTube)](https://www.youtube.com/watch?v=D9Uq0zo-NrU) — the canonical visionsets + LUT walkthrough.
- [Change Look of Map / LuT Edit (YouTube)](https://www.youtube.com/watch?v=p6iXtx50bxI) — practical LUT authoring.
- See related fields on the SSI page ([`asset-pipeline/05-ssi.md`](../asset-pipeline/05-ssi.md)) — `lensFlare`, sun cookie, exposure that interact visually with the chosen visionset.
