package ai

import (
	"context"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Client struct {
	genaiClient *genai.Client
}

func NewClient(apiKey string) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}
	log.Println("Gemini AI Client initialized successfully.")

	return &Client{genaiClient: client}, nil
}

func (c *Client) CallModel(prompt string) (string, error) {
	if c.genaiClient == nil {
		return "", fmt.Errorf("AI client not initialized")
	}

	ctx := context.Background()
	model := c.genaiClient.GenerativeModel("gemini-2.5-flash")
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))

	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		return fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0]), nil
	}

	return "", fmt.Errorf("no text content found in response")

}

func (c *Client) Close() error {
	if c.genaiClient != nil {
		log.Println("Closing Gemini AI Client.")
		return c.genaiClient.Close()
	}
	return nil
}
