VERSION := $(shell git describe --tags --abbrev=0)

OS_BUILDS = linux darwin windows

APP := gotp
BUILD_DIR := "build"

SHA_FILE = sha256.txt
LIC_FILE = LICENSE

GO_LDFLAGS = "-s -extldflags=-static"

.PHONY: help clean build dist linux darwin windows archive checksum

help:
	@echo ""
	@grep -E '^[sa-zA-Z_-]+:.*#' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*# "; printf "  %-15s %s\n", "Command", "Description"; printf "  %-15s %s\n", "-------", "-----------"} {printf "  %-15s %s\n", $$1, $$2}'
	@echo ""

build: clean $(OS_BUILDS) # Build to BUILD_DIR/{linux,darwin,windows}

dist: clean $(OS_BUILDS) archive checksum # Build bins, create archives, and checksums

linux: # Build bin to BUILD_DIR/linux
	@GOOS=linux CGO_ENABLED=0 go build -ldflags=$(GO_LDFLAGS) -o $(BUILD_DIR)/linux/$(APP)
	@echo "output build/linux/$(APP)"

darwin: # Build bin to BUILD_DIR/darwin
	@GOOS=darwin go build -ldflags=$(GO_LDFLAGS) -o $(BUILD_DIR)/darwin/$(APP)
	@echo "output build/darwin/$(APP)"

windows: # Build bin to BUILD_DIR/windows
	@GOOS=windows go build -ldflags=$(GO_LDFLAGS) -o $(BUILD_DIR)/windows/$(APP).exe
	@echo "output build/windows/$(APP).exe"

archive: # Create archives for distribution
	@echo "Creating tar/zip archives"
	@for os in $(OS_BUILDS); do \
			echo "Creating zip archive for $$os"; \
		if [ "$$os" = "windows" ]; then \
			cp -a $(LIC_FILE) $(BUILD_DIR)/$$os/$(LIC_FILE).txt; \
			zip -r $(BUILD_DIR)/$(APP)-$$os-v$(VERSION).zip -j $(BUILD_DIR)/$$os/; \
		else \
			cp -a $(LIC_FILE) $(BUILD_DIR)/$$os/; \
			tar -cf $(BUILD_DIR)/$(APP)-$$os-v$(VERSION).tar -C $(BUILD_DIR)/$$os $(APP) $(LIC_FILE); \
		fi \
	done

checksum: # Create sha256sum(s) for distribution
	@echo "Creating checksum hashes"
	@for os in $(OS_BUILDS); do \
			if [ "$$os" = "windows" ]; then \
				sha256sum $(BUILD_DIR)/$(APP)-$$os-v$(VERSION).zip > $(BUILD_DIR)/$(APP)-$$os-v$(VERSION)-$(SHA_FILE); \
			else \
				sha256sum $(BUILD_DIR)/$(APP)-$$os-v$(VERSION).tar > $(BUILD_DIR)/$(APP)-$$os-v$(VERSION)-$(SHA_FILE); \
			fi \
	done

clean: # Remove BUILD_DIR
	@if [ -d $(BUILD_DIR) ]; then \
		echo "Removing $(BUILD_DIR) dir"; \
		rm -rf $(BUILD_DIR); \
	fi \

