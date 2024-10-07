# Makefile for Trail Finder

# Variables
BINARY_NAME=trailfinder
OUTPUT_DIR=bin

# Default target
all: build

# Build the project
build:
	go build -o $(OUTPUT_DIR)/$(BINARY_NAME) -ldflags '-extldflags "-static"' main.go

# Run tests
test:
	go test ./...

# Clean up build artifacts
clean:
	rm -f $(OUTPUT_DIR)/$(BINARY_NAME)  # Clean from the output directory

.PHONY: all build test clean
