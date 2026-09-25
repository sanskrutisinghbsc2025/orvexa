package cmd

import (
	"fmt"
	"orvexa/internal/ui"
	"orvexa/internal/vcs"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(historyCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all saved versions (commits)",
	Run: func(cmd *cobra.Command, args []string) {
		r, err := vcs.Open()
		if err != nil {
			ui.PrintError(err.Error())
			return
		}

		logs, err := r.Log()
		if err != nil {
			ui.PrintInfo("No versions found yet.")
			return
		}

		ui.PrintHeader("Version History")
		logs.ForEach(func(c *object.Commit) error {
			date := c.Author.When.Format("Jan 02, 2006")
			hash := c.Hash.String()[:7]
			fmt.Printf("  ")
			ui.PrintKV("["+hash+"]", c.Message+" ("+date+")")
			return nil
		})
	},
}

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show a visual timeline of your resume evolution",
	Run: func(cmd *cobra.Command, args []string) {
		r, err := vcs.Open()
		if err != nil {
			ui.PrintError(err.Error())
			return
		}

		logs, err := r.Log()
		if err != nil {
			ui.PrintInfo("No history found.")
			return
		}

		ui.PrintHeader("Resume Network Timeline")
		fmt.Println()

		logs.ForEach(func(c *object.Commit) error {
			date := c.Author.When.Format("2006-01-02 15:04")
			hash := c.Hash.String()[:7]

			fmt.Printf("  ")
			ui.PrintSuccess("[" + hash + "] " + date)
			fmt.Printf("  ┃  \n")
			fmt.Printf("  ┗━▶ %s\n\n", c.Message)

			return nil
		})
	},
}
