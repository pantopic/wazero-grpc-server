package wazero_grpc_server

import (
	"google.golang.org/grpc/codes"
)

type Error struct {
	code codes.Code
	msg  string
}

func (e Error) Error() string {
	return e.msg
}
