// Package tavily provides a Go client for the Tavily AI-powered search and web content extraction API.
//
// Usage:
//
//	client := tavily.New("tvly-your-api-key")
//	resp, err := client.Search(ctx, &tavily.SearchParams{
//	    Query: "Go programming language",
//	})
package tavily

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	json "github.com/go-json-experiment/json"
)

const DefaultBaseURL = "https://api.tavily.com"

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// Option configures the Client.
type Option func(*Client)

// WithBaseURL sets a custom API base URL.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = strings.TrimSuffix(url, "/")
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

// New creates a new Tavily API client.
// If apiKey is empty, it reads from TAVILY_API_KEY environment variable.
func New(apiKey string, opts ...Option) *Client {
	if apiKey == "" {
		apiKey = os.Getenv("TAVILY_API_KEY")
	}

	c := &Client{
		baseURL:    DefaultBaseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) do(ctx context.Context, endpoint string, reqBody any, respBody any) error {
	if c.apiKey == "" {
		return &APIError{
			StatusCode: 401,
			Message:    "missing API key - provide via parameter or TAVILY_API_KEY environment variable",
		}
	}

	var body io.Reader
	if reqBody != nil {
		data, err := json.Marshal(reqBody, json.DefaultOptionsV2())
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}
		body = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+endpoint, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return parseAPIError(resp.StatusCode, respData)
	}

	if respBody != nil {
		if err := json.Unmarshal(respData, respBody, json.DefaultOptionsV2()); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

func parseAPIError(statusCode int, respData []byte) error {
	var errorResp struct {
		Detail struct {
			Error string `json:"error"`
		} `json:"detail"`
	}

	message := "unknown error"
	if json.Unmarshal(respData, &errorResp, json.DefaultOptionsV2()) == nil && errorResp.Detail.Error != "" {
		message = errorResp.Detail.Error
	}

	return &APIError{
		StatusCode: statusCode,
		Message:    message,
	}
}
