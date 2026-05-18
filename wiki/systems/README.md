# Gameplay Systems

Game mechanics that are *bigger than scripting* — they thread through assets, scripting, mapping, audio, and FX. Each page collects how the system is wired in BO3 and where the integration points live across pipelines.

These are stubs. Content will land as we work through `zm_test`'s build of each system. Research starting points: the YouTube transcript corpus (`F:/victor/bo3 tools - sources/.Tutorials/youtube_transcripts/bo3_mapping_videos.json`) — filter `theme` for the page topic.

## Pages

- [Perks](./perks.md) — vending machines, perk slots, HarryBo21's pack, custom perks.
- [Pack-a-Punch](./pack-a-punch.md) — vanilla PaP, custom PaP machines, `_up_up` second tier.
- [Mystery Box](./mystery-box.md) — magic box, weapon karusel, custom box weapon lists.
- [Wallbuys](./wallbuys.md) — wall weapon triggers, chalk drawings.
- [Wonder weapons](./wonder-weapons.md) — Ray Gun, Thundergun, Magmagat, custom wonder weapons.
- [Traps](./traps.md) — fire traps, electric traps, custom hazards.
- [Traversals](./traversals.md) — zombie jump-up, jump-down, mantle, climb.
- [Easter eggs](./easter-eggs.md) — soul chests, shootables, song eggs, secret doors, buyable ending.
- [Power-ups](./power-ups.md) — Max Ammo, Insta-Kill, Nuke, custom drops.

## Why a separate section

Each of these is a player-facing system that touches multiple pipelines. Documenting them under "scripting" hides the asset-pipeline and mapping work; documenting them under "asset-pipeline" hides the gameplay logic. They earn their own section.
