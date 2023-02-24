.DEFAULT_GOAL := help

.PHONY: run
run: ## Run the application
	@echo "> Running the scribbler..."
	go run cmd/main.go

.PHONY: tidy
tidy: ## Clean and format Go code
	@echo "> Tidying..."
	go mod tidy
	go fmt ./...
	@echo "> Done!"

.PHONY: fmt
fmt: ## Format Go code
	go fmt ./...

.PHONY: lint-host
lint-host: ## Run golangci-lint directly on host
	@echo "> Linting..."
	golangci-lint run -c .golangci.yml -v
	@echo "> Done!"

.PHONY: help
help: ## Show this help
	@echo "make run - Run the application"
	@echo "make tidy - Clean and format Go code"
	@echo "make fmt - Format Go code"
	@echo "make lint-host - Run golangci-lint directly on host"
	@echo "make help - Show this help"