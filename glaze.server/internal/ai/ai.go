package ai

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"glaze/logger"

	"go.uber.org/zap"
	"google.golang.org/genai"
)

type Models string

const (
	GEMINI_FLASH      Models = "gemini-2.5-flash"
	GEMINI_FLASH_LITE Models = "gemini-2.5-flash-lite"
)

type GenAIClient struct {
	g *genai.Client
}

func NewClient(ctx context.Context, apiKey string) (*GenAIClient, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("api key is required")
	}

	genAIClient, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		logger.Logger.Error("Error creating GenAI client", zap.Error(err))
		return nil, err
	}

	g := &GenAIClient{
		genAIClient,
	}

	return g, nil
}

func (g *GenAIClient) GetStructuredContent(ctx context.Context, prompt string, model Models) (string, error) {
	// Create the chat session
	chat, err := g.g.Chats.Create(ctx, string(model), nil, nil)
	if err != nil {
		logger.Logger.Error("Error creating chat", zap.Error(err))
		return "", err
	}

	// Send the message
	result, err := chat.SendMessage(ctx, genai.Part{Text: prompt})
	if err != nil {
		logger.Logger.Error("Error sending message", zap.Error(err))
		return "", err
	}

	// 1. Check if there are any candidates in the response
	if len(result.Candidates) == 0 || result.Candidates[0].Content == nil {
		return "", errors.New("no response generated from the model")
	}

	// 2. Extract the text from the first part of the first candidate
	// Most simple text prompts return a single candidate with a single part
	var sb strings.Builder
	for _, part := range result.Candidates[0].Content.Parts {
		if part.Text != "" {
			sb.WriteString(part.Text)
		}
	}

	responseText := sb.String()
	if responseText == "" {
		return "", errors.New("model returned an empty response")
	}

	return responseText, nil
}

func debugPrint[T any](r *T) string {
	response, err := json.MarshalIndent(*r, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	return string(response)
}
