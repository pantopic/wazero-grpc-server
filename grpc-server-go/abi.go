package grpc_server

import (
	"sort"
	"strings"
	"unsafe"
)

var (
	methodMax uint32 = 256
	methodLen uint32
	msgMax    uint32 = 1.5 * 1024 * 1024
	msgLen    uint32
	errCode   uint32
	method    = make([]byte, int(methodMax))
	msg       = make([]byte, int(msgMax))
	meta      = make([]uint32, 7)

	services = map[string]*service{}
)

type service struct {
	handlers map[string]handler
}

type handler func([]byte) ([]byte, error)

//export grpc
func grpc() (res uint32) {
	meta[0] = uint32(uintptr(unsafe.Pointer(&methodMax)))
	meta[1] = uint32(uintptr(unsafe.Pointer(&methodLen)))
	meta[2] = uint32(uintptr(unsafe.Pointer(&msgMax)))
	meta[3] = uint32(uintptr(unsafe.Pointer(&msgLen)))
	meta[4] = uint32(uintptr(unsafe.Pointer(&errCode)))
	meta[5] = uint32(uintptr(unsafe.Pointer(&msg[0])))
	meta[6] = uint32(uintptr(unsafe.Pointer(&method[0])))
	msg = msg[:0]
	var serviceNames []string
	for k := range services {
		serviceNames = append(serviceNames, k)
	}
	sort.Strings(serviceNames)
	for _, name := range serviceNames {
		msg = append(msg, []byte("/"+name+"/")...)
		var methods []string
		for k := range services[name].handlers {
			methods = append(methods, k)
		}
		sort.Strings(methods)
		msg = append(msg, []byte(strings.Join(methods, ","))...)
	}
	return uint32(uintptr(unsafe.Pointer(&meta[0])))
}

func setMsg(b []byte) {
	copy(msg[:len(b)], b)
	msgLen = uint32(len(b))
}

func getMsg() []byte {
	return msg[:msgLen]
}

func getMethod() []byte {
	return method[:methodLen]
}

//export grpcCall
func grpcCall() {
	parts := strings.Split(string(getMethod()), "/")
	if len(parts) < 3 {
		errCode = errCodeUnrecognized
		return
	}
	service, ok := services[parts[1]]
	if !ok {
		errCode = errCodeNotImplemented
		return
	}
	handler, ok := service.handlers[parts[2]]
	if !ok {
		errCode = errCodeNotImplemented
		return
	}
	b, err := handler(getMsg())
	if err != nil {
		if err, ok := err.(Error); ok {
			errCode = err.code
		} else {
			errCode = errCodeUnknown
		}
		setMsg([]byte(err.Error()))
		return
	}
	errCode = errCodeEmpty
	setMsg(b)
}

// Fix for lint rule `unusedfunc`
var _ = grpc
var _ = grpcCall
