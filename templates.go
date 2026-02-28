package sdk

import (
	"context"
	"net/url"
)

// ListTemplates returns all templates for the authenticated user.
func (c *Client) ListTemplates(ctx context.Context) ([]Template, error) {
	var templates []Template
	err := c.request(ctx, "GET", "/v1/templates", nil, &templates)
	return templates, err
}

// GetTemplate returns a template by ID.
func (c *Client) GetTemplate(ctx context.Context, templateID string) (Template, error) {
	var tmpl Template
	err := c.request(ctx, "GET", "/v1/templates/"+url.PathEscape(templateID), nil, &tmpl)
	return tmpl, err
}

// CreateTemplate creates a new template.
func (c *Client) CreateTemplate(ctx context.Context, req CreateTemplate) (Template, error) {
	var tmpl Template
	err := c.request(ctx, "POST", "/v1/templates", req, &tmpl)
	return tmpl, err
}

// UpdateTemplate updates a template.
func (c *Client) UpdateTemplate(ctx context.Context, templateID string, req UpdateTemplate) (Template, error) {
	var tmpl Template
	err := c.request(ctx, "PATCH", "/v1/templates/"+url.PathEscape(templateID), req, &tmpl)
	return tmpl, err
}

// DeleteTemplate deletes a template.
func (c *Client) DeleteTemplate(ctx context.Context, templateID string) error {
	return c.request(ctx, "DELETE", "/v1/templates/"+url.PathEscape(templateID), nil, nil)
}

// PreviewTemplate renders a template with a pigeon's data.
func (c *Client) PreviewTemplate(ctx context.Context, req PreviewTemplateRequest) (PreviewTemplateResponse, error) {
	var resp PreviewTemplateResponse
	err := c.request(ctx, "POST", "/v1/templates/preview", req, &resp)
	return resp, err
}
