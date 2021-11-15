LIST_ALL := $(shell go list ./... | grep -v vendor | grep -v mocks)

all: lint test

.PHONY: lint
lint: ## Lint the files
	@go fmt ${LIST_ALL}
	@golangci-lint version
	@golangci-lint run

.PHONY: test
test: ## Run unit tests
	@go test -short -count 1 -v ./...

.PHONY: race
race: ## Run data race detector
	@go test -race -short -count 1 -v ./...

.PHONY: coverage
coverage: ## Generate coverage report
	@go-acc ./...
	@go tool cover -func=coverage.txt

.PHONY: depcheck
depcheck: ## Check dependencies for vulnerabilities
	@go list -json -deps ./... | nancy sleuth

.PHONY: upgrade
upgrade: ## Upgrade the dependencies
	@go get -u -t ./...
	@go mod tidy
	@go mod vendor

.PHONY: licensecheck
licensecheck: ## Check dependencies for forbidden licenses
	@go-licenses check ./...

.PHONY: clean
clean: ## Remove outdated file and empty cache
	@rm -rf "$(go env GOCACHE)"
	@rm -f coverage.*
	@rm -f ./bin/tilores-cli

.PHONY: help
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
