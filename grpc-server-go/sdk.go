package grpc_server

import (
	"iter"
)

type Service struct {
	handlers map[string]any
}

type handler func([]byte) ([]byte, error)
type clientStream func(iter.Seq[[]byte]) ([]byte, error)
type serverStream func([]byte) (iter.Seq[[]byte], error)
type bidirectionalStream func(iter.Seq[[]byte]) (iter.Seq[[]byte], error)

func NewService(name string) *Service {
	services[name] = &Service{
		handlers: make(map[string]any),
	}
	return services[name]
}

func (s *Service) Unary(name string, fn handler) {
	s.handlers[name] = fn
}

func (s *Service) ClientStream(name string, fn clientStream) {
	s.handlers[name] = fn
}

func (s *Service) ServerStream(name string, fn serverStream) {
	s.handlers[name] = fn
}

// func (s *Service) BidirectionalStream(name string, fn bidirectionalStream) {
// 	s.handlers[name] = fn
// }
