VERSION := $(shell cat ./VERSION)
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

ARCH := amd64
BUILD_DIR := build
PKG_DIR := build/389-ds-exporter_$(VERSION)_$(ARCH)

build:
	go build -ldflags "-X 'main.Version=$(VERSION)' -X 'main.CommitHash=$(COMMIT)' -X 'main.BuildTime=$(BUILD_TIME)'" -o $(BUILD_DIR)/389-ds-exporter

deb: build
	@echo "Creating package directories"
	mkdir -p package/deb/$(PKG_DIR)/DEBIAN
	mkdir -p package/deb/$(PKG_DIR)/usr/local/bin
	mkdir -p package/deb/$(PKG_DIR)/etc/389-ds-exporter
	mkdir -p package/deb/$(PKG_DIR)/etc/systemd/system
	mkdir -p package/deb/$(PKG_DIR)/etc/logrotate.d/
	@echo "Generating control file"
	sed "s/@VERSION@/$(VERSION)/" package/deb/control.in > package/deb/$(PKG_DIR)/DEBIAN/control
	@echo "Copy build files"
	cp $(BUILD_DIR)/389-ds-exporter package/deb/$(PKG_DIR)/usr/local/bin
	chmod 755 package/deb/$(PKG_DIR)/usr/local/bin

	cp config.yml package/deb/$(PKG_DIR)/etc/389-ds-exporter/config.yml
	chmod 644 package/deb/$(PKG_DIR)/etc/389-ds-exporter/config.yml

	cp package/deb/389-ds-exporter.service package/deb/$(PKG_DIR)/etc/systemd/system/389-ds-exporter.service
	chmod 644 package/deb/$(PKG_DIR)/etc/systemd/system/389-ds-exporter.service

	cp package/deb/389-ds-exporter.logrotate package/deb/$(PKG_DIR)/etc/logrotate.d/389-ds-exporter
	chmod 644 package/deb/$(PKG_DIR)/etc/logrotate.d/389-ds-exporter

	cp package/deb/postinst package/deb/$(PKG_DIR)/DEBIAN/postinst
	chmod 755 package/deb/$(PKG_DIR)/DEBIAN/postinst

	cp package/deb/prerm package/deb/$(PKG_DIR)/DEBIAN/prerm
	chmod 755 package/deb/$(PKG_DIR)/DEBIAN/prerm

	@echo Building deb package
	fakeroot dpkg-deb --build package/deb/$(PKG_DIR)/ $(BUILD_DIR)/389-ds-exporter_$(VERSION)_$(ARCH).deb

clean:
	rm -rf build/
	rm -rf package/deb/build
