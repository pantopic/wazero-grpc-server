package wazero_grpc_server

import (
	"context"
	"fmt"

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
	pm, ok := m.(proto.Message)
	if !ok {
		return fmt.Errorf("Invalid message")
	}
	msg, err := proto.Marshal(pm)
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
	pm, ok := m.(proto.Message)
	if !ok {
		return fmt.Errorf("Invalid message")
	}
	select {
	case <-cs.ready:
		err = proto.Unmarshal(msg(cs.mod, cs.meta), pm)
	case <-cs.ctx.Done():
	}
	return
}
