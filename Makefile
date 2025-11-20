.PHONY: build clean help

# Variables
BINARY_NAME=mcp-server
GOOS=linux
GOARCH=amd64
BUILD_DIR=.

# Default target
all: build

# Build for Linux AMD64
build:
	@echo "Building $(BINARY_NAME) for $(GOOS)/$(GOARCH)..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BUILD_DIR)/$(BINARY_NAME) main.go
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
	@echo "Clean complete"

# Display help
help:
	@echo "Available targets:"
	@echo "  build   - Build the application for Linux AMD64 (default)"
	@echo "  clean   - Remove build artifacts"
	@echo "  help    - Display this help message"
