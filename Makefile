test:
	go test ./...

generate:
	go generate ./...

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

lint_fix:
	golangci-lint run --fix ./...

test_cover:
	go test ./... -cover

run_example:
	go build -o ./bin/example ./examples/
	./bin/example

all: fmt generate lint test