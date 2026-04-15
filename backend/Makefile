# ══════════════════════════════════════════════════════════════════════════════
#  Makefile — Golang backend
# ══════════════════════════════════════════════════════════════════════════════

APP_NAME   := golang-backend
MAIN_PKG   := ./cmd/server
BIN_DIR    := ./bin
BIN        := $(BIN_DIR)/$(APP_NAME)

.DEFAULT_GOAL := help

## ── Development ──────────────────────────────────────────────────────────────

.PHONY: run
run: ## Run the server with hot reload (requires 'air')
	@air -c .air.toml

.PHONY: dev
dev: ## Run the server directly via go run
	@go run $(MAIN_PKG)/main.go

## ── Build ────────────────────────────────────────────────────────────────────

.PHONY: build
build: ## Compile the binary into ./bin/
	@mkdir -p $(BIN_DIR)
	@go build -ldflags="-s -w" -o $(BIN) $(MAIN_PKG)/main.go
	@echo "✅  Built $(BIN)"

## ── Quality ──────────────────────────────────────────────────────────────────

.PHONY: test
test: ## Run all tests
	@go test -v -race -count=1 ./...

.PHONY: test-cover
test-cover: ## Run tests and open HTML coverage report
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅  Coverage report: coverage.html"

.PHONY: lint
lint: ## Run golangci-lint (requires 'golangci-lint')
	@golangci-lint run ./...

.PHONY: vet
vet: ## Run go vet
	@go vet ./...

.PHONY: fmt
fmt: ## Format all Go files
	@gofmt -w .

## ── Dependencies ─────────────────────────────────────────────────────────────

.PHONY: tidy
tidy: ## Tidy go modules
	@go mod tidy

.PHONY: deps
deps: ## Download all dependencies
	@go mod download

## ── Database ─────────────────────────────────────────────────────────────────

.PHONY: migrate
migrate: ## Apply SQL migrations (runs all *.sql in migrations/)
	@echo "Running migrations..."
	@for f in migrations/*.sql; do \
		echo "  ▶ $$f"; \
		psql "$$DATABASE_URL" -f "$$f"; \
	done

## ── Cleanup ──────────────────────────────────────────────────────────────────

.PHONY: clean
clean: ## Remove build artifacts
	@rm -rf $(BIN_DIR) coverage.out coverage.html

## ── Help ─────────────────────────────────────────────────────────────────────

.PHONY: help
help: ## Show this help message
	@echo ""
	@echo "  $(APP_NAME) — available targets:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
