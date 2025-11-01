---
description: Check beads and plugin versions
---

# Version Check

> Check the installed versions of beads components and verify compatibility.

**Note:** The MCP server automatically checks beads CLI version >= 0.9.0 on startup. This command provides detailed version info and update instructions.

Use the beads MCP tools to:

1. Run `beads version` via bash to get the CLI version
2. Check the plugin version (0.9.2)
3. Compare versions and report any mismatches

Display:

- beads CLI version (from `beads version`)
- Plugin version (0.9.2)
- MCP server version (0.9.2)
- MCP server status (from `stats` tool or connection test)
- Compatibility status (✔ compatible or WARNING: update needed)

If versions are mismatched, provide instructions:

- Update beads CLI: `curl -fsSL https://raw.githubusercontent.com/shaneholloman/beads/main/install.sh | bash`
- Update plugin: `/plugin update beads`
- Restart Claude Code after updating

Suggest checking for updates if the user is on an older version.
