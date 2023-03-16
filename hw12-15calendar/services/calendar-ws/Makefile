build:
	go build -o ../../bin/calendar-ws -ldflags '-w'

run:
	go run -race main.go -config ../../configs/calendar-ws.yml

test:
	go test -v -race ./... -cover

lint:
	golangci-lint run ./...

