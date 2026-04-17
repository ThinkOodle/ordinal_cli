// Package client provides an HTTP client for the Ordinal API with
// authentication, rate limiting, and error handling.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	// DefaultBaseURL is the base URL for the Ordinal API.
	DefaultBaseURL = "https://app.tryordinal.com/api/v1"

	// DefaultTimeout is the default HTTP client timeout.
	DefaultTimeout = 30 * time.Second

	// maxRetries is the maximum number of retries for rate-limited requests.
	maxRetries = 5
)

// initialBackoffValue is the starting backoff duration for retries. It is a
// var so tests can shrink it to make retry timing deterministic.
var initialBackoffValue = 2 * time.Second

// Client is an HTTP client for the Ordinal API.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	verbose    bool
}

// Option configures a Client.
type Option func(*Client)

// WithBaseURL sets a custom base URL (useful for testing).
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// WithVerbose enables verbose output for debugging.
func WithVerbose(v bool) Option {
	return func(c *Client) {
		c.verbose = v
	}
}

// New creates a new Ordinal API client.
func New(apiKey string, opts ...Option) *Client {
	c := &Client{
		baseURL: DefaultBaseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// APIError represents an error response from the Ordinal API.
type APIError struct {
	StatusCode int
	Message    string
	Code       string
	Body       string
}

// Error returns the error message.
func (e *APIError) Error() string {
	if e.Message != "" {
		if e.Code != "" {
			return fmt.Sprintf("api error (status %d, code %s): %s", e.StatusCode, e.Code, e.Message)
		}
		return fmt.Sprintf("api error (status %d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("api error (status %d): %s", e.StatusCode, e.Body)
}

// apiErrorResponse is the JSON structure of Ordinal API error responses.
type apiErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"error"`
	Message string `json:"message"`
}

// Get performs a GET request to the given path with optional query parameters.
func (c *Client) Get(path string, query url.Values) ([]byte, error) {
	u := c.baseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	return c.do(req)
}

// Post performs a POST request to the given path with a JSON body.
func (c *Client) Post(path string, body interface{}) ([]byte, error) {
	return c.doJSON(http.MethodPost, path, body)
}

// Put performs a PUT request to the given path with a JSON body.
func (c *Client) Put(path string, body interface{}) ([]byte, error) {
	return c.doJSON(http.MethodPut, path, body)
}

// Patch performs a PATCH request to the given path with a JSON body.
func (c *Client) Patch(path string, body interface{}) ([]byte, error) {
	return c.doJSON(http.MethodPatch, path, body)
}

// Delete performs a DELETE request to the given path.
func (c *Client) Delete(path string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodDelete, c.baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	return c.do(req)
}

// doJSON marshals the body to JSON and performs the request.
func (c *Client) doJSON(method, path string, body interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, fmt.Errorf("encoding request body: %w", err)
		}
	}

	req, err := http.NewRequest(method, c.baseURL+path, &buf)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	return c.do(req)
}

// isIdempotent reports whether a request method is safe to retry automatically
// after a transport-level failure. Non-idempotent methods (POST/PATCH) are not
// retried on transport errors because the server may have processed the first
// request before the response was lost, and a retry would duplicate the side
// effect.
func isIdempotent(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodPut, http.MethodDelete, http.MethodOptions:
		return true
	}
	return false
}

// do executes the request with authentication and retry logic. Retries are
// restricted to:
//   - 429 Too Many Requests responses (always safe; the request was rejected
//     before it could take effect).
//   - Transport errors on idempotent methods.
//
// Non-idempotent methods (POST, PATCH) are NOT retried on transport errors:
// the server may have already processed the request, and retrying would
// duplicate side effects like creating or modifying resources.
func (c *Client) do(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")

	if c.verbose {
		fmt.Printf(">> %s %s\n", req.Method, req.URL.String())
	}

	idempotent := isIdempotent(req.Method)
	var lastErr error
	backoff := initialBackoffValue

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			if c.verbose {
				fmt.Printf(">> retry %d/%d after %s\n", attempt, maxRetries, backoff)
			}
			time.Sleep(backoff)
			backoff *= 2

			if req.GetBody != nil {
				body, err := req.GetBody()
				if err != nil {
					return nil, fmt.Errorf("re-reading request body for retry: %w", err)
				}
				req.Body = body
			}
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			if !idempotent {
				return nil, fmt.Errorf("executing request: %w", err)
			}
			lastErr = fmt.Errorf("executing request: %w", err)
			continue
		}

		data, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			if !idempotent {
				return nil, fmt.Errorf("reading response body: %w", err)
			}
			lastErr = fmt.Errorf("reading response body: %w", err)
			continue
		}

		if c.verbose {
			fmt.Printf("<< %d %s\n", resp.StatusCode, http.StatusText(resp.StatusCode))
			if len(data) > 0 {
				fmt.Printf("<< %s\n", string(data))
			}
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
				if seconds, err := strconv.Atoi(retryAfter); err == nil {
					backoff = time.Duration(seconds) * time.Second
				}
			}
			lastErr = &APIError{
				StatusCode: resp.StatusCode,
				Body:       string(data),
				Message:    "rate limit exceeded",
			}
			continue
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			apiErr := &APIError{
				StatusCode: resp.StatusCode,
				Body:       string(data),
			}
			var errResp apiErrorResponse
			if json.Unmarshal(data, &errResp) == nil {
				if errResp.Error.Message != "" {
					apiErr.Message = errResp.Error.Message
					apiErr.Code = errResp.Error.Code
				} else if errResp.Message != "" {
					apiErr.Message = errResp.Message
				}
			}
			return nil, apiErr
		}

		return data, nil
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", maxRetries, lastErr)
}
