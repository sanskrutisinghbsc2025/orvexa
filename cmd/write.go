package cmd

import (
	"context"
	"fmt"
	"strings"

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

RESUME DATA:
%s
`, string(resumeData))

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

// Models are tried in order.
// If the first model is temporarily unavailable, Orvexa moves to the next one.
var fallbackModels = []string{
	"gemini-3.8-flash",
	"gemini-3.7-flash",
	"gemini-3.6-flash",
	"gemini-3.5-flash",
	"gemini-3.5-flash-lite",
}

func generateWithFallback(
	ctx context.Context,
	client *genai.Client,
	prompt string,
) (*genai.GenerateContentResponse, string, error) {

	var lastErr error

	for _, model := range fallbackModels {
		fmt.Printf("[AI] Trying %s...\n", model)

		resp, err := client.Models.GenerateContent(
			ctx,
			model,
			genai.Text(prompt),
			nil,
		)

		if err == nil {
			return resp, model, nil
		}

		lastErr = err

		// Only fail over for temporary/service/rate-limit errors.
		if !isTransientAIError(err) {
			return nil, "", err
		}

		fmt.Printf("[AI] %s unavailable. Trying next model...\n", model)
	}

	return nil, "", fmt.Errorf(
		"all fallback models failed; last error: %v",
		lastErr,
	)
}

func isTransientAIError(err error) bool {
	msg := strings.ToLower(err.Error())

	transientErrors := []string{
		"429",
		"resource_exhausted",
		"rate limit",
		"too many requests",
		"500",
		"internal",
		"503",
		"unavailable",
		"service unavailable",
		"timeout",
	}

	for _, item := range transientErrors {
		if strings.Contains(msg, item) {
			return true
		}
	}

	return false
}