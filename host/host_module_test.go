package wazero_grpc_server

import (
	"bytes"
	"context"
	_ "embed"
	"net"
	"testing"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/pantopic/wazero-grpc-server/host/pb"
)

//go:embed test\.wasm
var testWasm []byte

//go:embed test-easy\.wasm
var testWasmEasy []byte

//go:embed test-lite\.wasm
var testWasmLite []byte

func TestHostModule(t *testing.T) {
	var (
		ctx = context.Background()
		out = &bytes.Buffer{}
	)
	r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig().
		WithMemoryLimitPages(256).
		WithMemoryCapacityFromMax(true))
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	var hostModule *hostModule
	t.Run(`register`, func(t *testing.T) {
		hostModule = New(
			WithCtxKeyMeta(`test_key_meta`),
			WithCtxKeyServer(`test_key_server`),
		)
		hostModule.Register(ctx, r)
	})

	for _, tc := range []struct {
		name string
		wasm []byte
	}{
		{`testWasm`, testWasm},
		{`testWasmEasy`, testWasmEasy},
		{`testWasmLite`, testWasmLite},
	} {
		t.Run(tc.name, func(t *testing.T) {
			compiled, err := r.CompileModule(ctx, tc.wasm)
			if err != nil {
				panic(err)
			}
			cfg := wazero.NewModuleConfig().WithStdout(out)
			mod, err := r.InstantiateModule(ctx, compiled, cfg.WithName(tc.name))
			if err != nil {
				t.Fatalf(`%v`, err)
			}

			ctx, err = hostModule.InitContext(ctx, mod)
			if err != nil {
				t.Fatalf(`%v`, err)
			}
			meta := get[*meta](ctx, hostModule.ctxKeyMeta)
			if readUint32(mod, meta.ptrMethodMax) != 256 {
				t.Errorf("incorrect maximum method length: %#v", meta)
			}

			s := grpc.NewServer()
			addr := `:9001`
			ctx = hostModule.RegisterService(ctx, s, mod)
			lis, err := net.Listen(`tcp`, addr)
			if err != nil {
				t.Fatalf(`%v`, err)
			}
			go func() {
				if err := s.Serve(lis); err != nil {
					panic(err)
				}
			}()
			defer s.Stop()
			conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				t.Fatalf(`%v`, err)
			}
			client := pb.NewTestServiceClient(conn)
			req := &pb.TestRequest{
				Foo: 2,
			}
			res, err := client.Test(ctx, req)
			if err != nil {
				t.Fatalf(`%v`, err)
			}
			if res.Bar != req.Foo {
				t.Fatalf(`Incorrect response value in test response: %d`, res.Bar)
			}
			res2, err := client.Retest(ctx, &pb.RetestRequest{
				Bar: 11,
			})
			if err != nil {
				t.Fatalf(`%v`, err)
			}
			if res2.Foo != 11 {
				t.Fatalf(`Incorrect response value in retest response: %d`, res2.Foo)
			}
		})
	}
}

//go:embed test\.prod\.wasm
var testWasmProd []byte

//go:embed test-easy\.prod\.wasm
var testWasmProdEasy []byte

//go:embed test-lite\.prod\.wasm
var testWasmProdLite []byte

func BenchmarkHostModule(b *testing.B) {
	var ctx = context.Background()
	var out = &bytes.Buffer{}
	r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig().
		WithMemoryLimitPages(256).
		WithMemoryCapacityFromMax(true))
	wasi_snapshot_preview1.MustInstantiate(ctx, r)
	hostModule := New()
	hostModule.Register(ctx, r)
	for _, tc := range []struct {
		name string
		wasm []byte
	}{
		{`testWasm`, testWasm},
		{`testWasmEasy`, testWasmEasy},
		{`testWasmLite`, testWasmLite},
		{`testWasmProd`, testWasmProd},
		{`testWasmProdEasy`, testWasmProdEasy},
		{`testWasmProdLite`, testWasmProdLite},
	} {
		compiled, err := r.CompileModule(ctx, tc.wasm)
		if err != nil {
			panic(err)
		}
		cfg := wazero.NewModuleConfig().WithStdout(out)
		mod, _ := r.InstantiateModule(ctx, compiled, cfg.WithName(tc.name))
		ctx, _ = hostModule.InitContext(ctx, mod)
		s := grpc.NewServer()
		addr := `:9001`
		ctx = hostModule.RegisterService(ctx, s, mod)
		lis, _ := net.Listen(`tcp`, addr)
		go func() {
			if err := s.Serve(lis); err != nil {
				panic(err)
			}
		}()
		conn, _ := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		client := pb.NewTestServiceClient(conn)
		req := &pb.TestRequest{Foo: 20}
		var res *pb.TestResponse
		b.Run(tc.name, func(b *testing.B) {
			for b.Loop() {
				res, _ = client.Test(ctx, req)
			}
		})
		if res.Bar != 20 {
			b.Fatalf(`Nope`)
		}
		s.Stop()
	}
}
