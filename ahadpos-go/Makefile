.PHONY: all build run clean test deps

# Variables
APP_NAME=ahadpos-api
CMD_DIR=cmd/api
BUILD_DIR=build
GO_FILES=$(shell find . -name '*.go' -not -path './vendor/*')

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

all: clean deps build

## deps: Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

## build: Build the application
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME) ./$(CMD_DIR)
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

## run: Run the application
run:
	@echo "Running $(APP_NAME)..."
	$(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME) ./$(CMD_DIR)
	./$(BUILD_DIR)/$(APP_NAME)

## test: Run tests
test:
	$(GOTEST) -v ./...

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	@echo "Clean complete"

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
