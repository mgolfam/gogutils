MODULE := github.com/mgolfam/gogutils
GO     := go

.PHONY: all build test tidy vet fmt lint clean

all: test

## Build the module (no binary output, just verify it compiles)
build:
	$(GO) build ./...

## Run unit tests
test:
	$(GO) test ./...

## Format code
fmt:
	$(GO) fmt ./...

## Run go vet
vet:
	$(GO) vet ./...

## Basic lint alias (fmt + vet + test)
lint: fmt vet test

## Sync go.mod/go.sum
tidy:
	$(GO) mod tidy

## Clean build artifacts (if any)
clean:
	@echo "Nothing to clean (module-based project)"


