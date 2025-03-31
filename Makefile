.PHONY: all run build clean test deps docker-up docker-down

# Application name
APP_NAME = go-starter

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GORUN = $(GOCMD) run
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod

# Main entry point
MAIN_PATH = ./cmd/main.go

# Build directory
BUILD_DIR = ./build

# Environment variables
export SERVER_PORT ?= 8080
export DATABASE_URL ?= postgres://todos:todos@db:5432/todos?sslmode=disable
export KAFKA_BROKERS ?= kafka:9092
export KAFKA_TOPIC ?= todos

# Default target
all: deps test build

# Run the application
run:
	@echo "Running application..."
	$(GORUN) $(MAIN_PATH)

# Build the application
build:
	@echo "Building application..."
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOGET) github.com/gorilla/mux
	$(GOGET) github.com/joho/godotenv
	$(GOGET) github.com/lib/pq
	$(GOGET) github.com/segmentio/kafka-go
	$(GOMOD) tidy

# Start Docker services
docker-up:
	@echo "Starting Docker services..."
	docker-compose up -d

# Stop Docker services
docker-down:
	@echo "Stopping Docker services..."
	docker-compose down

# Start the development environment
dev: docker-up run

# Stop the development environment
stop: docker-down

# Show help
help:
	@echo "Available targets:"
	@echo "  make          : Download dependencies, run tests, and build the application"
	@echo "  make run      : Run the application"
	@echo "  make build    : Build the application"
	@echo "  make clean    : Clean build artifacts"
	@echo "  make test     : Run tests"
	@echo "  make deps     : Download dependencies"
	@echo "  make docker-up: Start Docker services"
	@echo "  make docker-down: Stop Docker services"
	@echo "  make dev      : Start the development environment"
	@echo "  make stop     : Stop the development environment"
	@echo "  make help     : Show this help message" 