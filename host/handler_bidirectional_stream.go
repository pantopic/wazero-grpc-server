package wazero_grpc_server

import (
	"context"
	"io"
	"log"
	"sync"

	"github.com/tetratelabs/wazero/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"

	"github.com/pantopic/wazero-pool"
)

func newHandlerFactoryBidirectionalStream(m *hostModule) handlerFactory {
	return func(ctx context.Context, pool wazeropool.Instance, meta *meta, method string) grpc.ClientStream {
		next := make(chan []byte)
		s := &handlerBidirectionalStream{ctx, pool, meta, method, make(chan resp), make(chan msgErr), next, false, false}
		s.ctx = context.WithValue(s.ctx, m.ctxKeyMeta, meta)
		s.ctx = context.WithValue(s.ctx, m.ctxKeyNext, next)
		s.ctx = context.WithValue(s.ctx, m.ctxKeySend, s.send)
		return s
	}
}

type handlerBidirectionalStream struct {
	ctx      context.Context
	pool     wazeropool.Instance
	meta     *meta
	method   string
	chanSend chan resp
	data     chan msgErr
	next     chan []byte
	initSend bool
	initRecv bool
}

func (h *handlerBidirectionalStream) Header() (md metadata.MD, err error) {
	return
}

func (h *handlerBidirectionalStream) Trailer() (md metadata.MD) {
	return
}

func (h *handlerBidirectionalStream) CloseSend() (err error) {
	close(h.next)
	return
}

func (h *handlerBidirectionalStream) Context() context.Context {
	return h.ctx
}

func (h *handlerBidirectionalStream) SendMsg(m any) (err error) {
	msg, err := proto.Marshal(m.(proto.Message))
	if err != nil {
		panic(err)
	}
	if !h.initSend {
		h.initSend = true
		go func() {
			var err error
			// Start recv module instance
			h.pool.Run(func(mod api.Module) {
				setMethod(mod, h.meta, []byte(h.method))
				setErrCode(mod, h.meta, codes.OK)
				_, err = mod.ExportedFunction("__grpc_server_bidirectional_recv").Call(h.ctx)
				if err != nil {
					log.Println(err)
				}
				err = getError(mod, h.meta)
			})
			if err != nil {
				h.send(nil, err)
			}
		}()
	}
	h.next <- msg
	return
}

func (h *handlerBidirectionalStream) send(msg []byte, err error) {
	d := msgErr{append([]byte{}, msg...), err, &sync.WaitGroup{}}
	d.wg.Add(1)
	h.data <- d
	d.wg.Wait()
}

func (h *handlerBidirectionalStream) RecvMsg(m any) (err error) {
	if !h.initRecv {
		h.initRecv = true
		go func() {
			var r msgErr
			// Start send module instance
			h.pool.Run(func(mod api.Module) {
				setMethod(mod, h.meta, []byte(h.method))
				setErrCode(mod, h.meta, codes.OK)
				_, err := mod.ExportedFunction("__grpc_server_bidirectional_send").Call(h.ctx)
				if err != nil {
					log.Println(err)
				}
				r.err = getError(mod, h.meta)
			})
			if r.err != nil {
				h.send(nil, err)
			}
		}()
	}
	select {
	case d, ok := <-h.data:
		if !ok {
			return io.EOF
		}
		defer d.wg.Done()
		if d.err != nil {
			log.Printf(`data err %v`, d.err)
			return &Error{
				msg: d.err.Error(),
			}
		}
		err = proto.Unmarshal(d.msg, m.(proto.Message))
	case <-h.ctx.Done():
	}
	return
}
