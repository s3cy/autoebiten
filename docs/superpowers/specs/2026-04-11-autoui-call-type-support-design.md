# autoui.call Extended Type Support Design

**Date**: 2026-04-11
**Package**: `autoui` - EbitenUI automation helper for autoebiten

## Overview

Extend `autoui.call` to support `any`/`interface{}` types and enum types, enabling interaction with List, ListComboButton, SelectComboButton selection methods and WidgetState/Visibility enum methods.

This design focuses on High and Medium priority blocked methods. RadioGroup, TabBook, and List content operations (AddEntry, RemoveEntry, SetEntries, UpdateEntry) are reserved for future iterations.

## Problem Statement

### Current Limitation

The `InvokeMethod` function in `autoui/caller.go` uses a strict whitelist that blocks:

1. `any`/`interface{}` parameters and return values
2. Custom types with non-empty `PkgPath()` (including enum types like `widget.WidgetState`)
3. Slice return types (e.g., `[]any`)
4. Multiple parameters

### Blocked Methods (High/Medium Priority)

| Widget | Method | Signature | Blocked Reason |
|--------|--------|-----------|----------------|
| List | `SelectedEntry() any` | `func() any` | `any` return not captured |
| List | `Entries() []any` | `func() []any` | Slice return not captured |
| List | `SetSelectedEntry(any)` | `func(any)` | `any` param not whitelisted |
| ListComboButton | `SelectedEntry() interface{}` | `func() interface{}` | `interface{}` return |
| ListComboButton | `SetSelectedEntry(interface{})` | `func(interface{})` | `interface{}` param |
| SelectComboButton | `SelectedEntry() interface{}` | `func() interface{}` | `interface{}` return |
| SelectComboButton | `SetSelectedEntry(interface{})` | `func(interface{})` | `interface{}` param |
| Button | `State() WidgetState` | `func(widget.WidgetState)` | Custom type return |
| Button | `SetState(WidgetState)` | `func(widget.WidgetState)` | Custom type param |
| Checkbox | `State() WidgetState` | `func(widget.WidgetState)` | Custom type return |
| Checkbox | `SetState(WidgetState)` | `func(widget.WidgetState)` | Custom type param |
| Widget | `GetVisibility() Visibility` | `func() widget.Visibility` | Custom type return |
| Widget | `SetVisibility(Visibility)` | `func(widget.Visibility)` | Custom type param |

### Scope

**In scope**:
- Selection widgets: List, ListComboButton, SelectComboButton
- WidgetState enum: Button, Checkbox state control
- Visibility enum: Widget visibility control

**Out of scope (future iterations)**:
- RadioGroup: `Active()`, `SetActive()` - requires Container traversal
- TabBook: `Tab()`, `SetTab()` - requires TabBookTab handling
- List content: `AddEntry()`, `RemoveEntry()`, `SetEntries()`, `UpdateEntry()`

## Design Solution

### Key Insight: Enum Type Reflection

Go's reflection reveals that enum types like `WidgetState int` have:

- `Kind() == reflect.Int` (underlying type)
- `ConvertibleTo(int) == true`
- Can convert int → enum via `reflect.Value.Convert()`

This enables treating enum types as integers for JSON serialization/deserialization.

### Extended Whitelist

**Return values** - enable:
1. `any`/`interface{}` returns → serialize to JSON
2. `[]any` returns → serialize to JSON array
3. Enum types (underlying `int/float/string/bool`) → return as underlying type

**Input parameters** - enable:
1. `any`/`interface{}` parameters → accept JSON value
2. Enum types (underlying `int/float/string/bool`) → convert from basic type
3. Slice parameters → still blocked (reserved for future)

### Proxy Handlers (Minimal Set)

Two proxy handlers for selection by index:

**`SelectEntryByIndex(int)`**:
- Calls `Entries()` to get entry list
- Validates index bounds (0 to len(entries)-1)
- Calls `SetSelectedEntry(entries[index])`
- Works for List

**`SelectedEntryIndex() int`**:
- Calls `SelectedEntry()` to get current selection
- Calls `Entries()` to get entry list
- Finds index by comparing entries
- Returns index (-1 if no selection)
- Works for List

## Architecture & Data Flow

```
autoui.call Request
    {"target":"type=List","method":"SelectEntryByIndex","args":[2]}
                        │
                        ▼
            handleCallCommand
              1. Parse JSON request
              2. Find target widget via Finder
              3. Call InvokeMethod
                        │
                        ▼
                InvokeMethod
    ┌─────────────────────────────────────────┐
    │ Step 1: Check proxy registry             │
    │   If found: execute proxy, return result │
    └─────────────────────────────────────────┘
    ┌─────────────────────────────────────────┐
    │ Step 2: Use reflection (if no proxy)     │
    │   a. isWhitelistedSignature(methodType)  │
    │      - Allow any/interface{}             │
    │      - Allow enum types                  │
    │   b. convertArgs(args, methodType)       │
    │   c. method.Call(convertedArgs)          │
    │   d. captureReturn(results[0])           │
    └─────────────────────────────────────────┘
                        │
                        ▼
                CallResponse
    {"success":true,"result":null}
    {"success":true,"result":2}
    {"success":true,"result":["A","B","C"]}
```

## Implementation Details

### Modified Files

- `autoui/caller.go`: Extend whitelist, add return capture
- `autoui/handlers.go`: Add Result field to CallResponse

### New File

- `autoui/proxy.go`: Proxy handler implementations

## CLI Usage Examples

### Selection Operations (List)

```bash
# Get entries
autoui.call --request '{"target":"type=List","method":"Entries"}'
# {"success":true,"result":["Option A","Option B","Option C"]}

# Select by index
autoui.call --request '{"target":"type=List","method":"SelectEntryByIndex","args":[2]}'

# Verify selection index
autoui.call --request '{"target":"type=List","method":"SelectedEntryIndex"}'
# {"success":true,"result":2}
```

### WidgetState Operations (Button/Checkbox)

```bash
# Get button state (enum → int)
autoui.call --request '{"target":"type=Button","method":"State"}'
# {"success":true,"result":1}

# Set button state (int → enum)
autoui.call --request '{"target":"type=Button","method":"SetState","args":[2]}'
```

### Visibility Operations (Widget)

```bash
# Get visibility
autoui.call --request '{"target":"id=my-widget","method":"GetVisibility"}'
# {"success":true,"result":0}

# Set visibility
autoui.call --request '{"target":"id=my-widget","method":"SetVisibility","args":[1]}'
```

## Testing Strategy

- Unit tests: extended whitelist, proxy handlers
- Integration tests: full workflows

## Future Enhancements

- RadioGroup support
- TabBook support
- List content operations
- String-based enum proxies
