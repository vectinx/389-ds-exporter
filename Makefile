VERSION := $(shell cat ./VERSION)
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

BUILD_DIR := build

build:
	go build -ldflags "-X 'main.Version=$(VERSION)' -X 'main.CommitHash=$(COMMIT)' -X 'main.BuildTime=$(BUILD_TIME)'" -o $(BUILD_DIR)/389-ds-exporter

clean:
	rm -rf build/
