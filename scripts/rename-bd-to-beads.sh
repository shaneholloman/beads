#!/bin/bash
set -e

# Script to rename beads → beads throughout the codebase
# Handles edge cases: hashes, function calls, word boundaries
#
# Files EXCLUDED (will NOT be modified):
#   - go.mod, go.sum (contain dependency hashes with 'bd' in them)
#   - .git/ directory (git internals)
#   - .beads/*.db (binary database files)
#   - Binary executables
#
# Uses word boundaries (\b) to avoid corrupting:
#   - Git commit hashes (e.g., 5f936abd7ae8)
#   - Base64 checksums (e.g., klgBDR4)
#   - Middle of other words

echo "Step 1: Renaming command references (bd → beads)..."
echo ""

# Get list of files to process (excluding go.mod, go.sum, .db files)
rg --files -g "*.go" -g "*.md" -g "*.sh" -g "*.yml" -g "*.yaml" -g "*.json" -g "*.toml" -g "*.jsonl" \
  | grep -v "go.mod" \
  | grep -v "go.sum" \
  | grep -v "\.db$" \
  | xargs sed -i '' \
    -e 's/\bbd init\b/beads init/g' \
    -e 's/\bbd ready\b/beads ready/g' \
    -e 's/\bbd create\b/beads create/g' \
    -e 's/\bbd list\b/beads list/g' \
    -e 's/\bbd show\b/beads show/g' \
    -e 's/\bbd update\b/beads update/g' \
    -e 's/\bbd close\b/beads close/g' \
    -e 's/\bbd dep\b/beads dep/g' \
    -e 's/\bbd sync\b/beads sync/g' \
    -e 's/\bbd export\b/beads export/g' \
    -e 's/\bbd import\b/beads import/g' \
    -e 's/\bbd label\b/beads label/g' \
    -e 's/\bbd delete\b/beads delete/g' \
    -e 's/\bbd merge\b/beads merge/g' \
    -e 's/\bbd compact\b/beads compact/g' \
    -e 's/\bbd migrate\b/beads migrate/g' \
    -e 's/\bbd doctor\b/beads doctor/g' \
    -e 's/\bbd info\b/beads info/g' \
    -e 's/\bbd stats\b/beads stats/g' \
    -e 's/\bbd blocked\b/beads blocked/g' \
    -e 's/\bbd version\b/beads version/g' \
    -e 's/\bbd onboard\b/beads onboard/g' \
    -e 's/\bbd quickstart\b/beads quickstart/g' \
    -e 's/\bbd config\b/beads config/g' \
    -e 's/\bbd daemons\b/beads daemons/g' \
    -e 's/\bbd daemon\b/beads daemon/g' \
    -e 's/\bbd duplicates\b/beads duplicates/g' \
    -e 's/\bbd validate\b/beads validate/g' \
    -e 's/\bbd restore\b/beads restore/g' \
    -e 's/\bbd reopen\b/beads reopen/g' \
    -e 's/\bbd edit\b/beads edit/g' \
    -e 's/\bbd epic\b/beads epic/g' \
    -e 's/\bbd detect-pollution\b/beads detect-pollution/g' \
    -e 's/\bbd repair-deps\b/beads repair-deps/g' \
    -e 's/\bbd rename-prefix\b/beads rename-prefix/g' \
    -e 's/\bbd renumber\b/beads renumber/g' \
    -e 's/\bbd stale\b/beads stale/g' \
    -e 's/\bbd comments\b/beads comments/g'

echo "Step 2: Renaming issue ID prefixes (bd- → beads-)..."
echo ""

# Rename issue IDs
rg --files -g "*.go" -g "*.md" -g "*.sh" -g "*.txt" -g "*.jsonl" \
  | grep -v "go.mod" \
  | grep -v "go.sum" \
  | grep -v "\.db$" \
  | xargs sed -i '' \
    -e 's/\bbd-\([a-f0-9]\{4,8\}\)\b/beads-\1/g' \
    -e 's/\bbd-\([0-9]\+\)\b/beads-\1/g'

echo "Step 3: Renaming binary paths..."
echo ""

# Rename binary references (./bd, /bd, etc.)
rg --files -g "*.go" -g "*.md" -g "*.sh" -g "*.yml" -g "*.yaml" -g "*.jsonl" \
  | grep -v "go.mod" \
  | grep -v "go.sum" \
  | grep -v "\.db$" \
  | xargs sed -i '' \
    -e 's|/bd\b|/beads|g' \
    -e 's|\./bd\b|./beads|g' \
    -e 's|"beads"|"beads"|g' \
    -e 's|-o beads |-o beads |g'

echo "Step 4: Renaming directory references..."
echo ""

# cmd/beads → cmd/beads in all references
rg --files -g "*.go" -g "*.md" -g "*.sh" -g "*.yml" -g "*.yaml" -g "*.jsonl" \
  | grep -v "go.mod" \
  | grep -v "go.sum" \
  | grep -v "\.db$" \
  | xargs sed -i '' -e 's|cmd/beads|cmd/beads|g'

echo "Step 5: Renaming example directory..."
echo ""

if [ -d "examples/bd-example-extension-go" ]; then
  mv examples/bd-example-extension-go examples/beads-example-extension-go
  echo "  Renamed examples/bd-example-extension-go → examples/beads-example-extension-go"
fi

# Update references to the old example path
rg --files -g "*.go" -g "*.md" -g "*.jsonl" \
  | grep -v "go.mod" \
  | grep -v "go.sum" \
  | xargs sed -i '' -e 's|bd-example-extension-go|beads-example-extension-go|g'

echo "Step 6: Renaming cmd/beads directory to cmd/beads..."
echo ""

if [ -d "cmd/beads" ]; then
  mv cmd/beads cmd/beads
  echo "  Renamed cmd/beads → cmd/beads"
fi

echo ""
echo "Done! beads → beads rename complete."
echo ""
echo "Verification steps:"
echo "  1. go build -o beads ./cmd/beads"
echo "  2. ./beads version"
echo "  3. go test ./..."
echo "  4. rg '\\bbd\\b' | grep -v 'beads' | head -20  # Check for remaining 'bd' refs"
echo ""
echo "If everything looks good:"
echo "  git add -A"
echo "  git commit -m 'refactor: Rename beads → beads for semantic consistency'"
