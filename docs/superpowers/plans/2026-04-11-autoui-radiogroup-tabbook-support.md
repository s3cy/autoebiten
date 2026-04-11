# autoui RadioGroup & TabBook Support Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Extend autoui to support RadioGroup and TabBook widget automation with HTML-like tree output and index/label-based selection.

**Architecture:** RadioGroup requires user registration (not discoverable in tree); TabBook uses tree traversal with reflection for private `tabs` field. Both use reflection for private field access and proxy handlers for automation.

**Tech Stack:** Go, EbitenUI, unsafe.Pointer for reflection, sync.RWMutex for thread safety

---

## File Structure

| File | Purpose | Status |
|------|---------|--------|
| `autoui/internal/reflection.go` | Private field access utilities | NEW |
| `autoui/registry.go` | RadioGroup registration API | MODIFY (add RadioGroup registry) |
| `autoui/proxy_tabbook.go` | TabBook proxy handlers | NEW |
| `autoui/proxy_radiogroup.go` | RadioGroup proxy handlers | NEW |
| `autoui/proxy.go` | Add TabBook handlers to registry | MODIFY |
| `autoui/tree.go` | Inject RadioGroups & TabBookTab children | MODIFY |
| `autoui/xml.go` | Handle synthetic RadioGroup nodes | MODIFY |
| `autoui/handlers.go` | Special `radiogroup=` target prefix | MODIFY |
| `autoui/internal/widgetstate.go` | TabBook active_tab label | MODIFY |

---

### Task 1: Reflection Utilities

**Files:**
- Create: `autoui/internal/reflection.go`
- Create: `autoui/internal/reflection_test.go`

**Goal:** Create safe private field access utilities for RadioGroup and TabBook.

- [ ] **Step 1: Write the failing test for getRadioGroupElements**

```go
// autoui/internal/reflection_test.go
package internal_test

import (
	"testing"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/s3cy/autoebiten/autoui/internal"
)

func TestGetRadioGroupElements(t *testing.T) {
	// Create buttons for radio group
	btn1 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	btn2 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	btn3 := widget.NewButton(widget.ButtonOpts.ToggleMode())

	// Create radio group
	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1, btn2, btn3))

	// Extract elements via reflection
	elements := internal.GetRadioGroupElements(rg)

	if len(elements) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(elements))
	}

	// Verify elements are the same buttons
	if elements[0] != btn1 {
		t.Error("element 0 is not btn1")
	}
	if elements[1] != btn2 {
		t.Error("element 1 is not btn2")
	}
	if elements[2] != btn3 {
		t.Error("element 2 is not btn3")
	}
}

func TestGetRadioGroupElementsEmpty(t *testing.T) {
	rg := widget.NewRadioGroup()
	elements := internal.GetRadioGroupElements(rg)

	if len(elements) != 0 {
		t.Errorf("expected 0 elements for empty RadioGroup, got %d", len(elements))
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./autoui/internal -run TestGetRadioGroupElements -v`
Expected: FAIL with "undefined: internal.GetRadioGroupElements"

- [ ] **Step 3: Write minimal implementation for reflection.go**

```go
// autoui/internal/reflection.go
package internal

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/ebitenui/ebitenui/widget"
)

// GetRadioGroupElements returns the elements slice from a RadioGroup using reflection.
func GetRadioGroupElements(rg *widget.RadioGroup) []widget.RadioGroupElement {
	field := getPrivateField(rg, "elements")
	elements := make([]widget.RadioGroupElement, field.Len())
	for i := 0; i < field.Len(); i++ {
		elements[i] = field.Index(i).Interface().(widget.RadioGroupElement)
	}
	return elements
}

// GetTabBookTabs returns the tabs slice from a TabBook using reflection.
func GetTabBookTabs(tb *widget.TabBook) []*widget.TabBookTab {
	field := getPrivateField(tb, "tabs")
	tabs := make([]*widget.TabBookTab, field.Len())
	for i := 0; i < field.Len(); i++ {
		tabs[i] = field.Index(i).Interface().(*widget.TabBookTab)
	}
	return tabs
}

// GetTabBookTabLabel returns the label string from a TabBookTab using reflection.
func GetTabBookTabLabel(tab *widget.TabBookTab) string {
	field := getPrivateField(tab, "label")
	return field.String()
}

// getPrivateField returns a reflect.Value for a private field using unsafe.
// obj must be a pointer to the struct containing the field.
func getPrivateField(obj interface{}, fieldName string) reflect.Value {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("getPrivateField: obj must be a pointer, got %T", obj))
	}
	v = v.Elem()

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		panic(fmt.Sprintf("getPrivateField: field '%s' not found in %T", fieldName, obj))
	}

	// Use unsafe to bypass visibility
	// Create a new accessible value from the unsafe pointer
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./autoui/internal -run TestGetRadioGroupElements -v`
Expected: PASS

- [ ] **Step 5: Write failing test for GetTabBookTabs**

```go
// Add to autoui/internal/reflection_test.go

func TestGetTabBookTabs(t *testing.T) {
	// Create TabBookTab instances (minimal setup)
	tab1 := widget.NewTabBookTab(widget.TabBookTabOpts.Label("General"))
	tab2 := widget.NewTabBookTab(widget.TabBookTabOpts.Label("Settings"))

	// Note: TabBook requires full validation which needs theme setup
	// For this test, we create a TabBook and manually set tabs field via reflection
	// This tests the getter function, not the full TabBook initialization

	tb := widget.NewTabBook()
	// Set tabs via reflection for testing
	setPrivateField(tb, "tabs", []*widget.TabBookTab{tab1, tab2})

	tabs := internal.GetTabBookTabs(tb)

	if len(tabs) != 2 {
		t.Fatalf("expected 2 tabs, got %d", len(tabs))
	}

	if tabs[0] != tab1 {
		t.Error("tab 0 is not tab1")
	}
	if tabs[1] != tab2 {
		t.Error("tab 1 is not tab2")
	}
}

func TestGetTabBookTabLabel(t *testing.T) {
	tab := widget.NewTabBookTab(widget.TabBookTabOpts.Label("Settings"))
	label := internal.GetTabBookTabLabel(tab)

	if label != "Settings" {
		t.Errorf("expected label 'Settings', got '%s'", label)
	}
}

// Helper for testing: setPrivateField
func setPrivateField(obj interface{}, fieldName string, value interface{}) {
	v := reflect.ValueOf(obj).Elem()
	field := v.FieldByName(fieldName)
	fieldPtr := unsafe.Pointer(field.UnsafeAddr())
	realField := reflect.NewAt(field.Type(), fieldPtr).Elem()
	realField.Set(reflect.ValueOf(value))
}
```

- [ ] **Step 6: Run test to verify it fails**

Run: `go test ./autoui/internal -run TestGetTabBookTabs -v`
Expected: FAIL with "undefined: internal.GetTabBookTabs"

- [ ] **Step 7: Verify tests pass (GetTabBookTabs already implemented in Step 3)**

Run: `go test ./autoui/internal -run TestGetTabBook -v`
Expected: PASS

- [ ] **Step 8: Commit reflection utilities**

```bash
git add autoui/internal/reflection.go autoui/internal/reflection_test.go
git commit -m "feat(autoui): add reflection utilities for private field access

Add getPrivateField, GetRadioGroupElements, GetTabBookTabs, and
GetTabBookTabLabel for accessing private fields in EbitenUI widgets.
These enable RadioGroup and TabBook enumeration via unsafe reflection."
```

---

### Task 2: RadioGroup Registry

**Files:**
- Modify: `autoui/registry.go`
- Create: `autoui/registry_test.go`

**Goal:** Add RadioGroup registration API with thread-safe storage.

- [ ] **Step 1: Write the failing test for RegisterRadioGroup**

```go
// autoui/registry_test.go
package autoui_test

import (
	"testing"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/s3cy/autoebiten/autoui"
)

func TestRegisterRadioGroup(t *testing.T) {
	btn1 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	btn2 := widget.NewButton(widget.ButtonOpts.ToggleMode())

	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1, btn2))

	autoui.RegisterRadioGroup("test-group", rg)

	// Verify retrieval
	retrieved := autoui.GetRadioGroup("test-group")
	if retrieved != rg {
		t.Error("retrieved RadioGroup does not match registered")
	}

	// Clean up
	autoui.UnregisterRadioGroup("test-group")
	retrieved = autoui.GetRadioGroup("test-group")
	if retrieved != nil {
		t.Error("RadioGroup still exists after unregister")
	}
}

func TestGetRadioGroupNotFound(t *testing.T) {
	retrieved := autoui.GetRadioGroup("nonexistent")
	if retrieved != nil {
		t.Error("expected nil for nonexistent RadioGroup")
	}
}

func TestRegisterRadioGroupReplace(t *testing.T) {
	btn1 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	rg1 := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1))
	rg2 := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1))

	autoui.RegisterRadioGroup("replace-group", rg1)
	autoui.RegisterRadioGroup("replace-group", rg2)

	retrieved := autoui.GetRadioGroup("replace-group")
	if retrieved != rg2 {
		t.Error("RadioGroup was not replaced")
	}

	autoui.UnregisterRadioGroup("replace-group")
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./autoui -run TestRegisterRadioGroup -v`
Expected: FAIL with "undefined: autoui.RegisterRadioGroup"

- [ ] **Step 3: Add RadioGroup registry to registry.go**

```go
// Add to autoui/registry.go (after existing uiReference section)

// radioGroupRegistry stores registered RadioGroup instances.
// Access is protected by mutex for thread safety.
var radioGroupRegistry = struct {
	sync.RWMutex
	groups map[string]*widget.RadioGroup
}{
	groups: make(map[string]*widget.RadioGroup),
}

// RegisterRadioGroup registers a RadioGroup with a unique name.
// If a RadioGroup with the same name exists, it will be replaced.
// Use this to make RadioGroups discoverable via autoui.call and autoui.tree.
func RegisterRadioGroup(name string, rg *widget.RadioGroup) {
	radioGroupRegistry.Lock()
	radioGroupRegistry.groups[name] = rg
	radioGroupRegistry.Unlock()
}

// UnregisterRadioGroup removes a registered RadioGroup by name.
// Call this when the RadioGroup is no longer needed.
func UnregisterRadioGroup(name string) {
	radioGroupRegistry.Lock()
	delete(radioGroupRegistry.groups, name)
	radioGroupRegistry.Unlock()
}

// GetRadioGroup returns a registered RadioGroup by name, or nil if not found.
func GetRadioGroup(name string) *widget.RadioGroup {
	radioGroupRegistry.RLock()
	defer radioGroupRegistry.RUnlock()
	return radioGroupRegistry.groups[name]
}

// GetRegisteredRadioGroups returns all registered RadioGroup names.
// Used by tree traversal to inject RadioGroups into output.
func GetRegisteredRadioGroups() []string {
	radioGroupRegistry.RLock()
	defer radioGroupRegistry.RUnlock()

	names := make([]string, 0, len(radioGroupRegistry.groups))
	for name := range radioGroupRegistry.groups {
		names = append(names, name)
	}
	return names
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./autoui -run TestRegisterRadioGroup -v`
Expected: PASS

- [ ] **Step 5: Commit RadioGroup registry**

```bash
git add autoui/registry.go autoui/registry_test.go
git commit -m "feat(autoui): add RadioGroup registration API

Add RegisterRadioGroup, UnregisterRadioGroup, GetRadioGroup, and
GetRegisteredRadioGroups for managing RadioGroup instances.
RadioGroups require registration since they are not discoverable
in the widget tree."
```

---

### Task 3: TabBook Proxy Handlers

**Files:**
- Create: `autoui/proxy_tabbook.go`
- Create: `autoui/proxy_tabbook_test.go`
- Modify: `autoui/proxy.go`

**Goal:** Add proxy handlers for TabBook: Tabs, TabIndex, SetTabByIndex, TabLabel, SetTabByLabel.

- [ ] **Step 1: Write the failing test for handleTabs**

```go
// autoui/proxy_tabbook_test.go
package autoui_test

import (
	"testing"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/s3cy/autoebiten/autoui"
	"github.com/ebitenui/ebitenui/image"
	"image/color"
	"github.com/hajimehoshi/ebiten/v2"
)

func TestHandleTabs(t *testing.T) {
	// Create minimal TabBook with tabs
	tab1 := widget.NewTabBookTab(widget.TabBookTabOpts.Label("General"))
	tab2 := widget.NewTabBookTab(widget.TabBookTabOpts.Label("Settings"))
	tab2.Disabled = true

	// Create TabBook (requires button image for validation)
	buttonImage := &widget.ButtonImage{
		Idle:    createTestNineSlice(100, 30, color.Gray{150}),
		Pressed: createTestNineSlice(100, 30, color.Gray{100}),
	}

	tb := widget.NewTabBook(
		widget.TabBookOpts.Tabs(tab1, tab2),
		widget.TabBookOpts.TabButtonImage(buttonImage),
	)
	tb.Validate()

	// Call proxy handler
	result, err := autoui.InvokeMethod(tb, "Tabs", nil)
	if err != nil {
		t.Fatalf("Tabs failed: %v", err)
	}

	tabs := result.([]autoui.TabInfo)
	if len(tabs) != 2 {
		t.Fatalf("expected 2 tabs, got %d", len(tabs))
	}

	if tabs[0].Label != "General" {
		t.Errorf("expected tab 0 label 'General', got '%s'", tabs[0].Label)
	}
	if tabs[1].Disabled != true {
		t.Errorf("expected tab 1 to be disabled")
	}
}

func createTestNineSlice(w, h int, c color.Color) *image.NineSlice {
	img := ebiten.NewImage(w, h)
	img.Fill(c)
	return image.NewNineSliceSimple(img, 0, 0)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./autoui -run TestHandleTabs -v`
Expected: FAIL with "method 'Tabs' not found" or proxy not registered

- [ ] **Step 3: Create proxy_tabbook.go**

```go
// autoui/proxy_tabbook.go
package autoui

import (
	"fmt"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/s3cy/autoebiten/autoui/internal"
)

// TabInfo represents information about a TabBook tab.
type TabInfo struct {
	Index    int    `json:"index"`
	Label    string `json:"label"`
	Disabled bool   `json:"disabled"`
}

// handleTabs returns the list of tabs in a TabBook.
func handleTabs(w widget.PreferredSizeLocateableWidget, args []any) (any, error) {
	tb, ok := w.(*widget.TabBook)
	if !ok {
		return nil, fmt.Errorf("Tabs requires widget of type *TabBook, got %T", w)
	}

	tabs := internal.GetTabBookTabs(tb)
	result := make([]TabInfo, len(tabs))

	for i, tab := range tabs {
		result[i] = TabInfo{
			Index:    i,
			Label:    internal.GetTabBookTabLabel(tab),
			Disabled: tab.Disabled,
		}
	}

	return result, nil
}

// handleTabIndex returns the index of the currently active tab.
func handleTabIndex(w widget.PreferredSizeLocateableWidget, args []any) (any, error) {
	tb, ok := w.(*widget.TabBook)
	if !ok {
		return nil, fmt.Errorf("TabIndex requires widget of type *TabBook, got %T", w)
	}

	activeTab := tb.Tab()
	if activeTab == nil {
		return int64(-1), nil
	}

	tabs := internal.GetTabBookTabs(tb)
	for i, tab := range tabs {
		if tab == activeTab {
			return int64(i), nil
		}
	}

	return int64(-1), nil
}

// handleSetTabByIndex sets the active tab by index.
func handleSetTabByIndex(w widget.PreferredSizeLocateableWidget, args []any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("SetTabByIndex requires 1 argument (index)")
	}

	index := int(args[0].(float64))

	tb, ok := w.(*widget.TabBook)
	if !ok {
		return nil, fmt.Errorf("SetTabByIndex requires widget of type *TabBook, got %T", w)
	}

	tabs := internal.GetTabBookTabs(tb)
	if index < 0 || index >= len(tabs) {
		return nil, fmt.Errorf("index %d out of range (0-%d)", index, len(tabs)-1)
	}

	tb.SetTab(tabs[index])
	return nil, nil
}

// handleTabLabel returns the label of the currently active tab.
func handleTabLabel(w widget.PreferredSizeLocateableWidget, args []any) (any, error) {
	tb, ok := w.(*widget.TabBook)
	if !ok {
		return nil, fmt.Errorf("TabLabel requires widget of type *TabBook, got %T", w)
	}

	activeTab := tb.Tab()
	if activeTab == nil {
		return "", nil
	}

	return internal.GetTabBookTabLabel(activeTab), nil
}

// handleSetTabByLabel sets the active tab by matching label.
func handleSetTabByLabel(w widget.PreferredSizeLocateableWidget, args []any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("SetTabByLabel requires 1 argument (label)")
	}

	label := args[0].(string)

	tb, ok := w.(*widget.TabBook)
	if !ok {
		return nil, fmt.Errorf("SetTabByLabel requires widget of type *TabBook, got %T", w)
	}

	tabs := internal.GetTabBookTabs(tb)
	for _, tab := range tabs {
		if internal.GetTabBookTabLabel(tab) == label {
			if tab.Disabled {
				return nil, fmt.Errorf("tab with label '%s' is disabled", label)
			}
			tb.SetTab(tab)
			return nil, nil
		}
	}

	return nil, fmt.Errorf("tab with label '%s' not found", label)
}
```

- [ ] **Step 4: Register TabBook handlers in proxy.go**

```go
// Add to proxyHandlers map in autoui/proxy.go

var proxyHandlers = map[string]ProxyHandler{
	"SelectEntryByIndex": handleSelectEntryByIndex,
	"SelectedEntryIndex": handleSelectedEntryIndex,
	// TabBook handlers
	"Tabs":           handleTabs,
	"TabIndex":       handleTabIndex,
	"SetTabByIndex":  handleSetTabByIndex,
	"TabLabel":       handleTabLabel,
	"SetTabByLabel":  handleSetTabByLabel,
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./autoui -run TestHandleTabs -v`
Expected: PASS

- [ ] **Step 6: Write tests for remaining TabBook handlers**

```go
// Add to autoui/proxy_tabbook_test.go

func TestHandleTabIndex(t *testing.T) {
	tab1 := widget.NewTabBookTab(widget.TabBookTabOpts.Label("General"))
	tab2 := widget.NewTabBookTab(widget.TabBookTabOpts.Label("Settings"))

	buttonImage := &widget.ButtonImage{
		Idle:    createTestNineSlice(100, 30, color.Gray{150}),
		Pressed: createTestNineSlice(100, 30, color.Gray{100}),
	}

	tb := widget.NewTabBook(
		widget.TabBookOpts.Tabs(tab1, tab2),
		widget.TabBookOpts.TabButtonImage(buttonImage),
		widget.TabBookOpts.InitialTab(tab1),
	)
	tb.Validate()

	result, err := autoui.InvokeMethod(tb, "TabIndex", nil)
	if err != nil {
		t.Fatalf("TabIndex failed: %v", err)
	}

	index := result.(int64)
	if index != 0 {
		t.Errorf("expected tab index 0, got %d", index)
	}
}

func TestHandleSetTabByIndex(t *testing.T) {
	tab1 := widget.NewTabBookTab(widget.TabBookTabOpts.Label("General"))
	tab2 := widget.NewTabBookTab(widget.TabBookTabOpts.Label("Settings"))

	buttonImage := &widget.ButtonImage{
		Idle:    createTestNineSlice(100, 30, color.Gray{150}),
		Pressed: createTestNineSlice(100, 30, color.Gray{100}),
	}

	tb := widget.NewTabBook(
		widget.TabBookOpts.Tabs(tab1, tab2),
		widget.TabBookOpts.TabButtonImage(buttonImage),
		widget.TabBookOpts.InitialTab(tab1),
	)
	tb.Validate()

	_, err := autoui.InvokeMethod(tb, "SetTabByIndex", []any{1.0})
	if err != nil {
		t.Fatalf("SetTabByIndex failed: %v", err)
	}

	// Verify active tab changed
	result, _ := autoui.InvokeMethod(tb, "TabIndex", nil)
	index := result.(int64)
	if index != 1 {
		t.Errorf("expected tab index 1 after SetTabByIndex, got %d", index)
	}
}

func TestHandleTabLabel(t *testing.T) {
	tab1 := widget.NewTabBookTab(widget.TabBookTabOpts.Label("General"))
	tab2 := widget.NewTabBookTab(widget.TabBookTabOpts.Label("Settings"))

	buttonImage := &widget.ButtonImage{
		Idle:    createTestNineSlice(100, 30, color.Gray{150}),
		Pressed: createTestNineSlice(100, 30, color.Gray{100}),
	}

	tb := widget.NewTabBook(
		widget.TabBookOpts.Tabs(tab1, tab2),
		widget.TabBookOpts.TabButtonImage(buttonImage),
		widget.TabBookOpts.InitialTab(tab2),
	)
	tb.Validate()

	result, err := autoui.InvokeMethod(tb, "TabLabel", nil)
	if err != nil {
		t.Fatalf("TabLabel failed: %v", err)
	}

	label := result.(string)
	if label != "Settings" {
		t.Errorf("expected label 'Settings', got '%s'", label)
	}
}

func TestHandleSetTabByLabel(t *testing.T) {
	tab1 := widget.NewTabBookTab(widget.TabBookTabOpts.Label("General"))
	tab2 := widget.NewTabBookTab(widget.TabBookTabOpts.Label("Settings"))

	buttonImage := &widget.ButtonImage{
		Idle:    createTestNineSlice(100, 30, color.Gray{150}),
		Pressed: createTestNineSlice(100, 30, color.Gray{100}),
	}

	tb := widget.NewTabBook(
		widget.TabBookOpts.Tabs(tab1, tab2),
		widget.TabBookOpts.TabButtonImage(buttonImage),
		widget.TabBookOpts.InitialTab(tab1),
	)
	tb.Validate()

	_, err := autoui.InvokeMethod(tb, "SetTabByLabel", []any{"Settings"})
	if err != nil {
		t.Fatalf("SetTabByLabel failed: %v", err)
	}

	// Verify active tab changed
	result, _ := autoui.InvokeMethod(tb, "TabLabel", nil)
	label := result.(string)
	if label != "Settings" {
		t.Errorf("expected label 'Settings' after SetTabByLabel, got '%s'", label)
	}
}
```

- [ ] **Step 7: Run all TabBook tests**

Run: `go test ./autoui -run TestHandleTab -v`
Expected: PASS for all tests

- [ ] **Step 8: Commit TabBook proxy handlers**

```bash
git add autoui/proxy_tabbook.go autoui/proxy_tabbook_test.go autoui/proxy.go
git commit -m "feat(autoui): add TabBook proxy handlers

Add Tabs, TabIndex, SetTabByIndex, TabLabel, and SetTabByLabel
proxy handlers for TabBook automation. Uses reflection to access
private tabs and label fields."
```

---

### Task 4: RadioGroup Proxy Handlers

**Files:**
- Create: `autoui/proxy_radiogroup.go`
- Create: `autoui/proxy_radiogroup_test.go`
- Modify: `autoui/handlers.go`

**Goal:** Add proxy handlers for RadioGroup: Elements, ActiveIndex, SetActiveByIndex, ActiveLabel, SetActiveByLabel.

- [ ] **Step 1: Write the failing test for handleRadioGroupElements**

```go
// autoui/proxy_radiogroup_test.go
package autoui_test

import (
	"testing"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/s3cy/autoebiten/autoui"
)

func TestInvokeRadioGroupMethod(t *testing.T) {
	btn1 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	btn2 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	btn3 := widget.NewButton(widget.ButtonOpts.ToggleMode())

	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1, btn2, btn3))

	// Register RadioGroup
	autoui.RegisterRadioGroup("test-group", rg)

	// Test via invokeRadioGroupMethod (not InvokeMethod, since RadioGroup isn't a widget)
	result, err := autoui.InvokeRadioGroupMethod(rg, "Elements", nil)
	if err != nil {
		t.Fatalf("Elements failed: %v", err)
	}

	elements := result.([]autoui.RadioGroupElementInfo)
	if len(elements) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(elements))
	}

	autoui.UnregisterRadioGroup("test-group")
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./autoui -run TestInvokeRadioGroupMethod -v`
Expected: FAIL with "undefined: autoui.InvokeRadioGroupMethod"

- [ ] **Step 3: Create proxy_radiogroup.go**

```go
// autoui/proxy_radiogroup.go
package autoui

import (
	"fmt"
	"reflect"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/s3cy/autoebiten/autoui/internal"
)

// RadioGroupElementInfo represents information about a RadioGroup element.
type RadioGroupElementInfo struct {
	Type   string `json:"type"`
	Active bool   `json:"active"`
	Label  string `json:"label,omitempty"`
}

// radioGroupProxyHandlers maps method names to RadioGroup handlers.
var radioGroupProxyHandlers = map[string]RadioGroupHandler{
	"Elements":        handleRadioGroupElements,
	"ActiveIndex":     handleRadioGroupActiveIndex,
	"SetActiveByIndex": handleRadioGroupSetActiveByIndex,
	"ActiveLabel":     handleRadioGroupActiveLabel,
	"SetActiveByLabel": handleRadioGroupSetActiveByLabel,
}

// RadioGroupHandler handles a named operation on a RadioGroup.
type RadioGroupHandler func(rg *widget.RadioGroup, args []any) (any, error)

// InvokeRadioGroupMethod invokes a method on a RadioGroup using proxy handlers.
func InvokeRadioGroupMethod(rg *widget.RadioGroup, method string, args []any) (any, error) {
	handler, ok := radioGroupProxyHandlers[method]
	if !ok {
		return nil, fmt.Errorf("method '%s' not supported for RadioGroup", method)
	}
	return handler(rg, args)
}

// handleRadioGroupElements returns information about all elements in the group.
func handleRadioGroupElements(rg *widget.RadioGroup, args []any) (any, error) {
	elements := internal.GetRadioGroupElements(rg)
	active := rg.Active()

	result := make([]RadioGroupElementInfo, len(elements))
	for i, elem := range elements {
		info := RadioGroupElementInfo{
			Active: elem == active,
		}

		// Extract type and label based on concrete type
		switch v := elem.(type) {
		case *widget.Button:
			info.Type = "Button"
			if text := v.Text(); text != nil {
				info.Label = text.Label
			}
		case *widget.Checkbox:
			info.Type = "Checkbox"
			if text := v.Text(); text != nil {
				info.Label = text.Label
			}
		default:
			info.Type = reflect.TypeOf(elem).String()
		}

		result[i] = info
	}

	return result, nil
}

// handleRadioGroupActiveIndex returns the index of the active element.
func handleRadioGroupActiveIndex(rg *widget.RadioGroup, args []any) (any, error) {
	elements := internal.GetRadioGroupElements(rg)
	active := rg.Active()

	for i, elem := range elements {
		if elem == active {
			return int64(i), nil
		}
	}

	return int64(-1), nil
}

// handleRadioGroupSetActiveByIndex sets the active element by index.
func handleRadioGroupSetActiveByIndex(rg *widget.RadioGroup, args []any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("SetActiveByIndex requires 1 argument (index)")
	}

	index := int(args[0].(float64))
	elements := internal.GetRadioGroupElements(rg)

	if index < 0 || index >= len(elements) {
		return nil, fmt.Errorf("index %d out of range (0-%d)", index, len(elements)-1)
	}

	rg.SetActive(elements[index])
	return nil, nil
}

// handleRadioGroupActiveLabel returns the label of the active element.
func handleRadioGroupActiveLabel(rg *widget.RadioGroup, args []any) (any, error) {
	active := rg.Active()
	if active == nil {
		return "", nil
	}

	switch v := active.(type) {
	case *widget.Button:
		if text := v.Text(); text != nil {
			return text.Label, nil
		}
	case *widget.Checkbox:
		if text := v.Text(); text != nil {
			return text.Label, nil
		}
	}

	return "", nil
}

// handleRadioGroupSetActiveByLabel sets the active element by matching label.
func handleRadioGroupSetActiveByLabel(rg *widget.RadioGroup, args []any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("SetActiveByLabel requires 1 argument (label)")
	}

	label := args[0].(string)
	elements := internal.GetRadioGroupElements(rg)

	for _, elem := range elements {
		var elemLabel string

		switch v := elem.(type) {
		case *widget.Button:
			if text := v.Text(); text != nil {
				elemLabel = text.Label
			}
		case *widget.Checkbox:
			if text := v.Text(); text != nil {
				elemLabel = text.Label
			}
		}

		if elemLabel == label {
			rg.SetActive(elem)
			return nil, nil
		}
	}

	return nil, fmt.Errorf("element with label '%s' not found", label)
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./autoui -run TestInvokeRadioGroupMethod -v`
Expected: PASS

- [ ] **Step 5: Write tests for remaining RadioGroup handlers**

```go
// Add to autoui/proxy_radiogroup_test.go

func TestHandleRadioGroupActiveIndex(t *testing.T) {
	btn1 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	btn2 := widget.NewButton(widget.ButtonOpts.ToggleMode())

	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1, btn2))

	autoui.RegisterRadioGroup("test-index", rg)

	result, err := autoui.InvokeRadioGroupMethod(rg, "ActiveIndex", nil)
	if err != nil {
		t.Fatalf("ActiveIndex failed: %v", err)
	}

	// Default active should be first element (index 0)
	index := result.(int64)
	if index != 0 {
		t.Errorf("expected active index 0, got %d", index)
	}

	autoui.UnregisterRadioGroup("test-index")
}

func TestHandleRadioGroupSetActiveByIndex(t *testing.T) {
	btn1 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	btn2 := widget.NewButton(widget.ButtonOpts.ToggleMode())

	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1, btn2))

	autoui.RegisterRadioGroup("test-set", rg)

	_, err := autoui.InvokeRadioGroupMethod(rg, "SetActiveByIndex", []any{1.0})
	if err != nil {
		t.Fatalf("SetActiveByIndex failed: %v", err)
	}

	// Verify active changed
	result, _ := autoui.InvokeRadioGroupMethod(rg, "ActiveIndex", nil)
	index := result.(int64)
	if index != 1 {
		t.Errorf("expected active index 1, got %d", index)
	}

	autoui.UnregisterRadioGroup("test-set")
}

func TestHandleRadioGroupSetActiveByLabel(t *testing.T) {
	btn1 := widget.NewButton(widget.ButtonOpts.Text("Option A", nil, &widget.ButtonTextColor{Idle: color.White}))
	btn2 := widget.NewButton(widget.ButtonOpts.Text("Option B", nil, &widget.ButtonTextColor{Idle: color.White}))

	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1, btn2))

	autoui.RegisterRadioGroup("test-label", rg)

	_, err := autoui.InvokeRadioGroupMethod(rg, "SetActiveByLabel", []any{"Option B"})
	if err != nil {
		t.Fatalf("SetActiveByLabel failed: %v", err)
	}

	result, _ := autoui.InvokeRadioGroupMethod(rg, "ActiveLabel", nil)
	label := result.(string)
	if label != "Option B" {
		t.Errorf("expected active label 'Option B', got '%s'", label)
	}

	autoui.UnregisterRadioGroup("test-label")
}
```

- [ ] **Step 6: Run all RadioGroup tests**

Run: `go test ./autoui -run TestHandleRadioGroup -v`
Expected: PASS

- [ ] **Step 7: Commit RadioGroup proxy handlers**

```bash
git add autoui/proxy_radiogroup.go autoui/proxy_radiogroup_test.go
git commit -m "feat(autoui): add RadioGroup proxy handlers

Add Elements, ActiveIndex, SetActiveByIndex, ActiveLabel, and
SetActiveByLabel proxy handlers for RadioGroup automation.
RadioGroups must be registered via RegisterRadioGroup before use."
```

---

### Task 5: Tree Traversal - TabBookTab Injection

**Files:**
- Modify: `autoui/tree.go`
- Modify: `autoui/xml.go`
- Modify: `autoui/internal/widgetstate.go`

**Goal:** Inject TabBookTab children into TabBook tree output with label and content.

- [ ] **Step 1: Write the failing test for TabBook tree output**

```go
// Add to autoui/tree_test.go

func TestWalkTreeTabBookInjection(t *testing.T) {
	// Create TabBook with tabs
	tab1 := widget.NewTabBookTab(widget.TabBookTabOpts.Label("General"))
	tab1Content := widget.NewButton(widget.ButtonOpts.Text("Save", nil, &widget.ButtonTextColor{Idle: color.White}))
	tab1.AddChild(tab1Content)

	tab2 := widget.NewTabBookTab(widget.TabBookTabOpts.Label("Settings"))
	tab2.Disabled = true

	buttonImage := &widget.ButtonImage{
		Idle:    createTestNineSlice(100, 30, color.Gray{150}),
		Pressed: createTestNineSlice(100, 30, color.Gray{100}),
	}

	tb := widget.NewTabBook(
		widget.TabBookOpts.Tabs(tab1, tab2),
		widget.TabBookOpts.TabButtonImage(buttonImage),
	)
	tb.Validate()

	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout()),
	)
	root.AddChild(tb)

	// Walk tree
	widgets := WalkTree(root)

	// Find TabBook in result
	var tabBookInfo *WidgetInfo
	for i, w := range widgets {
		if w.Type == "TabBook" {
			tabBookInfo = &widgets[i]
			break
		}
	}

	if tabBookInfo == nil {
		t.Fatal("TabBook not found in tree")
	}

	// Check that TabBook has injected state for tabs
	// The tabs should be represented in the tree structure via XML marshaling
	// Here we verify the widgetinfo contains the expected tab data
	xmlData, err := MarshalWidgetTreeXML(widgets)
	if err != nil {
		t.Fatalf("failed to marshal tree: %v", err)
	}

	// Check XML contains TabBookTab elements
	if !strings.Contains(string(xmlData), "<TabBookTab") {
		t.Error("XML should contain TabBookTab elements")
	}
	if !strings.Contains(string(xmlData), "label=\"General\"") {
		t.Error("XML should contain tab label 'General'")
	}
	if !strings.Contains(string(xmlData), "disabled=\"true\"") {
		t.Error("XML should show tab2 as disabled")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./autoui -run TestWalkTreeTabBookInjection -v`
Expected: FAIL - TabBookTab elements not present in XML

- [ ] **Step 3: Modify tree.go for TabBook handling**

```go
// Modify walkTreeRecursive in autoui/tree.go

func walkTreeRecursive(w widget.PreferredSizeLocateableWidget, result *[]WidgetInfo) {
	// Extract info for current widget
	info := ExtractWidgetInfo(w)
	*result = append(*result, info)

	// Handle TabBook specially - inject TabBookTab children
	if tb, ok := w.(*widget.TabBook); ok {
		injectTabBookTabs(tb, result)
		return // TabBook's normal children are handled via tabs
	}

	// If this is a Container, traverse children
	if container, ok := w.(widget.Containerer); ok {
		children := container.Children()
		for _, child := range children {
			walkTreeRecursive(child, result)
		}
	}
}

// injectTabBookTabs traverses TabBook tabs and adds them to the result.
func injectTabBookTabs(tb *widget.TabBook, result *[]WidgetInfo) {
	tabs := internal.GetTabBookTabs(tb)
	for _, tab := range tabs {
		// Create synthetic WidgetInfo for TabBookTab
		tabInfo := WidgetInfo{
			Widget:   tab, // TabBookTab embeds Container which implements PreferredSizeLocateableWidget
			Type:     "TabBookTab",
			Rect:     tab.GetWidget().Rect,
			Visible:  tab.GetWidget().IsVisible(),
			Disabled: tab.Disabled,
			Addr:     fmt.Sprintf("0x%x", reflect.ValueOf(tab).Pointer()),
			State: map[string]string{
				"label":    internal.GetTabBookTabLabel(tab),
				"disabled": fmt.Sprintf("%v", tab.Disabled),
			},
		}
		*result = append(*result, tabInfo)

		// Traverse tab's children (TabBookTab embeds Container)
		walkTreeRecursive(tab, result)
	}
}
```

- [ ] **Step 4: Modify xml.go for TabBook parent handling**

The existing `buildNodeTree` function handles parent-child relationships. TabBookTab needs to have TabBook as its parent.

```go
// Modify buildNodeTree in autoui/xml.go to handle TabBookTab parent relationship

func buildNodeTree(root *WidgetNode, widgets []WidgetInfo) {
	if len(widgets) <= 1 {
		return
	}

	widgetToNode := map[*widget.Widget]*WidgetNode{}
	widgetToNode[widgets[0].Widget.GetWidget()] = root

	for _, info := range widgets[1:] {
		node := widgetInfoToNode(info)
		baseWidget := info.Widget.GetWidget()
		widgetToNode[baseWidget] = node

		// Special handling for TabBookTab - parent is TabBook, not Container
		if info.Type == "TabBookTab" {
			// Find parent TabBook by traversing back
			parentWidget := baseWidget.Parent()
			if parentNode, ok := widgetToNode[parentWidget]; ok {
				parentNode.Children = append(parentNode.Children, node)
			} else {
				root.Children = append(root.Children, node)
			}
			continue
		}

		// Normal parent lookup
		parentWidget := baseWidget.Parent()
		if parentNode, ok := widgetToNode[parentWidget]; ok {
			parentNode.Children = append(parentNode.Children, node)
		} else {
			root.Children = append(root.Children, node)
		}
	}
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./autoui -run TestWalkTreeTabBookInjection -v`
Expected: PASS

- [ ] **Step 6: Update TabBook state extraction in widgetstate.go**

```go
// Modify extractTabBookState in autoui/internal/widgetstate.go

func extractTabBookState(tb *widget.TabBook, result map[string]string) {
	if tab := tb.Tab(); tab != nil {
		result["active_tab"] = GetTabBookTabLabel(tab)
	}
}
```

- [ ] **Step 7: Commit TabBook tree injection**

```bash
git add autoui/tree.go autoui/xml.go autoui/internal/widgetstate.go autoui/tree_test.go
git commit -m "feat(autoui): inject TabBookTab children in tree output

Modify WalkTree to inject synthetic TabBookTab elements with label
and disabled attributes. Update XML building to handle TabBookTab
parent relationships. Update widgetstate to show active tab label."
```

---

### Task 6: Tree Traversal - RadioGroup Injection

**Files:**
- Modify: `autoui/tree.go`
- Modify: `autoui/xml.go`
- Modify: `autoui/tree_test.go`

**Goal:** Append registered RadioGroups to tree output at `<UI>` root level.

- [ ] **Step 1: Write the failing test for RadioGroup tree output**

```go
// Add to autoui/tree_test.go

func TestWalkTreeRadioGroupInjection(t *testing.T) {
	btn1 := widget.NewButton(widget.ButtonOpts.Text("Option A", nil, &widget.ButtonTextColor{Idle: color.White}))
	btn2 := widget.NewButton(widget.ButtonOpts.Text("Option B", nil, &widget.ButtonTextColor{Idle: color.White}))

	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1, btn2))

	// Register RadioGroup
	RegisterRadioGroup("test-group", rg)

	root := widget.NewContainer()

	// Walk tree - RadioGroups should be appended
	widgets := WalkTree(root)

	// Find RadioGroup in result
	var radioGroupInfo *WidgetInfo
	for i, w := range widgets {
		if w.Type == "RadioGroup" {
			radioGroupInfo = &widgets[i]
			break
		}
	}

	if radioGroupInfo == nil {
		t.Fatal("RadioGroup not found in tree")
	}

	// Verify RadioGroup has name attribute
	if radioGroupInfo.State["name"] != "test-group" {
		t.Errorf("expected RadioGroup name 'test-group', got '%s'", radioGroupInfo.State["name"])
	}

	// Verify XML output
	xmlData, err := MarshalWidgetTreeXML(widgets)
	if err != nil {
		t.Fatalf("failed to marshal tree: %v", err)
	}

	if !strings.Contains(string(xmlData), "<RadioGroup") {
		t.Error("XML should contain RadioGroup element")
	}
	if !strings.Contains(string(xmlData), "name=\"test-group\"") {
		t.Error("XML should contain RadioGroup name attribute")
	}

	UnregisterRadioGroup("test-group")
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./autoui -run TestWalkTreeRadioGroupInjection -v`
Expected: FAIL - RadioGroup not found in tree

- [ ] **Step 3: Modify WalkTree to append RadioGroups**

```go
// Modify WalkTree in autoui/tree.go

func WalkTree(root widget.PreferredSizeLocateableWidget) []WidgetInfo {
	if root == nil {
		return nil
	}

	var result []WidgetInfo
	walkTreeRecursive(root, &result)

	// Append registered RadioGroups
	injectRegisteredRadioGroups(&result)

	return result
}

// injectRegisteredRadioGroups appends registered RadioGroups to the result.
func injectRegisteredRadioGroups(result *[]WidgetInfo) {
	names := GetRegisteredRadioGroups()
	for _, name := range names {
		rg := GetRadioGroup(name)
		if rg == nil {
			continue
		}

		// Create synthetic RadioGroup WidgetInfo
		info := WidgetInfo{
			Widget: nil, // RadioGroup is not a widget
			Type:   "RadioGroup",
			Rect:   image.Rectangle{}, // No geometry
			Visible: true,
			Disabled: false,
			Addr:     "", // No pointer for synthetic element
			State: map[string]string{
				"name": name,
			},
		}

		// Get active index
		active := rg.Active()
		elements := internal.GetRadioGroupElements(rg)
		for i, elem := range elements {
			if elem == active {
				info.State["active_index"] = fmt.Sprintf("%d", i)
				break
			}
		}

		*result = append(*result, info)

		// Inject element widgets as children
		for _, elem := range elements {
			// Convert RadioGroupElement to WidgetInfo
			switch v := elem.(type) {
			case widget.PreferredSizeLocateableWidget:
				elemInfo := ExtractWidgetInfo(v)
				elemInfo.State["active"] = fmt.Sprintf("%v", elem == active)
				*result = append(*result, elemInfo)
			}
		}
	}
}
```

- [ ] **Step 4: Modify xml.go for RadioGroup handling**

```go
// Modify buildNodeTree in autoui/xml.go to handle RadioGroup

func buildNodeTree(root *WidgetNode, widgets []WidgetInfo) {
	if len(widgets) <= 1 {
		return
	}

	widgetToNode := map[*widget.Widget]*WidgetNode{}
	widgetToNode[widgets[0].Widget.GetWidget()] = root

	// Track RadioGroup nodes for element parent lookup
	radioGroupNodes := map[string]*WidgetNode{}

	for _, info := range widgets[1:] {
		node := widgetInfoToNode(info)

		// Special handling for synthetic RadioGroup
		if info.Type == "RadioGroup" {
			// RadioGroups are children of root UI node
			root.Children = append(root.Children, node)
			if name, ok := info.State["name"]; ok {
				radioGroupNodes[name] = node
			}
			widgetToNode[info.Widget.GetWidget()] = node // nil widget, but track for consistency
			continue
		}

		// Handle RadioGroup elements - they should be children of RadioGroup
		if active, hasActive := info.State["active"]; hasActive {
			// Find parent RadioGroup by checking previous widgets
			for i := len(widgets) - 1; i >= 0; i-- {
				if widgets[i].Type == "RadioGroup" {
					if rgNode, ok := widgetToNode[widgets[i].Widget.GetWidget()]; ok {
						rgNode.Children = append(rgNode.Children, node)
						break
					}
				}
			}
			continue
		}

		// Normal handling
		if info.Widget != nil {
			baseWidget := info.Widget.GetWidget()
			widgetToNode[baseWidget] = node

			parentWidget := baseWidget.Parent()
			if parentNode, ok := widgetToNode[parentWidget]; ok {
				parentNode.Children = append(parentNode.Children, node)
			} else {
				root.Children = append(root.Children, node)
			}
		} else {
			root.Children = append(root.Children, node)
		}
	}
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./autoui -run TestWalkTreeRadioGroupInjection -v`
Expected: PASS

- [ ] **Step 6: Commit RadioGroup tree injection**

```bash
git add autoui/tree.go autoui/xml.go autoui/tree_test.go
git commit -m "feat(autoui): inject registered RadioGroups in tree output

Append registered RadioGroups to WalkTree result with synthetic
WidgetInfo entries. RadioGroup elements (Button/Checkbox) are
included as children. RadioGroups appear at UI root level in XML."
```

---

### Task 7: Call Handler - RadioGroup Target Syntax

**Files:**
- Modify: `autoui/handlers.go`
- Modify: `autoui/handlers_test.go`

**Goal:** Support `radiogroup=name` target syntax in autoui.call command.

- [ ] **Step 1: Write the failing test for radiogroup target**

```go
// Add to autoui/handlers_test.go

func TestHandleCallCommandRadioGroup(t *testing.T) {
	btn1 := widget.NewButton(widget.ButtonOpts.Text("Option A", nil, &widget.ButtonTextColor{Idle: color.White}))
	btn2 := widget.NewButton(widget.ButtonOpts.Text("Option B", nil, &widget.ButtonTextColor{Idle: color.White}))

	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1, btn2))

	ui := &ebitenui.UI{
		Container: widget.NewContainer(),
	}

	RegisterWithOptions(ui, "autoui", nil)
	RegisterRadioGroup("test-call-group", rg)

	// Create mock command context
	ctx := createMockCommandContext(`{"target":"radiogroup=test-call-group","method":"ActiveIndex","args":[]}`)

	// This would normally be called by autoebiten command system
	// For testing, we call the handler directly
	handler := handleCallCommand(ui)
	handler(ctx)

	// Check response contains success
	response := ctx.GetResponse()
	if !strings.Contains(response, `"success":true`) {
		t.Errorf("expected success response, got: %s", response)
	}

	UnregisterRadioGroup("test-call-group")
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./autoui -run TestHandleCallCommandRadioGroup -v`
Expected: FAIL - "radiogroup=" prefix not handled

- [ ] **Step 3: Modify handleCallCommand in handlers.go**

```go
// Modify handleCallCommand in autoui/handlers.go

const (
	radioGroupPrefix = "radiogroup="
)

func handleCallCommand(ui *ebitenui.UI) func(ctx autoebiten.CommandContext) {
	return func(ctx autoebiten.CommandContext) {
		if ui == nil {
			ctx.Respond("error: UI not registered")
			return
		}

		request := ctx.Request()
		if request == "" {
			ctx.Respond("error: missing call request")
			return
		}

		var callReq CallRequest
		if err := json.Unmarshal([]byte(request), &callReq); err != nil {
			ctx.Respond("error: invalid JSON format: " + err.Error())
			return
		}

		if callReq.Target == "" {
			ctx.Respond("error: missing target query")
			return
		}

		if callReq.Method == "" {
			ctx.Respond("error: missing method name")
			return
		}

		// Check for RadioGroup target prefix
		if strings.HasPrefix(callReq.Target, radioGroupPrefix) {
			name := strings.TrimPrefix(callReq.Target, radioGroupPrefix)
			rg := GetRadioGroup(name)
			if rg == nil {
				ctx.Respond(fmt.Sprintf("error: RadioGroup '%s' not registered. Did you call autoui.RegisterRadioGroup?", name))
				return
			}

			result, err := InvokeRadioGroupMethod(rg, callReq.Method, callReq.Args)

			response := CallResponse{Success: err == nil}
			if err != nil {
				response.Error = err.Error()
			}
			if result != nil {
				response.Result = result
			}

			respData, _ := json.Marshal(response)
			ctx.Respond(string(respData))
			return
		}

		// Normal widget lookup (existing code)
		widgets := SnapshotTree(ui)

		var targetWidget *WidgetInfo
		var matching []WidgetInfo

		if len(callReq.Target) > 0 && callReq.Target[0] == '{' {
			matching = FindByQueryJSON(widgets, callReq.Target)
		} else {
			matching = FindByQuery(widgets, callReq.Target)
		}

		if len(matching) == 0 {
			ctx.Respond("error: no widget found matching target query")
			return
		}

		targetWidget = &matching[0]

		result, err := InvokeMethod(targetWidget.Widget, callReq.Method, callReq.Args)

		response := CallResponse{Success: err == nil}
		if err != nil {
			response.Error = err.Error()
		}
		if result != nil {
			response.Result = result
		}

		respData, _ := json.Marshal(response)
		ctx.Respond(string(respData))
	}
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./autoui -run TestHandleCallCommandRadioGroup -v`
Expected: PASS

- [ ] **Step 5: Commit RadioGroup call handler**

```bash
git add autoui/handlers.go autoui/handlers_test.go
git commit -m "feat(autoui): support radiogroup=name target syntax in call command

Handle 'radiogroup=' prefix in handleCallCommand to lookup registered
RadioGroups and invoke methods via InvokeRadioGroupMethod. Provides
clear error message when RadioGroup is not registered."
```

---

### Task 8: Integration Tests

**Files:**
- Modify: `autoui/integration_test.go`

**Goal:** Add full workflow integration tests for RadioGroup and TabBook.

- [ ] **Step 1: Write RadioGroup integration test**

```go
// Add to autoui/integration_test.go

func TestRadioGroupIntegration(t *testing.T) {
	// Setup
	btn1 := widget.NewButton(widget.ButtonOpts.Text("Option A", nil, &widget.ButtonTextColor{Idle: color.White}))
	btn2 := widget.NewButton(widget.ButtonOpts.Text("Option B", nil, &widget.ButtonTextColor{Idle: color.White}))
	btn3 := widget.NewButton(widget.ButtonOpts.Text("Option C", nil, &widget.ButtonTextColor{Idle: color.White}))

	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1, btn2, btn3))

	ui := &ebitenui.UI{Container: widget.NewContainer()}
	RegisterWithOptions(ui, "autoui", nil)
	RegisterRadioGroup("options", rg)

	// Test 1: Get elements
	output, err := runCustomCommand("autoui.call", `{"target":"radiogroup=options","method":"Elements"}`)
	require.NoError(t, err)
	assert.Contains(t, output, `"success":true`)
	assert.Contains(t, output, "Option A")
	assert.Contains(t, output, "Option B")

	// Test 2: Get active index
	output, err = runCustomCommand("autoui.call", `{"target":"radiogroup=options","method":"ActiveIndex"}`)
	require.NoError(t, err)
	assert.Contains(t, output, `"success":true`)
	assert.Contains(t, output, `"result":0`) // Default is first

	// Test 3: Set active by index
	output, err = runCustomCommand("autoui.call", `{"target":"radiogroup=options","method":"SetActiveByIndex","args":[2]}`)
	require.NoError(t, err)
	assert.Contains(t, output, `"success":true`)

	// Test 4: Verify change
	output, err = runCustomCommand("autoui.call", `{"target":"radiogroup=options","method":"ActiveLabel"}`)
	require.NoError(t, err)
	assert.Contains(t, output, "Option C")

	// Test 5: Tree output includes RadioGroup
	output, err = runCustomCommand("autoui.tree", "")
	require.NoError(t, err)
	assert.Contains(t, output, "<RadioGroup")
	assert.Contains(t, output, "name=\"options\"")

	// Cleanup
	UnregisterRadioGroup("options")
}
```

- [ ] **Step 2: Write TabBook integration test**

```go
// Add to autoui/integration_test.go

func TestTabBookIntegration(t *testing.T) {
	// Setup tabs with content
	tab1 := widget.NewTabBookTab(widget.TabBookTabOpts.Label("General"))
	saveBtn := widget.NewButton(widget.ButtonOpts.Text("Save", nil, &widget.ButtonTextColor{Idle: color.White}))
	tab1.AddChild(saveBtn)

	tab2 := widget.NewTabBookTab(widget.TabBookTabOpts.Label("Settings"))
	tab2.Disabled = true

	buttonImage := &widget.ButtonImage{
		Idle:    createTestNineSlice(100, 30, color.Gray{150}),
		Pressed: createTestNineSlice(100, 30, color.Gray{100}),
	}

	tb := widget.NewTabBook(
		widget.TabBookOpts.Tabs(tab1, tab2),
		widget.TabBookOpts.TabButtonImage(buttonImage),
	)
	tb.Validate()

	root := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewRowLayout()))
	root.AddChild(tb)

	ui := &ebitenui.UI{Container: root}
	RegisterWithOptions(ui, "autoui", nil)

	// Test 1: Get tabs
	output, err := runCustomCommand("autoui.call", `{"target":"type=TabBook","method":"Tabs"}`)
	require.NoError(t, err)
	assert.Contains(t, output, `"success":true`)
	assert.Contains(t, output, "General")
	assert.Contains(t, output, "Settings")
	assert.Contains(t, output, `"disabled":true`)

	// Test 2: Get active tab label
	output, err = runCustomCommand("autoui.call", `{"target":"type=TabBook","method":"TabLabel"}`)
	require.NoError(t, err)
	assert.Contains(t, output, "General")

	// Test 3: Tree output includes TabBookTab
	output, err = runCustomCommand("autoui.tree", "")
	require.NoError(t, err)
	assert.Contains(t, output, "<TabBookTab")
	assert.Contains(t, output, "label=\"General\"")
	assert.Contains(t, output, "label=\"Settings\"")
	assert.Contains(t, output, "disabled=\"true\"")
}
```

- [ ] **Step 3: Run integration tests**

Run: `go test ./autoui -run TestRadioGroupIntegration -v`
Run: `go test ./autoui -run TestTabBookIntegration -v`
Expected: PASS for both

- [ ] **Step 4: Commit integration tests**

```bash
git add autoui/integration_test.go
git commit -m "test(autoui): add RadioGroup and TabBook integration tests

Add full workflow tests covering registration, method invocation,
and tree output for RadioGroup and TabBook automation."
```

---

### Task 9: Documentation Update

**Files:**
- Modify: `docs/autoui.md`
- Modify: `skills/using-autoebiten/references/autoui.md`

**Goal:** Document RadioGroup and TabBook automation in reference docs.

- [ ] **Step 1: Update widget-specific attributes table**

Add to the Widget-specific attributes table in both docs:

```
| RadioGroup | name, active_index, elements (synthetic) |
| TabBookTab | label, disabled, active (synthetic) |
```

- [ ] **Step 2: Add RadioGroup section**

```markdown
### RadioGroup Operations

**Note:** RadioGroup requires explicit registration. It is not discoverable via tree traversal.

**Registration:**
```go
settingsGroup := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1, btn2, btn3))
autoui.RegisterRadioGroup("settings-group", settingsGroup)
```

**CLI Usage:**
```bash
# Get elements
autoui.call --request '{"target":"radiogroup=settings-group","method":"Elements"}'

# Set active by index
autoui.call --request '{"target":"radiogroup=settings-group","method":"SetActiveByIndex","args":[1]}'

# Set active by label
autoui.call --request '{"target":"radiogroup=settings-group","method":"SetActiveByLabel","args":["Option B"]}'
```

**Available Methods:**
- `Elements()` - Return element list with type, label, active status
- `ActiveIndex()` - Return index of active element
- `SetActiveByIndex(int)` - Set active element by position
- `ActiveLabel()` - Return label of active element
- `SetActiveByLabel(string)` - Set active by matching label
```

- [ ] **Step 3: Add TabBook section**

```markdown
### TabBook Operations

TabBook is discoverable via tree traversal (no registration needed).

**CLI Usage:**
```bash
# Get tabs
autoui.call --request '{"target":"type=TabBook","method":"Tabs"}'

# Set active tab by index
autoui.call --request '{"target":"type=TabBook","method":"SetTabByIndex","args":[1]}'

# Set active tab by label
autoui.call --request '{"target":"type=TabBook","method":"SetTabByLabel","args":["Settings"]}'
```

**Available Methods:**
- `Tabs()` - Return tab list with label, disabled status
- `TabIndex()` - Return index of active tab
- `SetTabByIndex(int)` - Set active tab by position
- `TabLabel()` - Return label of active tab
- `SetTabByLabel(string)` - Set tab by matching label (fails if disabled)
```

- [ ] **Step 4: Commit documentation**

```bash
git add docs/autoui.md skills/using-autoebiten/references/autoui.md
git commit -m "docs(autoui): add RadioGroup and TabBook automation sections

Document RadioGroup registration requirement, available methods,
and CLI usage examples. Document TabBook discovery and methods."
```

---

## Self-Review Checklist

**1. Spec coverage:**
- [x] Reflection utilities - Task 1
- [x] RadioGroup registration - Task 2
- [x] RadioGroup proxy handlers (Elements, ActiveIndex, SetActiveByIndex, ActiveLabel, SetActiveByLabel) - Task 4
- [x] TabBook proxy handlers (Tabs, TabIndex, SetTabByIndex, TabLabel, SetTabByLabel) - Task 3
- [x] TabBook tree injection - Task 5
- [x] RadioGroup tree injection - Task 6
- [x] RadioGroup call handler target syntax - Task 7
- [x] Integration tests - Task 8
- [x] Documentation - Task 9

**2. Placeholder scan:**
- No TBD, TODO, or vague descriptions found
- All code blocks contain complete implementations
- All commands have expected outputs

**3. Type consistency:**
- `TabInfo` defined in Task 3, used consistently
- `RadioGroupElementInfo` defined in Task 4, used consistently
- `InvokeRadioGroupMethod` signature matches handler usage
- Proxy handler signatures match `ProxyHandler` type