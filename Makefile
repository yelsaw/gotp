VERSION := 0.0.1
OS_BUILDS = linux darwin windows

APP := gotp
BUILD_DIR := "build"

SHA_FILE = sha256.txt
LIC_FILE = LICENSE

GO_LDFLAGS = "-s -extldflags=-static"

# Builds all binaries to BUILD_DIR/{linux,darwin,windows}
build: clean $(OS_BUILDS)

# Builds all binaries, archives with LIC_FILE, and creates sha256sum
dist: clean $(OS_BUILDS) archive checksum

# Builds binary and outputs to BUILD_DIR/linux
linux:
	@GOOS=linux CGO_ENABLED=0 go build -ldflags=$(GO_LDFLAGS) -o build/linux/$(APP)
	@echo "output build/linux/$(APP)"

# Builds binary and outputs to BUILD_DIR/darwin
darwin:
	@GOOS=darwin go build -ldflags=$(GO_LDFLAGS) -o build/darwin/$(APP)
	@echo "output build/darwin/$(APP)"

# Builds binary and outputs to BUILD_DIR/windows
windows:
	@GOOS=windows go build -ldflags=$(GO_LDFLAGS) -o build/windows/$(APP).exe
	@echo "output build/windows/$(APP).exe"

# Creates archives for distribution
archive:
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

# Creates sha256sum for distribution
checksum:
	@echo "Creating checksum hashes"
	@for os in $(OS_BUILDS); do \
			if [ "$$os" = "windows" ]; then \
				sha256sum $(BUILD_DIR)/$(APP)-$$os-v$(VERSION).zip > $(BUILD_DIR)/$(APP)-$$os-v$(VERSION)-$(SHA_FILE); \
			else \
				sha256sum $(BUILD_DIR)/$(APP)-$$os-v$(VERSION).tar > $(BUILD_DIR)/$(APP)-$$os-v$(VERSION)-$(SHA_FILE); \
			fi \
	done

# Removes BUILD_DIR
clean:
	@echo "Removing $(BUILD_DIR) dir"
	@rm -rf $(BUILD_DIR)
