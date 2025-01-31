# Deepseek Go Client

[![Go Reference](https://pkg.go.dev/badge/github.com/yourusername/deepseek.svg)](https://pkg.go.dev/github.com/yourusername/deepseek)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/yourusername/deepseek/blob/main/LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/yourusername/deepseek)](https://golang.org/)

A modern, type-safe Go client for the Deepseek API, designed for reliability and ease of integration.

## Features

- **Full API Coverage**: Supports all Deepseek API endpoints
- **Type-Safe**: Generated from API schema definitions
- **Production Ready**:
  - Context support
  - Customizable HTTP client
  - Comprehensive error handling
  - Request/Response validation
- **Extensible Architecture**: Easy to add new API endpoints
- **Developer Friendly**:
  - Clear documentation
  - IDE-friendly design
  - Example implementations

## Installation

```bash
go get github.com/samanxzrpor/deepseek
```

## Requirements
  - Go 1.18+
  - Deepseek API key

## Quick Start
```go
package main

import (
	"context"
	"fmt"
	"os"
	
	"github.com/yourusername/deepseek"
)

func main() {
	// Initialize client
	client := deepseek.NewClient(deepseek.ClientConfig{
		APIKey: os.Getenv("DEEPSEEK_API_KEY"),
	})

	// Create chat completion
	resp, err := client.Chat.CreateCompletion(context.Background(), &deepseek.ChatCompletionRequest{
		Model: "deepseek-chat",
		Messages: []deepseek.Message{
			{Role: "user", Content: "Explain quantum computing in simple terms"},
		},
		Temperature: 0.7,
		MaxTokens:   300,
	})

	if err != nil {
		// Handle error
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
```
## Configuration
### Environment Variables
```bash
export DEEPSEEK_API_KEY="your-api-key-here"
```
## Client Options
```go
client := deepseek.NewClient(deepseek.ClientConfig{
	APIKey:     "your-api-key",
	BaseURL:    "https://api.deepseek.com/v1", // Default
	UserAgent:  "my-app/1.0",
	HTTPClient: &http.Client{Timeout: 60 * time.Second},
})
```
##Core Features
### Error Handling
```go
resp, err := client.Chat.CreateCompletion(ctx, request)
if err != nil {
	if apiErr, ok := err.(*deepseek.ErrorResponse); ok {
		fmt.Printf("API Error [%d]: %s", apiErr.HTTPStatusCode, apiErr.Error.Message)
		return
	}
	// Handle other errors
}
```
### Context Support
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

resp, err := client.Chat.CreateCompletion(ctx, request)
```
### Custom HTTP Client
```go
customClient := &http.Client{
	Timeout: 60 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        10,
		IdleConnTimeout:     30 * time.Second,
	},
}

client := deepseek.NewClient(deepseek.ClientConfig{
	APIKey:     apiKey,
	HTTPClient: customClient,
})
```
# API Reference
## Chat Completion
```go
type ChatCompletionRequest struct {
	Model            string    // Required
	Messages         []Message // Required
	Temperature      float64
	TopP             float64
	MaxTokens        int
	Stream           bool
	PresencePenalty  float64
	FrequencyPenalty float64
}

// Usage
resp, err := client.Chat.CreateCompletion(ctx, &ChatCompletionRequest{
	Model: "deepseek-chat",
	Messages: []Message{
		{Role: "system", Content: "You are a helpful assistant"},
		{Role: "user", Content: "What's the weather in London?"},
	},
})
```
## Advanced Usage
### Streaming Responses
```go
// Create streaming request
req := &deepseek.ChatCompletionRequest{
	Model:    "deepseek-chat",
	Messages: messages,
	Stream:   true,
}

stream, err := client.Chat.CreateCompletionStream(ctx, req)
if err != nil {
	return err
}

for {
	response, err := stream.Recv()
	if errors.Is(err, io.EOF) {
		break
	}
	if err != nil {
		return err
	}
	
	// Process partial response
	fmt.Printf(response.Choices[0].Delta.Content)
}
```
## Request Validation
```go
// Returns validation errors before making API call
err := req.Validate()
if err != nil {
	if validationErr, ok := err.(*deepseek.ValidationError); ok {
		for _, fieldError := range validationErr.Errors {
			fmt.Printf("%s: %s\n", fieldError.Field, fieldError.Message)
		}
	}
}
```
# Contributing
- Fork the repository
- Create a feature branch (git checkout -b feature/your-feature)
- Commit your changes (git commit -am 'Add some feature')
- Push to the branch (git push origin feature/your-feature)
- Open a Pull Request
