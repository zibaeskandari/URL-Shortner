APP_NAME := url_shortner
BIN_DIR  := bin
BIN_PATH := $(BIN_DIR)/$(APP_NAME)

MAIN_PKG := ./cmd/main.go

ENT_PKG      := ./internal/infrastructure/persistence/ent/schema
ATLAS_DIR    := ./internal/infrastructure/persistence/ent/migrations
DB_URL       := "postgres://postgres:pgpassword@localhost:5432/url_shortner?search_path=public&sslmode=disable"


.PHONY: help
help: ## Show this help message
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)


.PHONY: all
all: build

.PHONY: build
build: ## Build the application
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_PATH) $(MAIN_PKG)


.PHONY: test
test: ## Run all tests
	@go test ./...


.PHONY: test-unit
test-unit: ## Run unit tests
	@go test -tags=unit ./...


.PHONY: test-integration
test-integration: ## Run integration tests
	@go test -tags=integration ./...


.PHONY: dev
dev:  ## Run with air on dev.yaml config
	@air -- --config ./configs/dev.yaml


.PHONY: ent-generate
ent-generate:  ## Generate entgo files based on schema
	@go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/versioned-migration $(ENT_PKG)


.PHONY: migrate-new
migrate-new:  ## Creates a migration file with a given name
ifndef name
	$(error name is required, e.g. `make migrate-new name=add_users_table`)
endif
	@echo ">> Creating new migration: $(name)"
	@atlas migrate diff $(name) \
		--dir "file://$(ATLAS_DIR)" \
		--to "ent://$(ENT_PKG)" \
		--dev-url "docker://postgres/18.1-alpine/ent"


.PHONY: migrate-up
migrate-up:  ## Apply migrations to dev db
	@atlas migrate apply \
		--dir "file://$(ATLAS_DIR)" \
		--url $(DB_URL)

.PHONY: migrate-status
migrate-status:  ## Show migration status
	@atlas migrate status \
		--dir "file://$(ATLAS_DIR)" \
		--url $(DB_URL)
