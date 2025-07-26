SHELL := /bin/bash

# Binary name
BINARY_NAME=osc2hue

# Build the application
build:
	go build -o ${BINARY_NAME} .

# Run the application
run:
	go run .

# Clean build artifacts
clean:
	go clean
	rm -f ${BINARY_NAME}*

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Install dependencies
deps:
	go mod tidy
	go mod download

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o ${BINARY_NAME}-linux-amd64 .
	GOOS=windows GOARCH=amd64 go build -o ${BINARY_NAME}-windows-amd64.exe .
	GOOS=darwin GOARCH=amd64 go build -o ${BINARY_NAME}-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o ${BINARY_NAME}-darwin-arm64 .

# Install the binary
install:
	go install .

# Show help
help:
	@echo "Available commands:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage"
	@echo "  deps         - Install dependencies"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  build-all    - Build for multiple platforms"
	@echo "  install      - Install the binary"
	@echo "  help         - Show this help message"

.PHONY: build run clean test test-coverage deps fmt lint build-all install help
