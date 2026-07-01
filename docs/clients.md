# Connecting an MCP client to t7kb

After installing (see the [README](../README.md)), every client runs the same thing: **`t7kb mcp`** as a stdio MCP server, pointed at your installed binary. Replace `/path/to/t7kb` below (Windows: `C:\path\to\t7kb.exe`).

**Claude Code** — the plugin does this for you (`/t7kb:setup`); to wire it manually:
```bash
claude mcp add t7kb -- /path/to/t7kb mcp
```

**Codex** — `~/.codex/config.toml`:
```toml
[mcp_servers.t7kb]
command = "/path/to/t7kb"
args = ["mcp"]
```

**OpenCode** — `opencode.json`:
```json
{ "mcp": { "t7kb": { "type": "local", "command": ["/path/to/t7kb", "mcp"], "enabled": true } } }
```

**Copilot (VS Code)** — `.vscode/mcp.json`:
```json
{ "servers": { "t7kb": { "type": "stdio", "command": "/path/to/t7kb", "args": ["mcp"] } } }
```

**Cursor** — `.cursor/mcp.json`:
```json
{ "mcpServers": { "t7kb": { "command": "/path/to/t7kb", "args": ["mcp"] } } }
```

## Workspace guidance

Drop [`templates/AGENTS.md`](../templates/AGENTS.md) at your **BO3 mod-tools root** (the folder with `raw/`, `share_raw/`, `usermaps/`, `mods/` — not a specific map/mod folder) so the agent knows to use the knowledge base and how to query it well. AGENTS.md-aware tools (Codex, OpenCode, recent Cursor) load ancestor files, so one copy at the root covers every map/mod underneath; add a second, project-specific `AGENTS.md` inside a `usermaps/<map>` or `mods/<mod>` folder for that project's own conventions — it layers on top rather than replacing the root one.

Claude Code reads `CLAUDE.md`, not `AGENTS.md` — but the plugin's skills already cover this guidance without any file needed. If you also want it version-controlled or shared with a team that uses other tools, drop a one-line `CLAUDE.md` next to `AGENTS.md` at the same root:

```
@AGENTS.md
```

Claude Code walks the directory tree the same way, loading any `CLAUDE.md` from that root down to wherever a session starts. `/t7kb:setup` offers to set both files up automatically the first time it runs in a BO3 root that doesn't have them yet.

Editors with their own rules file want the `AGENTS.md` contents pasted in instead:

| Editor | Instruction file |
|---|---|
| GitHub Copilot | `.github/copilot-instructions.md` |
| Cursor (older) | `.cursor/rules/t7kb.mdc` |
| Windsurf | `.windsurf/rules/t7kb.md` |
| Cline | `.clinerules/t7kb.md` |
| Kiro | `.kiro/steering/t7kb.md` |
