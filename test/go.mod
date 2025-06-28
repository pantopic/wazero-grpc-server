module github.com/pantopic/wazero-grpc-server/test

go 1.24.3

replace github.com/pantopic/wazero-grpc-server/grpc-server-go v0.0.0 => ../grpc-server-go

require (
	github.com/pantopic/wazero-grpc-server/grpc-server-go v0.0.0
	google.golang.org/protobuf v1.36.6
)
