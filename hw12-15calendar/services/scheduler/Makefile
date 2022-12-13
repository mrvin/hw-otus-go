build:
	CGO_ENABLED=0 go build -o ../../bin/scheduler -ldflags '-w'

run:
	go run -race main.go -config ../../configs/scheduler.yml

test:
	go test -v -race ./... -cover

lint:
	golangci-lint run ./...

