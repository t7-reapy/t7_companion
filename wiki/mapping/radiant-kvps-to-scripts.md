# Radiant KVPs → Scripts

> Stub — content pending.

How values placed on entities in Radiant become readable from GSC / CSC. The contract is small but easy to get wrong; documenting the registration step that's the most common gap.

## TODO

### The standard KVPs

- `targetname` — the most common; identifies the entity for `GetEnt(name, "targetname")` and `GetEntArray(name, "targetname")`. Multiple entities can share a targetname; that's how `GetEntArray` works.
- `target` — points at another entity's targetname (entity-to-entity link).
- `script_noteworthy` — free-form string slot; convention is "second-order tag", e.g. distinguishing variants of an entity that share a targetname. (See `_rotating_object.csc` — uses it to pick the rotation axis.)
- `script_string` — same idea, different reserved key.
- `script_int`, `script_float` — numeric slots. (See `_rotating_object.csc` — uses `script_float` for the rotation period.)
- `script_parameters` — a parsed key=value blob. Less common; document the parser when we have a clean example.
- Custom KVPs — any other `key value` pair you set in Radiant; readable via `entity.<key>` *only after* the registration step.

### The registration step (the gap)

Custom KVPs (anything not in the standard list above) **don't appear on the entity by default**. They have to be registered in `bo3_kvp_db` (or whatever Reapy's source-of-truth file is — verify exact path and command). Without registration, `entity.my_custom_key` is `undefined` even though Radiant shows the value.

This is the single most-common "Radiant value isn't showing up in script" bug. Document the workflow:

1. Add the KVP in Radiant
2. Register it in `bo3_kvp_db`
3. Re-link / restart Radiant (verify which is needed)
4. Read `entity.<key>` from script

### The reading APIs

- `GetEnt(name, key)` — single entity by KVP value
- `GetEntArray(name, key)` — multiple entities matching a KVP value
- `entity.<key>` — direct field access for any standard or registered KVP
- Edge case: `localClientNum` parameter on the CSC equivalents (see `_rotating_object.csc` for the CSC pattern: `GetEntArray(localClientNum, "rotating_object", "targetname")`)

### Inspecting at runtime

- Standard pattern: `iprintln` the `entity.<key>` to confirm the value made it across.
- The `developer 1` dvar surfaces some KVP issues earlier; `developer 2` more.

## Cross-references

- Scripting [`api-reference.md`](../scripting/api-reference.md) for the canonical signatures of `GetEnt` / `GetEntArray`
- [`01-overview.md`](./01-overview.md) for where this fits in the Radiant ↔ script handshake
