package main

import (
	"google.golang.org/protobuf/proto"

	"github.com/pantopic/wazero-grpc-server/grpc-server-go"
	"github.com/pantopic/wazero-grpc-server/test/pb"
)

func main() {
	s := grpc_server.NewService(`test.TestService`)
	s.AddMethod(`Test`, protoWrap(test, &pb.TestRequest{}))
}

func test(req *pb.TestRequest) (res *pb.TestResponse, err error) {
	return &pb.TestResponse{
		Bar: req.Foo,
	}, nil
}

func protoWrap[Req proto.Message, Res proto.Message](fn func(Req) (Res, error), req Req) func([]byte) ([]byte, error) {
	return func(in []byte) (out []byte, err error) {
		err = proto.Unmarshal(in, req)
		if err != nil {
			return []byte(err.Error()), grpc_server.ErrMalformed
		}
		res, err := fn(req)
		if err != nil {
			return []byte(err.Error()), grpc_server.ErrUnexpected
		}
		out, err = proto.Marshal(res)
		if err != nil {
			return []byte(err.Error()), grpc_server.ErrMarshal
		}
		return
	}
}
