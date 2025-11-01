# Beads Quickstart

Get up and running with Beads in 2 minutes.

## Installation

```bash
cd ~/src/beads
go build -o beads ./cmd/beads
./beads --help
```

## Your First Issues

```bash
# Create a few issues
./beads create "Set up database" -p 1 -t task
./beads create "Create API" -p 2 -t feature
./beads create "Add authentication" -p 2 -t feature

# List them
./beads list
```

**Note:** Issue IDs are hash-based (e.g., `beads-a1b2`, `beads-f14c`) to prevent collisions when multiple agents/branches work concurrently.

## Hierarchical Issues (Epics)

For large features, use hierarchical IDs to organize work:

```bash
# Create epic (generates parent hash ID)
./beads create "Auth System" -t epic -p 1
# Returns: beads-a3f8e9

# Create child tasks (automatically get .1, .2, .3 suffixes)
./beads create "Design login UI" -p 1       # beads-a3f8e9.1
./beads create "Backend validation" -p 1    # beads-a3f8e9.2
./beads create "Integration tests" -p 1     # beads-a3f8e9.3

# View hierarchy
./beads dep tree beads-a3f8e9
```

Output:

```
Dependency tree for beads-a3f8e9:

→ beads-a3f8e9: Auth System [epic] [P1] (open)
  → beads-a3f8e9.1: Design login UI [P1] (open)
  → beads-a3f8e9.2: Backend validation [P1] (open)
  → beads-a3f8e9.3: Integration tests [P1] (open)
```

## Add Dependencies

```bash
# API depends on database
./beads dep add beads-2 beads-1

# Auth depends on API
./beads dep add beads-3 beads-2

# View the tree
./beads dep tree beads-3
```

Output:

```
Dependency tree for beads-3:

→ beads-3: Add authentication [P2] (open)
  → beads-2: Create API [P2] (open)
    → beads-1: Set up database [P1] (open)
```

## Find Ready Work

```bash
./beads ready
```

Output:

```
Ready work (1 issues with no blockers):

1. [P1] beads-1: Set up database
```

Only beads-1 is ready because beads-2 and beads-3 are blocked!

## Work the Queue

```bash
# Start working on beads-1
./beads update beads-1 --status in_progress

# Complete it
./beads close beads-1 --reason "Database setup complete"

# Check ready work again
./beads ready
```

Now beads-2 is ready!

## Track Progress

```bash
# See blocked issues
./beads blocked

# View statistics
./beads stats
```

## Database Location

By default: `~/.beads/default.db`

You can use project-specific databases:

```bash
./beads --db ./my-project.db create "Task"
```

## Migrating Databases

After upgrading beads, use `beads migrate` to check for and migrate old database files:

```bash
# Check for migration opportunities
./beads migrate --dry-run

# Migrate old databases to beads.db
./beads migrate

# Migrate and clean up old files
./beads migrate --cleanup --yes
```

## Next Steps

- Add labels: `./beads create "Task" -l "backend,urgent"`
- Filter ready work: `./beads ready --priority 1`
- Search issues: `./beads list --status open`
- Detect cycles: `./beads dep cycles`

See [README.md](../README.md) for full documentation.
