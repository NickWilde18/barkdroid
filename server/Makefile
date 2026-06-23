ROOT_DIR = $(shell pwd)

.PHONY: test lint build docker clean

# Run all tests
test:
	go test -v -count=1 -timeout 30s ./...

# Run linter
lint:
	golangci-lint run ./...

# Build binary
build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o barkdroid .

# Build Docker image
docker:
	docker build -t barkdroid .

# Run with docker compose
up:
	docker compose up -d

# Clean build artifacts
clean:
	rm -f barkdroid
	rm -rf data/
