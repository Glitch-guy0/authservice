.PHONY: help build run test lint clean

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=auth-service

# Build the application
build:
	$(GOBUILD) -o bin/$(BINARY_NAME) ./cmd/app

# Run the application in development mode
run:
	$(GOCMD) run ./cmd/app

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -coverprofile=coverage.out -v ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Install dependencies
deps:
	$(GOGET) -v ./...

# Lint the code
lint:
	golangci-lint run --timeout 3m

# Clean build files
clean:
	$(GOCLEAN)
	rm -f bin/$(BINARY_NAME)
	rm -f coverage.out coverage.html

# Show help
help:
	@echo 'Available commands:'
	@echo '  build         - Build the application'
	@echo '  run           - Run the application in development mode'
	@echo '  test          - Run tests'
	@echo '  test-coverage - Run tests with coverage report'
	@echo '  deps          - Install dependencies'
	@echo '  lint          - Run linter'
	@echo '  clean         - Clean build files'
	@echo '  help          - Show this help message'
