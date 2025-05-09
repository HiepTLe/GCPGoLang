.PHONY: build clean test run fmt lint docker-build

# Build settings
BINARY_NAME=gcpgolang
VERSION?=1.0.0
LDFLAGS=-ldflags "-X main.Version=${VERSION}"
GOARCH?=amd64

# Directory structure
SRC_DIR=.
BIN_DIR=./bin

# Default target
all: clean build

# Build the application
build:
	@echo "Building ${BINARY_NAME}..."
	@go build ${LDFLAGS} -o ${BINARY_NAME} ${SRC_DIR}

# Build for multiple platforms
build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p ${BIN_DIR}
	@GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BIN_DIR}/${BINARY_NAME}-darwin-${GOARCH} ${SRC_DIR}
	@GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BIN_DIR}/${BINARY_NAME}-linux-${GOARCH} ${SRC_DIR}
	@GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BIN_DIR}/${BINARY_NAME}-windows-${GOARCH}.exe ${SRC_DIR}

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f ${BINARY_NAME}
	@rm -rf ${BIN_DIR}

# Run the application 
run:
	@go run ${SRC_DIR}

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run code linting
lint:
	@echo "Running linter..."
	@go vet ./...
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Skipping extended linting."; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t ${BINARY_NAME}:${VERSION} .

# Run all checks before committing
pre-commit: fmt lint test

# Run a specific scan
scan:
	@./$(BINARY_NAME) misconfig-scanner --project=$(PROJECT_ID) $(SCAN_ARGS)

# Run a scan with Wiz integration
scan-wiz:
	@./$(BINARY_NAME) misconfig-scanner --project=$(PROJECT_ID) --wiz --wiz-client-id=$(WIZ_CLIENT_ID) --wiz-client-secret=$(WIZ_CLIENT_SECRET) $(SCAN_ARGS)

# Help target
help:
	@echo "GCPGoLang Makefile Help"
	@echo "----------------------"
	@echo "make              : Build the application"
	@echo "make build-all    : Build for multiple platforms"
	@echo "make clean        : Clean build artifacts"
	@echo "make run          : Run the application"
	@echo "make test         : Run tests"
	@echo "make fmt          : Format the code"
	@echo "make lint         : Run linting tools"
	@echo "make pre-commit   : Run all checks before committing"
	@echo "make docker-build : Build Docker image"
	@echo "make scan PROJECT_ID=your-project-id [SCAN_ARGS=\"--verbose\"] : Run a scan"
	@echo "make scan-wiz PROJECT_ID=your-project-id WIZ_CLIENT_ID=your-id WIZ_CLIENT_SECRET=your-secret [SCAN_ARGS=\"--verbose\"] : Run a scan with Wiz" 