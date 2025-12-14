package wazero_grpc_server

import (
	"context"
	"log"
	"strings"
	"sync"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"google.golang.org/grpc"

	"github.com/pantopic/wazero-pool"
)

const Name = "pantopic/wazero-grpc-server"

var (
	DefaultCtxKeyMeta = `wazero_grpc_server_meta`
	DefaultCtxKeyNext = `wazero_grpc_next`
	DefaultCtxKeySend = `wazero_grpc_send`
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

	module     api.Module
	ctxKeyMeta string
	ctxKeySend string
	ctxKeyNext string
}

func New(opts ...Option) *hostModule {
	p := &hostModule{
		ctxKeyMeta: DefaultCtxKeyMeta,
		ctxKeyNext: DefaultCtxKeyNext,
		ctxKeySend: DefaultCtxKeySend,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (p *hostModule) Name() string {
	return Name
}

// Register instantiates the host module, making it available to all module instances in this runtime
func (p *hostModule) Register(ctx context.Context, r wazero.Runtime) (err error) {
	builder := r.NewHostModuleBuilder(Name)
	register := func(name string, fn func(ctx context.Context, m api.Module, stack []uint64)) {
		builder = builder.NewFunctionBuilder().WithGoModuleFunction(api.GoModuleFunc(fn), nil, nil).Export(name)
	}
	for name, fn := range map[string]any{
		"Recv": func(ctx context.Context) (msg []byte, ok bool) {
			msg, ok = <-get[chan []byte](ctx, p.ctxKeyNext)
			return
		},
		"Send": func(ctx context.Context, msg []byte, err error) {
			get[func([]byte, error)](ctx, p.ctxKeySend)(msg, err) // TODO - Replace copy with blocking call
		},
	} {
		switch fn := fn.(type) {
		case func(context.Context) ([]byte, bool):
			register(name, func(ctx context.Context, m api.Module, stack []uint64) {
				meta := get[*meta](ctx, p.ctxKeyMeta)
				b, ok := fn(ctx)
				if !ok {
					setErrCode(m, meta, errCodeDone)
					return
				}
				setErrCode(m, meta, errCodeEmpty)
				setMsg(m, meta, b)
			})
		case func(context.Context, []byte, error):
			register(name, func(ctx context.Context, m api.Module, stack []uint64) {
				meta := get[*meta](ctx, p.ctxKeyMeta)
				fn(ctx, msg(m, meta), getError(m, meta))
			})
		default:
			log.Panicf("Method signature implementation missing: %#v", fn)
		}
	}
	p.module, err = builder.Instantiate(ctx)
	return
}

// InitContext retrieves the meta page from the wasm module
func (p *hostModule) InitContext(ctx context.Context, m api.Module) (context.Context, error) {
	stack, err := m.ExportedFunction(`__grpcServerInit`).Call(ctx)
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
	return context.WithValue(ctx, p.ctxKeyMeta, meta), nil
}

// ContextCopy populates dst context with the meta page from src context.
func (h *hostModule) ContextCopy(src, dst context.Context) context.Context {
	dst = context.WithValue(dst, h.ctxKeyMeta, get[*meta](src, h.ctxKeyMeta))
	return dst
}

// RegisterServices attaches the grpc service(s) to the grpc server
// Called once before server open, usually given a module instance pool
func (p *hostModule) RegisterServices(ctx context.Context, s *grpc.Server, pool wazeropool.Instance) error {
	mod := pool.Get()
	defer pool.Put(mod)
	meta := get[*meta](ctx, p.ctxKeyMeta)
	// Format: msg = "/package1.ServiceName/u.method1,u.method2,c.method3/service2.ServiceName/u.method1,s.method2"
	parts := strings.Split(string(msg(mod, meta)), "/")
	for i := 1; i+2 <= len(parts); i += 2 {
		p.registerService(s, pool, meta, parts[i], strings.Split(parts[i+1], ","))
	}
	return nil
}

func (p *hostModule) registerService(s *grpc.Server, pool wazeropool.Instance, meta *meta, serviceName string, methods []string) {
	h := &grpcHandler{pool, meta}
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
			d.Handler = h.handler(newHandlerUnary)
		case "c":
			d.Handler = h.handler(newHandlerClientStream)
		case "s":
			d.Handler = h.handler(newHandlerServerStream)
		}
		fakeDesc.Streams = append(fakeDesc.Streams, d)
	}
	s.RegisterService(fakeDesc, h)
}

func (p *hostModule) Stop() (err error) {
	return
}

func get[T any](ctx context.Context, key string) T {
	v := ctx.Value(key)
	if v == nil {
		log.Panicf("Context item missing %s", key)
	}
	return v.(T)
}

func errCode(m api.Module, meta *meta) uint32 {
	return readUint32(m, meta.ptrErrCode)
}

func setErrCode(m api.Module, meta *meta, code uint32) {
	writeUint32(m, meta.ptrErrCode, uint32(code))
}

func methodBuf(m api.Module, meta *meta) []byte {
	return read(m, meta.ptrMethod, 0, meta.ptrMethodCap)
}

func setMethod(m api.Module, meta *meta, method []byte) {
	copy(methodBuf(m, meta)[:len(method)], method)
	writeUint32(m, meta.ptrMethodLen, uint32(len(method)))
}

func msg(m api.Module, meta *meta) []byte {
	return read(m, meta.ptrMsg, meta.ptrMsgLen, meta.ptrMsgCap)
}

func msgBuf(m api.Module, meta *meta) []byte {
	return read(m, meta.ptrMsg, 0, meta.ptrMsgCap)
}

func setMsg(m api.Module, meta *meta, msg []byte) {
	copy(msgBuf(m, meta)[:len(msg)], msg)
	writeUint32(m, meta.ptrMsgLen, uint32(len(msg)))
}

func getError(m api.Module, meta *meta) error {
	if err, ok := errorsByCode[errCode(m, meta)]; ok {
		return err
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
