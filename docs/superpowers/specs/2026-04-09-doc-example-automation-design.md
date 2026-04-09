---
name: doc-example-automation
description: Automate documentation examples with real runs and output capture
type: project
---

# Doc Example Automation Design

**Goal:** Ensure documentation examples match actual implementation by running real commands and capturing outputs. Scripts are runnable by users for learning; outputs are captured for doc generation.

---

## Architecture

### File Structure

```
docs/
├── autoui.md                  # Generated documentation (top level)
├── commands.md
├── tutorial.md
└── generate/                  # All generation materials
    ├── autoui.md.tmpl         # Template file
    ├── autoui_examples/       # Autoui example artifacts
    │   ├── config.yaml        # Game dir + normalization rules
    │   ├── autoui_tree.sh     # Minimal runnable script
    │   ├── autoui_tree_out.txt
    │   ├── autoui_find_buttons.sh
    │   ├── autoui_find_buttons_out.txt
    │   └── ...
    ├── commands.md.tmpl
    ├── commands_examples/
    │   ├── config.yaml
    │   ├── input_press.sh
    │   ├── input_press_out.txt
    │   └── ...
    └── tutorial.md.tmpl
        └── tutorial_examples/
            └── ...

cmd/
├── docgen/main.go             # Generator: runs scripts, processes templates
├── docverify/main.go          # Verifier: compares outputs
Makefile
```

---

## Components

### 1. Example Scripts

**Location:** `docs/generate/<name>_examples/*.sh`

**Purpose:** Minimal commands users can run to reproduce examples.

**Format:**
```bash
#!/bin/bash
# To run this example:
#   1. cd examples/autoui && go build -o autoui_demo
#   2. autoebiten launch -- ./autoui_demo &
#   3. Run this script

autoebiten custom autoui.tree
```

**Principles:**
- Contain only the essential command(s)
- No setup/cleanup logic (handled by orchestrator)
- Users can read and run them directly
- Comments explain prerequisites

---

### 2. Output Files

**Location:** `docs/generate/<name>_examples/*_out.txt`

**Purpose:** Captured outputs from running scripts.

**Format:** Raw output (XML, JSON, text, etc.) - no metadata.

**Example:** `autoui_tree_out.txt`
```
<UI>
  <Container x="0" y="0" width="640" height="480">
    <Button x="100" y="50" id="submit-btn" disabled="false"/>
  </Container>
</UI>
```

---

### 3. Config Files

**Location:** `docs/generate/<name>_examples/config.yaml`

**Purpose:** Define game directory and normalization rules for dynamic content.

**Format:**
```yaml
game_dir: examples/autoui
normalize:
  - pattern: '_addr="[^"]*"'
    replace: '_addr="<ADDR>"'
  - pattern: '\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+'
    replace: '<TIMESTAMP>'
  - pattern: 'PID=\d+'
    replace: 'PID=<PID>'
```

**Why:** Some outputs have dynamic values (memory addresses, timestamps, PIDs) that change each run. Normalization replaces them with stable placeholders for comparison.

---

### 4. Template Files

**Location:** `docs/generate/*.md.tmpl`

**Purpose:** Markdown with `{{output}}` placeholders.

**Syntax:**
```markdown
**Output:**
{{output "autoui_examples/autoui_tree_out.txt" "xml"}}
```

- First arg: path relative to `docs/generate/`
- Second arg: code fence language (`xml`, `json`, `bash`, `text`)

---

### 5. docgen (Generator)

**Location:** `cmd/docgen/main.go`

**Purpose:** Run scripts, capture outputs, process templates.

**Workflow:**
1. Load config from example directory
2. Build game binary
3. Launch game via testkit
4. Run each script, capture output to `*_out.txt`
5. Shutdown game
6. Process all `.tmpl` files, generate final `.md` files

**Key implementation:**
- Use `testkit.Launch()` with minimal `*testing.T` for game lifecycle
- Use `testkit.WaitFor()` to wait for game ready (no fixed sleep)
- Custom template function `{{output}}` reads file and wraps in code fence

---

### 6. docverify (Verifier)

**Location:** `cmd/docverify/main.go`

**Purpose:** Re-run scripts, compare outputs to recorded files.

**Workflow:**
1. Load config
2. Build and launch game via testkit
3. Run each script
4. Normalize live output using config rules
5. Normalize expected output from `*_out.txt`
6. Compare normalized outputs
7. Report mismatches with full diff
8. Shutdown game

**Key implementation:**
- Apply regex normalization before comparison
- Print full normalized diff (not preview) for LLM to analyze
- Exit with error if any mismatch

---

## Makefile Targets

```makefile
.PHONY: docs-generate docs-verify docs-clean

docs-generate:
	go run ./cmd/docgen --generate docs/generate/autoui_examples
	go run ./cmd/docgen --generate docs/generate/commands_examples
	go run ./cmd/docgen --process

docs-verify:
	go run ./cmd/docverify docs/generate/autoui_examples
	go run ./cmd/docverify docs/generate/commands_examples

docs-clean:
	rm -f docs/generate/*_examples/*_out.txt
	rm -f docs/*.md
```

---

## Workflows

### Developer: Adding New Example

1. Write script: `docs/generate/autoui_examples/new_cmd.sh`
2. Add placeholder to template: `{{output "autoui_examples/new_cmd_out.txt" "xml"}}`
3. Run: `make docs-generate`
4. Commit: outputs and generated docs

### Developer: Checking for Drift

1. Run: `make docs-verify`
2. If mismatch reported:
   - Intentional (implementation changed): `make docs-generate`, commit
   - Bug: fix implementation, re-verify

### LLM: Same workflows work
- Full normalized diff output enables LLM to judge mismatches
- Config and scripts are readable for context
- No ambiguous partial output

### User: Learning from Examples

1. Build game: `cd examples/autoui && go build -o autoui_demo`
2. Launch: `autoebiten launch -- ./autoui_demo &`
3. Run script: `bash ../docs/generate/autoui_examples/autoui_tree.sh`
4. See real output in terminal
5. Cleanup: `autoebiten exit`

---

## Implementation Notes

### testkit Integration

Use `testkit.Launch()` for reliable game lifecycle:

```go
t := &testing.T{}
game := testkit.Launch(t, gameBin, testkit.WithTimeout(30*time.Second))

ready := game.WaitFor(func() bool {
    return game.Ping() == nil
}, 5*time.Second)

// Run scripts via autoebiten CLI (connects to running game)
// ...

game.Shutdown()
```

### Template Function

```go
func outputFn(baseDir string) func(string, string) (string, error) {
    return func(path, lang string) (string, error) {
        content, err := os.ReadFile(filepath.Join(baseDir, path))
        if err != nil {
            return "", err
        }
        return fmt.Sprintf("```%s\n%s\n```", lang, strings.TrimSpace(string(content))), nil
    }
}
```

### Normalization

```go
func normalize(s string, rules []NormalizeRule) string {
    for _, r := range rules {
        re := regexp.MustCompile(r.Pattern)
        s = re.ReplaceAllString(s, r.Replace)
    }
    return strings.TrimSpace(s)
}
```

---

## Why This Approach

1. **Scripts as tutorials** — Users can run them directly, see real behavior
2. **Outputs as contracts** — Recorded outputs verify implementation matches docs
3. **LLM-friendly** — Full diffs, readable configs, no ambiguity
4. **Minimal scripts** — Setup/cleanup in orchestrator, scripts stay clean
5. **testkit reliability** — Proper game lifecycle, wait-for instead of sleep
6. **Separate outputs** — Easy to version, review, update