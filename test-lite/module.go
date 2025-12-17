package main

import (
	"iter"

	proto "github.com/aperturerobotics/protobuf-go-lite"

	"github.com/pantopic/wazero-grpc-server/sdk-go"
	"github.com/pantopic/wazero-grpc-server/sdk-go/codes"
	"github.com/pantopic/wazero-grpc-server/sdk-go/status"
	"github.com/pantopic/wazero-grpc-server/test-lite/pb"
)

func main() {
	s := grpc_server.NewService(`test.TestService`)
	s.Unary(`Test`, protoWrap(test, &pb.TestRequest{}))
	s.Unary(`Retest`, protoWrap(retest, &pb.RetestRequest{}))
	s.Unary(`TestBytes`, protoWrap(testBytes, &pb.TestBytesRequest{}))
	s.ClientStream(`ClientStream`, protoWrapClientStream(clientStream, &pb.ClientStreamRequest{}))
	s.ServerStream(`ServerStream`, protoWrapServerStream(serverStream, &pb.ServerStreamRequest{}))
}

func test(req *pb.TestRequest) (res *pb.TestResponse, err error) {
	return &pb.TestResponse{Bar: req.Foo}, nil
}

func retest(req *pb.RetestRequest) (res *pb.RetestResponse, err error) {
	return &pb.RetestResponse{Foo: req.Bar}, nil
}

func testBytes(req *pb.TestBytesRequest) (res *pb.TestBytesResponse, err error) {
	return &pb.TestBytesResponse{Code: 1, Data: []byte(`ACK`)}, nil
}

func clientStream(reqs iter.Seq[*pb.ClientStreamRequest]) (res *pb.ClientStreamResponse, err error) {
	var n uint64
	for req := range reqs {
		n += req.Foo2
	}
	return &pb.ClientStreamResponse{Bar2: n}, nil
}

func serverStream(req *pb.ServerStreamRequest) (res iter.Seq[*pb.ServerStreamResponse], err error) {
	return func(yield func(*pb.ServerStreamResponse) bool) {
		var n uint64
		for range req.Foo3 {
			if !yield(&pb.ServerStreamResponse{Bar3: n}) {
				return
			}
			n++
		}
	}, nil
}

func protoWrap[ReqType proto.Message, ResType proto.Message](fn func(ReqType) (ResType, error), req ReqType) func([]byte) ([]byte, error) {
	return func(in []byte) (out []byte, err error) {
		err = req.UnmarshalVT(in)
		if err != nil {
			return []byte(err.Error()), status.New(codes.InvalidArgument, err.Error()).Err()
		}
		res, err := fn(req)
		if err != nil {
			return []byte(err.Error()), status.New(codes.Unknown, err.Error()).Err()
		}
		out, err = res.MarshalVT()
		if err != nil {
			return []byte(err.Error()), status.New(codes.InvalidArgument, err.Error()).Err()
		}
		return
	}
}

func protoWrapClientStream[ReqType proto.Message, ResType proto.Message](fn func(iter.Seq[ReqType]) (ResType, error), req ReqType) func(iter.Seq[[]byte]) ([]byte, error) {
	return func(in iter.Seq[[]byte]) (out []byte, err error) {
		var err2 error
		res, err := fn(func(yield func(ReqType) bool) {
			for b := range in {
				err2 = req.UnmarshalVT(b)
				if err2 != nil {
					out = []byte(err.Error())
					err2 = status.New(codes.InvalidArgument, err.Error()).Err()
					return
				}
				if !yield(req) {
					return
				}
			}
		})
		if err2 != nil {
			err = err2
			return
		}
		out, err = res.MarshalVT()
		if err != nil {
			return []byte(err.Error()), status.New(codes.InvalidArgument, err.Error()).Err()
		}
		return
	}
}

func protoWrapServerStream[ReqType proto.Message, ResType proto.Message](fn func(ReqType) (iter.Seq[ResType], error), req ReqType) func([]byte) (iter.Seq[[]byte], error) {
	return func(in []byte) (out iter.Seq[[]byte], err error) {
		err = req.UnmarshalVT(in)
		if err != nil {
			return nil, status.New(codes.InvalidArgument, err.Error()).Err()
		}
		all, err := fn(req)
		if err != nil {
			return nil, status.New(codes.Unknown, err.Error()).Err()
		}
		return func(yield func([]byte) bool) {
			for res := range all {
				out, err2 := res.MarshalVT()
				if err2 != nil {
					err2 = status.New(codes.InvalidArgument, err.Error()).Err()
					return
				}
				if !yield(out) {
					return
				}
			}
		}, nil
	}
}
