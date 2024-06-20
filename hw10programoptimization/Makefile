install-easyjson:
	go install github.com/mailru/easyjson/easyjson@latest
codegen:
	easyjson -all stats_optimization.go
bench:
	go test -bench=. -benchmem -benchtime 10s
test:
	go test -v -run TestGetDomainStat -cover
	go test -v -count=1 -timeout=30s -tags bench .
lint:
	golangci-lint run .
