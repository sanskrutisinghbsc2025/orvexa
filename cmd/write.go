package cmd

import (
	"context"
	"fmt"

	"orvexa/internal/resume"

	"github.com/spf13/cobra"
	"google.golang.org/genai"
)

func init() {
	rootCmd.AddCommand(writeCmd)
}

var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "Generate improved resume bullet points using AI",
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := getAPIKey()
		if apiKey == "" {
			fmt.Println("[ERROR] API Key not found. Run: orvexa config --key YOUR_KEY")
			return
		}

		fmt.Println("[AI] AI is improving your resume...")

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
Act as a professional technical resume writer.

Improve the resume content below.

For each important experience or project bullet:
- Make it concise and professional.
- Use strong action verbs.
- Highlight measurable impact where possible.
- Keep the claims truthful.
- Make the wording suitable for a technical resume.

Preserve the original information and do not invent experience,
qualifications, achievements, or metrics.

RESUME DATA:
%s
`, string(resumeData))

		// Use the shared fallback mechanism from ai.go.
		resp, model, err := generateWithFallback(ctx, client, prompt)
		if err != nil {
			fmt.Println("[ERROR] AI Error:", err)
			return
		}

		fmt.Printf("[AI] Response generated using %s\n", model)
		fmt.Println("\n--- AI-IMPROVED RESUME ---")
		fmt.Println(resp.Text())
		fmt.Println("\n------------------------------------------------")
	},
}
