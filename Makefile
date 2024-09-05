compile-headers:
	go run ./cmd/headers/main.go -tags headers
	gofmt -w ./types/headers.go

build: compile-headers
	go build .

test: compile-headers
	go test ./... -timeout 10s
