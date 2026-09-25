package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "orvexa",
	Short: "Orvexa: The Resume Versioning Network",
	Long: `Orvexa is a professional version control system for career data.
It manages resume iterations via Git-based branching, provides semantic
analysis, and automates high-fidelity PDF generation.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "[ERROR]", err)
		os.Exit(1)
	}
}

const BrandASCII = `
   ██████╗ ██████╗ ██╗   ██╗███████╗██╗  ██╗ █████╗
  ██╔═══██╗██╔══██╗██║   ██║██╔════╝╚██╗██╔╝██╔══██╗
  ██║   ██║██████╔╝██║   ██║█████╗   ╚███╔╝ ███████║
  ██║   ██║██╔══██╗╚██╗ ██╔╝██╔══╝   ██╔██╗ ██╔══██║
  ╚██████╔╝██║  ██║ ╚████╔╝ ███████╗██╔╝ ██╗██║  ██║
   ╚═════╝ ╚═╝  ╚═╝  ╚═══╝  ╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝

             >>> RESUME VERSIONING SYSTEM <<<
`

func printBrand() {
	fmt.Print(BrandASCII)
}
