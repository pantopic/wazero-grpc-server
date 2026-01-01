package grpc_server

import (
	"iter"
)

type Service struct {
	unary             map[string]unary
	clientStream      map[string]clientStream
	serverStream      map[string]serverStream
	bidirectionalRecv map[string]bidirectionalRecv
	bidirectionalSend map[string]bidirectionalSend
}

type unary func([]byte) ([]byte, error)
type clientStream func(iter.Seq[[]byte]) ([]byte, error)
type serverStream func([]byte) (iter.Seq[[]byte], error)
type bidirectionalRecv func(iter.Seq[[]byte]) error
type bidirectionalSend func() (iter.Seq[[]byte], error)

func NewService(name string) *Service {
	services[name] = &Service{
		unary:             make(map[string]unary),
		clientStream:      make(map[string]clientStream),
		serverStream:      make(map[string]serverStream),
		bidirectionalRecv: make(map[string]bidirectionalRecv),
		bidirectionalSend: make(map[string]bidirectionalSend),
	}
	return services[name]
}

func (s *Service) Unary(name string, fn unary) *Service {
	s.unary[name] = fn
	return s
}

func (s *Service) ClientStream(name string, fn clientStream) *Service {
	s.clientStream[name] = fn
	return s
}

func (s *Service) ServerStream(name string, fn serverStream) *Service {
	s.serverStream[name] = fn
	return s
}

func (s *Service) BidirectionalStream(name string, recv bidirectionalRecv, send bidirectionalSend) *Service {
	s.bidirectionalRecv[name] = recv
	s.bidirectionalSend[name] = send
	return s
}
