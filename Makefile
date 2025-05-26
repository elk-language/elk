header:
	go run ./cmd/headers/main.go -tags headers
	go fmt ./types/headers.go

generate:
	go generate ./...

elkify:
	go run ./cmd/elkify/main.go

fmt: generate header
	go fmt ./...

vet:
	go vet ./...

build: fmt
	go build -ldflags "-s -w"

test: header
	go test ./... -timeout 40s

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
