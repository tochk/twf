test:
	go test ./...

generate:
	go generate ./...

fmt:
	go fmt ./...

test_cover:
	go test ./... -cover

run:
	go build -o ./bin/example ./examples/
	./bin/example

all: fmt generate test run