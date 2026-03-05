package tavily

// APIError represents an error response from the Tavily API.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return e.Message
}

// IsBadRequest returns true if the error is due to invalid parameters (400).
func (e *APIError) IsBadRequest() bool {
	return e.StatusCode == 400
}

// IsUnauthorized returns true if the error is due to an invalid API key (401).
func (e *APIError) IsUnauthorized() bool {
	return e.StatusCode == 401
}

// IsRateLimit returns true if the error is due to rate limiting (429).
func (e *APIError) IsRateLimit() bool {
	return e.StatusCode == 429
}

// IsPlanLimitExceeded returns true if the plan usage limit is exceeded (432).
func (e *APIError) IsPlanLimitExceeded() bool {
	return e.StatusCode == 432
}

// IsPayGoLimitExceeded returns true if the pay-as-you-go limit is exceeded (433).
func (e *APIError) IsPayGoLimitExceeded() bool {
	return e.StatusCode == 433
}
