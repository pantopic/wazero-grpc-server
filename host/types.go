package wazero_grpc_server

import (
	"context"

	"github.com/tetratelabs/wazero/api"
	"google.golang.org/grpc"
)

type handlerFactory func(context.Context, api.Module, *meta, string) grpc.ClientStream
