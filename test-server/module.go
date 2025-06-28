package main

import (
	"google.golang.org/protobuf/proto"

	"github.com/pantopic/wazero-grpc/grpc-server-go"
)

func main() {
	s := grpc.NewService(`test.TestService`)
	s.AddMethod(`Test`, protoWrap(test, &TestRequest{}))
}

func test(req *TestRequest) (res *TestResponse, err error) {
	return &TestResponse{
		Bar: req.Foo,
	}, nil
}

func protoWrap[Req proto.Message, Res proto.Message](fn func(Req) (Res, error), req Req) func([]byte) ([]byte, error) {
	return func(in []byte) (out []byte, err error) {
		err = proto.Unmarshal(in, req)
		if err != nil {
			return []byte(err.Error()), grpc.ErrMalformed
		}
		res, err := fn(req)
		if err != nil {
			return []byte(err.Error()), grpc.ErrUnexpected
		}
		out, err = proto.Marshal(res)
		if err != nil {
			return []byte(err.Error()), grpc.ErrMarshal
		}
		return
	}
}
