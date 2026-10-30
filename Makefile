GO            ?= go
GOLANGCI_LINT ?= golangci-lint

.PHONY: run test lint fmt

run:
	$(GO) run . 

test:
	$(GO) test ./...

lint:
	$(GO) vet ./...
	$(GOLANGCI_LINT) run ./...

fmt:
	$(GOLANGCI_LINT) fmt ./...
