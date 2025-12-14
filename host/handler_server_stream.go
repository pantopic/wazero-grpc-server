package wazero_grpc_server

import (
	"context"
	"io"
	"sync"

	"github.com/tetratelabs/wazero/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

type msgErr struct {
	msg []byte
	err error
	wg  *sync.WaitGroup
}

func newHandlerServerStream(ctx context.Context, mod api.Module, meta *meta, method string) grpc.ClientStream {
	s := &handlerServerStream{ctx, mod, meta, method, make(chan msgErr)}
	s.ctx = context.WithValue(s.ctx, DefaultCtxKeyMeta, meta)
	s.ctx = context.WithValue(s.ctx, DefaultCtxKeySend, s.send)
	return s
}

type handlerServerStream struct {
	ctx    context.Context
	mod    api.Module
	meta   *meta
	method string
	data   chan msgErr
}

func (cs *handlerServerStream) Header() (md metadata.MD, err error) {
	return
}

func (cs *handlerServerStream) Trailer() (md metadata.MD) {
	return
}

func (cs *handlerServerStream) CloseSend() (err error) {
	close(cs.data)
	return
}

func (cs *handlerServerStream) Context() context.Context {
	return cs.ctx
}

func (cs *handlerServerStream) send(msg []byte, err error) {
	d := msgErr{msg, err, &sync.WaitGroup{}}
	d.wg.Add(1)
	cs.data <- d
	d.wg.Wait()
}

func (cs *handlerServerStream) SendMsg(m any) (err error) {
	msg, err := proto.Marshal(m.(proto.Message))
	if err != nil {
		panic(err)
	}
	setMethod(cs.mod, cs.meta, []byte(cs.method))
	setMsg(cs.mod, cs.meta, msg)
	cs.mod.ExportedFunction("__grpcServerCall").Call(cs.ctx)
	return
}

func (cs *handlerServerStream) RecvMsg(m any) (err error) {
	select {
	case d, ok := <-cs.data:
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
	case <-cs.ctx.Done():
	}
	return
}
