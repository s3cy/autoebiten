.PHONY: docs-generate docs-verify docs-clean build test

# Build all binaries
build:
	go build ./...
	go build ./cmd/autoebiten
	go build ./cmd/docgen
	go build ./cmd/docverify

# Run tests
test:
	go test -race ./...

# Generate all example outputs and process templates
docs-generate:
	go run ./cmd/docgen --generate docs/generate/autoui_examples
	go run ./cmd/docgen --generate docs/generate/commands_examples
	go run ./cmd/docgen --generate docs/generate/integration_examples
	go run ./cmd/docgen --process

# Verify outputs match recorded files
docs-verify:
	go run ./cmd/docverify docs/generate/autoui_examples
	go run ./cmd/docverify docs/generate/commands_examples
	go run ./cmd/docverify docs/generate/integration_examples

# Clean generated files
docs-clean:
	rm -f docs/generate/autoui_examples/*_out.txt
	rm -f docs/generate/commands_examples/*_out.txt
	rm -f docs/generate/integration_examples/*_out.txt
	rm -f docs/autoui.md docs/commands.md docs/integration.md
