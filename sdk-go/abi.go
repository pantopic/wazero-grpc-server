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
		for k := range services[name].unary {
			methods = append(methods, `u.`+k)
		}
		for k := range services[name].clientStream {
			methods = append(methods, `c.`+k)
		}
		for k := range services[name].serverStream {
			methods = append(methods, `s.`+k)
		}
		for k := range services[name].bidirectionalRecv {
			methods = append(methods, `b.`+k)
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

func getCallOpts() (service *Service, method string) {
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
	method = parts[2]
	return
}

//export __grpc_server_unary
func __grpc_server_unary() {
	service, method := getCallOpts()
	h, ok := service.unary[method]
	if !ok {
		errCode = codes.Unimplemented
		return
	}
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
}

//export __grpc_server_client_stream
func __grpc_server_client_stream() {
	service, method := getCallOpts()
	h, ok := service.clientStream[method]
	if !ok {
		errCode = codes.Unimplemented
		return
	}
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
}

//export __grpc_server_server_stream
func __grpc_server_server_stream() {
	service, method := getCallOpts()
	h, ok := service.serverStream[method]
	if !ok {
		errCode = codes.Unimplemented
		return
	}
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
}

//export __grpc_server_bidirectional_recv
func __grpc_server_bidirectional_recv() {
	service, method := getCallOpts()
	h, ok := service.bidirectionalRecv[method]
	if !ok {
		errCode = codes.Unimplemented
		return
	}
	err := h(func(yield func([]byte) bool) {
		for {
			grpcRecv()
			if errCode != codes.OK {
				return
			}
			if !yield(getMsg()) {
				return
			}
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
}

//export __grpc_server_bidirectional_send
func __grpc_server_bidirectional_send() {
	service, method := getCallOpts()
	h, ok := service.bidirectionalSend[method]
	if !ok {
		errCode = codes.Unimplemented
		return
	}
	all, err := h()
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
}

//go:wasm-module pantopic/wazero-grpc-server
//export __host_grpc_server_send
func grpcSend()

//go:wasm-module pantopic/wazero-grpc-server
//export __host_grpc_server_recv
func grpcRecv()

// Fix for lint rule `unusedfunc`
var _ = __grpc_server
var _ = __grpc_server_unary
var _ = __grpc_server_client_stream
var _ = __grpc_server_server_stream
var _ = __grpc_server_bidirectional_recv
var _ = __grpc_server_bidirectional_send
