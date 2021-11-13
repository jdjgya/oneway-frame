.PHONY: check-lint
check-lint:
	@golangci-lint run

.PHONY: unit-test
unit-test:
	@go test ./... -coverprofile=cp.out && go tool cover -func=cp.out && rm -f cp.out
