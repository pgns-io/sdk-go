package sdk

import "context"

// Refresh refreshes the access token using the httpOnly refresh cookie.
func (c *Client) Refresh(ctx context.Context) (AuthTokens, error) {
	var tokens AuthTokens
	if err := c.unauthRequest(ctx, "POST", "/v1/auth/refresh", nil, &tokens); err != nil {
		return AuthTokens{}, err
	}
	c.mu.Lock()
	c.accessToken = tokens.AccessToken
	c.mu.Unlock()
	if c.onTokenRefresh != nil {
		c.onTokenRefresh(tokens)
	}
	return tokens, nil
}

// Logout revokes the refresh token and clears stored credentials.
func (c *Client) Logout(ctx context.Context) error {
	err := c.request(ctx, "POST", "/v1/auth/logout", nil, nil)
	if err != nil {
		return err
	}
	c.mu.Lock()
	c.accessToken = ""
	c.mu.Unlock()
	return nil
}
