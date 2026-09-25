package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const CurrentVersion = "v1.0.1"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Orvexa",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Orvexa %s - The Resume Versioning Network\n", CurrentVersion)
	},
}
