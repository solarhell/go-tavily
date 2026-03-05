package tavily

import (
	"context"
	"fmt"
)

// SearchDepth controls the depth of the search.
type SearchDepth string

const (
	SearchDepthBasic    SearchDepth = "basic"
	SearchDepthAdvanced SearchDepth = "advanced"
	SearchDepthFast     SearchDepth = "fast"
	SearchDepthUltraFast SearchDepth = "ultra-fast"
)

// Topic represents the topic category for search.
type Topic string

const (
	TopicGeneral Topic = "general"
	TopicNews    Topic = "news"
	TopicFinance Topic = "finance"
)

// TimeRange represents the time range filter for search results.
type TimeRange string

const (
	TimeRangeDay   TimeRange = "day"
	TimeRangeWeek  TimeRange = "week"
	TimeRangeMonth TimeRange = "month"
	TimeRangeYear  TimeRange = "year"
)

// IncludeAnswer controls the LLM-generated answer in search results.
type IncludeAnswer string

const (
	IncludeAnswerBasic    IncludeAnswer = "basic"
	IncludeAnswerAdvanced IncludeAnswer = "advanced"
)

// IncludeRawContent controls cleaned content format in search results.
type IncludeRawContent string

const (
	IncludeRawContentText     IncludeRawContent = "text"
	IncludeRawContentMarkdown IncludeRawContent = "markdown"
)

// SearchParams is the request body for POST /search.
// Zero-value fields are omitted; the API uses server-side defaults.
type SearchParams struct {
	Query                    string             `json:"query"`
	SearchDepth              SearchDepth        `json:"search_depth,omitzero"`
	Topic                    Topic              `json:"topic,omitzero"`
	TimeRange                TimeRange          `json:"time_range,omitzero"`
	StartDate                string             `json:"start_date,omitzero"`
	EndDate                  string             `json:"end_date,omitzero"`
	MaxResults               uint64             `json:"max_results,omitzero"`
	ChunksPerSource          uint64             `json:"chunks_per_source,omitzero"`
	IncludeDomains           []string           `json:"include_domains,omitzero"`
	ExcludeDomains           []string           `json:"exclude_domains,omitzero"`
	IncludeAnswer            *IncludeAnswer     `json:"include_answer,omitzero"`
	IncludeRawContent        *IncludeRawContent `json:"include_raw_content,omitzero"`
	IncludeFavicon           *bool              `json:"include_favicon,omitzero"`
	Country                  string             `json:"country,omitzero"`
	AutoParameters           *bool              `json:"auto_parameters,omitzero"`
	ExactMatch               *bool              `json:"exact_match,omitzero"`
	IncludeUsage             *bool              `json:"include_usage,omitzero"`
}

// SearchResponse is the response from POST /search.
type SearchResponse struct {
	Query        string         `json:"query"`
	Answer       string         `json:"answer,omitzero"`
	ResponseTime float64        `json:"response_time"`
	Results      []SearchResult `json:"results"`
	Usage        *Usage         `json:"usage,omitzero"`
	RequestID    string         `json:"request_id,omitzero"`
}

// SearchResult represents a single search result.
type SearchResult struct {
	Title         string  `json:"title"`
	URL           string  `json:"url"`
	Content       string  `json:"content"`
	RawContent    string  `json:"raw_content,omitzero"`
	Score         float64 `json:"score"`
	PublishedDate string  `json:"published_date,omitzero"`
	Favicon       string  `json:"favicon,omitzero"`
}

// Usage represents credit usage information.
type Usage struct {
	Credits uint64 `json:"credits"`
}

// Search performs a web search via the Tavily API.
func (c *Client) Search(ctx context.Context, params *SearchParams) (*SearchResponse, error) {
	var resp SearchResponse
	if err := c.do(ctx, "/search", params, &resp); err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}
	return &resp, nil
}
