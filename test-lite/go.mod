module github.com/pantopic/wazero-grpc-server/test-lite

go 1.24.3

replace github.com/pantopic/wazero-grpc-server/grpc-server-go v0.0.0 => ../grpc-server-go

require (
	github.com/aperturerobotics/protobuf-go-lite v0.9.1
	github.com/pantopic/wazero-grpc-server/grpc-server-go v0.0.0
)
