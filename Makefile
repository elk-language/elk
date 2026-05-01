header: generate
	go run ./cmd/headers/main.go -tags headers
	go fmt ./types/headers.go

generate:
	go generate ./...

elkify:
	go run ./cmd/elkify/main.go

fmt: header
	go fmt ./...

vet:
	go vet ./...

build: fmt
	go build -ldflags "-s -w -X 'github.com/elk-language/elk/info.Version=$$(git describe --tags --exact-match 2>/dev/null || git branch --show-current)'"

go-test: header
	go test ./... -timeout 40s -tags debug

elk-test: header
	go run -tags debug ./cmd/elk test

test: go-test elk-test

repl: fmt
	go run -tags debug ./cmd/elk repl

typecheck: fmt
	go run -tags debug ./cmd/elk repl --typecheck

disassemble: fmt
	go run -tags debug ./cmd/elk repl --disassemble

transpile: fmt
	go run -tags debug ./cmd/elk repl --transpile

native: fmt
	go run -tags debug ./cmd/elk repl --native

lex: fmt
	go run -tags debug ./cmd/elk repl --lex

parse: fmt
	go run -tags debug ./cmd/elk repl --parse

expand: fmt
	go run -tags debug ./cmd/elk repl --expand

inspect: fmt
	go run -tags debug ./cmd/elk repl --inspect-stack

licenses:
	rm -rf licenses/
	go-licenses save ./... --save_path="licenses/"

tidy: licenses
	go mod tidy
