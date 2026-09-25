package cmd

import (
	"orvexa/internal/ui"
	"orvexa/internal/vcs"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current branch and unsaved changes",
	Run: func(cmd *cobra.Command, args []string) {
		r, err := vcs.Open()
		if err != nil {
			ui.PrintError(err.Error())
			return
		}

		ui.PrintHeader("Network Status")

		// 1. Get current branch
		branch, err := r.CurrentBranch()
		if err != nil {
			ui.PrintKV("Active Branch", "(initial branch)")
		} else {
			ui.PrintKV("Active Branch", branch)
		}

		// 2. Check for changes
		status, err := r.Status()
		if err != nil {
			ui.PrintError("Failed to get network status")
			return
		}

		if status.IsClean() {
			ui.PrintSuccess("Orvexa network is healthy and synchronized.")
			return
		}

		// resume.json may not exist in the status map when there
		// are changes only in other files, so check safely.
		resumeStatus, exists := status["resume.json"]

		if exists && resumeStatus.Worktree == git.Unmodified {
			ui.PrintSuccess("Orvexa network is healthy (metadata changes ignored).")
		} else {
			ui.PrintWarning("Uncommitted changes detected in the network.")
			ui.PrintInfo("Run 'orvexa commit' to protect this version.")
		}
	},
}
