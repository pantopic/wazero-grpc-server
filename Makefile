dev:
	@go build -ldflags="-s -w" -o _dist/pantopic && cd cmd/standalone && docker compose up --build

build:
	@go build -ldflags="-s -w" -o _dist/pantopic

wasm-server:
	@cd test-server && tinygo build -buildmode=wasi-legacy -target=wasi -opt=2 -gc=conservative -scheduler=none -o ../host/test-server.wasm

wasm-server-prod:
	@cd test-server && tinygo build -buildmode=wasi-legacy -target=wasi -opt=2 -gc=conservative -scheduler=none -o ../host/test-server.prod.wasm -no-debug

wasm-server-easy:
	@cd test-server-easy && tinygo build -buildmode=wasi-legacy -target=wasi -opt=2 -gc=conservative -scheduler=none -o ../host/test-server-easy.wasm

wasm-server-easy-prod:
	@cd test-server-easy && tinygo build -buildmode=wasi-legacy -target=wasi -opt=2 -gc=conservative -scheduler=none -o ../host/test-server-easy.prod.wasm -no-debug

wasm-server-lite:
	@cd test-server-lite && tinygo build -buildmode=wasi-legacy -target=wasi -opt=2 -gc=conservative -scheduler=none -o ../host/test-server-lite.wasm

wasm-server-lite-prod:
	@cd test-server-lite && tinygo build -buildmode=wasi-legacy -target=wasi -opt=2 -gc=conservative -scheduler=none -o ../host/test-server-lite.prod.wasm -no-debug

wasm-all: wasm-server wasm-server-easy wasm-server-lite wasm-server-prod wasm-server-easy-prod wasm-server-lite-prod

test:
	@cd host && go test .

cover:
	@mkdir -p _dist
	@cd host && go test . -coverprofile=../_dist/coverage.out -v
	@go tool cover -html=_dist/coverage.out -o _dist/coverage.html

gen:
	@protoc test-server/*.proto --go_out=test-server \
		--go_opt=paths=source_relative -I test-server

gen-install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

gen-lite:
	@ protoc test-server-lite/*.proto \
		--plugin protoc-gen-go-lite="${GOBIN}/protoc-gen-go-lite" \
		--go-lite_out=. \
		--go-lite_opt=features=marshal+unmarshal+size \
		--go-lite_opt=paths=source_relative

gen-lite-install:
	go install github.com/aperturerobotics/protobuf-go-lite/cmd/protoc-gen-go-lite@latest

cloc:
	@cloc . --exclude-dir=_example,_dist,internal,cmd --exclude-ext=pb.go

.PHONY: all test clean
