# integrate.AfterDraw() Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add integrate.AfterDraw() to handle all post-draw operations, enabling DrawHighlights in patch method without import cycles.

**Architecture:** Callback registry pattern - integrate exposes RegisterDrawHighlights(), autoui registers its callback during Register() call, patch calls AfterDraw() which invokes all registered callbacks.

**Tech Stack:** Go 1.25, ebiten v2.9.9

---

## File Structure

```
integrate/
├── integrate.go        # Modify: add RegisterDrawHighlights, AfterDraw
├── integrate_test.go   # Create: test AfterDraw and callback registration

autoui/
├── highlight.go        # Modify: add drawHighlightsCallback
├── register.go         # Modify: import integrate, call RegisterDrawHighlights

ebiten.patch            # Modify: change Capture to AfterDraw
```

---

### Task 1: Add RegisterDrawHighlights and AfterDraw to integrate package

**Files:**
- Modify: `integrate/integrate.go`
- Create: `integrate/integrate_test.go`

- [ ] **Step 1: Write the failing test for AfterDraw**

Create `integrate/integrate_test.go`:

```go
package integrate

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAfterDrawNoCallback tests AfterDraw when no callback is registered.
func TestAfterDrawNoCallback(t *testing.T) {
	// Reset callback to nil
	drawHighlightsFunc = nil

	// Create a dummy image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	// AfterDraw should not panic when no callback is registered
	AfterDraw(img)
	// No assertion needed - just verify it doesn't panic
}

// TestAfterDrawWithCallback tests AfterDraw invokes registered callback.
func TestAfterDrawWithCallback(t *testing.T) {
	// Track if callback was invoked
	callbackInvoked := false
	var receivedImage image.Image

	// Register a test callback
	RegisterDrawHighlights(func(screen image.Image) {
		callbackInvoked = true
		receivedImage = screen
	})

	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	// Call AfterDraw
	AfterDraw(img)

	// Verify callback was invoked with correct image
	assert.True(t, callbackInvoked, "Callback should be invoked")
	assert.Equal(t, img, receivedImage, "Callback should receive the same image")

	// Reset callback
	drawHighlightsFunc = nil
}

// TestRegisterDrawHighlights tests callback registration.
func TestRegisterDrawHighlights(t *testing.T) {
	// Reset
	drawHighlightsFunc = nil

	// Register first callback
	callback1 := func(screen image.Image) {}
	RegisterDrawHighlights(callback1)

	// Verify it's registered
	assert.NotNil(t, drawHighlightsFunc, "Callback should be registered")

	// Register second callback (should replace first)
	callback2 := func(screen image.Image) {}
	RegisterDrawHighlights(callback2)

	// Verify second callback is registered (simple pointer comparison won't work,
	// but we can verify it's not nil)
	assert.NotNil(t, drawHighlightsFunc, "Second callback should be registered")

	// Reset
	drawHighlightsFunc = nil
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./integrate/... -v`
Expected: FAIL with "AfterDraw undefined" or similar

- [ ] **Step 3: Add RegisterDrawHighlights and AfterDraw to integrate.go**

Modify `integrate/integrate.go`. The full file should look like:

```go
package integrate

import (
	"image"

	"github.com/s3cy/autoebiten/internal/input"
	"github.com/s3cy/autoebiten/internal/server"
)

type Key = input.Key
type MouseButton = input.MouseButton

// IsPatched indicates whether the game uses a patched version of Ebiten.
var IsPatched = false

// drawHighlightsFunc is the registered callback for drawing highlight overlays.
var drawHighlightsFunc func(screen image.Image)

// RegisterDrawHighlights registers a callback for drawing highlight overlays.
// The callback receives screen as image.Image; caller must type assert to *ebiten.Image.
// Called by autoui during Register() to enable highlights in patch method.
func RegisterDrawHighlights(fn func(screen image.Image)) {
	drawHighlightsFunc = fn
}

// AfterDraw handles all post-draw operations after game.Draw() completes.
// Currently handles: screenshots (ProcessScreenshots) and highlight overlays.
// Called by patched ebiten in DrawOffscreen().
func AfterDraw(screen image.Image) {
	server.ProcessScreenshots(screen)
	if drawHighlightsFunc != nil {
		drawHighlightsFunc(screen)
	}
}

// Capture processes screenshots for injection.
// Deprecated: Use AfterDraw instead. Kept for backward compatibility with library method.
func Capture(screen image.Image) {
	server.ProcessScreenshots(screen)
}

// Update runs the internal update loop.
func Update() bool {
	return server.Update()
}

// ... rest of existing functions unchanged ...
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./integrate/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add integrate/integrate.go integrate/integrate_test.go
git commit -m "feat(integrate): add RegisterDrawHighlights and AfterDraw for patch method highlights"
```

---

### Task 2: Add drawHighlightsCallback to autoui package

**Files:**
- Modify: `autoui/highlight.go`

- [ ] **Step 1: Add import for integrate package**

Add `"github.com/s3cy/autoebiten/integrate"` to imports in `autoui/highlight.go`:

```go
import (
	"image"
	"image/color"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/s3cy/autoebiten/integrate"
)
```

- [ ] **Step 2: Add drawHighlightsCallback function**

Add after `globalHighlightManager` declaration (around line 115):

```go
// drawHighlightsCallback is the callback registered with integrate for patch method.
// It type asserts image.Image to *ebiten.Image for vector drawing.
func drawHighlightsCallback(screen image.Image) {
	if ebiScreen, ok := screen.(*ebiten.Image); ok {
		globalHighlightManager.draw(ebiScreen)
	}
}
```

- [ ] **Step 3: Run tests to verify nothing is broken**

Run: `go test ./autoui/... -v`
Expected: All existing tests PASS

- [ ] **Step 4: Commit**

```bash
git add autoui/highlight.go
git commit -m "feat(autoui): add drawHighlightsCallback for integrate registration"
```

---

### Task 3: Register callback in autoui.Register()

**Files:**
- Modify: `autoui/register.go`
- Modify: `autoui/register_test.go`

- [ ] **Step 1: Write the failing test for callback registration**

Check existing `autoui/register_test.go` and add test for integrate callback registration:

```go
// Add to autoui/register_test.go if integration tests exist,
// or verify through integration_test.go that callback is properly registered.
```

Note: The callback registration happens inside RegisterWithPrefix. Since integrate.drawHighlightsFunc is a package-level variable, we can test it indirectly. For now, rely on integration tests.

- [ ] **Step 2: Add import and registration call to register.go**

Modify `autoui/register.go`:

Add import:
```go
import (
	"sync"

	"github.com/ebitenui/ebitenui"
	"github.com/s3cy/autoebiten"
	"github.com/s3cy/autoebiten/integrate"
)
```

Add registration in `RegisterWithPrefix` after storing UI reference:

```go
func RegisterWithPrefix(ui *ebitenui.UI, prefix string) {
	if ui == nil {
		panic("autoui.RegisterWithPrefix: UI cannot be nil")
	}

	// Store UI reference
	uiMu.Lock()
	uiReference = ui
	uiMu.Unlock()

	// Register highlight callback for patch method
	integrate.RegisterDrawHighlights(drawHighlightsCallback)

	// Register all command handlers
	registerCommands(prefix)
}
```

- [ ] **Step 3: Run tests to verify nothing is broken**

Run: `go test ./autoui/... -v`
Expected: All existing tests PASS

- [ ] **Step 4: Commit**

```bash
git add autoui/register.go
git commit -m "feat(autoui): register drawHighlightsCallback in RegisterWithPrefix"
```

---

### Task 4: Update ebiten.patch

**Files:**
- Modify: `ebiten.patch`

- [ ] **Step 1: Update patch to use AfterDraw instead of Capture**

Modify `ebiten.patch`, change line 33:

From:
```diff
+	integrate.Capture(g.offscreen)
```

To:
```diff
+	integrate.AfterDraw(g.offscreen)
```

The full section should look like:
```diff
 func (g *gameForUI) DrawOffscreen() error {
 	g.game.Draw(g.offscreen)
+	integrate.AfterDraw(g.offscreen)
 	if err := g.imageDumper.dump(g.offscreen, g.transparent); err != nil {
 		return err
 	}
```

- [ ] **Step 2: Commit**

```bash
git add ebiten.patch
git commit -m "refactor(patch): use integrate.AfterDraw instead of Capture"
```

---

### Task 5: Verify no import cycles and run full tests

**Files:**
- No file changes

- [ ] **Step 1: Build all packages to verify no import cycles**

Run: `go build ./...`
Expected: SUCCESS, no cycle errors

- [ ] **Step 2: Run all tests**

Run: `go test ./... -v`
Expected: All tests PASS

- [ ] **Step 3: Verify race detector**

Run: `go test -race ./...`
Expected: All tests PASS, no race conditions

---

### Task 6: Update documentation

**Files:**
- Modify: `README.md` (optional)

- [ ] **Step 1: Update README.md if needed**

The README already mentions `autoebiten.Capture(screen)` for library method. For patch method, the new AfterDraw is automatic. Consider adding note about extensibility:

In README.md, near the Capture section, note that AfterDraw handles both screenshots and any registered post-draw callbacks (like autoui highlights).

- [ ] **Step 2: Commit (if changes made)**

```bash
git add README.md
git commit -m "docs: update README with AfterDraw extensibility note"
```