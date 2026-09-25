package cmd

import (
	"context"
	"fmt"
	"orvexa/internal/resume"

	"github.com/spf13/cobra"
	"google.golang.org/genai"
)

func init() {
	rootCmd.AddCommand(scoreCmd)
}

var scoreCmd = &cobra.Command{
	Use:   "score",
	Short: "Get an AI score for your resume",
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := getAPIKey()
		if apiKey == "" {
			fmt.Println("[ERROR] API Key not found. Run: orvexa config --key YOUR_KEY")
			return
		}

		fmt.Println("[AI] AI is analyzing your resume...")

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
Analyze the following resume as an experienced technical recruiter.

Give:
1. Overall Resume Score (0-100)
2. Strengths
3. Weaknesses
4. Missing Skills
5. Specific Improvements

RESUME DATA:
%s
`, string(resumeData))

		resp, err := client.Models.GenerateContent(
			ctx,
			"gemini-3.5-flash",
			genai.Text(prompt),
			nil,
		)
		if err != nil {
			fmt.Println("[ERROR] AI Error:", err)
			return
		}

		fmt.Println("\n--- AI RESUME SCORE ---")
		fmt.Println(resp.Text())
		fmt.Println("\n------------------------------------------------")
	},
}
