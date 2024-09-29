header:
	go run ./cmd/headers/main.go -tags headers
	go fmt ./types/headers.go

fmt: header
	go fmt ./...

vet:
	go vet ./...

build: fmt vet
	go build .

test: header
	go test ./... -timeout 10s

repl: fmt vet
	go run . repl

typecheck: fmt vet
	go run . repl --typecheck

disassemble: fmt vet
	go run . repl --disassemble

lex: fmt vet
	go run . repl --lex

parse: fmt vet
	go run . repl --parse

inspect: fmt vet
	go run . repl --inspect-stack
