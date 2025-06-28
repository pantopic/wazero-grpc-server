# Wazero gRPC

A [wazero](https://pkg.go.dev/github.com/tetratelabs/wazero) host module, ABI and guest SDK providing [gRPC](https://grpc.io/) for WASI modules.

## Host Module

<!-- [![Go Reference](https://godoc.org/github.com/pantopic/wazero-grpc/host?status.svg)](https://godoc.org/github.com/pantopic/wazero-grpc/host) -->
<!-- [![Go Report Card](https://goreportcard.com/badge/github.com/pantopic/wazero-grpc/host)](https://goreportcard.com/report/github.com/pantopic/wazero-grpc/host) -->
<!-- [![Go Coverage](https://github.com/pantopic/wazero-grpc/wiki/host/coverage.svg)](https://raw.githack.com/wiki/pantopic/wazero-grpc/host/coverage.html) -->

First register the host module with the runtime

```go
import (
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"

	"github.com/pantopic/wazero-grpc/host"
)

func main() {
	ctx := context.Background()
	r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig())
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	module := wazero_grpc.New()
	module.Register(ctx, r)

	// ...
}
```

## Guest SDK (Go)

<!-- [![Go Reference](https://godoc.org/github.com/pantopic/wazero-grpc/grpc-go?status.svg)](https://godoc.org/github.com/pantopic/wazero-grpc/grpc-go) -->
<!-- [![Go Report Card](https://goreportcard.com/badge/github.com/pantopic/wazero-grpc/grpc-go)](https://goreportcard.com/report/github.com/pantopic/wazero-grpc/grpc-go) -->

Then you can import the guest SDK into your WASI module to receive gRPC requests from WASM.

```go
package main

import (
	"github.com/pantopic/wazero-grpc/grpc-server-go"
)

func main() {}

// ...
```

The [guest SDK](https://pkg.go.dev/github.com/pantopic/wazero-grpc/grpc-server-go) has no dependencies outside the Go std lib.
The guest SDK is serialization format agnostic in order to provide users with more control over performance and binary size.

See examples for serialization techniques:

- [test-server](/test-server) uses `protoc-go-gen` for serialization, resulting in a large WASM binary (789k)
- [test-server-easy](/test-server-easy) uses `easyproto` for manually-generated serialization, resulting in a small WASM binary (146k)
- [test-server-lite](/test-server-lite) uses `protobuf-go-lite` for auto-generated serialization also resulting in a small WASM binary (145k)

Any of these approaches and others will work for protobuf serialization.

## Roadmap

This project is in alpha. Breaking API changes should be expected until Beta.

- `v0.0.x` - Alpha
  - [ ] Stabilize API
- `v0.x.x` - Beta
  - [ ] Finalize API
  - [ ] Test in production
- `v1.x.x` - General Availability
  - [ ] Proven long term stability in production
