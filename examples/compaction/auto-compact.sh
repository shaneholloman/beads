#!/bin/bash
# Smart auto-compaction with thresholds
# Only compacts if there are enough eligible issues
#
# Usage: ./auto-compact.sh [--threshold N] [--tier 1|2]

# Default configuration
THRESHOLD=10  # Minimum eligible issues to trigger compaction
TIER=1
DRY_RUN=false

# Parse arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --threshold)
      THRESHOLD="$2"
      shift 2
      ;;
    --tier)
      TIER="$2"
      shift 2
      ;;
    --dry-run)
      DRY_RUN=true
      shift
      ;;
    *)
      echo "Unknown option: $1"
      echo "Usage: $0 [--threshold N] [--tier 1|2] [--dry-run]"
      exit 1
      ;;
  esac
done

# Check API key
if [ -z "$ANTHROPIC_API_KEY" ]; then
  echo "✘ Error: ANTHROPIC_API_KEY not set"
  exit 1
fi

# Check beads is installed
if ! command -v beads &> /dev/null; then
  echo "✘ Error: beads command not found"
  exit 1
fi

# Check eligible issues
echo "Checking eligible issues (Tier $TIER)..."
ELIGIBLE=$(beads compact --dry-run --all --tier "$TIER" --json 2>/dev/null | jq '. | length' || echo "0")

if [ -z "$ELIGIBLE" ] || [ "$ELIGIBLE" = "null" ]; then
  ELIGIBLE=0
fi

echo "Found $ELIGIBLE eligible issues (threshold: $THRESHOLD)"

if [ "$ELIGIBLE" -lt "$THRESHOLD" ]; then
  echo "⟶ Below threshold, skipping compaction"
  exit 0
fi

if [ "$DRY_RUN" = true ]; then
  echo "Dry run mode - showing candidates:"
  beads compact --dry-run --all --tier "$TIER"
  exit 0
fi

# Run compaction
echo "Compacting $ELIGIBLE issues (Tier $TIER)..."
beads compact --all --tier "$TIER"

# Show stats
echo
echo "Statistics:"
beads compact --stats

echo
echo "✔ Auto-compaction complete"
echo "Remember to commit: git add .beads/issues.jsonl issues.db && git commit -m 'Auto-compact'"
