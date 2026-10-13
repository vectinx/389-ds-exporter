VERSION := $(shell cat ./VERSION)
GIT_COMMIT := $(shell git rev-parse --short=7 HEAD 2>/dev/null || echo unknown)
GIT_DIRTY := $(shell git diff --quiet --ignore-submodules -- || echo -dirty)
COMMIT := $(GIT_COMMIT)$(GIT_DIRTY)
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
BUILD_USER := $(shell whoami)
BUILD_DATE := $(shell date -u +%Y-%m-%d)

DOCKER_REGISTRY ?= docker.io/vectinx/389-ds-exporter
DOCKER_TAG ?= $(VERSION)
BUILD_DIR := build
BINARY := $(BUILD_DIR)/389-ds-exporter

GO_SOURCES := $(shell find . -type f -name '*.go' \
	-not -path './vendor/*' \
	-not -path './.git/*')

.PHONY: all docker clean

all: $(BINARY)

$(BINARY): $(GO_SOURCES) go.mod go.sum VERSION
	mkdir -p $(BUILD_DIR)
	GOOS=linux \
	GOARCH=amd64 \
	CGO_ENABLED=0 \
	go build -trimpath -ldflags \
		"-s -w \
		-X github.com/prometheus/common/version.Version=$(VERSION) \
		-X github.com/prometheus/common/version.Revision=$(COMMIT) \
		-X github.com/prometheus/common/version.Branch=$(BRANCH) \
		-X github.com/prometheus/common/version.BuildUser=$(BUILD_USER) \
		-X github.com/prometheus/common/version.BuildDate=$(BUILD_DATE)" \
		-o $(BINARY) \
		./cmd/exporter

build: $(BINARY)

docker: build
	cp -rf docker/* $(BUILD_DIR)
	cd $(BUILD_DIR) && docker build . -t $(DOCKER_REGISTRY):$(DOCKER_TAG)

clean:
	rm -rf $(BUILD_DIR)
