# Project variables
BINARY_NAME := compactify
COVERAGE_FILE := coverage.out

.PHONY : all fmt fmt-fix vet test coverage test build clean init-hooks

all: fmt vet test build

fmt:
	@echo "🔍 Checking Go formatting..."
	@files="$$(find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -l {} +)"; \
	if [ -n "$$files" ]; then \
		echo "❌ The following files need formatting:"; \
		echo "$$files"; \
		echo "Run 'make fmt-fix' to format them."; \
		exit 1; \
	fi
	@echo "✅ Everything is formatted!"

fmt-fix:
	@echo "🧹 Formatting Go code..."
	@gofmt -w $$(find . -type f -name '*.go' -not -path './vendor/*')

vet:
	@echo "🔍 Running go vet..."
	go vet ./...

test:
	@echo "🧪 Running tests..."
	go test ./... --failfast

coverage:
	@echo "📊 Running tests with coverage..."
	go test ./... -coverprofile=$(COVERAGE_FILE)
	@go tool cover -html=$(COVERAGE_FILE) -o coverage.html
	@echo "🌐 Opening in browser..."
	@explorer.exe coverage.html 2>/dev/null || open coverage.html 2>/dev/null || xdg-open coverage.html 2>/dev/null || true

build:
	@echo "🚀 Building the binary..."
	go build -ldflags="-w -s -X 'github.com/felipesimis/go-compactify-cli/version.Version=local-dev'" -trimpath -o $(BINARY_NAME) .

clean:
	@echo "🧹 Cleaning up..."
	go clean
	@rm -f $(BINARY_NAME)
	@rm -f $(COVERAGE_FILE)
	@rm -f coverage.html

init-hooks:
	@echo "🪝 Installing Git Hooks..."
	@cd tools && go run github.com/evilmartians/lefthook/v2 install --force --reset-hooks-path
	@echo "✅ Git Hooks installed!"