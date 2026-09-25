package cmd

import (
	"fmt"
	"orvexa/internal/vcs"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().StringP("message", "m", "", "Commit message")
}

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Save current state of the resume",
	Run: func(cmd *cobra.Command, args []string) {
		msg, _ := cmd.Flags().GetString("message")
		if msg == "" {
			fmt.Println("[ERROR] Please provide a commit message: orvexa commit -m 'Updated skills'")
			return
		}

		r, err := vcs.Open()
		if err != nil {
			fmt.Println("[ERROR]", err)
			return
		}

		hash, err := r.Commit(msg)
		if err != nil {
			fmt.Println("[ERROR] Error committing:", err)
			return
		}

		fmt.Printf("[SUCCESS] Version Saved! [%s] %s\n", hash[:7], msg)
	},
}
