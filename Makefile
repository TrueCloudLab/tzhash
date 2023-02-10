#!/usr/bin/make -f
SHELL = bash

REPO ?= $(shell go list -m)
VERSION ?= $(shell git describe --tags --dirty --match "v*" --always --abbrev=8 2>/dev/null || cat VERSION 2>/dev/null || echo "develop")

BIN = bin
DIRS = $(BIN)

# List of binaries to build.
CMDS = $(notdir $(basename $(wildcard cmd/*)))
BINS = $(addprefix $(BIN)/, $(CMDS))

.PHONY: all help clean

# To build a specific binary, use it's name prefix with bin/ as a target
# For example `make bin/tzsum` will build only storage node binary
# Just `make` will build all possible binaries
all: $(DIRS) $(BINS)

# help target
include help.mk

$(BINS): $(DIRS) dep
	@echo "⇒ Build $@"
	CGO_ENABLED=0 \
	go build -v -trimpath \
	-ldflags "-X $(REPO)/misc.Version=$(VERSION)" \
	-o $@ ./cmd/$(notdir $@)

$(DIRS):
	@echo "⇒ Ensure dir: $@"
	@mkdir -p $@

# Pull go dependencies
dep:
	@printf "⇒ Download requirements: "
	CGO_ENABLED=0 \
	go mod download && echo OK
	@printf "⇒ Tidy requirements : "
	CGO_ENABLED=0 \
	go mod tidy -v && echo OK

# Run Unit Test with go test
test:
	@echo "⇒ Running go test"
	@go test ./...

# Run Unit Test with go test
test.generic:
	@echo "⇒ Running go test with generic tag"
	@go test ./... --tags=generic

# Print version
version:
	@echo $(VERSION)

clean:
	rm -rf vendor
	rm -rf .cache
	rm -rf $(BIN)
