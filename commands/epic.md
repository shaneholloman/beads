---
description: Epic management commands
argument-hint: [command]
---

Manage epics (large features composed of multiple issues).

## Available Commands

- **status**: Show epic completion status
  - Shows progress for each epic
  - Lists child issues and their states
  - Calculates completion percentage

- **close-eligible**: Close epics where all children are complete
  - Automatically closes epics when all child issues are done
  - Useful for bulk epic cleanup

## Epic Workflow

1. Create epic: `beads create "Large Feature" -t epic -p 1`
2. Link subtasks: `beads dep add beads-10 beads-20 --type parent-child` (epic beads-10 is parent of task beads-20)
3. Track progress: `beads epic status`
4. Auto-close when done: `beads epic close-eligible`

Epics use parent-child dependencies to track subtasks.
