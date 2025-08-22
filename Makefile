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
	go build -ldflags "-s -w"

go-test: header
	go test ./... -timeout 40s

elk-test: header
	go run ./cmd/elk test

test: go-test elk-test

repl: fmt
	go run ./cmd/elk repl

typecheck: fmt
	go run ./cmd/elk repl --typecheck

disassemble: fmt
	go run ./cmd/elk repl --disassemble

lex: fmt
	go run ./cmd/elk repl --lex

parse: fmt
	go run ./cmd/elk repl --parse

expand: fmt
	go run ./cmd/elk repl --expand

inspect: fmt
	go run ./cmd/elk repl --inspect-stack

3rd-party-licenses:
	rm -rf licenses/
	go-licenses save ./... --save_path="licenses/"

tidy: 3rd-party-licenses
	go mod tidy
