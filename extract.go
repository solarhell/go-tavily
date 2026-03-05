package tavily

import (
	"context"
	"fmt"
)

// ExtractDepth controls the depth of content extraction.
type ExtractDepth string

const (
	ExtractDepthBasic    ExtractDepth = "basic"
	ExtractDepthAdvanced ExtractDepth = "advanced"
)

// Format represents the output format for extracted content.
type Format string

const (
	FormatMarkdown Format = "markdown"
	FormatText     Format = "text"
)

// ExtractParams is the request body for POST /extract.
// Zero-value fields are omitted; the API uses server-side defaults.
type ExtractParams struct {
	URLs            []string     `json:"urls"`
	Query           string       `json:"query,omitzero"`
	ChunksPerSource uint64       `json:"chunks_per_source,omitzero"`
	ExtractDepth    ExtractDepth `json:"extract_depth,omitzero"`
	IncludeImages   *bool        `json:"include_images,omitzero"`
	IncludeFavicon  *bool        `json:"include_favicon,omitzero"`
	Format          Format       `json:"format,omitzero"`
	Timeout         float64      `json:"timeout,omitzero"`
	IncludeUsage    *bool        `json:"include_usage,omitzero"`
}

// ExtractResponse is the response from POST /extract.
type ExtractResponse struct {
	ResponseTime  float64               `json:"response_time"`
	Results       []ExtractResult       `json:"results"`
	FailedResults []ExtractFailedResult `json:"failed_results"`
	Usage         *Usage                `json:"usage,omitzero"`
	RequestID     string                `json:"request_id,omitzero"`
}

// ExtractResult represents a successful content extraction.
type ExtractResult struct {
	URL        string   `json:"url"`
	RawContent string   `json:"raw_content"`
	Images     []string `json:"images,omitzero"`
	Favicon    string   `json:"favicon,omitzero"`
}

// ExtractFailedResult represents a failed content extraction.
type ExtractFailedResult struct {
	URL   string `json:"url"`
	Error string `json:"error"`
}

// Extract extracts content from one or more URLs via the Tavily API.
func (c *Client) Extract(ctx context.Context, params *ExtractParams) (*ExtractResponse, error) {
	var resp ExtractResponse
	if err := c.do(ctx, "/extract", params, &resp); err != nil {
		return nil, fmt.Errorf("extract failed: %w", err)
	}
	return &resp, nil
}
