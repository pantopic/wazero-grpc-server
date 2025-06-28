package main

import (
	proto "github.com/aperturerobotics/protobuf-go-lite"

	"github.com/pantopic/wazero-grpc-server/grpc-server-go"
)

func main() {
	s := grpc_server.NewService(`test.TestService`)
	s.AddMethod(`Test`, protoWrap(test, &TestRequest{}))
}

func test(req *TestRequest) (res *TestResponse, err error) {
	return &TestResponse{
		Bar: req.Foo,
	}, nil
}

func protoWrap[Req proto.Message, Res proto.Message](fn func(Req) (Res, error), req Req) func([]byte) ([]byte, error) {
	return func(in []byte) (out []byte, err error) {
		err = req.UnmarshalVT(in)
		if err != nil {
			return []byte(err.Error()), grpc_server.ErrMalformed
		}
		res, err := fn(req)
		if err != nil {
			return []byte(err.Error()), grpc_server.ErrUnexpected
		}
		out, err = res.MarshalVT()
		if err != nil {
			return []byte(err.Error()), grpc_server.ErrMarshal
		}
		return
	}
}
