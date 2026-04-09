---
name: autoui-exists-command
description: Add autoui.exists command for widget waiting with wait-for
type: project
---

# autoui.exists Command Design

**Date:** 2026-04-09

## Problem

E2E tests need to wait for widgets to appear dynamically (dialogs, loading screens, async content). The existing `wait-for` command requires JSON responses, but autoui commands return XML.

## Solution

Add `autoui.exists` command that returns JSON instead of XML, making it compatible with `wait-for`.

## Design

### autoui.exists

**Purpose:** Check widget existence, returning JSON for use with `wait-for`.

**Command format:**
```bash
autoebiten custom autoui.exists --request "type=Dialog"
autoebiten custom autoui.exists --request '{"type":"TextInput","id":"name-input"}'
```

**Request parsing:** Same as `autoui.find` - supports both `key=value` and JSON format.

**Response format:**
```json
{"found":true,"count":1}
```
or
```json
{"found":false,"count":0}
```

**No error on not found:** Unlike `autoui.find`, this returns `{found:false,count:0}` instead of an error message. This allows use with `wait-for` which expects valid JSON.

### Usage with wait-for

**CLI:**
```bash
autoebiten wait-for 'custom:autoui.exists:type=Dialog == {"found":true}' --timeout 5s
```

**testkit:**
```go
game.WaitFor('custom:autoui.exists:type=Dialog == {"found":true}', "5s", "")
```

## Implementation

### Files to modify

1. `autoui/handlers.go`
   - Add `ExistsResponse` struct
   - Add `handleExistsCommand` handler function

2. `autoui/register.go`
   - Register `autoui.exists` command

3. `docs/autoui.md`
   - Add `autoui.exists` command section
   - Add SetText example under `autoui.call`
   - Add `autoui.exists` + `wait-for` usage section

### Implementation details

**ExistsResponse struct:**
```go
type ExistsResponse struct {
    Found bool `json:"found"`
    Count int  `json:"count"`
}
```

**handleExistsCommand:**
- Parse request using existing query parsing (key=value or JSON)
- Walk widget tree
- Filter widgets matching query
- Return JSON response (no error for empty results)

## Documentation Updates

### docs/autoui.md additions

1. **Quick Decision section:** Add `autoui.exists` entry
2. **Commands section:** New `autoui.exists` subsection with:
   - Purpose
   - Usage examples
   - Response format
3. **autoui.call section:** Add SetText example for TextInput
4. **New section:** "Waiting for Widgets" showing `autoui.exists` + `wait-for`

## SetText Example for autoui.call

```bash
# Set text in TextInput widget
autoebiten custom autoui.call --request '{"target":"id=name-input","method":"SetText","args":["Alice"]}'
```

**Response:**
```json
{"success":true}
```

**Why:** TextInput's `SetText(string)` method triggers `ChangedEvent` for game logic. Direct field setting would bypass events.

## Why Not autoui.type

Analysis of ebitenui setters showed:
- TextInput.SetText() fires ChangedEvent
- Checkbox.SetState() fires StateChangedEvent
- Button.SetText() updates internal state

Generic `autoui.set` would bypass these events. `autoui.call` with explicit method names is safer and already works.

## Test Coverage

- Unit tests for `handleExistsCommand` in `autoui/handlers_test.go`
- Integration test showing `autoui.exists` + `wait-for` in `autoui/integration_test.go`