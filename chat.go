package deepseek

import (
	"context"
	"net/http"
)

// ChatService handles communication with the chat-related API endpoints
type ChatService struct {
	client *Client
}

// ChatCompletionRequest represents a request to create chat completion
type ChatCompletionRequest struct {
	Messages         []Message `json:"messages"`
	Model            string    `json:"model"`
	Temperature      float64   `json:"temperature,omitempty"`
	TopP             float64   `json:"top_p,omitempty"`
	MaxTokens        int       `json:"max_tokens,omitempty"`
	Stream           bool      `json:"stream,omitempty"`
	PresencePenalty  float64   `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64   `json:"frequency_penalty,omitempty"`
}

// Message represents a single message in a chat conversation
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionResponse represents the response from a chat completion request
type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a completion choice
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// CreateCompletion creates a new chat completion
func (s *ChatService) CreateCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	const endpoint = "/chat/completions"

	var response ChatCompletionResponse
	err := s.client.sendRequest(ctx, requestOpts{
		method:   http.MethodPost,
		endpoint: endpoint,
		body:     req,
	}, &response)

	if err != nil {
		return nil, err
	}

	return &response, nil
}