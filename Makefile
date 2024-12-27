VERSION := $(shell git describe --tags --abbrev=0)

OS_BUILDS = linux darwin windows

APP := gotp
BUILD_DIR := build

SHA_FILE = gotp-sha256.txt
LIC_FILE = LICENSE

GO_LDFLAGS = "-s -extldflags=-static"

.PHONY: help build dist linux darwin windows archive checksum verify clean

help:
	@echo ""
	@grep -E '^[sa-zA-Z_-]+:.*#' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*# "; printf "  %-15s %s\n", "Command", "Description"; printf "  %-15s %s\n", "-------", "-----------"} {printf "  %-15s %s\n", $$1, $$2}'
	@echo ""

build: clean $(OS_BUILDS) # Build to BUILD_DIR/{linux,darwin,windows}

dist: clean $(OS_BUILDS) archive checksum # Build bins, create archives, and checksums

linux: # Build bin to BUILD_DIR/linux
	@echo "Building build/linux/$(APP)"
	@GOOS=linux CGO_ENABLED=0 go build -ldflags=$(GO_LDFLAGS) -o $(BUILD_DIR)/linux/$(APP)

darwin: # Build bin to BUILD_DIR/darwin
	@echo "Building build/darwin/$(APP)"
	@GOOS=darwin CGO_ENABLED=0 go build -ldflags=$(GO_LDFLAGS) -o $(BUILD_DIR)/darwin/$(APP)

windows: # Build bin to BUILD_DIR/windows
	@echo "Building build/windows/$(APP).exe"
	@GOOS=windows CGO_ENABLED=0 go build -ldflags=$(GO_LDFLAGS) -o $(BUILD_DIR)/windows/$(APP).exe

archive: # Create archives for distribution
	@echo "Creating tar.gz/zip archives"
	@for os in $(OS_BUILDS); do \
		echo "Archiving $$os archive"; \
		cp $(LIC_FILE) $(BUILD_DIR)/$$os/$(LIC_FILE).txt; \
		if [ "$$os" = "windows" ]; then \
			zip -r $(BUILD_DIR)/$(APP)-$$os-v$(VERSION).zip -j $(BUILD_DIR)/$$os/; \
		else \
			tar -czf $(BUILD_DIR)/$(APP)-$$os-v$(VERSION).tar.gz -C $(BUILD_DIR)/$$os $(APP) $(LIC_FILE).txt; \
		fi \
	done

	@for os in $(OS_BUILDS); do \
		rm -rf $(BUILD_DIR)/$$os; \
	done

checksum: # Create checksum for distribution
	@echo "Creating checksum hashes"
	@rm -f $(BUILD_DIR)/$(SHA_FILE)
	@touch $(BUILD_DIR)/$(SHA_FILE)
	@for os in $(OS_BUILDS); do \
			if [ "$$os" = "windows" ]; then \
				shasum -a 256 $(BUILD_DIR)/$(APP)-$$os-v$(VERSION).zip >> $(BUILD_DIR)/$(SHA_FILE); \
			else \
				shasum -a 256 $(BUILD_DIR)/$(APP)-$$os-v$(VERSION).tar.gz >> $(BUILD_DIR)/$(SHA_FILE); \
			fi \
	done
	@sed -i.del 's/build\///g' $(BUILD_DIR)/$(SHA_FILE) && rm -f $(BUILD_DIR)/$(SHA_FILE).del

verify: # Verify checksums
	@echo "Verifying checksum hashes"
	@cd $(BUILD_DIR) && shasum -a 256 -c $(SHA_FILE)

clean: # Remove BUILD_DIR
	@if [ -d $(BUILD_DIR) ]; then \
		echo "Removing $(BUILD_DIR) dir"; \
		rm -rf $(BUILD_DIR); \
	fi \

