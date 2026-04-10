// internal/docgen/cmd/generate.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/s3cy/autoebiten/internal/docgen"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run ./internal/docgen/cmd/generate.go <templates...>")
		fmt.Fprintln(os.Stderr, "  example: go run ./internal/docgen/cmd/generate.go docs/generate/*.md.gotmpl")
		os.Exit(1)
	}

	// Process each template
	for _, tmplPath := range os.Args[1:] {
		// Expand glob if needed
		if strings.Contains(tmplPath, "*") {
			matches, err := filepath.Glob(tmplPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error expanding glob %s: %v\n", tmplPath, err)
				continue
			}
			for _, m := range matches {
				processTemplate(m)
			}
		} else {
			processTemplate(tmplPath)
		}
	}
}

func processTemplate(tmplPath string) {
	fmt.Printf("Processing: %s\n", tmplPath)

	result, err := docgen.ProcessTemplate(tmplPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", tmplPath, err)
		os.Exit(1)
	}

	// Output path: docs/generate/X.md.gotmpl -> docs/X.md
	outPath := outputPath(tmplPath)

	// Ensure output directory exists
	outDir := filepath.Dir(outPath)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Write output
	if err := os.WriteFile(outPath, []byte(result), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", outPath, err)
		os.Exit(1)
	}

	fmt.Printf("Generated: %s -> %s\n", tmplPath, outPath)
}

func outputPath(tmplPath string) string {
	// docs/generate/commands.md.tmpl -> docs/commands.md
	base := filepath.Base(tmplPath)
	name := strings.TrimSuffix(base, ".md.tmpl")
	return filepath.Join("docs", name+".md")
}