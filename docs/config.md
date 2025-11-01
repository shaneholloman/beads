# Configuration System

beads has two complementary configuration systems:

1. **Tool-level configuration** (Viper): User preferences for tool behavior (flags, output format)
2. **Project-level configuration** (`beads config`): Integration data and project-specific settings

## Tool-Level Configuration (Viper)

### Overview

Tool preferences control how `beads` behaves globally or per-user. These are stored in config files or environment variables and managed by [Viper](https://github.com/spf13/viper).

**Configuration precedence** (highest to lowest):

1. Command-line flags (`--json`, `--no-daemon`, etc.)
2. Environment variables (`BEADS_JSON`, `BEADS_NO_DAEMON`, etc.)
3. Config file (`~/.config/beads/config.yaml` or `.beads/config.yaml`)
4. Defaults

### Config File Locations

Viper searches for `config.yaml` in these locations (in order):

1. `.beads/config.yaml` - Project-specific tool settings (version-controlled)
2. `~/.config/beads/config.yaml` - User-specific tool settings
3. `~/.beads/config.yaml` - Legacy user settings

### Supported Settings

Tool-level settings you can configure:

| Setting | Flag | Environment Variable | Default | Description |
|---------|------|---------------------|---------|-------------|
| `json` | `--json` | `BEADS_JSON` | `false` | Output in JSON format |
| `no-daemon` | `--no-daemon` | `BEADS_NO_DAEMON` | `false` | Force direct mode, bypass daemon |
| `no-auto-flush` | `--no-auto-flush` | `BEADS_NO_AUTO_FLUSH` | `false` | Disable auto JSONL export |
| `no-auto-import` | `--no-auto-import` | `BEADS_NO_AUTO_IMPORT` | `false` | Disable auto JSONL import |
| `db` | `--db` | `BEADS_DB` | (auto-discover) | Database path |
| `actor` | `--actor` | `BEADS_ACTOR` | `$USER` | Actor name for audit trail |
| `flush-debounce` | - | `BEADS_FLUSH_DEBOUNCE` | `5s` | Debounce time for auto-flush |
| `auto-start-daemon` | - | `BEADS_AUTO_START_DAEMON` | `true` | Auto-start daemon if not running |

### Example Config File

`~/.config/beads/config.yaml`:

```yaml
# Default to JSON output for scripting
json: true

# Disable daemon for single-user workflows
no-daemon: true

# Custom debounce for auto-flush (default 5s)
flush-debounce: 10s

# Auto-start daemon (default true)
auto-start-daemon: true
```

`.beads/config.yaml` (project-specific):

```yaml
# Project team prefers longer flush delay
flush-debounce: 15s
```

### Why Two Systems?

**Tool settings (Viper)** are user preferences:

- How should I see output? (`--json`)
- Should I use the daemon? (`--no-daemon`)
- How should the CLI behave?

**Project config (`beads config`)** is project data:

- What's our Jira URL?
- What are our Linear tokens?
- How do we map statuses?

This separation is correct: **tool settings are user-specific, project config is team-shared**.

Agents benefit from `beads config`'s structured CLI interface over manual YAML editing.

## Project-Level Configuration (`beads config`)

### Overview

Project configuration is:

- **Per-project**: Isolated to each `.beads/*.db` database
- **Version-control-friendly**: Stored in SQLite, queryable and scriptable
- **Machine-readable**: JSON output for automation
- **Namespace-based**: Organized by integration or purpose

## Commands

### Set Configuration

```sh
beads config set <key> <value>
beads config set --json <key> <value>  # JSON output
```

Examples:

```sh
beads config set jira.url "https://company.atlassian.net"
beads config set jira.project "PROJ"
beads config set jira.status_map.todo "open"
```

### Get Configuration

```sh
beads config get <key>
beads config get --json <key>  # JSON output
```

Examples:

```sh
beads config get jira.url
# Output: https://company.atlassian.net

beads config get --json jira.url
# Output: {"key":"jira.url","value":"https://company.atlassian.net"}
```

### List All Configuration

```sh
beads config list
beads config list --json  # JSON output
```

Example output:

```txt
Configuration:
  compact_tier1_days = 90
  compact_tier1_dep_levels = 2
  jira.project = PROJ
  jira.url = https://company.atlassian.net
```

JSON output:

```json
{
  "compact_tier1_days": "90",
  "compact_tier1_dep_levels": "2",
  "jira.project": "PROJ",
  "jira.url": "https://company.atlassian.net"
}
```

### Unset Configuration

```sh
beads config unset <key>
beads config unset --json <key>  # JSON output
```

Example:

```sh
beads config unset jira.url
```

## Namespace Convention

Configuration keys use dot-notation namespaces to organize settings:

### Core Namespaces

- `compact_*` - Compaction settings (see extending.md)
- `issue_prefix` - Issue ID prefix (managed by `beads init`)
- `max_collision_prob` - Maximum collision probability for adaptive hash IDs (default: 0.25)
- `min_hash_length` - Minimum hash ID length (default: 4)
- `max_hash_length` - Maximum hash ID length (default: 8)

### Integration Namespaces

Use these namespaces for external integrations:

- `jira.*` - Jira integration settings
- `linear.*` - Linear integration settings
- `github.*` - GitHub integration settings
- `custom.*` - Custom integration settings

### Example: Adaptive Hash ID Configuration

```sh
# Configure adaptive ID lengths (see docs/adaptive-ids.md)
# Default: 25% max collision probability
beads config set max_collision_prob "0.25"

# Start with 4-char IDs, scale up as database grows
beads config set min_hash_length "4"
beads config set max_hash_length "8"

# Stricter collision tolerance (1%)
beads config set max_collision_prob "0.01"

# Force minimum 5-char IDs for consistency
beads config set min_hash_length "5"
```

See [docs/adaptive-ids.md](adaptive-ids.md) for detailed documentation.

### Example: Jira Integration

```sh
# Configure Jira connection
beads config set jira.url "https://company.atlassian.net"
beads config set jira.project "PROJ"
beads config set jira.api_token "YOUR_TOKEN"

# Map beads statuses to Jira statuses
beads config set jira.status_map.open "To Do"
beads config set jira.status_map.in_progress "In Progress"
beads config set jira.status_map.closed "Done"

# Map beads issue types to Jira issue types
beads config set jira.type_map.bug "Bug"
beads config set jira.type_map.feature "Story"
beads config set jira.type_map.task "Task"
```

### Example: Linear Integration

```sh
# Configure Linear connection
beads config set linear.api_token "YOUR_TOKEN"
beads config set linear.team_id "team-123"

# Map statuses
beads config set linear.status_map.open "Backlog"
beads config set linear.status_map.in_progress "In Progress"
beads config set linear.status_map.closed "Done"
```

### Example: GitHub Integration

```sh
# Configure GitHub connection
beads config set github.org "myorg"
beads config set github.repo "myrepo"
beads config set github.token "YOUR_TOKEN"

# Map beads labels to GitHub labels
beads config set github.label_map.bug "bug"
beads config set github.label_map.feature "enhancement"
```

## Use in Scripts

Configuration is designed for scripting. Use `--json` for machine-readable output:

```sh
#!/bin/bash

# Get Jira URL
JIRA_URL=$(beads config get --json jira.url | jq -r '.value')

# Get all config and extract multiple values
beads config list --json | jq -r '.["jira.project"]'
```

Example Python script:

```python
import json
import subprocess

def get_config(key):
    result = subprocess.run(
        ["beads", "config", "get", "--json", key],
        capture_output=True,
        text=True
    )
    data = json.loads(result.stdout)
    return data["value"]

def list_config():
    result = subprocess.run(
        ["beads", "config", "list", "--json"],
        capture_output=True,
        text=True
    )
    return json.loads(result.stdout)

# Use in integration
jira_url = get_config("jira.url")
jira_project = get_config("jira.project")
```

## Best Practices

1. **Use namespaces**: Prefix keys with integration name (e.g., `jira.*`, `linear.*`)
2. **Hierarchical keys**: Use dots for structure (e.g., `jira.status_map.open`)
3. **Document your keys**: Add comments in integration scripts
4. **Security**: Store tokens in config, but add `.beads/*.db` to `.gitignore` (beads does this automatically)
5. **Per-project**: Configuration is project-specific, so each repo can have different settings

## Integration with beads Commands

Some beads commands automatically use configuration:

- `beads compact` uses `compact_tier1_days`, `compact_tier1_dep_levels`, etc.
- `beads init` sets `issue_prefix`

External integration scripts can read configuration to sync with Jira, Linear, GitHub, etc.

## See Also

- [README.md](../README.md) - Main documentation
- [extending.md](extending.md) - Database schema and compaction config
- [examples/adapters/](../examples/adapters) - Integration examples
