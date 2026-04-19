package grpc_server

type Option func()

func WithBufferCapMethod(cap uint32) Option {
	return func() {
		methodCap = cap
	}
}

func WithBufferCapMsg(cap uint32) Option {
	return func() {
		msgCap = cap
	}
}
