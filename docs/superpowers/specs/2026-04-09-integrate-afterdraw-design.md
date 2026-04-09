---
name: integrate-afterdraw
description: Add integrate.AfterDraw() to handle all post-draw operations, enabling DrawHighlights in patch method without import cycles
type: project
---

# integrate.AfterDraw() Design

**Why:** DrawHighlights() in autoui package imports ebiten, creating an import cycle when called from patched ebiten. Need a mechanism to inject highlight drawing into the patch flow without cycles.

**How to apply:** Use callback registry pattern - integrate package exposes registration function, autoui registers its callback during Register() call.

## Architecture

```
Import Graph (No Cycles):
  ebiten → integrate (patch injects AfterDraw call)
  autoui → integrate (registers callback via RegisterDrawHighlights)
  autoui → ebiten (for vector.StrokeRect, direct DrawHighlights)
  integrate → image (stdlib only)
```

## Changes

### integrate/integrate.go

Replace `Capture()` with `AfterDraw()` that handles all post-draw operations:

```go
var drawHighlightsFunc func(screen image.Image)

// RegisterDrawHighlights registers a callback for drawing highlight overlays.
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

// Capture is deprecated - use AfterDraw instead.
// Kept for backward compatibility with library method users.
func Capture(screen image.Image) {
    server.ProcessScreenshots(screen)
}
```

### ebiten.patch

Simplify `gameforui.go` to single AfterDraw call:

```diff
 func (g *gameForUI) DrawOffscreen() error {
     g.game.Draw(g.offscreen)
-    integrate.Capture(g.offscreen)
+    integrate.AfterDraw(g.offscreen)
     if err := g.imageDumper.dump(g.offscreen, g.transparent); err != nil {
         return err
     }
```

### autoui/highlight.go

Add callback implementation (internal, not exported):

```go
// drawHighlightsCallback is the callback registered with integrate.
// It type asserts image.Image to *ebiten.Image for vector drawing.
func drawHighlightsCallback(screen image.Image) {
    if ebiScreen, ok := screen.(*ebiten.Image); ok {
        globalHighlightManager.draw(ebiScreen)
    }
}
```

Keep `DrawHighlights(screen *ebiten.Image)` public for library method users who call it directly in their game's Draw() method.

### autoui/register.go

Register callback during Register() call:

```go
func Register(ui *ebitenui.UI) {
    integrate.RegisterDrawHighlights(drawHighlightsCallback)
    // ... existing registration logic
}

func RegisterWithPrefix(ui *ebitenui.UI, prefix string) {
    integrate.RegisterDrawHighlights(drawHighlightsCallback)
    // ... existing registration logic
}
```

## Usage

**Patch method:** Highlights automatically drawn when autoui.Register() is called:
```go
// In game initialization
autoui.Register(ui)  // Registers highlight callback with integrate

// In ebiten patch (automatic)
integrate.AfterDraw(screen)  // Calls ProcessScreenshots + drawHighlightsCallback
```

**Library method:** Call DrawHighlights directly:
```go
func (g *Game) Draw(screen *ebiten.Image) {
    g.ui.Draw(screen)
    autoui.DrawHighlights(screen)  // Direct call, no callback needed
}
```

## Extensibility

Future post-draw features can be added to AfterDraw() without modifying the patch:
- Debug overlays
- Performance metrics display
- Additional visual debugging tools

Pattern can extend to other lifecycle hooks:
- `integrate.BeforeUpdate()`
- `integrate.AfterUpdate()`

## Implementation Checklist

1. Update integrate/integrate.go: add RegisterDrawHighlights, AfterDraw, keep Capture for compat
2. Update ebiten.patch: change Capture to AfterDraw
3. Update autoui/highlight.go: add drawHighlightsCallback
4. Update autoui/register.go: call RegisterDrawHighlights in Register/RegisterWithPrefix
5. Update tests if needed
6. Verify no import cycles with `go build ./...`