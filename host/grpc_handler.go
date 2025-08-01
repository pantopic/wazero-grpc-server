// Borrowed heavily from mwitkow/grpc-proxy (Apache 2.0)
// See https://github.com/mwitkow/grpc-proxy/blob/master/proxy/handler.go

package wazero_grpc_server

import (
	"context"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/pantopic/wazero-pool"
)

type grpcHandler struct {
	pool wazeropool.Instance
	meta *meta
}

func (h *grpcHandler) handler(f handlerFactory) func(srv any, serverStream grpc.ServerStream) error {
	return func(srv any, serverStream grpc.ServerStream) error {
		fullMethodName, ok := grpc.MethodFromServerStream(serverStream)
		if !ok {
			return status.Errorf(codes.Internal, "lowLevelServerStream not exists in context")
		}
		ctx, cancel := context.WithCancel(serverStream.Context())
		defer cancel()
		mod := h.pool.Get()
		defer h.pool.Put(mod)
		clientStream := f(ctx, mod, h.meta, fullMethodName)
		s2cErrChan := h.forwardServerToWazero(serverStream, clientStream)
		c2sErrChan := h.forwardWazeroToServer(clientStream, serverStream)
		for range 2 {
			select {
			case s2cErr := <-s2cErrChan:
				if s2cErr == io.EOF {
					clientStream.CloseSend()
				} else {
					cancel()
					return status.Errorf(codes.Internal, "failed proxying s2c: %v", s2cErr)
				}
			case c2sErr := <-c2sErrChan:
				serverStream.SetTrailer(clientStream.Trailer())
				if c2sErr != io.EOF {
					return c2sErr
				}
				return nil
			}
		}
		return status.Errorf(codes.Internal, "gRPC proxying should never reach this stage.")
	}
}

func (h *grpcHandler) forwardWazeroToServer(src grpc.ClientStream, dst grpc.ServerStream) chan error {
	ret := make(chan error, 1)
	go func() {
		f := &emptypb.Empty{}
		for i := 0; ; i++ {
			if err := src.RecvMsg(f); err != nil {
				ret <- err
				break
			}
			if i == 0 {
				md, err := src.Header()
				if err != nil {
					ret <- err
					break
				}
				if err := dst.SendHeader(md); err != nil {
					ret <- err
					break
				}
			}
			if err := dst.SendMsg(f); err != nil {
				ret <- err
				break
			}
		}
	}()
	return ret
}

func (h *grpcHandler) forwardServerToWazero(src grpc.ServerStream, dst grpc.ClientStream) chan error {
	ret := make(chan error, 1)
	go func() {
		f := &emptypb.Empty{}
		for i := 0; ; i++ {
			if err := src.RecvMsg(f); err != nil {
				ret <- err
				break
			}
			if err := dst.SendMsg(f); err != nil {
				ret <- err
				break
			}
		}
	}()
	return ret
}
