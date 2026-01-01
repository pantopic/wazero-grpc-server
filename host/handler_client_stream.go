package wazero_grpc_server

import (
	"context"
	"io"
	"log"

	"github.com/tetratelabs/wazero/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"

	"github.com/pantopic/wazero-pool"
)

func newHandlerFactoryClientStream(m *hostModule) handlerFactory {
	return func(ctx context.Context, pool wazeropool.Instance, meta *meta, method string) grpc.ClientStream {
		next := make(chan []byte)
		ctx = context.WithValue(ctx, m.ctxKeyMeta, meta)
		ctx = context.WithValue(ctx, m.ctxKeyNext, next)
		return &handlerClientStream{ctx, pool, meta, method, make(chan resp), next, false}
	}
}

type handlerClientStream struct {
	ctx    context.Context
	pool   wazeropool.Instance
	meta   *meta
	method string
	send   chan resp
	next   chan []byte
	init   bool
}

func (h *handlerClientStream) Header() (md metadata.MD, err error) {
	return
}

func (h *handlerClientStream) Trailer() (md metadata.MD) {
	return
}

func (h *handlerClientStream) CloseSend() (err error) {
	close(h.next)
	return
}

func (h *handlerClientStream) Context() context.Context {
	return h.ctx
}

func (h *handlerClientStream) SendMsg(m any) (err error) {
	msg, err := proto.Marshal(m.(proto.Message))
	if err != nil {
		panic(err)
	}
	if !h.init {
		h.init = true
		go func() {
			var r resp
			h.pool.Run(func(mod api.Module) {
				// Special case for first message
				setMethod(mod, h.meta, []byte(h.method))
				setMsg(mod, h.meta, msg)
				setErrCode(mod, h.meta, codes.OK)
				_, err := mod.ExportedFunction("__grpc_server_client_stream").Call(h.ctx)
				if err != nil {
					log.Println(err)
				}
				r.err = getError(mod, h.meta)
				if r.err == nil {
					r.data = append(r.data, getMsg(mod, h.meta)...)
				}
			})
			h.send <- r
		}()
	} else {
		h.next <- msg
	}
	return
}

func (h *handlerClientStream) RecvMsg(m any) (err error) {
	select {
	case r, ok := <-h.send:
		if !ok {
			return io.EOF
		}
		close(h.send)
		if r.err != nil {
			return r.err
		}
		err = proto.Unmarshal(r.data, m.(proto.Message))
	case <-h.ctx.Done():
	}
	return
}
