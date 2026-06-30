---
name: bo3-knowledge
description: How to answer Black Ops 3 / BO3 / Treyarch mod-tools / custom-zombies modding questions using the t7kb knowledge base — GSC/CSC scripting, Radiant mapping, zombies mechanics, assets, FX, audio, lighting, and compile/linker errors.
---

# Answering BO3 modding questions with t7kb

You have a local knowledge base of the Black Ops 3 modding community via the **t7kb** MCP server — tools `search` (hybrid keyword + semantic) and `get` (full document by `doc_id`). For any non-trivial BO3 modding question, query it before answering from memory; the corpus is the authority on what BO3 modding actually contains, your training data is not.

_If the `t7kb` tools aren't available, the knowledge base isn't installed — run `/t7kb:setup` first._

## Query it well

- **Search broad, then narrow.** Issue several short, differently-phrased `search` queries (symptom-side, mechanism-side, exact-jargon-side). The full-text index is conjunctive, so a single phrasing misses the long tail.
- **Read full bodies.** `get` the top `doc_id`s — don't answer from snippets.
- **Weigh reliability.** Each result carries a `reliability` score. On conflict, prefer higher-reliability sources and surface the disagreement when it matters.
- **Cite.** When a claim comes from the kb, name the `source` + `url`.

## Verify shipped tokens against ground truth

The corpus is a starting point, not the final authority. For anything Treyarch **shipped** — exact function names, entity KVPs, asset fields, error strings, file paths — confirm against the raw mod-tools install (the game's own files under the BO3 root) before stating it as fact. Decompiled and community sources can be paraphrased or subtly wrong; the shipped files are ground truth. Drop any claim you can't ground.

## Don't invent

BO3 has its own vocabulary. Cross-game intuitions (other CoD titles, generic engine terms) are usually wrong here. If neither t7kb nor the raw install supports a function, KVP, or concept, don't assert it exists.
