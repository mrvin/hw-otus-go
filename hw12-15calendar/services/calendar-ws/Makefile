build:
	go build -o ../../bin/calendar-ws

run:
	go run -race main.go -config ../../configs/calendar-ws.yml

test:
	go test -race ./... -cover

lint:
	golangci-lint run ./...

