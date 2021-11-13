.PHONY: check-lint
check-lint:
	@echo 'check-lint TBD'

.PHONY: unit-test
unit-test:
	@go test ./... -coverprofile=cp.out && go tool cover -func=cp.out && rm -f cp.out
