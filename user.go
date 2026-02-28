package sdk

import "context"

// GetMe returns the authenticated user's profile.
func (c *Client) GetMe(ctx context.Context) (User, error) {
	var user User
	err := c.request(ctx, "GET", "/v1/me", nil, &user)
	return user, err
}

// UpdateMe updates the authenticated user's profile.
func (c *Client) UpdateMe(ctx context.Context, req UpdateProfileRequest) (User, error) {
	var user User
	err := c.request(ctx, "PATCH", "/v1/me", req, &user)
	return user, err
}

// GetStats returns aggregated dashboard statistics.
func (c *Client) GetStats(ctx context.Context) (DashboardStats, error) {
	var stats DashboardStats
	err := c.request(ctx, "GET", "/v1/stats", nil, &stats)
	return stats, err
}
