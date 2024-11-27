# Variables
BINARY_NAME := lazyblockchain
SRC_DIR := .
BUILD_DIR := ./bin

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC_DIR)
	@echo "Binary built at $(BUILD_DIR)/$(BINARY_NAME)"

# Run the binary
.PHONY: run
run:
	@echo "Running $(BINARY_NAME)..."
	@go run .


# Clean build artifacts
.PHONY: clean
clean:
	@echo "removing bin/lazyblockchain"
	@rm -rf $(BUILD_DIR)
	@echo "removed."


# Help message
.PHONY: help
help:
	@echo "Makefile for lazyblockchain"
	@echo "Usage:"
	@echo "  make build       - Build the binary"
	@echo "  make run         - Run the binary"
	@echo "  make clean       - Clean build artifacts"
