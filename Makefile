.PHONY: build test clean install lint fmt vet release-dry release dev-deps

# Variables
BINARY_NAME=claude-mux
MAIN_PATH=./cmd/claude-mux
VERSION=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X main.version=${VERSION} -X main.commit=$(shell git rev-parse HEAD) -X main.date=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)"

# Tool versions
GOIMPORTS_VERSION := v0.36.0
GOLANGCI_VERSION := v2.4.0
GOTESTSUM_VERSION := v1.12.3
GOSEC_VERSION := v2.22.8

# Install development dependencies
dev-deps:
	@echo "Installing development dependencies..."
	go install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_VERSION)
	go install gotest.tools/gotestsum@$(GOTESTSUM_VERSION)
	go install github.com/securego/gosec/v2/cmd/gosec@$(GOSEC_VERSION)
	@echo "Development dependencies installed!"

# Build the binary
build:
	go build ${LDFLAGS} -o ${BINARY_NAME} ${MAIN_PATH}

# Run tests with gotestsum if available, fallback to go test
test:
	@if command -v gotestsum > /dev/null; then \
		gotestsum --format testname -- -v -race -coverprofile=coverage.out ./...; \
	else \
		go test -v -race -coverprofile=coverage.out ./...; \
	fi

# Run security scan
security:
	@if command -v gosec > /dev/null; then \
		gosec -fmt sarif -out gosec-results.sarif ./...; \
		gosec ./...; \
	else \
		echo "gosec not installed. Run 'make dev-deps' first"; \
		exit 1; \
	fi

# Clean build artifacts
clean:
	rm -f ${BINARY_NAME}
	rm -f coverage.out coverage.html
	rm -rf dist/

# Install binary to GOPATH/bin
install: build
	go install ${LDFLAGS} ${MAIN_PATH}

# Run linter
lint:
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Run 'make dev-deps' first"; \
		exit 1; \
	fi

# Format code with goimports (or gofmt as fallback)
fmt:
	@if command -v goimports > /dev/null; then \
		goimports -w .; \
		go fmt ./...; \
	else \
		go fmt ./...; \
	fi

# Run go vet
vet:
	go vet ./...

# Run all checks
check: fmt vet lint security test

# Test release with goreleaser
release-dry:
	goreleaser release --snapshot --clean

# Create a new release (requires tag)
release:
	goreleaser release --clean

# Run the binary
run: build
	./${BINARY_NAME}

# Development mode with hot reload (requires air)
dev:
	air -c .air.toml

# Generate mocks (if needed in future)
mocks:
	go generate ./...

# Update dependencies
deps:
	go mod download
	go mod tidy

# Show help
help:
	@echo "Available targets:"
	@echo "  dev-deps    - Install development dependencies"
	@echo "  build       - Build the binary"
	@echo "  test        - Run tests with coverage"
	@echo "  clean       - Remove build artifacts"
	@echo "  install     - Install binary to GOPATH/bin"
	@echo "  lint        - Run golangci-lint"
	@echo "  fmt         - Format code"
	@echo "  vet         - Run go vet"
	@echo "  check       - Run all checks (fmt, vet, lint, test)"
	@echo "  release-dry - Test release with goreleaser"
	@echo "  release     - Create a new release"
	@echo "  deps        - Update dependencies"
	@echo "  help        - Show this help message"
