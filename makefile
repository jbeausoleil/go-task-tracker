# Project binary name
BINARY_NAME=task-cli

# Version info (commit and date will be injected at build time)
VERSION=0.0.1
COMMIT=$(shell git rev-parse --short HEAD || echo "dev")
DATE=$(shell date +%Y-%m-%dT%H:%M:%S%z)

# Go files
MAIN=main.go

# Default build for current OS/Arch
build:
	@echo "Building for current system..."
	go build -o $(BINARY_NAME) \
	-ldflags "-X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)' -X 'main.date=$(DATE)'" \
	$(MAIN)

# Build for macOS
build-mac:
	@echo "Building for macOS (amd64)..."
	GOOS=darwin GOARCH=amd64 go build -o dist/$(BINARY_NAME)-mac \
	-ldflags "-X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)' -X 'main.date=$(DATE)'" \
	$(MAIN)

# Build for Linux
build-linux:
	@echo "Building for Linux (amd64)..."
	GOOS=linux GOARCH=amd64 go build -o dist/$(BINARY_NAME)-linux \
	-ldflags "-X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)' -X 'main.date=$(DATE)'" \
	$(MAIN)

# Build both (cross-compile)
cross-compile: build-mac build-linux

# Clean up build files
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -rf dist