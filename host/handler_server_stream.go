package wazero_grpc_server

import (
	"context"
	"io"
	"sync"

	"github.com/tetratelabs/wazero/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"

	"github.com/pantopic/wazero-pool"
)

type msgErr struct {
	msg []byte
	err error
	wg  *sync.WaitGroup
}

func newHandlerFactoryServerStream(m *hostModule) handlerFactory {
	return func(ctx context.Context, pool wazeropool.Instance, meta *meta, method string) grpc.ClientStream {
		s := &handlerServerStream{ctx, pool, meta, method, make(chan msgErr)}
		s.ctx = context.WithValue(s.ctx, m.ctxKeyMeta, meta)
		s.ctx = context.WithValue(s.ctx, m.ctxKeySend, s.send)
		return s
	}
}

type handlerServerStream struct {
	ctx    context.Context
	pool   wazeropool.Instance
	meta   *meta
	method string
	data   chan msgErr
}

func (h *handlerServerStream) Header() (md metadata.MD, err error) {
	return
}

func (h *handlerServerStream) Trailer() (md metadata.MD) {
	return
}

func (h *handlerServerStream) CloseSend() (err error) {
	close(h.data)
	return
}

func (h *handlerServerStream) Context() context.Context {
	return h.ctx
}

func (h *handlerServerStream) send(msg []byte, err error) {
	d := msgErr{msg, err, &sync.WaitGroup{}}
	d.wg.Add(1)
	h.data <- d
	d.wg.Wait()
}

func (h *handlerServerStream) SendMsg(m any) (err error) {
	msg, err := proto.Marshal(m.(proto.Message))
	if err != nil {
		panic(err)
	}
	h.pool.Run(func(mod api.Module) {
		setMethod(mod, h.meta, []byte(h.method))
		setMsg(mod, h.meta, msg)
		mod.ExportedFunction("__grpc_server_server_stream").Call(h.ctx)
	})
	return
}

func (h *handlerServerStream) RecvMsg(m any) (err error) {
	select {
	case d, ok := <-h.data:
		if !ok {
			return io.EOF
		}
		defer d.wg.Done()
		if d.err != nil {
			return &Error{
				msg: d.err.Error(),
			}
		}
		err = proto.Unmarshal(d.msg, m.(proto.Message))
	case <-h.ctx.Done():
	}
	return
}
