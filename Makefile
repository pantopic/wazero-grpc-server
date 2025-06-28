dev:
	@go build -ldflags="-s -w" -o _dist/pantopic && cd cmd/standalone && docker compose up --build

build:
	@go build -ldflags="-s -w" -o _dist/pantopic

wasm:
	@cd test && tinygo build -buildmode=wasi-legacy -target=wasi -opt=2 -gc=conservative -scheduler=none -o ../host/test.wasm

wasm-prod:
	@cd test && tinygo build -buildmode=wasi-legacy -target=wasi -opt=2 -gc=conservative -scheduler=none -o ../host/test.prod.wasm -no-debug

wasm-easy:
	@cd test-easy && tinygo build -buildmode=wasi-legacy -target=wasi -opt=2 -gc=conservative -scheduler=none -o ../host/test-easy.wasm

wasm-easy-prod:
	@cd test-easy && tinygo build -buildmode=wasi-legacy -target=wasi -opt=2 -gc=conservative -scheduler=none -o ../host/test-easy.prod.wasm -no-debug

wasm-lite:
	@cd test-lite && tinygo build -buildmode=wasi-legacy -target=wasi -opt=2 -gc=conservative -scheduler=none -o ../host/test-lite.wasm

wasm-lite-prod:
	@cd test-lite && tinygo build -buildmode=wasi-legacy -target=wasi -opt=2 -gc=conservative -scheduler=none -o ../host/test-lite.prod.wasm -no-debug

wasm-all: wasm wasm-easy wasm-lite wasm-prod wasm-easy-prod wasm-lite-prod

test:
	@cd host && go test .

cover:
	@mkdir -p _dist
	@cd host && go test . -coverprofile=../_dist/coverage.out -v
	@go tool cover -html=_dist/coverage.out -o _dist/coverage.html

gen:
	@protoc test/*.proto --go_out=test \
		--go_opt=paths=source_relative -I test

gen-install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

gen-lite:
	@ protoc test-lite/*.proto \
		--plugin protoc-gen-go-lite="${GOBIN}/protoc-gen-go-lite" \
		--go-lite_out=. \
		--go-lite_opt=features=marshal+unmarshal+size \
		--go-lite_opt=paths=source_relative

gen-lite-install:
	go install github.com/aperturerobotics/protobuf-go-lite/cmd/protoc-gen-go-lite@latest

cloc:
	@cloc . --exclude-dir=_example,_dist,internal,cmd --exclude-ext=pb.go

.PHONY: all test clean
