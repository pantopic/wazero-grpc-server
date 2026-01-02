package main

import (
	"iter"

	"github.com/pantopic/wazero-pipe/sdk-go"

	"github.com/pantopic/wazero-grpc-server/sdk-go"
	"github.com/pantopic/wazero-grpc-server/sdk-go/codes"
	"github.com/pantopic/wazero-grpc-server/sdk-go/status"
	"github.com/pantopic/wazero-grpc-server/test-easy/pb"
)

var bridge *pipe.Pipe[uint64]

func main() {
	bridge = pipe.New[uint64]()
	grpc_server.NewService(`test.TestService`).
		Unary(`Test`, protoWrap(test, &pb.TestRequest{})).
		Unary(`Retest`, protoWrap(retest, &pb.RetestRequest{})).
		Unary(`TestBytes`, protoWrap(testBytes, &pb.TestBytesRequest{})).
		ClientStream(`ClientStream`, protoWrapClientStream(clientStream, &pb.ClientStreamRequest{})).
		ServerStream(`ServerStream`, protoWrapServerStream(serverStream, &pb.ServerStreamRequest{})).
		BidirectionalStream(`BidirectionalStream`,
			protoWrapBidirectionalRecv(bidirectionalStreamRecv, &pb.BidirectionalStreamRequest{}),
			protoWrapBidirectionalSend(bidirectionalStreamSend))
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

func bidirectionalStreamRecv(reqs iter.Seq[*pb.BidirectionalStreamRequest]) (err error) {
	for req := range reqs {
		bridge.Send(req.Foo4)
	}
	return nil
}

func bidirectionalStreamSend() (res iter.Seq[*pb.BidirectionalStreamResponse], err error) {
	return func(yield func(*pb.BidirectionalStreamResponse) bool) {
		var n uint64
		for {
			i, err := bridge.Recv()
			if err != nil {
				break
			}
			n += i
			if !yield(&pb.BidirectionalStreamResponse{Bar4: n}) {
				return
			}
		}
	}, nil
}

func protoWrap[ReqType pb.Message, ResType pb.Message](fn func(ReqType) (ResType, error), req ReqType) func([]byte) ([]byte, error) {
	return func(in []byte) (out []byte, err error) {
		err = req.Unmarshal(in)
		if err != nil {
			return []byte(err.Error()), status.New(codes.InvalidArgument, err.Error()).Err()
		}
		res, err := fn(req)
		if err != nil {
			return []byte(err.Error()), status.New(codes.Unknown, err.Error()).Err()
		}
		out = res.Marshal(out)
		return
	}
}

func protoWrapClientStream[ReqType pb.Message, ResType pb.Message](fn func(iter.Seq[ReqType]) (ResType, error), req ReqType) func(iter.Seq[[]byte]) ([]byte, error) {
	return func(in iter.Seq[[]byte]) (out []byte, err error) {
		var err2 error
		res, err := fn(func(yield func(ReqType) bool) {
			for b := range in {
				err2 = req.Unmarshal(b)
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
		return res.Marshal(out), nil
	}
}

func protoWrapServerStream[ReqType pb.Message, ResType pb.Message](fn func(ReqType) (iter.Seq[ResType], error), req ReqType) func([]byte) (iter.Seq[[]byte], error) {
	return func(in []byte) (out iter.Seq[[]byte], err error) {
		err = req.Unmarshal(in)
		if err != nil {
			return nil, status.New(codes.InvalidArgument, err.Error()).Err()
		}
		all, err := fn(req)
		if err != nil {
			return nil, status.New(codes.Unknown, err.Error()).Err()
		}
		return func(yield func([]byte) bool) {
			for res := range all {
				if !yield(res.Marshal(nil)) {
					return
				}
			}
		}, nil
	}
}

func protoWrapBidirectionalRecv[ReqType pb.Message](fn func(iter.Seq[ReqType]) error, req ReqType) func(iter.Seq[[]byte]) error {
	return func(in iter.Seq[[]byte]) (err error) {
		var err2 error
		err = fn(func(yield func(ReqType) bool) {
			for b := range in {
				err2 = req.Unmarshal(b)
				if err2 != nil {
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
		return
	}
}

func protoWrapBidirectionalSend[ResType pb.Message](fn func() (iter.Seq[ResType], error)) func() (iter.Seq[[]byte], error) {
	return func() (out iter.Seq[[]byte], err error) {
		all, err := fn()
		if err != nil {
			return nil, status.New(codes.Unknown, err.Error()).Err()
		}
		return func(yield func([]byte) bool) {
			for res := range all {
				if !yield(res.Marshal(nil)) {
					return
				}
			}
		}, nil
	}
}
