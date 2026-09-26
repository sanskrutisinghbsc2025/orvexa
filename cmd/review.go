package cmd

import (
	"context"
	"fmt"
	"orvexa/internal/resume"

	"github.com/spf13/cobra"
	"google.golang.org/genai"
)

func init() {
	rootCmd.AddCommand(reviewCmd)
	reviewCmd.Flags().StringP("role", "r", "General Software Engineer", "The specific job role you are targeting")
}

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Get an AI recruiter to critique your resume for a specific role",
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := getAPIKey()
		if apiKey == "" {
			fmt.Println("[ERROR] API Key not found. Run: orvexa config --key YOUR_KEY")
			return
		}

		targetRole, _ := cmd.Flags().GetString("role")

		fmt.Printf("[AI] AI Recruiter is analyzing your resume for the role: [%s]...\n", targetRole)

		resumeData, err := resume.ReadRaw()
		if err != nil {
			fmt.Println("[ERROR]", err)
			return
		}

		ctx := context.Background()

		client, err := genai.NewClient(ctx, &genai.ClientConfig{
			APIKey: apiKey,
		})
		if err != nil {
			fmt.Println("[ERROR] Error connecting to AI:", err)
			return
		}

		prompt := fmt.Sprintf(`
Act as a Senior Technical Recruiter who specializes in hiring for %s roles.
I am going to provide you with a resume in JSON format.

Critique this resume ONLY through the lens of a %s position.

Please provide:
1. A Role-Match Score (0-100).
2. Missing Keywords: What technologies or skills typical for a %s are missing?
3. Bullet Point Strengthening: Pick 2 bullet points and show how to make them more impressive for this specific role.
4. "Red Flags": Anything that would make a recruiter for this role hesitate.

RESUME DATA:
%s
`, targetRole, targetRole, targetRole, string(resumeData))

		resp, err := generateAIContent(ctx, client, prompt)
		if err != nil {
			fmt.Println("[ERROR] AI Error:", err)
			return
		}

		fmt.Printf("\n--- AI RECRUITER FEEDBACK FOR: %s ---\n", targetRole)
		fmt.Println(resp.Text())
		fmt.Println("\n------------------------------------------------")
	},
}
