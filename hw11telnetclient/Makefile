build:
	go build -ldflags '-w -s' -o go-telnet

test:
	go test -race -cover -v ./telnet/
	./test.sh

lint:
	golangci-lint run ./...

.PHONY: build test lint
