build:
	go build -ldflags '-w -s' -o go-envdir

test:
	go test -race -cover -v && \
	./test.sh

lint:
	golangci-lint run .

.PHONY: build test lint
