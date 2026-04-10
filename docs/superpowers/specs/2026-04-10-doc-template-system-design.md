---
name: Documentation Template System Rewrite
description: Go template system for autoebiten docs with inline config, direct command execution, and Go code extraction
type: project
---

# Documentation Template System Rewrite

**Date:** 2026-04-10
**Status:** Design Approved

---

## Overview

Rewrite the documentation generation system to use Go `text/template` with inline configuration, eliminating shell scripts and config files. Templates directly invoke autoebiten commands, verify variations, and extract Go code from source files.

**Note:** This design uses Go's native `text/template` package, not Jinja2. The syntax is inspired by Jinja but implemented as Go template FuncMap functions.

---

## Goals

1. **Single source of truth** - All config, commands, and outputs in template file
2. **No shell scripts** - Direct autoebiten invocation from templates
3. **No config.yaml** - Inline normalization rules
4. **Variation verification** - Show all variants, verify identical output
5. **Go code extraction** - Extract from source files with transforms (strip imports, rename, package change)
6. **Crash output** - Run actual crash scenarios, capture real proxy output

---

## Removed Files

| File Type | Current | After |
|-----------|---------|-------|
| Config files | `docs/generate/*/config.yaml` | ❌ Removed |
| Shell scripts | `docs/generate/*/*.sh` | ❌ Removed |
| Output files | `docs/generate/*/*_out.txt` | ❌ Removed |
| Templates | `docs/generate/*.md.tmpl` | → `docs/generate/*.md.gotmpl` |

---

## Architecture

### Makefile-Based Generation

```makefile
.PHONY: docs

docs:
	go run ./internal/docgen/cmd/generate.go docs/generate/*.md.gotmpl
```

Single Go command processes all templates. No separate docgen binary.

### Template Processing Pipeline

```
Template (.md.gotmpl)
    ↓ Parse template, extract {{ config }} block
    ↓ Build game binary (go build)
    ↓ Execute {{ launch_game }} → {{ command }} → {{ end_game }}
    ↓ Apply normalization to captured output
    ↓ {{ verify }} checks variation consistency
    ↓ {{ gocode }} extracts and transforms Go source
    ↓ Render final markdown
Generated doc (.md)
```

---

## Template Syntax

### Section Configuration

Replace `config.yaml` with inline config block:

```gotmpl
{{ config gameDir="examples/custom_commands" normalize=[
    {pattern: "PID=\\d+", replace: "PID=<PID>"},
    {pattern: "\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}", replace: "<TIMESTAMP>"},
    {pattern: "/[^\s]*\\.png", replace: "<PATH>.png"}
] }}
```

**Behavior:**
- Declared once per template, applies to all commands in that template
- `gameDir` - which example game to build and launch
- `normalize` - regex → replacement pairs for output normalization

---

### Command Execution

Direct autoebiten invocation:

```gotmpl
## ping

{{ launch_game }}
{{ command "ping" }}
{{ end_game }}
```

**Generated output:**
```text
OK: game is running
```

**Multiple commands:**
```gotmpl
{{ launch_game }}
{{ command "input" key="KeySpace" action="press" }}
{{ command "screenshot" output="test.png" }}
{{ end_game }}
```

**Command with flags:**
```gotmpl
{{ command "mouse" x=100 y=200 }}
{{ command "mouse" x=100 y=200 button="MouseButtonLeft" }}
{{ command "input" key="KeyW" action="hold" durationTicks=60 }}
```

---

### Variation Verification

Show all variations, verify identical output:

```gotmpl
## ping

{{ launch_game }}
{{ $variations := list
    (tuple "ping" dict)
    (tuple "ping" dict "pid" "<PID>")
    (tuple "ping" dict "timeout" "30s")
}}

**Examples:**
```bash
{{ range $variations }}
{{ $cmd := index . 0 }}
{{ $args := index . 1 }}
autoebiten {{ $cmd }}{{ range $k, $v := $args }} --{{ $k }} {{ $v }}{{ end }}
{{ end }}
```

{{ verify $variations }}

Output (same for all variations):
{{ lastOutput }}
{{ end_game }}
```

**Generated output:**
```bash
autoebiten ping
autoebiten ping --pid <PID>
autoebiten ping --timeout 30s
```

Output (same for all variations):
```text
OK: game is running
```

**Behavior:**
- `verify` runs all commands, applies normalization
- Fails generation if normalized outputs differ
- `lastOutput` contains verified result

---

### Go Code Extraction

Extract from source files with AST-based transforms:

```gotmpl
{{ gocode 
    file="examples/testkit/white_box_test.go"
    function="TestPlayerMovesRight"
    transforms=(list "stripImports" "rename:state_exporter.Game→Game" "package:mygame")
}}
```

**Transforms:**

| Transform | Description |
|-----------|-------------|
| `stripImports` | Remove import block |
| `rename:Old→New` | Replace type/variable names via AST |
| `package:Name` | Change package declaration |
| `placeholder` | Insert `// ...` where code removed |

**Multiple extractions:**
```gotmpl
{{ gocode file="examples/state_exporter/game.go" struct="GameState" }}
{{ gocode file="examples/state_exporter/game.go" function="NewGame" transforms=(list "stripImports") }}
```

**Implementation:**
- Uses `go/parser` and `go/ast` for precise extraction
- AST manipulation for transforms (reliable, not regex-based)
- Preserves formatting via `format.Node`

---

### Crash Output

Run crash scenarios, capture real proxy output:

```gotmpl
## Crash Diagnostic Testing

{{ launch_game gameDir="examples/crash_diagnostic" args=(list "--crash-before-rpc") }}
{{ $beforeCrash := command "ping" }}
{{ end_game }}

{{ launch_game gameDir="examples/crash_diagnostic" args=(list "--crash-after-rpc") }}
{{ command "ping" }}  <!-- returns OK initially -->
{{ delay "4s" }}
{{ $afterCrash := command "ping" }}
{{ end_game }}

**Test 1: Crash before RPC**
{{ $beforeCrash }}

**Test 2: Crash after RPC**
{{ $afterCrash }}
```

**Generated output:**
```text
**Test 1: Crash before RPC**
<log_diff>
--- snapshot ...
+++ current ...
@@ ... @@
+Starting crash diagnostic demo...
+Flags: crash-before-rpc=true, crash-after-rpc=false
+Initialization complete
+About to crash before RPC connection!
+panic: intentional crash before RPC connection
</log_diff>
<proxy_error>
game exited: exit status 2
</proxy_error>
Error: game not connected

**Test 2: Crash after RPC**
<log_diff>
--- snapshot ...
+++ current ...
@@ ... @@
+Game running... tick 60
+Game running... tick 120
+Game running... tick 180
+panic: intentional crash after RPC connection
</log_diff>
<proxy_error>
game exited: exit status 2
</proxy_error>
Error: game not connected
```

**Behavior:**
- `launch_game` starts game with crash flags via testkit
- `command "ping"` on crashed game returns proxy error
- `delay` waits for crash (after-rpc scenarios)
- Output captured and normalized inline

---

## Template Functions

### Config Functions

| Function | Purpose |
|----------|---------|
| `config` | Set gameDir and normalize rules for template |

### Game Control Functions

| Function | Purpose |
|----------|---------|
| `launch_game` | Build and start game via testkit, optional args |
| `end_game` | Shutdown game, cleanup process |
| `command` | Execute autoebiten CLI command, capture output |
| `delay` | Wait for specified duration (for crash timing) |

### Verification Functions

| Function | Purpose |
|----------|---------|
| `verify` | Run command variations, check identical output |
| `lastOutput` | Return last captured normalized output |

### Code Extraction Functions

| Function | Purpose |
|----------|---------|
| `gocode` | Extract function/struct from Go file with transforms |

### Utility Functions

| Function | Purpose |
|----------|---------|
| `list` | Create slice for Go template |
| `tuple` | Create tuple/slice |
| `dict` | Create map for Go template |

---

## Implementation

### internal/docgen/cmd/generate.go

Single entry point invoked by Makefile:

```go
package main

import (
    "os"
    "path/filepath"
    
    "github.com/s3cy/autoebiten/internal/docgen"
)

func main() {
    if len(os.Args) < 2 {
        fatal("usage: go run generate.go <templates...>")
    }
    
    for _, tmplPath := range os.Args[1:] {
        outPath := outputPath(tmplPath) // strip .gotmpl → .md in docs/
        
        result, err := docgen.ProcessTemplate(tmplPath)
        if err != nil {
            fatal(err)
        }
        
        os.WriteFile(outPath, []byte(result), 0644)
        fmt.Printf("Generated: %s → %s\n", tmplPath, outPath)
    }
}

func outputPath(tmplPath string) string {
    base := filepath.Base(tmplPath)
    name := strings.TrimSuffix(base, ".gotmpl")
    return filepath.Join("docs", name)
}
```

### internal/docgen/template.go

Core template processing with FuncMap:

```go
package docgen

import (
    "text/template"
)

func ProcessTemplate(tmplPath string) (string, error) {
    content, err := os.ReadFile(tmplPath)
    if err != nil {
        return "", err
    }
    
    tmpl := template.New(filepath.Base(tmplPath)).Funcs(FuncMap())
    
    // Parse and execute
    // ...
}

func FuncMap() template.FuncMap {
    return template.FuncMap{
        "config":      ConfigFunc,
        "launch_game": LaunchGameFunc,
        "end_game":    EndGameFunc,
        "command":     CommandFunc,
        "delay":       DelayFunc,
        "verify":      VerifyFunc,
        "lastOutput":  LastOutputFunc,
        "gocode":      GocodeFunc,
        "list":        ListFunc,
        "tuple":       TupleFunc,
        "dict":        DictFunc,
    }
}
```

### internal/docgen/game.go

Game lifecycle management using testkit:

```go
package docgen

import (
    "testing"
    "time"
    
    "github.com/s3cy/autoebiten/testkit"
)

type GameSession struct {
    game    *testkit.Game
    t       *testing.T
    config  *Config
    outputs []string
}

func LaunchGameFunc(gameDir string, args ...string) (*GameSession, error) {
    // Build binary
    binary := buildGame(gameDir)
    
    // Launch via testkit
    t := &testing.T{}
    opts := []testkit.Option{testkit.WithTimeout(30 * time.Second)}
    if len(args) > 0 {
        opts = append(opts, testkit.WithArgs(args...))
    }
    
    game := testkit.Launch(t, binary, opts...)
    
    // Wait for ready
    ready := game.WaitFor(func() bool {
        return game.Ping() == nil
    }, 5 * time.Second)
    
    if !ready {
        game.Shutdown()
        return nil, errors.New("game failed to start")
    }
    
    return &GameSession{game: game, t: t}, nil
}

func CommandFunc(session *GameSession, cmd string, flags ...map[string]any) string {
    // Invoke autoebiten CLI with flags
    // Capture output, normalize, store in session.outputs
    // Return normalized output
}
```

### internal/docgen/gocode.go

Go code extraction with AST transforms:

```go
package docgen

import (
    "go/parser"
    "go/ast"
    "go/format"
    "go/token"
)

func GocodeFunc(file, target string, transforms []string) (string, error) {
    fset := token.NewFileSet()
    astFile, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
    if err != nil {
        return "", err
    }
    
    // Find target (function or struct)
    node := findTarget(astFile, target)
    
    // Apply transforms via AST manipulation
    for _, t := range transforms {
        applyTransform(astFile, node, t)
    }
    
    // Format and return
    var buf bytes.Buffer
    format.Node(&buf, fset, node)
    return buf.String(), nil
}

func applyTransform(file *ast.File, node ast.Node, transform string) {
    switch {
    case transform == "stripImports":
        file.Decls = removeImports(file.Decls)
    case strings.HasPrefix(transform, "rename:"):
        // Parse Old→New, walk AST and replace identifiers
        parts := strings.Split(strings.TrimPrefix(transform, "rename:"), "→")
        renameIdentifiers(node, parts[0], parts[1])
    case strings.HasPrefix(transform, "package:"):
        // Change package name
        file.Name.Name = strings.TrimPrefix(transform, "package:")
    }
}
```

---

## Why No Separate docverify Tool?

Verification happens inline during generation via `{{ verify }}`. No separate `docverify` binary needed.

---

## File Structure After Rewrite

```
docs/
├── generate/
│   ├── commands.md.gotmpl
│   ├── testkit.md.gotmpl
│   ├── tutorial.md.gotmpl
│   ├── integration.md.gotmpl
│   └── autoui.md.gotmpl
├── commands.md         (generated)
├── testkit.md          (generated)
├── tutorial.md         (generated)
├── integration.md      (generated)
├── autoui.md           (generated)

internal/docgen/
├── cmd/
│   └── generate.go     (entry point)
├── template.go         (FuncMap, template processing)
├── game.go             (launch_game, command, end_game)
├── gocode.go           (gocode extraction with AST transforms)
├── verify.go           (verify variations)
├── config.go           (config parsing, normalization)
├── normalize.go        (regex normalization)

Makefile
```

---

## Migration Plan

1. Create new `internal/docgen/` package with FuncMap functions
2. Create `internal/docgen/cmd/generate.go` entry point
3. Rewrite `commands.md.tmpl` → `commands.md.gotmpl` using new syntax
4. Delete shell scripts and config files for commands section
5. Repeat for each doc section (testkit, tutorial, integration, autoui)
6. Update Makefile with `docs:` target
7. Remove old `cmd/docgen/main.go`
8. Verify generated docs match original docs

---

## Summary

| Feature | Implementation |
|---------|---------------|
| Config removal | `{{ config }}` inline |
| Shell scripts removal | `{{ command }}` direct invocation via testkit |
| Variations | `{{ verify }}` with template loop display |
| Go code extraction | `{{ gocode }}` with native AST parser |
| Crash output | `{{ launch_game args="--crash-before-rpc" }}` direct run |
| Normalization | Inline regex rules in config block |
| Tooling | Makefile + `go run ./internal/docgen/cmd/generate.go` |