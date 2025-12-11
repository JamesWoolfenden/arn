.PHONY: help build test lint fmt clean install-hooks

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the project
	go build -v ./...

test: ## Run tests
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

lint: ## Run golangci-lint
	golangci-lint run

fmt: ## Format code
	go fmt ./...
	goimports -w .

vet: ## Run go vet
	go vet ./...

tidy: ## Tidy go modules
	go mod tidy

clean: ## Clean build artifacts
	go clean
	rm -f coverage.txt

install-hooks: ## Install pre-commit hooks
	@command -v pre-commit >/dev/null 2>&1 || { echo "pre-commit not found. Install with: pip install pre-commit"; exit 1; }
	pre-commit install

bump:
	git push
	$(eval VERSION=$(shell git describe --tags --abbrev=0 | awk -F. '{OFS="."; $$NF+=1; print $0}'))
	git tag -a $(VERSION) -m "new release"
	git push origin $(VERSION)

psbump:
	git push
	powershell -command "./bump.ps1"

update:
	go get -u
	go mod tidy
	pre-commit autoupdate

all: fmt vet lint test build ## Run all checks and build
