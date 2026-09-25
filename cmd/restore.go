package cmd

import (
	"fmt"
	"orvexa/internal/vcs"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(restoreCmd)
	restoreCmd.Flags().BoolP("force", "f", false, "Force restore even if there are unsaved changes")
}

var restoreCmd = &cobra.Command{
	Use:   "restore [hash]",
	Short: "Restore resume.json to a previous version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shortHash := args[0]
		force, _ := cmd.Flags().GetBool("force")

		r, err := vcs.Open()
		if err != nil {
			fmt.Println("[ERROR]", err)
			return
		}

		fullHash, err := r.Restore(shortHash, force)
		if err != nil {
			fmt.Println("[ERROR]", err)
			return
		}

		fmt.Printf("[SUCCESS] Network restored to version [%s].\n", fullHash[:7])
	},
}
