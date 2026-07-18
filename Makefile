GO      ?= go
PKGS    := ./...
COVER   := coverage.out

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## Run all unit tests
	$(GO) test $(PKGS)

.PHONY: race
race: ## Run tests with the race detector
	$(GO) test -race -count=1 $(PKGS)

.PHONY: cover
cover: ## Run tests and report coverage
	$(GO) test -covermode=atomic -coverprofile=$(COVER) $(PKGS)
	$(GO) tool cover -func=$(COVER) | tail -1

.PHONY: cover-html
cover-html: cover ## Open the HTML coverage report
	$(GO) tool cover -html=$(COVER)

.PHONY: vet
vet: ## Run go vet
	$(GO) vet $(PKGS)

.PHONY: fmt
fmt: ## Format all Go source
	$(GO) fmt $(PKGS)

.PHONY: tidy
tidy: ## Sync go.mod / go.sum
	$(GO) mod tidy

.PHONY: build
build: ## Compile all packages
	$(GO) build $(PKGS)

.PHONY: check
check: fmt vet test ## Format, vet, and test

.PHONY: clean
clean: ## Remove build and coverage artifacts
	rm -f $(COVER)
	$(GO) clean
