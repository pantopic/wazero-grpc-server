package grpc_server

func NewService(name string) *Service {
	services[name] = &Service{
		handlers: make(map[string]handler),
	}
	return services[name]
}

func (s *Service) AddMethod(name string, fn func([]byte) ([]byte, error)) {
	s.handlers[name] = fn
}
