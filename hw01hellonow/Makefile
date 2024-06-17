build:
	go build -ldflags '-w -s' -o hellonow

run:
	go run -race main.go

test:
	go test -race . -cover -v

lint:
	golangci-lint run .

.PHONY: build run test lint
