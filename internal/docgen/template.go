package docgen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// OutputFunc returns a template function that reads a file and wraps it in a code fence.
// The path is relative to baseDir.
func OutputFunc(baseDir string) func(string, string) (string, error) {
	return func(path, lang string) (string, error) {
		fullPath := filepath.Join(baseDir, path)
		content, err := readFile(fullPath)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("```%s\n%s\n```", lang, strings.TrimSpace(content)), nil
	}
}

// readFile reads a file and returns its content.
func readFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read output file %s: %w", path, err)
	}
	return string(data), nil
}

// ProcessTemplate reads a template file, processes it with the output function, and returns the result.
func ProcessTemplate(templatePath, baseDir string) (string, error) {
	content, err := readFile(templatePath)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(template.FuncMap{
		"output": OutputFunc(baseDir),
	}).Parse(content)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, nil); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return result.String(), nil
}
