# `beads`

> Issue Tracker for agents

## What is it?

Dependency-aware issue tracker for AI coding agents. Issues chain together like beads through four dependency types (blocks, related, parent-child, discovered-from). Local SQLite database syncs via git through JSONL export/import, acting like a distributed database without any server infrastructure.

**WARNING: Alpha Status** - Core features work but expect API changes before 1.0.

## Features

- **Dependency tracking** - Chain issues together, detect cycles, visualize trees
- **Ready work detection** - Find unblocked issues automatically
- **Git-synced** - JSONL stored in git, auto-sync between database and files
- **Hash-based IDs** - Collision-resistant for multi-agent/multi-branch workflows
- **Agent-optimized** - JSON output, MCP server, CLI designed for programmatic use
- **Multi-project** - Auto-discovers databases by directory, complete isolation
- **Extensible** - Add your own SQLite tables alongside core schema

## Installation

**Quick install:**

```sh
curl -fsSL https://raw.githubusercontent.com/shaneholloman/beads/main/scripts/install.sh | bash
```

**Homebrew:**

```sh
brew tap shaneholloman/beads
brew install beads
```

**From source:**

```sh
git clone https://github.com/shaneholloman/beads.git
cd beads
go build -o beads ./cmd/beads
sudo cp beads /usr/local/bin/beads  # Or: cp beads ~/bin/beads
```

**Other platforms:** See [installing.md](docs/installing.md) for Windows, Arch Linux, IDE integration.

## Quick Start

**Initialize in your project:**

```sh
cd your-project
beads init
```

**Tell your AI agent:**

```sh
echo "Use 'beads' for issue tracking. Run 'beads onboard' for instructions." AGENTS.md
```

Your agent handles the rest - creating issues, tracking dependencies, finding ready work.

**Manual usage:**

```sh
beads ready                           # Show unblocked work
beads create "Fix bug" -p 1 -t bug    # Create issue
beads show beads-a1b2                    # View details
beads dep tree beads-a1b2                # Visualize dependencies
beads close beads-a1b2 --reason "Done"   # Mark complete
```

## Core Concepts

**Dependencies:**

- `blocks` - Hard blocker (affects ready work)
- `related` - Soft connection
- `parent-child` - Epic/subtask hierarchy
- `discovered-from` - Work found during execution

**Priorities:** 0 (critical) to 4 (backlog)

**Types:** bug, feature, task, epic, chore

**Hash IDs:** Collision-resistant identifiers (beads-a1b2, beads-f14c) instead of sequential numbers. See [hash-id-design.md](docs/hash-id-design.md) for details.

## Git Workflow

beads auto-syncs with git:

```sh
beads create "Fix bug" -p 1
# After 5 seconds: exports to .beads/issues.jsonl

git add .beads/issues.jsonl
git commit -m "Working on fix"
git push

# On other machine:
git pull
beads ready  # Auto-imports updated JSONL
```

Install git hooks for instant sync:

```sh
cd examples/git-hooks && ./install.sh
```

See [advanced.md](docs/advanced.md) for merge conflict handling and daemon configuration.

## Development

**Build and test:**

```sh
# Build local binary
go build -o beads ./cmd/beads

# Test locally
./beads version
./beads ready

# Run tests
go test ./...

# Install system-wide
cp ./beads ~/bin/beads
# or
sudo cp ./beads /usr/local/bin/beads
```

**Version management:**

```sh
./scripts/bump-version.sh 0.24.0 --commit  # Bump all version files
git push origin main
```

**Before committing:**

```sh
go test ./...                    # All tests pass
go fmt ./...                     # Format code
golangci-lint run ./...          # Check linting (see docs/linting.md for baseline)
```

See [scripts/README.md](scripts/README.md) for release process.

## Documentation

- **[installing.md](docs/installing.md)** - Complete installation guide
- **[quickstart.md](docs/quickstart.md)** - Interactive tutorial
- **[advanced.md](docs/advanced.md)** - Daemon config, merge strategies, prefix renaming
- **[hash-id-design.md](docs/hash-id-design.md)** - Hash ID system and collision math
- **[labels.md](docs/labels.md)** - Label system guide
- **[extending.md](docs/extending.md)** - Database extension patterns
- **[faq.md](docs/faq.md)** - Frequently asked questions
- **[troubleshooting.md](docs/troubleshooting.md)** - Common issues

## License

MIT
