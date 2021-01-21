BIN := kubetools

GO := GO111MODULE=on GOFLAGS=-mod=vendor go
GOARCH := $(shell go env GOARCH)
GOOS := $(shell go env GOOS)

VERSION := $(shell git describe --always --dirty --long)
BUILD_COMMIT := $(shell git rev-parse HEAD)
BUILD_TIMESTAMP := $(shell git show -s --format=%ct HEAD)

LDFLAGS := "-X github.com/vetyy/kubetools/version.Version=$(VERSION) \
	-X github.com/vetyy/kubetools/version.CommitID=$(BUILD_COMMIT) \
	-X github.com/vetyy/kubetools/version.CommitTimestamp=$(BUILD_TIMESTAMP)"

.PHONY: clean build test test-html lint

clean:
	rm -rf vendor/

build:
	$(GO) build -o $(BIN) -ldflags $(LDFLAGS) ./cmd/kubetools/...
	@hash notify-send 2> /dev/null && notify-send --app-name kubetools 'Build complete.' || echo 'Build complete.'

test:
	$(GO) test -cover ./... -v

test-html:
	$(GO) test -coverprofile=cov.out ./... -v
	$(GO) tool cover -html=cov.out

lint:
	@golangci-lint run ./...

deps:
	$(GO) mod vendor
