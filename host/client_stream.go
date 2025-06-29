package wazero_grpc_server

import (
	"context"
	"io"

	"github.com/tetratelabs/wazero/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

func newClientStream(ctx context.Context, mod api.Module, meta *meta, method string) grpc.ClientStream {
	return &clientStream{ctx, mod, meta, method, make(chan bool)}
}

type clientStream struct {
	ctx    context.Context
	mod    api.Module
	meta   *meta
	method string
	ready  chan bool
}

func (cs *clientStream) Header() (md metadata.MD, err error) {
	return
}

func (cs *clientStream) Trailer() (md metadata.MD) {
	return
}

func (cs *clientStream) CloseSend() (err error) {
	return
}

func (cs *clientStream) Context() context.Context {
	return cs.ctx
}

func (cs *clientStream) SendMsg(m any) (err error) {
	msg, err := proto.Marshal(m.(proto.Message))
	if err != nil {
		panic(err)
	}
	setMethod(cs.mod, cs.meta, []byte(cs.method))
	setMsg(cs.mod, cs.meta, msg)
	cs.mod.ExportedFunction("grpcCall").Call(cs.ctx)
	cs.ready <- true
	return
}

func (cs *clientStream) RecvMsg(m any) (err error) {
	select {
	case _, ok := <-cs.ready:
		if !ok {
			return io.EOF
		}
		if ferr := getError(cs.mod, cs.meta); ferr != nil {
			ferr.msg += `: ` + string(msg(cs.mod, cs.meta))
			return ferr
		}
		b := msg(cs.mod, cs.meta)
		err = proto.Unmarshal(b, m.(proto.Message))
		close(cs.ready)
	case <-cs.ctx.Done():
	}
	return
}
