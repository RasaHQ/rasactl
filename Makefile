SHELL := /bin/bash

PREFIX = rasactl

OS=$(or ${GOOS},${GOOS},linux)
ARCH=$(or ${GOARCH},${GOARCH},amd64)

CURRENTDIR = $(shell pwd)
SOURCEDIR = $(CURRENTDIR)

PATH := $(CURRENTDIR)/bin:$(PATH)

VERSION?=$(shell git describe --tags --always)

LD_FLAGS = -ldflags "-X 'github.com/RasaHQ/rasactl/pkg/version.VERSION=$(VERSION)' -s -w"

.PHONY: clean build

build: dist/rasactl

dist/rasactl:
	mkdir -p $(@D)
	GOOS=${OS} GOARCH=${GOARCH} go build $(LD_FLAGS) -o $(@D)

clean:
	rm -rf dist
