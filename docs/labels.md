# Labels in Beads

Labels provide flexible, multi-dimensional categorization for issues beyond the structured fields (status, priority, type). Use labels for cross-cutting concerns, technical metadata, and contextual tagging without schema changes.

## Design Philosophy

**When to use labels vs. structured fields:**

- **Structured fields** (status, priority, type) → Core workflow state
  - Status: Where the issue is in the workflow (`open`, `in_progress`, `blocked`, `closed`)
  - Priority: How urgent (0-4)
  - Type: What kind of work (`bug`, `feature`, `task`, `epic`, `chore`)

- **Labels** → Everything else
  - Technical metadata (`backend`, `frontend`, `api`, `database`)
  - Domain/scope (`auth`, `payments`, `search`, `analytics`)
  - Effort estimates (`small`, `medium`, `large`)
  - Quality gates (`needs-review`, `needs-tests`, `breaking-change`)
  - Team/ownership (`team-infra`, `team-product`)
  - Release tracking (`v1.0`, `v2.0`, `backport-candidate`)

## Quick Start

```bash
# Add labels when creating issues
beads create "Fix auth bug" -t bug -p 1 -l auth,backend,urgent

# Add labels to existing issues
beads label add beads-42 security
beads label add beads-42 breaking-change

# List issue labels
beads label list beads-42

# Remove a label
beads label remove beads-42 urgent

# List all labels in use
beads label list-all

# Filter by labels (AND - must have ALL)
beads list --label backend,auth

# Filter by labels (OR - must have AT LEAST ONE)
beads list --label-any frontend,backend

# Combine filters
beads list --status open --priority 1 --label security
```

## Common Label Patterns

### 1. Technical Component Labels

Identify which part of the system:

```bash
backend
frontend
api
database
infrastructure
cli
ui
mobile
```

**Example:**

```bash
beads create "Add GraphQL endpoint" -t feature -p 2 -l backend,api
beads create "Update login form" -t task -p 2 -l frontend,auth,ui
```

### 2. Domain/Feature Area

Group by business domain:

```bash
auth
payments
search
analytics
billing
notifications
reporting
admin
```

**Example:**

```bash
beads list --label payments --status open  # All open payment issues
beads list --label-any auth,security       # Security-related work
```

### 3. Size/Effort Estimates

Quick effort indicators:

```bash
small     # < 1 day
medium    # 1-3 days
large     # > 3 days
```

**Example:**

```bash
# Find small quick wins
beads ready --json | jq '.[] | select(.labels[] == "small")'
```

### 4. Quality Gates

Track what's needed before closing:

```bash
needs-review
needs-tests
needs-docs
breaking-change
```

**Example:**

```bash
beads label add beads-42 needs-review
beads list --label needs-review --status in_progress
```

### 5. Release Management

Track release targeting:

```bash
v1.0
v2.0
backport-candidate
release-blocker
```

**Example:**

```bash
beads list --label v1.0 --status open    # What's left for v1.0?
beads label add beads-42 release-blocker
```

### 6. Team/Ownership

Indicate ownership or interest:

```bash
team-infra
team-product
team-mobile
needs-triage
help-wanted
```

**Example:**

```bash
beads list --assignee alice --label team-infra
beads create "Memory leak in cache" -t bug -p 1 -l team-infra,help-wanted
```

### 7. Special Markers

Process or workflow flags:

```bash
auto-generated     # Created by automation
discovered-from    # Found during other work (also a dep type)
technical-debt
good-first-issue
duplicate
wontfix
```

**Example:**

```bash
beads create "TODO: Refactor parser" -t chore -p 3 -l technical-debt,auto-generated
```

## Filtering by Labels

### AND Filtering (--label)

All specified labels must be present:

```bash
# Issues that are BOTH backend AND urgent
beads list --label backend,urgent

# Open bugs that need review AND tests
beads list --status open --type bug --label needs-review,needs-tests
```

### OR Filtering (--label-any)

At least one specified label must be present:

```bash
# Issues in frontend OR backend
beads list --label-any frontend,backend

# Security or auth related
beads list --label-any security,auth
```

### Combining AND/OR

Mix both filters for complex queries:

```bash
# Backend issues that are EITHER urgent OR a blocker
beads list --label backend --label-any urgent,release-blocker

# Frontend work that needs BOTH review and tests, but in any component
beads list --label needs-review,needs-tests --label-any frontend,ui,mobile
```

## Workflow Examples

### Triage Workflow

```bash
# Create untriaged issue
beads create "Crash on login" -t bug -p 1 -l needs-triage

# During triage, add context
beads label add beads-42 auth
beads label add beads-42 backend
beads label add beads-42 urgent
beads label remove beads-42 needs-triage

# Find untriaged issues
beads list --label needs-triage
```

### Quality Gate Workflow

```bash
# Start work
beads update beads-42 --status in_progress

# Mark quality requirements
beads label add beads-42 needs-tests
beads label add beads-42 needs-docs

# Before closing, verify
beads label list beads-42
# ... write tests and docs ...
beads label remove beads-42 needs-tests
beads label remove beads-42 needs-docs

# Close when gates satisfied
beads close beads-42
```

### Release Planning

```bash
# Tag issues for v1.0
beads label add beads-42 v1.0
beads label add beads-43 v1.0
beads label add beads-44 v1.0

# Track v1.0 progress
beads list --label v1.0 --status closed    # Done
beads list --label v1.0 --status open      # Remaining
beads stats  # Overall progress

# Mark critical items
beads label add beads-45 v1.0
beads label add beads-45 release-blocker
```

### Component-Based Work Distribution

```bash
# Backend team picks up work
beads ready --json | jq '.[] | select(.labels[]? == "backend")'

# Frontend team finds small tasks
beads list --status open --label frontend,small

# Find help-wanted items for new contributors
beads list --label help-wanted,good-first-issue
```

## Label Management

### Listing Labels

```bash
# Labels on a specific issue
beads label list beads-42

# All labels in database with usage counts
beads label list-all

# JSON output for scripting
beads label list-all --json
```

Output:

```json
[
  {"label": "auth", "count": 5},
  {"label": "backend", "count": 12},
  {"label": "frontend", "count": 8}
]
```

### Bulk Operations

Add labels in batch during creation:

```bash
beads create "Issue" -l label1,label2,label3
```

Script to add label to multiple issues:

```bash
# Add "needs-review" to all in_progress issues
beads list --status in_progress --json | jq -r '.[].id' | while read id; do
  beads label add "$id" needs-review
done
```

Remove label from multiple issues:

```bash
# Remove "urgent" from closed issues
beads list --status closed --label urgent --json | jq -r '.[].id' | while read id; do
  beads label remove "$id" urgent
done
```

## Integration with Git Workflow

Labels are automatically synced to `.beads/issues.jsonl` along with all issue data:

```bash
# Make changes
beads create "Fix bug" -l backend,urgent
beads label add beads-42 needs-review

# Auto-exported after 5 seconds (or use git hooks for immediate export)
git add .beads/issues.jsonl
git commit -m "Add backend issue"

# After git pull, labels are auto-imported
git pull
beads list --label backend  # Fresh data including labels
```

## Markdown Import/Export

Labels are preserved when importing from markdown:

```markdown
# Fix Authentication Bug

### Type
bug

### Priority
1

### Labels
auth, backend, urgent, needs-review

### Description
Users can't log in after recent deployment.
```

```bash
beads create -f issue.md
# Creates issue with all four labels
```

## Best Practices

### 1. Establish Conventions Early

Document your team's label taxonomy:

```bash
# Add to project README or CONTRIBUTING.md
- Use lowercase, hyphen-separated (e.g., `good-first-issue`)
- Prefix team labels (e.g., `team-infra`, `team-product`)
- Use consistent size labels (`small`, `medium`, `large`)
```

### 2. Don't Overuse Labels

Labels are flexible, but too many can cause confusion. Prefer:

- 5-10 core technical labels (`backend`, `frontend`, `api`, etc.)
- 3-5 domain labels per project
- Standard process labels (`needs-review`, `needs-tests`)
- Release labels as needed

### 3. Clean Up Unused Labels

Periodically review:

```bash
beads label list-all
# Remove obsolete labels from issues
```

### 4. Use Labels for Filtering, Not Search

Labels are for categorization, not free-text search:

- ✔ Good: `backend`, `auth`, `urgent`
- ✘ Bad: `fix-the-login-bug`, `john-asked-for-this`

### 5. Combine with Dependencies

Labels + dependencies = powerful organization:

```bash
# Epic with labeled subtasks
beads create "Auth system rewrite" -t epic -p 1 -l auth,v2.0
beads create "Implement JWT" -t task -p 1 -l auth,backend --deps parent-child:beads-42
beads create "Update login UI" -t task -p 1 -l auth,frontend --deps parent-child:beads-42

# Find all v2.0 auth work
beads list --label auth,v2.0
```

## AI Agent Usage

Labels are especially useful for AI agents managing complex workflows:

```bash
# Auto-label discovered work
beads create "Found TODO in auth.go" -t task -p 2 -l auto-generated,technical-debt

# Filter for agent review
beads list --label needs-review --status in_progress --json

# Track automation metadata
beads label add beads-42 ai-generated
beads label add beads-42 needs-human-review
```

Example agent workflow:

```bash
# Agent discovers issues during refactor
beads create "Extract validateToken function" -t chore -p 2 \
  -l technical-debt,backend,auth,small \
  --deps discovered-from:beads-10

# Agent marks work for review
beads update beads-42 --status in_progress
# ... agent does work ...
beads label add beads-42 needs-review
beads label add beads-42 ai-generated

# Human reviews and approves
beads label remove beads-42 needs-review
beads label add beads-42 approved
beads close beads-42
```

## Advanced Patterns

### Component Matrix

Track issues across multiple dimensions:

```bash
# Backend + auth + high priority
beads list --label backend,auth --priority 1

# Any frontend work that's small
beads list --label-any frontend,ui --label small

# Critical issues across all components
beads list --priority 0 --label-any backend,frontend,infrastructure
```

### Sprint Planning

```bash
# Label issues for sprint
for id in beads-42 beads-43 beads-44 beads-45; do
  beads label add "$id" sprint-12
done

# Track sprint progress
beads list --label sprint-12 --status closed    # Velocity
beads list --label sprint-12 --status open      # Remaining
beads stats | grep "In Progress"                # Current WIP
```

### Technical Debt Tracking

```bash
# Mark debt
beads create "Refactor legacy parser" -t chore -p 3 -l technical-debt,large

# Find debt to tackle
beads list --label technical-debt --label small
beads list --label technical-debt --priority 1  # High-priority debt
```

### Breaking Change Coordination

```bash
# Identify breaking changes
beads label add beads-42 breaking-change
beads label add beads-42 v2.0

# Find all breaking changes for next major release
beads list --label breaking-change,v2.0

# Ensure they're documented
beads list --label breaking-change --label needs-docs
```

## Troubleshooting

### Labels Not Showing in List

Labels require explicit fetching. The `beads list` command shows issues but not labels in human output (only in JSON).

```bash
# See labels in JSON
beads list --json | jq '.[] | {id, labels}'

# See labels for specific issue
beads show beads-42 --json | jq '.labels'
beads label list beads-42
```

### Label Filtering Not Working

Check label names for exact matches (case-sensitive):

```bash
# These are different labels:
beads label add beads-42 Backend    # Capital B
beads list --label backend       # Won't match

# List all labels to see exact names
beads label list-all
```

### Syncing Labels with Git

Labels are included in `.beads/issues.jsonl` export. If labels seem out of sync:

```bash
# Force export
beads export -o .beads/issues.jsonl

# After pull, force import
beads import -i .beads/issues.jsonl
```

## See Also

- [README.md](../README.md) - Main documentation
- [AGENTS.md](../AGENTS.md) - AI agent integration guide
- [AGENTS.md](../AGENTS.md) - Team workflow patterns
- [advanced.md](advanced.md) - JSONL format details
