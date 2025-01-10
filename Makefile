# Variables
PROJECT_NAME := golang_web_service
CMD_DIR := cmd
BUILD_DIR := bin

# Binaries
ORION_BINARY := $(BUILD_DIR)/orion
ORIONCLI_BINARY := $(BUILD_DIR)/orioncli

# Go Settings
GO := go
GOFLAGS := -mod=readonly

.PHONY: all clean build run server client

# Default target
all: build

# Build targets
build: $(ORION_BINARY) $(ORIONCLI_BINARY)

$(ORION_BINARY): $(CMD_DIR)/orion/main.go
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $@ $<

$(ORIONCLI_BINARY): $(CMD_DIR)/orioncli/main.go
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $@ $<

# Run server
run: server

server: $(ORION_BINARY)
	$(ORION_BINARY)

# Run client
client: $(ORIONCLI_BINARY)
	$(ORIONCLI_BINARY)

# Clean build artifacts
clean:
	@rm -rf $(BUILD_DIR)