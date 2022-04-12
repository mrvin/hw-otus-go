bench:
	go test -bench=. -benchmem -benchtime 10s
test:
	go test -v -count=1 -timeout=30s -tags bench .
lint:
	golangci-lint run .
cover:
	go test . -cover

