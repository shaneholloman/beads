# Beads Scripts

Utility scripts for maintaining the beads project.

## release.sh (The Easy Button)

**One-command release** from version bump to local installation.

### Usage

```sh
# Full release (does everything)
./scripts/release.sh 0.0.0

# Preview what would happen
./scripts/release.sh 0.0.0 --dry-run
```

### What It Does

This master script automates the **entire release process**:

1. ✔ Kills running daemons (avoids version conflicts)
2. ✔ Runs tests and linting
3. ✔ Bumps version in all files
4. ✔ Commits and pushes version bump
5. ✔ Creates and pushes git tag
6. ✔ Updates Homebrew formula
7. ✔ Upgrades local brew installation
8. ✔ Verifies everything works

**After this script completes, your system is running the new version!**

### Examples

```sh
# Release version 0.0.0
./scripts/release.sh 0.0.0

# Preview a release (no changes made)
./scripts/release.sh 1.0.0 --dry-run
```

### Prerequisites

- Clean git working directory
- All changes committed
- golangci-lint installed
- Homebrew installed (for local upgrade)
- Push access to shaneholloman/beads and shaneholloman/homebrew-beads

### Output

The script provides colorful, step-by-step progress output:

- 🟡 Yellow: Current step
- 🟢 Green: Step completed
- 🔴 Red: Errors
- 🔵 Blue: Section headers

### What Happens Next

After the script finishes:

- GitHub Actions builds binaries for all platforms (~5 minutes)
- PyPI package is published automatically
- Users can `brew upgrade beads` to get the new version
- GitHub Release is created with binaries and changelog

---

## bump-version.sh

Bumps the version number across all beads components in a single command.

### Usage

```sh
# Show usage
./scripts/bump-version.sh

# Update versions (shows diff, no commit)
./scripts/bump-version.sh 0.0.0

# Update versions and auto-commit
./scripts/bump-version.sh 0.0.0 --commit
```

### What It Does

Updates version in all these files:

- `cmd/beads/version.go` - beads CLI version constant
- `.claude-plugin/plugin.json` - Plugin version
- `.claude-plugin/marketplace.json` - Marketplace plugin version
- `adapters/mcp/pyproject.toml` - MCP server version
- `README.md` - Alpha status version
- `plugin.md` - Version requirements

### Features

- **Validates** semantic versioning format (MAJOR.MINOR.PATCH)
- **Verifies** all versions match after update
- **Shows** git diff of changes
- **Auto-commits** with standardized message (optional)
- **Cross-platform** compatible (macOS and Linux)

### Examples

```sh
# Bump to 0.0.0 and review changes
./scripts/bump-version.sh 0.0.0
# Review the diff, then manually commit

# Bump to 1.0.0 and auto-commit
./scripts/bump-version.sh 1.0.0 --commit
git push origin main
```

### Why This Script Exists

Previously, version bumps only updated `cmd/beads/version.go`, leaving other components out of sync. This script ensures all version numbers stay consistent across the project.

### Safety

- Checks for uncommitted changes before proceeding
- Refuses to auto-commit if there are existing uncommitted changes
- Validates version format before making any changes
- Verifies all versions match after update
- Shows diff for review before commit

## Future Scripts

Additional maintenance scripts may be added here as needed.
