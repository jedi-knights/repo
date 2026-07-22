BIN     := bin/repo
PKG     := ./...
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
LDFLAGS := -X github.com/jedi-knights/repo/internal/cli.Version=$(VERSION) \
           -X github.com/jedi-knights/repo/internal/cli.Commit=$(COMMIT)

.PHONY: build test lint tidy vet run clean

build:
	@mkdir -p bin
	go build -ldflags "$(LDFLAGS)" -o $(BIN) ./cmd/repo

test:
	go test $(PKG)

vet:
	go vet $(PKG)

lint: vet
	go build -o /dev/null ./...

tidy:
	go mod tidy

run: build
	$(BIN) --help

clean:
	rm -rf bin
