# Contributing

How this doc tree gets written, and how external community resources flow in without polluting it.

## Goals

- **Self-dependent.** A reader gets the answer here, not by chasing five external links. External links are allowed, but the page should still stand on its own if the link 404s. The long-term aim is for this repository to be the source of truth newcomers reach for first — which only works if the content lives here.
- **Truth-grounded.** Every claim traces back to source code, official files, or hands-on verification. Community write-ups are allowed when the writer hasn't (yet) verified the claim themselves, but they're flagged as such.
- **Living.** Drip-fed by topic, plannotated, never declared "final".

## Trust ranking — sources of truth

When two sources disagree, the higher one wins:

1. **Treyarch source** in `share/raw/` (GSC/CSC), or decompiled from the community. Plus observable and retro-engineered engine behaviour.
2. **Community tooling output** — Greyhound, HydraX, acts, NevisX, Husky. Authoritative for what they expose, not always for *why*.
3. **GSCode JSON** ([Blakintosh/gscode](https://github.com/Blakintosh/gscode), `t7_api_*.json`) for engine builtin signatures.
4. **Hands-on experience** from building real maps in this engine.
5. **Community write-ups** — YouTube transcripts, Discord, forum threads, wiki pages. Not authority by default, but if a claim is *proven* by a community source it can stand here until disproven.

If a claim only has #5 backing it, it must either be verified against #1–#4 or marked `> ⚠️ uncertain` with a hint on how to verify.

## How external corpora are used

Resources are scraped (or planned to be) into [github.com/McReaper/t7_knowledge](https://github.com/McReaper/t7_knowledge/):

- **YouTube transcripts** — 300+ videos, theme-indexed in `bo3_mapping_videos.json`.
- **Discord exports** — *planned*.
- **Forum + wiki dumps** — *planned*.

Their role is **not** to be quoted, cited, or linked from doc pages. Their role is:

- **Surface gaps.** What topics did we forget to plan for? (The YouTube tagging already surfaced 9 missing categories.)
- **Challenge content.** Once a page is drafted, scan the matching corpus slice — does community knowledge contradict or complicate what we wrote? Treat hits as prompts to verify against #1–#4, not as edits to apply directly.
- **Suggest depth.** Frequency overall — many videos *or* heavy Discord chatter *or* recurring forum questions on a topic + a 2-paragraph page = signal that the page is too thin.
- **Find recipes.** Spot the common community patterns, then verify the mechanism in source before documenting.

External links in the doc are allowed when they point at:

- A canonical Treyarch source file path
- An official Activision / Treyarch reference
- A tool the reader genuinely needs to install (its GitHub repo, not a tutorial about it)

External links are **not** used to substitute for explanation we should be writing here.

## The per-page loop

1. Pick a page or section. Order of priority: asset-pipeline → scripting → mapping → lighting → audio → systems → lua-lui → community.
2. Claude asks 1–5 short focused questions to clarify scope and shadow spots.
3. Reapy answers short.
4. Claude drip-writes the `.md` file(s) — by topic, not whole-page-from-one-answer.
5. Reapy runs the **[plannotator plugin](https://github.com/backnotprop/plannotator)** (`/plannotator-annotate <path>`) for structured per-line feedback.
6. Multi-pass until clean.
7. **Corpus pass** — Claude pulls the matching corpus slice and surfaces gaps/contradictions in chat (never directly into the doc). The exact CLI flow for this still needs to be settled — see the planned `/corpus-challenge` skill below.
8. Apply edits if warranted, re-plannotate if non-trivial. Move on.

The corpus pass (step 7) will be formalized as a dedicated skill (`/corpus-challenge <path>`) once the Discord and forum/wiki dumps exist and the slicing logic is non-trivial. Until then, ad-hoc.

## Commit policy

- **Never auto-commit.** Reapy reviews every diff and triggers commits explicitly.
- Conventional commits. Parenthetical scopes optional — `docs:` alone is fine, `docs(scripting): ...` when useful.
- Stub files commit fine — they're scaffolding, not claims.

## Handling contradictions

When the corpus contradicts a doc page and neither side can be verified against #1–#4, **keep the existing claim and mark it loudly with `> ⚠️ uncertain`** plus a verification hint. Pulling the claim entirely loses information; the warning is enough to keep the reader honest.

## Cross-references

- Page template, conventions, and section order: [`docs/README.md`](./README.md).
- Routing for fresh Claude sessions, including the Q&A loop and external resource paths: `CLAUDE.md` at repo root.
