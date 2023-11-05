MAIN_SRC_DIR=.
DIST_DIR=dist
BINARY_NAME=progen
BINARY_LINUX=$(BINARY_NAME)_linux
BINARY_WIN=$(BINARY_NAME)_win.exe
GOARCH="amd64"

# Build the application
all: build

# format the go source files
fmt:
	goimports -w .

build:
	@echo 'Building MacOS binary'
	GOARCH=${GOARCH} GOOS=darwin go build -o ${DIST_DIR}/${BINARY_NAME}-darwin-${GOARCH} ${MAIN_SRC_DIR}
	@echo 'Building Linux binary'
	GOARCH=${GOARCH} GOOS=linux go build -o ${DIST_DIR}/${BINARY_NAME}-linux-${GOARCH} ${MAIN_SRC_DIR}
	@echo 'Building Windows binary'
	GOARCH=${GOARCH} GOOS=windows go build -o ${DIST_DIR}/${BINARY_NAME}-windows-${GOARCH}.exe ${MAIN_SRC_DIR}

# Run the application
run:
	@go run main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./...

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f progen
	@rm -rf dist

.PHONY: all build run test clean
