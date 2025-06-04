.PHONY: build run test clean dev lint fmt vet coverage help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=rss-feed-to-csv
BINARY_PATH=./bin/$(BINARY_NAME)

# Default target
all: test build

## build: Build the application
build:
	@echo "Building..."
	@mkdir -p bin
	@$(GOBUILD) -o $(BINARY_PATH) -v ./cmd/server

## run: Run the application
run: build
	@echo "Running..."
	@$(BINARY_PATH)

## test: Run all tests
test:
	@echo "Testing..."
	@$(GOTEST) -v ./...

## test-coverage: Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	@$(GOTEST) -coverprofile=coverage.out ./...
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## clean: Clean build files
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_PATH)
	@rm -f coverage.out coverage.html
	@$(GOCMD) clean

## dev: Run with hot reload using air
dev:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Installing air..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

## lint: Run linters
lint:
	@echo "Running linters..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi

## fmt: Format code
fmt:
	@echo "Formatting code..."
	@$(GOCMD) fmt ./...

## vet: Run go vet
vet:
	@echo "Running go vet..."
	@$(GOCMD) vet ./...

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	@$(GOMOD) download

## tidy: Tidy go modules
tidy:
	@echo "Tidying modules..."
	@$(GOMOD) tidy

## docker-build: Build docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -f build/docker/Dockerfile -t $(BINARY_NAME):latest .

## docker-run: Run docker container
docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --rm $(BINARY_NAME):latest

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^##' Makefile | sed -e 's/## /  /'