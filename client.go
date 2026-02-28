package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"
)

// Client is a Go client for the pgns webhook relay API.
type Client struct {
	baseURL        string
	apiKey         string
	accessToken    string
	httpClient     *http.Client
	onTokenRefresh func(AuthTokens)
	mu             sync.Mutex
}

// Option configures a [Client].
type Option func(*Client)

// WithAPIKey sets the API key for authentication.
func WithAPIKey(key string) Option {
	return func(c *Client) { c.apiKey = key }
}

// WithAccessToken sets the JWT access token for authentication.
func WithAccessToken(token string) Option {
	return func(c *Client) { c.accessToken = token }
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

// WithTokenRefreshHandler registers a callback invoked after a token refresh.
func WithTokenRefreshHandler(fn func(AuthTokens)) Option {
	return func(c *Client) { c.onTokenRefresh = fn }
}

// NewClient creates a new pgns API client.
func NewClient(baseURL string, opts ...Option) *Client {
	jar, _ := cookiejar.New(nil)
	c := &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{Jar: jar},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// SetAccessToken replaces the current JWT access token.
func (c *Client) SetAccessToken(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.accessToken = token
}

// SetAPIKey replaces the current API key.
func (c *Client) SetAPIKey(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.apiKey = key
}

// doRequest executes an HTTP request with the given auth header.
func (c *Client) doRequest(ctx context.Context, method, path string, bodyBytes []byte, authHeader string) (*http.Response, error) {
	var body io.Reader
	if bodyBytes != nil {
		body = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("pgns: create request: %w", err)
	}

	if bodyBytes != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	return c.httpClient.Do(req)
}

// authHeader returns the current Authorization header value.
func (c *Client) authHeader() string {
	if c.apiKey != "" {
		return "Bearer " + c.apiKey
	}
	if c.accessToken != "" {
		return "Bearer " + c.accessToken
	}
	return ""
}

// refreshToken performs a token refresh and updates the stored access token.
func (c *Client) refreshToken(ctx context.Context) (AuthTokens, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	resp, err := c.doRequest(ctx, http.MethodPost, "/v1/auth/refresh", nil, "")
	if err != nil {
		return AuthTokens{}, err
	}
	defer func() { _ = resp.Body.Close() }()

	var tokens AuthTokens
	if err := handleResponse(resp, &tokens); err != nil {
		return AuthTokens{}, err
	}

	c.accessToken = tokens.AccessToken
	if c.onTokenRefresh != nil {
		c.onTokenRefresh(tokens)
	}
	return tokens, nil
}

// request performs an authenticated request with automatic 401 retry.
func (c *Client) request(ctx context.Context, method, path string, reqBody any, result any) error {
	var bodyBytes []byte
	if reqBody != nil {
		var err error
		bodyBytes, err = json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("pgns: marshal request: %w", err)
		}
	}

	resp, err := c.doRequest(ctx, method, path, bodyBytes, c.authHeader())
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	// Handle 401 with token refresh (JWT mode only).
	if resp.StatusCode == http.StatusUnauthorized && c.apiKey == "" {
		_ = resp.Body.Close()

		tokens, refreshErr := c.refreshToken(ctx)
		if refreshErr != nil {
			c.mu.Lock()
			c.accessToken = ""
			c.mu.Unlock()
			return &PigeonsError{Message: "session expired", StatusCode: 401}
		}

		resp, err = c.doRequest(ctx, method, path, bodyBytes, "Bearer "+tokens.AccessToken)
		if err != nil {
			return err
		}
		defer func() { _ = resp.Body.Close() }()
	}

	return handleResponse(resp, result)
}

// unauthRequest performs a request without auth headers (for login/signup/etc.).
func (c *Client) unauthRequest(ctx context.Context, method, path string, reqBody any, result any) error {
	var bodyBytes []byte
	if reqBody != nil {
		var err error
		bodyBytes, err = json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("pgns: marshal request: %w", err)
		}
	}

	resp, err := c.doRequest(ctx, method, path, bodyBytes, "")
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	return handleResponse(resp, result)
}

// handleResponse reads the response body and unmarshals it or returns an error.
func handleResponse(resp *http.Response, result any) error {
	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("pgns: read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var apiErr struct {
			Error string `json:"error"`
		}
		if json.Unmarshal(respBody, &apiErr) == nil && apiErr.Error != "" {
			return &PigeonsError{Message: apiErr.Error, StatusCode: resp.StatusCode}
		}
		return &PigeonsError{Message: resp.Status, StatusCode: resp.StatusCode}
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("pgns: unmarshal response: %w", err)
		}
	}
	return nil
}
