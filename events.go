package sdk

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultRetryDelay = 3 * time.Second

// EventOption configures [Client.ListenEvents].
type EventOption func(*eventConfig)

type eventConfig struct {
	roostID string
	onError func(error)
}

// WithRoostID restricts the event stream to a single roost.
func WithRoostID(id string) EventOption {
	return func(cfg *eventConfig) { cfg.roostID = id }
}

// WithErrorHandler registers a callback for connection errors.
func WithErrorHandler(fn func(error)) EventOption {
	return func(cfg *eventConfig) { cfg.onError = fn }
}

// ListenEvents connects to the SSE event stream and calls onEvent for each
// event received. It automatically reconnects on failure with a 3-second
// delay. The function blocks until ctx is cancelled.
func (c *Client) ListenEvents(ctx context.Context, onEvent func(data string), opts ...EventOption) error {
	var cfg eventConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	eventURL := c.baseURL + "/v1/events"
	if cfg.roostID != "" {
		eventURL += "?roost_id=" + url.QueryEscape(cfg.roostID)
	}

	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		err := c.streamEvents(ctx, eventURL, onEvent)
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err != nil && cfg.onError != nil {
			cfg.onError(err)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(defaultRetryDelay):
		}
	}
}

func (c *Client) streamEvents(ctx context.Context, eventURL string, onEvent func(string)) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, eventURL, nil)
	if err != nil {
		return fmt.Errorf("pgns: create SSE request: %w", err)
	}

	req.Header.Set("Accept", "text/event-stream")
	auth := c.authHeader()
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("pgns: SSE connect: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return &PigeonsError{
			Message:    fmt.Sprintf("SSE connect failed: %s", resp.Status),
			StatusCode: resp.StatusCode,
		}
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data:") {
			data := strings.TrimSpace(line[5:])
			onEvent(data)
		}
	}
	return scanner.Err()
}
