build:
	go build -o ../../bin/calendar-ws -ldflags '-w -s'

run:
	go run -race main.go

test:
	go test -race ./... -cover

lint:
	golangci-lint run ./...

.PHONY: build run test lint
