---
name: setup
description: Install or update the t7kb Black Ops 3 knowledge base and register its MCP server. Run this once after installing the plugin (invoked as /t7kb:setup).
---

# Install t7kb

One-time setup: download the t7kb tool + database, then register it as an MCP server so it auto-starts in future sessions. Do these steps for the user's OS, reporting what you run.

## 1. Run the installer

Detect the OS and run the matching command. It downloads the binary, the embedding model, and the ~0.9 GB database archive into one folder, and prints the install path.

- **Linux / macOS / WSL:**
  ```bash
  curl -fsSL https://raw.githubusercontent.com/t7-reapy/t7_companion/main/install/install.sh | bash
  ```
  Installs to `~/.t7kb`; binary at `~/.t7kb/t7kb`.

- **Windows (PowerShell):**
  ```powershell
  irm https://raw.githubusercontent.com/t7-reapy/t7_companion/main/install/install.ps1 | iex
  ```
  Installs to `%LOCALAPPDATA%\t7kb`; binary at `%LOCALAPPDATA%\t7kb\t7kb.exe`.

## 2. Register the MCP server

Use the **absolute path** to the installed binary from step 1 (include `.exe` on Windows — a path without the extension fails to spawn on Windows):

```bash
claude mcp add t7kb -- /absolute/path/to/t7kb mcp
```

Example paths: `~/.t7kb/t7kb` (Linux/macOS, expand `~` to the real home) or `C:\Users\<you>\AppData\Local\t7kb\t7kb.exe` (Windows).

## 3. Confirm

Tell the user setup is done and that the `t7kb` MCP tools (`search`, `get`) become available in the **next session** (or after `/reload-plugins`). The 3.5 GB database unpacks itself automatically the first time the server runs.

If the install or registration fails, report the exact error — don't claim success.
