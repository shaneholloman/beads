# GitHub Issues to beads Importer

Import issues from GitHub repositories into `beads`.

## Overview

This tool converts GitHub Issues to beads's JSONL format, supporting both:

1. **GitHub API** - Fetch issues directly from a repository
2. **JSON Export** - Parse manually exported GitHub issues

## Features

- ✔ **Fetch from GitHub API** - Direct import from any public/private repo
- ✔ **JSON file import** - Parse exported GitHub issues JSON
- ✔ **Label mapping** - Auto-map GitHub labels to beads priority/type
- ✔ **Preserve metadata** - Keep assignees, timestamps, descriptions
- ✔ **Cross-references** - Convert `#123` references to dependencies
- ✔ **External links** - Preserve URLs back to original GitHub issues
- ✔ **Filter PRs** - Automatically excludes pull requests

## Installation

No dependencies required! Uses Python 3 standard library.

For API access, set up a GitHub token:

```bash
# Create token at: https://github.com/settings/tokens
# Permissions needed: public_repo (or repo for private repos)

export GITHUB_TOKEN=ghp_your_token_here
```

**Security Note:** Use the `GITHUB_TOKEN` environment variable instead of `--token` flag when possible. The `--token` flag may appear in shell history and process listings.

## Usage

### From GitHub API

```bash
# Fetch all issues from a repository
uv run gh2jsonl.py --repo owner/repo | beads import

# Save to file first (recommended)
uv run gh2jsonl.py --repo owner/repo > issues.jsonl
beads import -i issues.jsonl --dry-run  # Preview
beads import -i issues.jsonl             # Import

# Fetch only open issues
uv run gh2jsonl.py --repo owner/repo --state open

# Fetch only closed issues
uv run gh2jsonl.py --repo owner/repo --state closed
```

### From JSON File

Export issues from GitHub (via API or manually), then:

```bash
# Single issue
curl -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/repos/owner/repo/issues/123 > issue.json

uv run gh2jsonl.py --file issue.json | beads import

# Multiple issues
curl -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/repos/owner/repo/issues > issues.json

uv run gh2jsonl.py --file issues.json | beads import
```

### Custom Options

```bash
# Use custom prefix (instead of 'beads')
uv run gh2jsonl.py --repo owner/repo --prefix myproject

# Start numbering from specific ID
uv run gh2jsonl.py --repo owner/repo --start-id 100

# Pass token directly (instead of env var)
uv run gh2jsonl.py --repo owner/repo --token ghp_...
```

## Label Mapping

The script maps GitHub labels to beads fields:

### Priority Mapping

| GitHub Labels              | beads Priority |
| -------------------------- | -------------- |
| `critical`, `p0`, `urgent` | 0 (Critical)   |
| `high`, `p1`, `important`  | 1 (High)       |
| (default)                  | 2 (Medium)     |
| `low`, `p3`, `minor`       | 3 (Low)        |
| `backlog`, `p4`, `someday` | 4 (Backlog)    |

### Type Mapping

| GitHub Labels                          | beads Type |
| -------------------------------------- | ---------- |
| `bug`, `defect`                        | bug        |
| `feature`, `enhancement`               | feature    |
| `epic`, `milestone`                    | epic       |
| `chore`, `maintenance`, `dependencies` | chore      |
| (default)                              | task       |

### Status Mapping

| GitHub State | GitHub Labels                       | beads Status |
| ------------ | ----------------------------------- | ------------ |
| closed       | (any)                               | closed       |
| open         | `in progress`, `in-progress`, `wip` | in_progress  |
| open         | `blocked`                           | blocked      |
| open         | (default)                           | open         |

### Labels

All other labels are preserved in the `labels` field. Labels used for mapping (priority, type, status) are filtered out to avoid duplication.

## Field Mapping

| GitHub Field     | beads Field                        | Notes                  |
| ---------------- | ---------------------------------- | ---------------------- |
| `number`         | (internal mapping)                 | GH#123 → beads-1, etc. |
| `title`          | `title`                            | Direct copy            |
| `body`           | `description`                      | Direct copy            |
| `state`          | `status`                           | See status mapping     |
| `labels`         | `priority`, `issue_type`, `labels` | See label mapping      |
| `assignee.login` | `assignee`                         | First assignee only    |
| `created_at`     | `created_at`                       | ISO 8601 timestamp     |
| `updated_at`     | `updated_at`                       | ISO 8601 timestamp     |
| `closed_at`      | `closed_at`                        | ISO 8601 timestamp     |
| `html_url`       | `external_ref`                     | Link back to GitHub    |

## Cross-References

Issue references in the body text are converted to dependencies:

**GitHub:**

```markdown
This depends on #123 and fixes #456.
See also owner/other-repo#789.
```

**Result:**

- If GH#123 was imported, creates `related` dependency to its beads ID
- If GH#456 was imported, creates `related` dependency to its beads ID
- Cross-repo references (#789) are ignored (unless those issues were also imported)

**Note:** Dependency records use `"issue_id": ""` format, which the beads importer automatically fills. This matches the behavior of the markdown-to-jsonl converter.

## Examples

### Example 1: Import Active Issues

```bash
# Import only open issues for active work
export GITHUB_TOKEN=ghp_...
uv run gh2jsonl.py --repo mycompany/myapp --state open > open-issues.jsonl

# Preview
cat open-issues.jsonl | jq .

# Import
beads import -i open-issues.jsonl
beads ready  # See what's ready to work on
```

### Example 2: Full Repository Migration

```bash
# Import all issues (open and closed)
uv run gh2jsonl.py --repo mycompany/myapp > all-issues.jsonl

# Preview import (check for new issues and updates)
beads import -i all-issues.jsonl --dry-run

# Import issues
beads import -i all-issues.jsonl

# View stats
beads stats
```

### Example 3: Partial Import from JSON

```bash
# Manually export specific issues via GitHub API
gh api repos/owner/repo/issues?labels=p1,bug > high-priority-bugs.json

# Import
uv run gh2jsonl.py --file high-priority-bugs.json | beads import
```

## Customization

The script is intentionally simple to customize for your workflow:

### 1. Adjust Label Mappings

Edit `map_priority()`, `map_issue_type()`, and `map_status()` to match your label conventions:

```python
def map_priority(self, labels: List[str]) -> int:
    label_names = [label.get("name", "").lower() if isinstance(label, dict) else label.lower() for label in labels]

    # Add your custom mappings
    if any(l in label_names for l in ["sev1", "emergency"]):
        return 0
    # ... etc
```

### 2. Add Custom Fields

Map additional GitHub fields to beads:

```python
def convert_issue(self, gh_issue: Dict[str, Any]) -> Dict[str, Any]:
    # ... existing code ...

    # Add milestone to design field
    if gh_issue.get("milestone"):
        issue["design"] = f"Milestone: {gh_issue['milestone']['title']}"

    return issue
```

### 3. Enhanced Dependency Detection

Parse more dependency patterns from body text:

```python
def extract_dependencies_from_body(self, body: str) -> List[str]:
    # ... existing code ...

    # Add: "Blocks: #123, #456"
    blocks_pattern = r'Blocks:\s*((?:#\d+(?:\s*,\s*)?)+)'
    # ... etc
```

## Limitations

- **Single assignee**: GitHub supports multiple assignees, beads supports one
- **No milestones**: GitHub milestones aren't mapped (consider using design field)
- **Simple cross-refs**: Only basic `#123` patterns detected
- **No comments**: Issue comments aren't imported (only the body)
- **No reactions**: GitHub reactions/emoji aren't imported
- **No projects**: GitHub project board info isn't imported

## API Rate Limits

GitHub API has rate limits:

- **Authenticated**: 5,000 requests/hour
- **Unauthenticated**: 60 requests/hour

This script uses 1 request per 100 issues (pagination), so:

- Can fetch ~500,000 issues/hour (authenticated)
- Can fetch ~6,000 issues/hour (unauthenticated)

For large repositories (>1000 issues), authentication is recommended.

**Note:** The script automatically includes a `User-Agent` header (required by GitHub) and provides actionable error messages when rate limits are exceeded, including the reset timestamp.

## Troubleshooting

### "GitHub token required"

Set the `GITHUB_TOKEN` environment variable:

```bash
export GITHUB_TOKEN=ghp_your_token_here
```

Or pass directly:

```bash
uv run gh2jsonl.py --repo owner/repo --token ghp_...
```

### "GitHub API error: 404"

- Check repository name format: `owner/repo`
- Check repository exists and is accessible
- For private repos, ensure token has `repo` scope

### "GitHub API error: 403"

- Rate limit exceeded (wait or use authentication)
- Token doesn't have required permissions
- Repository requires different permissions

### Issue numbers don't match

This is expected! GitHub issue numbers (e.g., #123) are mapped to beads IDs (e.g., beads-1) based on import order. The original GitHub URL is preserved in `external_ref`.

## See Also

- [beads README](../../README.md) - Main documentation
- [Markdown Import Example](../markdown-to-jsonl/) - Import from markdown
- [TEXT_FORMATS.md](../../TEXT_FORMATS.md) - Understanding beads's JSONL format
- [JSONL Import Guide](../../README.md#import) - Import collision handling
