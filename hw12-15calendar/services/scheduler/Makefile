build:
	go build -o ../../bin/scheduler

run:
	go run -race main.go -config ../../configs/scheduler.yml

test:
	go test -race ./... -cover

lint:
	golangci-lint run ./...

