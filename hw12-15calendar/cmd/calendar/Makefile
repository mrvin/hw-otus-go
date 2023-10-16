build:
	go build -o ../../bin/calendar -ldflags '-w -s'

run:
	go run -race main.go

test:
	go test -race ../../internal/calendar/httpserver/handlers -cover

lint:
	golangci-lint run ./...

.PHONY: build run test lint
