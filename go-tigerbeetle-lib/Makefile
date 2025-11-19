.PHONY: help test test-coverage fmt lint vet build clean examples
.PHONY: docker-start docker-stop docker-restart docker-logs docker-clean docker-status
.PHONY: cluster-start cluster-stop cluster-logs cluster-clean cluster-status

# Default target
help:
	@echo "TigerBeetle Go Library - Available Commands"
	@echo "============================================"
	@echo ""
	@echo "Docker Commands (Single Node):"
	@echo "  docker-start    - Start TigerBeetle single node"
	@echo "  docker-stop     - Stop TigerBeetle"
	@echo "  docker-restart  - Restart TigerBeetle"
	@echo "  docker-logs     - View TigerBeetle logs"
	@echo "  docker-clean    - Stop and remove all data"
	@echo "  docker-status   - Check TigerBeetle status"
	@echo ""
	@echo "Docker Commands (Cluster):"
	@echo "  cluster-start   - Start 3-node TigerBeetle cluster"
	@echo "  cluster-stop    - Stop cluster"
	@echo "  cluster-logs    - View cluster logs"
	@echo "  cluster-clean   - Stop cluster and remove all data"
	@echo "  cluster-status  - Check cluster status"
	@echo ""
	@echo "Development Commands:"
	@echo "  test            - Run all tests"
	@echo "  test-coverage   - Run tests with coverage report"
	@echo "  fmt             - Format code"
	@echo "  lint            - Run linter (requires golangci-lint)"
	@echo "  vet             - Run go vet"
	@echo "  build           - Build all packages"
	@echo "  examples        - Build all examples"
	@echo "  clean           - Clean build artifacts"
	@echo "  all             - Run fmt, vet, lint, test"
	@echo ""

# Run all tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format code
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run

# Run go vet
vet:
	go vet ./...

# Build all packages
build:
	go build ./...

# Build examples
examples:
	@echo "Building examples..."
	@for example in examples/*.go; do \
		echo "Building $$example..."; \
		go build -o /dev/null $$example || exit 1; \
	done
	@echo "All examples built successfully"

# Clean build artifacts
clean:
	go clean ./...
	rm -f coverage.txt coverage.html

# Run all checks
all: fmt vet test
	@echo "All checks passed!"

# Install dependencies
deps:
	go mod download
	go mod tidy

# Update dependencies
update-deps:
	go get -u ./...
	go mod tidy

# ============================================
# Docker Commands
# ============================================

# Single node commands
docker-start:
	@echo "Starting TigerBeetle single node..."
	docker-compose up -d
	@echo "Waiting for TigerBeetle to be ready..."
	@sleep 2
	@docker-compose ps

docker-stop:
	@echo "Stopping TigerBeetle..."
	docker-compose down

docker-restart: docker-stop docker-start

docker-logs:
	docker-compose logs -f tigerbeetle

docker-clean:
	@echo "Stopping TigerBeetle and removing all data..."
	docker-compose down -v
	@echo "All data removed."

docker-status:
	@docker-compose ps

# Cluster commands
cluster-start:
	@echo "Starting TigerBeetle 3-node cluster..."
	docker-compose -f docker-compose.cluster.yml up -d
	@echo "Waiting for cluster to be ready..."
	@sleep 3
	@docker-compose -f docker-compose.cluster.yml ps

cluster-stop:
	@echo "Stopping TigerBeetle cluster..."
	docker-compose -f docker-compose.cluster.yml down

cluster-logs:
	docker-compose -f docker-compose.cluster.yml logs -f

cluster-clean:
	@echo "Stopping cluster and removing all data..."
	docker-compose -f docker-compose.cluster.yml down -v
	@echo "All cluster data removed."

cluster-status:
	@docker-compose -f docker-compose.cluster.yml ps

# Quick development workflow
dev: docker-start test
	@echo "Development environment ready!"
	@echo "TigerBeetle is running on port 3000"
