.PHONY: build run clean deploy-local

# Binary name
BINARY_NAME=mcp-compliance

# Build directory
BUILD_DIR=bin

# Local deployment directory
LOCAL_DEPLOY_DIR=$(HOME)/.mcp-compliance/bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Main build target
build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/mcp-compliance

# Run the server
run: build
	$(BUILD_DIR)/$(BINARY_NAME)

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Run tests
test:
	$(GOTEST) -v ./...

# Get dependencies
deps:
	$(GOGET) -v ./...

# Build for all platforms
build-all: build-linux build-windows build-macos

# Build for Linux
build-linux:
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/mcp-compliance

# Build for Windows
build-windows:
	mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/mcp-compliance

# Build for macOS
build-macos:
	mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/mcp-compliance

# Deploy locally on macOS
deploy-local: build-macos
	mkdir -p $(LOCAL_DEPLOY_DIR)
	cp $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(LOCAL_DEPLOY_DIR)/$(BINARY_NAME)
	chmod +x $(LOCAL_DEPLOY_DIR)/$(BINARY_NAME)
	@echo "Deployed to $(LOCAL_DEPLOY_DIR)/$(BINARY_NAME)" 