VERSION := $(shell git tag -l | tail -1)

OS_BUILDS = linux darwin windows

APP := gotp
BUILD_DIR := build
DIST_DIR := dist

SHA_ALGO ?= 256

SHA_FILE = gotp-sha$(SHA_ALGO).txt
LIC_FILE = LICENSE

GO_LDFLAGS = "-s -extldflags=-static"

MAKE := $(MAKE) --no-print-directory

.PHONY: help build dist linux darwin windows archive checksum verify clean

help:
	@echo ""
	@awk '/^[a-zA-Z_-]+:.*#/ {split($$0, a, ":.*# "); printf "  %-15s %s\n", a[1], a[2]}' $(MAKEFILE_LIST) | \
	awk 'BEGIN {printf "  %-15s %s\n", "Command", "Description"; printf "  %-15s %s\n", "-------", "-----------"} {print}'
	@echo ""

build: # Build to BUILD_DIR/{linux,darwin,windows}
	@$(MAKE) clean && $(MAKE) $(OS_BUILDS)

dist: # Build bins, create archives, and checksums
	@$(MAKE) build && $(MAKE) archive && $(MAKE) checksum

linux: # Build bin to BUILD_DIR/linux
	@echo "Building build/linux/$(APP)"
	@env GOOS=linux CGO_ENABLED=0 go build -ldflags=$(GO_LDFLAGS) -o $(BUILD_DIR)/linux/$(APP)

darwin: # Build bin to BUILD_DIR/darwin
	@echo "Building build/darwin/$(APP)"
	@env GOOS=darwin CGO_ENABLED=0 go build -ldflags=$(GO_LDFLAGS) -o $(BUILD_DIR)/darwin/$(APP)

windows: # Build bin to BUILD_DIR/windows
	@echo "Building build/windows/$(APP).exe"
	@env GOOS=windows CGO_ENABLED=0 go build -ldflags=$(GO_LDFLAGS) -o $(BUILD_DIR)/windows/$(APP).exe

archive: # Create archives for distribution
	@echo "Creating tar.gz/zip archives"
	@for os in $(OS_BUILDS); do \
		mkdir -p $(DIST_DIR); \
		echo "Archiving $$os archive"; \
		cp $(LIC_FILE) $(BUILD_DIR)/$$os/$(LIC_FILE).txt; \
		if [ "$$os" = "windows" ]; then \
			zip -r $(BUILD_DIR)/$(APP)-$$os-v$(VERSION).zip -j $(BUILD_DIR)/$$os/; \
			mv $(BUILD_DIR)/$(APP)-$$os-v$(VERSION).zip $(DIST_DIR)/; \
		else \
			tar -czf $(BUILD_DIR)/$(APP)-$$os-v$(VERSION).tar.gz -C $(BUILD_DIR)/$$os $(APP) $(LIC_FILE).txt; \
			mv $(BUILD_DIR)/$(APP)-$$os-v$(VERSION).tar.gz $(DIST_DIR)/; \
		fi \
	done

	@for os in $(OS_BUILDS); do \
		rm -rf $(BUILD_DIR)/$$os; \
	done

checksum: # Create checksum for distribution
	@echo "Creating checksum $(SHA_ALGO) hashes"
	@rm -f $(DIST_DIR)/$(SHA_FILE)
	@touch $(DIST_DIR)/$(SHA_FILE)
	@for os in $(OS_BUILDS); do \
			if [ "$$os" = "windows" ]; then \
				shasum -a $(SHA_ALGO) $(DIST_DIR)/$(APP)-$$os-v$(VERSION).zip >> $(DIST_DIR)/$(SHA_FILE); \
			else \
				shasum -a $(SHA_ALGO) $(DIST_DIR)/$(APP)-$$os-v$(VERSION).tar.gz >> $(DIST_DIR)/$(SHA_FILE); \
			fi \
	done
	@perl -pi -e 's/$(DIST_DIR)\///g' $(DIST_DIR)/$(SHA_FILE)
	@rm -rf $(BUILD_DIR)

verify: # Verify checksums
	@echo "Verifying checksum hashes"
	@cd $(DIST_DIR) && shasum -a $(SHA_ALGO) -c $(SHA_FILE)

clean: # Remove DIST_DIR BUILD_DIR
	@echo "Removing $(DIST_DIR) and $(BUILD_DIR) dirs"
	@rm -rf $(DIST_DIR) $(BUILD_DIR)

