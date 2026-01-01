package wazero_grpc_server

import (
	"context"
	"io"

	"github.com/tetratelabs/wazero/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"

	"github.com/pantopic/wazero-pool"
)

func newHandlerFactoryUnary(_ *hostModule) handlerFactory {
	return func(ctx context.Context, pool wazeropool.Instance, meta *meta, method string) grpc.ClientStream {
		return &handlerUnary{ctx, pool, meta, method, make(chan resp, 10)}
	}
}

type handlerUnary struct {
	ctx    context.Context
	pool   wazeropool.Instance
	meta   *meta
	method string
	send   chan resp
}

type resp struct {
	data []byte
	err  error
}

func (h *handlerUnary) Header() (md metadata.MD, err error) {
	return
}

func (h *handlerUnary) Trailer() (md metadata.MD) {
	return
}

func (h *handlerUnary) CloseSend() (err error) {
	return
}

func (h *handlerUnary) Context() context.Context {
	return h.ctx
}

func (h *handlerUnary) SendMsg(m any) (err error) {
	data, err := proto.Marshal(m.(proto.Message))
	if err != nil {
		panic(err)
	}
	var r resp
	h.pool.Run(func(mod api.Module) {
		setMethod(mod, h.meta, []byte(h.method))
		setMsg(mod, h.meta, data)
		setErrCode(mod, h.meta, codes.OK)
		if _, err = mod.ExportedFunction("__grpc_server_unary").Call(h.ctx); err != nil {
			return
		}
		r.err = getError(mod, h.meta)
		if r.err == nil {
			r.data = append(r.data, getMsg(mod, h.meta)...)
		}
	})
	h.send <- r
	return
}

func (h *handlerUnary) RecvMsg(m any) (err error) {
	select {
	case resp, ok := <-h.send:
		if !ok {
			return io.EOF
		}
		close(h.send)
		if resp.err != nil {
			return resp.err
		}
		err = proto.Unmarshal(resp.data, m.(proto.Message))
	case <-h.ctx.Done():
	}
	return
}
