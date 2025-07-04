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

	services = map[string]*Service{}
)

//export __grpc
func __grpc() (res uint32) {
	meta[0] = uint32(uintptr(unsafe.Pointer(&methodMax)))
	meta[1] = uint32(uintptr(unsafe.Pointer(&methodLen)))
	meta[2] = uint32(uintptr(unsafe.Pointer(&method[0])))
	meta[3] = uint32(uintptr(unsafe.Pointer(&msgMax)))
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
			case handler:
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

//export grpcCall
func grpcCall() {
	m := string(getMethod())
	parts := strings.Split(m, "/")
	if len(parts) != 3 {
		errCode = errCodeInvalid
		setMsg([]byte(`Invalid method: ` + m))
		return
	}
	service, ok := services[parts[1]]
	if !ok {
		errCode = errCodeNotImplemented
		return
	}
	h, ok := service.handlers[parts[2]]
	if !ok {
		errCode = errCodeNotImplemented
		return
	}
	switch h := h.(type) {
	case handler:
		b, err := h(getMsg())
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
	case clientStream:
		b, err := h(func(yield func([]byte) bool) {
			for {
				if errCode != errCodeEmpty {
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
				errCode = err.code
			} else {
				errCode = errCodeUnknown
			}
			setMsg([]byte(err.Error()))
			return
		}
		errCode = errCodeEmpty
		setMsg(b)
	case serverStream:
		all, err := h(getMsg())
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
		for m := range all {
			setMsg(m)
			grpcSend()
		}
	case bidirectionalStream:
		errCode = errCodeNotImplemented
		return
	}
}

//go:wasm-module grpc
//export Recv
func grpcRecv()

//go:wasm-module grpc
//export Send
func grpcSend()

// Fix for lint rule `unusedfunc`
var _ = __grpc
var _ = grpcCall
