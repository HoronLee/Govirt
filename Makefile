# 定义变量
BINARY_NAME := govirt
SRC_DIR := .
VERSION := $(shell git describe --tags --abbrev=0 --match 'v*')
COMMIT := $(shell git rev-parse --short HEAD)
EXTERNAL_VERSION ?= $(VERSION)
OUTPUT_DIR := bin

# 支持的目标平台
PLATFORMS := windows/amd64 linux/amd64 linux/arm64 darwin/amd64 darwin/arm64

# 默认目标
.PHONY: all
all: $(PLATFORMS)

# 构建多平台二进制文件
.PHONY: $(PLATFORMS)
$(PLATFORMS):
	$(eval GOOS := $(word 1, $(subst /, ,$@)))
	$(eval GOARCH := $(word 2, $(subst /, ,$@)))
	$(eval OUTPUT_NAME := $(BINARY_NAME)-$(EXTERNAL_VERSION)-$(GOOS)-$(GOARCH)$(if $(filter $(GOOS), windows),.exe))
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(OUTPUT_DIR)/$(OUTPUT_NAME) $(SRC_DIR)

# 清理构建文件
.PHONY: clean
clean:
	rm -rf $(OUTPUT_DIR)

# 运行二进制文件 (仅适用于当前平台构建)
.PHONY: run
run: build
	./$(OUTPUT_DIR)/$(BINARY_NAME)-$(EXTERNAL_VERSION)-$(GOOS)-$(GOARCH)$(if $(filter $(GOOS), windows),.exe)

# 安装依赖
.PHONY: deps
deps:
	go mod tidy

# 开发构建
.PHONY: dev
dev:
	go build -o $(OUTPUT_DIR)/$(BINARY_NAME)-dev
	./$(OUTPUT_DIR)/$(BINARY_NAME)-dev

# 帮助信息
.PHONY: help
help:
	@echo "可用目标:"
	@echo "  all          构建所有平台的二进制文件"
	@echo "  clean        清理构建文件"
	@echo "  run          运行当前平台的二进制文件"
	@echo "  deps         安装依赖"
	@echo "  dev          开发构建并运行"
	@echo "  help         显示帮助信息"