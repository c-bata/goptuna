VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' \
           -X 'main.revision=$(REVISION)'

GOIMPORTS ?= GO111MODULE=on goimports
GOCILINT ?= GO111MODULE=on golangci-lint
GO ?= GO111MODULE=on go
GODOC ?= GO111MODULE=on godoc

.DEFAULT_GOAL := help

PKGS := $(shell go list ./...)
SOURCES := $(shell find . -name "*.go" -not -name '*_test.go')
ENV := GO111MODULE=on

.PHONY: setup
setup:  ## Setup for required tools.
	go get -u golang.org/x/tools/cmd/goimports
	go get -u golang.org/x/tools/cmd/stringer
	go get golang.org/x/tools/cmd/godoc
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint


.PHONY: fmt
fmt: $(SOURCES) ## Formatting source codes.
	@$(GO) goimports -w $^

.PHONY: lint
lint: ## Run golint and go vet.
	@$(GOCILINT) run --no-config --disable-all --enable=goimports --enable=gocyclo --enable=govet --enable=misspell --enable=golint ./...

.PHONY: test
test:  ## Run tests with race condition checking.
	@$(GO) test -race ./...

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
