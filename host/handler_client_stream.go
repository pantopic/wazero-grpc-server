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
)

func newHandlerClientStream(ctx context.Context, mod api.Module, meta *meta, method string) grpc.ClientStream {
	next := make(chan []byte)
	ctx = context.WithValue(ctx, DefaultCtxKeyMeta, meta)
	ctx = context.WithValue(ctx, DefaultCtxKeyNext, next)
	return &handlerClientStream{ctx, mod, meta, method, next, make(chan bool), false}
}

type handlerClientStream struct {
	ctx    context.Context
	mod    api.Module
	meta   *meta
	method string
	next   chan []byte
	done   chan bool
	init   bool
}

func (cs *handlerClientStream) Header() (md metadata.MD, err error) {
	return
}

func (cs *handlerClientStream) Trailer() (md metadata.MD) {
	return
}

func (cs *handlerClientStream) CloseSend() (err error) {
	close(cs.next)
	return
}

func (cs *handlerClientStream) Context() context.Context {
	return cs.ctx
}

func (cs *handlerClientStream) SendMsg(m any) (err error) {
	msg, err := proto.Marshal(m.(proto.Message))
	if err != nil {
		panic(err)
	}
	if !cs.init {
		cs.init = true
		// Special case for first message
		setMethod(cs.mod, cs.meta, []byte(cs.method))
		setMsg(cs.mod, cs.meta, msg)
		setErrCode(cs.mod, cs.meta, uint32(codes.OK))
		go func() {
			_, err := cs.mod.ExportedFunction("__grpcServerCall").Call(cs.ctx)
			if err != nil {
				log.Println(err)
			}
			cs.done <- true
		}()
	} else {
		cs.next <- msg
	}
	return
}

func (cs *handlerClientStream) RecvMsg(m any) (err error) {
	select {
	case _, ok := <-cs.done:
		if !ok {
			return io.EOF
		}
		if ferr := getError(cs.mod, cs.meta); ferr != nil {
			ferr.(*Error).msg += `: ` + string(msg(cs.mod, cs.meta))
			return ferr
		}
		b := msg(cs.mod, cs.meta)
		err = proto.Unmarshal(b, m.(proto.Message))
		close(cs.done)
	case <-cs.ctx.Done():
	}
	return
}
