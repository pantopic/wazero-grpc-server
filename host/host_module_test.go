package wazero_grpc_server

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"net"
	"testing"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/pantopic/wazero-grpc-server/host/pb"
	"github.com/pantopic/wazero-pool"
)

//go:embed test-easy\.wasm
var testWasmEasy []byte

//go:embed test-easy\.prod\.wasm
var testWasmEasyProd []byte

//go:embed test-lite\.wasm
var testWasmLite []byte

//go:embed test-lite\.prod\.wasm
var testWasmLiteProd []byte

func TestHostModule(t *testing.T) {
	var (
		ctx = context.Background()
		out = &bytes.Buffer{}
	)
	r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig().
		WithMemoryLimitPages(256)) // 16 MB
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	var hostModule *hostModule
	t.Run(`register`, func(t *testing.T) {
		hostModule = New()
		hostModule.Register(ctx, r)
	})

	port := 9000
	for _, tc := range []struct {
		name string
		wasm []byte
	}{
		{`testWasmEasy`, testWasmEasy},
		{`testWasmLite`, testWasmLite},
		{`testWasmEasyProd`, testWasmEasyProd},
		{`testWasmLiteProd`, testWasmLiteProd},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := grpc.NewServer()
			cfg := wazero.NewModuleConfig().WithStdout(out)
			pool, err := wazeropool.New(ctx, r, tc.wasm, cfg)
			if err != nil {
				t.Fatalf(`%v`, err)
			}
			ctx, err = hostModule.RegisterServices(ctx, s, pool)
			if err != nil {
				t.Fatalf(`%v`, err)
			}
			meta := get[*meta](ctx, hostModule.ctxKeyMeta)
			mod := pool.Get()
			if readUint32(mod, meta.ptrMethodMax) != 256 {
				t.Errorf("incorrect maximum method length: %#v", meta)
			}
			pool.Put(mod)
			port++
			addr := fmt.Sprintf(`:%d`, port)
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
				Foo: 1,
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
			cs, err := client.ClientStream(ctx)
			if err != nil {
				t.Fatalf(`%v`, err)
			}
			for range 50 {
				if err = cs.Send(&pb.ClientStreamRequest{
					Foo2: 2,
				}); err != nil {
					t.Fatalf(`%v`, err)
				}
			}
			res3, err := cs.CloseAndRecv()
			if err != nil {
				t.Fatalf(`%v`, err)
			}
			if res3.Bar2 != 100 {
				t.Fatalf(`Incorrect response value in ClientStream response: %d`, res3.Bar2)
			}
		})
	}
}

func BenchmarkHostModule(b *testing.B) {
	var ctx = context.Background()
	var out = &bytes.Buffer{}
	r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig().
		WithMemoryLimitPages(256)) // 16 MB
	wasi_snapshot_preview1.MustInstantiate(ctx, r)
	hostModule := New()
	hostModule.Register(ctx, r)
	b.Run(`linear`, func(b *testing.B) {
		for _, tc := range []struct {
			name string
			wasm []byte
		}{
			{`testWasmEasy`, testWasmEasy},
			{`testWasmLite`, testWasmLite},
			{`testWasmEasyProd`, testWasmEasyProd},
			{`testWasmLiteProd`, testWasmLiteProd},
		} {
			s := grpc.NewServer()
			cfg := wazero.NewModuleConfig().WithStdout(out)
			pool, err := wazeropool.New(ctx, r, tc.wasm, cfg)
			if err != nil {
				b.Fatalf(`%v`, err)
			}
			addr := `:9001`
			ctx, err = hostModule.RegisterServices(ctx, s, pool)
			if err != nil {
				b.Fatalf(`%v`, err)
			}
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
	})
	for _, n := range []int{0, 2, 4, 8, 16} {
		b.Run(fmt.Sprintf(`parallel-%d`, n), func(b *testing.B) {
			for _, tc := range []struct {
				name string
				wasm []byte
			}{
				{`testWasmEasyProd`, testWasmEasyProd},
				{`testWasmLiteProd`, testWasmLiteProd},
			} {
				s := grpc.NewServer()
				cfg := wazero.NewModuleConfig().WithStdout(out)
				pool, err := wazeropool.New(ctx, r, tc.wasm, cfg, wazeropool.WithLimit(n))
				if err != nil {
					b.Fatalf(`%v`, err)
				}
				addr := `:9001`
				ctx, err = hostModule.RegisterServices(ctx, s, pool)
				if err != nil {
					b.Fatalf(`%v`, err)
				}
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
					b.SetParallelism(n)
					b.RunParallel(func(pb *testing.PB) {
						for pb.Next() {
							res, _ = client.Test(ctx, req)
						}
					})
				})
				if res.Bar != 20 {
					b.Fatalf(`Nope`)
				}
				s.Stop()
			}
		})
	}
}
