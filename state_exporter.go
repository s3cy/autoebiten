package autoebiten

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var StateExporterPathPrefix = ".state."

// ErrPathNotFound is returned when a state query path cannot be resolved.
// This occurs when:
//   - The path references a non-existent field
//   - An array index is out of bounds
//   - A map key does not exist
//   - The path traverses a nil pointer
var ErrPathNotFound = errors.New("path not found")

/*
RegisterStateExporter registers a custom command that exposes game
state via reflection-based path queries. The path uses dot notation:
  - "Player.X" - access struct field
  - "Inventory.0.Name" - access array/slice index
  - "Skills.Sword" - access map key

# Game state can be retrieved by calling the testkit.StateQuery function

In game:

	type Player struct {
		X      float64
		Y      float64
		Health int
	}

	type Game struct {
		Player Player
	}

	func main() {
		var g := Game{Player: Player{X: 0, Y: 0, Health: 100}}
		RegisterStateExporter("mystate", &g)
	}

In test:

	func TestPlayerMovement(t *testing.T) {
	    game := testkit.Launch(t, "./mygame")
	    defer game.Shutdown()

		game.HoldKey(ebiten.KeyD, 10)
	    x, _ := game.StateQuery("mystate", "Player.X")
	    assert.Equal(t, 10, x)
	}
*/
func RegisterStateExporter(name string, root any) {
	Register(StateExporterPathPrefix+name, stateExporter(root))
}

func stateExporter(root any) func(CommandContext) {
	return func(ctx CommandContext) {
		path := ctx.Request()
		if path == "" {
			ctx.Respond(`{"error":"path required"}`)
			return
		}

		value, err := navigatePath(root, path)
		if err != nil {
			if err == ErrPathNotFound {
				ctx.Respond(`{"error":"path not found"}`)
				return
			}
			ctx.Respond(fmt.Sprintf(`{"error":"%s"}`, err.Error()))
			return
		}

		// Encode result as JSON
		data, err := json.Marshal(value)
		if err != nil {
			ctx.Respond(fmt.Sprintf(`{"error":"failed to marshal: %s"}`, err.Error()))
			return
		}

		ctx.Respond(string(data))
	}
}

// navigatePath navigates a value using dot-notation path.
func navigatePath(root any, path string) (any, error) {
	parts := parsePath(path)
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty path")
	}

	v := reflect.ValueOf(root)

	// Handle pointer at root
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, ErrPathNotFound
		}
		v = v.Elem()
	}

	for _, part := range parts {
		if !v.IsValid() {
			return nil, ErrPathNotFound
		}

		switch v.Kind() {
		case reflect.Ptr:
			if v.IsNil() {
				return nil, ErrPathNotFound
			}
			v = v.Elem()
			// Re-process this part with the dereferenced value
			v, _ = getFieldByName(v, part)

		case reflect.Struct:
			field, err := getFieldByName(v, part)
			if err != nil {
				return nil, err
			}
			v = field

		case reflect.Map:
			key := reflect.ValueOf(part)
			val := v.MapIndex(key)
			if !val.IsValid() {
				return nil, ErrPathNotFound
			}
			v = val

		case reflect.Slice, reflect.Array:
			index, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid array index: %s", part)
			}
			if index < 0 || index >= v.Len() {
				return nil, ErrPathNotFound
			}
			v = v.Index(index)

		default:
			return nil, ErrPathNotFound
		}
	}

	if !v.IsValid() {
		return nil, ErrPathNotFound
	}

	// Convert reflect.Value to interface{}
	return v.Interface(), nil
}

// parsePath splits a dot-notation path into components.
// Handles escaped dots (\.) as literal dots in field names.
func parsePath(path string) []string {
	var parts []string
	var current strings.Builder
	escaped := false

	for _, r := range path {
		if escaped {
			current.WriteRune(r)
			escaped = false
			continue
		}

		if r == '\\' {
			escaped = true
			continue
		}

		if r == '.' {
			parts = append(parts, current.String())
			current.Reset()
			continue
		}

		current.WriteRune(r)
	}

	// Add final part
	if current.Len() > 0 || len(parts) > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

// getFieldByName gets a struct field by name, handling case-insensitive lookup.
func getFieldByName(v reflect.Value, name string) (reflect.Value, error) {
	t := v.Type()

	// Try exact match first
	field, ok := t.FieldByName(name)
	if ok {
		return v.FieldByIndex(field.Index), nil
	}

	// Try case-insensitive match
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if strings.EqualFold(field.Name, name) {
			return v.Field(i), nil
		}
	}

	return reflect.Value{}, ErrPathNotFound
}

// StateQueryPath represents a parsed state query path.
// It can be used for repeated queries against different roots.
type StateQueryPath struct {
	parts []string
}

// ParseStateQueryPath parses a state query path for reuse.
func ParseStateQueryPath(path string) (*StateQueryPath, error) {
	parts := parsePath(path)
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty path")
	}
	return &StateQueryPath{parts: parts}, nil
}

// Query executes the parsed path against a root value.
func (p *StateQueryPath) Query(root any) (any, error) {
	return navigatePath(root, strings.Join(p.parts, "."))
}
