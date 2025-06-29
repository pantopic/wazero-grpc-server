// Borrowed heavily from mwitkow/grpc-proxy
// See https://github.com/mwitkow/grpc-proxy/blob/master/proxy/handler.go

package wazero_grpc_server

import (
	"github.com/tetratelabs/wazero/api"
	"google.golang.org/grpc"
)

func registerService(s *grpc.Server, m api.Module, meta *meta, serviceName string, methods []string) {
	h := &grpcHandler{m, meta}
	fakeDesc := &grpc.ServiceDesc{
		ServiceName: serviceName,
		HandlerType: (*any)(nil),
	}
	for _, m := range methods {
		streamDesc := grpc.StreamDesc{
			StreamName:    m,
			Handler:       h.handler,
			ServerStreams: true,
			ClientStreams: true,
		}
		fakeDesc.Streams = append(fakeDesc.Streams, streamDesc)
	}
	s.RegisterService(fakeDesc, h)
}
