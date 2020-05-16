test:
	go test ./...

generate:
	go generate ./...

fmt:
	go fmt ./...

run:
	go build -o ./bin/example ./examples/
	./bin/example

all: fmt generate test run