test:
	go test ./...

generate:
	go generate ./...

fmt:
	go fmt ./...

run:
	go run examples/main.go

all: fmt generate test run