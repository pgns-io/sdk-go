package sdk

import (
	"context"
	"net/url"
)

// ListApiKeys returns all API keys for the authenticated user.
func (c *Client) ListApiKeys(ctx context.Context) ([]ApiKeyResponse, error) {
	var keys []ApiKeyResponse
	err := c.request(ctx, "GET", "/v1/api-keys", nil, &keys)
	return keys, err
}

// GetApiKey returns an API key by ID (does not return the full key value).
func (c *Client) GetApiKey(ctx context.Context, keyID string) (ApiKeyResponse, error) {
	var key ApiKeyResponse
	err := c.request(ctx, "GET", "/v1/api-keys/"+url.PathEscape(keyID), nil, &key)
	return key, err
}

// CreateApiKey creates a new API key. The full key is only in this response.
func (c *Client) CreateApiKey(ctx context.Context, req *CreateApiKeyRequest) (ApiKeyCreatedResponse, error) {
	var body any
	if req != nil {
		body = req
	} else {
		body = struct{}{}
	}
	var result ApiKeyCreatedResponse
	err := c.request(ctx, "POST", "/v1/api-keys", body, &result)
	return result, err
}

// UpdateApiKey renames an API key.
func (c *Client) UpdateApiKey(ctx context.Context, keyID string, req UpdateApiKeyRequest) (ApiKeyResponse, error) {
	var key ApiKeyResponse
	err := c.request(ctx, "PATCH", "/v1/api-keys/"+url.PathEscape(keyID), req, &key)
	return key, err
}

// DeleteApiKey permanently revokes and deletes an API key.
func (c *Client) DeleteApiKey(ctx context.Context, keyID string) error {
	return c.request(ctx, "DELETE", "/v1/api-keys/"+url.PathEscape(keyID), nil, nil)
}
