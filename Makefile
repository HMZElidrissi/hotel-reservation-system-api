# Variables
BIN := hotel-reservation-system-api

# Default target
.PHONY: all
all: build

# Load the environment variables
.env:
	@echo "Loading environment variables..."
	@source .env

# Build the application
.PHONY: build
build:
	@echo "Building the application..."
	@go build -o $(BIN) ./cmd/api

# Run the application
.PHONY: run
run: build
	@echo "Running the application..."
	@./$(BIN)

# Test the application
.PHONY: test
test:
	@echo "Running tests..."
	@go test ./...

# Clean up build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -f $(EXECUTABLE)

# Format the code
.PHONY: format
fmt:
	@echo "Formatting the code..."
	@go fmt ./...

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod tidy

# Lint the code
.PHONY: lint
lint:
	@echo "Linting the code..."
	@golangci-lint run

# Generate .env file from .env.example if it does not exist
.PHONY: init
init:
	@echo "Initializing the project..."
	@[ -f .env ] || cp .env.example .env