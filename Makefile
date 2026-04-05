BINARY_NAME = pkm
BUILD_DIR   = bin
MODULE      = github.com/mswiente/pkm.ai

.PHONY: build install test test-verbose fmt vet clean

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/pkm

install:
	go install ./cmd/pkm

test:
	go test ./...

test-verbose:
	go test -v ./...

fmt:
	gofmt -w .

vet:
	go vet ./...

clean:
	rm -rf $(BUILD_DIR)

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/pkm

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/pkm
