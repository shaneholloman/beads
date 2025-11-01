# CLI Reference

Complete command reference for beads (beads) CLI tool. All commands support `--json` flag for structured output.

## Contents

- [Quick Reference](#quick-reference)
- [Global Flags](#global-flags)
- [Core Commands](#core-commands)
  - [beads ready](#beads-ready) - Find unblocked work
  - [beads create](#beads-create) - Create new issues
  - [beads update](#beads-update) - Update issue status, priority, assignee
  - [beads close](#beads-close) - Close completed work
  - [beads show](#beads-show) - Show issue details
  - [beads list](#beads-list) - List issues with filters
- [Dependency Commands](#dependency-commands)
  - [beads dep add](#beads-dep-add) - Create dependencies
  - [beads dep tree](#beads-dep-tree) - Visualize dependency trees
  - [beads dep cycles](#beads-dep-cycles) - Detect circular dependencies
- [Monitoring Commands](#monitoring-commands)
  - [beads stats](#beads-stats) - Project statistics
  - [beads blocked](#beads-blocked) - Find blocked work
- [Data Management Commands](#data-management-commands)
  - [beads export](#beads-export) - Export database to JSONL
  - [beads import](#beads-import) - Import issues from JSONL
- [Setup Commands](#setup-commands)
  - [beads init](#beads-init) - Initialize database
  - [beads quickstart](#beads-quickstart) - Show quick start guide
- [Common Workflows](#common-workflows)
- [JSON Output](#json-output)
- [Database Auto-Discovery](#database-auto-discovery)
- [Git Integration](#git-integration)
- [Tips](#tips)

## Quick Reference

| Command | Purpose | Key Flags |
|---------|---------|-----------|
| `beads ready` | Find unblocked work | `--priority`, `--assignee`, `--limit`, `--json` |
| `beads list` | List all issues with filters | `--status`, `--priority`, `--type`, `--assignee` |
| `beads show` | Show issue details | `--json` |
| `beads create` | Create new issue | `-t`, `-p`, `-d`, `--design`, `--acceptance` |
| `beads update` | Update existing issue | `--status`, `--priority`, `--design` |
| `beads close` | Close completed issue | `--reason` |
| `beads dep add` | Add dependency | `--type` (blocks, related, parent-child, discovered-from) |
| `beads dep tree` | Visualize dependency tree | (no flags) |
| `beads dep cycles` | Detect circular dependencies | (no flags) |
| `beads stats` | Get project statistics | `--json` |
| `beads blocked` | Find blocked issues | `--json` |
| `beads export` | Export issues to JSONL | `--json` |
| `beads import` | Import issues from JSONL | `--resolve-collisions` |
| `beads init` | Initialize beads in directory | `--prefix` |
| `beads quickstart` | Show quick start guide | (no flags) |

## Global Flags

Available for all commands:

```bash
--json                 # Output in JSON format
--db /path/to/db       # Specify database path (default: auto-discover)
--actor "name"         # Actor name for audit trail
--no-auto-flush        # Disable automatic JSONL sync
--no-auto-import       # Disable automatic JSONL import
```

## Core Commands

### beads ready

Find tasks with no blockers - ready to be worked on.

```bash
beads ready                      # All ready work
beads ready --json               # JSON format
beads ready --priority 0         # Only priority 0 (critical)
beads ready --assignee alice     # Only assigned to alice
beads ready --limit 5            # Limit to 5 results
```

**Use at session start** to see available work.

---

### beads create

Create a new issue with optional metadata.

```bash
beads create "Title"
beads create "Title" -t bug -p 0
beads create "Title" -d "Description"
beads create "Title" --design "Design notes"
beads create "Title" --acceptance "Definition of done"
beads create "Title" --assignee alice
```

**Flags**:

- `-t, --type`: task (default), bug, feature, epic, chore
- `-p, --priority`: 0-3 (default: 2)
- `-d, --description`: Issue description
- `--design`: Design notes
- `--acceptance`: Acceptance criteria
- `--assignee`: Who should work on this

---

### beads update

Update an existing issue's metadata.

```bash
beads update issue-123 --status in_progress
beads update issue-123 --priority 0
beads update issue-123 --design "Decided to use Redis"
beads update issue-123 --acceptance "Tests passing"
```

**Status values**: open, in_progress, blocked, closed

---

### beads close

Close (complete) an issue.

```bash
beads close issue-123
beads close issue-123 --reason "Implemented in PR #42"
beads close issue-1 issue-2 issue-3 --reason "Bulk close"
```

**Note**: Closed issues remain in database for history.

---

### beads show

Show detailed information about a specific issue.

```bash
beads show issue-123
beads show issue-123 --json
```

Shows: all fields, dependencies, dependents, audit history.

---

### beads list

List all issues with optional filters.

```bash
beads list                          # All issues
beads list --status open            # Only open
beads list --priority 0             # Critical
beads list --type bug               # Only bugs
beads list --assignee alice         # By assignee
beads list --status closed --limit 10  # Recent completions
```

---

## Dependency Commands

### beads dep add

Add a dependency between issues.

```bash
beads dep add from-issue to-issue                      # blocks (default)
beads dep add from-issue to-issue --type blocks
beads dep add from-issue to-issue --type related
beads dep add epic-id task-id --type parent-child
beads dep add original-id found-id --type discovered-from
```

**Dependency types**:

1. **blocks**: from-issue blocks to-issue (hard blocker)
2. **related**: Soft link (no blocking)
3. **parent-child**: Epic/subtask hierarchy
4. **discovered-from**: Tracks origin of discovery

---

### beads dep tree

Visualize full dependency tree for an issue.

```bash
beads dep tree issue-123
```

Shows all dependencies and dependents in tree format.

---

### beads dep cycles

Detect circular dependencies.

```bash
beads dep cycles
```

Finds dependency cycles that would prevent work from being ready.

---

## Monitoring Commands

### beads stats

Get project statistics.

```bash
beads stats
beads stats --json
```

Returns: total, open, in_progress, closed, blocked, ready, avg lead time.

---

### beads blocked

Get blocked issues with blocker information.

```bash
beads blocked
beads blocked --json
```

Use to identify bottlenecks when ready list is empty.

---

## Data Management Commands

### beads export

Export all issues to JSONL format.

```bash
beads export > issues.jsonl
beads export --json  # Same output, explicit flag
```

**Use cases:**

- Manual backup before risky operations
- Sharing issues across databases
- Version control / git tracking
- Data migration or analysis

**Note**: beads auto-exports to `.beads/*.jsonl` after each operation (5s debounce). Manual export is rarely needed.

---

### beads import

Import issues from JSONL format.

```bash
beads import < issues.jsonl
beads import -i issues.jsonl --dry-run  # Preview changes
```

**Behavior with hash-based IDs (v0.20.1+):**

- Same ID = update operation (hash IDs remain stable)
- Different issues get different hash IDs (no collisions)
- Import automatically applies updates to existing issues

**Use `--dry-run` to preview:**

```bash
beads import -i issues.jsonl --dry-run
# Shows: new issues, updates, exact matches
```

**Use cases:**

- **Syncing after git pull** - daemon auto-imports, manual rarely needed
- **Merging databases** - import issues from another database
- **Restoring from backup** - reimport JSONL to restore state

---

## Setup Commands

### beads init

Initialize beads in current directory.

```bash
beads init                    # Auto-detect prefix
beads init --prefix api       # Custom prefix
```

Creates `.beads/` directory and database.

---

### beads quickstart

Show comprehensive quick start guide.

```bash
beads quickstart
```

Displays built-in reference for command syntax and workflows.

---

## Common Workflows

### Session Start

```bash
beads ready --json
beads show issue-123
beads update issue-123 --status in_progress
```

### Discovery During Work

```bash
beads create "Found: bug in auth" -t bug
beads dep add current-issue new-issue --type discovered-from
```

### Completing Work

```bash
beads close issue-123 --reason "Implemented with tests passing"
beads ready  # See what unblocked
```

### Planning Epic

```bash
beads create "OAuth Integration" -t epic
beads create "Set up credentials" -t task
beads create "Implement flow" -t task

beads dep add oauth-epic oauth-creds --type parent-child
beads dep add oauth-epic oauth-flow --type parent-child
beads dep add oauth-creds oauth-flow  # creds blocks flow

beads dep tree oauth-epic
```

---

## JSON Output

All commands support `--json` for structured output:

```bash
beads ready --json
beads show issue-123 --json
beads list --status open --json
beads stats --json
```

Use when parsing programmatically or extracting specific fields.

---

## Database Auto-Discovery

beads finds database in this order:

1. `--db` flag: `beads ready --db /path/to/db.db`
2. `$BEADS_DB` environment variable
3. `.beads/*.db` in current directory or ancestors
4. `~/.beads/default.db` as fallback

**Project-local** (`.beads/`): Project-specific work, git-tracked

**Global fallback** (`~/.beads/`): Cross-project tracking, personal tasks

---

## Git Integration

beads automatically syncs with git:

- **After each operation**: Exports to JSONL (5s debounce)
- **After git pull**: Imports from JSONL if newer than DB

**Files**:

- `.beads/*.jsonl` - Source of truth (git-tracked)
- `.beads/*.db` - Local cache (gitignored)

### Git Integration Troubleshooting

**Problem: `.gitignore` ignores entire `.beads/` directory**

**Symptom**: JSONL file not tracked in git, can't commit beads

**Cause**: Incorrect `.gitignore` pattern blocks everything

**Fix**:

```bash
# Check .gitignore
cat .gitignore | grep beads

# ✘ WRONG (ignores everything including JSONL):
.beads/

# ✔ CORRECT (ignores only SQLite cache):
.beads/*.db
.beads/*.db-*
```

**After fixing**: Remove the `.beads/` line and add the specific patterns. Then `git add .beads/issues.jsonl`.

---

### Permission Troubleshooting

**Problem: beads commands prompt for permission despite whitelist**

**Symptom**: `beads` commands ask for confirmation even with `Bash(beads:*)` in settings.local.json

**Root Cause**: Wildcard patterns in settings.local.json don't actually work - not for beads, not for git, not for any Bash commands. This is a general Claude Code limitation, not beads-specific.

**How It Actually Works**:

- Individual command approvals (like `Bash(beads ready)`) DO persist across sessions
- These are stored server-side by Claude Code, not in local config files
- Commands like `git status` work without prompting because they've been individually approved many times, creating the illusion of a working wildcard pattern

**Permanent Solution**:

1. Trigger each beads subcommand you use frequently (see command list below)
2. When prompted, click "Yes, and don't ask again" (NOT "Allow this time")
3. That specific command will be permanently approved across all future sessions

**Common beads Commands to Approve**:

```bash
beads ready
beads list
beads stats
beads blocked
beads export
beads version
beads quickstart
beads dep cycles
beads --help
beads [command] --help  # For any subcommand help
```

**Note**: Dynamic commands with arguments (like `beads show <issue-id>`, `beads create "title"`) must be approved per-use since arguments vary. Only static commands can be permanently whitelisted.

---

## Tips

**Use JSON for parsing**:

```bash
beads ready --json | jq '.[0].id'
```

**Bulk operations**:

```bash
beads close issue-1 issue-2 issue-3 --reason "Sprint complete"
```

**Quick filtering**:

```bash
beads list --status open --priority 0 --type bug
```

**Built-in help**:

```bash
beads quickstart       # Comprehensive guide
beads create --help    # Command-specific help
```
