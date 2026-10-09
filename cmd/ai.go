package cmd

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

// Models are tried in order.
var fallbackModels = []string{
	"gemini-3.8-flash",
	"gemini-3.7-flash",
	"gemini-3.6-flash",
	"gemini-3.5-flash",
	"gemini-3.5-flash-lite",
	"gemini-flash-latest",
}

// generateAIContent keeps the existing interface used by score and review.
func generateAIContent(
	ctx context.Context,
	client *genai.Client,
	prompt string,
) (*genai.GenerateContentResponse, error) {
	resp, _, err := generateWithFallback(ctx, client, prompt)
	return resp, err
}

// generateWithFallback tries each configured model until one succeeds.
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
			fmt.Printf("[AI] Successfully used %s\n", model)
			return resp, model, nil
		}

		lastErr = err

		if !isTransientAIError(err) {
			return nil, "", fmt.Errorf(
				"model %s failed with a non-retryable error: %w",
				model, err,
			)
		}

		fmt.Printf("[AI] %s is temporarily unavailable: %v\n", model, err)
		fmt.Println("[AI] Trying the next fallback model...")
	}

	return nil, "", fmt.Errorf(
		"all %d fallback models failed; last error: %w",
		len(fallbackModels), lastErr,
	)
}

// isTransientAIError identifies errors for which another model may help.
func isTransientAIError(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())

	transientErrors := []string{
		"429",
		"resource_exhausted",
		"rate limit",
		"too many requests",
		"500",
		"internal",
		"502",
		"503",
		"504",
		"unavailable",
		"high demand",
		"overloaded",
		"timeout",
		"deadline exceeded",
	}

	for _, item := range transientErrors {
		if strings.Contains(msg, item) {
			return true
		}
	}

	return false
}
