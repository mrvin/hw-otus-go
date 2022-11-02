build:
	go build -o ../../bin/sender

run:
	go run -race main.go -config ../../configs/sender.yml

test:
	go test -race ./... -cover

lint:
	golangci-lint run ./...

