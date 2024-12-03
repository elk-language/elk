header:
	go run ./cmd/headers/main.go -tags headers
	go fmt ./types/headers.go

fmt: header
	go fmt ./...

vet:
	go vet ./...

build: fmt
	go build .

test: header
	go test ./... -timeout 10s

repl: fmt
	go run . repl

typecheck: fmt
	go run . repl --typecheck

disassemble: fmt
	go run . repl --disassemble

lex: fmt
	go run . repl --lex

parse: fmt
	go run . repl --parse

inspect: fmt
	go run . repl --inspect-stack
