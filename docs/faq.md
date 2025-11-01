# Frequently Asked Questions

Common questions about beads (beads) and how to use it effectively.

## General Questions

### What is beads?

beads is a lightweight, git-based issue tracker designed for AI coding agents. It provides dependency-aware task management with automatic sync across machines via git.

### Why not just use GitHub Issues?

GitHub Issues + gh CLI can approximate some features, but fundamentally cannot replicate what AI agents need:

**Key Differentiators:**

1. **Typed Dependencies with Semantics**
   - beads: Four types (`blocks`, `related`, `parent-child`, `discovered-from`) with different behaviors
   - GH: Only "blocks/blocked by" links, no semantic enforcement, no `discovered-from` for agent work discovery

2. **Deterministic Ready-Work Detection**
   - beads: `beads ready` computes transitive blocking offline in ~10ms, no network required
   - GH: No built-in "ready" concept; would require custom GraphQL + sync service + ongoing maintenance

3. **Git-First, Offline, Branch-Scoped Task Memory**
   - beads: Works offline, issues live on branches, mergeable with code via `beads import --resolve-collisions`
   - GH: Cloud-first, requires network/auth, global per-repo, no branch-scoped task state

4. **AI-Resolvable Conflicts & Duplicate Merge**
   - beads: Automatic collision resolution, duplicate merge with dependency consolidation and reference rewriting
   - GH: Manual close-as-duplicate, no safe bulk merge, no cross-reference updates

5. **Extensible Local Database**
   - beads: Add SQL tables and join with issue data locally (see [extending.md](extending.md))
   - GH: No local database; would need to mirror data externally

6. **Agent-Native APIs**
   - beads: Consistent `--json` on all commands, dedicated MCP server with auto workspace detection
   - GH: Mixed JSON/text output, GraphQL requires custom queries, no agent-focused MCP layer

**When to use each:** GitHub Issues excels for human teams in web UI with cross-repo dashboards and integrations. beads excels for AI agents needing offline, git-synchronized task memory with graph semantics and deterministic queries.

See [GitHub issue #125](https://github.com/shaneholloman/beads/issues/125) for detailed comparison.

### How is this different from Taskwarrior?

Taskwarrior is excellent for personal task management, but beads is built for AI agents:

- **Explicit agent semantics**: `discovered-from` dependency type, `beads ready` for queue management
- **JSON-first design**: Every command has `--json` output
- **Git-native sync**: No sync server setup required
- **Merge-friendly JSONL**: One issue per line, AI-resolvable conflicts
- **Extensible SQLite**: Add your own tables without forking

### Can I use beads without AI agents?

Absolutely! beads is a great CLI issue tracker for humans too. The `beads ready` command is useful for anyone managing dependencies. Think of it as "Taskwarrior meets git."

### Is this production-ready?

**Current status: Alpha (v0.9.11)**

beads is in active development and being dogfooded on real projects. The core functionality (create, update, dependencies, ready work, collision resolution) is stable and well-tested. However:

- **WARNING: Alpha software** - No 1.0 release yet
- **WARNING: API may change** - Command flags and JSONL format may evolve before 1.0
- ✔ **Safe for development** - Use for development/internal projects
- ✔ **Data is portable** - JSONL format is human-readable and easy to migrate
- **Rapid iteration** - Expect frequent updates and improvements

**When to use beads:**

- ✔ AI-assisted development workflows
- ✔ Internal team projects
- ✔ Personal productivity with dependency tracking
- ✔ Experimenting with agent-first tools

**When to wait:**

- ✘ Mission-critical production systems (wait for 1.0)
- ✘ Large enterprise deployments (wait for stability guarantees)
- ✘ Long-term archival (though JSONL makes migration easy)

Follow the repo for updates and the path to 1.0!

## Usage Questions

### Why hash-based IDs? Why not sequential?

**Hash IDs eliminate collisions** when multiple agents or branches create issues concurrently.

**The problem with sequential IDs:**

```sh
# Branch A creates beads-10
git checkout -b feature-auth
beads create "Add OAuth"  # Sequential ID: beads-10

# Branch B also creates beads-10
git checkout -b feature-payments
beads create "Add Stripe"  # Collision! Same sequential ID: beads-10

# Merge conflict!
git merge feature-auth   # Two different issues, same ID
```

**Hash IDs solve this:**

```sh
# Branch A
beads create "Add OAuth"  # Hash ID: beads-a1b2 (from random UUID)

# Branch B
beads create "Add Stripe"  # Hash ID: beads-f14c (different UUID, different hash)

# Clean merge!
git merge feature-auth   # No collision, different IDs
```

**Progressive length scaling:**

- 4 chars (0-500 issues): `beads-a1b2`
- 5 chars (500-1,500 issues): `beads-f14c3`
- 6 chars (1,500+ issues): `beads-3e7a5b`

beads automatically extends hash length as your database grows to maintain low collision probability.

### What are hierarchical child IDs?

**Hierarchical IDs** (e.g., `beads-a3f8e9.1`, `beads-a3f8e9.2`) provide human-readable structure for epics and their subtasks.

**Example:**

```sh
# Create epic (generates parent hash)
beads create "Auth System" -t epic -p 1
# Returns: beads-a3f8e9

# Create children (auto-numbered .1, .2, .3)
beads create "Login UI" -p 1       # beads-a3f8e9.1
beads create "Validation" -p 1     # beads-a3f8e9.2
beads create "Tests" -p 1          # beads-a3f8e9.3
```

**Benefits:**

- Parent hash ensures unique namespace (no cross-epic collisions)
- Sequential child IDs are human-friendly
- Up to 3 levels of nesting supported
- Clear visual grouping in issue lists

**When to use:**

- Epics with multiple related tasks
- Large features with sub-features
- Work breakdown structures

**When NOT to use:**

- Simple one-off tasks (use regular hash IDs)
- Cross-cutting dependencies (use `beads dep add` instead)

### Should I run beads init or have my agent do it?

**Either works!** But use the right flag:

**Humans:**

```sh
beads init  # Interactive - prompts for git hooks
```

**Agents:**

```sh
beads init --quiet  # Non-interactive - auto-installs hooks, no prompts
```

**Workflow for humans:**

```sh
# Clone existing project with beads:
git clone <repo>
cd <repo>
beads init  # Auto-imports from .beads/issues.jsonl

# Or initialize new project:
cd ~/my-project
beads init  # Creates .beads/, sets up daemon
git add .beads/
git commit -m "Initialize beads"
```

**Workflow for agents setting up repos:**

```sh
git clone <repo>
cd <repo>
beads init --quiet  # No prompts, auto-installs hooks
beads ready --json  # Start using beads normally
```

### Do I need to run export/import manually?

**No! Sync is automatic by default.**

beads automatically:

- **Exports** to JSONL after CRUD operations (5-second debounce)
- **Imports** from JSONL when it's newer than DB (e.g., after `git pull`)

**How auto-import works:** The first beads command after `git pull` detects that `.beads/issues.jsonl` is newer than the database and automatically imports it. There's no background daemon watching for changes - the check happens when you run a beads command.

**Optional**: For immediate export (no 5-second wait) and guaranteed import after git operations, install the git hooks:

```sh
cd examples/git-hooks && ./install.sh
```

**Disable auto-sync** if needed:

```sh
beads --no-auto-flush create "Issue"   # Disable auto-export
beads --no-auto-import list            # Disable auto-import check
```

### What if my database feels stale after git pull?

Just run any beads command - it will auto-import:

```sh
git pull
beads ready     # Automatically imports fresh data from git
beads list      # Also triggers auto-import if needed
beads sync      # Explicit sync command for manual control
```

The auto-import check is fast (<5ms) and only imports when the JSONL file is newer than the database. If you want guaranteed immediate sync without waiting for the next command, use the git hooks (see `examples/git-hooks/`).

### Can I track issues for multiple projects?

**Yes! Each project is completely isolated.** beads uses project-local databases:

```sh
cd ~/project1 && beads init --prefix proj1
cd ~/project2 && beads init --prefix proj2
```

Each project gets its own `.beads/` directory with its own database and JSONL file. beads auto-discovers the correct database based on your current directory (walks up like git).

**Multi-project scenarios work seamlessly:**

- Multiple agents working on different projects simultaneously → No conflicts
- Same machine, different repos → Each finds its own `.beads/*.db` automatically
- Agents in subeadsirectories → beads walks up to find the project root (like git)
- **Per-project daemons** → Each project gets its own daemon at `.beads/beads.sock` (LSP model)

**Limitation:** Issues cannot reference issues in other projects. Each database is isolated by design. If you need cross-project tracking, initialize beads in a parent directory that contains both projects.

**Example:** Multiple agents, multiple projects, same machine:

```sh
# Agent 1 working on web app
cd ~/work/webapp && beads ready --json    # Uses ~/work/webapp/.beads/webapp.db

# Agent 2 working on API
cd ~/work/api && beads ready --json       # Uses ~/work/api/.beads/api.db

# No conflicts! Completely isolated databases and daemons.
```

**Architecture:** beads uses per-project daemons (like LSP/language servers) for complete database isolation. See [advanced.md#architecture-daemon-vs-mcp-vs-beads](advanced.md#architecture-daemon-vs-mcp-vs-beads).

### What happens if two agents work on the same issue?

The last agent to export/commit wins. This is the same as any git-based workflow. To prevent conflicts:

- Have agents claim work with `beads update <id> --status in_progress`
- Query by assignee: `beads ready --assignee agent-name`
- Review git diffs before merging

For true multi-agent coordination, you'd need additional tooling (like locks or a coordination server). beads handles the simpler case: multiple humans/agents working on different tasks, syncing via git.

### Why JSONL instead of JSON?

- ✔ **Git-friendly**: One line per issue = clean diffs
- ✔ **Mergeable**: Concurrent appends rarely conflict
- ✔ **Human-readable**: Easy to review changes
- ✔ **Scriptable**: Use `jq`, `grep`, or any text tools
- ✔ **Portable**: Export/import between databases

See [advanced.md](advanced.md) for detailed analysis.

### How do I handle merge conflicts?

When two developers create new issues:

```diff
 {"id":"beads-1","title":"First issue",...}
 {"id":"beads-2","title":"Second issue",...}
+{"id":"beads-3","title":"From branch A",...}
+{"id":"beads-4","title":"From branch B",...}
```

Git may show a conflict, but resolution is simple: **keep both lines** (both changes are compatible).

**With hash-based IDs (v0.20.1+), same-ID scenarios are updates, not collisions:**

If you import an issue with the same ID but different fields, beads treats it as an update to the existing issue. This is normal behavior - hash IDs remain stable, so same ID = same issue being updated.

For git conflicts where the same issue was modified on both branches, manually resolve the JSONL conflict (usually keeping the newer `updated_at` timestamp), then `beads import` will apply the update.

## Migration Questions

### How do I migrate from GitHub Issues / Jira / Linear?

We don't have automated migration tools yet, but you can:

1. Export issues from your current tracker (usually CSV or JSON)
2. Write a simple script to convert to beads's JSONL format
3. Import with `beads import -i issues.jsonl`

See [examples/](../examples) for scripting patterns. Contributions welcome!

### Can I export back to GitHub Issues / Jira?

Not yet built-in, but you can:

1. Export from beads: `beads export -o issues.jsonl --json`
2. Write a script to convert JSONL to your target format
3. Use the target system's API to import

The [config.md](./config.md) guide shows how to store integration settings. Contributions for standard exporters welcome!

## Performance Questions

### How does beads handle scale?

beads uses SQLite, which handles millions of rows efficiently. For a typical project with thousands of issues:

- Commands complete in <100ms
- Full-text search is instant
- Dependency graphs traverse quickly
- JSONL files stay small (one line per issue)

For extremely large projects (100k+ issues), you might want to filter exports or use multiple databases per component.

### What if my JSONL file gets too large?

Use compaction to remove old closed issues:

```sh
# Preview what would be compacted
beads compact --dry-run --all

# Compact issues closed more than 90 days ago
beads compact --days 90
```

Or split your project into multiple databases:

```sh
cd ~/project/frontend && beads init --prefix fe
cd ~/project/backend && beads init --prefix be
```

## Use Case Questions

### Can I use beads for non-code projects?

Sure! beads is just an issue tracker. Use it for:

- Writing projects (chapters as issues, dependencies as outlines)
- Research projects (papers, experiments, dependencies)
- Home projects (renovations with blocking tasks)
- Any workflow with dependencies

The agent-friendly design works for any AI-assisted workflow.

### Can I use beads with multiple AI agents simultaneously?

Yes! Each agent can:

1. Query ready work: `beads ready --assignee agent-name`
2. Claim issues: `beads update <id> --status in_progress --assignee agent-name`
3. Create discovered work: `beads create "Found issue" --deps discovered-from:<parent-id>`
4. Sync via git commits

beads's git-based sync means agents work independently and merge their changes like developers do.

### Does beads work offline?

Yes! beads is designed for offline-first operation:

- All queries run against local SQLite database
- No network required for any commands
- Sync happens via git push/pull when you're online
- Full functionality available without internet

This makes beads ideal for:

- Working on planes/trains
- Unstable network connections
- Air-gapped environments
- Privacy-sensitive projects

## Technical Questions

### What dependencies does beads have?

beads is a single static binary with no runtime dependencies:

- **Language**: Go 1.24+
- **Database**: SQLite (embedded, pure Go driver)
- **Optional**: Git (for sync across machines)

That's it! No PostgreSQL, no Redis, no Docker, no node_modules.

### Can I extend beads's database?

Yes! See [extending.md](extending.md) for how to:

- Add custom tables to the SQLite database
- Join with issue data
- Build custom queries
- Create integrations

### Does beads support Windows?

Yes! beads has native Windows support (v0.9.0+):

- No MSYS or MinGW required
- PowerShell install script
- Works with Windows paths and filesystem
- Daemon uses TCP instead of Unix sockets

See [installing.md](./installing.md#windows-11) for details.

### Can I use beads with git worktrees?

Yes, but with limitations. The daemon doesn't work correctly with worktrees, so use `--no-daemon` mode:

```sh
export BEADS_NO_DAEMON=1
beads ready
beads create "Fix bug" -p 1
```

See [advanced.md#git-worktrees](advanced.md#git-worktrees) for details.

### What's the difference between SQLite corruption and ID collisions?

beads handles two distinct types of integrity issues:

**1. Logical Consistency (Collision Resolution)**

The hash/fingerprint/collision architecture prevents:

- **ID collisions**: Same ID assigned to different issues (e.g., from parallel workers or branch merges)
- **Wrong prefix bugs**: Issues created with incorrect prefix due to config mismatch
- **Merge conflicts**: Branch divergence creating conflicting JSONL content

**Solution**: `beads import --resolve-collisions` automatically remaps colliding IDs and updates all references.

**2. Physical SQLite Corruption**

SQLite database file corruption can occur from:

- **Disk/hardware failures**: Power loss, disk errors, filesystem corruption
- **Concurrent writes**: Multiple processes writing to the same database file simultaneously
- **Container scenarios**: Shared database volumes with multiple containers

**Solution**: Reimport from JSONL (which survives in git history):

```sh
mv .beads/*.db .beads/*.db.backup
beads init
beads import -i .beads/issues.jsonl
```

**Key Difference**: Collision resolution fixes logical issues in the data. Physical corruption requires restoring from the JSONL source of truth.

**When to use in-memory mode (`--no-db`)**: For multi-process/container scenarios where SQLite's file locking isn't sufficient. The in-memory backend loads from JSONL at startup and writes back after each command, avoiding shared database state entirely.

## Getting Help

### Where can I get more help?

- **Documentation**: [README.md](../README.md), [quickstart.md](./quickstart.md), [advanced.md](advanced.md)
- **Troubleshooting**: [troubleshooting.md](./troubleshooting.md)
- **Examples**: [examples/](../examples)
- **GitHub Issues**: [Report bugs or request features](https://github.com/shaneholloman/beads/issues)
- **GitHub Discussions**: [Ask questions](https://github.com/shaneholloman/beads/discussions)

### How can I contribute?

Contributions are welcome! See [CONTRIBUTING.md](../CONTRIBUTING.md) for:

- Code contribution guidelines
- How to run tests
- Development workflow
- Issue and PR templates

### Where's the roadmap?

The roadmap lives in beads itself! Run:

```sh
beads list --priority 0 --priority 1 --json
```

Or check the GitHub Issues for feature requests and planned improvements.
