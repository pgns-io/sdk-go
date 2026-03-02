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
