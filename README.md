# Wazero gRPC Server

A [wazero](https://pkg.go.dev/github.com/tetratelabs/wazero) host module, ABI and guest SDK enabling WASI modules to register as [gRPC](https://grpc.io/) services.

## Host Module

[![Go Reference](https://godoc.org/github.com/pantopic/wazero-grpc/host?status.svg)](https://godoc.org/github.com/pantopic/wazero-grpc/host)
[![Go Report Card](https://goreportcard.com/badge/github.com/pantopic/wazero-grpc/host)](https://goreportcard.com/report/github.com/pantopic/wazero-grpc/host)
[![Go Coverage](https://github.com/pantopic/wazero-grpc/wiki/host/coverage.svg)](https://raw.githack.com/wiki/pantopic/wazero-grpc/host/coverage.html)

First register the host module with the runtime

```go
import (
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"

	"github.com/pantopic/wazero-grpc-server/host"
)

func main() {
	ctx := context.Background()
	r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig())
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	module := wazero_grpc_server.New()
	module.Register(ctx, r)

	// ...
}
```

## Guest SDK (Go)

[![Go Reference](https://godoc.org/github.com/pantopic/wazero-grpc/grpc-go?status.svg)](https://godoc.org/github.com/pantopic/wazero-grpc/grpc-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/pantopic/wazero-grpc/grpc-go)](https://goreportcard.com/report/github.com/pantopic/wazero-grpc/grpc-go)

Then you can import the guest SDK into your WASI module to receive gRPC requests from WASM.

```go
package main

import (
	"google.golang.org/protobuf/proto"

	"github.com/pantopic/wazero-grpc-server/grpc-server-go"
	"github.com/pantopic/wazero-grpc-server/test/pb"
)

func main() {
	s := grpc_server.NewService(`test.TestService`)
	s.AddMethod(`Test`, protoWrap(test, &pb.TestRequest{}))
	s.AddMethod(`Retest`, protoWrap(retest, &pb.RetestRequest{}))
}

func test(req *pb.TestRequest) (res *pb.TestResponse, err error) {
	return &pb.TestResponse{
		Bar: req.Foo,
	}, nil
}

func retest(req *pb.RetestRequest) (res *pb.RetestResponse, err error) {
	return &pb.RetestResponse{
		Foo: req.Bar,
	}, nil
}

// ...
```

The [guest SDK](https://pkg.go.dev/github.com/pantopic/wazero-grpc/grpc-server-go) has no dependencies outside the Go std lib.
The guest SDK is serialization agnostic in order to provide users with more control over performance, compile time and binary size.

See examples for protobuf message serialization options:

- [test](/test) - `protoc-go-gen` auto-generated serialization (`792kb` binary, `15s` build time)
- [test-easy](/test-easy) - `easyproto` manually-generated serialization (`147kb` binary, `3s` build time)
- [test-lite](/test-lite) - `protobuf-go-lite` auto-generated serialization (`146kb` binary, `3s` build time)

Any of these options and others will work for protobuf serialization.

## Roadmap

This project is in alpha. Breaking API changes should be expected until Beta.

- `v0.0.x` - Alpha
  - [ ] Stabilize API
- `v0.x.x` - Beta
  - [ ] Finalize API
  - [ ] Test in production
- `v1.x.x` - General Availability
  - [ ] Proven long term stability in production
