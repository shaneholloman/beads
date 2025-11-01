---
description: List issues with optional filters
argument-hint: [--status] [--priority] [--type] [--assignee] [--label]
---

List beads issues with optional filtering.

## Filters

- **--status, -s**: Filter by status (open, in_progress, blocked, closed)
- **--priority, -p**: Filter by priority (0-4: 0=critical, 1=high, 2=medium, 3=low, 4=backlog)
- **--type, -t**: Filter by type (bug, feature, task, epic, chore)
- **--assignee, -a**: Filter by assignee
- **--label, -l**: Filter by labels (comma-separated, must have ALL labels)
- **--title**: Filter by title text (case-insensitive substring match)
- **--limit, -n**: Limit number of results

## Examples

- `beads list --status open --priority 1`: High priority open issues
- `beads list --type bug --assignee alice`: Alice's assigned bugs
- `beads list --label backend,needs-review`: Backend issues needing review
- `beads list --title "auth"`: Issues with "auth" in the title

## Output Formats

- Default: Human-readable table
- `--json`: JSON format for scripting
- `--format digraph`: Graph format for golang.org/x/tools/cmd/digraph
- `--format dot`: Graphviz DOT format
