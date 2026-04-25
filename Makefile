.PHONY: wire generate run build test clean dev install-tools http-wire migration-wire migrate benchmark-scheduling predictive-train ml-install ml-export-dataset ml-train-db ml-train-csv ml-predict-db ml-predict-csv ml-demo

PYTHON ?= python3
ML_DIR := /Users/hant/golang/doan/ml/at_risk_prediction
ML_DEMO_DATASET := $(ML_DIR)/data/raw/at_risk_dataset_demo.csv
ML_MPLCONFIGDIR := /tmp/matplotlib
ML_XDG_CACHE_HOME := /tmp

# Install required tools
install-tools:
	@echo "Installing Wire..."
	go install github.com/google/wire/cmd/wire@latest
	@echo "Installing Swag..."
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Installing GORM..."
	go get -u gorm.io/gorm
	go get -u gorm.io/driver/postgres
	go get -u github.com/google/uuid
	go mod tidy
	@echo "Tools installed successfully!"

# Generate wire for HTTP server
http-wire:
	@echo "Generating wire for HTTP server..."
	cd cmd/http && wire

# Generate wire for migration CLI
migration-wire:
	@echo "Generating wire for migration CLI..."
	cd cmd/cli/migration && wire

# Generate all wire injections
wire: http-wire migration-wire
	@echo "All wire code generated successfully!"

# Generate swagger documentation
swagger:
	@echo "Generating Swagger docs..."
	swag init -g cmd/http/main.go -o cmd/http/docs
	@echo "Swagger docs generated!"

# Generate all (wire + swagger)
generate: wire swagger
	@echo "Running go generate..."
	go generate ./...
	@echo "All code generated successfully!"

# Run HTTP server (with generation)
run: generate
	@echo "Running HTTP server..."
	go run cmd/http/main.go cmd/http/wire_gen.go

# Run migration (up)
migrate: migration-wire
	@echo "Running migration..."
	go run cmd/cli/migration/main.go cmd/cli/migration/wire_gen.go

# Seed database with sample data
seed:
	@echo "Seeding database..."
	go run cmd/cli/seed/main.go
	@echo "Database seeded successfully!"

# Run synthetic scheduling benchmark study
benchmark-scheduling:
	@echo "Running scheduling benchmark study..."
	go run ./cmd/cli/scheduling_benchmark

predictive-train:
	@echo "Running minimal AT_RISK training..."
	go run ./cmd/cli/predictive_train

# Install Python dependencies for AT_RISK ML project
ml-install:
	@echo "Installing Python dependencies for AT_RISK ML project..."
	$(PYTHON) -m pip install -r $(ML_DIR)/requirements.txt

# Export normalized dataset from DB for ML pipeline
ml-export-dataset:
	@echo "Exporting AT_RISK dataset from PostgreSQL..."
	cd $(ML_DIR) && MPLCONFIGDIR=$(ML_MPLCONFIGDIR) XDG_CACHE_HOME=$(ML_XDG_CACHE_HOME) $(PYTHON) scripts/export_dataset.py --dataset-name at_risk_dataset

# Train/evaluate ML models directly from DB
ml-train-db:
	@echo "Training AT_RISK ML models from PostgreSQL..."
	cd $(ML_DIR) && MPLCONFIGDIR=$(ML_MPLCONFIGDIR) XDG_CACHE_HOME=$(ML_XDG_CACHE_HOME) $(PYTHON) scripts/train_from_db.py --dataset-name at_risk_dataset_db

# Train/evaluate ML models from demo CSV
ml-train-csv:
	@echo "Training AT_RISK ML models from demo CSV..."
	cd $(ML_DIR) && MPLCONFIGDIR=$(ML_MPLCONFIGDIR) XDG_CACHE_HOME=$(ML_XDG_CACHE_HOME) $(PYTHON) scripts/train_from_csv.py --input $(ML_DEMO_DATASET) --dataset-name at_risk_dataset_demo

# Generate latest prediction artifact from DB
ml-predict-db:
	@echo "Generating AT_RISK predictions from PostgreSQL..."
	cd $(ML_DIR) && MPLCONFIGDIR=$(ML_MPLCONFIGDIR) XDG_CACHE_HOME=$(ML_XDG_CACHE_HOME) $(PYTHON) scripts/predict_from_db.py

# Generate latest prediction artifact from demo CSV
ml-predict-csv:
	@echo "Generating AT_RISK predictions from demo CSV..."
	cd $(ML_DIR) && MPLCONFIGDIR=$(ML_MPLCONFIGDIR) XDG_CACHE_HOME=$(ML_XDG_CACHE_HOME) $(PYTHON) scripts/predict_from_csv.py --input $(ML_DEMO_DATASET)

# Full offline demo flow for ML artifacts
ml-demo: ml-train-csv ml-predict-csv
	@echo "AT_RISK demo artifacts generated successfully."

# Run migration down
migrate-down: migration-wire
	@echo "Rolling back migration..."
	go run cmd/cli/migration/main.go cmd/cli/migration/wire_gen.go --down

# Build HTTP server
build: generate
	@echo "Building HTTP server..."
	mkdir -p bin
	go build -o bin/http-server cmd/http/*.go
	@echo "Binary created at: bin/http-server"

# Build migration CLI
build-migration: migration-wire
	@echo "Building migration CLI..."
	mkdir -p bin
	go build -o bin/migration cmd/cli/migration/*.go
	@echo "Binary created at: bin/migration"

# Build all binaries
build-all: build build-migration
	@echo "All binaries built successfully!"

# Run tests
test:
	@echo "Running tests..."
	go test -v -race ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean generated files and binaries
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f cmd/http/wire_gen.go
	rm -f cmd/cli/migration/wire_gen.go
	rm -f coverage.out coverage.html
	@echo "Cleaned!"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies downloaded!"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatted!"

# Development mode (auto-generate and run)
dev: generate
	@echo "Starting development server..."
	go run cmd/http/main.go cmd/http/wire_gen.go

# Quick start (one command to rule them all)
start: deps generate dev

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker build -t doan-app -f deploy/Dockerfile .

docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 doan-app

# Docker compose for local development
docker-local-up:
	@echo "Starting local Docker services..."
	docker-compose -f tools/docker-compose.local.yml up -d

docker-local-down:
	@echo "Stopping local Docker services..."
	docker-compose -f tools/docker-compose.local.yml down

# Help command
help:
	@echo "Available commands:"
	@echo "  make install-tools    - Install required tools (wire, swag)"
	@echo "  make generate        - Generate all code (wire + swagger)"
	@echo "  make wire           - Generate all wire injection code"
	@echo "  make swagger        - Generate swagger documentation"
	@echo "  make run            - Generate and run HTTP server"
	@echo "  make build          - Build HTTP server binary"
	@echo "  make build-all      - Build all binaries"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make clean          - Clean generated files and binaries"
	@echo "  make deps           - Download dependencies"
	@echo "  make fmt            - Format code"
	@echo "  make dev            - Development mode (auto-generate and run)"
	@echo "  make start          - Quick start (deps + generate + dev)"
	@echo ""
	@echo "ML commands:"
	@echo "  make ml-install      - Install Python dependencies for ML project"
	@echo "  make ml-export-dataset - Export normalized AT_RISK dataset from PostgreSQL"
	@echo "  make ml-train-db     - Train/evaluate ML models from PostgreSQL"
	@echo "  make ml-train-csv    - Train/evaluate ML models from demo CSV"
	@echo "  make ml-predict-db   - Generate prediction artifact from PostgreSQL"
	@echo "  make ml-predict-csv  - Generate prediction artifact from demo CSV"
	@echo "  make ml-demo         - Run offline demo flow (train + predict from CSV)"
	@echo ""
	@echo "Migration commands:"
	@echo "  make migrate        - Run database migrations"
	@echo "  make migrate-down   - Rollback migrations"
	@echo ""
	@echo "Docker commands:"
	@echo "  make docker-build      - Build Docker image"
	@echo "  make docker-local-up   - Start local Docker services"
	@echo "  make docker-local-down - Stop local Docker services"
	@echo ""
	@echo "Quick start: make install-tools && make deps && make dev"
