package sdk

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// ListPigeonsOptions configures the ListPigeons call.
type ListPigeonsOptions struct {
	RoostID string
	Limit   int
	Cursor  string
}

// ListPigeons returns a paginated list of pigeons.
func (c *Client) ListPigeons(ctx context.Context, opts *ListPigeonsOptions) (PaginatedPigeons, error) {
	var params []string
	if opts != nil {
		if opts.RoostID != "" {
			params = append(params, "roost_id="+url.QueryEscape(opts.RoostID))
		}
		if opts.Limit > 0 {
			params = append(params, fmt.Sprintf("limit=%d", opts.Limit))
		}
		if opts.Cursor != "" {
			params = append(params, "cursor="+url.QueryEscape(opts.Cursor))
		}
	}
	path := "/v1/pigeons"
	if len(params) > 0 {
		path += "?" + strings.Join(params, "&")
	}

	var result PaginatedPigeons
	err := c.request(ctx, "GET", path, nil, &result)
	return result, err
}

// GetPigeon returns a single pigeon by ID.
func (c *Client) GetPigeon(ctx context.Context, pigeonID string) (Pigeon, error) {
	var pigeon Pigeon
	err := c.request(ctx, "GET", "/v1/pigeons/"+url.PathEscape(pigeonID), nil, &pigeon)
	return pigeon, err
}

// ListDeliveriesOptions configures the GetPigeonDeliveries call.
type ListDeliveriesOptions struct {
	Limit  int
	Cursor string
}

// GetPigeonDeliveries returns delivery attempts for a pigeon.
func (c *Client) GetPigeonDeliveries(ctx context.Context, pigeonID string, opts *ListDeliveriesOptions) (PaginatedDeliveryAttempts, error) {
	var params []string
	if opts != nil {
		if opts.Limit > 0 {
			params = append(params, fmt.Sprintf("limit=%d", opts.Limit))
		}
		if opts.Cursor != "" {
			params = append(params, "cursor="+url.QueryEscape(opts.Cursor))
		}
	}
	path := "/v1/pigeons/" + url.PathEscape(pigeonID) + "/deliveries"
	if len(params) > 0 {
		path += "?" + strings.Join(params, "&")
	}

	var result PaginatedDeliveryAttempts
	err := c.request(ctx, "GET", path, nil, &result)
	return result, err
}

// ReplayPigeon re-delivers a pigeon to all active destinations.
func (c *Client) ReplayPigeon(ctx context.Context, pigeonID string) (ReplayResponse, error) {
	var result ReplayResponse
	err := c.request(ctx, "POST", "/v1/pigeons/"+url.PathEscape(pigeonID)+"/replay", nil, &result)
	return result, err
}
