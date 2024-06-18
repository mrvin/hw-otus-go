build:
	go build -ldflags '-w -s' -o go-cp

test:
	go test -race -cover -v && \
	./test.sh

lint:
	golangci-lint run .

.PHONY: build test lint
