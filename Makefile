# Makefile for idl-parser project

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

BINARY_NAME=idl-parser

.PHONY: all
all: test

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

.PHONY: test
test:
	$(GOTEST) ./...

.PHONY: test-verbose
test-verbose:
	$(GOTEST) -v ./...

.PHONY: test-coverage
test-coverage:
	$(GOTEST) -cover ./...

.PHONY: test-coverage-html
test-coverage-html:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

.PHONY: build
build:
	$(GOBUILD) -o $(BINARY_NAME) .

.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) tidy

.PHONY: fmt
fmt:
	$(GOCMD) fmt ./...

.PHONY: vet
vet:
	$(GOCMD) vet ./...

.PHONY: check
check: fmt vet test
