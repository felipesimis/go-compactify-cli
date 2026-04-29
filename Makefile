# Project variables
BINARY_NAME := compactify
COVERAGE_FILE := coverage.out

.PHONY : all fmt vet test coverage test clean

all: fmt vet test build

fmt:
	@echo "🧹 Running go fmt..."
	@go fmt ./...

vet:
	@echo "🔍 Running go vet..."
	@go vet ./...

test:
	@echo "🧪 Running tests..."
	@go test ./... --failfast

coverage:
	@echo "📊 Running tests with coverage..."
	@go test ./... -coverprofile=$(COVERAGE_FILE)
	@go tool cover -html=$(COVERAGE_FILE) -o coverage.html
	@echo "🌐 Opening in browser..."
	@explorer.exe coverage.html 2>/dev/null || open coverage.html 2>/dev/null || xdg-open coverage.html 2>/dev/null || true

build:
	@echo "🚀 Building the binary..."
	@go build -ldflags="-w -s -X 'github.com/felipesimis/go-compactify-cli/version.Version=local-dev'" -trimpath -o $(BINARY_NAME) .

clean:
	@echo "🧹 Cleaning up..."
	@go clean
	@rm -f $(BINARY_NAME)
	@rm -f $(COVERAGE_FILE)
	@rm -f coverage.html