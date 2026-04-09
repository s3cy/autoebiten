---
name: autoui-bugfixes-2026-04-09
description: Fix 5 bugs in autoui package for XPath matching, race conditions, hierarchy, slices, and button text
type: project
---

# autoui Bug Fixes Design

Date: 2026-04-09

## Overview

Fix 5 bugs in the autoui package affecting widget querying, thread safety, and data extraction.

## Bug 1: XPath Node Matching by Position Only

**File:** `autoui/xpath.go:48-94`

**Problem:** When mapping XPath results back to WidgetInfo, matching uses only type + x,y coordinates. If two widgets of the same type overlap at the same position, only one matches.

**Solution:** Add `_addr` hidden attribute (widget pointer address) for exact matching.

### Implementation

1. **Add `_addr` to WidgetNode output** (`autoui/xml.go`)
   - In `widgetInfoToNode`, add `_addr` attribute with widget pointer address
   - Format: `_addr="0x14000abc0"` (hex string)
   - Hidden attribute (internal use, not for user queries)

2. **Match by `_addr` only in XPath** (`autoui/xpath.go`)
   - `xmlNodesToWidgetInfo` extracts `_addr` from XML node
   - Find matching WidgetInfo by comparing `_addr` string with widget pointer
   - Remove position-based matching logic entirely

3. **Update WidgetInfo** (`autoui/tree.go`)
   - Add `Addr string` field to WidgetInfo struct
   - Populate in `ExtractWidgetInfo` from widget pointer address

### Code Changes

```go
// tree.go - WidgetInfo struct
type WidgetInfo struct {
    Widget     widget.PreferredSizeLocateableWidget
    Type       string
    Rect       image.Rectangle
    Visible    bool
    Disabled   bool
    State      map[string]string
    CustomData map[string]string
    Addr       string  // Widget pointer address for unique identification
}

// tree.go - ExtractWidgetInfo
func ExtractWidgetInfo(w widget.PreferredSizeLocateableWidget) WidgetInfo {
    // ... existing code ...
    info.Addr = fmt.Sprintf("0x%x", reflect.ValueOf(w).Pointer())
    return info
}

// xml.go - widgetInfoToNode
func widgetInfoToNode(info WidgetInfo) *WidgetNode {
    // ... existing code ...
    node.Attrs["_addr"] = info.Addr
    return node
}

// xpath.go - xmlNodesToWidgetInfo
func xmlNodesToWidgetInfo(nodes []*xmlquery.Node, widgets []WidgetInfo) []WidgetInfo {
    result := make([]WidgetInfo, 0, len(nodes))
    for _, node := range nodes {
        addr := node.SelectAttr("_addr")
        if addr == "" {
            continue
        }
        for _, w := range widgets {
            if w.Addr == addr {
                result = append(result, w)
                break
            }
        }
    }
    return result
}
```

### Documentation Update

Add note in `docs/autoui.md`:

> **Note:** `_addr` is an internal attribute for widget identification. It changes between runs, so don't hard-code it in tests. Instead, read `_addr` dynamically from `autoui.tree` output and use it for subsequent queries if needed.

---

## Bug 2: Race Condition on Widget Tree

**File:** `autoui/handlers.go`, `autoui/register.go`

**Problem:** `uiReference` is protected by mutex, but `WalkTree` traverses the widget tree without locking. If the UI modifies the tree during traversal, race conditions occur.

**Solution:** Add `SnapshotTree()` that copies widget tree under RLock before traversal.

### Implementation

1. **Add SnapshotTree function** (`autoui/tree.go`)
   - Acquire RLock on uiMu
   - Call WalkTree with container reference
   - Release RLock
   - Return copied WidgetInfo slice

2. **Update handlers** (`autoui/handlers.go`)
   - Replace all `WalkTree(ui.Container)` calls with `SnapshotTree(ui)`
   - Affects: handleTreeCommand, handleAtCommand, handleFindCommand, handleXPathCommand, handleCallCommand, handleHighlightCommand

3. **Add UI getter for SnapshotTree** (`autoui/register.go`)
   - Add `GetUIContainer()` that returns container under RLock
   - Or pass ui to SnapshotTree directly

### Code Changes

```go
// tree.go
func SnapshotTree(ui *ebitenui.UI) []WidgetInfo {
    if ui == nil || ui.Container == nil {
        return nil
    }
    return WalkTree(ui.Container)
}

// handlers.go - example change
func handleTreeCommand(ui *ebitenui.UI) func(ctx autoebiten.CommandContext) {
    return func(ctx autoebiten.CommandContext) {
        if ui == nil {
            ctx.Respond("error: UI not registered")
            return
        }
        widgets := SnapshotTree(ui)  // Changed from WalkTree(ui.Container)
        // ... rest of handler
    }
}
```

**Note:** `uiMu` only protects `uiReference` (the pointer to the UI). The UI container's internal state (children list) has no mutex in ebitenui. This means:

- RLock prevents `uiReference` from being changed during traversal
- But `container.AddChild()` could still race with `container.Children()` if called from another goroutine

**Documented limitation:** autoui commands should be called from the main thread where Ebiten runs (typical game engine pattern). The RLock provides partial protection against concurrent `Register` calls, but full thread safety requires single-threaded UI access.

---

## Bug 3: Orphan Widgets Lose Hierarchy

**File:** `autoui/xml.go:121-146` in `buildNodeTree()`

**Problem:** When `find`/`xpath` returns filtered results without parent containers, orphan widgets are added directly to the root, losing hierarchy.

**Solution:** Follow XPath convention — return matched widgets directly without hierarchy reconstruction.

### Implementation

1. **Change output format** (`autoui/xml.go`)
   - Rename `MarshalWidgetTreeXML` to `MarshalWidgetsXML`
   - Single widget: return `<Button .../>` (no wrapper)
   - Multiple widgets: return elements concatenated
   - No hierarchy reconstruction

2. **Update handlers** (`autoui/handlers.go`)
   - Use `MarshalWidgetsXML` for find/xpath handlers
   - Adjust output logic for single vs multiple matches

3. **Keep tree command unchanged**
   - `handleTreeCommand` still uses `MarshalWidgetTreeXML` for full hierarchy
   - Rename old function to `MarshalWidgetTreeXML` for tree command only
   - Add new `MarshalWidgetsXML` for filtered results

### Code Changes

```go
// xml.go - new function
func MarshalWidgetsXML(widgets []WidgetInfo) ([]byte, error) {
    if len(widgets) == 0 {
        return nil, nil
    }

    var buf bytes.Buffer
    for _, info := range widgets {
        node := widgetInfoToNode(info)
        data, err := xml.Marshal(node)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal widget: %w", err)
        }
        buf.Write(data)
    }
    return buf.Bytes(), nil
}

// handlers.go - handleFindCommand
func handleFindCommand(ui *ebitenui.UI) func(ctx autoebiten.CommandContext) {
    return func(ctx autoebiten.CommandContext) {
        // ... existing logic ...
        matching := FindByQuery(widgets, request)
        if len(matching) == 0 {
            ctx.Respond("error: no widgets found matching query")
            return
        }
        xmlData, err := MarshalWidgetsXML(matching)  // Changed
        ctx.Respond(string(xmlData))
    }
}
```

### Documentation Update

Add note in `docs/autoui.md`:

> `autoui.find` and `autoui.xpath` return matched widgets directly without wrapper elements. For full hierarchy, use `autoui.tree`.

---

## Bug 4: No Slice/Array Support in CustomData

**File:** `autoui/internal/customdata.go:30-65`

**Problem:** Slices and arrays fall through to default case and get stringified as `[item1 item2]`. Not queryable via XPath.

**Solution:** Flatten slices with indexed keys.

### Implementation

1. **Add slice/array handling** (`autoui/internal/customdata.go`)
   - Add case for `reflect.Slice` and `reflect.Array`
   - Iterate elements, output with indexed keys: `tags.0`, `tags.1`, etc.

2. **Handle nested slices**
   - If element is a slice, recursively flatten
   - Keys become `prefix.0.0`, `prefix.0.1`, etc.

### Code Changes

```go
// customdata.go - add case in ExtractCustomData switch
case reflect.Slice, reflect.Array:
    for i := 0; i < v.Len(); i++ {
        elem := v.Index(i)
        key := fmt.Sprintf("%d", i)
        if prefix != "" {
            key = prefix + "." + key
        }
        // Handle nested slices/structs
        if elem.Kind() == reflect.Slice || elem.Kind() == reflect.Array {
            extractSliceElements(elem, result, key)
        } else if elem.Kind() == reflect.Struct {
            extractStructFields(elem, result, key)
        } else {
            result[key] = valueToString(elem)
        }
    }

// customdata.go - helper function
func extractSliceElements(v reflect.Value, result map[string]string, prefix string) {
    for i := 0; i < v.Len(); i++ {
        elem := v.Index(i)
        key := fmt.Sprintf("%s.%d", prefix, i)
        if elem.Kind() == reflect.Struct {
            extractStructFields(elem, result, key)
        } else if elem.Kind() == reflect.Slice || elem.Kind() == reflect.Array {
            extractSliceElements(elem, result, key)
        } else {
            result[key] = valueToString(elem)
        }
    }
}
```

### Test Cases

Add tests for:
- `[]string{"fire", "ice"}` → `{"0": "fire", "1": "ice"}`
- `[]int{1, 2, 3}` → `{"0": "1", "1": "2", "2": "3"}`
- Nested slices: `[][]string{{"a", "b"}, {"c"}}` → `{"0.0": "a", "0.1": "b", "1.0": "c"}`

### Documentation Update

Add example in `docs/autoui.md`:

```go
type PlayerData struct {
    Tags []string `ae:"tags"`
}

// Output XML: tags.0="fire" tags.1="ice"
// XPath: //Button[@tags.0='fire']
```

---

## Bug 5: Button Text Not Extracted Without Fonts

**File:** `autoui/internal/widgetstate.go:47-67` in `extractButtonState()`

**Problem:** `btn.Text()` returns nil when no font is set, so text attribute is missing. The example game can't use `text=Submit` queries.

**Solution:** Document the limitation. No code changes.

### Implementation

1. **Update documentation** (`docs/autoui.md`)
   - Add note in Button widget-specific attributes section
   - Explain that `text` requires `ButtonOpts.Text()` with font

2. **No code changes**
   - Current behavior is correct: nil text means no text
   - Users should use `id` from CustomData for queries

### Documentation Update

Add in `docs/autoui.md` under "Widget-specific attributes":

> | Widget | Attributes |
> |--------|------------|
> | Button | text (requires font), state, toggle, focused |
>
> **Note:** Button `text` attribute requires font setup via `ButtonOpts.Text()`. For buttons without fonts, use `id` from CustomData for queries:
>
> ```go
> btn.GetWidget().CustomData = map[string]string{"id": "submit-btn"}
> ```
>
> ```bash
> autoebiten custom autoui.find --request "id=submit-btn"
> ```

---

## Test Updates

### Existing Tests to Update

1. `autoui/xpath_test.go`
   - Tests that rely on position matching need update
   - Add tests for `_addr` based matching

2. `autoui/xml_test.go`
   - Add test for `_addr` in output
   - Add test for `MarshalWidgetsXML` format

3. `autoui/internal/customdata_test.go`
   - Add slice/array tests

### New Tests to Add

1. Test `_addr` uniqueness in XPath matching
2. Test overlapping widgets with same position
3. Test `MarshalWidgetsXML` output format
4. Test slice flattening in CustomData

---

## Backward Compatibility

### Breaking Changes

1. **XPath output format** — Users parsing `autoui.find`/`autoui.xpath` output may need to update parsers if they expect `<UI>` wrapper.

2. **`_addr` attribute** — New attribute appears in all widget XML output. May affect users who parse all attributes.

### Non-Breaking Changes

1. Slice support — adds new key format, doesn't break existing
2. Race condition fix — internal change, user-facing behavior unchanged
3. Button text documentation — no code change

---

## Verification Steps

After implementation:

```bash
# Run all tests
go test -race ./autoui/...

# Build and test example
cd examples/autoui
go build -o autoui_demo
autoebiten launch -- ./autoui_demo &
sleep 2

# Test commands
autoebiten custom autoui.tree
autoebiten custom autoui.find --request "type=Button"
autoebiten custom autoui.xpath --request "//Button"

# Check _addr appears in output
autoebiten custom autoui.tree | grep "_addr"

# Test slice support (if example has slices)
autoebiten custom autoui.find --request "player_id=p001"
```