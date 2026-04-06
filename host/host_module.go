package wazero_grpc_server

import (
	"context"
	"log"
	"strings"
	"sync"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pantopic/wazero-pool"
)

const Name = "pantopic/wazero-grpc-server"

var (
	ctxKeyMeta = Name + `/meta`
	ctxKeyNext = Name + `/next`
	ctxKeySend = Name + `/send`
)

type meta struct {
	ptrMethodCap uint32
	ptrMethodLen uint32
	ptrMethod    uint32
	ptrMsgCap    uint32
	ptrMsgLen    uint32
	ptrMsg       uint32
	ptrErrCode   uint32
}

type hostModule struct {
	sync.RWMutex

	module api.Module
}

func New(opts ...Option) (h *hostModule) {
	h = &hostModule{}
	for _, opt := range opts {
		opt(h)
	}
	return
}

func (h *hostModule) Name() string {
	return Name
}

// Register instantiates the host module, making it available to all module instances in this runtime
func (h *hostModule) Register(ctx context.Context, r wazero.Runtime) (err error) {
	builder := r.NewHostModuleBuilder(Name)
	register := func(name string, fn func(ctx context.Context, m api.Module, stack []uint64)) {
		builder = builder.NewFunctionBuilder().WithGoModuleFunction(api.GoModuleFunc(fn), nil, nil).Export(name)
	}
	for name, fn := range map[string]any{
		"__grpc_server_send": func(ctx context.Context, msg []byte, err error) {
			get[func([]byte, error)](ctx, ctxKeySend)(msg, err)
		},
	} {
		switch fn := fn.(type) {
		case func(context.Context, []byte, error):
			register(name, func(ctx context.Context, m api.Module, stack []uint64) {
				meta := get[*meta](ctx, ctxKeyMeta)
				fn(ctx, getMsgCopy(m, meta), getError(m, meta))
			})
		default:
			log.Panicf("Method signature implementation missing: %#v", fn)
		}
	}
	h.module, err = builder.Instantiate(ctx)
	return
}

// InitContext retrieves the meta page from the wasm module
func (h *hostModule) InitContext(ctx context.Context, m api.Module) (context.Context, error) {
	stack, err := m.ExportedFunction(`__grpc_server`).Call(ctx)
	if err != nil {
		return ctx, err
	}
	meta := &meta{}
	ptr := uint32(stack[0])
	for i, v := range []*uint32{
		&meta.ptrMethodCap,
		&meta.ptrMethodLen,
		&meta.ptrMethod,
		&meta.ptrMsgCap,
		&meta.ptrMsgLen,
		&meta.ptrMsg,
		&meta.ptrErrCode,
	} {
		*v = readUint32(m, ptr+uint32(4*i))
	}
	return context.WithValue(ctx, ctxKeyMeta, meta), nil
}

// ContextCopy populates dst context with the meta page from src context.
func (h *hostModule) ContextCopy(dst, src context.Context) context.Context {
	dst = context.WithValue(dst, ctxKeyMeta, get[*meta](src, ctxKeyMeta))
	if v := src.Value(ctxKeyNext); v != nil {
		dst = context.WithValue(dst, ctxKeyNext, v.(chan []byte))
	}
	if v := src.Value(ctxKeySend); v != nil {
		dst = context.WithValue(dst, ctxKeySend, v.(func([]byte, error)))
	}
	return dst
}

type ctxCopyFunc func(dst, src context.Context) context.Context

// RegisterServices attaches the grpc service(s) to the grpc server
// Called once before server open, usually given a module instance pool
func (h *hostModule) RegisterServices(ctx context.Context, s *grpc.Server, pool wazeropool.Instance, copy ...ctxCopyFunc) error {
	ctx = wazeropool.ContextSet(ctx, pool)
	copy = append(copy, wazeropool.ContextCopy)
	meta := get[*meta](ctx, ctxKeyMeta)
	pool.Run(func(mod api.Module) {
		// Format: msg = "/package1.ServiceName/u.method1,c.method2/service2.ServiceName/s.method1,b.method2"
		parts := strings.Split(string(getMsg(mod, meta)), "/")
		for i := 1; i+2 <= len(parts); i += 2 {
			h.registerService(s, pool, meta, parts[i], strings.Split(parts[i+1], ","), ctx, copy...)
		}
	})
	return nil
}

func (h *hostModule) registerService(s *grpc.Server, pool wazeropool.Instance, meta *meta, serviceName string, methods []string, ctx context.Context, copy ...ctxCopyFunc) {
	handler := &grpcHandler{pool, meta, ctx, copy}
	fakeDesc := &grpc.ServiceDesc{
		ServiceName: serviceName,
		HandlerType: (*any)(nil),
	}
	for _, m := range methods {
		parts := strings.Split(m, ".")
		if len(parts) < 2 {
			log.Panicf(`%s %#v`, methods, parts)
		}
		var d = grpc.StreamDesc{
			StreamName:    parts[1],
			ServerStreams: true,
			ClientStreams: true,
		}
		switch parts[0] {
		case "u":
			d.Handler = handler.handle(newHandlerFactoryUnary(h))
		case "c":
			d.Handler = handler.handle(newHandlerFactoryClientStream(h))
		case "s":
			d.Handler = handler.handle(newHandlerFactoryServerStream(h))
		case "b":
			d.Handler = handler.handle(newHandlerFactoryBidirectionalStream(h))
		}
		fakeDesc.Streams = append(fakeDesc.Streams, d)
	}
	s.RegisterService(fakeDesc, handler)
}

func (h *hostModule) Stop() (err error) {
	return
}

func get[T any](ctx context.Context, key string) T {
	v := ctx.Value(key)
	if v == nil {
		log.Panicf("Context item missing %s", key)
	}
	return v.(T)
}

func getErrCode(m api.Module, meta *meta) codes.Code {
	return codes.Code(readUint32(m, meta.ptrErrCode))
}

func setErrCode(m api.Module, meta *meta, code codes.Code) {
	writeUint32(m, meta.ptrErrCode, uint32(code))
}

func methodBuf(m api.Module, meta *meta) []byte {
	return read(m, meta.ptrMethod, 0, meta.ptrMethodCap)
}

func setMethod(m api.Module, meta *meta, method []byte) {
	copy(methodBuf(m, meta)[:len(method)], method)
	writeUint32(m, meta.ptrMethodLen, uint32(len(method)))
}

func getMsg(m api.Module, meta *meta) []byte {
	return read(m, meta.ptrMsg, meta.ptrMsgLen, meta.ptrMsgCap)
}

func getMsgCopy(m api.Module, meta *meta) []byte {
	return append([]byte(nil), getMsg(m, meta)...)
}

func msgBuf(m api.Module, meta *meta) []byte {
	return read(m, meta.ptrMsg, 0, meta.ptrMsgCap)
}

func setMsg(m api.Module, meta *meta, msg []byte) {
	copy(msgBuf(m, meta)[:len(msg)], msg)
	writeUint32(m, meta.ptrMsgLen, uint32(len(msg)))
}

func getError(m api.Module, meta *meta) error {
	c := getErrCode(m, meta)
	if c != codes.OK {
		return status.New(c, string(getMsg(m, meta))).Err()
	}
	return nil
}

func read(m api.Module, ptr, ptrLen, ptrCap uint32) (buf []byte) {
	buf, ok := m.Memory().Read(ptr, readUint32(m, ptrCap))
	if !ok {
		log.Panicf("Memory.Read(%d, %d) out of range", ptr, ptrCap)
	}
	return buf[:readUint32(m, ptrLen)]
}

func readUint32(m api.Module, ptr uint32) (val uint32) {
	val, ok := m.Memory().ReadUint32Le(ptr)
	if !ok {
		log.Panicf("Memory.Read(%d) out of range", ptr)
	}
	return
}

func writeUint32(m api.Module, ptr uint32, val uint32) {
	if ok := m.Memory().WriteUint32Le(ptr, val); !ok {
		log.Panicf("Memory.Read(%d) out of range", ptr)
	}
}
