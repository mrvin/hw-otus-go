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
	openssl req -nodes -x509 -newkey rsa:4096 \
		-keyout cert/serverKey.pem -out cert/serverCert.pem -days 365 \
		-subj "/C=RU/L=Saint Petersburg/O=Calendar Corp./OU=Calendar/CN=calendar.com/emailAddress=v.v.vinogradovv@gmail.com"
	openssl req -nodes -x509 -newkey rsa:4096 \
		-keyout cert/clientKey.pem -out cert/clientCert.pem -days 365 \
		-subj "/C=RU/L=Saint Petersburg/O=Calendar Corp./OU=Calendar/CN=calendar.com/emailAddress=v.v.vinogradovv@gmail.com"

.PHONY: test build build-calendar build-calendar-ws up down lint codegen certgen

