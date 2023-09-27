test:
	go test -race ./internal/storage/... -cover
	cd cmd/calendar && make test
build:
	docker-compose -f deployments/docker-compose.yaml build
build-calendar:
	docker-compose -f deployments/docker-compose.yaml build calendar
build-calendar-ws:
	docker-compose -f deployments/docker-compose.yaml build calendar-ws
up:
	docker-compose -f deployments/docker-compose.yaml --env-file deployments/postgres.env up
down:
	docker-compose -f deployments/docker-compose.yaml down
lint:
	golangci-lint run ./internal/...
codegen:
	go generate ./...
certgen:
	openssl req -x509 -newkey rsa:4096 -nodes \
		-keyout cert/key.pem -out cert/cert.pem \
		-subj "/C=RU/L=Saint Petersburg/O=Calendar Corp./OU=Calendar/CN=calendar.com/emailAddress=v.v.vinogradovv@gmail.com"

.PHONY: test build build-calendar build-calendar-ws up down lint codegen certgen

