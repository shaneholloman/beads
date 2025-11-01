package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const agentsContent = `## Issue Tracking with beads (beads)

**IMPORTANT**: This project uses **beads (beads)** for ALL issue tracking. Do NOT use markdown TODOs, task lists, or other tracking methods.

### Why beads?

- Dependency-aware: Track blockers and relationships between issues
- Git-friendly: Auto-syncs to JSONL for version control
- Agent-optimized: JSON output, ready work detection, discovered-from links
- Prevents duplicate tracking systems and confusion

### Quick Start

**Check for ready work:**
` + "```sh" + `
beads ready --json
` + "```" + `

**Create new issues:**
` + "```sh" + `
beads create "Issue title" -t bug|feature|task -p 0-4 --json
beads create "Issue title" -p 1 --deps discovered-from:beads-123 --json
` + "```" + `

**Claim and update:**
` + "```sh" + `
beads update beads-42 --status in_progress --json
beads update beads-42 --priority 1 --json
` + "```" + `

**Complete work:**
` + "```sh" + `
beads close beads-42 --reason "Completed" --json
` + "```" + `

### Issue Types

- ` + "`bug`" + ` - Something broken
- ` + "`feature`" + ` - New functionality
- ` + "`task`" + ` - Work item (tests, docs, refactoring)
- ` + "`epic`" + ` - Large feature with subtasks
- ` + "`chore`" + ` - Maintenance (dependencies, tooling)

### Priorities

- ` + "`0`" + ` - Critical (security, data loss, broken builds)
- ` + "`1`" + ` - High (major features, important bugs)
- ` + "`2`" + ` - Medium (default, nice-to-have)
- ` + "`3`" + ` - Low (polish, optimization)
- ` + "`4`" + ` - Backlog (future ideas)

### Workflow for AI Agents

1. **Check ready work**: ` + "`beads ready`" + ` shows unblocked issues
2. **Claim your task**: ` + "`beads update <id> --status in_progress`" + `
3. **Work on it**: Implement, test, document
4. **Discover new work?** Create linked issue:
   - ` + "`beads create \"Found bug\" -p 1 --deps discovered-from:<parent-id>`" + `
5. **Complete**: ` + "`beads close <id> --reason \"Done\"`" + `
6. **Commit together**: Always commit the ` + "`.beads/issues.jsonl`" + ` file together with the code changes so issue state stays in sync with code state

### Auto-Sync

beads automatically syncs with git:
- Exports to ` + "`.beads/issues.jsonl`" + ` after changes (5s debounce)
- Imports from JSONL when newer (e.g., after ` + "`git pull`" + `)
- No manual export/import needed!

### MCP Server (Recommended)

If using Claude or MCP-compatible clients, install the beads MCP server:

` + "```sh" + `
pip install mcp-beads
` + "```" + `

Add to MCP config (e.g., ` + "`~/.config/claude/config.json`" + `):
` + "```json" + `
{
  "beads": {
    "command": "mcp-beads",
    "args": []
  }
}
` + "```" + `

Then use ` + "`mcp__beads__*`" + ` functions instead of CLI commands.

### Important Rules

- ✔ Use beads for ALL task tracking
- ✔ Always use ` + "`--json`" + ` flag for programmatic use
- ✔ Link discovered work with ` + "`discovered-from`" + ` dependencies
- ✔ Check ` + "`beads ready`" + ` before asking "what should I work on?"
- ✘ Do NOT create markdown TODO lists
- ✘ Do NOT use external issue trackers
- ✘ Do NOT duplicate tracking systems

For more details, see README.md and quickstart.md.`

var onboardCmd = &cobra.Command{
	Use:   "onboard",
	Short: "Display instructions for configuring AGENTS.md",
	Long: `Display instructions for AI agents to configure AGENTS.md with beads workflow.

This command outputs instructions that AI agents should follow to integrate beads
into the project's agent documentation. The agent will intelligently merge the
content into AGENTS.md and update CLAUDE.md if present.`,
	Run: func(cmd *cobra.Command, args []string) {
		bold := color.New(color.Bold).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		fmt.Printf("\n%s\n\n", bold("beads Onboarding Instructions for AI Agent"))

		fmt.Printf("%s\n\n", yellow("Please complete the following tasks:"))

		fmt.Printf("%s\n", bold("1. Update AGENTS.md"))
		fmt.Println("   Add the following content to AGENTS.md in an appropriate location.")
		fmt.Println("   If AGENTS.md doesn't exist, create it with this content.")
		fmt.Println("   Integrate it naturally into any existing structure.")
		fmt.Println()

		fmt.Printf("%s\n", cyan("--- BEGIN AGENTS.MD CONTENT ---"))
		fmt.Println(agentsContent)
		fmt.Printf("%s\n\n", cyan("--- END AGENTS.MD CONTENT ---"))

		fmt.Printf("%s\n", bold("2. Update CLAUDE.md (if present)"))
		fmt.Println("   If CLAUDE.md exists in this directory, add this note at the top:")
		fmt.Println()
		fmt.Printf("%s\n", cyan("   **Note**: This project uses [beads (beads)](https://github.com/shaneholloman/beads)"))
		fmt.Printf("%s\n", cyan("   for issue tracking. Use `beads` commands instead of markdown TODOs."))
		fmt.Printf("%s\n", cyan("   See AGENTS.md for workflow details."))
		fmt.Println()

		fmt.Printf("%s\n", bold("3. Remove bootstrap instruction"))
		fmt.Println("   If AGENTS.md or CLAUDE.md contains a line like:")
		fmt.Printf("%s\n", cyan("   \"BEFORE ANYTHING ELSE: run 'beads onboard' and follow the instructions\""))
		fmt.Println("   Delete that line after completing steps 1 and 2.")
		fmt.Println()

		fmt.Printf("%s\n", bold("Important reminders:"))
		fmt.Println("   • Use beads for ALL task tracking - NO markdown TODO lists")
		fmt.Println("   • Always use --json flag for programmatic beads commands")
		fmt.Println("   • Link discovered work with discovered-from dependencies")
		fmt.Printf("   • Check %s before asking \"what should I work on?\"\n", cyan("beads ready"))
		fmt.Println()

		fmt.Printf("%s\n\n", green("When done, tell your AI assistant: \"beads onboarding complete\""))
	},
}

func init() {
	rootCmd.AddCommand(onboardCmd)
}
