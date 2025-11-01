---
description: Reopen closed issues
argument-hint: [issue-ids...] [--reason]
---

Reopen one or more closed issues.

Sets status to 'open' and clears the closed_at timestamp. Emits a Reopened event.

## Usage

- **Reopen single**: `beads reopen beads-42`
- **Reopen multiple**: `beads reopen beads-42 beads-43 beads-44`
- **With reason**: `beads reopen beads-42 --reason "Found regression"`

More explicit than `beads update --status open` - specifically designed for reopening workflow.

Common reasons for reopening:

- Regression found
- Requirements changed
- Incomplete implementation
- New information discovered
