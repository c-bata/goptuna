VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' \
           -X 'main.revision=$(REVISION)'

NO_EMBED_PKGS := $(shell go list -e ./... | grep -v -e dashboard -e cmd)
GOIMPORTS ?= goimports
GOCILINT ?= golangci-lint
GO ?= go
GODOC ?= godoc

.DEFAULT_GOAL := help

SOURCES := $(shell find . -name "*.go" | grep -v -e "sobol/direction_numbers.go" -e "dashboard/statik/statik.go")

.PHONY: fmt
fmt: $(SOURCES) ## Formatting source codes.
	@$(GOIMPORTS) -w $^

.PHONY: lint
lint: ## Run golint and go vet.
	@$(GOCILINT) run --no-config --disable-all --enable=goimports --enable=gocyclo --enable=govet --enable=misspell --enable=golint ./...

.PHONY: test
test:  ## Run tests with race condition checking.
	@$(GO) test -race $(NO_EMBED_PKGS)

.PHONY: bench
bench:  ## Run benchmarks.
	@$(GO) test -bench=. -run=- -benchmem ./...

.PHONY: coverage
cover:  ## Run the tests.
	@$(GO) test -coverprofile=coverage.o ./...
	@$(GO) tool cover -func=coverage.o

.PHONY: godoc
godoc: ## Run godoc http server
	@echo "Please open http://localhost:6060/pkg/github.com/c-bata/goptuna/"
	$(GODOC) -http=localhost:6060

.PHONY: generate
generate: ## Run go generate
	@$(GO) generate ./...

.PHONY: build
build: ## Build example command lines.
	mkdir -p ./bin/
	$(GO) build -o ./bin/goptuna -ldflags "$(LDFLAGS)" cmd/main.go
	./_examples/build.sh

.PHONY: help
help: ## Show help text
	@echo "Commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'
