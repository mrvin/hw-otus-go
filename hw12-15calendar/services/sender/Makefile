build:
	CGO_ENABLED=0 go build -o ../../bin/sender  -ldflags '-w'

run:
	go run -race main.go -config ../../configs/sender.yml

test:
	go test -v -race ./... -cover

lint:
	golangci-lint run ./...

