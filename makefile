.PHONY: test check
test:
	go test -race -coverprofile=coverage.out -timeout 30s github.com/AleksandrMac/ucaller/test

check:
	golangci-lint run