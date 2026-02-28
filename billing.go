package sdk

import "context"

// CreateCheckout creates a Stripe checkout session.
func (c *Client) CreateCheckout(ctx context.Context, req CheckoutRequest) (CheckoutResponse, error) {
	var resp CheckoutResponse
	err := c.request(ctx, "POST", "/v1/billing/checkout", req, &resp)
	return resp, err
}

// CreatePortal creates a Stripe customer portal session.
func (c *Client) CreatePortal(ctx context.Context, req PortalRequest) (PortalResponse, error) {
	var resp PortalResponse
	err := c.request(ctx, "POST", "/v1/billing/portal", req, &resp)
	return resp, err
}

// GetBillingStatus returns the current billing status and limits.
func (c *Client) GetBillingStatus(ctx context.Context) (BillingStatus, error) {
	var status BillingStatus
	err := c.request(ctx, "GET", "/v1/billing/status", nil, &status)
	return status, err
}
