# Claude Desktop MCP Server for Beads

> **Note**: The beads MCP server is now fully implemented! See [adapters/mcp](../../adapters/mcp/) for the production implementation.

## What This Provides

An MCP server that exposes beads functionality to Claude Desktop and other MCP clients, allowing Claude to:

- Query ready work
- Create and update issues
- Manage dependencies
- Track discovered work

## Quick Start

Install the beads MCP server:

```sh
uv tool install mcp-beads
```

Add to your Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "beads": {
      "command": "mcp-beads"
    }
  }
}
```

Restart Claude Desktop and you're done! Claude can now manage your beads issues.

## Full Documentation

See the [mcp-beads README](../../adapters/mcp/README.md) for:

- Installation instructions
- Configuration options
- Environment variables
- Development guide

---

## Original Design Documentation (Historical)

## Planned Features

```typescript
// MCP server will expose these tools to Claude:

// Find ready work
{
  "name": "beads_ready_work",
  "description": "Find issues with no blocking dependencies",
  "parameters": {
    "limit": "number",
    "priority": "number (0-4)",
    "assignee": "string"
  }
}

// Create issue
{
  "name": "beads_create_issue",
  "description": "Create a new issue",
  "parameters": {
    "title": "string",
    "description": "string",
    "priority": "number (0-4)",
    "type": "bug|feature|task|epic|chore"
  }
}

// Update issue
{
  "name": "beads_update_issue",
  "description": "Update issue status or fields",
  "parameters": {
    "id": "string",
    "status": "open|in_progress|blocked|closed",
    "priority": "number",
    "assignee": "string"
  }
}

// Add dependency
{
  "name": "beads_add_dependency",
  "description": "Link issues with dependencies",
  "parameters": {
    "from": "string (issue ID)",
    "to": "string (issue ID)",
    "type": "blocks|related|parent-child|discovered-from"
  }
}
```

## How It Would Work

1. Install the MCP server for Claude Desktop
2. Claude gains access to beads commands
3. During coding sessions, Claude can:
   - Check for ready work: "What should I work on next?"
   - Create issues: "I found a bug, let me track it"
   - Link discoveries: "This is related to the issue we're working on"
   - Update status: "Marking beads-5 as completed"

## Example Conversation

```txt
User: What's ready to work on?

Claude: Let me check the ready work queue...
[Calls beads_ready_work()]

I found 3 ready issues:
1. beads-12: Fix authentication bug (priority 0, bug)
2. beads-15: Add dark mode toggle (priority 1, feature)
3. beads-18: Update documentation (priority 2, task)

Which would you like to work on?

User: Let's do beads-12

Claude: Great! I'll mark beads-12 as in progress.
[Calls beads_update_issue("beads-12", status="in_progress")]

Let me investigate the authentication code...
[Works on the issue]

I found the bug in auth.js:45 - there's a race condition. While fixing this, I noticed we're also missing rate limiting. Should I track that as a separate issue?

User: Yes please

Claude: [Calls beads_create_issue("Add rate limiting to auth endpoints", type="task", priority=1)]
Created beads-23: Add rate limiting to auth endpoints

[Calls beads_add_dependency("beads-23", "beads-12", type="discovered-from")]
Linked beads-23 as discovered from beads-12

Now I'll complete the auth bug fix...
```

## Implementation Guide

To implement this MCP server:

1. Create a Node.js/TypeScript project
2. Use the MCP SDK: `npm install @anthropic/mcp-sdk`
3. Implement tool handlers that call `beads` commands via child_process
4. Register with Claude Desktop's config

Example skeleton:

```typescript
import { Server } from "@anthropic/mcp-sdk";
import { exec } from "child_process";
import { promisify } from "util";

const execAsync = promisify(exec);

const server = new Server({
  name: "beads",
  version: "1.0.0"
});

// Register ready work tool
server.tool("beads_ready_work", async (params) => {
  const { stdout } = await execAsync(
    `beads ready --json --limit ${params.limit || 10}`
  );
  return JSON.parse(stdout);
});

// Register create issue tool
server.tool("beads_create_issue", async (params) => {
  const { stdout } = await execAsync(
    `beads create "${params.title}" -d "${params.description}" -p ${params.priority} -t ${params.type} --json`
  );
  return JSON.parse(stdout);
});

// ... more tools ...

server.start();
```

## Installation (Future)

```sh
# Install the MCP server
npm install -g mcp-beads-server

# Configure Claude Desktop
# Add to ~/Library/Application Support/Claude/claude_desktop_config.json
{
  "mcpServers": {
    "beads": {
      "command": "mcp-beads-server",
      "args": []
    }
  }
}

# Restart Claude Desktop
```

## Alternative: Direct beads Usage

Until the MCP server is available, you can instruct Claude to use beads directly:

```markdown
# In your CLAUDE.md or project instructions:

We use Beads (beads) for issue tracking. Available commands:

- `beads ready --json` - Find ready work
- `beads create "title" -p 1 -t bug --json` - Create issue
- `beads update beads-1 --status in_progress --json` - Update status
- `beads dep add beads-2 beads-1 --type discovered-from` - Link issues
- `beads close beads-1 --reason "Done" --json` - Complete work

All commands support --json for parsing. Please use beads to track work during our sessions.
```

## Contributing

Interested in building this MCP server? We welcome contributions!

See [CONTRIBUTING.md](../../CONTRIBUTING.md) for guidelines.

## See Also

- [MCP Documentation](https://docs.anthropic.com/claude/docs/model-context-protocol)
- [MCP SDK](https://github.com/anthropics/mcp-sdk)
- [Claude Desktop](https://claude.ai/desktop)
- [../python-agent/](../python-agent/) - Python implementation pattern
