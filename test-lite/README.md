# test-lite

The `test-lite` package uses [protobuf-go-lite](github.com/aperturerobotics/protobuf-go-lite) for protobuf serialization.

This results in a WASM binary that is 5x smaller than the standard `protoc-go-gen` for this basic use case.
