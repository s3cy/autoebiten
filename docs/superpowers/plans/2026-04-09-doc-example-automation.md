# Doc Example Automation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Create a system that automates documentation examples by running real commands, capturing outputs, and generating docs from templates.

**Architecture:** Three components: (1) Example scripts that users can run directly, (2) docgen tool that runs scripts and processes templates, (3) docverify tool that compares outputs. Uses testkit for reliable game lifecycle management.

**Tech Stack:** Go, text/template, gopkg.in/yaml.v3, testkit package

---

## Task 1: Create Directory Structure

**Files:**
- Create: `docs/generate/` directory
- Create: `docs/generate/autoui_examples/` directory

- [ ] **Step 1: Create docs/generate directory**

```bash
mkdir -p docs/generate
```

- [ ] **Step 2: Create autoui_examples subdirectory**

```bash
mkdir -p docs/generate/autoui_examples
```

- [ ] **Step 3: Verify structure**

```bash
ls -la docs/generate/
```

Expected output shows `autoui_examples` directory.

---

## Task 2: Create Config Loader Package

**Files:**
- Create: `internal/docgen/config.go`
- Create: `internal/docgen/config_test.go`

- [ ] **Step 1: Write failing test for config loading**

```go
// internal/docgen/config_test.go
package docgen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	// Create temp config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `game_dir: examples/autoui
normalize:
  - pattern: '_addr="[^"]*"'
    replace: '_addr="<ADDR>"'
  - pattern: '\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+'
    replace: '<TIMESTAMP>'
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	config, err := LoadConfig(configPath)
	require.NoError(t, err)

	assert.Equal(t, "examples/autoui", config.GameDir)
	assert.Len(t, config.Normalize, 2)
	assert.Equal(t, `_addr="[^"]*"`, config.Normalize[0].Pattern)
	assert.Equal(t, "_addr=\"<ADDR>\"", config.Normalize[0].Replace)
}

func TestLoadConfigNotFound(t *testing.T) {
	_, err := LoadConfig("/nonexistent/config.yaml")
	assert.Error(t, err)
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/docgen -run TestLoadConfig -v
```

Expected: FAIL with "package docgen not found" or similar.

- [ ] **Step 3: Create package directory**

```bash
mkdir -p internal/docgen
```

- [ ] **Step 4: Write config types**

```go
// internal/docgen/config.go
package docgen

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// NormalizeRule defines a regex pattern and replacement for output normalization.
type NormalizeRule struct {
	Pattern string `yaml:"pattern"`
	Replace string `yaml:"replace"`
}

// Config defines the example generation settings for a documentation section.
type Config struct {
	GameDir   string          `yaml:"game_dir"`
	Normalize []NormalizeRule `yaml:"normalize"`
}

// LoadConfig reads a config.yaml file and returns the parsed Config.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}
```

- [ ] **Step 5: Run test to verify it passes**

```bash
go test ./internal/docgen -run TestLoadConfig -v
```

Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add internal/docgen/
git commit -m "feat(docgen): add config loader for example generation"
```

---

## Task 3: Create Normalization Function

**Files:**
- Modify: `internal/docgen/config.go`
- Modify: `internal/docgen/config_test.go`

- [ ] **Step 1: Write failing test for normalization**

```go
// Add to internal/docgen/config_test.go

func TestNormalize(t *testing.T) {
	rules := []NormalizeRule{
		{Pattern: `_addr="[^"]*"`, Replace: `_addr="<ADDR>"`},
		{Pattern: `\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+`, Replace: `<TIMESTAMP>`},
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normalize address",
			input:    `<Button _addr="0x14000123456" id="btn"/>`,
			expected: `<Button _addr="<ADDR>" id="btn"/>`,
		},
		{
			name:     "normalize timestamp",
			input:    `2026-04-09 14:30:15.123 game started`,
			expected: `<TIMESTAMP> game started`,
		},
		{
			name:     "normalize multiple",
			input:    `<Button _addr="0x14000123456"/> at 2026-04-09 14:30:15.123`,
			expected: `<Button _addr="<ADDR>"/> at <TIMESTAMP>`,
		},
		{
			name:     "no matches",
			input:    `plain text`,
			expected: `plain text`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Normalize(tt.input, rules)
			assert.Equal(t, tt.expected, result)
		})
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/docgen -run TestNormalize -v
```

Expected: FAIL with "Normalize not defined"

- [ ] **Step 3: Write normalize function**

```go
// Add to internal/docgen/config.go

import (
	"regexp"
	"strings"
)

// Normalize applies all normalization rules to the input string.
func Normalize(s string, rules []NormalizeRule) string {
	for _, r := range rules {
		re := regexp.MustCompile(r.Pattern)
		s = re.ReplaceAllString(s, r.Replace)
	}
	return strings.TrimSpace(s)
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./internal/docgen -run TestNormalize -v
```

Expected: PASS

- [ ] **Step 5: Run all config tests**

```bash
go test ./internal/docgen -v
```

Expected: All PASS

- [ ] **Step 6: Commit**

```bash
git add internal/docgen/
git commit -m "feat(docgen): add normalization function for dynamic output values"
```

---

## Task 4: Create Template Processing Package

**Files:**
- Create: `internal/docgen/template.go`
- Create: `internal/docgen/template_test.go`

- [ ] **Step 1: Write failing test for output function**

```go
// internal/docgen/template_test.go
package docgen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutputFunc(t *testing.T) {
	// Create temp output file
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test_out.txt")

	content := `<UI>
  <Button id="btn"/>
</UI>`
	err := os.WriteFile(outputPath, []byte(content), 0644)
	require.NoError(t, err)

	// Create output function with base dir
	outputFn := OutputFunc(tmpDir)

	// Test reading file
	result, err := outputFn("test_out.txt", "xml")
	require.NoError(t, err)

	expected := "```xml\n<UI>\n  <Button id=\"btn\"/>\n</UI>\n```"
	assert.Equal(t, expected, result)
}

func TestOutputFuncNotFound(t *testing.T) {
	outputFn := OutputFunc("/nonexistent")
	_, err := outputFn("missing.txt", "text")
	assert.Error(t, err)
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/docgen -run TestOutputFunc -v
```

Expected: FAIL

- [ ] **Step 3: Write output function**

```go
// internal/docgen/template.go
package docgen

import (
	"fmt"
	"path/filepath"
	"strings"
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
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./internal/docgen -run TestOutputFunc -v
```

Expected: PASS

- [ ] **Step 5: Write test for template processing**

```go
// Add to internal/docgen/template_test.go

func TestProcessTemplate(t *testing.T) {
	// Create temp files
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "example_out.txt")
	templatePath := filepath.Join(tmpDir, "test.md.tmpl")

	outputContent := `<Button id="btn"/>`
	err := os.WriteFile(outputPath, []byte(outputContent), 0644)
	require.NoError(t, err)

	templateContent := `# Test
**Output:**
{{output "example_out.txt" "xml"}}
`
	err = os.WriteFile(templatePath, []byte(templateContent), 0644)
	require.NoError(t, err)

	// Process template
	result, err := ProcessTemplate(templatePath, tmpDir)
	require.NoError(t, err)

	expected := `# Test
**Output:**
```xml
<Button id="btn"/>
```
`
	assert.Equal(t, expected, result)
}
```

- [ ] **Step 6: Run test to verify it fails**

```bash
go test ./internal/docgen -run TestProcessTemplate -v
```

Expected: FAIL

- [ ] **Step 7: Write process template function**

```go
// Add to internal/docgen/template.go

import (
	"os"
	"text/template"
)

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
```

- [ ] **Step 8: Run test to verify it passes**

```bash
go test ./internal/docgen -run TestProcessTemplate -v
```

Expected: PASS

- [ ] **Step 9: Run all docgen tests**

```bash
go test ./internal/docgen -v
```

Expected: All PASS

- [ ] **Step 10: Commit**

```bash
git add internal/docgen/
git commit -m "feat(docgen): add template processing with output function"
```

---

## Task 5: Create docgen CLI Tool

**Files:**
- Create: `cmd/docgen/main.go`

- [ ] **Step 1: Create cmd/docgen directory**

```bash
mkdir -p cmd/docgen
```

- [ ] **Step 2: Write docgen main.go with generate command**

```go
// cmd/docgen/main.go
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/s3cy/autoebiten/internal/docgen"
	"github.com/s3cy/autoebiten/testkit"
)

func main() {
	if len(os.Args) < 2 {
		processAllTemplates()
		return
	}

	switch os.Args[1] {
	case "--generate":
		if len(os.Args) < 3 {
			fatal(fmt.Errorf("usage: docgen --generate <example_dir>"))
		}
		generateExamples(os.Args[2])
	case "--process":
		processAllTemplates()
	default:
		fatal(fmt.Errorf("unknown command: %s", os.Args[1]))
	}
}

func generateExamples(exampleDir string) {
	// Load config
	configPath := filepath.Join(exampleDir, "config.yaml")
	config, err := docgen.LoadConfig(configPath)
	if err != nil {
		fatal(err)
	}

	// Build game binary
	gameBin := filepath.Join(config.GameDir, "autoui_demo")
	fmt.Printf("Building: %s\n", config.GameDir)
	buildCmd := exec.Command("go", "build", "-o", "autoui_demo", ".")
	buildCmd.Dir = config.GameDir
	if err := buildCmd.Run(); err != nil {
		fatal(fmt.Errorf("failed to build game: %w", err))
	}

	// Launch game using testkit
	fmt.Println("Launching game...")
	t := &testing.T{}
	game := testkit.Launch(t, gameBin, testkit.WithTimeout(30*time.Second))

	// Wait for game to be ready
	ready := game.WaitFor(func() bool {
		return game.Ping() == nil
	}, 5*time.Second)
	if !ready {
		game.Shutdown()
		fatal(fmt.Errorf("game failed to start"))
	}
	fmt.Println("Game ready")

	// Find and run all scripts
	scripts := findScripts(exampleDir)
	for _, script := range scripts {
		name := strings.TrimSuffix(filepath.Base(script), ".sh")
		outputFile := filepath.Join(exampleDir, name+"_out.txt")

		fmt.Printf("Running: %s\n", script)
		output := runScript(script)

		// Save output
		if err := os.WriteFile(outputFile, []byte(output), 0644); err != nil {
			fatal(fmt.Errorf("failed to write output: %w", err))
		}
		fmt.Printf("Generated: %s\n", outputFile)
	}

	// Shutdown game
	game.Shutdown()
	fmt.Println("Game stopped")
}

func processAllTemplates() {
	generateDir := "docs/generate"

	// Find all .tmpl files
	templates, err := filepath.Glob(filepath.Join(generateDir, "*.tmpl"))
	if err != nil {
		fatal(err)
	}

	for _, tmplPath := range templates {
		// Output file: strip .tmpl suffix, place in docs/
		base := filepath.Base(tmplPath)
		outName := strings.TrimSuffix(base, ".tmpl")
		outPath := filepath.Join("docs", outName)

		fmt.Printf("Processing: %s\n", tmplPath)

		result, err := docgen.ProcessTemplate(tmplPath, generateDir)
		if err != nil {
			fatal(err)
		}

		if err := os.WriteFile(outPath, []byte(result), 0644); err != nil {
			fatal(fmt.Errorf("failed to write %s: %w", outPath, err))
		}
		fmt.Printf("Generated: %s -> %s\n", tmplPath, outPath)
	}
}

func findScripts(dir string) []string {
	files, err := filepath.Glob(filepath.Join(dir, "*.sh"))
	if err != nil {
		fatal(err)
	}
	return files
}

func runScript(script string) string {
	cmd := exec.Command("bash", script)
	output, err := cmd.Output()
	if err != nil {
		// Include stderr in output for debugging
		if exitErr, ok := err.(*exec.ExitError); ok {
			return string(exitErr.Stderr)
		}
		return fmt.Sprintf("ERROR: %v", err)
	}
	return string(output)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
```

- [ ] **Step 3: Build docgen**

```bash
go build ./cmd/docgen
```

Expected: Success (no errors)

- [ ] **Step 4: Test docgen --process command**

```bash
./docgen --process
```

Expected: "Processing: ..." messages for any existing templates (or "no templates found" if none exist yet).

- [ ] **Step 5: Commit**

```bash
git add cmd/docgen/
git commit -m "feat(docgen): add CLI tool for example generation and template processing"
```

---

## Task 6: Create docverify CLI Tool

**Files:**
- Create: `cmd/docverify/main.go`

- [ ] **Step 1: Create cmd/docverify directory**

```bash
mkdir -p cmd/docverify
```

- [ ] **Step 2: Write docverify main.go**

```go
// cmd/docverify/main.go
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/s3cy/autoebiten/internal/docgen"
	"github.com/s3cy/autoebiten/testkit"
)

func main() {
	if len(os.Args) < 2 {
		fatal(fmt.Errorf("usage: docverify <example_dir>"))
	}
	verifyExamples(os.Args[1])
}

func verifyExamples(exampleDir string) {
	// Load config
	configPath := filepath.Join(exampleDir, "config.yaml")
	config, err := docgen.LoadConfig(configPath)
	if err != nil {
		fatal(err)
	}

	// Build game binary
	gameBin := filepath.Join(config.GameDir, "autoui_demo")
	fmt.Printf("Building: %s\n", config.GameDir)
	buildCmd := exec.Command("go", "build", "-o", "autoui_demo", ".")
	buildCmd.Dir = config.GameDir
	if err := buildCmd.Run(); err != nil {
		fatal(fmt.Errorf("failed to build game: %w", err))
	}

	// Launch game using testkit
	fmt.Println("Launching game...")
	t := &testing.T{}
	game := testkit.Launch(t, gameBin, testkit.WithTimeout(30*time.Second))

	// Wait for game to be ready
	ready := game.WaitFor(func() bool {
		return game.Ping() == nil
	}, 5*time.Second)
	if !ready {
		game.Shutdown()
		fatal(fmt.Errorf("game failed to start"))
	}
	fmt.Println("Game ready")

	// Verify each script
	hasMismatch := false
	scripts := findScripts(exampleDir)

	for _, script := range scripts {
		name := strings.TrimSuffix(filepath.Base(script), ".sh")
		outputFile := filepath.Join(exampleDir, name+"_out.txt")

		fmt.Printf("Verifying: %s\n", name)

		// Run script
		liveOutput := runScript(script)

		// Normalize live output
		liveNorm := docgen.Normalize(liveOutput, config.Normalize)

		// Load and normalize expected
		expected, err := os.ReadFile(outputFile)
		if err != nil {
			fmt.Printf("  MISSING: %s\n", outputFile)
			hasMismatch = true
			continue
		}
		expectedNorm := docgen.Normalize(string(expected), config.Normalize)

		// Compare
		if liveNorm == expectedNorm {
			fmt.Println("  OK: matches")
		} else {
			fmt.Println("  MISMATCH: output differs")
			printDiff(expectedNorm, liveNorm)
			hasMismatch = true
		}
	}

	// Cleanup
	game.Shutdown()
	fmt.Println("Game stopped")

	if hasMismatch {
		fmt.Println("\nFAILED: Some outputs do not match")
		fmt.Println("Run: make docs-generate to update outputs")
		os.Exit(1)
	}
	fmt.Println("\nSUCCESS: All outputs match")
}

func findScripts(dir string) []string {
	files, err := filepath.Glob(filepath.Join(dir, "*.sh"))
	if err != nil {
		fatal(err)
	}
	return files
}

func runScript(script string) string {
	cmd := exec.Command("bash", script)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return string(exitErr.Stderr)
		}
		return fmt.Sprintf("ERROR: %v", err)
	}
	return string(output)
}

func printDiff(expected, live string) {
	fmt.Println("=== Full diff (normalized) ===")

	expectedLines := strings.Split(expected, "\n")
	liveLines := strings.Split(live, "\n")

	maxLen := max(len(expectedLines), len(liveLines))
	for i := 0; i < maxLen; i++ {
		if i < len(expectedLines) && i < len(liveLines) {
			if expectedLines[i] != liveLines[i] {
				fmt.Printf("  @@ line %d @@\n", i+1)
				fmt.Printf("  -%s\n", expectedLines[i])
				fmt.Printf("  +%s\n", liveLines[i])
			}
		} else if i < len(expectedLines) {
			fmt.Printf("  -%s\n", expectedLines[i])
		} else {
			fmt.Printf("  +%s\n", liveLines[i])
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
```

- [ ] **Step 3: Build docverify**

```bash
go build ./cmd/docverify
```

Expected: Success (no errors)

- [ ] **Step 4: Commit**

```bash
git add cmd/docverify/
git commit -m "feat(docverify): add CLI tool for verifying example outputs"
```

---

## Task 7: Create autoui Example Config

**Files:**
- Create: `docs/generate/autoui_examples/config.yaml`

- [ ] **Step 1: Write config.yaml**

```yaml
# docs/generate/autoui_examples/config.yaml
game_dir: examples/autoui
normalize:
  # Normalize memory addresses (e.g., _addr="0x14000123456")
  - pattern: '_addr="[^"]*"'
    replace: '_addr="<ADDR>"'
  
  # Normalize timestamps in game logs
  - pattern: '\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+'
    replace: '<TIMESTAMP>'
  
  # Normalize PIDs
  - pattern: 'PID=\d+'
    replace: 'PID=<PID>'
```

- [ ] **Step 2: Commit**

```bash
git add docs/generate/autoui_examples/config.yaml
git commit -m "docs(docgen): add config for autoui examples"
```

---

## Task 8: Create autoui Example Scripts

**Files:**
- Create: `docs/generate/autoui_examples/autoui_tree.sh`
- Create: `docs/generate/autoui_examples/autoui_find_buttons.sh`
- Create: `docs/generate/autoui_examples/autoui_find_submit.sh`
- Create: `docs/generate/autoui_examples/autoui_call_click.sh`
- Create: `docs/generate/autoui_examples/autoui_exists_button.sh`

- [ ] **Step 1: Write autoui_tree.sh**

```bash
#!/bin/bash
# docs/generate/autoui_examples/autoui_tree.sh
# To run this example:
#   1. cd examples/autoui && go build -o autoui_demo
#   2. autoebiten launch -- ./autoui_demo &
#   3. Run this script

autoebiten custom autoui.tree
```

- [ ] **Step 2: Write autoui_find_buttons.sh**

```bash
#!/bin/bash
# docs/generate/autoui_examples/autoui_find_buttons.sh
# Find all Button widgets

autoebiten custom autoui.find --request "type=Button"
```

- [ ] **Step 3: Write autoui_find_submit.sh**

```bash
#!/bin/bash
# docs/generate/autoui_examples/autoui_find_submit.sh
# Find Submit button by id

autoebiten custom autoui.find --request "id=submit-btn"
```

- [ ] **Step 4: Write autoui_call_click.sh**

```bash
#!/bin/bash
# docs/generate/autoui_examples/autoui_call_click.sh
# Click the Submit button

autoebiten custom autoui.call --request '{"target":"id=submit-btn","method":"Click","args":[]}'
```

- [ ] **Step 5: Write autoui_exists_button.sh**

```bash
#!/bin/bash
# docs/generate/autoui_examples/autoui_exists_button.sh
# Check if Button widgets exist

autoebiten custom autoui.exists --request "type=Button"
```

- [ ] **Step 6: Make scripts executable**

```bash
chmod +x docs/generate/autoui_examples/*.sh
```

- [ ] **Step 7: Commit**

```bash
git add docs/generate/autoui_examples/*.sh
git commit -m "docs(docgen): add autoui example scripts"
```

---

## Task 9: Create autoui Template File

**Files:**
- Create: `docs/generate/autoui.md.tmpl`

- [ ] **Step 1: Write autoui.md.tmpl (simplified version)**

```markdown
# autoui Reference

> Purpose: EbitenUI widget automation for CLI, tests, and LLM agents
> Audience: CLI users, test writers, LLM agents automating EbitenUI games

---

## Quick Decision

**Querying widgets:**
```
├─ Need widget at coordinates? → autoui.at command
├─ Need widgets by attribute? → autoui.find command
├─ Need complex query? → autoui.xpath command
├─ Need to check existence? → autoui.exists command (returns JSON)
└─ Need full tree? → autoui.tree command
```

**Acting on widgets:**
```
├─ Need to click/interact? → autoui.call command
├─ Need to set text? → autoui.call SetText method
└─ Need visual debugging? → autoui.highlight command
```

---

## Overview

autoui provides EbitenUI automation via CLI commands:

- Widget tree inspection → XML export
- Widget search → coordinates, attributes, XPath
- Method invocation → reflection-based calls
- Visual debugging → highlight rectangles

**Key concepts:**

1. **WidgetInfo**: Internal representation (type, rect, state, customData)
2. **XML output**: Widget tree as XML (type as element name)
3. **ae tags**: Custom attributes via struct field tags

**Integration:**
```go
ui := ebitenui.UI{Container: root}
autoui.Register(&ui)  // Registers autoui.* commands
```

---

## Commands

### autoui.tree

Export full widget tree as XML.

**Usage:**
```bash
autoebiten custom autoui.tree
```

**Output:**
{{output "autoui_examples/autoui_tree_out.txt" "xml"}}

---

### autoui.find

Find widgets by attribute (AND logic for multiple criteria).

**Usage:**
```bash
autoebiten custom autoui.find --request "type=Button"
```

**Output:**
{{output "autoui_examples/autoui_find_buttons_out.txt" "xml"}}

---

### autoui.exists

Check if widgets matching a query exist. Returns JSON for use with `wait-for`.

**Usage:**
```bash
autoebiten custom autoui.exists --request "type=Button"
```

**Output:**
{{output "autoui_examples/autoui_exists_button_out.txt" "json"}}

---

### autoui.call

Invoke method on widget.

**Usage:**
```bash
autoebiten custom autoui.call --request '{"target":"id=submit-btn","method":"Click","args":[]}'
```

**Output:**
{{output "autoui_examples/autoui_call_click_out.txt" "json"}}
```

- [ ] **Step 2: Commit**

```bash
git add docs/generate/autoui.md.tmpl
git commit -m "docs(docgen): add autoui template file"
```

---

## Task 10: Create Makefile

**Files:**
- Create: `Makefile`

- [ ] **Step 1: Write Makefile**

```makefile
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
	go run ./cmd/docgen --process

# Verify outputs match recorded files
docs-verify:
	go run ./cmd/docverify docs/generate/autoui_examples

# Clean generated files
docs-clean:
	rm -f docs/generate/autoui_examples/*_out.txt
	rm -f docs/autoui.md
```

- [ ] **Step 2: Test Makefile**

```bash
make build
```

Expected: Success (no errors)

- [ ] **Step 3: Commit**

```bash
git add Makefile
git commit -m "feat: add Makefile with docs generation targets"
```

---

## Task 11: Generate and Verify autoui Examples

**Files:**
- Generate: `docs/generate/autoui_examples/*_out.txt`
- Generate: `docs/autoui.md`

- [ ] **Step 1: Run docs-generate**

```bash
make docs-generate
```

Expected: Builds game, runs scripts, saves outputs, processes template.

- [ ] **Step 2: Check generated outputs exist**

```bash
ls -la docs/generate/autoui_examples/*_out.txt
```

Expected: Shows `autoui_tree_out.txt`, `autoui_find_buttons_out.txt`, etc.

- [ ] **Step 3: Check generated docs**

```bash
ls -la docs/autoui.md
```

Expected: File exists.

- [ ] **Step 4: Read generated autoui.md to verify template processing**

```bash
head -50 docs/autoui.md
```

Expected: Shows markdown with actual outputs embedded (not `{{output}}` placeholders).

- [ ] **Step 5: Run docs-verify**

```bash
make docs-verify
```

Expected: "SUCCESS: All outputs match"

- [ ] **Step 6: Commit generated files**

```bash
git add docs/generate/autoui_examples/*_out.txt docs/autoui.md
git commit -m "docs: generate autoui examples and documentation"
```

---

## Task 12: Final Verification

- [ ] **Step 1: Run full test suite**

```bash
go test -race ./...
```

Expected: All tests pass.

- [ ] **Step 2: Build all binaries**

```bash
make build
```

Expected: Success.

- [ ] **Step 3: Verify docs-verify passes**

```bash
make docs-verify
```

Expected: "SUCCESS: All outputs match"

- [ ] **Step 4: Clean and regenerate to test full workflow**

```bash
make docs-clean
make docs-generate
make docs-verify
```

Expected: All steps succeed.

- [ ] **Step 5: Final commit**

```bash
git add -A
git status
```

Expected: Clean working tree (all changes committed).

---

## Summary

This plan creates:

1. **internal/docgen** package with config loading, normalization, and template processing
2. **cmd/docgen** CLI tool that runs scripts and generates docs
3. **cmd/docverify** CLI tool that verifies outputs match
4. **docs/generate/** structure with examples, config, and templates
5. **Makefile** targets for docs-generate, docs-verify, docs-clean

The system allows:
- Users to run example scripts directly for learning
- Developers to generate outputs and docs with `make docs-generate`
- CI/LLM to verify outputs with `make docs-verify`
- Normalization of dynamic values for stable comparison