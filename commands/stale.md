---
description: Show orphaned claims and dead executors
argument-hint: [--release] [--threshold]
---

# Stale Issues

> Show issues stuck in_progress with execution_state where the executor is dead or stopped.

Helps identify orphaned work that needs manual recovery.

## Stale Detection

An issue is stale if:

- It has an execution_state (claimed by an executor)
- AND the executor status is 'stopped'
- OR the executor's last_heartbeat is older than threshold

Default threshold: 300 seconds (5 minutes)

## Usage

- **List stale issues**: `beads stale`
- **Custom threshold**: `beads stale --threshold 600` (10 minutes)
- **Auto-release**: `beads stale --release` (automatically release all stale issues)

Useful for parallel execution systems where workers may crash or get stopped.
