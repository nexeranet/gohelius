-include .env
export

# load env to shell
# export $(grep -v '^#' .env | xargs)

BUILD_NAME_API=gohelius

.PHONY: prepare
prepare:
	@go clean
	@go fmt ./...

.PHONY: test
test:
	@go clean -testcache
	@go test ./...

.PHONY: lint
lint:
	@golangci-lint cache clean
	@golangci-lint run ./...
