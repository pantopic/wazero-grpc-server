module github.com/pantopic/wazero-grpc-server/test-easy

go 1.24.3

replace github.com/pantopic/wazero-grpc-server/grpc-server-go v0.0.0 => ../grpc-server-go

require (
	github.com/VictoriaMetrics/easyproto v0.1.4
	github.com/pantopic/wazero-grpc-server/grpc-server-go v0.0.0
)
