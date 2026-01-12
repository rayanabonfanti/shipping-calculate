.PHONY: tidy build run test test-coverage test-coverage-check test-race fmt vet lint validate pre-commit-check security-check check-signed-commits verify-commits all-checks coverage help

# Variables
BINARY_NAME=shipping-calculator
MAIN_PATH=./cmd/api
COVERAGE_FILE=coverage/coverage.out
COVERAGE_THRESHOLD=80

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@echo "  make tidy                  - Run go mod tidy"
	@echo "  make build                 - Build the application"
	@echo "  make run                   - Run the application"
	@echo "  make test                  - Run all tests"
	@echo "  make test-coverage         - Run tests with coverage report"
	@echo "  make test-coverage-check   - Run tests and validate 80% minimum coverage"
	@echo "  make test-race             - Run tests with race detector"
	@echo "  make fmt                   - Format code with gofmt"
	@echo "  make vet                   - Run go vet"
	@echo "  make lint                  - Run golangci-lint (if installed)"
	@echo "  make validate              - Run fmt, vet, and tests"
	@echo "  make pre-commit-check      - Run pre-commit hooks manually"
	@echo "  make security-check        - Run security hooks only"
	@echo "  make check-signed-commits  - Check if commits are signed (last 10 commits)"
	@echo "  make verify-commits        - Verify commit signatures in detail"
	@echo "  make all-checks            - Run all validations (recommended before commit)"

tidy: ## Run go mod tidy
	@echo "Running go mod tidy..."
	go mod tidy
	@echo "Done!"

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete! Binary: bin/$(BINARY_NAME)"

run: ## Run the application
	@echo "Running $(BINARY_NAME)..."
	go run $(MAIN_PATH)/main.go

test: ## Run all tests
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@mkdir -p coverage
	@go test -coverprofile=$(COVERAGE_FILE) ./internal/... ./telemetry/...
	@echo ""
	@echo "Coverage report:"
	@go tool cover -func=$(COVERAGE_FILE)
	@echo ""
	@COVERAGE=$$(go tool cover -func=$(COVERAGE_FILE) | grep total | awk '{print $$3}'); \
	echo "Total coverage: $$COVERAGE"

test-coverage-check: ## Run tests and validate 80% minimum coverage
	@echo "Running tests with coverage validation..."
	@mkdir -p coverage
	@go test -coverprofile=$(COVERAGE_FILE) ./internal/... ./telemetry/...
	@echo ""
	@echo "Coverage report:"
	@go tool cover -func=$(COVERAGE_FILE)
	@echo ""
	@COVERAGE=$$(go tool cover -func=$(COVERAGE_FILE) | grep total | awk '{print $$3}' | sed 's/%//'); \
	if [ -z "$$COVERAGE" ]; then \
		echo "Error: Could not calculate coverage"; \
		exit 1; \
	fi; \
	echo "Total coverage: $$COVERAGE%"; \
	COVERAGE_INT=$$(echo $$COVERAGE | cut -d. -f1); \
	if [ $$COVERAGE_INT -lt $(COVERAGE_THRESHOLD) ]; then \
		echo "ERROR: Coverage $$COVERAGE% is below threshold of $(COVERAGE_THRESHOLD)%"; \
		exit 1; \
	else \
		echo "SUCCESS: Coverage $$COVERAGE% meets threshold of $(COVERAGE_THRESHOLD)%"; \
	fi

test-race: ## Run tests with race detector
	@echo "Running tests with race detector..."
	go test -race -v ./...

fmt: ## Format code with gofmt
	@echo "Formatting code with gofmt..."
	@if [ $$(gofmt -l . | wc -l) -ne 0 ]; then \
		echo "The following files need formatting:"; \
		gofmt -l .; \
		echo "Run 'make fmt-fix' to fix them automatically"; \
		exit 1; \
	else \
		echo "All files are properly formatted"; \
	fi

fmt-fix: ## Fix formatting issues automatically
	@echo "Fixing code formatting..."
	gofmt -w .
	@echo "Formatting complete!"

vet: ## Run go vet
	@echo "Running go vet..."
	@if go vet ./...; then \
		echo "go vet: no issues found"; \
	else \
		echo "go vet: issues found"; \
		exit 1; \
	fi

lint: ## Run golangci-lint (if installed)
	@echo "Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint is not installed. Install it with:"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi

validate: ## Run fmt, vet, and tests
	@echo "Running validation checks..."
	@$(MAKE) fmt
	@$(MAKE) vet
	@$(MAKE) test
	@echo "All validations passed!"

pre-commit-check: ## Run pre-commit hooks manually
	@echo "Running pre-commit checks..."
	@$(MAKE) fmt
	@$(MAKE) vet
	@$(MAKE) test-coverage-check
	@echo "Pre-commit checks completed!"

security-check: ## Run security hooks only
	@echo "Running security checks..."
	@echo "Checking for known vulnerabilities in dependencies..."
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./...; \
	else \
		echo "govulncheck is not installed. Install it with:"; \
		echo "  go install golang.org/x/vuln/cmd/govulncheck@latest"; \
		echo "Skipping vulnerability check..."; \
	fi
	@echo "Security checks completed!"

check-signed-commits: ## Check if commits are signed (last 10 commits)
	@echo "Checking commit signatures (last 10 commits)..."
	@UNSIGNED=$$(git log --pretty=format:"%H %G?" -10 | grep -v " G$$" | wc -l | tr -d ' '); \
	if [ $$UNSIGNED -gt 0 ]; then \
		echo "WARNING: Found unsigned commits in the last 10 commits"; \
		git log --pretty=format:"%h %s %G?" -10 | grep -v " G$$"; \
		exit 1; \
	else \
		echo "SUCCESS: All commits in the last 10 are signed"; \
	fi

verify-commits: ## Verify commit signatures in detail
	@echo "Verifying commit signatures in detail..."
	@echo "Last 10 commits:"
	@git log --pretty=format:"%h - %s [%G?]" -10
	@echo ""
	@UNSIGNED=$$(git log --pretty=format:"%H %G?" -10 | grep -v " G$$" | wc -l | tr -d ' '); \
	if [ $$UNSIGNED -gt 0 ]; then \
		echo "Found $$UNSIGNED unsigned commit(s)"; \
		exit 1; \
	else \
		echo "All commits are properly signed"; \
	fi

all-checks: ## Run all validations (recommended before commit)
	@echo "Running all validation checks..."
	@echo "=================================="
	@$(MAKE) tidy
	@echo ""
	@echo "=================================="
	@$(MAKE) fmt
	@echo ""
	@echo "=================================="
	@$(MAKE) vet
	@echo ""
	@echo "=================================="
	@$(MAKE) test-race
	@echo ""
	@echo "=================================="
	@$(MAKE) test-coverage-check
	@echo ""
	@echo "=================================="
	@if command -v golangci-lint >/dev/null 2>&1; then \
		$(MAKE) lint; \
		echo ""; \
		echo "=================================="; \
	fi
	@$(MAKE) security-check
	@echo ""
	@echo "=================================="
	@echo "All checks completed successfully!"

coverage: test-coverage-check ## Alias for test-coverage-check
