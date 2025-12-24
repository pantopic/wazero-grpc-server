package grpc_server

import (
	"sort"
	"strings"
	"unsafe"

	"github.com/pantopic/wazero-grpc-server/sdk-go/codes"
)

var (
	meta             = make([]uint32, 7)
	methodCap uint32 = 256
	methodLen uint32
	msgCap    uint32 = 1.5 * 1024 * 1024
	msgLen    uint32
	errCode   codes.Code
	method    = make([]byte, int(methodCap))
	msg       = make([]byte, int(msgCap))
	services  = map[string]*Service{}
)

//export __grpc_server
func __grpc_server() (res uint32) {
	meta[0] = uint32(uintptr(unsafe.Pointer(&methodCap)))
	meta[1] = uint32(uintptr(unsafe.Pointer(&methodLen)))
	meta[2] = uint32(uintptr(unsafe.Pointer(&method[0])))
	meta[3] = uint32(uintptr(unsafe.Pointer(&msgCap)))
	meta[4] = uint32(uintptr(unsafe.Pointer(&msgLen)))
	meta[5] = uint32(uintptr(unsafe.Pointer(&msg[0])))
	meta[6] = uint32(uintptr(unsafe.Pointer(&errCode)))
	var serviceNames []string
	for k := range services {
		serviceNames = append(serviceNames, k)
	}
	sort.Strings(serviceNames)
	msg = msg[:0]
	for _, name := range serviceNames {
		msg = append(msg, []byte("/"+name+"/")...)
		var methods []string
		for k, h := range services[name].handlers {
			switch h.(type) {
			case unary:
				methods = append(methods, `u.`+k)
			case clientStream:
				methods = append(methods, `c.`+k)
			case serverStream:
				methods = append(methods, `s.`+k)
			}
		}
		sort.Strings(methods)
		msg = append(msg, []byte(strings.Join(methods, ","))...)
	}
	msgLen = uint32(len(msg))
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

func getHandler() (h any, errCode codes.Code) {
	m := string(getMethod())
	parts := strings.Split(m, "/")
	if len(parts) != 3 {
		errCode = codes.InvalidArgument
		setMsg([]byte(`Invalid method: ` + m))
		return
	}
	service, ok := services[parts[1]]
	if !ok {
		errCode = codes.Unimplemented
		return
	}
	h, ok = service.handlers[parts[2]]
	if !ok {
		errCode = codes.Unimplemented
		return
	}
	return
}

//export __grpc_server_call
func __grpc_server_call() {
	h, err := getHandler()
	if err > 0 {
		return
	}
	switch h := h.(type) {
	case unary:
		b, err := h(getMsg())
		if err != nil {
			if err2, ok := err.(Error); ok {
				errCode = err2.Code()
				setMsg([]byte(err2.Message()))
			} else {
				errCode = codes.Unknown
				setMsg([]byte(err.Error()))
			}
			return
		}
		errCode = codes.OK
		setMsg(b)
	case clientStream:
		b, err := h(func(yield func([]byte) bool) {
			for {
				if errCode != codes.OK {
					return
				}
				m := getMsg()
				if !yield(m) {
					return
				}
				grpcRecv()
			}
		})
		if err != nil {
			if err, ok := err.(Error); ok {
				errCode = err.Code()
			} else {
				errCode = codes.Unknown
			}
			setMsg([]byte(err.Error()))
			return
		}
		errCode = codes.OK
		setMsg(b)
	case serverStream:
		all, err := h(getMsg())
		if err != nil {
			if err, ok := err.(Error); ok {
				errCode = err.Code()
			} else {
				errCode = codes.Unknown
			}
			setMsg([]byte(err.Error()))
			return
		}
		errCode = codes.OK
		for m := range all {
			setMsg(m)
			grpcSend()
		}
	case bidirectionalStream:
		errCode = codes.Unimplemented
		return
	}
}

//go:wasm-module pantopic/wazero-grpc-server
//export Send
func grpcSend()

//go:wasm-module pantopic/wazero-grpc-server
//export Recv
func grpcRecv()

// Fix for lint rule `unusedfunc`
var _ = __grpc_server
var _ = __grpc_server_call
