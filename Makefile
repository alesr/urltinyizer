.DEFAULT_GOAL := help

NAME = $(shell basename $(PWD))
AWAIT_DB_SCRIPT = wait-for-it.sh

.PHONY: help
help:
	@echo "------------------------------------------------------------------------"
	@echo "${NAME}"
	@echo "------------------------------------------------------------------------"
	@grep -E '^[a-zA-Z0-9_/%\-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the project
	@GOOS=linux go build -o $(NAME) main.go

.PHONY: run
run: build ## Run the project on a Docker container
	@[ -f $(AWAIT_DB_SCRIPT) ] || curl -o $(AWAIT_DB_SCRIPT) "https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh"
	@chmod +x $(AWAIT_DB_SCRIPT)
	@docker-compose -f docker-compose.yml up db urltinyizer --force-recreate --build

.PHONY: lint
lint: ## Run linter
	@go vet ./...
	@go fmt ./...

.PHONY: test-it
test-it: ## Run integration tests
	@docker-compose -f docker-compose.yml up -d db
	@sleep 3
	@go test -v -tags=integration -race -vet=all -count=1 -timeout 60s ./...
	@docker-compose -f docker-compose.yml down

.PHONY: test-unit
test-unit: ## Run unit tests
	@go test -v -race -vet=all -count=1 -timeout 60s ./...

.PHONY: test
test: lint test-unit test-it ## Run all tests


