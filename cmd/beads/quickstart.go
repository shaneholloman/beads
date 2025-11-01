package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var quickstartCmd = &cobra.Command{
	Use:   "quickstart",
	Short: "Quick start guide for beads",
	Long:  `Display a quick start guide showing common beads workflows and patterns.`,
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		bold := color.New(color.Bold).SprintFunc()

		fmt.Printf("\n%s\n\n", bold("beads - Dependency-Aware Issue Tracker"))
		fmt.Printf("Issues chained together like beads.\n\n")

		fmt.Printf("%s\n", bold("GETTING STARTED"))
		fmt.Printf("  %s   Initialize beads in your project\n", cyan("beads init"))
		fmt.Printf("            Creates .beads/ directory with project-specific database\n")
		fmt.Printf("            Auto-detects prefix from directory name (e.g., myapp-1, myapp-2)\n\n")

		fmt.Printf("  %s   Initialize with custom prefix\n", cyan("beads init --prefix api"))
		fmt.Printf("            Issues will be named: api-1, api-2, ...\n\n")

		fmt.Printf("%s\n", bold("CREATING ISSUES"))
		fmt.Printf("  %s\n", cyan("beads create \"Fix login bug\""))
		fmt.Printf("  %s\n", cyan("beads create \"Add auth\" -p 0 -t feature"))
		fmt.Printf("  %s\n\n", cyan("beads create \"Write tests\" -d \"Unit tests for auth\" --assignee alice"))

		fmt.Printf("%s\n", bold("VIEWING ISSUES"))
		fmt.Printf("  %s       List all issues\n", cyan("beads list"))
		fmt.Printf("  %s  List by status\n", cyan("beads list --status open"))
		fmt.Printf("  %s  List by priority (0-4, 0=highest)\n", cyan("beads list --priority 0"))
		fmt.Printf("  %s       Show issue details\n\n", cyan("beads show beads-1"))

		fmt.Printf("%s\n", bold("MANAGING DEPENDENCIES"))
		fmt.Printf("  %s     Add dependency (beads-2 blocks beads-1)\n", cyan("beads dep add beads-1 beads-2"))
		fmt.Printf("  %s  Visualize dependency tree\n", cyan("beads dep tree beads-1"))
		fmt.Printf("  %s      Detect circular dependencies\n\n", cyan("beads dep cycles"))

		fmt.Printf("%s\n", bold("DEPENDENCY TYPES"))
		fmt.Printf("  %s  Task B must complete before task A\n", yellow("blocks"))
		fmt.Printf("  %s  Soft connection, doesn't block progress\n", yellow("related"))
		fmt.Printf("  %s  Epic/subtask hierarchical relationship\n", yellow("parent-child"))
		fmt.Printf("  %s  Auto-created when AI discovers related work\n\n", yellow("discovered-from"))

		fmt.Printf("%s\n", bold("READY WORK"))
		fmt.Printf("  %s       Show issues ready to work on\n", cyan("beads ready"))
		fmt.Printf("            Ready = status is 'open' AND no blocking dependencies\n")
		fmt.Printf("            Perfect for agents to claim next work!\n\n")

		fmt.Printf("%s\n", bold("UPDATING ISSUES"))
		fmt.Printf("  %s\n", cyan("beads update beads-1 --status in_progress"))
		fmt.Printf("  %s\n", cyan("beads update beads-1 --priority 0"))
		fmt.Printf("  %s\n\n", cyan("beads update beads-1 --assignee bob"))

		fmt.Printf("%s\n", bold("CLOSING ISSUES"))
		fmt.Printf("  %s\n", cyan("beads close beads-1"))
		fmt.Printf("  %s\n\n", cyan("beads close beads-2 beads-3 --reason \"Fixed in PR #42\""))

		fmt.Printf("%s\n", bold("DATABASE LOCATION"))
		fmt.Printf("  beads automatically discovers your database:\n")
		fmt.Printf("    1. %s flag\n", cyan("--db /path/to/db.db"))
		fmt.Printf("    2. %s environment variable\n", cyan("$BEADS_DB"))
		fmt.Printf("    3. %s in current directory or ancestors\n", cyan(".beads/*.db"))
		fmt.Printf("    4. %s as fallback\n\n", cyan("~/.beads/default.db"))

		fmt.Printf("%s\n", bold("AGENT INTEGRATION"))
		fmt.Printf("  beads is designed for AI-supervised workflows:\n")
		fmt.Printf("    • Agents create issues when discovering new work\n")
		fmt.Printf("    • %s shows unblocked work ready to claim\n", cyan("beads ready"))
		fmt.Printf("    • Use %s flags for programmatic parsing\n", cyan("--json"))
		fmt.Printf("    • Dependencies prevent agents from duplicating effort\n\n")

		fmt.Printf("%s\n", bold("DATABASE EXTENSION"))
		fmt.Printf("  Applications can extend beads's SQLite database:\n")
		fmt.Printf("    • Add your own tables (e.g., %s)\n", cyan("myapp_executions"))
		fmt.Printf("    • Join with %s table for powerful queries\n", cyan("issues"))
		fmt.Printf("    • See database extension docs for integration patterns:\n")
		fmt.Printf("      %s\n\n", cyan("https://github.com/shaneholloman/beads/blob/main/extending.md"))

		fmt.Printf("%s\n", bold("GIT WORKFLOW (AUTO-SYNC)"))
		fmt.Printf("  beads automatically keeps git in sync:\n")
		fmt.Printf("    • %s Export to JSONL after CRUD operations (5s debounce)\n", green("✔"))
		fmt.Printf("    • %s Import from JSONL when newer than DB (after %s)\n", green("✔"), cyan("git pull"))
		fmt.Printf("    • %s Works seamlessly across machines and team members\n", green("✔"))
		fmt.Printf("    • No manual export/import needed!\n")
		fmt.Printf("  Disable with: %s or %s\n\n", cyan("--no-auto-flush"), cyan("--no-auto-import"))

		fmt.Printf("%s\n", green("Ready to start!"))
		fmt.Printf("Run %s to create your first issue.\n\n", cyan("beads create \"My first issue\""))
	},
}

func init() {
	rootCmd.AddCommand(quickstartCmd)
}
