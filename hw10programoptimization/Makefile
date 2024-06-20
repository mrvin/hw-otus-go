install-easyjson:
	go install github.com/mailru/easyjson/easyjson@latest
codegen:
	easyjson -all stats_optimization.go
bench:
	go test -bench=. -benchmem -benchtime 10s -cpuprofile cpu.prof -memprofile mem.prof
test:
	go test -v -run TestGetDomainStat -cover
	go test -v -count=1 -timeout=30s -tags bench .
pprof:
	go tool pprof -http=":8090" cpu.prof
lint:
	golangci-lint run .
