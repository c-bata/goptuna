.DEFAULT_GOAL := help

PKGS := $(shell go list ./...)
SOURCES := $(shell find . -name "*.go" -not -name '*_test.go')

.PHONY: setup
setup:  ## Setup for required tools.
	go get -u golang.org/x/lint/golint
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/client9/misspell/cmd/misspell
	go get -u golang.org/x/tools/cmd/stringer
	go get golang.org/x/tools/cmd/godoc


.PHONY: fmt
fmt: $(SOURCES) ## Formatting source codes.
	@goimports -w $^

.PHONY: lint
lint: ## Run golint and go vet.
	@golint -set_exit_status=1 $(PKGS)
	@go vet $(PKGS)
	@misspell $(SOURCES)

.PHONY: test
test:  ## Run tests with race condition checking.
	@go test -race ./...

.PHONY: bench
bench:  ## Run benchmarks.
	@go test -bench=. -run=- -benchmem ./...

.PHONY: coverage
cover:  ## Run the tests.
	@go test -coverprofile=coverage.o ./...
	@go tool cover -func=coverage.o

.PHONY: godoc
godoc: ## Run godoc http server
	@echo "Please open http://localhost:6060/pkg/github.com/c-bata/goptuna/"
	godoc -http=localhost:6060

.PHONY: generate
generate: ## Run go generate
	@go generate ./...

.PHONY: build
build: ## Build example command lines.
	./_examples/build.sh

.PHONY: help
help: ## Show help text
	@echo "Commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'
