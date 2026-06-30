# Universal installer for t7kb on Windows — the BO3-modding knowledge search tool.
# Downloads the latest release (binary + embedding model) and the database into
# one directory; any MCP-capable agent (Claude Code, Codex, OpenCode, Copilot,
# Cursor) then points at "<dir>\t7kb.exe mcp". Re-runnable.
#
#   irm https://raw.githubusercontent.com/t7-reapy/t7_companion/main/install/install.ps1 | iex
#   .\install.ps1 [-Target <dir>]    # default: %LOCALAPPDATA%\t7kb
param([string]$Target = "$env:LOCALAPPDATA\t7kb")

$ErrorActionPreference = "Stop"
$repo = "t7-reapy/t7_companion"
$base = "https://github.com/$repo/releases/latest/download"

New-Item -ItemType Directory -Force -Path $Target | Out-Null
$tmp = Join-Path $env:TEMP "t7kb_install.zip"

Write-Host "Installing t7kb to $Target"
Write-Host "  downloading t7kb_windows_amd64.zip ..."
Invoke-WebRequest -Uri "$base/t7kb_windows_amd64.zip" -OutFile $tmp
Expand-Archive -Path $tmp -DestinationPath $Target -Force
Remove-Item $tmp -Force

Write-Host "  downloading t7kb.db.zip (~0.9 GB) ..."
# Left zipped on purpose: the binary unpacks it once on first run.
Invoke-WebRequest -Uri "$base/t7kb.db.zip" -OutFile (Join-Path $Target "t7kb.db.zip")

$exe = Join-Path $Target "t7kb.exe"
Write-Host @"

t7kb installed.  Binary: $exe

Register it as an MCP server in your agent (command + args):
    command: $exe
    args:    ["mcp"]

The database unpacks automatically the first time the server runs.
See the README for per-client config (Claude Code, Codex, OpenCode, Copilot, Cursor).
"@
