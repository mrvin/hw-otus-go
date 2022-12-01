build:
	go build -o ../../bin/calendar

run:
	go run -race main.go -config ../../configs/calendar.yml

test:
	go test -v -race ./... -cover

lint:
	golangci-lint run ./...

