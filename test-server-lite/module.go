package main

import (
	proto "github.com/aperturerobotics/protobuf-go-lite"

	"github.com/pantopic/wazero-grpc/grpc-server-go"
)

func main() {
	s := grpc.NewService(`test.Test`)
	s.AddMethod(`Test`, protoWrap(test, &TestRequest{}))
}

func protoWrap(fn func(proto.Message) (proto.Message, error), req proto.Message) func([]byte) ([]byte, error) {
	return func(in []byte) (out []byte, err error) {
		err = req.UnmarshalVT(in)
		if err != nil {
			return []byte(err.Error()), grpc.ErrMalformed
		}
		res, err := fn(req)
		if err != nil {
			return []byte(err.Error()), grpc.ErrUnexpected
		}
		out, err = res.MarshalVT()
		if err != nil {
			return []byte(err.Error()), grpc.ErrMarshal
		}
		return
	}
}

func test(req proto.Message) (res proto.Message, err error) {
	req = req.(*TestRequest)
	res = &TestResponse{}
	return
}
