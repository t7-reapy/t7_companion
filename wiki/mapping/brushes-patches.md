# Brushes & Patches

> Stub — content pending.

The two geometry primitives Radiant exposes: **brushes** (convex 3D solids) and **patches** (curved/planar surfaces with bezier control points). Most BSP geometry is brushes; patches are reached for organic shapes brushes can't do.

## TODO

### Brushes
- Convex requirement and what happens when you violate it (slivers, light leaks)
- Subtraction: when the CSG (Constructive Solid Geometry) operation actually helps vs. when it makes a mess
- Carving — practical use vs. risk
- Cutting brushes (clipper tool) — the workhorse
- Free scaling vs constrained scaling — when each is OK
- Brush textures live on a separate page: [`brush-textures.md`](./brush-textures.md)

### Patches
- Primitive vs curve patches — when each shape is right
- Vertex editing, density, free scaling
- Patches don't block AI by default — clip brushes still needed (link to AI pathing)
- Texture blending across patch seams
- Patch density vs runtime cost

### Detail prop placement
- See [`terrain-detailing.md`](./terrain-detailing.md) for the detail/decal/prop layer that sits on top of base geometry.

## Common gotchas

- Light leaks at brush slivers — almost always a non-convex brush or two brushes that almost-but-not-quite meet
- Bad collision after cutting — collision uses simplified geometry, can drift from the visual
- Patch density too high — VRAM + perf; resist the urge to crank it past what the camera can resolve
- Free scaling a brush off-grid — the brush vertices snap to grid silently and the visual changes (see also Radiant grid sizes)

## Cross-references

- [`brush-textures.md`](./brush-textures.md) — tool-textures (caulk, nodraw, clip, hint, skip)
- [`zoning-portals-umbra.md`](./zoning-portals-umbra.md) — how portals interact with brush geometry
- [`ai-pathing.md`](./ai-pathing.md) — clip brushes for AI navigation constraints
- [`prefabs.md`](./prefabs.md) — bundling brush/patch assemblies for reuse
