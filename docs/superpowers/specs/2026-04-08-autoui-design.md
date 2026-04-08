# autoui Design Document

**Date**: 2026-04-08  
**Package**: `autoui` - EbitenUI automation helper for autoebiten

## Overview

`autoui` is a helper package that enables autoebiten to interact with ebitenui widget trees. It provides CLI-accessible commands for inspecting UI state, finding widgets by attributes, invoking widget methods, and visual debugging.

This package bridges the gap between autoebiten's RPC-based automation and ebitenui's widget hierarchy, enabling LLM-assisted E2E testing without hardcoded coordinates.

## Architecture

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   CLI/User  │────▶│ autoebiten  │────▶│   autoui    │
│             │     │   Register  │     │   Commands  │
└─────────────┘     └─────────────┘     └──────┬──────┘
                                               │
                           ┌───────────────────┼───────────────────┐
                           │                   │                   │
                           ▼                   ▼                   ▼
                    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
                    │ Tree Walker │    │   Finder    │    │   Caller    │
                    │  (traverse) │    │  (search)   │    │  (reflect)  │
                    └─────────────┘    └─────────────┘    └─────────────┘
                           │                   │                   │
                           └───────────────────┼───────────────────┘
                                               │
                                               ▼
                                        ┌─────────────┐
                                        │  ebitenui   │
                                        │   UI Tree   │
                                        └─────────────┘
```

### Components

1. **Registration**: Registers commands with autoebiten during game init
2. **Tree Walker**: Traverses widget hierarchy, extracts widget info
3. **Finder**: Searches widgets by attribute patterns
4. **Caller**: Invokes widget methods via reflection
5. **Highlighter**: Visual feedback for widget location

## Data Structures

### WidgetNode

Represents a widget in the tree. Used for XML marshaling.

```go
type WidgetNode struct {
    XMLName  xml.Name     `xml:"-"`              // Dynamic based on widget type (Button, Container, etc.)
    X        int          `xml:"x,attr"`
    Y        int          `xml:"y,attr"`
    Width    int          `xml:"width,attr"`
    Height   int          `xml:"height,attr"`
    Visible  bool         `xml:"visible,attr"`
    Disabled bool         `xml:"disabled,attr"`
    // Dynamic attributes from custom_data via reflection
    // Widget-specific state (text, state, value, etc.)
    Children []WidgetNode `xml:"Widget"`         // Nested child widgets
}

// MarshalXML implements custom marshaling to set element name from widget type
func (n WidgetNode) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    start.Name.Local = n.XMLName.Local
    // ... encode attributes from struct fields and dynamic maps
    return e.EncodeElement(n, start)
}
```

### WidgetInfo

Internal representation before XML conversion.

```go
type WidgetInfo struct {
    Widget     widget.PreferredSizeLocateableWidget
    Type       string
    Rect       image.Rectangle
    Visible    bool
    Disabled   bool
    State      map[string]string  // type-specific state
    CustomData map[string]string  // from custom_data via reflection
}
```

## XML Format

Standard XML with attributes. Widget type becomes element name (Button, Container, TextInput, etc.).

```xml
<?xml version="1.0" encoding="UTF-8"?>
<UI>
  <Container x="0" y="0" width="800" height="600" visible="true" disabled="false">
    <Button x="10" y="10" width="100" height="30" 
            visible="true" disabled="false"
            text="Start Game" state="idle"/>
    <Container x="0" y="50" width="800" height="550" visible="true">
      <TextInput x="20" y="70" width="200" height="30"
                 visible="true" disabled="false"
                 text="Hello World" cursor="5" focused="true"/>
    </Container>
  </Container>
</UI>
```

### Widget-Specific Attributes

| Widget Type | Extra Attributes |
|-------------|------------------|
| Button | `text`, `state` (idle/pressed/hover/disabled), `toggle` |
| TextInput | `text`, `cursor`, `selection_start`, `selection_end`, `focused` |
| Checkbox | `checked`, `state` |
| Slider | `value`, `min`, `max` |
| ProgressBar | `value`, `min`, `max` |
| Label | `text` |
| Container | (layout info if available) |

## Commands

All commands are registered via `autoebiten.Register()` and accessible via CLI.

### 1. `autoui.tree`

Returns full widget tree as XML.

**Request**: (none, or optional filter)  
**Response**: XML string

```bash
autoebiten custom autoui.tree
```

Optional: Filter by type
```bash
autoebiten custom autoui.tree --request "type=Button"
```

### 2. `autoui.at`

Returns widget at specific screen coordinates.

**Request**: `x,y` or JSON `{"x": 100, "y": 200}`  
**Response**: XML fragment (widget + children)

```bash
autoebiten custom autoui.at --request "100,200"
# or
autoebiten custom autoui.at --request '{"x":100,"y":200}'
```

### 3. `autoui.find`

Find widgets matching attribute pattern.

**Request**: Query string or JSON
**Response**: XML list of matching widgets

```bash
# Find by ID attribute
autoebiten custom autoui.find --request "id=start-btn"

# Find all buttons
autoebiten custom autoui.find --request "type=Button"

# Find visible buttons with specific text
autoebiten custom autoui.find --request '{"type":"Button","visible":"true","text":"Start"}'

# Find focused widget
autoebiten custom autoui.find --request "focused=true"
```

**Query syntax**:
- Simple: `key=value` (exact match)
- JSON object for multiple criteria (AND logic)
- Special keys: `type`, `x`, `y`, `width`, `height`, `visible`, `disabled`

### 4. `autoui.call`

Invoke method on widget via reflection.

**Request**: JSON `{ "target": "<query>", "method": "<name>", "args": [...] }`  
**Response**: Success/error message

```bash
# Click button by ID
autoebiten custom autoui.call --request '{"target":"id=start-btn","method":"Click"}'

# Focus a widget
autoebiten custom autoui.call --request '{"target":"id=username","method":"Focus","args":[true]}'

# Set slider value
autoebiten custom autoui.call --request '{"target":"id=volume","method":"SetCurrentValue","args":[0.5]}'
```

**Supported method signatures**:
- `func()` - No args
- `func(bool)` - Boolean
- `func(int)` - Integer
- `func(float64)` - Float
- `func(string)` - String

**Safety**: Only methods on exported types with specific signatures are invocable.

### 5. `autoui.highlight`

Visually highlight widget(s) for debugging.

**Request**: Query string (same as `find`) or `x,y` coordinates  
**Response**: Success confirmation

```bash
# Highlight by ID
autoebiten custom autoui.highlight --request "id=start-btn"

# Highlight at position
autoebiten custom autoui.highlight --request "100,200"

# Clear highlights
autoebiten custom autoui.highlight --request "clear"
```

Highlight appears as a colored border (red) around the widget for 3 seconds or until cleared.

## Custom Data Contract

Custom data attached to widgets via `WidgetOpts.CustomData()` is automatically flattened into XML attributes.

### Supported Types

**1. Simple values** (become attribute directly):
```go
widget.WidgetOpts.CustomData("my-id")
// XML: custom_data="my-id"
```

**2. Struct with xml tags**:
```go
type WidgetMeta struct {
    ID      string `xml:"id,attr"`
    Name    string `xml:"name,attr"`
    Section string `xml:"section,attr"`
}

widget.WidgetOpts.CustomData(WidgetMeta{
    ID: "start-btn", 
    Name: "Start Button",
    Section: "main",
})
```

**XML output**:
```xml
<Button x="10" y="10" ... id="start-btn" name="Start Button" section="main"/>
```

**3. Struct without tags** (field names lowercased):
```go
type SimpleMeta struct {
    ID   string
    Name string
}
// XML: id="..." name="..."
```

**4. Map[string]string**:
```go
widget.WidgetOpts.CustomData(map[string]string{
    "id": "btn1",
    "test-id": "start-button",
})
```

### Extraction Rules

1. If `custom_data` is nil → no extra attributes
2. If struct → use `xml` tags if present, else field names
3. If map[string]string → each key becomes attribute
4. If string/number/bool → becomes `custom_data` attribute
5. Nested structs → flattened with dot notation: `parent.child="value"`

## Method Invocation

The `autoui.call` command uses reflection to safely invoke widget methods.

### Method Discovery

```go
// Get widget type name for logging
widgetType := reflect.TypeOf(widget).Elem().Name()

// Find method by name
method := reflect.ValueOf(widget).MethodByName(methodName)
if !method.IsValid() {
    return fmt.Errorf("method %s not found on %s", methodName, widgetType)
}
```

### Argument Conversion

```go
func convertArg(arg interface{}, targetType reflect.Type) (reflect.Value, error) {
    switch targetType.Kind() {
    case reflect.Bool:
        return reflect.ValueOf(arg.(bool)), nil
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        return reflect.ValueOf(int64(arg.(float64))), nil
    case reflect.Float32, reflect.Float64:
        return reflect.ValueOf(arg.(float64)), nil
    case reflect.String:
        return reflect.ValueOf(arg.(string)), nil
    default:
        return reflect.Value{}, fmt.Errorf("unsupported type: %v", targetType)
    }
}
```

### Whitelist Approach

For safety, only allow methods with whitelisted signatures:
- No pointer returns (only void or simple values)
- No channel/function arguments
- Maximum 3 arguments

## Highlighting

Visual debugging aid. Implementation:

```go
type Highlight struct {
    Rect      image.Rectangle
    Color     color.Color
    ExpiresAt time.Time
}

var highlights []Highlight

func (h *HighlightManager) Draw(screen *ebiten.Image) {
    now := time.Now()
    // Remove expired
    highlights = filterExpired(highlights, now)
    // Draw remaining
    for _, h := range highlights {
        vector.StrokeRect(screen, float32(h.Rect.Min.X), float32(h.Rect.Min.Y),
            float32(h.Rect.Dx()), float32(h.Rect.Dy()), 2, h.Color, false)
    }
}
```

Integration with user's game:
```go
// In game's Draw() method, before or after UI.Draw()
autoui.DrawHighlights(screen)
```

Or auto-injected if using patch method.

## Usage Examples

### Basic Setup

```go
package main

import (
    "github.com/ebitenui/ebitenui"
    "github.com/ebitenui/ebitenui/widget"
    "github.com/s3cy/autoebiten"
    "github.com/s3cy/autoebiten/autoui"
)

func main() {
    // Create UI
    root := widget.NewContainer(
        widget.ContainerOpts.Layout(widget.NewRowLayout()),
    )
    
    // Create button with identifiable custom data
    btn := widget.NewButton(
        widget.ButtonOpts.Text("Start Game", face, colors),
        widget.WidgetOpts.CustomData(struct {
            ID string `xml:"id,attr"`
        }{ID: "start-btn"}),
    )
    root.AddChild(btn)
    
    ui := ebitenui.UI{
        Container: root,
    }
    
    // Register autoui commands
    autoui.Register(&ui)
    
    // Run game
    // ...
}
```

### E2E Test Flow

```go
func TestGameFlow(t *testing.T) {
    game := testkit.Launch(t, "./mygame")
    defer game.Shutdown()
    
    // 1. Get UI tree to understand structure
    tree := game.Custom("autoui.tree", "")
    t.Logf("UI Tree:\n%s", tree)
    
    // 2. Find button by ID
    result := game.Custom("autoui.find", "id=start-btn")
    // Parse XML to get position
    
    // 3. Click the button
    game.Custom("autoui.call", `{"target":"id=start-btn","method":"Click"}`)
    
    // 4. Verify state changed
    tree = game.Custom("autoui.tree", "")
    // Assert button state changed to "pressed" or navigate to new screen
}
```

### LLM-Assisted Testing

**Prompt**: "I need to test the main menu. What's on the screen?"

**autoui.tree output**:
```xml
<Container x="0" y="0" width="800" height="600">
  <Button x="300" y="200" width="200" height="50" 
          id="start-btn" text="Start Game" state="idle"/>
  <Button x="300" y="270" width="200" height="50" 
          id="settings-btn" text="Settings" state="idle"/>
  <Button x="300" y="340" width="200" height="50" 
          id="quit-btn" text="Quit" state="idle"/>
</Container>
```

**LLM response**: "I see a main menu with 3 buttons. Clicking Start Game..."

**Action**:
```bash
autoebiten custom autoui.call --request '{"target":"id=start-btn","method":"Click"}'
```

## API Surface

### Public Functions

```go
// Register registers all autoui commands with autoebiten.
// Must be called after UI is constructed.
func Register(ui *ebitenui.UI)

// RegisterWithPrefix registers commands with custom prefix.
// Useful to avoid conflicts with existing custom commands.
func RegisterWithPrefix(ui *ebitenui.UI, prefix string)

// DrawHighlights renders highlight overlays.
// Call in your game's Draw method.
func DrawHighlights(screen *ebiten.Image)

// SetHighlightDuration changes how long highlights persist.
// Default: 3 seconds
func SetHighlightDuration(d time.Duration)
```

### Configuration

```go
// Optional configuration before Register()
autoui.SetHighlightDuration(5 * time.Second)
autoui.Register(ui)
```

## Testing Strategy

### Unit Tests

1. **Tree Walker**: Verify all widget types are handled
2. **Custom Data Extraction**: Test struct, map, primitive types
3. **Finder**: Test query matching (exact, partial, multiple criteria)
4. **Caller**: Test reflection invocation with various signatures

### Integration Tests

1. **Full workflow**: Create UI → register → query → click → verify
2. **Error cases**: Invalid queries, missing widgets, type mismatches
3. **Performance**: Large UI trees (>1000 widgets)

### E2E Tests

1. **Real game scenario**: Menu navigation, form input, dialog interaction
2. **Screenshot comparison**: Highlighted widgets visible in screenshots

## Future Enhancements

1. **XPath queries**: `find --request "//Button[@id='start']"`
2. **CSS-style selectors**: `find --request "Button#start-btn"`
3. **Widget comparison**: Diff two tree states
4. **Recording/playback**: Record interactions as autoui.call sequences
5. **Accessibility tree**: Export compatible with screen readers

## Dependencies

- `github.com/s3cy/autoebiten` - RPC registration
- `github.com/ebitenui/ebitenui` - UI framework
- `encoding/xml` - Standard library
- `reflect` - Standard library
- `image` - Standard library

## Open Questions

1. Should `autoui.tree` include hidden widgets? (configurable?)
2. Should method invocation be async or blocking?
3. Maximum tree depth limit for XML output?
4. Support for widget state snapshots/diffs?
