package grpc_server

type Option func()

func WithBufferCap(method, msg uint32) Option {
	return func() {
		methodCap = method
		msgCap = msg
	}
}
