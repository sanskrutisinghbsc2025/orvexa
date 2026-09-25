package cmd

import (
	"fmt"
	"orvexa/internal/vcs"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	Use:   "sync [branch]",
	Short: "Sync (rebase) current branch with another branch",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetBranch := args[0]

		r, err := vcs.Open()
		if err != nil {
			fmt.Println("[ERROR]", err)
			return
		}

		currentBranch, _ := r.CurrentBranch()
		fmt.Printf("[INFO] Syncing %s with updates from %s...\n", currentBranch, targetBranch)

		out, err := r.Sync(targetBranch)
		if err != nil {
			fmt.Println("[ERROR] Sync conflict detected!")
			fmt.Println(string(out))
			fmt.Println("[INFO] Use standard git tools to resolve conflicts, or 'git rebase --abort'.")
			return
		}

		fmt.Printf("[SUCCESS] %s is now up-to-date with %s.\n", currentBranch, targetBranch)
	},
}
