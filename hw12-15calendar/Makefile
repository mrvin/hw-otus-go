test:
	go test -race ./internal/storage/... -cover && \
	go test -race ./internal/server/... -cover

build:
	docker-compose -f deployments/docker-compose.yaml build

up:
	docker-compose -f deployments/docker-compose.yaml up

down:
	docker-compose -f deployments/docker-compose.yaml down

lint:
	golangci-lint run ./internal/...

codegen:
	go generate ./...

