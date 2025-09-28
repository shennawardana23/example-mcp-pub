# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
BINARY_NAME=example-mcp-pub
BINARY_UNIX=$(BINARY_NAME)_unix

# Air parameters
AIR_BIN=$(shell if [ -f ~/go/bin/air ]; then echo ~/go/bin/air; elif command -v air >/dev/null 2>&1; then echo air; else echo ""; fi)

# Build flags
LDFLAGS=-ldflags "-w -s"

# Default target
.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: all
all: clean build test ## Clean, build and test

.PHONY: build
build: ## Build the binary
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v ./...

.PHONY: build-linux
build-linux: ## Build the binary for Linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_UNIX) -v ./...

.PHONY: test
test: ## Run tests
	$(GOTEST) -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	$(GOTEST) -race -coverprofile=coverage.out -covermode=atomic ./...

.PHONY: test-coverage-html
test-coverage-html: test-coverage ## Generate HTML coverage report
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

.PHONY: clean
clean: ## Clean build artifacts
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f coverage.out
	rm -f coverage.html

.PHONY: run
run: build ## Build and run the application
	./$(BINARY_NAME)

.PHONY: deps
deps: ## Download dependencies
	$(GOMOD) download
	$(GOMOD) verify

.PHONY: deps-update
deps-update: ## Update dependencies
	$(GOMOD) tidy
	$(GOGET) -u ./...

.PHONY: fmt
fmt: ## Format Go code
	$(GOFMT) ./...

.PHONY: vet
vet: ## Run go vet
	$(GOCMD) vet ./...

.PHONY: lint
lint: ## Run golangci-lint (requires golangci-lint to be installed)
	golangci-lint run

.PHONY: install-air
install-air: ## Install Air for hot reloading
	$(GOGET) -u github.com/air-verse/air
	$(GOCMD) install github.com/air-verse/air@latest

.PHONY: air-init
air-init: ## Initialize Air configuration
	$(AIR_BIN) init init

.PHONY: dev
dev: ## Start development server with Air hot reload
	@if [ -z "$(AIR_BIN)" ]; then \
		echo "Air is not installed. Installing Air..."; \
		$(MAKE) install-air; \
	fi
	@if [ ! -f .air.toml ]; then \
		echo "Air configuration not found. Creating default config..."; \
		$(MAKE) air-init; \
	fi
	$(AIR_BIN)

.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t $(BINARY_NAME) .

.PHONY: docker-run
docker-run: ## Run Docker container
	docker run -p 8080:8080 $(BINARY_NAME)

.PHONY: check
check: fmt vet test ## Run all checks (format, vet, test)