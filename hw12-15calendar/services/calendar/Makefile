build:
	go build -o ../../bin/calendar -ldflags '-w'

run:
	go run -race main.go -config ../../configs/calendar.yml

test:
	go test -race ./httpserver/handler -cover

lint:
	golangci-lint run ./...

PHONY: build run test lint
