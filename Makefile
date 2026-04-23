.PHONY: build install clean run deps test help

# Variables
BINARY_NAME=curlkeeper
INSTALL_PATH=/usr/local/bin
GO=go

# Default target
all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build -o $(BINARY_NAME) -ldflags="-s -w"
	@echo "Build complete!"

# Install to system
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	sudo mv $(BINARY_NAME) $(INSTALL_PATH)/
	@echo "Installation complete! Run '$(BINARY_NAME)' to start."

# Uninstall from system
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Uninstall complete!"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	@echo "Clean complete!"

# Run the application
run:
	$(GO) run .

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "Dependencies ready!"

# Run tests
test:
	$(GO) test -v ./...

# Format code
fmt:
	$(GO) fmt ./...

# Run linter
lint:
	golangci-lint run

# Show help
help:
	@echo "Curl Keeper - Makefile commands:"
	@echo ""
	@echo "  make build      - Build the application"
	@echo "  make install    - Build and install to $(INSTALL_PATH)"
	@echo "  make uninstall  - Remove from system"
	@echo "  make clean      - Remove build artifacts"
	@echo "  make run        - Run without building"
	@echo "  make deps       - Download and tidy dependencies"
	@echo "  make test       - Run tests"
	@echo "  make fmt        - Format code"
	@echo "  make lint       - Run linter"
	@echo "  make help       - Show this help message"
	@echo ""
