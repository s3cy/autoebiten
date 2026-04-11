# autoui RadioGroup & TabBook Support Design

**Date**: 2026-04-11
**Package**: `autoui` - EbitenUI automation helper for autoebiten
**Depends on**: `2026-04-11-autoui-call-type-support-design.md` (merged)

## Overview

Extend autoui to support RadioGroup and TabBook widget automation. Both widgets hide their child elements in private fields, requiring reflection for enumeration and special handling for discovery.

**Goal: HTML-like tree structure for inspection and automation:**
```xml
<RadioGroup name="settings-group">
  <Button active="true" text="Option A" .../>
  <Button text="Option B" .../>
</RadioGroup>

<TabBook active_tab="Settings">
  <TabBookTab label="General" disabled="false">
    <!-- Tab content -->
  </TabBookTab>
  <TabBookTab label="Settings" disabled="false">
    <!-- Tab content -->
  </TabBookTab>
</TabBook>
```

## Problem Statement

### RadioGroup

**Discovery problem:**
- RadioGroup is NOT a widget (`PreferredSizeLocateableWidget`) - no `GetWidget()`, `Render()`, etc.
- It's a controller that coordinates Button/Checkbox widgets
- EbitenUI examples commonly discard the return value: `widget.NewRadioGroup(...)` (no assignment)
- Cannot be found via tree traversal

**Enumeration problem:**
- `elements []RadioGroupElement` is a private field
- No public getter for element list
- `Active()` returns `RadioGroupElement` interface (Button or Checkbox)
- `SetActive(RadioGroupElement)` requires actual element reference

### TabBook

**Discovery:**
- TabBook IS a widget - appears in tree traversal
- Can be found via `FindByQuery("type=TabBook")`

**Enumeration problem:**
- `tabs []*TabBookTab` is a private field
- No public getter for tab list
- `Tab()` returns `*TabBookTab`
- `SetTab(*TabBookTab)` requires actual tab reference
- `TabBookTab.label` is a private field (need reflection for label)

## Design Solution

### Strategy Summary

| Widget | Discovery | Enumeration | Target Syntax |
|--------|-----------|-------------|---------------|
| RadioGroup | User registration | Reflection on `elements` | `radiogroup=name` |
| TabBook | Tree traversal | Reflection on `tabs` | `type=TabBook` |

**Why different discovery:**
- RadioGroup must be registered because it never appears in widget tree
- TabBook is discoverable but needs reflection for tab enumeration

---

### RadioGroup Design

#### User Code Pattern

**Current EbitenUI pattern (discarded):**
```go
widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1, btn2, btn3))
```

**New pattern (store & register):**
```go
settingsGroup := widget.NewRadioGroup(
    widget.RadioGroupOpts.Elements(btn1, btn2, btn3),
    widget.RadioGroupOpts.ChangedHandler(func(args *widget.RadioGroupChangedEventArgs) {
        fmt.Println("Selected:", args.Active.(*widget.Button).Text().Label)
    }),
)
autoui.RegisterRadioGroup("settings-group", settingsGroup)
```

**Registration API:**
```go
// Register a RadioGroup with a unique name
func RegisterRadioGroup(name string, rg *widget.RadioGroup)

// Unregister when RadioGroup is no longer needed
func UnregisterRadioGroup(name string)

// Get registered RadioGroup by name (for proxy handlers)
func GetRadioGroup(name string) *widget.RadioGroup
```

#### Registry Storage

```go
// radioGroupRegistry stores registered RadioGroup instances.
// Access is protected by mutex for thread safety.
var radioGroupRegistry = struct {
    sync.RWMutex
    groups map[string]*widget.RadioGroup
}{
    groups: make(map[string]*widget.RadioGroup),
}
```

#### Tree Injection

RadioGroups are injected into tree output as virtual elements:

```xml
<!-- Injected from registry -->
<RadioGroup name="settings-group" active_index="0">
  <Button _addr="0x..." active="true" text="Option A" .../>
  <Button _addr="0x..." text="Option B" .../>
  <Button _addr="0x..." text="Option C" .../>
</RadioGroup>
```

**Implementation:**
- Modify `WalkTree` to append registered RadioGroups at the end of traversal
- RadioGroups appear as children of `<UI>` root element, not inside specific containers
- For each RadioGroup, create synthetic `<RadioGroup>` element with `name` attribute
- Children are actual Button/Checkbox widgets extracted via reflection
- Note: Elements may already appear in tree elsewhere (they're normal widgets); RadioGroup wrapper just shows grouping context

#### Proxy Handlers

| Handler | Description | Return Type |
|---------|-------------|-------------|
| `Elements() []WidgetInfo` | Return element list with WidgetInfo | JSON array |
| `ActiveIndex() int` | Index of active element (-1 if none) | int64 |
| `SetActiveByIndex(int)` | Set active by position index | nil |
| `ActiveLabel() string` | Label of active element | string |
| `SetActiveByLabel(string)` | Set active by matching label | nil |

**Handler implementation pattern:**
```go
func handleRadioGroupSetActiveByIndex(rg *widget.RadioGroup, args []any) (any, error) {
    if len(args) != 1 {
        return nil, fmt.Errorf("SetActiveByIndex requires 1 argument")
    }

    index := int(args[0].(float64))
    elements := getRadioGroupElements(rg) // reflection

    if index < 0 || index >= len(elements) {
        return nil, fmt.Errorf("index %d out of range", index)
    }

    rg.SetActive(elements[index])
    return nil, nil
}
```

---

### TabBook Design

#### Discovery

TabBook is discoverable via normal tree traversal:

```bash
autoui.find --request "type=TabBook"
autoui.xpath --request "//TabBook"
```

#### Tree Injection

During tree traversal, when encountering TabBook:
- Use reflection to access `tabs []*TabBookTab`
- For each tab, emit `<TabBookTab>` child element with:
  - `label` from private `label` field (reflection)
  - `disabled` from public `Disabled` field
  - Children: traverse tab's embedded Container

```xml
<TabBook active_tab="Settings" _addr="0x...">
  <TabBookTab label="General" disabled="false">
    <Button text="Save" .../>
    <TextInput text="" .../>
  </TabBookTab>
  <TabBookTab label="Settings" disabled="false">
    <Checkbox checked="true" .../>
  </TabBookTab>
</TabBook>
```

#### Proxy Handlers

| Handler | Description | Return Type |
|---------|-------------|-------------|
| `Tabs() []TabInfo` | Return tab list with label, disabled | JSON array |
| `TabIndex() int` | Index of active tab | int64 |
| `SetTabByIndex(int)` | Set active tab by position | nil |
| `TabLabel() string` | Label of active tab | string |
| `SetTabByLabel(string)` | Set tab by matching label | nil |

**TabInfo structure:**
```go
type TabInfo struct {
    Index    int    `json:"index"`
    Label    string `json:"label"`
    Disabled bool   `json:"disabled"`
}
```

---

### Reflection Utilities

Private field access using `unsafe` + `reflect`:

```go
// internal/reflection.go

// getPrivateField returns the value of a private field using unsafe.
// obj must be a pointer to the struct.
func getPrivateField(obj interface{}, fieldName string) reflect.Value {
    v := reflect.ValueOf(obj)
    if v.Kind() != reflect.Ptr {
        panic("obj must be a pointer")
    }
    v = v.Elem()

    field := v.FieldByName(fieldName)
    if !field.IsValid() {
        panic(fmt.Sprintf("field %s not found", fieldName))
    }

    // Bypass visibility using unsafe
    return reflect.NewAt(field.Type(),
        unsafe.Pointer(field.UnsafeAddr())).Elem()
}

// getRadioGroupElements returns the elements slice from a RadioGroup.
func getRadioGroupElements(rg *widget.RadioGroup) []widget.RadioGroupElement {
    field := getPrivateField(rg, "elements")
    // Convert []RadioGroupElement slice
    elements := make([]widget.RadioGroupElement, field.Len())
    for i := 0; i < field.Len(); i++ {
        elements[i] = field.Index(i).Interface().(widget.RadioGroupElement)
    }
    return elements
}

// getTabBookTabs returns the tabs slice from a TabBook.
func getTabBookTabs(tb *widget.TabBook) []*widget.TabBookTab {
    field := getPrivateField(tb, "tabs")
    tabs := make([]*widget.TabBookTab, field.Len())
    for i := 0; i < field.Len(); i++ {
        tabs[i] = field.Index(i).Interface().(*widget.TabBookTab)
    }
    return tabs
}

// getTabBookTabLabel returns the label from a TabBookTab.
func getTabBookTabLabel(tab *widget.TabBookTab) string {
    field := getPrivateField(tab, "label")
    return field.String()
}
```

---

### Call Handler Integration

Modify `handleCallCommand` to support RadioGroup target syntax:

```go
// Special target prefixes
const (
    radioGroupPrefix = "radiogroup="
    tabBookPrefix    = "tabbook="
)

func handleCallCommand(ui *ebitenui.UI) func(ctx autoebiten.CommandContext) {
    return func(ctx autoebiten.CommandContext) {
        // ... existing parsing ...

        // Check for special target prefixes
        if strings.HasPrefix(callReq.Target, radioGroupPrefix) {
            name := strings.TrimPrefix(callReq.Target, radioGroupPrefix)
            rg := GetRadioGroup(name)
            if rg == nil {
                ctx.Respond("error: RadioGroup not found: " + name)
                return
            }
            result, err := invokeRadioGroupMethod(rg, callReq.Method, callReq.Args)
            // ... response handling ...
            return
        }

        if strings.HasPrefix(callReq.Target, tabBookPrefix) {
            name := strings.TrimPrefix(callReq.Target, tabBookPrefix)
            // Future: Find TabBook by name attribute (from ae tag or registration)
            // Currently: TabBook uses normal type=TabBook discovery
            ctx.Respond("error: tabbook=name syntax not yet implemented, use type=TabBook")
            return
        }

        // Existing: normal widget lookup via FindByQuery
        // ...
    }
}
```

**RadioGroup method invocation:**
```go
func invokeRadioGroupMethod(rg *widget.RadioGroup, method string, args []any) (any, error) {
    // Check RadioGroup proxy registry
    if handler := getRadioGroupProxyHandler(method); handler != nil {
        return handler(rg, args)
    }

    // Direct method call via reflection
    // (Active, SetActive require RadioGroupElement, not useful for CLI)
    return nil, fmt.Errorf("method %s not supported for RadioGroup", method)
}
```

---

## Architecture & Data Flow

```
autoui.call Request
  {"target":"radiogroup=settings-group","method":"SetActiveByIndex","args":[1]}
        │
        ▼
  handleCallCommand
    1. Parse JSON request
    2. Detect "radiogroup=" prefix
    3. GetRadioGroup("settings-group") from registry
        │
        ▼
  invokeRadioGroupMethod
    1. Check RadioGroup proxy registry
    2. handleRadioGroupSetActiveByIndex(rg, [1])
        │
        ▼
  handleRadioGroupSetActiveByIndex
    1. getRadioGroupElements(rg) → [btn1, btn2, btn3]
    2. Validate index (0-2)
    3. rg.SetActive(elements[1])
        │
        ▼
  CallResponse
  {"success":true}
```

```
autoui.tree Request
        │
        ▼
  SnapshotTree(ui)
    1. WalkTree(ui.Container) → normal widgets
    2. Append registered RadioGroups
    3. For each TabBook: inject TabBookTab children
        │
        ▼
  MarshalWidgetTreeXML
    <UI>
      <Container>
        <Button .../>
        <RadioGroup name="settings-group">...</RadioGroup>
        <TabBook><TabBookTab>...</TabBookTab></TabBook>
      </Container>
    </UI>
```

---

## Implementation Files

| File | Purpose | Changes |
|------|---------|---------|
| `autoui/registry.go` | RadioGroup registration | NEW |
| `autoui/tree.go` | Tree traversal | MODIFY - inject RadioGroups & TabBookTab |
| `autoui/handlers.go` | Call command | MODIFY - special target prefixes |
| `autoui/proxy_radiogroup.go` | RadioGroup proxy handlers | NEW |
| `autoui/proxy_tabbook.go` | TabBook proxy handlers | NEW |
| `autoui/proxy.go` | Proxy registry | MODIFY - add RadioGroup/TabBook registries |
| `autoui/internal/reflection.go` | Private field access | NEW |
| `autoui/internal/widgetstate.go` | State extraction | MODIFY - TabBook active_tab label |
| `autoui/xml.go` | XML marshaling | MODIFY - RadioGroup & TabBookTab elements |

---

## CLI Usage Examples

### RadioGroup Operations

```bash
# Register RadioGroup in game code:
settingsGroup := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1, btn2, btn3))
autoui.RegisterRadioGroup("settings-group", settingsGroup)

# CLI: Get elements
autoui.call --request '{"target":"radiogroup=settings-group","method":"Elements"}'
# {"success":true,"result":[{"type":"Button","text":"Option A","active":true},...]}

# CLI: Set active by index
autoui.call --request '{"target":"radiogroup=settings-group","method":"SetActiveByIndex","args":[2]}'

# CLI: Get active index
autoui.call --request '{"target":"radiogroup=settings-group","method":"ActiveIndex"}'
# {"success":true,"result":2}

# CLI: Get active label
autoui.call --request '{"target":"radiogroup=settings-group","method":"ActiveLabel"}'
# {"success":true,"result":"Option C"}
```

### TabBook Operations

```bash
# CLI: Get tabs (TabBook found via tree)
autoui.call --request '{"target":"type=TabBook","method":"Tabs"}'
# {"success":true,"result":[{"index":0,"label":"General","disabled":false},...]}

# CLI: Set tab by index
autoui.call --request '{"target":"type=TabBook","method":"SetTabByIndex","args":[1]}'

# CLI: Get active tab label
autoui.call --request '{"target":"type=TabBook","method":"TabLabel"}'
# {"success":true,"result":"Settings"}

# CLI: Set tab by label
autoui.call --request '{"target":"type=TabBook","method":"SetTabByLabel","args":["Settings"]}'
```

### Tree Output

```bash
autoui.tree
```
```xml
<UI>
  <Container>
    <Button id="btn1" .../>
    <TabBook active_tab="General">
      <TabBookTab label="General" disabled="false">
        <Button text="Save" .../>
      </TabBookTab>
      <TabBookTab label="Settings" disabled="false">
        <Checkbox checked="true" .../>
      </TabBookTab>
    </TabBook>
  </Container>
  <!-- RadioGroups appended at root level -->
  <RadioGroup name="settings-group" active_index="0">
    <Button text="Option A" state="checked" .../>
    <Button text="Option B" state="unchecked" .../>
    <Button text="Option C" state="unchecked" .../>
  </RadioGroup>
</UI>
```

**Note:** RadioGroups appear at `<UI>` root level, not inside containers. Elements (Button/Checkbox) are normal widgets that may also appear elsewhere in the tree; the RadioGroup wrapper shows their grouping relationship.

---

## Testing Strategy

### Unit Tests

- `reflection_test.go`: Private field access utilities
- `proxy_radiogroup_test.go`: RadioGroup proxy handlers
- `proxy_tabbook_test.go`: TabBook proxy handlers
- `registry_test.go`: RadioGroup registration/unregistration

### Integration Tests

- Full workflow: register RadioGroup → tree output → call methods
- TabBook discovery → enumerate tabs → switch tabs
- Edge cases: empty RadioGroup, disabled tabs, missing registration

---

## Risk Assessment

### Reflection on Private Fields

**Risk:** EbitenUI internal changes could break reflection.

**Mitigation:**
- Document exact field names accessed: `elements`, `tabs`, `label`
- Add validation that field exists before access
- Panic recovery with clear error message
- Monitor EbitenUI releases for structural changes

### RadioGroup Registration Required

**Risk:** Users forget to register RadioGroup, automation fails.

**Mitigation:**
- Clear documentation emphasizing registration requirement
- Example code showing registration pattern
- Error message when RadioGroup not found: "RadioGroup 'name' not registered. Did you call autoui.RegisterRadioGroup?"

### Thread Safety

**Risk:** Concurrent access to RadioGroup registry.

**Mitigation:**
- Use `sync.RWMutex` for registry access
- Follow existing pattern from `treeRWLock`

---

## Future Enhancements

- **Auto-registration hook**: Intercept `widget.NewRadioGroup` calls globally
- **Named TabBook lookup**: `tabbook=name` target syntax via ae tags
- **TabBook registration**: Optional registration for explicit naming (like RadioGroup)
- **Checkbox state in RadioGroup**: Show `checked="true/false"` for Checkbox elements