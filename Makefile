# Ultimate SDD Framework Makefile

.PHONY: all build clean test install release help

# Default target
all: build

# Build for current platform
build:
	@echo "🔨 Building Ultimate SDD Framework..."
	@go build -o sdd ./cmd/sdd
	@echo "✅ Built: sdd ($(shell du -h sdd | cut -f1))"

# Build optimized release binary
build-release:
	@echo "🔨 Building optimized release binary..."
	@go build \
		-ldflags "-s -w" \
		-o sdd \
		./cmd/sdd
	@echo "✅ Built optimized release: sdd ($(shell du -h sdd | cut -f1))"

# Build for all platforms
build-all:
	@echo "🔨 Building for all platforms..."
	@./scripts/build.sh all

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf build/ sdd sdd.exe *.sha256
	@go clean ./...

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report: coverage.html"

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy

# Install the binary
install: build
	@echo "📦 Installing to /usr/local/bin..."
	@sudo cp sdd /usr/local/bin/
	@echo "✅ Installed: /usr/local/bin/sdd"

# Create release archives
release: build-all
	@echo "📦 Creating release archives..."
	@./scripts/build.sh archive

# Development setup
dev-setup: deps
	@echo "🚀 Setting up development environment..."
	@go install github.com/cosmtrek/air@latest
	@echo "✅ Development environment ready"

# Run with hot reload (requires air)
dev:
	@echo "🔄 Starting development server with hot reload..."
	@air

# Format code
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "🔍 Linting code..."
	@golangci-lint run

# Security scan
security:
	@echo "🔒 Running security scan..."
	@gosec ./...

# Check for vulnerabilities
vuln:
	@echo "🔍 Checking for vulnerabilities..."
	@govulncheck ./...

# Show help
help:
	@echo "🚀 Ultimate SDD Framework - Build System"
	@echo ""
	@echo "Usage: make [TARGET]"
	@echo ""
	@echo "Targets:"
	@echo "  build         Build for current platform"
	@echo "  build-release Build optimized release binary"
	@echo "  build-all     Build for all supported platforms"
	@echo "  clean         Clean build artifacts"
	@echo "  test          Run tests"
	@echo "  test-coverage Run tests with coverage report"
	@echo "  deps          Install dependencies"
	@echo "  install       Install binary to /usr/local/bin"
	@echo "  release       Create release archives"
	@echo "  dev-setup     Setup development environment"
	@echo "  dev           Run with hot reload"
	@echo "  fmt           Format code"
	@echo "  lint          Lint code (requires golangci-lint)"
	@echo "  security      Run security scan (requires gosec)"
	@echo "  vuln          Check for vulnerabilities (requires govulncheck)"
	@echo "  help          Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make build        # Quick build for development"
	@echo "  make build-all    # Cross-platform builds"
	@echo "  make test         # Run test suite"
	@echo "  make install      # Install system-wide"
	@echo "  make release      # Create distributable packages"

# Show version
version:
	@echo "Ultimate SDD Framework $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)"

# Development shortcuts
run: build
	@echo "🚀 Running Ultimate SDD Framework..."
	@./sdd --help

check: fmt lint test security
	@echo "✅ All checks passed!"