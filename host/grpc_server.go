package wazero_grpc_server

import (
	"context"
	"log"
	"strings"
	"sync"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"google.golang.org/grpc"
)

var (
	DefaultCtxKeyMeta   = `wazero_grpc_server_meta_key`
	DefaultCtxKeyServer = `wazero_grpc_server`
)

type meta struct {
	ptrMethodMax uint32
	ptrMethodLen uint32
	ptrMsgMax    uint32
	ptrMsgLen    uint32
	ptrErrCode   uint32
	ptrMethod    uint32
	ptrMsg       uint32
}

type module struct {
	sync.RWMutex

	module       api.Module
	ctxKeyMeta   string
	ctxKeyServer string
}

func New(opts ...Option) *module {
	p := &module{
		ctxKeyMeta:   DefaultCtxKeyMeta,
		ctxKeyServer: DefaultCtxKeyServer,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Register instantiates the host module, making it available to all module instances in this runtime
// Called once after a runtime is created, usually on startup
func (p *module) Register(ctx context.Context, r wazero.Runtime) (err error) {
	builder := r.NewHostModuleBuilder("grpc")
	// TODO - Add callbacks (i.e. Watch message emission)
	p.module, err = builder.Instantiate(ctx)
	return
}

// InitContext populates the meta page in context for a given module instance
// Called per module instance immediately after module instantiation
func (p *module) InitContext(ctx context.Context, m api.Module) (context.Context, error) {
	stack, err := m.ExportedFunction(`grpc`).Call(ctx)
	if err != nil {
		return ctx, err
	}
	meta := &meta{}
	ptr := uint32(stack[0])
	meta.ptrMethodMax, _ = m.Memory().ReadUint32Le(ptr)
	meta.ptrMethodLen, _ = m.Memory().ReadUint32Le(ptr + 4)
	meta.ptrMethod, _ = m.Memory().ReadUint32Le(ptr + 8)
	meta.ptrMsgMax, _ = m.Memory().ReadUint32Le(ptr + 12)
	meta.ptrMsgLen, _ = m.Memory().ReadUint32Le(ptr + 16)
	meta.ptrMsg, _ = m.Memory().ReadUint32Le(ptr + 20)
	meta.ptrErrCode, _ = m.Memory().ReadUint32Le(ptr + 24)
	return context.WithValue(ctx, p.ctxKeyMeta, meta), nil
}

// RegisterService attaches the grpc service(s) to the grpc server
// Called once before server open, usually given a module instance pool
func (p *module) RegisterService(ctx context.Context, s *grpc.Server, m api.Module) context.Context {
	meta := get[*meta](ctx, p.ctxKeyMeta)
	// msg = "/service.1.name/method1,method2,method3/service.2.name/method1,method2"
	parts := strings.Split(string(msg(m, meta)), "/")
	for i := 1; i+2 <= len(parts); i += 2 {
		registerService(s, m, meta, parts[i], strings.Split(parts[i+1], ","))
	}
	return context.WithValue(ctx, p.ctxKeyServer, s)
}

func (p *module) Stop() (err error) {
	return
}

func (p *module) server(ctx context.Context) *grpc.Server {
	return get[*grpc.Server](ctx, p.ctxKeyServer)
}

func get[T any](ctx context.Context, key string) T {
	v := ctx.Value(key)
	if v == nil {
		log.Panicf("Context item missing %s", key)
	}
	return v.(T)
}

func method(m api.Module, meta *meta) []byte {
	return read(m, meta.ptrMethod, meta.ptrMethodLen, meta.ptrMethodMax)
}

func errCode(m api.Module, meta *meta) uint32 {
	return readUint32(m, meta.ptrErrCode)
}

func methodBuf(m api.Module, meta *meta) []byte {
	return read(m, meta.ptrMethod, 0, meta.ptrMethodMax)
}

func setMethod(m api.Module, meta *meta, method []byte) {
	copy(methodBuf(m, meta)[:len(method)], method)
	writeUint32(m, meta.ptrMethodLen, uint32(len(method)))
}

func msg(m api.Module, meta *meta) []byte {
	return read(m, meta.ptrMsg, meta.ptrMsgLen, meta.ptrMsgMax)
}

func msgBuf(m api.Module, meta *meta) []byte {
	return read(m, meta.ptrMsg, 0, meta.ptrMsgMax)
}

func setMsg(m api.Module, meta *meta, msg []byte) {
	copy(msgBuf(m, meta)[:len(msg)], msg)
	writeUint32(m, meta.ptrMsgLen, uint32(len(msg)))
}

func getError(m api.Module, meta *meta) *Error {
	if err, ok := errorsByCode[errCode(m, meta)]; ok {
		return err
	}
	return nil
}

func read(m api.Module, ptrData, ptrLen, ptrMax uint32) (buf []byte) {
	buf, ok := m.Memory().Read(ptrData, readUint32(m, ptrMax))
	if !ok {
		log.Panicf("Memory.Read(%d, %d) out of range", ptrData, ptrLen)
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
