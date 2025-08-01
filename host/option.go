package wazero_grpc_server

type Option func(*hostModule)

func WithCtxKeyMeta(key string) Option {
	return func(p *hostModule) {
		p.ctxKeyMeta = key
	}
}
