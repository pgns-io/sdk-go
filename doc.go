// Package sdk provides a Go client for the pgns webhook relay API.
//
// The client supports two authentication modes:
//   - API key — pass [WithAPIKey] for server-side usage.
//   - JWT — pass [WithAccessToken]. Expired tokens are refreshed automatically on 401.
//
// # Quick start
//
//	client := sdk.NewClient("https://api.pgns.io", sdk.WithAPIKey("pk_live_..."))
//	roosts, err := client.ListRoosts(ctx)
package sdk
