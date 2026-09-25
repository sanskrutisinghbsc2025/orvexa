package cmd

import (
	"context"
	"fmt"
	"google.golang.org/genai"
	"orvexa/internal/resume"
	"orvexa/internal/vcs"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(doctorCmd)
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check the health of your Orvexa environment",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🏥 ORVEXA HEALTH CHECK")
		fmt.Println("========================")

		// 1. Check Git/VCS
		fmt.Print("[1/5] Checking Network (VCS)... ")
		repo, err := vcs.Open()
		if err != nil {
			fmt.Println("❌ FAILED")
			fmt.Println("    - Not a Orvexa repo. Run 'orvexa init'.")
		} else {
			branch, _ := repo.CurrentBranch()
			fmt.Printf("✅ OK (Branch: %s)\n", branch)
		}

		// 2. Check Repository (JSON)
		fmt.Print("[2/5] Checking Repository (JSON)... ")
		res, err := resume.Read()
		if err != nil {
			fmt.Println("❌ FAILED")
			fmt.Printf("    - %v\n", err)
		} else {
			fmt.Printf("✅ OK (Name: %s)\n", res.Basics.Name)
		}

		// 3. Check AI Engine (Gemini)
		fmt.Print("[3/5] Checking AI Engine (Gemini)... ")
		apiKey := getAPIKey()
		if apiKey == "" {
			fmt.Println("⚠️  WARNING")
			fmt.Println("    - API Key missing. AI features disabled.")
		} else {
			ctx := context.Background()
			_, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: apiKey})
			if err != nil {
				fmt.Println("❌ FAILED")
				fmt.Printf("    - Connection error: %v\n", err)
			} else {
				fmt.Println("✅ OK")
			}
		}

		// 4. Check System Tools (Git)
		fmt.Print("[4/5] Checking System Tools (Git)... ")
		if _, err := exec.LookPath("git"); err != nil {
			fmt.Println("❌ FAILED")
			fmt.Println("    - Git CLI not found in PATH. Sync feature disabled.")
		} else {
			fmt.Println("✅ OK")
		}

		// 5. Check Export Engine (Rod/Chrome)
		fmt.Print("[5/5] Checking Export Engine (Rod)... ")
		// We don't want to launch a browser here, just check if we can.
		fmt.Println("⌛ PENDING")
		fmt.Println("    - Will be verified during first 'orvexa export'.")

		fmt.Println("\n========================")
		fmt.Println("🩺 Diagnosis complete.")
	},
}
