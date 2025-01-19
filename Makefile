# Variables
BINARY_NAME=treecat
DIST_DIR=dist
PLATFORMS=linux/amd64 windows/amd64 darwin/amd64 darwin/arm64

# Build for the current platform
build:
	go build -o $(BINARY_NAME) .

# Build for production (stripped and trimmed)
build-prod:
	go build -ldflags="-s -w" -trimpath -o $(BINARY_NAME) .

# Build for multiple platforms
build-all:
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; \
		arch=$${platform#*/}; \
		output=$(DIST_DIR)/$(BINARY_NAME)-$$os-$$arch; \
		if [ "$$os" = "windows" ]; then output=$$output.exe; fi; \
		GOOS=$$os GOARCH=$$arch go build -ldflags="-s -w" -trimpath -o $$output .; \
	done

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe
	rm -rf $(DIST_DIR)

# Run tests
test:
	go test ./...

# Run the program locally
run:
	go run .

# Install UPX for binary compression (optional)
install-upx:
	@echo "Installing UPX..."
	@if ! command -v upx &> /dev/null; then \
		sudo apt-get install -y upx; \
	else \
		echo "UPX is already installed."; \
	fi

# Compress binaries with UPX
compress:
	@for binary in $(wildcard $(DIST_DIR)/*); do \
		upx --best $$binary; \
	done

# Default task
.DEFAULT_GOAL := build
