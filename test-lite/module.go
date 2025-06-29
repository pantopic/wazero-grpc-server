package main

import (
	proto "github.com/aperturerobotics/protobuf-go-lite"

	"github.com/pantopic/wazero-grpc-server/grpc-server-go"
	"github.com/pantopic/wazero-grpc-server/test-lite/pb"
)

func main() {
	s := grpc_server.NewService(`test.TestService`)
	s.AddMethod(`Test`, protoWrap(test, &pb.TestRequest{}))
	s.AddMethod(`Retest`, protoWrap(retest, &pb.RetestRequest{}))
}

func test(req *pb.TestRequest) (res *pb.TestResponse, err error) {
	return &pb.TestResponse{
		Bar: req.Foo,
	}, nil
}

func retest(req *pb.RetestRequest) (res *pb.RetestResponse, err error) {
	return &pb.RetestResponse{
		Foo: req.Bar,
	}, nil
}

func protoWrap[ReqType proto.Message, ResType proto.Message](fn func(ReqType) (ResType, error), req ReqType) func([]byte) ([]byte, error) {
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
