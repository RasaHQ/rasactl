SHELL := /bin/bash

PREFIX = rasactl

CURRENTDIR = $(shell pwd)
SOURCEDIR = $(CURRENTDIR)
GOIMPORTS = $(GOBIN)/goimports
PATH := $(CURRENTDIR)/bin:$(PATH)

VERSION?=$(shell git describe --tags --always)
ARCH = $(shell uname -p)
LD_FLAGS = -ldflags "-X 'github.com/RasaHQ/rasactl/pkg/version.VERSION=$(VERSION)' -s -w"

# go option
PKG        := ./...
TAGS       :=
TESTS      := .
TESTFLAGS  :=
GOFLAGS    :=

.PHONY: clean build

build: dist/rasactl

dist/rasactl:
	mkdir -p $(@D)
	source <(go env)
	go build $(LD_FLAGS) -o $(@D)

.PHONY: test
test: build
ifeq ($(ARCH),s390x)
test: TESTFLAGS += -v
else
test: TESTFLAGS += -race -v
endif
test: test-unit

.PHONY: test-unit
test-unit:
	@echo
	@echo "== Running unit tests =="
	go test $(GOFLAGS) -run $(TESTS) $(PKG) $(TESTFLAGS)

.PHONY: test-style
test-style:
	golangci-lint run

$(GOIMPORTS):
	(cd /; go get -u golang.org/x/tools/cmd/goimports)

.PHONY: format
format: $(GOIMPORTS)
	go list -f '{{.Dir}}' ./... | xargs $(GOIMPORTS) -w -local github.com/RasaHQ/rasactl

clean:
	rm -rf dist

## Release
.PHONY: changelog
changelog:  ## Generate changelog
	git-chglog --next-tag $(VERSION) -o CHANGELOG.md

.PHONY: release
release: changelog   ## Release a new tag
	git add CHANGELOG.md
	git commit -m "chore: update changelog for $(VERSION)"
