export SHELL = /bin/bash

# TOOLING

LINTER_VERSION ?= 1.23.6
LINTER_EXE := golangci-lint
LINTER_FOLDER := $(monorepo_root)/bin
LINTER ?= $(LINTER_FOLDER)/$(LINTER_EXE)
LINT_CONFIG_PATH := $(monorepo_root)/.golangci.yml
install:
	@if [ "`$(LINTER) --version | awk '{print $$4}'`" != $(LINTER_VERSION) ]; then \
		curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LINTER_FOLDER) v$(LINTER_VERSION); fi

# TESTING STEPS

lint:
	golangci-lint run

test:
	go test ./...
