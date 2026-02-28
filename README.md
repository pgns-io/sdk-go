# pgns Go SDK

Go client library for the [pgns](https://pgns.io) webhook relay API.

## Installation

```bash
go get github.com/pgns-io/sdk-go
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"

	pgns "github.com/pgns-io/sdk-go"
)

func main() {
	client := pgns.NewClient("your-api-key")

	// List roosts
	roosts, err := client.Roosts.List(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(roosts)

	// Send a pigeon
	err = client.Pigeons.Send(context.Background(), "rst_abc123", map[string]any{
		"event": "user.created",
		"data":  map[string]any{"id": 1},
	})
	if err != nil {
		panic(err)
	}
}
```

## Documentation

Full documentation is available at [docs.pgns.io/sdks/go](https://docs.pgns.io/sdks/go).

API reference: [pkg.go.dev/github.com/pgns-io/sdk-go](https://pkg.go.dev/github.com/pgns-io/sdk-go)

## License

MIT
