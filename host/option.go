package wazero_grpc_server

type Option func(*module)

func WithCtxKeyMeta(key string) Option {
	return func(p *module) {
		p.ctxKeyMeta = key
	}
}
func WithCtxKeyServer(key string) Option {
	return func(p *module) {
		p.ctxKeyServer = key
	}
}
