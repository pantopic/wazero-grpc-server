package main

import (
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

func protoWrap[Req Message, Res Message](fn func(Req) (Res, error), req Req) func([]byte) ([]byte, error) {
	return func(in []byte) (out []byte, err error) {
		err = req.Unmarshal(in)
		if err != nil {
			return []byte(err.Error()), grpc_server.ErrMalformed
		}
		res, err := fn(req)
		if err != nil {
			return []byte(err.Error()), grpc_server.ErrUnexpected
		}
		out = res.Marshal(out)
		return
	}
}
