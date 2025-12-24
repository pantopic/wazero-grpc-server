package grpc_server

import (
	"iter"
)

type Service struct {
	handlers map[string]any
}

type unary func([]byte) ([]byte, error)
type clientStream func(iter.Seq[[]byte]) ([]byte, error)
type serverStream func([]byte) (iter.Seq[[]byte], error)
type bidirectionalStream func(iter.Seq[[]byte]) (iter.Seq[[]byte], error)

func NewService(name string) *Service {
	services[name] = &Service{
		handlers: make(map[string]any),
	}
	return services[name]
}

func (s *Service) Unary(name string, fn unary) *Service {
	s.handlers[name] = fn
	return s
}

func (s *Service) ClientStream(name string, fn clientStream) *Service {
	s.handlers[name] = fn
	return s
}

func (s *Service) ServerStream(name string, fn serverStream) *Service {
	s.handlers[name] = fn
	return s
}

func (s *Service) BidirectionalStream(name string, fn bidirectionalStream) *Service {
	s.handlers[name] = fn
	return s
}
