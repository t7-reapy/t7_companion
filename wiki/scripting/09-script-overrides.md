# Recipe — Script Overrides & Merge Workflow

> Stub — content pending.

The mechanism for replacing a Treyarch shared script with your own version, and the mostly-manual workflow for keeping your override in sync when the underlying source updates (or, more relevantly, when you copy from a community pack that diverges from stock).

## TODO

### How the override works (the engine side)
- Any script in `share/raw/scripts/zm/<file>.gsc` can be **shadowed** by placing a file with the same relative path in `usermaps/zm_test/scripts/zm/<file>.gsc`
- The linker prefers the usermap version
- **Relative paths matter** — must mirror the share path exactly
- Some scripts also need their reference removed from `zm_patch.csv` first (verify which / when — Reapy hit this empirically, undocumented)
- Discovered through trial and error, not in any official doc

### When to override vs when to add a new file
- Override if: you need to *change* an engine behavior that's hardcoded in the original
- Add new file if: you need *additional* behavior that doesn't conflict with the original
- Most map-side work should be additive; overrides are a last resort

### The merge workflow
- Stock script + your changes vs. community pack's stock + their changes vs. your wanted state
- The "shared script replacement" merge problem: when a community pack also overrides the same file
- TODO: document the actual workflow Reapy uses — three-way diff? manual diff? specific tool?

### Specific files Reapy has overridden in `zm_test`
- TODO: enumerate the override files under `usermaps/zm_test/scripts/zm/` and group by why they exist
- Cross-link to the systems pages where each override is meaningful

### `zm_patch.csv` — the reference removal step
- What `zm_patch.csv` actually is (verify)
- When commenting out a reference is required for the override to take effect
- The "edit zm_patch.csv → re-link → still broken → also override the file" failure cycle

## Common gotchas

- Override path doesn't match the share path exactly → linker doesn't shadow → silent failure
- Forgetting to remove the `zm_patch.csv` line → conflict at link time
- Diverging from stock without writing down what changed → impossible to merge upstream changes later
- A community pack also overriding the same file → unpredictable which override wins; you may need to fold both sets of changes into a single override

## Cross-references

- [`01-overview.md`](./01-overview.md) — the GSC compilation model that makes overrides possible
- Asset pipeline: [`07-community-packs.md`](../asset-pipeline/07-community-packs.md) — community packs ship overrides too
