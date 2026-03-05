package tavily

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("with api key", func(t *testing.T) {
		client := New("tvly-test-key")
		if client.apiKey != "tvly-test-key" {
			t.Errorf("apiKey = %v, want tvly-test-key", client.apiKey)
		}
		if client.baseURL != DefaultBaseURL {
			t.Errorf("baseURL = %v, want %v", client.baseURL, DefaultBaseURL)
		}
	})

	t.Run("with options", func(t *testing.T) {
		custom := &http.Client{}
		client := New("tvly-test-key",
			WithBaseURL("https://custom.api.com/"),
			WithHTTPClient(custom),
		)
		if client.baseURL != "https://custom.api.com" {
			t.Errorf("baseURL = %v, want https://custom.api.com", client.baseURL)
		}
		if client.httpClient != custom {
			t.Error("httpClient not set correctly")
		}
	})

	t.Run("empty api key", func(t *testing.T) {
		client := New("")
		if client.apiKey != "" {
			t.Errorf("apiKey = %v, want empty", client.apiKey)
		}
	})
}

func TestAPIError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
		checkFunc  func(*APIError) bool
	}{
		{"bad request", 400, "Invalid parameters", (*APIError).IsBadRequest},
		{"unauthorized", 401, "Invalid API key", (*APIError).IsUnauthorized},
		{"rate limit", 429, "Rate limit exceeded", (*APIError).IsRateLimit},
		{"plan limit", 432, "Plan limit exceeded", (*APIError).IsPlanLimitExceeded},
		{"pay-go limit", 433, "Pay-go limit exceeded", (*APIError).IsPayGoLimitExceeded},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &APIError{StatusCode: tt.statusCode, Message: tt.message}
			if err.Error() != tt.message {
				t.Errorf("Error() = %v, want %v", err.Error(), tt.message)
			}
			if !tt.checkFunc(err) {
				t.Errorf("check function returned false for status %d", tt.statusCode)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if !strings.Contains(r.Header.Get("Authorization"), "Bearer") {
			t.Error("Expected Authorization header with Bearer token")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"query": "test query",
			"response_time": 0.5,
			"images": [],
			"results": [
				{
					"title": "Test Result",
					"url": "https://example.com",
					"content": "Test content",
					"score": 0.95
				}
			]
		}`))
	}))
	defer server.Close()

	client := New("tvly-test-key", WithBaseURL(server.URL))
	ctx := context.Background()

	result, err := client.Search(ctx, &SearchParams{Query: "test query"})
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if result.Query != "test query" {
		t.Errorf("query = %v, want test query", result.Query)
	}
	if len(result.Results) != 1 {
		t.Fatalf("results count = %v, want 1", len(result.Results))
	}
	if result.Results[0].Title != "Test Result" {
		t.Errorf("result title = %v, want Test Result", result.Results[0].Title)
	}
}

func TestSearchWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"query": "test query",
			"answer": "Test answer",
			"response_time": 0.5,
			"results": []
		}`))
	}))
	defer server.Close()

	client := New("tvly-test-key", WithBaseURL(server.URL))
	answerMode := IncludeAnswerBasic

	result, err := client.Search(context.Background(), &SearchParams{
		Query:         "test query",
		SearchDepth:   SearchDepthAdvanced,
		Topic:         TopicNews,
		MaxResults:    10,
		IncludeAnswer: &answerMode,
	})
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if result.Answer != "Test answer" {
		t.Errorf("answer = %v, want Test answer", result.Answer)
	}
}

func TestExtract(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"response_time": 0.5,
			"results": [
				{
					"url": "https://example.com",
					"raw_content": "Test content",
					"images": ["https://example.com/image.jpg"]
				}
			],
			"failed_results": []
		}`))
	}))
	defer server.Close()

	client := New("tvly-test-key", WithBaseURL(server.URL))

	result, err := client.Extract(context.Background(), &ExtractParams{
		URLs: []string{"https://example.com"},
	})
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	if len(result.Results) != 1 {
		t.Fatalf("results count = %v, want 1", len(result.Results))
	}
	if result.Results[0].URL != "https://example.com" {
		t.Errorf("result URL = %v, want https://example.com", result.Results[0].URL)
	}
}

func TestErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"detail": {"error": "Invalid API key provided"}}`))
	}))
	defer server.Close()

	client := New("invalid-key", WithBaseURL(server.URL))
	_, err := client.Search(context.Background(), &SearchParams{Query: "test"})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("Expected *APIError, got %T", err)
	}
	if !apiErr.IsUnauthorized() {
		t.Error("Expected unauthorized error")
	}
	if !strings.Contains(apiErr.Message, "Invalid API key") {
		t.Errorf("Expected 'Invalid API key' in message, got %v", apiErr.Message)
	}
}

func TestMissingAPIKey(t *testing.T) {
	client := New("")
	_, err := client.Search(context.Background(), &SearchParams{Query: "test"})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("Expected *APIError, got %T", err)
	}
	if !apiErr.IsUnauthorized() {
		t.Error("Expected unauthorized error")
	}
}

func BenchmarkSearch(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"query": "test", "response_time": 0.5, "images": [], "results": []}`))
	}))
	defer server.Close()

	client := New("tvly-test-key", WithBaseURL(server.URL))
	ctx := context.Background()

	for b.Loop() {
		_, err := client.Search(ctx, &SearchParams{Query: "benchmark test"})
		if err != nil {
			b.Fatal(err)
		}
	}
}
