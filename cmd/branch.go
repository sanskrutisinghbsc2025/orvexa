package cmd

import (
	"orvexa/internal/ui"
	"orvexa/internal/vcs"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(branchCmd)
	branchCmd.AddCommand(branchCreateCmd)
	branchCmd.AddCommand(branchSwitchCmd)
	branchCmd.AddCommand(branchListCmd)
	branchCmd.AddCommand(branchDeleteCmd)
}

var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "Manage resume branches (e.g. for different job roles)",
}

var branchListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all resume branches",
	Run: func(cmd *cobra.Command, args []string) {
		r, err := vcs.Open()
		if err != nil {
			ui.PrintError(err.Error())
			return
		}

		current, _ := r.CurrentBranch()
		branches, err := r.ListBranches()
		if err != nil {
			ui.PrintError("Failed to list branches")
			return
		}

		ui.PrintHeader("Resume Branches")
		for _, name := range branches {
			ui.PrintBranch(name, name == current)
		}
	},
}

var branchCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new resume branch",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		r, err := vcs.Open()
		if err != nil {
			ui.PrintError(err.Error())
			return
		}

		name := args[0]
		err = r.CreateBranch(name)
		if err != nil {
			ui.PrintError("Failed to create branch: " + err.Error())
		} else {
			ui.PrintSuccess("Branch '" + name + "' created and active.")
		}
	},
}

var branchSwitchCmd = &cobra.Command{
	Use:   "switch [name]",
	Short: "Switch to a different resume branch",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		r, err := vcs.Open()
		if err != nil {
			ui.PrintError(err.Error())
			return
		}

		name := args[0]
		err = r.SwitchBranch(name)
		if err != nil {
			ui.PrintError("Branch '" + name + "' does not exist.")
		} else {
			ui.PrintSuccess("Switched to branch '" + name + "'.")
		}
	},
}

var branchDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a resume branch",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		r, err := vcs.Open()
		if err != nil {
			ui.PrintError(err.Error())
			return
		}

		name := args[0]
		err = r.DeleteBranch(name)
		if err != nil {
			ui.PrintError(err.Error())
		} else {
			ui.PrintSuccess("Branch '" + name + "' deleted.")
		}
	},
}
