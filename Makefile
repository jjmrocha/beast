.PHONY: help build clean test bench lint deps tidy
.DEFAULT_GOAL := help

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  build        Build beast into ./bin"
	@echo "  clean        Remove the build output in ./bin"
	@echo "  test         Run all tests"
	@echo "  bench        Run benchmarks"
	@echo "  lint         Run golangci-lint"
	@echo "  deps         Update dependencies"
	@echo "  tidy         Tidy go.mod"
	@echo ""
	@echo "CI runs on GitHub Actions (.github/workflows/ci.yml)."

build:
	go build -o bin/beast .

clean:
	rm -rf bin

test:
	go test -race ./...

bench:
	go test -bench=. ./...

lint:
	golangci-lint run

deps:
	go get -u ./...

tidy:
	go mod tidy
