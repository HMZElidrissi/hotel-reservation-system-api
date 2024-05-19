# Variables
BIN := server

# Build the application
build:
	@echo "Building the application..."
	@go build -o $(BIN) cmd/server

# Run the application
run: build
	@echo "Running the application..."
	@./$(BIN)