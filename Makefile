.PHONY: docs docs-generate docs-verify docs-clean build test

# Build all binaries
build:
	go build ./...
	go build ./cmd/autoebiten

# Run tests
test:
	go test -race ./...

# Documentation generation (template-based)
docs:
	go run ./cmd/docgen/main.go docs/generate/*.md.tmpl
