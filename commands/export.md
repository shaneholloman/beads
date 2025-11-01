---
description: Export issues to JSONL format
argument-hint: [-o output-file]
---

# Export Issues

> Export all issues to JSON Lines format (one JSON object per line).

## Usage

- **To stdout**: `beads export`
- **To file**: `beads export -o issues.jsonl`
- **Filter by status**: `beads export --status open`

Issues are sorted by ID for consistent diffs, making git diffs readable.

## Automatic Export

The daemon automatically exports to `.beads/issues.jsonl` after any CRUD operation (5-second debounce). Manual export is rarely needed unless you need a custom output location or filtered export.

Export is used for:

- Git version control
- Backup
- Sharing issues between repositories
- Data migration
