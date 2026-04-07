# Carriage Return Interpretation + Mutex-Protected Output Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make the log file reflect "what the terminal showed" (interpreting `\r` as line overwrite) while ensuring safe concurrent access for write, copy, and diff operations.

**Architecture:** Custom `CarriageReturnWriter` interprets `\r` before writing to `OutputManager`. `OutputManager` holds mutex to serialize write/diff/snapshot operations across concurrent goroutines.

**Tech Stack:** Go 1.25, standard library (`sync`, `io`, `os`), stretchr/testify for testing

---

## File Structure

| File | Responsibility |
|------|----------------|
| `internal/output/manager.go` | New: `CarriageReturnWriter` + `OutputManager` |
| `internal/output/manager_test.go` | New: Tests for `CarriageReturnWriter` and `OutputManager` |
| `internal/output/output.go` | Modify: Remove `GenerateDiff`, `CopyFile`, keep file path helpers |
| `internal/output/output_test.go` | Modify: Remove tests for deleted functions |
| `internal/cli/launch.go` | Modify: Use `OutputManager` instead of direct file write |
| `internal/proxy/server.go` | Modify: Use `OutputManager.DiffAndUpdateSnapshot()` |

---

### Task 1: Create CarriageReturnWriter

**Files:**
- Create: `internal/output/manager.go`
- Create: `internal/output/manager_test.go`

- [ ] **Step 1: Write the failing test for CarriageReturnWriter**

Create `internal/output/manager_test.go`:

```go
package output

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCarriageReturnWriterBasicOverwrite(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single overwrite",
			input:    "Loading 10%\rLoading 20%\n",
			expected: "Loading 20%\n",
		},
		{
			name:     "multiple overwrites",
			input:    "0%\r25%\r50%\r75%\r100%\n",
			expected: "100%\n",
		},
		{
			name:     "no overwrite just newline",
			input:    "Hello World\n",
			expected: "Hello World\n",
		},
		{
			name:     "incomplete line flushed at end",
			input:    "Loading 50%",
			expected: "Loading 50%\n",
		},
		{
			name:     "overwrite then incomplete line",
			input:    "Loading 0%\rLoading 100%",
			expected: "Loading 100%\n",
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "multiple lines",
			input:    "Line1\nLine2\rLine2Modified\nLine3\n",
			expected: "Line1\nLine2Modified\nLine3\n",
		},
		{
			name:     "multiple carriage returns before newline",
			input:    "abc\rdef\rghi\n",
			expected: "ghi\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := &bytes.Buffer{}
			writer := NewCarriageReturnWriter(dst)
			writer.Write([]byte(tt.input))
			writer.Flush()
			assert.Equal(t, tt.expected, dst.String())
		})
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/output/... -run TestCarriageReturnWriterBasicOverwrite -v`
Expected: FAIL with "NewCarriageReturnWriter not defined" or similar

- [ ] **Step 3: Write minimal implementation for CarriageReturnWriter**

Create `internal/output/manager.go` with `CarriageReturnWriter`:

```go
package output

import (
	"io"
)

// CarriageReturnWriter interprets \r as "move cursor to line start, overwrite".
// It writes to dst only the final visual state of each line.
type CarriageReturnWriter struct {
	dst     io.Writer
	lineBuf []byte
}

// NewCarriageReturnWriter creates a new carriage return interpreting writer.
func NewCarriageReturnWriter(dst io.Writer) *CarriageReturnWriter {
	return &CarriageReturnWriter{
		dst:     dst,
		lineBuf: make([]byte, 0, 1024),
	}
}

// Write processes input bytes, interpreting \r as line overwrite.
func (w *CarriageReturnWriter) Write(data []byte) (int, error) {
	for _, b := range data {
		switch b {
		case '\r':
			// Move cursor to line start: clear line buffer
			w.lineBuf = w.lineBuf[:0]
		case '\n':
			// Flush completed line
			if len(w.lineBuf) > 0 {
				w.dst.Write(w.lineBuf)
			}
			w.dst.Write([]byte{'\n'})
			w.lineBuf = w.lineBuf[:0]
		default:
			w.lineBuf = append(w.lineBuf, b)
		}
	}
	return len(data), nil
}

// Flush writes any remaining incomplete line (with newline appended).
func (w *CarriageReturnWriter) Flush() error {
	if len(w.lineBuf) > 0 {
		if _, err := w.dst.Write(w.lineBuf); err != nil {
			return err
		}
		if _, err := w.dst.Write([]byte{'\n'}); err != nil {
			return err
		}
		w.lineBuf = w.lineBuf[:0]
	}
	return nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/output/... -run TestCarriageReturnWriterBasicOverwrite -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/output/manager.go internal/output/manager_test.go
git commit -m "feat(output): add CarriageReturnWriter to interpret \\r as line overwrite"
```

---

### Task 2: Create OutputManager

**Files:**
- Modify: `internal/output/manager.go`
- Modify: `internal/output/manager_test.go`

- [ ] **Step 1: Write the failing test for OutputManager.Write**

Add to `internal/output/manager_test.go`:

```go
func TestOutputManagerWrite(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	snapPath := filepath.Join(tmpDir, "test-snapshot.log")

	logFile, err := os.Create(logPath)
	require.NoError(t, err)

	manager := NewOutputManager(logFile, logPath, snapPath)

	// Write some data
	n, err := manager.Write([]byte("Hello\n"))
	assert.NoError(t, err)
	assert.Equal(t, 6, n)

	// Read file to verify
	content, err := os.ReadFile(logPath)
	require.NoError(t, err)
	assert.Equal(t, "Hello\n", string(content))
}
```

Add imports:
```go
import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/output/... -run TestOutputManagerWrite -v`
Expected: FAIL with "NewOutputManager not defined"

- [ ] **Step 3: Write minimal OutputManager.Write implementation**

Add to `internal/output/manager.go`:

```go
import (
	"io"
	"os"
	"os/exec"
	"sync"
)

// OutputManager coordinates mutex-protected write, diff, and snapshot operations.
type OutputManager struct {
	mu           sync.Mutex
	logFile      *os.File
	outputPath   string
	snapshotPath string
}

// NewOutputManager creates a new output manager.
func NewOutputManager(logFile *os.File, outputPath, snapshotPath string) *OutputManager {
	return &OutputManager{
		logFile:      logFile,
		outputPath:   outputPath,
		snapshotPath: snapshotPath,
	}
}

// Write writes data to the log file with mutex protection and flush.
func (m *OutputManager) Write(data []byte) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	n, err := m.logFile.Write(data)
	if err != nil {
		return n, err
	}
	m.logFile.Sync()
	return n, nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/output/... -run TestOutputManagerWrite -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/output/manager.go internal/output/manager_test.go
git commit -m "feat(output): add OutputManager.Write with mutex protection"
```

---

### Task 3: Implement OutputManager.DiffAndUpdateSnapshot

**Files:**
- Modify: `internal/output/manager.go`
- Modify: `internal/output/manager_test.go`

- [ ] **Step 1: Write the failing test for DiffAndUpdateSnapshot**

Add to `internal/output/manager_test.go`:

```go
func TestOutputManagerDiffAndUpdateSnapshot(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	snapPath := filepath.Join(tmpDir, "test-snapshot.log")

	// Create initial snapshot
	os.WriteFile(snapPath, []byte("old line\n"), 0600)

	// Create log with new content
	logFile, err := os.Create(logPath)
	require.NoError(t, err)
	logFile.WriteString("new line\n")
	logFile.Close()

	// Reopen log file for manager
	logFile, err = os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND, 0600)
	require.NoError(t, err)

	manager := NewOutputManager(logFile, logPath, snapPath)

	diff, err := manager.DiffAndUpdateSnapshot()
	require.NoError(t, err)
	assert.Contains(t, diff, "new line")

	// Verify snapshot was updated
	snapContent, err := os.ReadFile(snapPath)
	require.NoError(t, err)
	assert.Equal(t, "new line\n", string(snapContent))
}

func TestOutputManagerDiffAndUpdateSnapshotNoChanges(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	snapPath := filepath.Join(tmpDir, "test-snapshot.log")

	// Same content in both
	os.WriteFile(snapPath, []byte("same content\n"), 0600)
	logFile, err := os.Create(logPath)
	require.NoError(t, err)
	logFile.WriteString("same content\n")
	logFile.Close()

	logFile, err = os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND, 0600)
	require.NoError(t, err)

	manager := NewOutputManager(logFile, logPath, snapPath)

	diff, err := manager.DiffAndUpdateSnapshot()
	require.NoError(t, err)
	assert.Empty(t, diff) // No diff when identical
}

func TestOutputManagerDiffAndUpdateSnapshotEmptyFiles(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	snapPath := filepath.Join(tmpDir, "test-snapshot.log")

	// Empty log file
	logFile, err := os.Create(logPath)
	require.NoError(t, err)

	manager := NewOutputManager(logFile, logPath, snapPath)

	diff, err := manager.DiffAndUpdateSnapshot()
	require.NoError(t, err)
	assert.Empty(t, diff) // Empty diff when both empty
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/output/... -run TestOutputManagerDiffAndUpdateSnapshot -v`
Expected: FAIL with "DiffAndUpdateSnapshot not defined"

- [ ] **Step 3: Write DiffAndUpdateSnapshot implementation**

Add to `internal/output/manager.go`:

```go
// DiffAndUpdateSnapshot generates diff between snapshot and log, then copies log to snapshot.
// Returns the diff string. Must be called with mutex already held or will acquire it.
func (m *OutputManager) DiffAndUpdateSnapshot() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Flush log file to ensure diff sees latest content
	m.logFile.Sync()

	// Generate diff using diff -u command
	cmd := exec.Command("bash", "-c",
		"diff -u '"+m.snapshotPath+"' '"+m.outputPath+"' 2>/dev/null || true")

	output, err := cmd.Output()
	if err != nil {
		// diff returns exit code 1 when files differ (expected)
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return string(output), nil
		}
		return "", err
	}

	// Copy log to snapshot
	content, err := os.ReadFile(m.outputPath)
	if err != nil {
		if os.IsNotExist(err) {
			content = nil
		} else {
			return string(output), err
		}
	}

	if err := writeSnapshot(m.snapshotPath, content); err != nil {
		return string(output), err
	}

	return string(output), nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/output/... -run TestOutputManagerDiffAndUpdateSnapshot -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/output/manager.go internal/output/manager_test.go
git commit -m "feat(output): add OutputManager.DiffAndUpdateSnapshot"
```

---

### Task 4: Add concurrency test for OutputManager

**Files:**
- Modify: `internal/output/manager_test.go`

- [ ] **Step 1: Write the failing concurrency test**

Add to `internal/output/manager_test.go`:

```go
func TestOutputManagerConcurrentWrite(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	snapPath := filepath.Join(tmpDir, "test-snapshot.log")

	logFile, err := os.Create(logPath)
	require.NoError(t, err)

	manager := NewOutputManager(logFile, logPath, snapPath)

	// Write concurrently
	done := make(chan bool, 10)
	for i := range 10 {
		go func(n int) {
			manager.Write([]byte("line " + strconv.Itoa(n) + "\n"))
			done <- true
		}(i)
	}

	// Wait for all writes
	for range 10 {
		<-done
	}

	// Read file and verify all lines present
	content, err := os.ReadFile(logPath)
	require.NoError(t, err)

	for i := range 10 {
		assert.Contains(t, string(content), "line "+strconv.Itoa(i))
	}
}
```

Add import: `"strconv"`

- [ ] **Step 2: Run test to verify it passes (mutex should already work)**

Run: `go test ./internal/output/... -run TestOutputManagerConcurrentWrite -v`
Expected: PASS (mutex already in implementation)

- [ ] **Step 3: Commit**

```bash
git add internal/output/manager_test.go
git commit -m "test(output): add concurrency test for OutputManager"
```

---

### Task 5: Remove GenerateDiff and CopyFile from output.go

**Files:**
- Modify: `internal/output/output.go`
- Modify: `internal/output/output_test.go`

- [ ] **Step 1: Remove GenerateDiff and CopyFile from output.go**

Edit `internal/output/output.go`, remove lines 72-110:
```go
// CopyFile copies src to dst, creating dst if needed.
// Returns error if src doesn't exist (not os.IsNotExist).
func CopyFile(dst, src string) error { ... }

// GenerateDiff generates a unified diff between snapshot and current file paths.
// Uses diff -u directly on files with sed to handle carriage returns.
func GenerateDiff(snapshotPath, currentPath string) string { ... }
```

Also remove `"os/exec"` import if no longer needed.

Keep `writeSnapshot` function (used by `DiffAndUpdateSnapshot`).

- [ ] **Step 2: Remove corresponding tests from output_test.go**

Edit `internal/output/output_test.go`, remove:
- `TestGenerateDiff` function (lines 129-207)
- `TestGenerateDiffWithCarriageReturn` function (lines 209-232)

- [ ] **Step 3: Run tests to verify**

Run: `go test ./internal/output/... -v`
Expected: PASS (all remaining tests pass)

- [ ] **Step 4: Commit**

```bash
git add internal/output/output.go internal/output/output_test.go
git commit -m "refactor(output): remove GenerateDiff and CopyFile, moved to OutputManager"
```

---

### Task 6: Integrate OutputManager into launch.go

**Files:**
- Modify: `internal/cli/launch.go`

- [ ] **Step 1: Modify launch.go to use OutputManager**

In `LaunchCommand` struct, add field:
```go
type LaunchCommand struct {
	options      *LaunchOptions
	outputFiles  *output.FilePath
	outputMgr    *output.OutputManager  // Add this
	gameProc     *os.Process
	proxyServer  *proxy.Server
	proxyHandler *proxy.Handler
	listener     net.Listener
	gameExited   chan struct{}
	done         chan struct{}
}
```

In `Run` method, after creating log file:
```go
// Create log file
logFile, err := output.CreateLogFile(lc.outputFiles.Log)
if err != nil {
	lc.terminateGame()
	return fmt.Errorf("failed to create log file: %w", err)
}

// Create OutputManager
lc.outputMgr = output.NewOutputManager(logFile, lc.outputFiles.Log, lc.outputFiles.Snapshot)

// Tee stdout/stderr through CarriageReturnWriter to OutputManager
stdoutWriter := output.NewCarriageReturnWriter(lc.outputMgr)
stderrWriter := output.NewCarriageReturnWriter(lc.outputMgr)
go lc.teeOutput(stdoutPipe, os.Stdout, stdoutWriter)
go lc.teeOutput(stderrPipe, os.Stderr, stderrWriter)
```

- [ ] **Step 2: Modify teeOutput signature**

Change `teeOutput` method:
```go
// teeOutput copies data from src to both dst1 (terminal) and dst2 (managed writer).
func (lc *LaunchCommand) teeOutput(src io.Reader, dst1 *os.File, dst2 io.Writer) {
	reader := bufio.NewReader(src)
	for {
		data, err := reader.ReadBytes('\n')
		if len(data) > 0 {
			dst1.Write(data)      // Terminal gets raw bytes (it interprets \r)
			dst2.Write(data)      // CarriageReturnWriter + OutputManager
		}
		if err != nil {
			if err == io.EOF {
				// Flush any remaining data at stream end
				remaining, _ := reader.ReadBytes('\n')
				if len(remaining) > 0 {
					dst1.Write(remaining)
					dst2.Write(remaining)
				}
				// Flush the CarriageReturnWriter
				if flusher, ok := dst2.(interface{ Flush() error }); ok {
					flusher.Flush()
				}
			}
			break
		}
	}
}
```

- [ ] **Step 3: Run tests to verify**

Run: `go test ./internal/cli/... -v`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add internal/cli/launch.go
git commit -m "feat(cli): integrate OutputManager and CarriageReturnWriter in launch"
```

---

### Task 7: Integrate OutputManager into proxy/server.go

**Files:**
- Modify: `internal/proxy/server.go`

- [ ] **Step 1: Modify Server struct to use OutputManager**

Edit `internal/proxy/server.go`:
```go
// Server wraps game RPC calls and captures output.
type Server struct {
	gameClient  GameClient
	outputMgr   *output.OutputManager  // Replace outputFiles with outputMgr
	mu          sync.Mutex             // Keep for ForwardRequest serialization
}
```

Update `NewServer`:
```go
// NewServer creates a new proxy server.
func NewServer(gameClient GameClient, outputMgr *output.OutputManager) *Server {
	return &Server{
		gameClient: gameClient,
		outputMgr:  outputMgr,
	}
}
```

- [ ] **Step 2: Modify ForwardRequest to use DiffAndUpdateSnapshot**

Edit `ForwardRequest`:
```go
func (s *Server) ForwardRequest(req *rpc.RPCRequest) (*Response, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Forward request to game
	gameResp, err := s.gameClient.SendRequest(req)
	if err != nil {
		return nil, fmt.Errorf("game request failed: %w", err)
	}

	// Generate diff and update snapshot
	diff, err := s.outputMgr.DiffAndUpdateSnapshot()
	if err != nil {
		return nil, fmt.Errorf("failed to generate diff: %w", err)
	}

	// Build proxy response
	resp := &Response{
		RPCResponse: *gameResp,
		Output:      diff,
	}

	return resp, nil
}
```

- [ ] **Step 3: Remove Cleanup method dependency on outputFiles**

The `Cleanup` method needs file paths. Pass them separately:
```go
type Server struct {
	gameClient  GameClient
	outputMgr   *output.OutputManager
	outputFiles *output.FilePath  // Keep for Cleanup
	mu          sync.Mutex
}

func NewServer(gameClient GameClient, outputMgr *output.OutputManager, outputFiles *output.FilePath) *Server {
	return &Server{
		gameClient:  gameClient,
		outputMgr:   outputMgr,
		outputFiles: outputFiles,
	}
}
```

- [ ] **Step 4: Update launch.go to pass outputFiles to NewServer**

In `launch.go`:
```go
// Create proxy server
lc.proxyServer = proxy.NewServer(gameClient, lc.outputMgr, lc.outputFiles)
```

- [ ] **Step 5: Run tests to verify**

Run: `go test ./internal/proxy/... -v`
Expected: Some tests may fail due to signature changes

- [ ] **Step 6: Fix proxy tests**

Update `internal/proxy/server_test.go` to use new signatures. For each test that creates a Server:
```go
// Create temp log file
logFile, err := os.Create(logPath)
require.NoError(t, err)

outputMgr := output.NewOutputManager(logFile, logPath, snapPath)
server := NewServer(mock, outputMgr, paths)
```

- [ ] **Step 7: Run tests to verify**

Run: `go test ./internal/proxy/... -v`
Expected: PASS

- [ ] **Step 8: Commit**

```bash
git add internal/proxy/server.go internal/proxy/server_test.go
git commit -m "feat(proxy): use OutputManager.DiffAndUpdateSnapshot in ForwardRequest"
```

---

### Task 8: Integration test

**Files:**
- Create: `e2e/output_manager_test.go`

- [ ] **Step 1: Write integration test**

Create `e2e/output_manager_test.go`:
```go
package e2e

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/s3cy/autoebiten/internal/output"
)

func TestCarriageReturnWriterIntegration(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "game.log")
	snapPath := filepath.Join(tmpDir, "snapshot.log")

	logFile, err := os.Create(logPath)
	require.NoError(t, err)

	outputMgr := output.NewOutputManager(logFile, logPath, snapPath)
	writer := output.NewCarriageReturnWriter(outputMgr)

	// Simulate game output with progress bar
	writer.Write([]byte("Loading: 0%"))
	writer.Write([]byte("\rLoading: 25%"))
	writer.Write([]byte("\rLoading: 50%"))
	writer.Write([]byte("\rLoading: 75%"))
	writer.Write([]byte("\rLoading: 100%\n"))
	writer.Write([]byte("Done!\n"))
	writer.Flush()

	// Verify log file shows final visual state
	content, err := os.ReadFile(logPath)
	require.NoError(t, err)

	// Should show final progress state and Done
	assert.Contains(t, string(content), "Loading: 100%\n")
	assert.Contains(t, string(content), "Done!\n")
	assert.NotContains(t, string(content), "Loading: 0%")
	assert.NotContains(t, string(content), "Loading: 25%")

	// Generate diff and update snapshot
	diff, err := outputMgr.DiffAndUpdateSnapshot()
	require.NoError(t, err)
	assert.Contains(t, diff, "Loading: 100%")

	// Verify snapshot was updated
	snapContent, err := os.ReadFile(snapPath)
	require.NoError(t, err)
	assert.Equal(t, string(content), string(snapContent))
}

func TestOutputManagerConcurrentWriteAndDiff(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "game.log")
	snapPath := filepath.Join(tmpDir, "snapshot.log")

	logFile, err := os.Create(logPath)
	require.NoError(t, err)

	outputMgr := output.NewOutputManager(logFile, logPath, snapPath)
	writer := output.NewCarriageReturnWriter(outputMgr)

	// Write in background (simulating game output)
	writeDone := make(chan bool)
	go func() {
		for i := range 5 {
			writer.Write([]byte("Progress: " + strconv.Itoa(i*20) + "%\r"))
			time.Sleep(10 * time.Millisecond)
		}
		writer.Write([]byte("Progress: 100%\nComplete!\n"))
		writer.Flush()
		writeDone <- true
	}()

	// Wait for writing to complete
	<-writeDone

	// Now diff
	diff, err := outputMgr.DiffAndUpdateSnapshot()
	require.NoError(t, err)
	assert.Contains(t, diff, "Complete!")
}
```

Add imports: `"strconv"`

- [ ] **Step 2: Run integration tests**

Run: `go test ./e2e/... -run TestCarriageReturn -v`
Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add e2e/output_manager_test.go
git commit -m "test(e2e): add integration tests for carriage return handling"
```

---

### Task 9: Run full test suite

- [ ] **Step 1: Run all tests with race detection**

Run: `go test -race ./...`
Expected: PASS (no race conditions)

- [ ] **Step 2: Run coverage check**

Run: `go test -cover ./internal/output/...`
Expected: 80%+ coverage for new code

- [ ] **Step 3: Commit if any minor fixes needed**

If fixes were needed:
```bash
git add -A
git commit -m "fix: address test failures/race conditions"
```

---

## Self-Review Checklist

1. **Spec coverage:**
   - ✅ `CarriageReturnWriter` interprets `\r` (Task 1)
   - ✅ `OutputManager.Write` with mutex (Task 2)
   - ✅ `OutputManager.DiffAndUpdateSnapshot` (Task 3)
   - ✅ Remove old functions (Task 5)
   - ✅ Integrate into launch.go (Task 6)
   - ✅ Integrate into proxy/server.go (Task 7)
   - ✅ Tests for all components (Tasks 1, 2, 3, 4, 8)

2. **Placeholder scan:** No TBD/TODO in plan

3. **Type consistency:**
   - `NewCarriageReturnWriter(dst io.Writer)` → matches `OutputManager.Write`
   - `NewOutputManager(logFile, outputPath, snapshotPath)` → matches usage in launch.go
   - `DiffAndUpdateSnapshot()` returns `(string, error)` → matches ForwardRequest usage

---

**Plan complete.** Ready for execution.