# AI Pathing

> Stub — content pending.

How AI navigates the map: path nodes, traversals, cover, navmesh. The infrastructure is largely Radiant-side (you place the nodes); script reaches into it via the standard locomotion APIs.

## TODO

### Path nodes
- `path_node` entities — placement rules (line of sight, density, edges of geometry)
- Node connectivity — who-can-see-whom is the graph
- Node islands — disconnected node groups silently break AI
- Node placement debugging: `developer 1` + `path_node_show 1` (verify exact dvar) → in-game visualization

### Clip brushes for AI
- `clip_zombie` / `clip_ai` — restrict where zombies can pathfind
- `clip_player` — symmetric for player movement
- The "tower of clip" pattern — stacking clip brushes to fence off sky/ceiling

### Traversals
- `traversal_node` entity types: jump-up, jump-down, mantle (verify exact list against share/raw)
- Volume sizing rules — the start/end origins must satisfy AI navigation constraints
- Negotiation animation contract: anim plays from start origin to end origin, AI commits at the start
- Common failure: traversal placed near a corner where the AI can't get a clean line into it
- Cross-link to systems/traversals.md for the gameplay system view

### Cover
- `cover_*` nodes (verify which exist for zombies vs MP)
- Mostly for MP AI; zombies barely use cover

### Navmesh
- BO3's navmesh is generated from BSP geometry + clip brushes during Compile
- Reading navmesh debug visualization

## Common gotchas

- Stuck zombies = pathing dead end; usually a clip-brush gap or an unconnected node
- Zombies refusing to use a traversal = volume too thin / geometry obstruction / negotiation anim doesn't fit
- The first cerberus in `zm_test` was tweaked because dog spawn density was too hard to path through (commit history)
- Pathing only works inside `zone_volume` brushes that are currently active

## Cross-references

- [`brushes-patches.md`](./brushes-patches.md) — clip brushes are brushes with the right tool-texture
- [`brush-textures.md`](./brush-textures.md) — `clip_zombie` etc.
- [`zoning-portals-umbra.md`](./zoning-portals-umbra.md) — pathing is bounded by zone volumes
- Systems: [`traversals.md`](../systems/traversals.md) — the gameplay-system view of jumps/mantles
