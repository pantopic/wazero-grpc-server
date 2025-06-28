package grpc_server

func NewService(name string) *service {
	services[name] = &service{
		handlers: make(map[string]handler),
	}
	return services[name]
}

func (s *service) AddMethod(name string, fn func([]byte) ([]byte, error)) {
	s.handlers[name] = fn
}
