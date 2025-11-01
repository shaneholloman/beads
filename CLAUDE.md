
<!-- beads integration note -->
**Note**: This project uses [beads (beads)](https://github.com/shaneholloman/beads) for issue tracking. Use `beads` commands or the beads MCP server instead of markdown TODOs. See AGENTS.md for workflow details.
<!-- /beads integration note -->

# Instructions for Claude

This file has been moved to **AGENTS.md** to support all AI agents, not just Claude.

Please refer to [AGENTS.md](./AGENTS.md) for complete instructions on working with the beads project.

## Critical Directive: Never Second-Guess Shane's Strong Assertions

When Shane makes a strong, direct assertion (especially about installing tools, packages, or technical decisions):

**DO**: Execute the request immediately without questioning
**DON'T**: Argue, assert it doesn't exist, or contradict without investigating first

If Shane says "install X" or "use Y", assume he knows what he's talking about and do it. If something genuinely doesn't work after attempting it, then report the actual error.

**Never waste time with confident assertions about what exists or doesn't exist without verification.**

This applies especially to:
- Package installations (e.g., "install ty" - just do it, don't say it doesn't exist)
- Tool usage (e.g., "use uv add" - just do it)
- Technical decisions (e.g., "remove hatchling" - just do it)

Shane has more context than you do. Trust his judgment and execute.
