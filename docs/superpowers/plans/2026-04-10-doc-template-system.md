# Documentation Template System Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Rewrite documentation generation to use Go templates with inline config, direct command execution, and Go code extraction.

**Architecture:** Go `text/template` with custom FuncMap functions. Templates store config inline, launch games via testkit, execute commands directly, and capture normalized output. No shell scripts or config.yaml files.

**Tech Stack:** Go stdlib (`text/template`, `go/parser`, `go/ast`, `go/format`), existing testkit package

---

## File Structure

```
internal/docgen/
├── cmd/
│   └── generate.go      # Entry point - processes templates
├── context.go           # Template execution context (state management)
├── config.go            # Config parsing from template
├── normalize.go         # Regex-based output normalization
├── template.go          # FuncMap and template processing
├── game.go              # launch_game, command, end_game, delay
├── verify.go            # verifyOutputs function
├── gocode.go            # Go code extraction with AST transforms

docs/generate/
├── commands.md.gotmpl   # Rewritten from commands.md.tmpl
├── testkit.md.gotmpl    # Rewritten from testkit.md.tmpl
├── tutorial.md.gotmpl   # Rewritten from tutorial.md.tmpl
├── integration.md.gotmpl
├── autoui.md.gotmpl

Makefile                 # Add docs: target
```

---

## Task 1: Template Context

**Files:**
- Create: `internal/docgen/context.go`
- Create: `internal/docgen/context_test.go`

**Why:** Go templates need shared state between function calls. The context holds the current game session, config, and outputs.

- [ ] **Step 1: Write the failing test**

```go
// internal/docgen/context_test.go
package docgen

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestNewContext(t *testing.T) {
    ctx := NewContext()
    assert.Nil(t, ctx.GameSession)
    assert.Nil(t, ctx.Config)
}

func TestContextSetConfig(t *testing.T) {
    ctx := NewContext()
    cfg := &Config{GameDir: "examples/test"}
    ctx.SetConfig(cfg)
    assert.Equal(t, "examples/test", ctx.Config.GameDir)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/docgen/... -run TestContext -v`
Expected: FAIL with "undefined: Context" or similar

- [ ] **Step 3: Write minimal implementation**

```go
// internal/docgen/context.go
package docgen

// Context holds template execution state.
type Context struct {
    GameSession *GameSession
    Config      *Config
    outputs     []string
}

// NewContext creates a new template context.
func NewContext() *Context {
    return &Context{}
}

// SetConfig sets the configuration for the template.
func (c *Context) SetConfig(cfg *Config) {
    c.Config = cfg
}

// AddOutput stores a captured output.
func (c *Context) AddOutput(output string) {
    c.outputs = append(c.outputs, output)
}

// GetOutputs returns all captured outputs.
func (c *Context) GetOutputs() []string {
    return c.outputs
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/docgen/... -run TestContext -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/docgen/context.go internal/docgen/context_test.go
git commit -m "feat(docgen): add template context for state management"
```

---

## Task 2: Config Parsing

**Files:**
- Create: `internal/docgen/config.go`
- Create: `internal/docgen/config_test.go`
- Modify: `internal/docgen/config.go` (replace existing)

**Why:** Parse inline config from template, replacing config.yaml files.

- [ ] **Step 1: Write the failing test**

```go
// internal/docgen/config_test.go
package docgen

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestConfigDefaults(t *testing.T) {
    cfg := &Config{}
    assert.Equal(t, "", cfg.GameDir)
    assert.Empty(t, cfg.Normalize)
}

func TestNormalizeRule(t *testing.T) {
    rules := []NormalizeRule{
        {Pattern: `PID=\d+`, Replace: "PID=<PID>"},
        {Pattern: `\d{4}-\d{2}-\d{2}`, Replace: "<TIMESTAMP>"},
    }
    cfg := &Config{Normalize: rules}
    assert.Len(t, cfg.Normalize, 2)
    assert.Equal(t, `PID=\d+`, cfg.Normalize[0].Pattern)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/docgen/... -run TestConfig -v`
Expected: FAIL (types don't match existing)

- [ ] **Step 3: Replace existing config.go**

```go
// internal/docgen/config.go
package docgen

// NormalizeRule defines a regex pattern and replacement.
type NormalizeRule struct {
    Pattern string
    Replace string
}

// Config holds template configuration.
type Config struct {
    GameDir   string
    Normalize []NormalizeRule
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/docgen/... -run TestConfig -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/docgen/config.go internal/docgen/config_test.go
git commit -m "feat(docgen): simplify config struct for inline use"
```

---

## Task 3: Normalization

**Files:**
- Create: `internal/docgen/normalize.go`
- Create: `internal/docgen/normalize_test.go`
- Modify: `internal/docgen/config.go` (move Normalize function)

**Why:** Apply regex rules to normalize command output (replace PIDs, timestamps, paths).

- [ ] **Step 1: Write the failing test**

```go
// internal/docgen/normalize_test.go
package docgen

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
    rules := []NormalizeRule{
        {Pattern: `PID=\d+`, Replace: "PID=<PID>"},
        {Pattern: `\d{4}-\d{2}-\d{2}`, Replace: "<TIMESTAMP>"},
    }

    input := "PID=12345 started at 2026-04-10"
    expected := "PID=<PID> started at <TIMESTAMP>"

    result := Normalize(input, rules)
    assert.Equal(t, expected, result)
}

func TestNormalizeEmptyRules(t *testing.T) {
    input := "unchanged output"
    result := Normalize(input, nil)
    assert.Equal(t, input, result)
}

func TestNormalizeTrimsWhitespace(t *testing.T) {
    input := "  output with spaces  \n"
    result := Normalize(input, nil)
    assert.Equal(t, "output with spaces", result)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/docgen/... -run TestNormalize -v`
Expected: FAIL with "undefined: Normalize"

- [ ] **Step 3: Write minimal implementation**

```go
// internal/docgen/normalize.go
package docgen

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

Run: `go test ./internal/docgen/... -run TestNormalize -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/docgen/normalize.go internal/docgen/normalize_test.go
git commit -m "feat(docgen): add regex-based output normalization"
```

---

## Task 4: Game Session

**Files:**
- Create: `internal/docgen/game.go`
- Create: `internal/docgen/game_test.go`

**Why:** Manage game lifecycle - build, launch, execute commands, shutdown. Wraps testkit with docgen-specific functionality.

- [ ] **Step 1: Write the failing test**

```go
// internal/docgen/game_test.go
package docgen

import (
    "os"
    "path/filepath"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestBuildGame(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    // Build the simple example
    gameDir := "examples/simple"
    binaryPath, err := buildGame(gameDir)
    require.NoError(t, err)
    defer os.Remove(binaryPath)
    
    // Verify binary exists
    _, err = os.Stat(binaryPath)
    assert.NoError(t, err)
}

func TestGameSessionLifecycle(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    ctx := NewContext()
    ctx.SetConfig(&Config{GameDir: "examples/simple"})
    
    // Launch
    session, err := LaunchGame(ctx)
    require.NoError(t, err)
    require.NotNil(t, session)
    
    // Verify game is running
    assert.NotNil(t, session.game)
    
    // Cleanup
    err = EndGame(session)
    assert.NoError(t, err)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/docgen/... -run TestGame -v`
Expected: FAIL with undefined functions

- [ ] **Step 3: Write minimal implementation**

```go
// internal/docgen/game.go
package docgen

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "time"
    
    "github.com/s3cy/autoebiten/testkit"
)

// GameSession wraps testkit.Game with docgen-specific state.
type GameSession struct {
    game       *testkit.Game
    t          *testing.T
    socketPath string
    ctx        *Context
}

// buildGame compiles a game binary and returns its path.
func buildGame(gameDir string) (string, error) {
    binaryName := filepath.Base(gameDir) + "_docgen"
    binaryPath := filepath.Join(gameDir, binaryName)
    
    cmd := exec.Command("go", "build", "-o", binaryName, ".")
    cmd.Dir = gameDir
    if output, err := cmd.CombinedOutput(); err != nil {
        return "", fmt.Errorf("failed to build game: %w\n%s", err, output)
    }
    
    return binaryPath, nil
}

// LaunchGame builds and starts a game, returning a session.
func LaunchGame(ctx *Context, args ...string) (*GameSession, error) {
    if ctx.Config == nil || ctx.Config.GameDir == "" {
        return nil, fmt.Errorf("config.GameDir not set")
    }
    
    // Build binary
    binaryPath, err := buildGame(ctx.Config.GameDir)
    if err != nil {
        return nil, err
    }
    
    // Create test context (testkit requires *testing.T)
    t := &testing.T{}
    
    // Build options
    opts := []testkit.Option{testkit.WithTimeout(30 * time.Second)}
    if len(args) > 0 {
        opts = append(opts, testkit.WithArgs(args...))
    }
    
    // Launch via testkit
    game := testkit.Launch(t, binaryPath, opts...)
    
    // Wait for ready
    ready := game.WaitFor(func() bool {
        return game.Ping() == nil
    }, 5*time.Second)
    
    if !ready {
        game.Shutdown()
        os.Remove(binaryPath)
        return nil, fmt.Errorf("game failed to start")
    }
    
    return &GameSession{
        game:       game,
        t:          t,
        socketPath: game.SocketPath(),
        ctx:        ctx,
    }, nil
}

// EndGame shuts down the game and cleans up.
func EndGame(session *GameSession) error {
    if session == nil || session.game == nil {
        return nil
    }
    
    session.game.Shutdown()
    
    // Remove built binary
    if session.ctx.Config != nil {
        binaryName := filepath.Base(session.ctx.Config.GameDir) + "_docgen"
        binaryPath := filepath.Join(session.ctx.Config.GameDir, binaryName)
        os.Remove(binaryPath)
    }
    
    return nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/docgen/... -run TestGame -v`
Expected: PASS (may take a few seconds to build/run game)

- [ ] **Step 5: Commit**

```bash
git add internal/docgen/game.go internal/docgen/game_test.go
git commit -m "feat(docgen): add game session management with testkit"
```

---

## Task 5: Command Execution

**Files:**
- Modify: `internal/docgen/game.go`
- Modify: `internal/docgen/game_test.go`

**Why:** Execute autoebiten commands and capture normalized output.

- [ ] **Step 1: Write the failing test**

Add to `internal/docgen/game_test.go`:

```go
func TestExecuteCommand(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    ctx := NewContext()
    ctx.SetConfig(&Config{
        GameDir: "examples/simple",
        Normalize: []NormalizeRule{
            {Pattern: `PID=\d+`, Replace: "PID=<PID>"},
        },
    })
    
    session, err := LaunchGame(ctx)
    require.NoError(t, err)
    defer EndGame(session)
    
    // Execute ping command
    output, err := ExecuteCommand(session, "ping", nil)
    require.NoError(t, err)
    
    // Should contain "OK" or "running"
    assert.Contains(t, output, "OK")
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/docgen/... -run TestExecuteCommand -v`
Expected: FAIL with "undefined: ExecuteCommand"

- [ ] **Step 3: Write minimal implementation**

Add to `internal/docgen/game.go`:

```go
import (
    "bytes"
    "os/exec"
    "strings"
)

// ExecuteCommand runs an autoebiten command and returns normalized output.
func ExecuteCommand(session *GameSession, cmdName string, flags map[string]any) (string, error) {
    if session == nil || session.game == nil {
        return "", fmt.Errorf("no active game session")
    }
    
    // Build command arguments
    args := []string{cmdName}
    for k, v := range flags {
        args = append(args, fmt.Sprintf("--%s", k), fmt.Sprintf("%v", v))
    }
    
    // Execute autoebiten CLI
    cmd := exec.Command("autoebiten", args...)
    cmd.Env = append(os.Environ(), "AUTOEBITEN_SOCKET="+session.socketPath)
    
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    
    // Run command (don't fail on non-zero exit - capture output)
    cmd.Run()
    
    // Combine stdout and stderr
    output := stdout.String() + stderr.String()
    
    // Apply normalization
    if session.ctx != nil && session.ctx.Config != nil {
        output = Normalize(output, session.ctx.Config.Normalize)
    }
    
    // Store in context
    if session.ctx != nil {
        session.ctx.AddOutput(output)
    }
    
    return output, nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/docgen/... -run TestExecuteCommand -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/docgen/game.go internal/docgen/game_test.go
git commit -m "feat(docgen): add command execution with normalization"
```

---

## Task 6: Delay Function

**Files:**
- Modify: `internal/docgen/game.go`
- Modify: `internal/docgen/game_test.go`

**Why:** Wait for crash scenarios to complete (game crashes after N ticks).

- [ ] **Step 1: Write the failing test**

Add to `internal/docgen/game_test.go`:

```go
func TestDelay(t *testing.T) {
    start := time.Now()
    Delay("100ms")
    elapsed := time.Since(start)
    
    assert.GreaterOrEqual(t, elapsed.Milliseconds(), int64(100))
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/docgen/... -run TestDelay -v`
Expected: FAIL with "undefined: Delay"

- [ ] **Step 3: Write minimal implementation**

Add to `internal/docgen/game.go`:

```go
import "time"

// Delay pauses execution for the specified duration.
// Used for crash scenarios where game crashes after N seconds.
func Delay(duration string) {
    d, err := time.ParseDuration(duration)
    if err != nil {
        d = 1 * time.Second // default fallback
    }
    time.Sleep(d)
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/docgen/... -run TestDelay -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/docgen/game.go internal/docgen/game_test.go
git commit -m "feat(docgen): add delay function for crash timing"
```

---

## Task 7: Verification

**Files:**
- Create: `internal/docgen/verify.go`
- Create: `internal/docgen/verify_test.go`

**Why:** Compare multiple command outputs to ensure variations produce identical results.

- [ ] **Step 1: Write the failing test**

```go
// internal/docgen/verify_test.go
package docgen

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestVerifyOutputs(t *testing.T) {
    // All identical - should pass
    err := VerifyOutputs("OK: running", "OK: running", "OK: running")
    assert.NoError(t, err)
}

func TestVerifyOutputsMismatch(t *testing.T) {
    // Different outputs - should fail
    err := VerifyOutputs("OK: running", "Error: failed", "OK: running")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "output mismatch")
}

func TestVerifyOutputsEmpty(t *testing.T) {
    // No outputs - should pass
    err := VerifyOutputs()
    assert.NoError(t, err)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/docgen/... -run TestVerify -v`
Expected: FAIL with "undefined: VerifyOutputs"

- [ ] **Step 3: Write minimal implementation**

```go
// internal/docgen/verify.go
package docgen

import (
    "fmt"
)

// VerifyOutputs checks that all provided outputs are identical.
// Returns an error if any outputs differ.
func VerifyOutputs(outputs ...string) error {
    if len(outputs) == 0 {
        return nil
    }
    
    first := outputs[0]
    for i, output := range outputs {
        if output != first {
            return fmt.Errorf("output mismatch: output[0]=%q != output[%d]=%q", first, i, output)
        }
    }
    
    return nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/docgen/... -run TestVerify -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/docgen/verify.go internal/docgen/verify_test.go
git commit -m "feat(docgen): add output verification for command variations"
```

---

## Task 8: Go Code Extraction

**Files:**
- Create: `internal/docgen/gocode.go`
- Create: `internal/docgen/gocode_test.go`

**Why:** Extract functions/structs from Go source files with AST-based transforms.

- [ ] **Step 1: Write the failing test**

```go
// internal/docgen/gocode_test.go
package docgen

import (
    "os"
    "path/filepath"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestExtractFunction(t *testing.T) {
    // Create a temp Go file
    tmpDir := t.TempDir()
    tmpFile := filepath.Join(tmpDir, "test.go")
    
    content := `package test

import "fmt"

func HelloWorld() {
    fmt.Println("hello")
}
`
    err := os.WriteFile(tmpFile, []byte(content), 0644)
    require.NoError(t, err)
    
    // Extract function
    result, err := ExtractGoCode(tmpFile, "HelloWorld", nil)
    require.NoError(t, err)
    
    assert.Contains(t, result, "func HelloWorld()")
    assert.Contains(t, result, `fmt.Println("hello")`)
}

func TestExtractFunctionStripImports(t *testing.T) {
    tmpDir := t.TempDir()
    tmpFile := filepath.Join(tmpDir, "test.go")
    
    content := `package test

import "fmt"

func HelloWorld() {
    fmt.Println("hello")
}
`
    err := os.WriteFile(tmpFile, []byte(content), 0644)
    require.NoError(t, err)
    
    result, err := ExtractGoCode(tmpFile, "HelloWorld", []string{"stripImports"})
    require.NoError(t, err)
    
    assert.Contains(t, result, "func HelloWorld()")
    assert.NotContains(t, result, `import "fmt"`)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/docgen/... -run TestExtract -v`
Expected: FAIL with "undefined: ExtractGoCode"

- [ ] **Step 3: Write minimal implementation**

```go
// internal/docgen/gocode.go
package docgen

import (
    "bytes"
    "fmt"
    "go/ast"
    "go/format"
    "go/parser"
    "go/token"
    "strings"
)

// ExtractGoCode extracts a function or struct from a Go file.
// Supported transforms: stripImports, rename:Old→New, package:Name
func ExtractGoCode(filePath, target string, transforms []string) (string, error) {
    fset := token.NewFileSet()
    astFile, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
    if err != nil {
        return "", fmt.Errorf("failed to parse file: %w", err)
    }
    
    // Find target declaration
    var targetNode ast.Node
    for _, decl := range astFile.Decls {
        switch d := decl.(type) {
        case *ast.FuncDecl:
            if d.Name.Name == target {
                targetNode = d
            }
        case *ast.GenDecl:
            for _, spec := range d.Specs {
                if ts, ok := spec.(*ast.TypeSpec); ok && ts.Name.Name == target {
                    targetNode = ts
                }
            }
        }
    }
    
    if targetNode == nil {
        return "", fmt.Errorf("target %q not found in %s", target, filePath)
    }
    
    // Apply transforms to the whole file (for package/import changes)
    for _, transform := range transforms {
        applyTransform(astFile, transform)
    }
    
    // Format the target node
    var buf bytes.Buffer
    if err := format.Node(&buf, fset, targetNode); err != nil {
        return "", fmt.Errorf("failed to format: %w", err)
    }
    
    result := buf.String()
    
    // Apply rename transforms to the result string
    for _, transform := range transforms {
        if strings.HasPrefix(transform, "rename:") {
            parts := strings.Split(strings.TrimPrefix(transform, "rename:"), "→")
            if len(parts) == 2 {
                result = strings.ReplaceAll(result, parts[0], parts[1])
            }
        }
    }
    
    return result, nil
}

func applyTransform(file *ast.File, transform string) {
    switch {
    case transform == "stripImports":
        // Remove import declarations
        var decls []ast.Decl
        for _, d := range file.Decls {
            if gd, ok := d.(*ast.GenDecl); ok && gd.Tok == token.IMPORT {
                continue
            }
            decls = append(decls, d)
        }
        file.Decls = decls
        
    case strings.HasPrefix(transform, "package:"):
        name := strings.TrimPrefix(transform, "package:")
        file.Name.Name = name
    }
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/docgen/... -run TestExtract -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/docgen/gocode.go internal/docgen/gocode_test.go
git commit -m "feat(docgen): add Go code extraction with AST transforms"
```

---

## Task 9: Template FuncMap

**Files:**
- Create: `internal/docgen/template.go`
- Create: `internal/docgen/template_test.go`
- Remove: `internal/docgen/template.go` (old version)

**Why:** Wire up all functions into a Go template FuncMap. This is the core integration point.

- [ ] **Step 1: Write the failing test**

```go
// internal/docgen/template_test.go
package docgen

import (
    "strings"
    "testing"
    "text/template"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestFuncMapHasAllFunctions(t *testing.T) {
    fm := FuncMap()
    
    // Core functions
    assert.Contains(t, fm, "config")
    assert.Contains(t, fm, "launch_game")
    assert.Contains(t, fm, "end_game")
    assert.Contains(t, fm, "command")
    assert.Contains(t, fm, "delay")
    assert.Contains(t, fm, "verifyOutputs")
    assert.Contains(t, fm, "gocode")
    
    // Utility functions
    assert.Contains(t, fm, "list")
    assert.Contains(t, fm, "dict")
}

func TestProcessSimpleTemplate(t *testing.T) {
    tmplContent := `Hello {{ .Name }}!`
    
    result, err := ProcessTemplateString(tmplContent, map[string]string{"Name": "World"})
    require.NoError(t, err)
    assert.Equal(t, "Hello World!", result)
}

func TestUtilityFunctions(t *testing.T) {
    fm := FuncMap()
    
    // Test list
    listFunc := fm["list"].(func(...any) []any)
    list := listFunc("a", "b", "c")
    assert.Equal(t, []any{"a", "b", "c"}, list)
    
    // Test dict
    dictFunc := fm["dict"].(func(...any) map[string]any)
    d := dictFunc("key1", "val1", "key2", "val2")
    assert.Equal(t, map[string]any{"key1": "val1", "key2": "val2"}, d)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/docgen/... -run TestFuncMap -v`
Expected: FAIL with type mismatches or missing functions

- [ ] **Step 3: Write minimal implementation**

```go
// internal/docgen/template.go
package docgen

import (
    "bytes"
    "fmt"
    "os"
    "path/filepath"
    "text/template"
)

// globalContext is the shared context for template execution.
var globalContext = NewContext()

// FuncMap returns the template function map.
func FuncMap() template.FuncMap {
    return template.FuncMap{
        // Config
        "config": configFunc,
        
        // Game control
        "launch_game": launchGameFunc,
        "end_game":    endGameFunc,
        "command":     commandFunc,
        "delay":       delayFunc,
        
        // Verification
        "verifyOutputs": verifyOutputsFunc,
        
        // Code extraction
        "gocode": gocodeFunc,
        
        // Utilities
        "list": listFunc,
        "dict": dictFunc,
    }
}

// configFunc sets the template configuration.
func configFunc(gameDir string, normalize ...map[string]any) (string, error) {
    cfg := &Config{GameDir: gameDir}
    
    // Parse normalize rules if provided
    for _, n := range normalize {
        if rules, ok := n["rules"].([]any); ok {
            for _, r := range rules {
                if rule, ok := r.(map[string]any); ok {
                    pattern, _ := rule["pattern"].(string)
                    replace, _ := rule["replace"].(string)
                    cfg.Normalize = append(cfg.Normalize, NormalizeRule{
                        Pattern: pattern,
                        Replace: replace,
                    })
                }
            }
        }
    }
    
    globalContext.SetConfig(cfg)
    return "", nil
}

// launchGameFunc launches a game session.
func launchGameFunc(args ...string) (string, error) {
    session, err := LaunchGame(globalContext, args...)
    if err != nil {
        return "", err
    }
    globalContext.GameSession = session
    return "", nil
}

// endGameFunc shuts down the current game session.
func endGameFunc() (string, error) {
    if globalContext.GameSession == nil {
        return "", nil
    }
    err := EndGame(globalContext.GameSession)
    globalContext.GameSession = nil
    return "", err
}

// commandFunc executes a command and returns the output.
func commandFunc(cmdName string, flags ...map[string]any) (string, error) {
    if globalContext.GameSession == nil {
        return "", fmt.Errorf("no active game session - call launch_game first")
    }
    
    var flagMap map[string]any
    if len(flags) > 0 {
        flagMap = flags[0]
    }
    
    return ExecuteCommand(globalContext.GameSession, cmdName, flagMap)
}

// delayFunc waits for the specified duration.
func delayFunc(duration string) (string, error) {
    Delay(duration)
    return "", nil
}

// verifyOutputsFunc verifies all outputs are identical.
func verifyOutputsFunc(outputs ...string) (string, error) {
    return "", VerifyOutputs(outputs...)
}

// gocodeFunc extracts Go code from a file.
func gocodeFunc(filePath, target string, transforms ...[]string) (string, error) {
    var t []string
    if len(transforms) > 0 {
        t = transforms[0]
    }
    return ExtractGoCode(filePath, target, t)
}

// listFunc creates a slice.
func listFunc(items ...any) []any {
    return items
}

// dictFunc creates a map.
func dictFunc(items ...any) map[string]any {
    d := make(map[string]any)
    for i := 0; i < len(items); i += 2 {
        if i+1 < len(items) {
            key, ok := items[i].(string)
            if ok {
                d[key] = items[i+1]
            }
        }
    }
    return d
}

// ProcessTemplate processes a template file and returns the result.
func ProcessTemplate(tmplPath string) (string, error) {
    content, err := os.ReadFile(tmplPath)
    if err != nil {
        return "", fmt.Errorf("failed to read template: %w", err)
    }
    
    // Reset global context for each template
    globalContext = NewContext()
    
    return ProcessTemplateString(string(content), nil)
}

// ProcessTemplateString processes a template string.
func ProcessTemplateString(tmplContent string, data any) (string, error) {
    tmpl, err := template.New("doc").Funcs(FuncMap()).Parse(tmplContent)
    if err != nil {
        return "", fmt.Errorf("failed to parse template: %w", err)
    }
    
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        return "", fmt.Errorf("failed to execute template: %w", err)
    }
    
    return buf.String(), nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/docgen/... -run TestFuncMap -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/docgen/template.go internal/docgen/template_test.go
git commit -m "feat(docgen): add template FuncMap with all functions"
```

---

## Task 10: Entry Point

**Files:**
- Create: `internal/docgen/cmd/generate.go`

**Why:** CLI entry point for Makefile to invoke. Processes all templates and writes output.

- [ ] **Step 1: Write the code**

```go
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
        // Skip if path contains '*' (shell already expanded)
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
    // docs/generate/commands.md.gotmpl -> docs/commands.md
    base := filepath.Base(tmplPath)
    name := strings.TrimSuffix(base, ".gotmpl")
    return filepath.Join("docs", name)
}
```

- [ ] **Step 2: Verify it compiles**

Run: `go build ./internal/docgen/cmd/generate.go`
Expected: No errors

- [ ] **Step 3: Commit**

```bash
git add internal/docgen/cmd/generate.go
git commit -m "feat(docgen): add generate entry point for Makefile"
```

---

## Task 11: Makefile Integration

**Files:**
- Modify: `Makefile`

**Why:** Add `make docs` target to regenerate documentation.

- [ ] **Step 1: Check current Makefile**

Run: `cat Makefile 2>/dev/null || echo "No Makefile exists"`
Note the existing targets.

- [ ] **Step 2: Add docs target**

```makefile
# Add to Makefile

.PHONY: docs

# Generate documentation from templates
docs:
	go run ./internal/docgen/cmd/generate.go docs/generate/*.md.gotmpl
```

- [ ] **Step 3: Test the target**

Run: `make docs`
Expected: Processes all templates (may fail if templates don't exist yet)

- [ ] **Step 4: Commit**

```bash
git add Makefile
git commit -m "feat: add docs target to Makefile"
```

---

## Task 12: Migrate commands.md Template

**Files:**
- Create: `docs/generate/commands.md.gotmpl`
- Remove: `docs/generate/commands.md.tmpl`
- Remove: `docs/generate/commands_examples/*.sh`
- Remove: `docs/generate/commands_examples/*.txt`
- Remove: `docs/generate/commands_examples/config.yaml`

**Why:** First template migration - prove the new system works end-to-end.

- [ ] **Step 1: Create new template**

```gotmpl
# CLI Commands Reference

> Purpose: Complete reference for all autoebiten CLI commands
> Audience: Developers using the CLI for game automation

---

## Global Flags

These flags apply to all commands:

| Flag | Description |
|------|-------------|
| --pid, -p | Target game process PID (auto-detected if not specified) |

If multiple games are running and --pid is not specified, autoebiten will list available games and exit with an error.

---

{{ config gameDir="examples/custom_commands" normalize=[
    {pattern: "PID=\\d+", replace: "PID=<PID>"},
    {pattern: "\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}", replace: "<TIMESTAMP>"},
    {pattern: "/[^\s]*\\.png", replace: "<PATH>.png"}
] }}

## ping

Check if game is running and responsive.

```bash
autoebiten ping
```

{{ launch_game }}
Output:
{{ command "ping" }}
{{ end_game }}

---

## version

Print CLI and game library versions.

```bash
autoebiten version
```

{{ launch_game }}
Output:
{{ command "version" }}
{{ end_game }}
```

Save as `docs/generate/commands.md.gotmpl` (continue with more commands from original template).

- [ ] **Step 2: Test generation**

Run: `go run ./internal/docgen/cmd/generate.go docs/generate/commands.md.gotmpl`
Expected: Generates `docs/commands.md`

- [ ] **Step 3: Verify output matches original**

Run: `diff docs/commands.md docs/original_docs/commands.md`
Expected: Minor differences acceptable (formatting, exact output)

- [ ] **Step 4: Delete old files**

```bash
rm docs/generate/commands.md.tmpl
rm -rf docs/generate/commands_examples/
```

- [ ] **Step 5: Commit**

```bash
git add docs/generate/commands.md.gotmpl
git add -u docs/generate/  # Stage deletions
git commit -m "feat(docs): migrate commands.md to new template system"
```

---

## Task 13: Migrate Remaining Templates

**Files:**
- Create: `docs/generate/testkit.md.gotmpl`
- Create: `docs/generate/tutorial.md.gotmpl`
- Create: `docs/generate/integration.md.gotmpl`
- Create: `docs/generate/autoui.md.gotmpl`
- Remove: Old templates and example directories

**Why:** Complete migration of all documentation templates.

- [ ] **Step 1: Migrate testkit.md**

Similar to Task 12, but also includes Go code extraction:

```gotmpl
{{ gocode file="examples/testkit/white_box_test.go" function="TestPlayerMovesRight" transforms=(list "stripImports" "package:mygame") }}
```

- [ ] **Step 2: Migrate tutorial.md**

Include crash scenarios:

```gotmpl
{{ launch_game gameDir="examples/crash_diagnostic" args=(list "--crash-before-rpc") }}
{{ $before := command "ping" }}
{{ end_game }}
```

- [ ] **Step 3: Migrate integration.md and autoui.md**

Follow same pattern.

- [ ] **Step 4: Verify all docs generate correctly**

Run: `make docs`
Expected: All 5 docs generated successfully

- [ ] **Step 5: Commit**

```bash
git add docs/generate/*.gotmpl
git add -u docs/generate/
git commit -m "feat(docs): migrate all templates to new system"
```

---

## Task 14: Cleanup

**Files:**
- Remove: `cmd/docgen/main.go`

**Why:** Old docgen tool no longer needed.

- [ ] **Step 1: Verify new system works**

Run: `make docs && git status --short`
Expected: All docs regenerated, no errors

- [ ] **Step 2: Remove old docgen**

```bash
rm -rf cmd/docgen/
```

- [ ] **Step 3: Update go.mod if needed**

Run: `go mod tidy`

- [ ] **Step 4: Commit**

```bash
git add -u cmd/
git commit -m "chore: remove old docgen tool"
```

---

## Task 15: Final Verification

**Files:**
- Verify all docs match originals

**Why:** Ensure the migration was successful.

- [ ] **Step 1: Regenerate all docs**

Run: `make docs`

- [ ] **Step 2: Compare with originals**

Run: `for f in docs/original_docs/*.md; do echo "=== $f ==="; diff "$f" "docs/$(basename $f)" || true; done`

Expected: Minor differences acceptable (formatting, exact output values)

- [ ] **Step 3: Run all tests**

Run: `go test ./internal/docgen/... -v`

Expected: All tests pass

- [ ] **Step 4: Final commit if needed**

```bash
git add -A
git commit -m "docs: complete template system migration"
```

---

## Summary

| Phase | Tasks | Description |
|-------|-------|-------------|
| Core Infrastructure | 1-4 | Context, config, normalization |
| Game Control | 5-6 | Command execution, delay |
| Verification | 7 | Output comparison |
| Code Extraction | 8 | Go AST parsing |
| Integration | 9-11 | FuncMap, entry point, Makefile |
| Migration | 12-13 | Convert templates |
| Cleanup | 14-15 | Remove old files, verify |