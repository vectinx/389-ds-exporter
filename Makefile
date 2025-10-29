VERSION := $(shell cat ./VERSION)
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

DOCKER_REGISTRY ?= docker.io/vectinx/389-ds-exporter
DOCKER_TAG ?= $(VERSION)
BUILD_DIR := build

build:
	GOOS=linux \
	GOARCH=amd64 \
	CGO_ENABLED=0 \
	go build -trimpath -ldflags "-X 'main.Version=$(VERSION)' -X 'main.CommitHash=$(COMMIT)' -X 'main.BuildTime=$(BUILD_TIME)'" -o $(BUILD_DIR)/389-ds-exporter

docker: build
	mkdir -p $(BUILD_DIR)
	cp -rf docker/* $(BUILD_DIR)
	cd $(BUILD_DIR) && docker build . -t $(DOCKER_REGISTRY):$(DOCKER_TAG)

clean:
	rm -rf build/
