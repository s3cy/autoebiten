# CLI State Query and WaitFor Feature Design

Date: 2026-04-06

## Overview

Add `state` and `wait-for` commands to the CLI, mirroring testkit's `StateQuery` and `WaitFor` capabilities. These commands enable AI agents and scripts to query game state and wait for conditions to be met, supporting more complex automation workflows.

## CLI Commands

### `state` command

Query game state via registered exporters.

```
autoebiten state --name <exporter_name> --path <dot_path>
```

**Flags:**
- `--name` (required): State exporter name registered via `autoebiten.RegisterStateExporter`
- `--path` (required): Dot-notation path (e.g., `Player.X`, `Inventory.0.Name`)

**Output:** `OK: <value>` (wraps result with success prefix, consistent with other commands)

**Examples:**

```bash
autoebiten state --name gamestate --path Player.Health
# Output: OK: 100

autoebiten state --name gamestate --path Player.Position
# Output: OK: {"X":10,"Y":20}

autoebiten state --name gamestate --path Inventory.0.Name
# Output: OK: "Sword"
```

### `wait-for` command

Poll until condition matches or timeout expires.

```
autoebiten wait-for --condition "<condition>" --timeout <duration> [--interval <duration>]
```

**Flags:**
- `--condition` (required): Condition string
- `--timeout` (required): Maximum wait duration (e.g., `10s`, `5m`)
- `--interval` (optional): Poll interval (default `100ms`)

**Output on success:** `OK: condition met after <duration>` (e.g., `OK: condition met after 2.3s`)

**Output on timeout:** Error message with timeout duration

**Condition Format:**

```
<type>:<name>:<path> <operator> <value>
```

- `type`: `state` or `custom`
- `name`: Exporter name or custom command name
- `path`: Dot-notation path (for state type) or request string passed to custom command (for custom type)
- `operator`: `==`, `!=`, `<`, `>`, `<=`, `>=`
- `value`: JSON value (number, string, boolean)

**Operators:**

| Operator | Description |
|----------|-------------|
| `==` | Equal |
| `!=` | Not equal |
| `<` | Less than (numbers only) |
| `>` | Greater than (numbers only) |
| `<=` | Less than or equal (numbers only) |
| `>=` | Greater than or equal (numbers only) |

**Examples:**

```bash
autoebiten wait-for --condition "state:gamestate:Player.Health == 100" --timeout 10s
# Output: OK: condition met after 2.3s

autoebiten wait-for --condition "state:gamestate:Player.X > 50" --timeout 30s
# Output: OK: condition met after 15.2s

autoebiten wait-for --condition "custom:getStatus:ready != false" --timeout 5s --interval 200ms
# Output: OK: condition met after 1.5s
```

## JSON Script Support

Add `state` and `wait` commands to script format.

**Schema Additions:**

```json
{
  "version": "1.0",
  "commands": [
    {"input": {"action": "press", "key": "KeyW"}},
    {"state": {"name": "gamestate", "path": "Player.X"}},
    {"wait": {"condition": "state:gamestate:Player.Health == 100", "timeout": "10s"}},
    {"screenshot": {"output": "shot.png"}}
  ]
}
```

**StateCmd:**

```go
type StateCmd struct {
    Name string `json:"name"`
    Path string `json:"path"`
}
```

**WaitCmd:**

```go
type WaitCmd struct {
    Condition string `json:"condition"`
    Timeout   string `json:"timeout"`
    Interval  string `json:"interval,omitempty"` // default: "100ms"
}
```

**Execution Behavior:**

- `state`: Execute query, continue to next command (no blocking)
- `wait`: Poll until condition met or timeout, block subsequent commands

## Architecture

### File Structure

**New files:**

```
internal/cli/state.go      - State command implementation
internal/cli/wait.go       - WaitFor command implementation
```

**Modified files:**

```
cmd/autoebiten/main.go     - Add state and wait-for command definitions
internal/script/ast.go     - Add StateCmd and WaitCmd structs, update CommandSchema
internal/script/executor.go - Add state and wait command handlers
```

### Implementation Details

**`internal/cli/state.go`:**

```go
func (e *CommandExecutor) RunStateCommand(name, path string) error {
    // Internally calls custom ".state.<name>" with path as request
    customName := autoebiten.StateExporterPathPrefix + name
    
    // Get response and wrap with OK: prefix using writer.Success
    params := &rpc.CustomParams{Name: customName, Request: path}
    req, _ := rpc.BuildRequest("custom", params)
    resp, err := rpc.SendRequestSocket(req)
    if err != nil { return err }
    if resp.Error != nil { return fmt.Errorf("rpc error: %s", resp.Error.Message) }
    
    var result rpc.CustomResult
    json.Unmarshal(resp.Result, &result)
    e.writer.Success(result.Response)
    return nil
}
```

**`internal/cli/wait.go`:**

```go
type Condition struct {
    Type     string   // "state" or "custom"
    Name     string   // exporter or command name
    Path     string   // path or request
    Operator string   // ==, !=, <, >, <=, >=
    Value    any      // expected value (parsed from JSON)
}

func ParseCondition(s string) (*Condition, error) {
    // Parse "type:name:path operator value"
}

func CheckCondition(queried any, op string, expected any) (bool, error) {
    // Compare values based on operator
}

func (e *CommandExecutor) RunWaitForCommand(condition, timeout, interval string) error {
    // Parse condition
    // Create timeout context
    // Poll loop: query -> check -> sleep(interval) until met or timeout
}
```

**Condition Parsing Logic:**

1. Split string on first space → query part + operator/value part
2. Split query part on `:` → type, name, path
3. Find operator in remaining part (==, !=, <, >, <=, >=)
4. Parse value as JSON (number, string, bool)

**Condition Checking Logic:**

| Queried Type | Supported Operators |
|--------------|---------------------|
| Number | All (`==`, `!=`, `<`, `>`, `<=`, `>=`) |
| String | `==`, `!=` only |
| Boolean | `==`, `!=` only |
| Object/Array | Unsupported (error) |

**`internal/script/executor.go`:**

```go
func (e *Executor) executeCommand(cmd CommandSchema) error {
    switch {
    case cmd.State != nil:
        return e.stateFunc(cmd.State.Name, cmd.State.Path)
    case cmd.Wait != nil:
        return e.waitFunc(cmd.Wait.Condition, cmd.Wait.Timeout, cmd.Wait.Interval)
    // ... existing cases
    }
}

func (e *Executor) SetStateFunc(fn func(name, path string) error) {
    e.stateFunc = fn
}

func (e *Executor) SetWaitFunc(fn func(condition, timeout, interval string) error) {
    e.waitFunc = fn
}
```

## RPC Protocol

No changes required. Both commands use existing RPC infrastructure:

- `state`: Uses `custom` method with `.state.<name>` prefix
- `wait-for`: CLI-side polling using `custom` method

## Error Handling

| Scenario | Error Message |
|----------|---------------|
| Invalid condition format | `invalid condition format: expected "type:name:path operator value"` |
| Unsupported operator for value type | `operator < not supported for string values` |
| Timeout exceeded | `timeout exceeded after 10s waiting for condition` |
| Query returns object/array for comparison | `cannot compare objects/arrays, use specific path to a primitive value` |
| State exporter not found | `state exporter "gamestate" not found` |

## Testing

**Unit tests:**
- `internal/cli/state_test.go`: State command execution
- `internal/cli/wait_test.go`: Condition parsing, condition checking, wait-for command
- `internal/script/ast_test.go`: StateCmd and WaitCmd parsing
- `internal/script/executor_test.go`: Script execution with state/wait commands

**Integration tests:**
- State query against running game with registered exporter
- Wait-for with various conditions and timeouts
- Script execution combining state/wait with other commands

## Success Criteria

1. `autoebiten state --name <n> --path <p>` outputs `OK: <value>`
2. `autoebiten wait-for --condition <c> --timeout <t>` polls until condition met
3. Scripts support `state` and `wait` commands
4. All comparison operators work correctly for supported value types
5. Timeout handling with clear error messages
6. 80%+ test coverage for new code