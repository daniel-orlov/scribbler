.DEFAULT_GOAL := help

.PHONY: run
run: ## Run the application
	@echo "> Running the scribbler..."
	go run cmd/main.go

.PHONY: cluster-up
cluster-up: ## Run the cluster
	@echo "> Running the cluster..."
	docker-compose -f deploy/docker-compose.yml up -d

.PHONY: cluster-down
cluster-down: ## Stop the cluster
	@echo "> Stopping the cluster..."
	docker-compose -f deploy/docker-compose.yml down

.PHONY: test
test: ## Run tests
	@echo "> Testing..."
	go test -v ./...

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

	@echo "make test - Run tests"
	@echo "make tidy - Clean and format Go code"
	@echo "make fmt - Format Go code"
	@echo "make lint-host - Run golangci-lint directly on host"
	@echo "make help - Show this help"