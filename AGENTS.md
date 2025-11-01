# Instructions for AI Agents Working on Beads

## Project Overview

This is **beads** (command: `beads`), an issue tracker designed for AI-supervised coding workflows. We dogfood our own tool!

## Human Setup vs Agent Usage

**IMPORTANT:** If you need to initialize beads, use the `--quiet` flag:

```bash
beads init --quiet  # Non-interactive, auto-installs git hooks, no prompts
```

**Why `--quiet`?** Regular `beads init` has interactive prompts (git hooks) that confuse agents. The `--quiet` flag makes it fully non-interactive:

- Automatically installs git hooks
- No prompts for user input
- Safe for agent-driven repo setup

**If the human already initialized:** Just use beads normally with `beads create`, `beads ready`, `beads update`, `beads close`, etc.

**If you see "database not found":** Run `beads init --quiet` yourself, or ask the human to run `beads init`.

## Issue Tracking

**CRITICAL: This repository uses beads (beads) for ALL project work tracking.**

### What Must Be Tracked in beads

ANY work that affects the codebase, documentation, or project state:

- Bug fixes
- New features
- Refactoring
- Documentation changes
- Test additions
- Configuration updates
- Build/CI changes

If it goes into git history, it must be a beads issue first.

### Ephemeral vs Project Work

**Project work (use beads):** Anything that becomes part of the project

- Example: "Rename docs to lowercase-with-dashes" → Create `beads` issue

**Ephemeral session planning (use native agent tools):** Temporary task breakdown for complex work

- Example: Breaking "Implement auth system" into 5 implementation steps
- These are throwaway plans that help YOU work, not project history

**When in doubt, use beads.** The escape hatch exists for genuinely ephemeral planning, not to avoid proper tracking.

### MCP Server (Recommended)

**RECOMMENDED**: Use the MCP (Model Context Protocol) server for the best experience! The beads MCP server provides native integration with Claude and other MCP-compatible AI assistants.

**Installation:**

```bash
# Install the MCP server
uv tool install beads-mcp

# Add to your MCP settings (e.g., Claude Desktop config)
{
  "beads": {
    "command": "beads-mcp",
    "args": []
  }
}
```

**Benefits:**

- Native function calls instead of shell commands
- Automatic workspace detection
- Better error handling and validation
- Structured JSON responses
- No need for `--json` flags

**All beads commands are available as MCP functions** with the prefix `mcp__beads-*__`. For example:

- `beads ready` → `mcp__beads__ready()`
- `beads create` → `mcp__beads__create(title="...", priority=1)`
- `beads update` → `mcp__beads__update(issue_id="beads-42", status="in_progress")`

See `adapters/mcp/README.md` for complete documentation.

### Multi-Repo Configuration (MCP Server)

**RECOMMENDED: Use a single MCP server for all beads projects** - it automatically routes to per-project local daemons.

**Setup (one-time):**

```bash
# MCP config in ~/.config/amp/settings.json or Claude Desktop config:
{
  "beads": {
    "command": "beads-mcp",
    "args": []
  }
}
```

**How it works (LSP model):**
The single MCP server instance automatically:

1. Checks for local daemon socket (`.beads/beads.sock`) in your current workspace
2. Routes requests to the correct **per-project daemon** based on working directory
3. Auto-starts the local daemon if not running (with exponential backoff)
4. **Each project gets its own isolated daemon** serving only its database

**Architecture:**

```
MCP Server (one instance)
    ↓
Per-Project Daemons (one per workspace)
    ↓
SQLite Databases (complete isolation)
```

**Why per-project daemons?**

- ✔ Complete database isolation between projects
- ✔ No cross-project pollution or git worktree conflicts
- ✔ Simpler mental model: one project = one database = one daemon
- ✔ Follows LSP (Language Server Protocol) architecture

**Note:** The daemon **auto-starts automatically** when you run any `beads` command (v0.9.11+). To disable auto-start, set `BEADS_AUTO_START_DAEMON=false`.

**Version Management:** beads automatically handles daemon version mismatches (v0.16.0+):

- When you upgrade beads, old daemons are automatically detected and restarted
- Version compatibility is checked on every connection
- No manual intervention required after upgrades
- Works transparently with MCP server and CLI
- Use `beads daemons health` to check for version mismatches
- Use `beads daemons killall` to force-restart all daemons if needed

**Alternative (not recommended): Multiple MCP Server Instances**
If you must use separate MCP servers:

```json
{
  "beads-webapp": {
    "command": "beads-mcp",
    "env": {
      "BEADS_WORKING_DIR": "/Users/you/projects/webapp"
    }
  },
  "beads-api": {
    "command": "beads-mcp",
    "env": {
      "BEADS_WORKING_DIR": "/Users/you/projects/api"
    }
  }
}
```

**WARNING: Problem**: AI may select the wrong MCP server for your workspace, causing commands to operate on the wrong database.

### CLI Quick Reference

If you're not using the MCP server, here are the CLI commands:

```bash
# Check database path and daemon status
beads info --json

# Find ready work (no blockers)
beads ready --json

# Create new issue
beads create "Issue title" -t bug|feature|task -p 0-4 -d "Description" --json

# Create with explicit ID (for parallel workers)
beads create "Issue title" --id worker1-100 -p 1 --json

# Create with labels
beads create "Issue title" -t bug -p 1 -l bug,critical --json

# Create multiple issues from markdown file
beads create -f feature-plan.md --json

# Create epic with hierarchical child tasks
beads create "Auth System" -t epic -p 1 --json         # Returns: beads-a3f8e9
beads create "Login UI" -p 1 --json                     # Auto-assigned: beads-a3f8e9.1
beads create "Backend validation" -p 1 --json           # Auto-assigned: beads-a3f8e9.2
beads create "Tests" -p 1 --json                        # Auto-assigned: beads-a3f8e9.3

# Update one or more issues
beads update <id> [<id>...] --status in_progress --json
beads update <id> [<id>...] --priority 1 --json

# Edit issue fields in $EDITOR (HUMANS ONLY - not for agents)
beads edit <id>                    # Edit description
beads edit <id> --title            # Edit title
beads edit <id> --design           # Edit design notes
beads edit <id> --notes            # Edit notes
beads edit <id> --acceptance       # Edit acceptance criteria

# Link discovered work (old way)
beads dep add <discovered-id> <parent-id> --type discovered-from

# Create and link in one command (new way)
beads create "Issue title" -t bug -p 1 --deps discovered-from:<parent-id> --json

# Label management (supports multiple IDs)
beads label add <id> [<id>...] <label> --json
beads label remove <id> [<id>...] <label> --json
beads label list <id> --json
beads label list-all --json

# Filter issues by label
beads list --label bug,critical --json

# Complete work (supports multiple IDs)
beads close <id> [<id>...] --reason "Done" --json

# Reopen closed issues (supports multiple IDs)
beads reopen <id> [<id>...] --reason "Reopening" --json

# Show dependency tree
beads dep tree <id>

# Get issue details (supports multiple IDs)
beads show <id> [<id>...] --json

# Rename issue prefix (e.g., from 'knowledge-work-' to 'kw-')
beads rename-prefix kw- --dry-run  # Preview changes
beads rename-prefix kw- --json     # Apply rename

# Restore compacted issue from git history
beads restore <id>  # View full history at time of compaction

# Import issues from JSONL
beads import -i .beads/issues.jsonl --dry-run      # Preview changes
beads import -i .beads/issues.jsonl                # Import and update issues
beads import -i .beads/issues.jsonl --dedupe-after # Import + detect duplicates

# Find and merge duplicate issues
beads duplicates                                          # Show all duplicates
beads duplicates --auto-merge                             # Automatically merge all
beads duplicates --dry-run                                # Preview merge operations

# Merge specific duplicate issues
beads merge <source-id...> --into <target-id> --json      # Consolidate duplicates
beads merge beads-42 beads-43 --into beads-41 --dry-run            # Preview merge

# Migrate databases after version upgrade
beads migrate                                             # Detect and migrate old databases
beads migrate --dry-run                                   # Preview migration
beads migrate --cleanup --yes                             # Migrate and remove old files
```

### Managing Daemons

beads runs a background daemon per workspace for auto-sync and RPC operations. Use `beads daemons` to manage multiple daemons:

```bash
# List all running daemons
beads daemons list --json

# Check health (version mismatches, stale sockets)
beads daemons health --json

# Stop a specific daemon
beads daemons stop /path/to/workspace --json
beads daemons stop 12345 --json  # By PID

# View daemon logs
beads daemons logs /path/to/workspace -n 100
beads daemons logs 12345 -f  # Follow mode

# Stop all daemons
beads daemons killall --json
beads daemons killall --force --json  # Force kill if graceful fails
```

**When to use:**

- **After upgrading beads**: Run `beads daemons health` to check for version mismatches, then `beads daemons killall` to restart all daemons with the new version
- **Debugging**: Use `beads daemons logs <workspace>` to view daemon logs
- **Cleanup**: `beads daemons list` auto-removes stale sockets

**Troubleshooting:**

- **Stale sockets**: `beads daemons list` auto-cleans them
- **Version mismatch**: `beads daemons killall` then let daemons auto-start on next command
- **Daemon won't stop**: `beads daemons killall --force`

See [commands/daemons.md](commands/daemons.md) for detailed documentation.

### Event-Driven Daemon Mode (Experimental)

**NEW in v0.16+**: The daemon supports an experimental event-driven mode that replaces 5-second polling with instant reactivity.

**Benefits:**

- **<500ms latency** (vs ~5000ms with polling)
- **~60% less CPU usage** (no continuous polling)
- **Instant sync** on mutations and file changes
- **Dropped events safety net** prevents data loss

**How it works:**

- **FileWatcher** monitors `.beads/issues.jsonl` and `.git/refs/heads` using platform-native APIs:
  - Linux: `inotify`
  - macOS: `FSEvents` (via kqueue)
  - Windows: `ReadDirectoryChangesW`
- **Mutation events** from RPC operations (create, update, close) trigger immediate export
- **Debouncer** batches rapid changes (500ms window) to avoid export storms
- **Polling fallback** if fsnotify unavailable (e.g., network filesystems)

**Opt-In (Phase 1):**

Event-driven mode is opt-in during Phase 1. To enable:

```bash
# Enable event-driven mode for a single daemon
BEADS_DAEMON_MODE=events beads daemon start

# Or set globally in your shell profile
export BEADS_DAEMON_MODE=events

# Restart all daemons to apply
beads daemons killall
# Next beads command will auto-start daemon with new mode
```

**Available modes:**

- `poll` (default) - Traditional 5-second polling, stable and battle-tested
- `events` - New event-driven mode, experimental but thoroughly tested

**Troubleshooting:**

If the watcher fails to start:

- Check daemon logs: `beads daemons logs /path/to/workspace -n 100`
- Look for "File watcher unavailable" warnings
- Common causes:
  - Network filesystem (NFS, SMB) - fsnotify may not work
  - Container environment - may need privileged mode
  - Resource limits - check `ulimit -n` (open file descriptors)

**Fallback behavior:**

- If `BEADS_DAEMON_MODE=events` but watcher fails, daemon falls back to polling automatically
- Set `BEADS_WATCHER_FALLBACK=false` to disable fallback and require fsnotify

**Disable polling fallback:**

```bash
# Require fsnotify, fail if unavailable
BEADS_WATCHER_FALLBACK=false BEADS_DAEMON_MODE=events beads daemon start
```

**Switch back to polling:**

```bash
# Explicitly use polling mode
BEADS_DAEMON_MODE=poll beads daemon start

# Or unset to use default
unset BEADS_DAEMON_MODE
beads daemons killall  # Restart with default (poll) mode
```

**Future (Phase 2):** Event-driven mode will become the default once it's proven stable in production use.

### Workflow

1. **Check for ready work**: Run `beads ready` to see what's unblocked
2. **Claim your task**: `beads update <id> --status in_progress`
3. **Work on it**: Implement, test, document
4. **Discover new work**: If you find bugs or TODOs, create issues:
   - Old way (two commands): `beads create "Found bug in auth" -t bug -p 1 --json` then `beads dep add <new-id> <current-id> --type discovered-from`
   - New way (one command): `beads create "Found bug in auth" -t bug -p 1 --deps discovered-from:<current-id> --json`
5. **Complete**: `beads close <id> --reason "Implemented"`
6. **Sync at end of session**: `beads sync` (see "Agent Session Workflow" below)

### Issue Types

- `bug` - Something broken that needs fixing
- `feature` - New functionality
- `task` - Work item (tests, docs, refactoring)
- `epic` - Large feature composed of multiple issues (supports hierarchical children)
- `chore` - Maintenance work (dependencies, tooling)

**Hierarchical children:** Epics can have child issues with dotted IDs (e.g., `beads-a3f8e9.1`, `beads-a3f8e9.2`). Children are auto-numbered sequentially. Up to 3 levels of nesting supported. The parent hash ensures unique namespace - no coordination needed between agents working on different epics.

### Priorities

- `0` - Critical (security, data loss, broken builds)
- `1` - High (major features, important bugs)
- `2` - Medium (nice-to-have features, minor bugs)
- `3` - Low (polish, optimization)
- `4` - Backlog (future ideas)

### Dependency Types

- `blocks` - Hard dependency (issue X blocks issue Y)
- `related` - Soft relationship (issues are connected)
- `parent-child` - Epic/subtask relationship
- `discovered-from` - Track issues discovered during work

Only `blocks` dependencies affect the ready work queue.

### Duplicate Detection & Merging

AI agents should proactively detect and merge duplicate issues to keep the database clean:

**Automated duplicate detection:**

```bash
# Find all content duplicates in the database
beads duplicates

# Automatically merge all duplicates
beads duplicates --auto-merge

# Preview what would be merged
beads duplicates --dry-run

# During import
beads import -i issues.jsonl --dedupe-after
```

**Detection strategies:**

1. **Before creating new issues**: Search for similar existing issues

   ```bash
   beads list --json | grep -i "authentication"
   beads show beads-41 beads-42 --json  # Compare candidates
   ```

2. **Periodic duplicate scans**: Review issues by type or priority

   ```bash
   beads list --status open --priority 1 --json  # High-priority issues
   beads list --issue-type bug --json             # All bugs
   ```

3. **During work discovery**: Check for duplicates when filing discovered-from issues

   ```bash
   # Before: beads create "Fix auth bug" --deps discovered-from:beads-100
   # First: beads list --json | grep -i "auth bug"
   # Then decide: create new or link to existing
   ```

**Merge workflow:**

```bash
# Step 1: Identify duplicates (beads-42 and beads-43 duplicate beads-41)
beads show beads-41 beads-42 beads-43 --json

# Step 2: Preview merge to verify
beads merge beads-42 beads-43 --into beads-41 --dry-run

# Step 3: Execute merge
beads merge beads-42 beads-43 --into beads-41 --json

# Step 4: Verify result
beads dep tree beads-41  # Check unified dependency tree
beads show beads-41 --json  # Verify merged content
```

**What gets merged:**

- ✔ All dependencies from source → target
- ✔ Text references updated across ALL issues (descriptions, notes, design, acceptance criteria)
- ✔ Source issues closed with "Merged into beads-X" reason
- ✘ Source issue content NOT copied (target keeps its original content)

**Important notes:**

- Merge preserves target issue completely; only dependencies/references migrate
- If source issues have valuable content, manually copy it to target BEFORE merging
- Cannot merge in daemon mode yet (beads-190); use `--no-daemon` flag
- Operation cannot be undone (but git history preserves the original)

**Best practices:**

- Merge early to prevent dependency fragmentation
- Choose the oldest or most complete issue as merge target
- Add labels like `duplicate` to source issues before merging (for tracking)
- File a discovered-from issue if you found duplicates during work:

  ```bash
  beads create "Found duplicates during beads-X" -p 2 --deps discovered-from:beads-X --json
  ```

## Development Guidelines

### Code Standards

- **Go version**: 1.21+
- **Linting**: `golangci-lint run ./...` (baseline warnings documented in [docs/linting.md](docs/linting.md))
- **Testing**: All new features need tests (`go test ./...`)
- **Documentation**: Update relevant .md files

### File Organization

```
beads/
├── cmd/beads/              # CLI commands
├── internal/
│   ├── types/           # Core data types
│   └── storage/         # Storage layer
│       └── sqlite/      # SQLite implementation
├── examples/            # Integration examples
└── *.md                 # Documentation
```

### Before Committing

1. **Run tests**: `go test ./...`
2. **Run linter**: `golangci-lint run ./...` (ignore baseline warnings)
3. **Update docs**: If you changed behavior, update README.md or other docs
4. **Commit**: Issues auto-sync to `.beads/issues.jsonl` and import after pull

### Git Workflow

**Auto-sync provides batching!** beads automatically:

- **Exports** to JSONL after CRUD operations (30-second debounce for batching)
- **Imports** from JSONL when it's newer than DB (e.g., after `git pull`)
- **Daemon commits/pushes** every 5 seconds (if `--auto-commit` / `--auto-push` enabled)

The 30-second debounce provides a **transaction window** for batch operations - multiple issue changes within 30 seconds get flushed together, avoiding commit spam.

### Agent Session Workflow

**IMPORTANT for AI agents:** When you finish making issue changes, always run:

```bash
beads sync
```

This immediately:

1. Exports pending changes to JSONL (no 30s wait)
2. Commits to git
3. Pulls from remote
4. Imports any updates
5. Pushes to remote

**Example agent session:**

```bash
# Make multiple changes (batched in 30-second window)
beads create "Fix bug" -p 1
beads create "Add tests" -p 1
beads update beads-42 --status in_progress
beads close beads-40 --reason "Completed"

# Force immediate sync at end of session
beads sync

# Now safe to end session - everything is committed and pushed
```

**Why this matters:**

- Without `beads sync`, changes sit in 30-second debounce window
- User might think you pushed but JSONL is still dirty
- `beads sync` forces immediate flush/commit/push

**Alternative**: Install git hooks for automatic flush on commit:

```bash
# One-time setup
./examples/git-hooks/install.sh
```

This installs:

- **pre-commit** - Flushes pending changes immediately before commit (bypasses 30s debounce)
- **post-merge** - Imports updated JSONL after pull/merge (guaranteed sync)

See [examples/git-hooks/README.md](examples/git-hooks/README.md) for details.

### Git Worktrees

**WARNING: Important Limitation:** Daemon mode does not work correctly with `git worktree`.

**The Problem:**
Git worktrees share the same `.git` directory and thus share the same `.beads` database. The daemon doesn't know which branch each worktree has checked out, which can cause it to commit/push to the wrong branch.

**What you lose without daemon mode:**

- **Auto-sync** - No automatic commit/push of changes (use `beads sync` manually)
- **MCP server** - The beads-mcp server requires daemon mode for multi-repo support
- **Background watching** - No automatic detection of remote changes

**Solutions for Worktree Users:**

1. **Use `--no-daemon` flag** (recommended):

   ```bash
   beads --no-daemon ready
   beads --no-daemon create "Fix bug" -p 1
   beads --no-daemon update beads-42 --status in_progress
   ```

2. **Disable daemon via environment variable** (for entire worktree session):

   ```bash
   export BEADS_NO_DAEMON=1
   beads ready  # All commands use direct mode
   ```

3. **Disable auto-start** (less safe, still warns):

   ```bash
   export BEADS_AUTO_START_DAEMON=false
   ```

**Automatic Detection:**
beads automatically detects when you're in a worktree and shows a prominent warning if daemon mode is active. The `--no-daemon` mode works correctly with worktrees since it operates directly on the database without shared state.

**Why It Matters:**
The daemon maintains its own view of the current working directory and git state. When multiple worktrees share the same `.beads` database, the daemon may commit changes intended for one branch to a different branch, leading to confusion and incorrect git history.

### Handling Git Merge Conflicts

**With hash-based IDs (v0.20.1+), ID collisions are eliminated!** Different issues get different hash IDs, so most git merges succeed cleanly.

**When git merge conflicts occur:**
Git conflicts in `.beads/beads.jsonl` happen when the same issue is modified on both branches (different timestamps/fields). This is a **same-issue update conflict**, not an ID collision.

**Resolution:**

```bash
# After git merge creates conflict
git checkout --theirs .beads/beads.jsonl  # Accept remote version
# OR
git checkout --ours .beads/beads.jsonl    # Keep local version
# OR manually resolve in editor

# Import the resolved JSONL
beads import -i .beads/beads.jsonl

# Commit the merge
git add .beads/beads.jsonl
git commit
```

**beads automatically handles updates** - same ID with different content is a normal update operation. No special flags needed.

### Advanced: Intelligent Merge Tools

For Git merge conflicts in `.beads/issues.jsonl`, consider using **[beads-merge](https://github.com/neongreen/mono/tree/main/beads-merge)** - a specialized merge tool by @neongreen that:

- Matches issues across conflicted JSONL files
- Merges fields intelligently (e.g., combines labels, picks newer timestamps)
- Resolves conflicts automatically where possible
- Leaves remaining conflicts for manual resolution
- Works as a Git/jujutsu merge driver

**Beads-merge** helps with intelligent field-level merging during git merge. After resolving, just `beads import` to update your database.

## Current Project Status

Run `beads stats` to see overall progress.

### Active Areas

- **Core CLI**: Mature, but always room for polish
- **Examples**: Growing collection of agent integrations
- **Documentation**: Comprehensive but can always improve
- **MCP Server**: Implemented at `adapters/mcp/` with Claude Code plugin
- **Migration Tools**: Planned (see beads-6)

### 1.0 Milestone

We're working toward 1.0. Key blockers tracked in beads. Run:

```bash
beads dep tree beads-8  # Show 1.0 epic dependencies
```

## Exclusive Lock Protocol (Advanced)

**For external tools that need full database control** (e.g., CI/CD, deterministic execution systems):

The beads daemon respects exclusive locks via `.beads/.exclusive-lock` file. When this lock exists:

- Daemon skips all operations for the locked database
- External tool has complete control over git sync and database operations
- Stale locks (dead process) are automatically cleaned up

**Use case:** Tools like VibeCoder that need deterministic execution without daemon interference.

See [exclusive-lock.md](docs/exclusive-lock.md) for:

- Lock file format (JSON schema)
- Creating and releasing locks (Go/shell examples)
- Stale lock detection behavior
- Integration testing guidance

**Quick example:**

```bash
# Create lock
echo '{"holder":"my-tool","pid":'$$',"hostname":"'$(hostname)'","started_at":"'$(date -u +%Y-%m-%dT%H:%M:%SZ)'","version":"1.0.0"}' > .beads/.exclusive-lock

# Do work...
beads create "My issue" -p 1

# Release lock
rm .beads/.exclusive-lock
```

## Common Tasks

### Adding a New Command

1. Create file in `cmd/beads/`
2. Add to root command in `cmd/beads/main.go`
3. Implement with Cobra framework
4. Add `--json` flag for agent use
5. Add tests in `cmd/beads/*_test.go`
6. Document in README.md

### Adding Storage Features

1. Update schema in `internal/storage/sqlite/schema.go`
2. Add migration if needed
3. Update `internal/types/types.go` if new types
4. Implement in `internal/storage/sqlite/sqlite.go`
5. Add tests
6. Update export/import in `cmd/beads/export.go` and `cmd/beads/import.go`

### Adding Examples

1. Create directory in `examples/`
2. Add README.md explaining the example
3. Include working code
4. Link from `examples/README.md`
5. Mention in main README.md

## Questions?

- Check existing issues: `beads list`
- Look at recent commits: `git log --oneline -20`
- Read the docs: README.md, advanced.md, extending.md
- Create an issue if unsure: `beads create "Question: ..." -t task -p 2`

## Important Files

- **README.md** - Main documentation (keep this updated!)
- **extending.md** - Database extension guide
- **advanced.md** - JSONL format analysis
- **CONTRIBUTING.md** - Contribution guidelines
- **SECURITY.md** - Security policy

## Pro Tips for Agents

- Always use `--json` flags for programmatic use
- **Always run `beads sync` at end of session** to flush/commit/push immediately
- Link discoveries with `discovered-from` to maintain context
- Check `beads ready` before asking "what next?"
- Auto-sync batches changes in 30-second window - use `beads sync` to force immediate flush
- Use `--no-auto-flush` or `--no-auto-import` to disable automatic sync if needed
- Use `beads dep tree` to understand complex dependencies
- Priority 0-1 issues are usually more important than 2-4
- Use `--dry-run` to preview import changes before applying
- Hash IDs eliminate collisions - same ID with different content is a normal update
- Use `--id` flag with `beads create` to partition ID space for parallel workers (e.g., `worker1-100`, `worker2-500`)

## Building and Testing

```bash
# Build
go build -o beads ./cmd/beads

# Test
go test ./...

# Test with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run locally
./beads init --prefix test
./beads create "Test issue" -p 1
./beads ready
```

## Version Management

**IMPORTANT**: When the user asks to "bump the version" or mentions a new version number (e.g., "bump to 0.9.3"), use the version bump script:

```bash
# Preview changes (shows diff, doesn't commit)
./scripts/bump-version.sh 0.9.3

# Auto-commit the version bump
./scripts/bump-version.sh 0.9.3 --commit
git push origin main
```

**What it does:**

- Updates ALL version files (CLI, plugin, MCP server, docs) in one command
- Validates semantic versioning format
- Shows diff preview
- Verifies all versions match after update
- Creates standardized commit message

**User will typically say:**

- "Bump to 0.9.3"
- "Update version to 1.0.0"
- "Rev the project to 0.9.4"
- "Increment the version"

**You should:**

1. Run `./scripts/bump-version.sh <version> --commit`
2. Push to GitHub
3. Confirm all versions updated correctly

**Files updated automatically:**

- `cmd/beads/version.go` - CLI version
- `.claude-plugin/plugin.json` - Plugin version
- `.claude-plugin/marketplace.json` - Marketplace version
- `adapters/mcp/pyproject.toml` - MCP server version
- `README.md` - Documentation version
- `plugin.md` - Version requirements

**Why this matters:** We had version mismatches (beads-66) when only `version.go` was updated. This script prevents that by updating all components atomically.

See `scripts/README.md` for more details.

## Release Process (Maintainers)

1. Bump version with `./scripts/bump-version.sh <version> --commit`
2. Update CHANGELOG.md with release notes
3. Run full test suite: `go test ./...`
4. Push version bump: `git push origin main`
5. Tag release: `git tag v<version>`
6. Push tag: `git push origin v<version>`
7. GitHub Actions handles the rest

---

**Remember**: We're building this tool to help AI agents like you! If you find the workflow confusing or have ideas for improvement, create an issue with your feedback.

Happy coding!

<!-- beads onboard section -->
## Issue Tracking with beads (beads)

**IMPORTANT**: This project uses **beads (beads)** for ALL issue tracking. Do NOT use markdown TODOs, task lists, or other tracking methods.

### Why beads?

- Dependency-aware: Track blockers and relationships between issues
- Git-friendly: Auto-syncs to JSONL for version control
- Agent-optimized: JSON output, ready work detection, discovered-from links
- Prevents duplicate tracking systems and confusion

### Quick Start

**FIRST TIME?** Just run `beads init` - it auto-imports issues from git:

```bash
beads init --prefix beads
```

**Check for ready work:**

```bash
beads ready --json
```

**Create new issues:**

```bash
beads create "Issue title" -t bug|feature|task -p 0-4 --json
beads create "Issue title" -p 1 --deps discovered-from:beads-123 --json
```

**Claim and update:**

```bash
beads update beads-42 --status in_progress --json
beads update beads-42 --priority 1 --json
```

**Complete work:**

```bash
beads close beads-42 --reason "Completed" --json
```

### Issue Types

- `bug` - Something broken
- `feature` - New functionality
- `task` - Work item (tests, docs, refactoring)
- `epic` - Large feature with subtasks
- `chore` - Maintenance (dependencies, tooling)

### Priorities

- `0` - Critical (security, data loss, broken builds)
- `1` - High (major features, important bugs)
- `2` - Medium (default, nice-to-have)
- `3` - Low (polish, optimization)
- `4` - Backlog (future ideas)

### Workflow for AI Agents

1. **Check ready work**: `beads ready` shows unblocked issues
2. **Claim your task**: `beads update <id> --status in_progress`
3. **Work on it**: Implement, test, document
4. **Discover new work?** Create linked issue:
   - `beads create "Found bug" -p 1 --deps discovered-from:<parent-id>`
5. **Complete**: `beads close <id> --reason "Done"`

### Auto-Sync

beads automatically syncs with git:

- Exports to `.beads/issues.jsonl` after changes (5s debounce)
- Imports from JSONL when newer (e.g., after `git pull`)
- No manual export/import needed!

### MCP Server (Recommended)

If using Claude or MCP-compatible clients, install the beads MCP server:

```bash
uv tool install beads-mcp
```

Add to MCP config (e.g., `~/.config/claude/config.json`):

```json
{
  "beads": {
    "command": "beads-mcp",
    "args": []
  }
}
```

Then use `mcp__beads__*` functions instead of CLI commands.

### Important Rules

**Project Work (MUST use beads):**

- ✔ Create beads issue for ANY code/docs/config changes
- ✔ Always use `--json` flag for programmatic use
- ✔ Link discovered work with `discovered-from` dependencies
- ✔ Check `beads ready` before asking "what should I work on?"
- ✔ Run `beads sync` at end of session to commit/push

**What NOT to do:**

- ✘ Do NOT create markdown TODO lists
- ✘ Do NOT use external issue trackers
- ✘ Do NOT duplicate tracking systems
- ✘ Do NOT commit work without a beads issue

**Ephemeral Session Planning (escape hatch):**

- ✔ Native agent tools (like TodoWrite) are OK for temporary task breakdown
- ✔ Use for complex multi-step work that helps YOU plan, not project history
- ✔ When in doubt, use beads

For more details, see README.md and quickstart.md.
<!-- /beads onboard section -->
