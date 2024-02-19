BUILD_VERSION?=1.0.5

BUILD_NAME:=remote
BUILD_DATE:=$(shell date '+%Y%m%d.%H%M%S')
BUILD_DIR:=build
DEST_DIR:=dest
SRC_FILE:=*.go

ifdef IS_RELEASE
	SOFT_VERSION:=$(BUILD_VERSION)
else
	SOFT_VERSION:=$(BUILD_VERSION)-$(BUILD_DATE)
endif

LDFLAGS:=-ldflags "-s -w -X remote_exec/util.VERSION=$(SOFT_VERSION)"

.PHONY: all clean prepare test win_x86 linux_x86 linux_arm package

all:clean win_x86 linux_x86 linux_arm package

prepare:
	@echo "====> prepare build environment"
	- mkdir $(DEST_DIR) $(BUILD_DIR)
	go mod tidy

test: prepare
	@echo "====> run unit case"
	- go test ./...

win_x86: prepare test
	@echo "====> build windows amd64"
	export GO111MODULE=on \
		&& export CGO_ENABLED=0 \
		&& export GOOS=windows \
		&& export GOARCH=amd64 \
		&& go build $(LDFLAGS) -o $(BUILD_DIR)/$(BUILD_NAME)_windows_amd64.exe $(SRC_FILE)
	upx $(BUILD_DIR)/$(BUILD_NAME)_windows_amd64.exe

linux_x86: prepare test
	@echo "====> build linux amd64"
	export GO111MODULE=on \
		&& export CGO_ENABLED=0 \
		&& export GOOS=linux \
		&& export GOARCH=amd64 \
		&& go build $(LDFLAGS) -o $(BUILD_DIR)/$(BUILD_NAME)_linux_amd64 $(SRC_FILE)
	upx $(BUILD_DIR)/$(BUILD_NAME)_linux_amd64

linux_arm: prepare test
	@echo "====> build linux arm64"
	export GO111MODULE=on \
		&& export CGO_ENABLED=0 \
		&& export GOOS=linux \
		&& export GOARCH=arm64 \
		&& go build $(LDFLAGS) -o $(BUILD_DIR)/$(BUILD_NAME)_linux_arm64 $(SRC_FILE)
	upx $(BUILD_DIR)/$(BUILD_NAME)_linux_arm64

package :
	@echo "====> archive artifact"
	cp config.yaml $(BUILD_DIR)/
	cp command.yaml $(BUILD_DIR)/
	cp README.md $(BUILD_DIR)/
	cd $(BUILD_DIR) && zip $(BUILD_NAME)-$(SOFT_VERSION).zip * && cd ..
	mv $(BUILD_DIR)/$(BUILD_NAME)-$(SOFT_VERSION).zip $(DEST_DIR)/

clean:
	@echo "====> clean temporary"
	- rm -rf $(BUILD_DIR)
	- rm -rf $(DEST_DIR)
	go clean