package deepseek

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultBaseURL    = "https://api.deepseek.com/v1"
	defaultTimeout    = 30 * time.Second
	defaultUserAgent  = "deepseek-go/1.0"
	contentTypeJSON   = "application/json"
	headerAPIKey      = "Authorization"
	headerContentType = "Content-Type"
	headerUserAgent   = "User-Agent"
)

// ClientConfig holds configuration for the Deepseek client
type ClientConfig struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
	UserAgent  string
}

// Client is the main API client structure
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	userAgent  string

	// API Services
	Chat *ChatService
}

// NewClient creates a new Deepseek API client
func NewClient(config ClientConfig) *Client {
	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{Timeout: defaultTimeout}
	}

	if config.BaseURL == "" {
		config.BaseURL = defaultBaseURL
	}

	if config.UserAgent == "" {
		config.UserAgent = defaultUserAgent
	}

	c := &Client{
		apiKey:     config.APIKey,
		baseURL:    config.BaseURL,
		httpClient: config.HTTPClient,
		userAgent:  config.UserAgent,
	}

	// Initialize services
	c.Chat = &ChatService{client: c}

	return c
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	HTTPStatusCode int    `json:"-"`
	APIError       struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"error"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("api error: [%s] %s (status %d)", e.APIError.Code, e.APIError.Message, e.HTTPStatusCode)
}

// requestOpts holds parameters for building requests
type requestOpts struct {
	method   string
	endpoint string
	body     interface{}
}

// sendRequest handles the HTTP request/response cycle
func (c *Client) sendRequest(ctx context.Context, opts requestOpts, result interface{}) error {
	var body io.Reader
	if opts.body != nil {
		jsonBody, err := json.Marshal(opts.body)
		if err != nil {
			return fmt.Errorf("marshal request body failed: %w", err)
		}
		body = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, opts.method, c.baseURL+opts.endpoint, body)
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}

	// Set headers
	req.Header.Set(headerAPIKey, "Bearer "+c.apiKey)
	req.Header.Set(headerContentType, contentTypeJSON)
	req.Header.Set(headerUserAgent, c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		errorResp := ErrorResponse{HTTPStatusCode: resp.StatusCode}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return fmt.Errorf("failed to decode error response (status %d): %w", resp.StatusCode, err)
		}
		return &errorResp
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("decode response failed: %w", err)
		}
	}

	return nil
}