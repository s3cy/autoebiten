# autoui.call Extended Type Support Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Extend autoui.call to support any/interface{} types and enum types, enabling List selection and WidgetState/Visibility operations.

**Architecture:** Extend whitelist in caller.go, add proxy handlers in proxy.go, update CallResponse to capture returns.

**Tech Stack:** Go, reflection, encoding/json, ebitenui widget package

---

## File Structure

| File | Responsibility |
|------|----------------|
| `autoui/caller.go` | Extended whitelist, return capture, enum conversion |
| `autoui/handlers.go` | CallResponse struct with Result field |
| `autoui/proxy.go` (new) | Proxy handler registry and implementations |
| `autoui/caller_test.go` | Tests for extended whitelist |
| `autoui/proxy_test.go` (new) | Tests for proxy handlers |

---

### Task 1: Extend isWhitelistedSignature for Return Values

**Files:**
- Modify: `autoui/caller.go:53-78`
- Test: `autoui/caller_test.go`

- [ ] **Step 1: Write the failing test for any return**

```go
// autoui/caller_test.go - add new test
func TestIsWhitelistedSignature_AnyReturn(t *testing.T) {
    // Create a mock function type: func() any
    fn := func() any { return "test" }
    fnType := reflect.TypeOf(fn)

    if !isWhitelistedSignature(fnType) {
        t.Error("Expected func() any to be whitelisted for returns")
    }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd autoui && go test -run TestIsWhitelistedSignature_AnyReturn -v`
Expected: FAIL - "isWhitelistedSignature not exported" or test doesn't compile

- [ ] **Step 3: Export isWhitelistedSignature for testing (internal test)**

Create `autoui/caller_export_test.go`:

```go
package autoui

// Export for testing
var IsWhitelistedSignature = isWhitelistedSignature
```

Update test to use exported name:

```go
package autoui_test

func TestIsWhitelistedSignature_AnyReturn(t *testing.T) {
    fn := func() any { return "test" }
    fnType := reflect.TypeOf(fn)

    if !autoui.IsWhitelistedSignature(fnType) {
        t.Error("Expected func() any to be whitelisted for returns")
    }
}
```

- [ ] **Step 4: Run test to verify it fails with actual logic**

Run: `cd autoui && go test -run TestIsWhitelistedSignature_AnyReturn -v`
Expected: FAIL - "Expected func() any to be whitelisted"

- [ ] **Step 5: Extend isWhitelistedSignature to allow any/interface{} returns**

Modify `autoui/caller.go` - update `isWhitelistedSignature`:

```go
func isWhitelistedSignature(t reflect.Type) bool {
    numIn := t.NumIn()
    numOut := t.NumOut()

    // Check return values
    if numOut > 1 {
        return false
    }

    if numOut == 1 {
        returnType := t.Out(0)

        // Allow error
        if returnType.Implements(reflect.TypeOf((*error)(nil)).Elem()) {
            return true
        }

        // Allow any/interface{} (empty PkgPath, interface Kind)
        if returnType.Kind() == reflect.Interface && returnType.PkgPath() == "" {
            return true
        }

        // Allow slices of basic types or any
        if returnType.Kind() == reflect.Slice {
            elemType := returnType.Elem()
            if elemType.Kind() == reflect.Interface && elemType.PkgPath() == "" {
                return true // []any
            }
            if elemType.PkgPath() == "" {
                switch elemType.Kind() {
                case reflect.Bool, reflect.Int, reflect.Int32, reflect.Int64,
                    reflect.Float32, reflect.Float64, reflect.String:
                    return true
                }
            }
        }

        // Allow types with underlying basic Kind (enums)
        switch returnType.Kind() {
        case reflect.Bool, reflect.Int, reflect.Int32, reflect.Int64,
            reflect.Float32, reflect.Float64, reflect.String:
            return true
        }

        return false
    }

    // Check input parameters (existing logic)
    switch numIn {
    case 0:
        return true
    case 1:
        paramType := t.In(0)
        if paramType.PkgPath() != "" && !(paramType.Kind() == reflect.Interface && paramType.PkgPath() == "") {
            return false
        }
        switch paramType.Kind() {
        case reflect.Bool, reflect.Int, reflect.Int32, reflect.Int64,
            reflect.Float32, reflect.Float64, reflect.String:
            return true
        case reflect.Interface:
            if paramType.PkgPath() == "" {
                return true // any/interface{}
            }
            return false
        default:
            return false
        }
    default:
        return false
    }
}
```

- [ ] **Step 6: Run test to verify it passes**

Run: `cd autoui && go test -run TestIsWhitelistedSignature_AnyReturn -v`
Expected: PASS

- [ ] **Step 7: Write test for []any return**

```go
// autoui/caller_test.go - add new test
func TestIsWhitelistedSignature_SliceAnyReturn(t *testing.T) {
    fn := func() []any { return []any{"a", "b"} }
    fnType := reflect.TypeOf(fn)

    if !autoui.IsWhitelistedSignature(fnType) {
        t.Error("Expected func() []any to be whitelisted")
    }
}
```

- [ ] **Step 8: Run test to verify it passes**

Run: `cd autoui && go test -run TestIsWhitelistedSignature_SliceAnyReturn -v`
Expected: PASS

- [ ] **Step 9: Write test for enum return**

```go
// autoui/caller_test.go - add new test
func TestIsWhitelistedSignature_EnumReturn(t *testing.T) {
    // WidgetState is defined as: type WidgetState int
    fn := func() widget.WidgetState { return widget.WidgetChecked }
    fnType := reflect.TypeOf(fn)

    if !autoui.IsWhitelistedSignature(fnType) {
        t.Error("Expected func() WidgetState (enum) to be whitelisted")
    }
}
```

- [ ] **Step 10: Run test to verify it passes**

Run: `cd autoui && go test -run TestIsWhitelistedSignature_EnumReturn -v`
Expected: PASS

- [ ] **Step 11: Write test for any parameter**

```go
// autoui/caller_test.go - add new test
func TestIsWhitelistedSignature_AnyParam(t *testing.T) {
    fn := func(any) {}
    fnType := reflect.TypeOf(fn)

    if !autoui.IsWhitelistedSignature(fnType) {
        t.Error("Expected func(any) to be whitelisted for params")
    }
}
```

- [ ] **Step 12: Run test to verify it passes**

Run: `cd autoui && go test -run TestIsWhitelistedSignature_AnyParam -v`
Expected: PASS

- [ ] **Step 13: Write test for enum parameter**

```go
// autoui/caller_test.go - add new test
func TestIsWhitelistedSignature_EnumParam(t *testing.T) {
    fn := func(widget.WidgetState) {}
    fnType := reflect.TypeOf(fn)

    if !autoui.IsWhitelistedSignature(fnType) {
        t.Error("Expected func(WidgetState) to be whitelisted for enum params")
    }
}
```

- [ ] **Step 14: Run test to verify it passes**

Run: `cd autoui && go test -run TestIsWhitelistedSignature_EnumParam -v`
Expected: PASS

- [ ] **Step 15: Commit**

```bash
git add autoui/caller.go autoui/caller_export_test.go autoui/caller_test.go
git commit -m "$(cat <<'EOF'
feat(autoui): extend whitelist for any/interface{} and enum types

- Allow any/interface{} parameters and return values
- Allow []any and basic type slice returns  
- Allow enum types (underlying int/float/string/bool) for params and returns
- Add exported test helper for isWhitelistedSignature

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
EOF
)"
```

---

### Task 2: Extend convertArg for Enum and Interface{} Conversion

**Files:**
- Modify: `autoui/caller.go:116-156`
- Test: `autoui/caller_test.go`

- [ ] **Step 1: Write the failing test for int → enum conversion**

```go
// autoui/caller_test.go - add new test
func TestConvertArg_IntToEnum(t *testing.T) {
    targetType := reflect.TypeOf(widget.WidgetState(0))
    arg := float64(1) // JSON unmarshals numbers to float64

    result, err := autoui.ConvertArg(arg, targetType)
    if err != nil {
        t.Errorf("ConvertArg failed: %v", err)
    }

    // Should convert to WidgetState(1)
    if result.Int() != 1 {
        t.Errorf("Expected int value 1, got %d", result.Int())
    }
}
```

- [ ] **Step 2: Create export helper for convertArg**

Add to `autoui/caller_export_test.go`:

```go
// Export convertArg for testing
var ConvertArg = convertArg
```

- [ ] **Step 3: Run test to verify it fails**

Run: `cd autoui && go test -run TestConvertArg_IntToEnum -v`
Expected: FAIL - "unsupported target type" or similar

- [ ] **Step 4: Extend convertArg to handle convertible types**

Modify `autoui/caller.go` - update `convertArg`:

```go
func convertArg(arg any, targetType reflect.Type) (reflect.Value, error) {
    if arg == nil {
        return reflect.Value{}, fmt.Errorf("nil argument not supported")
    }

    argValue := reflect.ValueOf(arg)
    argType := argValue.Type()

    // Direct type match
    if argType == targetType {
        return argValue, nil
    }

    // If arg is convertible to target type (int → enum, etc.)
    if argType.ConvertibleTo(targetType) {
        return argValue.Convert(targetType), nil
    }

    // Handle interface{} targets - accept any non-nil value
    if targetType.Kind() == reflect.Interface && targetType.PkgPath() == "" {
        return argValue, nil
    }

    // Handle conversions based on target type's underlying kind
    switch targetType.Kind() {
    case reflect.Bool:
        if argType.Kind() == reflect.Bool {
            return argValue, nil
        }
        return reflect.Value{}, fmt.Errorf("cannot convert %s to bool", argType)

    case reflect.Int, reflect.Int32, reflect.Int64:
        switch argType.Kind() {
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
            return argValue.Convert(targetType), nil
        case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
            return argValue.Convert(targetType), nil
        case reflect.Float32, reflect.Float64:
            return argValue.Convert(targetType), nil
        default:
            return reflect.Value{}, fmt.Errorf("cannot convert %s to %s", argType, targetType)
        }

    case reflect.Float32, reflect.Float64:
        switch argType.Kind() {
        case reflect.Float32, reflect.Float64:
            return argValue.Convert(targetType), nil
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
            return argValue.Convert(targetType), nil
        case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
            return argValue.Convert(targetType), nil
        default:
            return reflect.Value{}, fmt.Errorf("cannot convert %s to %s", argType, targetType)
        }

    case reflect.String:
        if argType.Kind() == reflect.String {
            return argValue, nil
        }
        return reflect.Value{}, fmt.Errorf("cannot convert %s to string", argType)

    case reflect.Interface:
        // For non-empty interface{} targets, just pass the value
        return argValue, nil

    default:
        return reflect.Value{}, fmt.Errorf("unsupported target type %s", targetType)
    }
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `cd autoui && go test -run TestConvertArg_IntToEnum -v`
Expected: PASS

- [ ] **Step 6: Write test for any/interface{} parameter**

```go
// autoui/caller_test.go - add new test
func TestConvertArg_AnyParam(t *testing.T) {
    // interface{} target should accept any value
    targetType := reflect.TypeOf((*interface{})(nil)).Elem()
    arg := "test string"

    result, err := autoui.ConvertArg(arg, targetType)
    if err != nil {
        t.Errorf("ConvertArg failed: %v", err)
    }

    if result.String() != "test string" {
        t.Errorf("Expected 'test string', got %v", result.Interface())
    }
}
```

- [ ] **Step 7: Run test to verify it passes**

Run: `cd autoui && go test -run TestConvertArg_AnyParam -v`
Expected: PASS

- [ ] **Step 8: Commit**

```bash
git add autoui/caller.go autoui/caller_export_test.go autoui/caller_test.go
git commit -m "$(cat <<'EOF'
feat(autoui): extend convertArg for enum and interface{} conversion

- Handle int → enum conversion via reflect.ConvertibleTo
- Accept any value for interface{} parameters
- Simplify conversion logic using targetType.Kind() for underlying types

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
EOF
)"
```

---

### Task 3: Add Return Value Capture to InvokeMethod

**Files:**
- Modify: `autoui/caller.go:17-45`
- Modify: `autoui/handlers.go:68-91`
- Test: `autoui/caller_test.go`

- [ ] **Step 1: Write the failing test for return value capture**

```go
// autoui/caller_test.go - add new test
func TestInvokeMethod_ReturnValue(t *testing.T) {
    textInput := widget.NewTextInput(
        widget.TextInputOpts.WidgetOpts(
            widget.WidgetOpts.MinSize(200, 30),
        ),
    )
    textInput.GetWidget().Rect = image.Rect(10, 10, 210, 40)
    textInput.SetText("Hello")

    info := autoui.ExtractWidgetInfo(textInput)

    result, err := autoui.InvokeMethodWithResult(info.Widget, "GetText", nil)
    if err != nil {
        t.Errorf("InvokeMethodWithResult failed: %v", err)
    }

    if result != "Hello" {
        t.Errorf("Expected 'Hello', got %v", result)
    }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd autoui && go test -run TestInvokeMethod_ReturnValue -v`
Expected: FAIL - "InvokeMethodWithResult not defined"

- [ ] **Step 3: Add InvokeMethodWithResult function**

Modify `autoui/caller.go` - add new function after InvokeMethod:

```go
// InvokeMethodWithResult invokes a method and returns the result.
// Extends InvokeMethod to capture return values for getters.
func InvokeMethodWithResult(w widget.PreferredSizeLocateableWidget, methodName string, args []any) (any, error) {
    if w == nil {
        return nil, fmt.Errorf("widget is nil")
    }

    method := reflect.ValueOf(w).MethodByName(methodName)
    if !method.IsValid() {
        return nil, fmt.Errorf("method '%s' not found on widget type %T", methodName, w)
    }

    methodType := method.Type()

    if !isWhitelistedSignature(methodType) {
        return nil, fmt.Errorf("method '%s' has non-whitelisted signature %s", methodName, methodType)
    }

    expectedArgs := methodType.NumIn()
    if len(args) != expectedArgs {
        return nil, fmt.Errorf("method '%s' expects %d arguments, got %d", methodName, expectedArgs, len(args))
    }

    convertedArgs, err := convertArgs(args, methodType)
    if err != nil {
        return nil, fmt.Errorf("argument conversion failed: %w", err)
    }

    results := method.Call(convertedArgs)

    // Capture return value
    if len(results) > 0 {
        ret := results[0]

        // Check for error return
        if errVal, ok := ret.Interface().(error); ok && errVal != nil {
            return nil, fmt.Errorf("method '%s' returned error: %w", methodName, errVal)
        }

        // Convert enum types to underlying type for JSON serialization
        if ret.Kind() != reflect.Interface {
            switch ret.Kind() {
            case reflect.Int, reflect.Int32, reflect.Int64:
                return ret.Int(), nil
            case reflect.Float32, reflect.Float64:
                return ret.Float(), nil
            case reflect.Bool:
                return ret.Bool(), nil
            case reflect.String:
                return ret.String(), nil
            case reflect.Slice:
                return ret.Interface(), nil
            default:
                return ret.Interface(), nil
            }
        }
        return ret.Interface(), nil
    }

    return nil, nil
}
```

- [ ] **Step 4: Export InvokeMethodWithResult**

Add to `autoui/caller_export_test.go` (or just make it public - it's already public):

The function is already public, so tests can call it directly via `autoui.InvokeMethodWithResult`.

- [ ] **Step 5: Run test to verify it passes**

Run: `cd autoui && go test -run TestInvokeMethod_ReturnValue -v`
Expected: PASS

- [ ] **Step 6: Write test for enum return conversion**

```go
// autoui/caller_test.go - add new test
func TestInvokeMethodWithResult_EnumReturn(t *testing.T) {
    buttonImage := &widget.ButtonImage{
        Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
        Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
        Disabled: createTestNineSlice(100, 30, color.RGBA{150, 150, 150, 255}),
    }

    btn := widget.NewButton(
        widget.ButtonOpts.Image(buttonImage),
    )
    btn.GetWidget().Rect = image.Rect(10, 10, 110, 40)
    btn.Press() // Set to pressed state

    info := autoui.ExtractWidgetInfo(btn)

    result, err := autoui.InvokeMethodWithResult(info.Widget, "State", nil)
    if err != nil {
        t.Errorf("InvokeMethodWithResult failed: %v", err)
    }

    // State() returns WidgetState, should be converted to int
    stateInt, ok := result.(int64)
    if !ok {
        t.Errorf("Expected int64, got %T", result)
    }
    if stateInt != 1 { // WidgetStatePressed = 1
        t.Errorf("Expected state 1 (pressed), got %d", stateInt)
    }
}
```

- [ ] **Step 7: Run test to verify it passes**

Run: `cd autoui && go test -run TestInvokeMethodWithResult_EnumReturn -v`
Expected: PASS

- [ ] **Step 8: Write test for []any return**

```go
// autoui/caller_test.go - add new test
func TestInvokeMethodWithResult_SliceReturn(t *testing.T) {
    // Create a List with entries
    entries := []any{"Entry A", "Entry B", "Entry C"}
    list := widget.NewList(
        widget.ListOpts.Entries(entries),
        widget.ListOpts.EntryLabelFunc(func(e any) string {
            return e.(string)
        }),
    )
    list.GetWidget().Rect = image.Rect(10, 10, 210, 200)

    info := autoui.ExtractWidgetInfo(list)

    result, err := autoui.InvokeMethodWithResult(info.Widget, "Entries", nil)
    if err != nil {
        t.Errorf("InvokeMethodWithResult failed: %v", err)
    }

    // Should return []any
    entriesResult, ok := result.([]any)
    if !ok {
        t.Errorf("Expected []any, got %T", result)
    }
    if len(entriesResult) != 3 {
        t.Errorf("Expected 3 entries, got %d", len(entriesResult))
    }
}
```

- [ ] **Step 9: Run test to verify it passes**

Run: `cd autoui && go test -run TestInvokeMethodWithResult_SliceReturn -v`
Expected: PASS

- [ ] **Step 10: Commit**

```bash
git add autoui/caller.go autoui/caller_test.go
git commit -m "$(cat <<'EOF'
feat(autoui): add InvokeMethodWithResult for return value capture

- New function captures return values from widget methods
- Converts enum types to underlying int/float for JSON serialization
- Returns []any slices directly for JSON array output
- Handles error returns appropriately

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
EOF
)"
```

---

### Task 4: Update CallResponse and handleCallCommand

**Files:**
- Modify: `autoui/handlers.go:60-68`
- Modify: `autoui/handlers.go:79-91`
- Test: `autoui/handlers_test.go`

- [ ] **Step 1: Update CallResponse struct**

Modify `autoui/handlers.go` - update `CallResponse`:

```go
// CallResponse represents the response from a method invocation.
type CallResponse struct {
    // Success indicates if the invocation was successful.
    Success bool `json:"success"`

    // Error contains the error message if invocation failed.
    Error string `json:"error,omitempty"`

    // Result contains the captured return value (if method has a return).
    Result any `json:"result,omitempty"`
}
```

- [ ] **Step 2: Write test for CallResponse with Result**

```go
// autoui/handlers_test.go - add new test
func TestCallResponse_WithResult(t *testing.T) {
    resp := autoui.CallResponse{
        Success: true,
        Result:  "test value",
    }

    data, err := json.Marshal(resp)
    if err != nil {
        t.Errorf("Marshal failed: %v", err)
    }

    expected := `{"success":true,"result":"test value"}`
    if string(data) != expected {
        t.Errorf("Expected %s, got %s", expected, string(data))
    }
}
```

- [ ] **Step 3: Run test to verify it passes**

Run: `cd autoui && go test -run TestCallResponse_WithResult -v`
Expected: PASS

- [ ] **Step 4: Update handleCallCommand to use InvokeMethodWithResult**

Modify `autoui/handlers.go` - update `handleCallCommand`:

```go
// handleCallCommand handles the "call" command which invokes a method on a widget.
// Request format: `{"target":"query","method":"name","args":[...]}`.
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

        widgets := SnapshotTree(ui)

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

        targetWidget := matching[0]

        // Use InvokeMethodWithResult to capture return values
        result, err := InvokeMethodWithResult(targetWidget.Widget, callReq.Method, callReq.Args)

        response := CallResponse{
            Success: err == nil,
        }
        if err != nil {
            response.Error = err.Error()
        }
        if result != nil {
            response.Result = result
        }

        respData, err := json.Marshal(response)
        if err != nil {
            ctx.Respond("error: failed to marshal response: " + err.Error())
            return
        }

        ctx.Respond(string(respData))
    }
}
```

- [ ] **Step 5: Run all handler tests to verify no breakage**

Run: `cd autoui && go test -run TestHandleCall -v`
Expected: PASS (existing tests)

- [ ] **Step 6: Commit**

```bash
git add autoui/handlers.go autoui/handlers_test.go
git commit -m "$(cat <<'EOF'
feat(autoui): add Result field to CallResponse and capture returns

- CallResponse now includes optional Result field for return values
- handleCallCommand uses InvokeMethodWithResult instead of InvokeMethod
- Return values serialized to JSON in response

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
EOF
)"
```

---

### Task 5: Create Proxy Handlers (proxy.go)

**Files:**
- Create: `autoui/proxy.go`
- Create: `autoui/proxy_test.go`

- [ ] **Step 1: Create proxy.go with registry structure**

```go
// autoui/proxy.go
package autoui

import (
    "fmt"
    "reflect"

    "github.com/ebitenui/ebitenui/widget"
)

// ProxyHandler handles a named operation on a widget, bypassing the whitelist.
// Returns (result, error) where result is nil for setters.
type ProxyHandler func(w widget.PreferredSizeLocateableWidget, args []any) (any, error)

// proxyHandlers maps proxy method names to their handlers.
var proxyHandlers = map[string]ProxyHandler{
    "SelectEntryByIndex": handleSelectEntryByIndex,
    "SelectedEntryIndex": handleSelectedEntryIndex,
}

// GetProxyHandler returns the handler for a proxy method name, or nil if not found.
func GetProxyHandler(name string) ProxyHandler {
    return proxyHandlers[name]
}
```

- [ ] **Step 2: Write failing test for SelectEntryByIndex**

```go
// autoui/proxy_test.go
package autoui_test

import (
    "image"
    "image/color"
    "testing"

    "github.com/ebitenui/ebitenui/widget"
    "github.com/s3cy/autoebiten/autoui"
)

func TestProxy_SelectEntryByIndex(t *testing.T) {
    entries := []any{"Entry A", "Entry B", "Entry C"}
    list := widget.NewList(
        widget.ListOpts.Entries(entries),
        widget.ListOpts.EntryLabelFunc(func(e any) string {
            return e.(string)
        }),
    )
    list.GetWidget().Rect = image.Rect(10, 10, 210, 200)

    info := autoui.ExtractWidgetInfo(list)

    // Select entry at index 1
    result, err := autoui.InvokeMethodWithResult(info.Widget, "SelectEntryByIndex", []any{float64(1)})
    if err != nil {
        t.Errorf("SelectEntryByIndex failed: %v", err)
    }

    // Verify selection changed
    selected := list.SelectedEntry()
    if selected != entries[1] {
        t.Errorf("Expected 'Entry B', got %v", selected)
    }
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `cd autoui && go test -run TestProxy_SelectEntryByIndex -v`
Expected: FAIL - proxy not integrated yet

- [ ] **Step 4: Add handleSelectEntryByIndex implementation**

Add to `autoui/proxy.go`:

```go
// handleSelectEntryByIndex selects a list entry by its position index.
// Works for widget.List.
func handleSelectEntryByIndex(w widget.PreferredSizeLocateableWidget, args []any) (any, error) {
    if len(args) != 1 {
        return nil, fmt.Errorf("SelectEntryByIndex requires 1 argument (index)")
    }

    index, ok := args[0].(float64)
    if !ok {
        return nil, fmt.Errorf("index must be a number, got %T", args[0])
    }
    idx := int(index)

    // Try List
    if list, ok := w.(*widget.List); ok {
        entries := list.Entries()
        if idx < 0 || idx >= len(entries) {
            return nil, fmt.Errorf("index %d out of range (0-%d)", idx, len(entries)-1)
        }
        list.SetSelectedEntry(entries[idx])
        return nil, nil
    }

    return nil, fmt.Errorf("SelectEntryByIndex requires widget of type *List, got %T", w)
}
```

- [ ] **Step 5: Integrate proxy check into InvokeMethodWithResult**

Modify `autoui/caller.go` - update `InvokeMethodWithResult` to check proxy first:

```go
func InvokeMethodWithResult(w widget.PreferredSizeLocateableWidget, methodName string, args []any) (any, error) {
    if w == nil {
        return nil, fmt.Errorf("widget is nil")
    }

    // Check proxy registry first
    if handler := GetProxyHandler(methodName); handler != nil {
        return handler(w, args)
    }

    // Use reflection for whitelisted methods
    method := reflect.ValueOf(w).MethodByName(methodName)
    // ... rest of existing implementation ...
}
```

- [ ] **Step 6: Run test to verify it passes**

Run: `cd autoui && go test -run TestProxy_SelectEntryByIndex -v`
Expected: PASS

- [ ] **Step 7: Write test for SelectedEntryIndex**

```go
// autoui/proxy_test.go - add new test
func TestProxy_SelectedEntryIndex(t *testing.T) {
    entries := []any{"Entry A", "Entry B", "Entry C"}
    list := widget.NewList(
        widget.ListOpts.Entries(entries),
        widget.ListOpts.EntryLabelFunc(func(e any) string {
            return e.(string)
        }),
    )
    list.GetWidget().Rect = image.Rect(10, 10, 210, 200)
    list.SetSelectedEntry(entries[1]) // Pre-select Entry B

    info := autoui.ExtractWidgetInfo(list)

    result, err := autoui.InvokeMethodWithResult(info.Widget, "SelectedEntryIndex", nil)
    if err != nil {
        t.Errorf("SelectedEntryIndex failed: %v", err)
    }

    index, ok := result.(int64)
    if !ok {
        t.Errorf("Expected int64, got %T", result)
    }
    if index != 1 {
        t.Errorf("Expected index 1, got %d", index)
    }
}
```

- [ ] **Step 8: Add handleSelectedEntryIndex implementation**

Add to `autoui/proxy.go`:

```go
// handleSelectedEntryIndex returns the index of the currently selected entry.
// Returns -1 if no entry is selected.
func handleSelectedEntryIndex(w widget.PreferredSizeLocateableWidget, args []any) (any, error) {
    if len(args) != 0 {
        return nil, fmt.Errorf("SelectedEntryIndex requires no arguments")
    }

    if list, ok := w.(*widget.List); ok {
        selected := list.SelectedEntry()
        entries := list.Entries()
        for i, e := range entries {
            if e == selected {
                return int64(i), nil
            }
        }
        return int64(-1), nil // No selection
    }

    return nil, fmt.Errorf("SelectedEntryIndex requires widget of type *List, got %T", w)
}
```

- [ ] **Step 9: Run test to verify it passes**

Run: `cd autoui && go test -run TestProxy_SelectedEntryIndex -v`
Expected: PASS

- [ ] **Step 10: Write test for out of bounds index**

```go
// autoui/proxy_test.go - add new test
func TestProxy_SelectEntryByIndex_OutOfBounds(t *testing.T) {
    entries := []any{"A", "B"}
    list := widget.NewList(
        widget.ListOpts.Entries(entries),
        widget.ListOpts.EntryLabelFunc(func(e any) string { return e.(string) }),
    )
    list.GetWidget().Rect = image.Rect(10, 10, 210, 200)

    info := autoui.ExtractWidgetInfo(list)

    _, err := autoui.InvokeMethodWithResult(info.Widget, "SelectEntryByIndex", []any{float64(5)})
    if err == nil {
        t.Error("Expected error for out of bounds index")
    }
    if !strings.Contains(err.Error(), "out of range") {
        t.Errorf("Expected 'out of range' error, got: %v", err)
    }
}
```

- [ ] **Step 11: Run test to verify it passes**

Run: `cd autoui && go test -run TestProxy_SelectEntryByIndex_OutOfBounds -v`
Expected: PASS

- [ ] **Step 12: Write test for wrong widget type**

```go
// autoui/proxy_test.go - add new test
func TestProxy_SelectEntryByIndex_WrongWidgetType(t *testing.T) {
    buttonImage := &widget.ButtonImage{
        Idle: createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
    }
    btn := widget.NewButton(widget.ButtonOpts.Image(buttonImage))
    btn.GetWidget().Rect = image.Rect(10, 10, 110, 40)

    info := autoui.ExtractWidgetInfo(btn)

    _, err := autoui.InvokeMethodWithResult(info.Widget, "SelectEntryByIndex", []any{float64(0)})
    if err == nil {
        t.Error("Expected error for wrong widget type")
    }
    if !strings.Contains(err.Error(), "requires widget of type *List") {
        t.Errorf("Expected type error, got: %v", err)
    }
}
```

- [ ] **Step 13: Run test to verify it passes**

Run: `cd autoui && go test -run TestProxy_SelectEntryByIndex_WrongWidgetType -v`
Expected: PASS

- [ ] **Step 14: Commit**

```bash
git add autoui/proxy.go autoui/proxy_test.go autoui/caller.go
git commit -m "$(cat <<'EOF'
feat(autoui): add proxy handlers for List selection by index

- SelectEntryByIndex(int): selects entry at position
- SelectedEntryIndex(): returns current selection position
- Proxy registry checked before reflection in InvokeMethodWithResult
- Bounds checking and type validation

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
EOF
)"
```

---

### Task 6: Integration Test for Full Workflow

**Files:**
- Modify: `autoui/integration_test.go`

- [ ] **Step 1: Write integration test for List selection workflow**

```go
// autoui/integration_test.go - add new test
func TestIntegration_ListSelectionWorkflow(t *testing.T) {
    // Create a List with string entries
    entries := []any{"Option A", "Option B", "Option C"}
    list := widget.NewList(
        widget.ListOpts.Entries(entries),
        widget.ListOpts.EntryLabelFunc(func(e any) string {
            return e.(string)
        }),
    )
    list.GetWidget().Rect = image.Rect(10, 10, 210, 200)

    info := autoui.ExtractWidgetInfo(list)

    // Step 1: Get entries
    result, err := autoui.InvokeMethodWithResult(info.Widget, "Entries", nil)
    if err != nil {
        t.Fatalf("Entries() failed: %v", err)
    }
    entriesResult := result.([]any)
    if len(entriesResult) != 3 {
        t.Errorf("Expected 3 entries, got %d", len(entriesResult))
    }

    // Step 2: Select entry by index
    _, err = autoui.InvokeMethodWithResult(info.Widget, "SelectEntryByIndex", []any{float64(2)})
    if err != nil {
        t.Fatalf("SelectEntryByIndex(2) failed: %v", err)
    }

    // Step 3: Verify selection index
    result, err = autoui.InvokeMethodWithResult(info.Widget, "SelectedEntryIndex", nil)
    if err != nil {
        t.Fatalf("SelectedEntryIndex() failed: %v", err)
    }
    if result.(int64) != 2 {
        t.Errorf("Expected index 2, got %d", result)
    }

    // Step 4: Verify selection value via SelectedEntry
    result, err = autoui.InvokeMethodWithResult(info.Widget, "SelectedEntry", nil)
    if err != nil {
        t.Fatalf("SelectedEntry() failed: %v", err)
    }
    if result != entries[2] {
        t.Errorf("Expected 'Option C', got %v", result)
    }
}
```

- [ ] **Step 2: Run integration test**

Run: `cd autoui && go test -run TestIntegration_ListSelectionWorkflow -v`
Expected: PASS

- [ ] **Step 3: Write integration test for WidgetState workflow**

```go
// autoui/integration_test.go - add new test
func TestIntegration_WidgetStateWorkflow(t *testing.T) {
    buttonImage := &widget.ButtonImage{
        Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
        Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
        Disabled: createTestNineSlice(100, 30, color.RGBA{150, 150, 150, 255}),
    }

    btn := widget.NewButton(widget.ButtonOpts.Image(buttonImage))
    btn.GetWidget().Rect = image.Rect(10, 10, 110, 40)

    info := autoui.ExtractWidgetInfo(btn)

    // Step 1: Get initial state (should be idle = 0)
    result, err := autoui.InvokeMethodWithResult(info.Widget, "State", nil)
    if err != nil {
        t.Fatalf("State() failed: %v", err)
    }
    if result.(int64) != 0 {
        t.Errorf("Expected initial state 0 (idle), got %d", result)
    }

    // Step 2: Set state to pressed (1)
    _, err = autoui.InvokeMethodWithResult(info.Widget, "SetState", []any{float64(1)})
    if err != nil {
        t.Fatalf("SetState(1) failed: %v", err)
    }

    // Step 3: Verify state changed
    result, err = autoui.InvokeMethodWithResult(info.Widget, "State", nil)
    if err != nil {
        t.Fatalf("State() failed: %v", err)
    }
    if result.(int64) != 1 {
        t.Errorf("Expected state 1 (pressed), got %d", result)
    }
}
```

- [ ] **Step 4: Run integration test**

Run: `cd autoui && go test -run TestIntegration_WidgetStateWorkflow -v`
Expected: PASS

- [ ] **Step 5: Run all tests**

Run: `cd autoui && go test -v`
Expected: All PASS

- [ ] **Step 6: Commit**

```bash
git add autoui/integration_test.go
git commit -m "$(cat <<'EOF'
test(autoui): add integration tests for selection and state workflows

- List selection: Entries → SelectEntryByIndex → SelectedEntryIndex
- WidgetState: State → SetState → State verification

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
EOF
)"
```

---

### Task 7: Final Verification and Documentation

**Files:**
- None (verification only)

- [ ] **Step 1: Run all autoui tests**

Run: `cd autoui && go test -v`
Expected: All PASS

- [ ] **Step 2: Run full project tests**

Run: `go test ./...`
Expected: All PASS

- [ ] **Step 3: Verify no regressions in existing tests**

Run: `cd autoui && go test -run TestInvokeMethod -v`
Expected: All existing InvokeMethod tests PASS

- [ ] **Step 4: Final commit with any remaining changes**

```bash
git status
# If clean, no action needed
git add -A && git commit -m "chore: final cleanup for autoui extended type support" || echo "Already clean"
```