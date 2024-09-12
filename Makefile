VERSION := 0.0.1
OS_BUILDS = linux darwin windows

APP := gotp
BUILD_DIR := "build"
SHA_FILE := $(APP)-sha256.txt
LIC_FILE := LICENSE

GO_LDFLAGS = "-s -extldflags=-static"

all: rmbuild $(OS_BUILDS) checksum archive

linux:
	@GOOS=linux CGO_ENABLED=0 go build -ldflags=$(GO_LDFLAGS) -o build/linux/$(APP)
	@echo "output build/linux/$(APP)"

darwin:
	@GOOS=darwin go build -ldflags=$(GO_LDFLAGS) -o build/darwin/$(APP)
	@echo "output build/darwin/$(APP)"

windows:
	@GOOS=windows go build -ldflags=$(GO_LDFLAGS) -o build/windows/$(APP).exe
	@echo "output build/windows/$(APP).exe"

archive:
	@echo "Creating tar/zip archives"
	@for os in $(OS_BUILDS); do \
		if [ "$$os" = "windows" ]; then \
			echo "Creating zip archive for $$os"; \
			cp -a $(LIC_FILE) $(BUILD_DIR)/$$os/$(LIC_FILE).txt; \
			zip -r $(BUILD_DIR)/$(APP)-$$os-v$(VERSION).zip -j $(BUILD_DIR)/$$os/; \
		else \
			echo "Creating tar archive for $$os"; \
			cp -a $(LIC_FILE) $(BUILD_DIR)/$$os/; \
			tar -cf $(BUILD_DIR)/$(APP)-$$os-v$(VERSION).tar -C $(BUILD_DIR)/$$os $(APP) $(LIC_FILE) $(SHA_FILE); \
		fi \
	done

checksum:
	@echo "Creating checksum hashes"
	@for os in $(OS_BUILDS); do \
		( \
			cd $(BUILD_DIR)/$$os && \
			if [ "$$os" = "windows" ]; then \
				sha256sum $(APP).exe > $(SHA_FILE); \
			else \
				sha256sum $(APP) > $(SHA_FILE); \
			fi \
		) \
	done

rmbuild:
	@rm -rf $(BUILD_DIR)
	@echo "Removing $(BUILD_DIR) dir"
