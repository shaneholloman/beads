# Troubleshooting beads

Common issues and solutions for beads users.

## Table of Contents

- [Installation Issues](#installation-issues)
- [Database Issues](#database-issues)
- [Git and Sync Issues](#git-and-sync-issues)
- [Ready Work and Dependencies](#ready-work-and-dependencies)
- [Performance Issues](#performance-issues)
- [Agent-Specific Issues](#agent-specific-issues)
- [Platform-Specific Issues](#platform-specific-issues)

## Installation Issues

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

## Database Issues

### `database is locked`

Another beads process is accessing the database, or SQLite didn't close properly. Solutions:

```sh
# Find and kill hanging processes
ps aux | grep beads
kill <pid>

# Remove lock files (safe if no beads processes running)
rm .beads/*.db-journal .beads/*.db-wal .beads/*.db-shm
```

**Note**: beads uses a pure Go SQLite driver (`modernc.org/sqlite`) for better portability. Under extreme concurrent load (100+ simultaneous operations), you may see "database is locked" errors. This is a known limitation of the pure Go implementation and does not affect normal usage. For very high concurrency scenarios, consider using the CGO-enabled driver or PostgreSQL (planned for future release).

### `beads init` fails with "directory not empty"

`.beads/` already exists. Options:

```sh
# Use existing database
beads list  # Should work if already initialized

# Or remove and reinitialize (DESTROYS DATA!)
rm -rf .beads/
beads init
```

### `failed to import: issue already exists`

You're trying to import issues that conflict with existing ones. Options:

```sh
# Skip existing issues (only import new ones)
beads import -i issues.jsonl --skip-existing

# Or clear database and re-import everything
rm .beads/*.db
beads import -i .beads/issues.jsonl
```

### Database corruption

**Important**: Distinguish between **logical consistency issues** (ID collisions, wrong prefixes) and **physical SQLite corruption**.

For **physical database corruption** (disk failures, power loss, filesystem errors):

```sh
# Check database integrity
sqlite3 .beads/*.db "PRAGMA integrity_check;"

# If corrupted, reimport from JSONL (source of truth in git)
mv .beads/*.db .beads/*.db.backup
beads init
beads import -i .beads/issues.jsonl
```

For **logical consistency issues** (ID collisions from branch merges, parallel workers):

```sh
# This is NOT corruption - use collision resolution instead
beads import -i .beads/issues.jsonl --resolve-collisions
```

See [FAQ](faq.md#whats-the-difference-between-sqlite-corruption-and-id-collisions) for the distinction.

### Multiple databases detected warning

If you see a warning about multiple `.beads` databases in the directory hierarchy:

```txt
╔══════════════════════════════════════════════════════════════════════════╗
║ WARNING: 2 beads databases detected in directory hierarchy               ║
╠══════════════════════════════════════════════════════════════════════════╣
║ Multiple databases can cause confusion and database pollution.           ║
║                                                                          ║
║ ▶ /path/to/project/.beads (15 issues)                                    ║
║   /path/to/parent/.beads (32 issues)                                     ║
║                                                                          ║
║ Currently using the closest database (▶). This is usually correct.       ║
║                                                                          ║
║ RECOMMENDED: Consolidate or remove unused databases to avoid confusion.  ║
╚══════════════════════════════════════════════════════════════════════════╝
```

This means beads found multiple `.beads` directories in your directory hierarchy. The `▶` marker shows which database is actively being used (usually the closest one to your current directory).

**Why this matters:**

- Can cause confusion about which database contains your work
- Easy to accidentally work in the wrong database
- May lead to duplicate tracking of the same work

**Solutions:**

1. **If you have nested projects** (intentional):
    - This is fine! beads is designed to support this
    - Just be aware which database you're using
    - Set `BEADS_DB` environment variable if you want to override the default selection

2. **If you have accidental duplicates** (unintentional):
    - Decide which database to keep
    - Export issues from the unwanted database: `cd <unwanted-dir> && beads export -o backup.jsonl`
    - Remove the unwanted `.beads` directory: `rm -rf <unwanted-dir>/.beads`
    - Optionally import issues into the main database if needed

3. **Override database selection**:

    ```sh
    # Temporarily use specific database
    BEADS_DB=/path/to/.beads/issues.db beads list

    # Or add to shell config for permanent override
    export BEADS_DB=/path/to/.beads/issues.db
    ```

**Note**: The warning only appears when beads detects multiple databases. If you see this consistently and want to suppress it, you're using the correct database (marked with `▶`).

## Git and Sync Issues

### Git merge conflict in `issues.jsonl`

When both sides add issues, you'll get conflicts. Resolution:

1. Open `.beads/issues.jsonl`
2. Look for `<<<<<<< HEAD` markers
3. Most conflicts can be resolved by **keeping both sides**
4. Each line is independent unless IDs conflict
5. For same-ID conflicts, keep the newest (check `updated_at`)

Example resolution:

```sh
# After resolving conflicts manually
git add .beads/issues.jsonl
git commit
beads import -i .beads/issues.jsonl  # Sync to SQLite
```

See [advanced.md](advanced.md) for detailed merge strategies.

### Git merge conflicts in JSONL

**With hash-based IDs (v0.20.1+), ID collisions don't occur.** Different issues get different hash IDs.

If git shows a conflict in `.beads/issues.jsonl`, it's because the same issue was modified on both branches:

```sh
# Preview what will be updated
beads import -i .beads/issues.jsonl --dry-run

# Resolve git conflict (keep newer version or manually merge)
git checkout --theirs .beads/issues.jsonl  # Or --ours, or edit manually

# Import updates the database
beads import -i .beads/issues.jsonl
```

See [advanced.md#handling-git-merge-conflicts](advanced.md#handling-git-merge-conflicts) for details.

### Permission denied on git hooks

Git hooks need execute permissions:

```sh
chmod +x .git/hooks/pre-commit
chmod +x .git/hooks/post-merge
chmod +x .git/hooks/post-checkout
```

Or use the installer: `cd examples/git-hooks && ./install.sh`

### Auto-sync not working

Check if auto-sync is enabled:

```sh
# Check if daemon is running
ps aux | grep "beads daemon"

# Manually export/import
beads export -o .beads/issues.jsonl
beads import -i .beads/issues.jsonl

# Install git hooks for guaranteed sync
cd examples/git-hooks && ./install.sh
```

If you disabled auto-sync with `--no-auto-flush` or `--no-auto-import`, remove those flags or use `beads sync` manually.

## Ready Work and Dependencies

### `beads ready` shows nothing but I have open issues

Those issues probably have open blockers. Check:

```sh
# See blocked issues
beads blocked

# Show dependency tree (default max depth: 50)
beads dep tree <issue-id>

# Limit tree depth to prevent deep traversals
beads dep tree <issue-id> --max-depth 10

# Remove blocking dependency if needed
beads dep remove <from-id> <to-id>
```

Remember: Only `blocks` dependencies affect ready work.

### Circular dependency errors

beads prevents dependency cycles, which break ready work detection. To fix:

```sh
# Detect all cycles
beads dep cycles

# Remove the dependency causing the cycle
beads dep remove <from-id> <to-id>

# Or redesign your dependency structure
```

### Dependencies not showing up

Check the dependency type:

```sh
# Show full issue details including dependencies
beads show <issue-id>

# Visualize the dependency tree
beads dep tree <issue-id>
```

Remember: Different dependency types have different meanings:

- `blocks` - Hard blocker, affects ready work
- `related` - Soft relationship, doesn't block
- `parent-child` - Hierarchical (child depends on parent)
- `discovered-from` - Work discovered during another issue

## Performance Issues

### Export/import is slow

For large databases (10k+ issues):

```sh
# Export only open issues
beads export --format=jsonl --status=open -o .beads/issues.jsonl

# Or filter by priority
beads export --format=jsonl --priority=0 --priority=1 -o critical.jsonl
```

Consider splitting large projects into multiple databases.

### Commands are slow

Check database size and consider compaction:

```sh
# Check database stats
beads stats

# Preview compaction candidates
beads compact --dry-run --all

# Compact old closed issues
beads compact --days 90
```

### Large JSONL files

If `.beads/issues.jsonl` is very large:

```sh
# Check file size
ls -lh .beads/issues.jsonl

# Remove old closed issues
beads compact --days 90

# Or split into multiple projects
cd ~/project/component1 && beads init --prefix comp1
cd ~/project/component2 && beads init --prefix comp2
```

## Agent-Specific Issues

### Agent creates duplicate issues

Agents may not realize an issue already exists. Prevention strategies:

- Have agents search first: `beads list --json | grep "title"`
- Use labels to mark auto-created issues: `beads create "..." -l auto-generated`
- Review and deduplicate periodically: `beads list | sort`
- Use `beads merge` to consolidate duplicates: `beads merge beads-2 --into beads-1`

### Agent gets confused by complex dependencies

Simplify the dependency structure:

```sh
# Check for overly complex trees
beads dep tree <issue-id>

# Remove unnecessary dependencies
beads dep remove <from-id> <to-id>

# Use labels instead of dependencies for loose relationships
beads label add <issue-id> related-to-feature-X
```

### Agent can't find ready work

Check if issues are blocked:

```sh
# See what's blocked
beads blocked

# See what's actually ready
beads ready --json

# Check specific issue
beads show <issue-id>
beads dep tree <issue-id>
```

### MCP server not working

Check installation and configuration:

```sh
# Verify MCP server is installed
pip list | grep mcp-beads

# Check MCP configuration
cat ~/Library/Application\ Support/Claude/claude_desktop_config.json

# Test CLI works
beads version
beads ready

# Check for daemon
ps aux | grep "beads daemon"
```

See [adapters/mcp/README.md](../adapters/mcp/README.md) for MCP-specific troubleshooting.

### Claude Code sandbox mode

**Issue:** Claude Code's sandbox restricts network access to a single socket, conflicting with beads's daemon and git operations.

**Solution:** Use the `--sandbox` flag:

```sh
# Sandbox mode disables daemon and auto-sync
beads --sandbox ready
beads --sandbox create "Fix bug" -p 1
beads --sandbox update beads-42 --status in_progress

# Or set individual flags
beads --no-daemon --no-auto-flush --no-auto-import <command>
```

**What sandbox mode does:**

- Disables daemon (uses direct SQLite mode)
- Disables auto-export to JSONL
- Disables auto-import from JSONL
- Allows beads to work in network-restricted environments

**Note:** You'll need to manually sync when outside the sandbox:

```sh
# After leaving sandbox, sync manually
beads sync
```

**Related:** See [Claude Code sandboxing documentation](https://www.anthropic.com/engineering/claude-code-sandboxing) for more about sandbox restrictions.

## Platform-Specific Issues

### Windows: Path issues

```pwsh
# Check if beads.exe is in PATH
where.exe beads

# Add Go bin to PATH (permanently)
[Environment]::SetEnvironmentVariable(
    "Path",
    $env:Path + ";$env:USERPROFILE\go\bin",
    [EnvironmentVariableTarget]::User
)

# Reload PATH in current session
$env:Path = [Environment]::GetEnvironmentVariable("Path", "User")
```

### Windows: Firewall blocking daemon

The daemon listens on loopback TCP. Allow `beads.exe` through Windows Firewall:

1. Open Windows Security → Firewall & network protection
2. Click "Allow an app through firewall"
3. Add `beads.exe` and enable for Private networks
4. Or disable firewall temporarily for testing

### macOS: Gatekeeper blocking execution

If macOS blocks beads:

```sh
# Remove quarantine attribute
xattr -d com.apple.quarantine /usr/local/bin/beads

# Or allow in System Preferences
# System Preferences → Security & Privacy → General → "Allow anyway"
```

### Linux: Permission denied

If you get permission errors:

```sh
# Make beads executable
chmod +x /usr/local/bin/beads

# Or install to user directory
mkdir -p ~/.local/bin
mv beads ~/.local/bin/
export PATH="$HOME/.local/bin:$PATH"
```

## Getting Help

If none of these solutions work:

1. **Check existing issues**: [GitHub Issues](https://github.com/shaneholloman/beads/issues)
2. **Enable debug logging**: `beads --verbose <command>`
3. **File a bug report**: Include:
   - beads version: `beads version`
   - OS and architecture: `uname -a`
   - Error message and full command
   - Steps to reproduce
4. **Join discussions**: [GitHub Discussions](https://github.com/shaneholloman/beads/discussions)

## Related Documentation

- **[README.md](../README.md)** - Core features and quick start
- **[advanced.md](advanced.md)** - Advanced features
- **[faq.md](faq.md)** - Frequently asked questions
- **[installing.md](installing.md)** - Installation guide
- **[advanced.md](advanced.md)** - JSONL format and merge strategies
