# Go Tavily Client

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/dl/)
[![Go Report Card](https://goreportcard.com/badge/github.com/solarhell/go-tavily?style=for-the-badge)](https://goreportcard.com/report/github.com/solarhell/go-tavily)
[![GoDoc](https://img.shields.io/badge/pkg.go.dev-reference-007d9c?style=for-the-badge&logo=go)](https://pkg.go.dev/github.com/solarhell/go-tavily)
[![CodeRabbit Pull Request Reviews](https://img.shields.io/coderabbit/prs/github/solarhell/go-tavily?utm_source=oss&utm_medium=github&utm_campaign=solarhell%2Fgo-tavily&labelColor=171717&color=FF570A&link=https%3A%2F%2Fcoderabbit.ai&label=CodeRabbit+Reviews)](https://coderabbit.ai)

A thin, type-safe Go client for the [Tavily API](https://docs.tavily.com). Built for Go 1.26+.

## Installation

```bash
go get github.com/solarhell/go-tavily
```

## Quick Start

```go
client := tavily.New("tvly-your-api-key")

resp, err := client.Search(ctx, &tavily.SearchParams{
    Query: "Go programming language",
})
```

## Search

```go
answerMode := tavily.IncludeAnswerAdvanced

resp, err := client.Search(ctx, &tavily.SearchParams{
    Query:         "AI news",
    SearchDepth:   tavily.SearchDepthAdvanced,
    Topic:         tavily.TopicNews,
    TimeRange:     tavily.TimeRangeWeek,
    MaxResults:    10,
    IncludeAnswer: &answerMode,
})
```

Zero-value fields are omitted from the request — the API uses its own server-side defaults.

## Extract

```go
resp, err := client.Extract(ctx, &tavily.ExtractParams{
    URLs:   []string{"https://example.com"},
    Format: tavily.FormatMarkdown,
})
```

## Client Options

```go
client := tavily.New("tvly-your-api-key",
    tavily.WithBaseURL("https://custom.api.com"),
    tavily.WithHTTPClient(&http.Client{Timeout: 45 * time.Second}),
)
```

If the API key is empty, it reads from the `TAVILY_API_KEY` environment variable.

## Error Handling

```go
resp, err := client.Search(ctx, &tavily.SearchParams{Query: "test"})
if err != nil {
    var apiErr *tavily.APIError
    if errors.As(err, &apiErr) {
        switch {
        case apiErr.IsUnauthorized():
            // invalid API key (401)
        case apiErr.IsRateLimit():
            // rate limited (429)
        case apiErr.IsPlanLimitExceeded():
            // plan usage limit (432)
        case apiErr.IsPayGoLimitExceeded():
            // pay-as-you-go limit (433)
        case apiErr.IsBadRequest():
            // invalid parameters (400)
        }
    }
}
```

## Testing

```bash
go test -v -race ./...
```

## Links

- [Tavily API Documentation](https://docs.tavily.com)
