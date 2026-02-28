package sdk

import "context"

// Signup creates a new account and stores the returned tokens.
func (c *Client) Signup(ctx context.Context, req SignupRequest) (AuthTokens, error) {
	var tokens AuthTokens
	if err := c.unauthRequest(ctx, "POST", "/v1/auth/signup", req, &tokens); err != nil {
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

// Login authenticates with email and password and stores the returned tokens.
func (c *Client) Login(ctx context.Context, req LoginRequest) (AuthTokens, error) {
	var tokens AuthTokens
	if err := c.unauthRequest(ctx, "POST", "/v1/auth/login", req, &tokens); err != nil {
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

// RequestMagicLink sends a magic-link email.
func (c *Client) RequestMagicLink(ctx context.Context, req MagicLinkRequest) (MagicLinkResponse, error) {
	var resp MagicLinkResponse
	err := c.unauthRequest(ctx, "POST", "/v1/auth/magic-link", req, &resp)
	return resp, err
}

// VerifyMagicLink exchanges a magic-link token for auth tokens.
func (c *Client) VerifyMagicLink(ctx context.Context, req MagicLinkVerifyRequest) (AuthTokens, error) {
	var tokens AuthTokens
	if err := c.unauthRequest(ctx, "POST", "/v1/auth/magic-link/verify", req, &tokens); err != nil {
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
