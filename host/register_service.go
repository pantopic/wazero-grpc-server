// Borrowed heavily from mwitkow/grpc-proxy
// See https://github.com/mwitkow/grpc-proxy/blob/master/proxy/handler.go

package wazero_grpc_server

import (
	"log"
	"strings"

	"google.golang.org/grpc"

	"github.com/pantopic/wazero-pool"
)

func registerService(s *grpc.Server, pool wazeropool.InstancePool, meta *meta, serviceName string, methods []string) {
	h := &grpcHandler{pool, meta}
	fakeDesc := &grpc.ServiceDesc{
		ServiceName: serviceName,
		HandlerType: (*any)(nil),
	}
	for _, m := range methods {
		parts := strings.Split(m, ".")
		if len(parts) < 2 {
			log.Panicf(`%s %#v`, methods, parts)
		}
		var d = grpc.StreamDesc{
			StreamName:    parts[1],
			ServerStreams: true,
			ClientStreams: true,
		}
		switch parts[0] {
		case "u":
			d.Handler = h.handler(newHandlerUnary)
		case "c":
			d.Handler = h.handler(newHandlerClientStream)
		case "s":
			d.Handler = h.handler(newHandlerServerStream)
		}
		fakeDesc.Streams = append(fakeDesc.Streams, d)
	}
	s.RegisterService(fakeDesc, h)
}
