package cmd

import (
	"orvexa/internal/resume"
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

		// Only resume.json matters: 'orvexa commit' never saves other files
		// (e.g. .gitignore), so they must not trigger the warning. An unchanged
		// resume.json is absent from the status map (Status.File() would report
		// it as Untracked, so use a plain lookup). Check both staging and
		// worktree so staged-only changes (e.g. after 'orvexa restore') count.
		fs, listed := status[resume.ResumeFile]

		if !listed || (fs.Staging == git.Unmodified && fs.Worktree == git.Unmodified) {
			ui.PrintSuccess("Orvexa network is healthy (metadata changes ignored).")
		} else {
			ui.PrintWarning("Uncommitted changes detected in the network.")
			ui.PrintInfo("Run 'orvexa commit' to protect this version.")
		}
	},
}
