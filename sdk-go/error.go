package grpc_server

import (
	"github.com/pantopic/wazero-grpc-server/sdk-go/codes"
)

type Error interface {
	Code() codes.Code
	Message() string
	Error() string
}
