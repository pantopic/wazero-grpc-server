package main

import (
	"google.golang.org/protobuf/proto"

	"github.com/pantopic/wazero-grpc/grpc-server-go"
)

func main() {
	s := grpc.NewService(`test.Test`)
	s.AddMethod(`Test`, protoWrap(test, &TestRequest{}))
}

func protoWrap(fn func(proto.Message) (proto.Message, error), req proto.Message) func([]byte) ([]byte, error) {
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

func test(req proto.Message) (res proto.Message, err error) {
	req = req.(*TestRequest)
	res = &TestResponse{}
	return
}
