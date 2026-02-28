package sdk

import (
	"context"
	"net/url"
)

// ListDestinations returns all destinations for a roost.
func (c *Client) ListDestinations(ctx context.Context, roostID string) ([]Destination, error) {
	var dests []Destination
	err := c.request(ctx, "GET", "/v1/roosts/"+url.PathEscape(roostID)+"/destinations", nil, &dests)
	return dests, err
}

// GetDestination returns a destination by ID.
func (c *Client) GetDestination(ctx context.Context, destinationID string) (Destination, error) {
	var dest Destination
	err := c.request(ctx, "GET", "/v1/destinations/"+url.PathEscape(destinationID), nil, &dest)
	return dest, err
}

// CreateDestination adds a new destination to a roost.
func (c *Client) CreateDestination(ctx context.Context, roostID string, req CreateDestination) (Destination, error) {
	var dest Destination
	err := c.request(ctx, "POST", "/v1/roosts/"+url.PathEscape(roostID)+"/destinations", req, &dest)
	return dest, err
}

// PauseDestination pauses or unpauses delivery to a destination.
func (c *Client) PauseDestination(ctx context.Context, destinationID string, isPaused bool) (PauseResponse, error) {
	var resp PauseResponse
	err := c.request(ctx, "PATCH", "/v1/destinations/"+url.PathEscape(destinationID)+"/pause", PauseInput{IsPaused: isPaused}, &resp)
	return resp, err
}

// DeleteDestination permanently deletes a destination.
func (c *Client) DeleteDestination(ctx context.Context, destinationID string) error {
	return c.request(ctx, "DELETE", "/v1/destinations/"+url.PathEscape(destinationID), nil, nil)
}
