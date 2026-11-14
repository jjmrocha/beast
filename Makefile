.PHONY: help build snapshot clean test bench lint vuln deps tidy
.DEFAULT_GOAL := help

VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo dev)
LDFLAGS := -X github.com/jjmrocha/beast/internal/command.version=$(VERSION)

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  build        Build beast into ./bin"
	@echo "  snapshot     Build the release files into ./dist without publishing"
	@echo "  clean        Remove the build output in ./bin and ./dist"
	@echo "  test         Run all tests"
	@echo "  bench        Run benchmarks"
	@echo "  lint         Run golangci-lint"
	@echo "  vuln         Run govulncheck"
	@echo "  deps         Update dependencies"
	@echo "  tidy         Tidy go.mod"
	@echo ""
	@echo "CI runs on GitHub Actions (.github/workflows/ci.yml)."

build:
	go build -ldflags "$(LDFLAGS)" -o bin/beast ./cmd/beast

snapshot:
	goreleaser release --snapshot --clean --skip=sign

clean:
	rm -rf bin dist

test:
	go test -race ./...

bench:
	go test -bench=. ./...

lint:
	golangci-lint run

vuln:
	govulncheck ./...

deps:
	go get -u ./...

tidy:
	go mod tidy
