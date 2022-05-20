build:
	go build -o ./bin/calendar ./cmd/calendar

test:
	go test -race ./internal/... # ./pkg/...

up:
	docker-compose -f deployments/docker-compose.yaml up

lint:
	golangci-lint run ./...

