# Sieve — Makefile

BINARY_NAME=sieve
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

.PHONY: all build run test lint vet fmt clean ci install bench vuln

## Build
all: ci build

build:
	@echo "==> Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) .

## Run
run:
	go run main.go testdata/sample.log

run-follow:
	go run main.go -f /tmp/test.log

## Quality
test:
	@echo "==> Running tests..."
	go test -race ./...

test-verbose:
	go test -v -race ./...

test-cover:
	@echo "==> Running tests with coverage..."
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

bench:
	@echo "==> Running benchmarks..."
	go test -bench=. -benchmem ./internal/parser/ ./internal/filter/

lint:
	@echo "==> Running linter..."
	golangci-lint run

vet:
	@echo "==> Running go vet..."
	go vet ./...

fmt:
	@echo "==> Formatting code..."
	gofmt -w .
	goimports -w .

fmt-check:
	@test -z "$$(gofmt -d .)" || (echo "Code not formatted. Run 'make fmt'" && exit 1)

## Security
vuln:
	@echo "==> Checking vulnerabilities..."
	govulncheck ./...

## Full CI pipeline
ci: fmt-check vet lint test
	@echo "==> All checks passed!"

## Install
install:
	go install $(LDFLAGS) .

## Cross-compile
build-all:
	@echo "==> Cross-compiling..."
	GOOS=linux   GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux   GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 .
	GOOS=darwin  GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin  GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe .
	@echo "==> Binaries in dist/"

## Clean
clean:
	rm -f $(BINARY_NAME)
	rm -rf dist/
	rm -f coverage.out coverage.html

## Help
help:
	@echo "Available targets:"
	@echo "  build        - Build the binary"
	@echo "  run          - Run with sample data"
	@echo "  test         - Run tests with race detector"
	@echo "  test-cover   - Run tests with coverage report"
	@echo "  bench        - Run benchmarks"
	@echo "  lint         - Run golangci-lint"
	@echo "  vet          - Run go vet"
	@echo "  fmt          - Format code"
	@echo "  ci           - Run full CI pipeline"
	@echo "  install      - Install to GOPATH/bin"
	@echo "  build-all    - Cross-compile for all platforms"
	@echo "  vuln         - Check for vulnerabilities"
	@echo "  clean        - Remove build artifacts"
