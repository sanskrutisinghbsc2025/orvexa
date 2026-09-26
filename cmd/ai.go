package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/genai"
)

const geminiModel = "gemini-3.8-flash"

// generateAIContent sends a request to Gemini and retries temporary
// service-unavailable errors automatically.
func generateAIContent(ctx context.Context, client *genai.Client, prompt string) (*genai.GenerateContentResponse, error) {
	var lastErr error

	for attempt := 1; attempt <= 3; attempt++ {
		resp, err := client.Models.GenerateContent(
			ctx,
			geminiModel,
			genai.Text(prompt),
			nil,
		)

		if err == nil {
			return resp, nil
		}

		lastErr = err

		errText := strings.ToLower(err.Error())

		// Retry temporary service/capacity errors.
		if strings.Contains(errText, "503") ||
			strings.Contains(errText, "unavailable") ||
			strings.Contains(errText, "high demand") {
			if attempt < 3 {
				wait := time.Duration(attempt*2) * time.Second
				fmt.Printf("[AI] Gemini is temporarily busy. Retrying in %d seconds... (%d/3)\n",
					int(wait.Seconds()), attempt+1)
				time.Sleep(wait)
				continue
			}
		}

		return nil, err
	}

	return nil, lastErr
}
