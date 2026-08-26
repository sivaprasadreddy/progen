MAIN_SRC_DIR=.
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
	@go run main.go

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
	@rm -rf dist

.PHONY: all build run test clean
