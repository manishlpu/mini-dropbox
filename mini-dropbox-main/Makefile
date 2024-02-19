SHELL := /bin/bash

GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test

BINARY_NAME=mini-dropbox
EXEC_ROOT=./cmd/dropbox

.PHONY: build
build:
	$(GOBUILD) -tags "$(TAGS)" -o $(BINARY_NAME) -v $(EXEC_ROOT)

.PHONY: help
help:
	$(GORUN) $(EXEC_ROOT) help

.PHONY: run
run:
	mkdir -p storage/
	$(GORUN) $(EXEC_ROOT) run

.PHONY: test
test:
	echo -e "\n\n Testing..."
	$(GOTEST) -v -race -cover --shuffle=on ./...
