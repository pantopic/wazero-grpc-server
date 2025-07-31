package grpc_server

//go:wasm-module pantopic/wazero-grpc-server
//export Recv
func grpcRecv()

//go:wasm-module pantopic/wazero-grpc-server
//export Send
func grpcSend()
