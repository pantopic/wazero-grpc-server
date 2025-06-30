# Wazero gRPC Server

A [wazero](https://pkg.go.dev/github.com/tetratelabs/wazero) host module, ABI and guest SDK enabling WASI modules to register as [gRPC](https://grpc.io/) services.

## Host Module

[![Go Reference](https://godoc.org/github.com/pantopic/wazero-grpc-server/host?status.svg)](https://godoc.org/github.com/pantopic/wazero-grpc-server/host)
[![Go Report Card](https://goreportcard.com/badge/github.com/pantopic/wazero-grpc-server/host)](https://goreportcard.com/report/github.com/pantopic/wazero-grpc-server/host)
[![Go Coverage](https://github.com/pantopic/wazero-grpc-server/wiki/host/coverage.svg)](https://raw.githack.com/wiki/pantopic/wazero-grpc-server/host/coverage.html)

First register the host module with the runtime

```go
import (
	"context"
	_ "embed"
	"net"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"google.golang.org/grpc"

	"github.com/pantopic/wazero-grpc-server/host"
)

//go:embed test\.wasm
var wasm []byte

func main() {
	ctx := context.Background()
	r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig())
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	module := wazero_grpc_server.New()
	module.Register(ctx, r)

	mod, _ := r.Instantiate(ctx, wasm)
	hostModule.RegisterService(ctx, s, mod)

	lis, _ := net.Listen(`tcp`, `:8000`)
	s := grpc.NewServer()
	s.Serve(lis)

	// ...
}
```

## Guest SDK (Go)

[![Go Reference](https://godoc.org/github.com/pantopic/wazero-grpc-server/grpc-server-go?status.svg)](https://godoc.org/github.com/pantopic/wazero-grpc-server/grpc-server-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/pantopic/wazero-grpc-server/grpc-server-go)](https://goreportcard.com/report/github.com/pantopic/wazero-grpc-server/grpc-server-go)

Then you can import the guest SDK into your WASI module to export your gRPC service description at runtime and receive gRPC requests in WASM.

```go
package main

import (
	proto "github.com/aperturerobotics/protobuf-go-lite"

	"github.com/pantopic/wazero-grpc-server/grpc-server-go"
	"github.com/pantopic/wazero-grpc-server/test-lite/pb"
)

func main() {
	s := grpc_server.NewService(`test.TestService`)
	s.Unary(`Test`, protoWrap(test, &pb.TestRequest{}))
	s.Unary(`Retest`, protoWrap(retest, &pb.RetestRequest{}))
}

func test(req *pb.TestRequest) (res *pb.TestResponse, err error) {
	return &pb.TestResponse{Bar: req.Foo}, nil
}

func retest(req *pb.RetestRequest) (res *pb.RetestResponse, err error) {
	return &pb.RetestResponse{Foo: req.Bar}, nil
}

func protoWrap[ReqType proto.Message, ResType proto.Message](fn func(ReqType) (ResType, error), req ReqType) func([]byte) ([]byte, error) {
	return func(in []byte) (out []byte, err error) {
		err = req.UnmarshalVT(in)
		if err != nil {
			return []byte(err.Error()), grpc_server.ErrMalformed
		}
		res, err := fn(req)
		if err != nil {
			return []byte(err.Error()), grpc_server.ErrUnexpected
		}
		out, err = res.MarshalVT()
		if err != nil {
			return []byte(err.Error()), grpc_server.ErrMarshal
		}
		return
	}
}
```

The [guest SDK](https://pkg.go.dev/github.com/pantopic/wazero-grpc-server/grpc-server-go) has no dependencies outside the Go std lib.
The guest SDK is serialization agnostic in order to provide users with more control over performance, compile time and binary size.

See examples for protobuf message serialization options:

- [test-easy](/test-easy) - `easyproto` manually-generated (`113kb` binary, `3s` build time)
- [test-lite](/test-lite) - `protobuf-go-lite` auto-generated (`113kb` binary, `3s` build time, recommended)

Any of these options and others can be used for protobuf serialization in WASM.

## Roadmap

This project is in alpha. Breaking API changes should be expected until Beta.

- `v0.0.x` - Alpha
  - [ ] Stabilize API
- `v0.x.x` - Beta
  - [ ] Finalize API
  - [ ] Test in production
- `v1.x.x` - General Availability
  - [ ] Proven long term stability in production
