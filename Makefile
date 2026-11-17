MAIN_SRC_DIR=./cmd/progen
DESKTOP_SRC_DIR=./cmd/desktop
DIST_DIR=dist
BINARY_NAME=progen
GOARCHES=amd64 arm64

# Build the application
all: build

# format the go source files
fmt:
	go fmt ./...

build:
	@for arch in ${GOARCHES}; do \
		echo "Building MacOS binary ($$arch)"; \
		GOARCH=$$arch GOOS=darwin go build -o ${DIST_DIR}/${BINARY_NAME}-darwin-$$arch ${MAIN_SRC_DIR}; \
		echo "Building Linux binary ($$arch)"; \
		GOARCH=$$arch GOOS=linux go build -o ${DIST_DIR}/${BINARY_NAME}-linux-$$arch ${MAIN_SRC_DIR}; \
		echo "Building Windows binary ($$arch)"; \
		GOARCH=$$arch GOOS=windows go build -o ${DIST_DIR}/${BINARY_NAME}-windows-$$arch.exe ${MAIN_SRC_DIR}; \
	done

# Run the application
run:
	@go run ${MAIN_SRC_DIR}

# Run the Progen Desktop app, rebuilding as the code changes
desktop_dev:
	go tool mygo dev ${DESKTOP_SRC_DIR}

# Build Progen Desktop.app and a DMG into cmd/desktop/build/
desktop_build:
	go tool mygo build ${DESKTOP_SRC_DIR}

# Run basic tests only
test:
	@echo "Testing..."
	go test -v -p 1 -short ./...

# Run all tests
all_tests:
	@echo "Testing..."
	go test -v -p 1 ./...

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -rf dist ${DESKTOP_SRC_DIR}/build

.PHONY: all build run desktop_dev desktop_build test all_tests clean
