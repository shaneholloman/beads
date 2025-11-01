# beads daemons - Daemon Management

Manage beads daemon processes across all repositories and worktrees.

## Synopsis

```bash
beads daemons <subcommand> [flags]
```

## Description

The `beads daemons` command provides tools for discovering, monitoring, and managing multiple beads daemon processes across your system. This is useful when working with multiple repositories or git worktrees.

## Subcommands

### list

List all running beads daemons with metadata.

```bash
beads daemons list [--search DIRS] [--json] [--no-cleanup]
```

**Flags:**

- `--search` - Directories to search for daemons (default: home, /tmp, cwd)
- `--json` - Output in JSON format
- `--no-cleanup` - Skip auto-cleanup of stale sockets

**Example:**

```bash
beads daemons list
beads daemons list --search /Users/me/projects --json
```

### health

Check health of all beads daemons and report issues.

```bash
beads daemons health [--search DIRS] [--json]
```

Reports:

- Stale sockets (dead processes)
- Version mismatches between daemon and CLI
- Unresponsive daemons

**Flags:**

- `--search` - Directories to search for daemons
- `--json` - Output in JSON format

**Example:**

```bash
beads daemons health
beads daemons health --json
```

### stop

Stop a specific daemon gracefully.

```bash
beads daemons stop <workspace-path|pid> [--json]
```

**Arguments:**

- `<workspace-path|pid>` - Workspace path or PID of daemon to stop

**Flags:**

- `--json` - Output in JSON format

**Example:**

```bash
beads daemons stop /Users/me/projects/myapp
beads daemons stop 12345
beads daemons stop /Users/me/projects/myapp --json
```

### logs

View logs for a specific daemon.

```bash
beads daemons logs <workspace-path|pid> [-f] [-n LINES] [--json]
```

**Arguments:**

- `<workspace-path|pid>` - Workspace path or PID of daemon

**Flags:**

- `-f, --follow` - Follow log output (like tail -f)
- `-n, --lines INT` - Number of lines to show from end (default: 50)
- `--json` - Output in JSON format

**Example:**

```bash
beads daemons logs /Users/me/projects/myapp
beads daemons logs 12345 -n 100
beads daemons logs /Users/me/projects/myapp -f
beads daemons logs 12345 --json
```

### killall

Stop all running beads daemons.

```bash
beads daemons killall [--search DIRS] [--force] [--json]
```

Uses escalating shutdown strategy:

1. RPC shutdown (2 second timeout)
2. SIGTERM (3 second timeout)
3. SIGKILL (1 second timeout)

**Flags:**

- `--search` - Directories to search for daemons
- `--force` - Use SIGKILL immediately if graceful shutdown fails
- `--json` - Output in JSON format

**Example:**

```bash
beads daemons killall
beads daemons killall --force
beads daemons killall --json
```

## Common Use Cases

### Version Upgrade

After upgrading beads, restart all daemons to use the new version:

```bash
beads daemons health  # Check for version mismatches
beads daemons killall # Stop all old daemons
# Daemons will auto-start with new version on next beads command
```

### Debugging

Check daemon status and view logs:

```bash
beads daemons list
beads daemons health
beads daemons logs /path/to/workspace -n 100
```

### Cleanup

Remove stale daemon sockets:

```bash
beads daemons list  # Auto-cleanup happens by default
beads daemons list --no-cleanup  # Skip cleanup
```

### Multi-Workspace Management

Discover daemons in specific directories:

```bash
beads daemons list --search /Users/me/projects
beads daemons health --search /Users/me/work
```

## Troubleshooting

### Stale Sockets

If you see stale sockets (dead process but socket file exists):

```bash
beads daemons list  # Auto-cleanup removes stale sockets
```

### Version Mismatch

If daemon version != CLI version:

```bash
beads daemons health  # Identify mismatched daemons
beads daemons killall # Stop all daemons
# Next beads command will auto-start new version
```

### Daemon Won't Stop

If graceful shutdown fails:

```bash
beads daemons killall --force  # Force kill with SIGKILL
```

### Can't Find Daemon

If daemon isn't discovered:

```bash
beads daemons list --search /path/to/workspace
```

Or check the socket manually:

```bash
ls -la /path/to/workspace/.beads/beads.sock
```

## See Also

- [beads daemon](daemon.md) - Start a daemon manually
- [AGENTS.md](../AGENTS.md) - Agent workflow guide
- [README.md](../README.md) - Main documentation
