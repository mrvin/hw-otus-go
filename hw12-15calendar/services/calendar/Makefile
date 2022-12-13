build:
	CGO_ENABLE=0 go build -o ../../bin/calendar -ldflags '-w'

run:
	go run -race main.go -config ../../configs/calendar.yml

test:
	go test -v -race ./... -cover

lint:
	golangci-lint run ./...

