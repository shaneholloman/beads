# Advanced beads Features

This guide covers advanced features for power users and specific use cases.

## Table of Contents

- [Renaming Prefix](#renaming-prefix)
- [Merging Duplicate Issues](#merging-duplicate-issues)
- [Git Worktrees](#git-worktrees)
- [Custom Git Hooks](#custom-git-hooks)
- [Extensible Database](#extensible-database)
- [Architecture: Daemon vs MCP vs Beads](#architecture-daemon-vs-mcp-vs-beads)

## Renaming Prefix

Change the issue prefix for all issues in your database. This is useful if your prefix is too long or you want to standardize naming.

```sh
# Preview changes without applying
beads rename-prefix kw- --dry-run

# Rename from current prefix to new prefix
beads rename-prefix kw-

# JSON output
beads rename-prefix kw- --json
```

The rename operation:

- Updates all issue IDs (e.g., `knowledge-work-1` → `kw-1`)
- Updates all text references in titles, descriptions, design notes, etc.
- Updates dependencies and labels
- Updates the counter table and config

**Prefix validation rules:**

- Max length: 8 characters
- Allowed characters: lowercase letters, numbers, hyphens
- Must start with a letter
- Must end with a hyphen (or will be trimmed to add one)
- Cannot be empty or just a hyphen

Example workflow:

```sh
# You have issues like knowledge-work-1, knowledge-work-2, etc.
beads list  # Shows knowledge-work-* issues

# Preview the rename
beads rename-prefix kw- --dry-run

# Apply the rename
beads rename-prefix kw-

# Now you have kw-1, kw-2, etc.
beads list  # Shows kw-* issues
```

## Duplicate Detection

Find issues with identical content using automated duplicate detection:

```sh
# Find all content duplicates in the database
beads duplicates

# Show duplicates in JSON format
beads duplicates --json

# Automatically merge all duplicates
beads duplicates --auto-merge

# Preview what would be merged
beads duplicates --dry-run

# Detect duplicates during import
beads import -i issues.jsonl --resolve-collisions --dedupe-after
```

**How it works:**

- Groups issues by content hash (title, description, design, acceptance criteria)
- Only groups issues with matching status (open with open, closed with closed)
- Chooses merge target by reference count (most referenced) or smallest ID
- Reports duplicate groups with suggested merge commands

**Example output:**

```
Found 3 duplicate group(s):

━━ Group 1: Fix authentication bug
→ beads-10 (open, P1, 5 references)
  beads-42 (open, P1, 0 references)
  Suggested: beads merge beads-42 --into beads-10

Run with --auto-merge to execute all suggested merges
```

**AI Agent Workflow:**

1. **Periodic scans**: Run `beads duplicates` to check for duplicates
2. **During import**: Use `--dedupe-after` to detect duplicates after collision resolution
3. **Auto-merge**: Use `--auto-merge` to automatically consolidate duplicates
4. **Manual review**: Use `--dry-run` to preview merges before executing

## Merging Duplicate Issues

Consolidate duplicate issues into a single issue while preserving dependencies and references:

```sh
# Merge beads-42 and beads-43 into beads-41
beads merge beads-42 beads-43 --into beads-41

# Merge multiple duplicates at once
beads merge beads-10 beads-11 beads-12 --into beads-10

# Preview merge without making changes
beads merge beads-42 beads-43 --into beads-41 --dry-run

# JSON output
beads merge beads-42 beads-43 --into beads-41 --json
```

**What the merge command does:**

1. **Validates** all issues exist and prevents self-merge
2. **Closes** source issues with reason `Merged into beads-X`
3. **Migrates** all dependencies from source issues to target
4. **Updates** text references across all issue descriptions, notes, design, and acceptance criteria

**Example workflow:**

```sh
# You discover beads-42 and beads-43 are duplicates of beads-41
beads show beads-41 beads-42 beads-43

# Preview the merge
beads merge beads-42 beads-43 --into beads-41 --dry-run

# Execute the merge
beads merge beads-42 beads-43 --into beads-41
# ✔ Merged 2 issue(s) into beads-41

# Verify the result
beads show beads-41  # Now has dependencies from beads-42 and beads-43
beads dep tree beads-41  # Shows unified dependency tree
```

**Important notes:**

- Source issues are permanently closed (status: `closed`)
- All dependencies pointing to source issues are redirected to target
- Text references like "see beads-42" are automatically rewritten to "see beads-41"
- Operation cannot be undone (but git history preserves the original state)
- Not yet supported in daemon mode (use `--no-daemon` flag)

**AI Agent Workflow:**

When agents discover duplicate issues, they should:

1. Search for similar issues: `beads list --json | grep "similar text"`
2. Compare issue details: `beads show beads-41 beads-42 --json`
3. Merge duplicates: `beads merge beads-42 --into beads-41`
4. File a discovered-from issue if needed: `beads create "Found duplicates during beads-X" --deps discovered-from:beads-X`

## Git Worktrees

**WARNING: Important Limitation:** Daemon mode does not work correctly with `git worktree`.

**The Problem:**
Git worktrees share the same `.git` directory and thus share the same `.beads` database. The daemon doesn't know which branch each worktree has checked out, which can cause it to commit/push to the wrong branch.

**What you lose without daemon mode:**

- **Auto-sync** - No automatic commit/push of changes (use `beads sync` manually)
- **MCP server** - The mcp-beads server requires daemon mode for multi-repo support
- **Background watching** - No automatic detection of remote changes

**Solutions for Worktree Users:**

1. **Use `--no-daemon` flag** (recommended):

   ```sh
   beads --no-daemon ready
   beads --no-daemon create "Fix bug" -p 1
   beads --no-daemon update beads-42 --status in_progress
   ```

2. **Disable daemon via environment variable** (for entire worktree session):

   ```sh
   export BEADS_NO_DAEMON=1
   beads ready  # All commands use direct mode
   ```

3. **Disable auto-start** (less safe, still warns):

   ```sh
   export BEADS_AUTO_START_DAEMON=false
   ```

**Automatic Detection:**
beads automatically detects when you're in a worktree and shows a prominent warning if daemon mode is active. The `--no-daemon` mode works correctly with worktrees since it operates directly on the database without shared state.

**Why It Matters:**
The daemon maintains its own view of the current working directory and git state. When multiple worktrees share the same `.beads` database, the daemon may commit changes intended for one branch to a different branch, leading to confusion and incorrect git history.

## Handling Git Merge Conflicts

**With hash-based IDs (v0.20.1+), ID collisions are eliminated.** Different issues get different hash IDs, so concurrent creation doesn't cause conflicts.

### Understanding Same-ID Scenarios

When you encounter the same ID during import, it's an **update operation**, not a collision:

- Hash IDs are content-based and remain stable across updates
- Same ID + different fields = normal update to existing issue
- beads automatically applies updates when importing

**Preview changes before importing:**

```sh
# After git merge or pull
beads import -i .beads/issues.jsonl --dry-run

# Output shows:
# Exact matches (idempotent): 15
# New issues: 5
# Updates: 3
#
# Issues to be updated:
#   beads-a3f2: Fix authentication (changed: priority, status)
#   beads-b8e1: Add feature (changed: description)
```

### Git Merge Conflicts

The conflicts you'll encounter are **git merge conflicts** in the JSONL file when the same issue was modified on both branches (different timestamps/fields). This is not an ID collision.

**Resolution:**

```sh
# After git merge creates conflict
git checkout --theirs .beads/beads.jsonl  # Accept remote version
# OR
git checkout --ours .beads/beads.jsonl    # Keep local version
# OR manually resolve in editor (keep line with newer updated_at)

# Import the resolved JSONL
beads import -i .beads/beads.jsonl

# Commit the merge
git add .beads/beads.jsonl
git commit
```

### Advanced: Intelligent Merge Tools

For Git merge conflicts in `.beads/issues.jsonl`, consider using **[beads-merge](https://github.com/neongreen/mono/tree/main/beads-merge)** - a specialized merge tool by @neongreen that:

- Matches issues across conflicted JSONL files
- Merges fields intelligently (e.g., combines labels, picks newer timestamps)
- Resolves conflicts automatically where possible
- Leaves remaining conflicts for manual resolution
- Works as a Git/jujutsu merge driver

After using beads-merge to resolve the git conflict, just run `beads import` to update your database.

## Custom Git Hooks

For immediate export (no 5-second wait) and guaranteed import after git operations, install the git hooks:

### Using the Installer

```sh
cd examples/git-hooks
./install.sh
```

### Manual Setup

Create `.git/hooks/pre-commit`:

```sh
#!/bin/bash
beads export -o .beads/issues.jsonl
git add .beads/issues.jsonl
```

Create `.git/hooks/post-merge`:

```sh
#!/bin/bash
beads import -i .beads/issues.jsonl
```

Create `.git/hooks/post-checkout`:

```sh
#!/bin/bash
beads import -i .beads/issues.jsonl
```

Make hooks executable:

```sh
chmod +x .git/hooks/pre-commit .git/hooks/post-merge .git/hooks/post-checkout
```

**Note:** Auto-sync is already enabled by default, so git hooks are optional. They're useful if you need immediate export or guaranteed import after git operations.

## Extensible Database

beads uses SQLite, which you can extend with your own tables and queries. This allows you to:

- Add custom metadata to issues
- Build integrations with other tools
- Implement custom workflows
- Create reports and analytics

**See [extending.md](./extending.md) for complete documentation:**

- Database schema and structure
- Adding custom tables
- Joining with issue data
- Example integrations
- Best practices

**Example use case:**

```sql
-- Add time tracking table
CREATE TABLE time_entries (
    id INTEGER PRIMARY KEY,
    issue_id TEXT NOT NULL,
    duration_minutes INTEGER NOT NULL,
    recorded_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(issue_id) REFERENCES issues(id)
);

-- Query total time per issue
SELECT i.id, i.title, SUM(t.duration_minutes) as total_minutes
FROM issues i
LEFT JOIN time_entries t ON i.id = t.issue_id
GROUP BY i.id;
```

## Architecture: Daemon vs MCP vs Beads

Understanding the role of each component:

### Beads (Core)

- **SQLite database** - The source of truth for all issues, dependencies, labels
- **Storage layer** - CRUD operations, dependency resolution, collision detection
- **Business logic** - Ready work calculation, merge operations, import/export
- **CLI commands** - Direct database access via `beads` command

### Local Daemon (Per-Project)

- **Lightweight RPC server** - Runs at `.beads/beads.sock` in each project
- **Auto-sync coordination** - Debounced export (5s), git integration, import detection
- **Process isolation** - Each project gets its own daemon for database safety
- **LSP model** - Similar to language servers, one daemon per workspace
- **No global daemon** - Removed in v0.16.0 to prevent cross-project pollution
- **Exclusive lock support** - External tools can prevent daemon interference (see [exclusive-lock.md](./exclusive-lock.md))

### MCP Server (Optional)

- **Protocol adapter** - Translates MCP calls to daemon RPC or direct CLI
- **Workspace routing** - Finds correct `.beads/beads.sock` based on working directory
- **Stateless** - Doesn't cache or store any issue data itself
- **Editor integration** - Makes beads available to Claude, Cursor, and other MCP clients
- **Single instance** - One MCP server can route to multiple project daemons

**Key principle**: The daemon and MCP server are thin layers. All heavy lifting (dependency graphs, collision resolution, merge logic) happens in the core beads storage layer.

**Why per-project daemons?**

- Complete database isolation between projects
- Git worktree safety (each worktree can disable daemon independently)
- No risk of committing changes to wrong branch
- Simpler mental model - one project, one database, one daemon
- Follows LSP/language server architecture patterns

## Next Steps

- **[README.md](../README.md)** - Core features and quick start
- **[troubleshooting.md](./troubleshooting.md)** - Common issues and solutions
- **[faq.md](./faq.md)** - Frequently asked questions
- **[config.md](./config.md)** - Configuration system guide
- **[extending.md](./extending.md)** - Database extension patterns
