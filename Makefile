# Makefile for Native Hosting Assignment

.PHONY: all build run test clean

# Default target
all: build

# Build the application
build:
	mkdir -p ./bin
	go build -o ./bin/main -v ./cmd/main.go

# Run the application
run:
	make build
	./bin/main

# Run tests (if any)
test:
	@echo "Running tests..."
	go test ./...

# Clean build artifacts and deployment directories
clean:
	@echo "Cleaning..."
	go clean
	rm -f ./bin/main
	rm -rf ./deployments
	@echo "Clean complete."


