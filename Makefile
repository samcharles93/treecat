BINARY_NAME=treecat
DIST_DIR=dist
PLATFORMS=linux/amd64 windows/amd64 darwin/amd64 darwin/arm64

VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT  ?= $(shell git rev-parse --short HEAD)
TIME    ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

build-prod:
    go build -ldflags="-s -w \
    -X main.Version=$(VERSION) \
    -X main.BuildTime=$(TIME) \
    -X main.GitCommit=$(COMMIT)" \
    -trimpath -o $(BINARY_NAME) .

build-all:
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; \
		arch=$${platform#*/}; \
		output=$(DIST_DIR)/$(BINARY_NAME)-$$os-$$arch; \
		if [ "$$os" = "windows" ]; then output=$$output.exe; fi; \
		GOOS=$$os GOARCH=$$arch go build -ldflags="-s -w" -trimpath -o $$output .; \
	done

clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe
	rm -rf $(DIST_DIR)

test:
	go test ./...

install-upx:
	@echo "Installing UPX..."
	@if ! command -v upx &> /dev/null; then \
		sudo apt-get install -y upx; \
	else \
		echo "UPX is already installed."; \
	fi

compress:
	@for binary in $(wildcard $(DIST_DIR)/*); do \
		upx --best $$binary; \
	done

.DEFAULT_GOAL := build-prod