# Traversals

> Stub — content pending.

The "jump up / jump down / mantle" system that lets zombies cross terrain features they can't navigate via normal pathing. Lives at the boundary between mapping (placing the volumes) and AI (the negotiation animation).

## TODO

- Traversal entity types in Radiant: jump-up, jump-down, mantle (verify exact list)
- The negotiation animation contract: what AI plays it, when it commits, when it bails
- Volume sizing rules so zombies actually pick the traversal
- Why traversals fail (typical: obstruction, mismatched start/end heights, missing path nodes nearby)
- Pathing interaction: does the traversal need a connected path node on each side?
- Custom traversal animations (if possible — verify)

## Research starting points

- YouTube corpus: `theme` includes `traversals`; metadata `zombie-jumps`.
- See also: planned [pathing](../mapping/) page (mapping section).
