# Release Process

Quick guide for releasing a new version of beads.

## The Easy Way (Recommended)

Use the fully automated release script:

```sh
./scripts/release.sh 0.9.3
```

This does **everything**: version bump, tests, git tag, Homebrew update, and local installation.

See [scripts/README.md](scripts/README.md#releasesh--the-easy-button) for details.

---

## The Manual Way

If you prefer step-by-step control:

### Pre-Release Checklist

1. **Kill all running daemons (CRITICAL)**:

   ```sh
   # Kill by process name
   pkill -f "beads.*daemon"
   
   # Verify no daemons are running
   pgrep -lf "beads.*daemon" || echo "No daemons running ✔"
   
   # Alternative: find and kill by socket
   find ~/.config -name "beads.sock" -type f 2>/dev/null | while read sock; do
     echo "Found daemon socket: $sock"
   done
   ```

   **Why this matters**: Old daemon versions can cause:
   - Auto-flush race conditions leaving working tree dirty after commits
   - Version mismatches between client (new) and daemon (old)
   - Confusing behavior where changes appear to sync incorrectly

2. **Run tests and build**:

   ```sh
   TMPDIR=/tmp go test ./...
   golangci-lint run ./...
   TMPDIR=/tmp go build -o beads ./cmd/beads
   ./beads version  # Verify it shows new version
   ```

3. **Skip local install** (avoid go install vs brew conflicts):
   - Use `./beads` directly from the repo for testing
   - Your system beads will be updated via brew after Homebrew formula update
   - Or temporarily: `alias beads="$PWD/beads"` if needed

4. **Update CHANGELOG.md**:
   - Add version heading: `## [0.9.X] - YYYY-MM-DD`
   - Summarize changes under: Added, Fixed, Changed, Performance, Community
   - Update Version History section
   - Add Upgrade Guide section if needed

5. **Commit changelog**:

   ```sh
   git add CHANGELOG.md
   git commit -m "Add 0.9.X release notes"
   ```

## Version Bump

Use the automated script to update all version files:

```sh
./scripts/bump-version.sh 0.9.X --commit
git push origin main
```

This updates:

- `cmd/beads/version.go`
- `.claude-plugin/plugin.json`
- `.claude-plugin/marketplace.json`
- `adapters/mcp/pyproject.toml`
- `adapters/mcp/src/mcp_beads/__init__.py`
- `README.md`
- `plugin.md`

**IMPORTANT**: After version bump, rebuild the local binary:

```sh
go build -o beads ./cmd/beads
./beads version  # Should show new version
```

## Publish to All Channels

### 1. Create Git Tag

```sh
git tag v0.9.X
git push origin main
git push origin v0.9.X
```

**That's it!** GitHub Actions automatically handles the rest:

- GoReleaser builds and publishes binaries to GitHub Releases
- PyPI publish job uploads the MCP server package to PyPI

### 2. GitHub Secrets Setup (One-Time)

The automation requires this secret to be configured:

**PYPI_API_TOKEN**: Your PyPI API token

1. Generate token at <https://pypi.org/manage/account/token/>
2. Add to GitHub at <https://github.com/shaneholloman/beads/settings/secrets/actions>
3. Name: `PYPI_API_TOKEN`
4. Value: `pypi-...` (your full token)

### 3. Manual PyPI Publish (If Needed)

If the automated publish fails, you can manually upload:

```sh
cd adapters/mcp

# Clean and rebuild
rm -rf dist/ build/ src/*.egg-info
uv build

# Upload to PyPI
TWINE_USERNAME=__token__ TWINE_PASSWORD=pypi-... uv tool run twine upload dist/*
```

See [adapters/mcp/pypi.md](adapters/mcp/pypi.md) for detailed PyPI instructions.

### 3. Update Homebrew Formula

**CRITICAL**: This step must be done AFTER pushing the git tag in step 1, otherwise the tarball won't exist yet.

The formula needs the SHA256 of the tag tarball:

```sh
# Compute SHA256 from tag (wait a few seconds after pushing tag if you get 404)
curl -sL https://github.com/shaneholloman/beads/archive/refs/tags/v0.9.X.tar.gz | shasum -a 256

# Navigate to tap repo (if already cloned) or clone it
cd /tmp/homebrew-beads || git clone https://github.com/shaneholloman/homebrew-beads /tmp/homebrew-beads

# Pull latest changes
cd /tmp/homebrew-beads
git pull

# Update Formula/beads.rb (replace version and SHA256):
# - url: https://github.com/shaneholloman/beads/archive/refs/tags/v0.9.X.tar.gz
# - sha256: <computed SHA256>

# Commit and push
git add Formula/beads.rb
git commit -m "Update beads formula to v0.9.X"
git push origin main

# CRITICAL: Verify the installation works
brew update
brew upgrade beads  # Or: brew reinstall beads
beads version  # MUST show v0.9.X - if not, the release is incomplete!
```

**WARNING: DO NOT SKIP THE VERIFICATION STEP ABOVE** - This ensures users can actually install the new version.

**Note**: Until this step is complete, users with Homebrew-installed beads will still have the old version.

**Note:** If you have an old beads binary from `go install` in your PATH, remove it to avoid conflicts:

```sh
# Find where beads is installed
which beads

# If it's in a Go toolchain path (e.g., ~/go/bin/beads or mise-managed Go), remove it
# Common locations:
rm ~/go/bin/beads                                    # Standard go install location
rm ~/.local/share/mise/installs/go/*/bin/beads      # mise-managed Go installs

# Verify you're using the correct version
which beads        # Should show /opt/homebrew/bin/beads or your package manager's path
beads version      # Should show the latest version
```

### 4. GitHub Releases (Automated)

**GoReleaser automatically creates releases when you push tags!**

The `.github/workflows/release.yml` workflow:

- Triggers on `v*` tags
- Builds cross-platform binaries (Linux, macOS, Windows for amd64/arm64)
- Generates checksums
- Creates GitHub release with binaries and changelog
- Publishes release automatically

Just push your tag and wait ~5 minutes:

```sh
git push origin v0.9.X
```

Monitor at: <https://github.com/shaneholloman/beads/actions>

The release will appear at: <https://github.com/shaneholloman/beads/releases>

## Post-Release

1. **Kill old daemons again**:

   ```sh
   pkill -f "beads.*daemon"
   pgrep -lf "beads.*daemon" || echo "No daemons running ✔"
   ```

   This ensures your local machine picks up the new version immediately.

2. **Verify installations**:

   ```sh
   # Homebrew
   brew update && brew upgrade beads && beads version
   
   # PyPI
   uv tool install --upgrade mcp-beads
   mcp-beads --help
   
   # Check daemon version matches client
   beads version --daemon  # Should match client version after first command
   ```

3. **Announce** (optional):
   - Project Discord/Slack
   - Twitter/social media
   - README badges

## Troubleshooting

### Stale dist/ directory

Always `rm -rf dist/` before `uv build` to avoid uploading old versions.

### PyPI version conflict

PyPI doesn't allow re-uploading same version. Increment version number even for fixes.

### Homebrew SHA256 mismatch

Wait a few seconds after pushing tag for GitHub to make tarball available, then recompute SHA256.

### Missing PyPI credentials

Set up API token at <https://pypi.org/manage/account/token/> and use `__token__` as username.

## Automation Status

✔ **Automated:**

- GitHub releases with binaries (GoReleaser + GitHub Actions)
- PyPI publish (automated via GitHub Actions)
- Cross-platform builds (Linux, macOS, Windows)
- Checksums and changelog generation

**TODO:**

- Auto-update Homebrew formula

## Related Documentation

- [CHANGELOG.md](CHANGELOG.md) - Release history
- [scripts/README.md](scripts/README.md) - Version bump script details
- [adapters/mcp/pypi.md](adapters/mcp/pypi.md) - Detailed PyPI guide
