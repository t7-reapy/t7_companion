# Adding or removing knowledge in `t7kb.db`

`t7kb.db` is the search index shipped with `t7_companion`. It is built from the upstream BO3 community corpus and published as a [GitHub Release asset](https://github.com/t7-reapy/t7_companion/releases). Every row carries enough metadata (source, license, upstream URL, author) to be added, modified, or removed surgically — without disturbing unrelated content. This page documents the workflow for three common operations.

## 1. Removing content (DMCA / opt-out / takedown)

If you are a rights-holder of any bundled content and want it removed or its disposition changed, contact the maintainer [McReaper](https://github.com/McReaper) directly via a **private channel** (DM, email — not a public issue, to avoid drawing attention to your takedown).

The maintainer applies the change in the upstream pipeline; the next `t7kb.db` release reflects it. There is no public removal list — out of respect for rights-holders, takedown traffic stays private.

You have **two options** depending on how aggressive the removal needs to be:

### Option A: Flip to pointer-only

The row stays in `t7kb.db` and remains searchable by title/summary, but the `body` column is emptied. The result is "here's where to find this — go read it upstream"; no verbatim copyrighted text remains.

This is the lighter option and the default for most takedown requests, because:

- Search behavior stays stable (the row's `doc_id` doesn't change).
- Consumers still get a discoverable pointer to your upstream content via `url` / `content_path`.
- The change is reversible (the upstream pipeline can re-fill the body if posture changes later).

### Option B: Full removal

The row is deleted entirely from `t7kb.db`. Search no longer returns it; its `doc_id` is gone.

Use this when even the *existence* of the entry shouldn't be discoverable in the shipped DB.

### What to include in the takedown request

- The `doc_id`(s) if you know them — or a description of the content / `source` prefix / path subtree you want covered.
- Whether you want **pointer-only** (Option A) or **full removal** (Option B).
- Your preferred attribution if the content stays as a pointer.

> [!NOTE]
> The mechanism is surgical because every row tracks its provenance. `doc_id`s of unaffected rows stay stable across removals — anything else that references those `doc_id`s (a wiki page's `## Sources` block, a saved query result, a cross-link in another row) keeps working.

## 2. Adding a new source

If you maintain a BO3 modding resource (forum, wiki, tool, video channel) and want it ingested into future `t7kb.db` releases:

1. Open an issue on this repo describing the source: URL, license, expected size, scraping feasibility.
2. The maintainer evaluates licensing posture (as-is shippable vs needs synthesis) and feasibility.
3. If accepted, the upstream pipeline adds an ingester for it; subsequent releases include the source.

There is no formal "submit a PR with content" path for `t7kb.db` rows — the DB is built deterministically from the upstream source-set, not edited by hand. **However, wiki content is fully PR'able** (see next section).

## 3. Correcting / improving existing content

| What | How |
|---|---|
| **Wiki page** (`wiki/*.md`) | **Open a PR.** Wiki pages are CC-BY-SA 4.0 and contributions are welcome under the same license. |
| **Skill file** (`skills/*.md`) | **Open a PR.** Same posture as wiki pages. |
| **Docs page** (`docs/*.md`) | **Open a PR.** Same posture. |
| **`t7kb.db` row from an as-is source** (gscode-api, source_scripts, etc.) | **Open an issue.** The fix lives upstream — once the upstream source updates, the next `t7kb.db` build picks it up. |
| **`t7kb.db` row for an attribution / metadata correction only** | **Open an issue.** The maintainer can land a small upstream patch and trigger a rebuild. |

> [!TIP]
> `t7kb.db` rebuilds are cheap. If you spot a wrong author name or a broken upstream URL on any row, just file an issue with the `doc_id` and the correction — it gets picked up in the next release without needing a full source re-scrape.

## 4. Reverting a takedown

If you previously requested pointer-only conversion and later change your mind, contact the maintainer again (same private channel). The same mechanism (one entry in the maintainer's removal config, one rebuild) restores the row's body content.

## Schema reference

For the actual `t7kb.db` schema and what fields each row carries, see [`data-model.md`](data-model.md).

## Contact

[McReaper](https://github.com/McReaper) (maintainer) — issue tracker for additions / corrections, private channel for takedowns.
