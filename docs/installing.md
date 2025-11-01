# Installing beads

Complete installation guide for all platforms.

## Quick Install (Recommended)

### Homebrew (macOS/Linux)

```sh
brew tap shaneholloman/beads
brew install beads
```

**Why Homebrew?**

- ✔ Simple one-command install
- ✔ Automatic updates via `brew upgrade`
- ✔ No need to install Go
- ✔ Handles PATH setup automatically

### Quick Install Script (All Platforms)

```sh
curl -fsSL https://raw.githubusercontent.com/shaneholloman/beads/main/scripts/install.sh | bash
```

The installer will:

- Detect your platform (macOS/Linux, amd64/arm64)
- Install via `go install` if Go is available
- Fall back to building from source if needed
- Guide you through PATH setup if necessary

## Platform-Specific Installation

### macOS

**Via Homebrew** (recommended):

```sh
brew tap shaneholloman/beads
brew install beads
```

**Via go install**:

```sh
go install github.com/shaneholloman/beads/cmd/beads@latest
```

**From source**:

```sh
git clone https://github.com/shaneholloman/beads
cd beads
go build -o beads ./cmd/beads
sudo mv beads /usr/local/bin/
```

### Linux

**Via Homebrew** (works on Linux too):

```sh
brew tap shaneholloman/beads
brew install beads
```

**Arch Linux** (AUR):

```sh
# Install from AUR
yay -S beads-git
# or
paru -S beads-git
```

Thanks to [@v4rgas](https://github.com/v4rgas) for maintaining the AUR package!

**Via go install**:

```sh
go install github.com/shaneholloman/beads/cmd/beads@latest
```

**From source**:

```sh
git clone https://github.com/shaneholloman/beads
cd beads
go build -o beads ./cmd/beads
sudo mv beads /usr/local/bin/
```

### Windows 11

Beads now ships with native Windows support—no MSYS or MinGW required.

**Prerequisites:**

- [Go 1.24+](https://go.dev/dl/) installed (add `%USERPROFILE%\go\bin` to your `PATH`)
- Git for Windows

**Via PowerShell script**:

```pwsh
irm https://raw.githubusercontent.com/shaneholloman/beads/main/install.ps1 | iex
```

**Via go install**:

```pwsh
go install github.com/shaneholloman/beads/cmd/beads@latest
```

**From source**:

```pwsh
git clone https://github.com/shaneholloman/beads
cd beads
go build -o beads.exe ./cmd/beads
Move-Item beads.exe $env:USERPROFILE\AppData\Local\Microsoft\WindowsApps\
```

**Verify installation**:

```pwsh
beads version
```

**Windows notes:**

- The background daemon listens on a loopback TCP endpoint recorded in `.beads\beads.sock`
- Keep that metadata file intact
- Allow `beads.exe` loopback traffic through any host firewall

## IDE and Editor Integrations

### Claude Code Plugin

For Claude Code users, the beads plugin provides slash commands and MCP tools.

**Prerequisites:**

1. First, install the beads CLI (see above)
2. Then install the plugin:

```sh
# In Claude Code
/plugin marketplace add shaneholloman/beads
/plugin install beads
# Restart Claude Code
```

The plugin includes:

- Slash commands: `/beads-ready`, `/beads-create`, `/beads-show`, `/beads-update`, `/beads-close`, etc.
- Full MCP server with all beads tools
- Task agent for autonomous execution

See [plugin.md](./plugin.md) for complete plugin documentation.

### MCP Server (For Sourcegraph Amp, Claude Desktop, and other MCP clients)

If you're using an MCP-compatible tool other than Claude Code:

```sh
# Using uv (recommended)
uv tool install mcp-beads

# Or using pip
uv tool install mcp-beads
```

**Configuration for Claude Desktop** (macOS):

Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "beads": {
      "command": "mcp-beads"
    }
  }
}
```

**Configuration for Sourcegraph Amp**:

Add to your MCP settings:

```json
{
  "beads": {
    "command": "mcp-beads",
    "args": []
  }
}
```

**What you get:**

- Full beads functionality exposed via MCP protocol
- Tools for creating, updating, listing, and closing issues
- Ready work detection and dependency management
- All without requiring Bash commands

See [adapters/mcp/README.md](../adapters/mcp/README.md) for detailed MCP server documentation.

## Verifying Installation

After installing, verify beads is working:

```sh
beads version
beads help
```

## Troubleshooting Installation

### `beads: command not found`

beads is not in your PATH. Either:

```sh
# Check if installed
go list -f {{.Target}} github.com/shaneholloman/beads/cmd/beads

# Add Go bin to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH="$PATH:$(go env GOPATH)/bin"

# Or reinstall
go install github.com/shaneholloman/beads/cmd/beads@latest
```

### `zsh: killed beads` or crashes on macOS

Some users report crashes when running `beads init` or other commands on macOS. This is typically caused by CGO/SQLite compatibility issues.

**Workaround:**

```sh
# Build with CGO enabled
CGO_ENABLED=1 go install github.com/shaneholloman/beads/cmd/beads@latest

# Or if building from source
git clone https://github.com/shaneholloman/beads
cd beads
CGO_ENABLED=1 go build -o beads ./cmd/beads
sudo mv beads /usr/local/bin/
```

If you installed via Homebrew, this shouldn't be necessary as the formula already enables CGO. If you're still seeing crashes with the Homebrew version, please [file an issue](https://github.com/shaneholloman/beads/issues).

## Next Steps

After installation:

1. **Initialize a project**: `cd your-project && beads init`
2. **Configure your agent**: Add beads instructions to `AGENTS.md` (see [README.md](../README.md#quick-start))
3. **Learn the basics**: Run `beads quickstart` for an interactive tutorial
4. **Explore examples**: Check out the [examples/](../examples) directory

## Updating beads

### Homebrew

```sh
brew upgrade beads
```

### go install

```sh
go install github.com/shaneholloman/beads/cmd/beads@latest
```

### From source

```sh
cd beads
git pull
go build -o beads ./cmd/beads
sudo mv beads /usr/local/bin/
```
