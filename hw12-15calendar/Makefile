build:
	go build -o ./bin/calendar ./cmd/calendar

run:
	go run cmd/calendar/main.go -config configs/config.yml

test:
	go test -race ./cmd/calendar/... -cover && \
	go test -race ./internal/storage/... -cover && \
	go test -race ./internal/server/... -cover

up:
	docker-compose -f deployments/docker-compose.yaml up

lint:
	golangci-lint run ./...

codegen:
	go generate ./...

