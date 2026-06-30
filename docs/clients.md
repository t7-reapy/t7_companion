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

Drop [`templates/AGENTS.md`](../templates/AGENTS.md) at your map project root so the agent knows to use the knowledge base and how to query it well. Claude Code, Codex, OpenCode, and recent Cursor auto-read `AGENTS.md`; editors with their own rules file want its contents pasted in instead:

| Editor | Instruction file |
|---|---|
| GitHub Copilot | `.github/copilot-instructions.md` |
| Cursor (older) | `.cursor/rules/t7kb.mdc` |
| Windsurf | `.windsurf/rules/t7kb.md` |
| Cline | `.clinerules/t7kb.md` |
| Kiro | `.kiro/steering/t7kb.md` |
