package main

import (
	"github.com/pantopic/wazero-grpc/grpc-server-go"
)

func main() {
	s := grpc.NewService(`test.Test`)
	s.AddMethod(`Test`, test)
}

func test(in []byte) (out []byte, err error) {
	var req TestRequest
	if err = req.Unmarshal(in); err != nil {
		return
	}
	res := &TestResponse{
		Bar: req.Foo,
	}
	out = res.Marshal(out)
	return
}
