# Recipe — Animations

> Stub — content pending. Paused recipe; resume when YouTube-corpus indexing is finalized.

How to play, scale, and chain animations from script: scripted models, doors, machines, NPCs, and cinematics.

## TODO

- `playanimation` / `setanim` / `useanimtree` — engine call surface and what each accepts
- Animstate vs anim alias vs raw animation reference
- Animation notetracks (the `notify` events embedded in the anim) — how to listen with `waittill`
- Looping vs one-shot, blending, root motion handling
- Door / machine / lid animations — the typical pattern for scripted props
- Cross-link to [callbacks](./04-callbacks.md) for `notify` plumbing
- Cross-link to [models](../asset-pipeline/03-models.md) for the asset side (xanim / xmodel pairing)

## Research starting points

- YouTube corpus: search `metadata` for `pose-manipulator` and titles mentioning "animation".
- t7_api_*.json (GSCode) for the exact builtin signatures.
