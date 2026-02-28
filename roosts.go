package sdk

import (
	"context"
	"net/url"
)

// ListRoosts returns all roosts for the authenticated user.
func (c *Client) ListRoosts(ctx context.Context) ([]Roost, error) {
	var roosts []Roost
	err := c.request(ctx, "GET", "/v1/roosts", nil, &roosts)
	return roosts, err
}

// GetRoost returns a roost by ID.
func (c *Client) GetRoost(ctx context.Context, roostID string) (Roost, error) {
	var roost Roost
	err := c.request(ctx, "GET", "/v1/roosts/"+url.PathEscape(roostID), nil, &roost)
	return roost, err
}

// CreateRoost creates a new roost (webhook endpoint).
func (c *Client) CreateRoost(ctx context.Context, req CreateRoost) (Roost, error) {
	var roost Roost
	err := c.request(ctx, "POST", "/v1/roosts", req, &roost)
	return roost, err
}

// UpdateRoost updates a roost's name, description, secret, or active state.
func (c *Client) UpdateRoost(ctx context.Context, roostID string, req UpdateRoost) (Roost, error) {
	var roost Roost
	err := c.request(ctx, "PATCH", "/v1/roosts/"+url.PathEscape(roostID), req, &roost)
	return roost, err
}

// DeleteRoost deletes a roost and all its destinations.
func (c *Client) DeleteRoost(ctx context.Context, roostID string) error {
	return c.request(ctx, "DELETE", "/v1/roosts/"+url.PathEscape(roostID), nil, nil)
}
