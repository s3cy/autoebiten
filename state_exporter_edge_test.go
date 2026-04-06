package autoebiten

import (
	"testing"
)

// Test struct with JSON tags
type taggedPlayer struct {
	Name   string `json:"player_name"`
	Health int    `json:"hp"`
	X      float64
}

// Test struct with interface field
type gameWithInterface struct {
	Entity any
}

// Test nested struct with tags
type taggedGameState struct {
	Player taggedPlayer `json:"player_data"`
	Level  int          `json:"current_level"`
}

// TestStateExporterJSONTags tests that JSON tags don't affect path navigation.
// You must use Go field names, not JSON tag names.
func TestStateExporterJSONTags(t *testing.T) {
	state := &taggedGameState{
		Player: taggedPlayer{
			Name:   "Hero",
			Health: 100,
			X:      10.5,
		},
		Level: 5,
	}

	// Query by Go field name - should work
	result, err := navigatePath(state, "Player.Name")
	if err != nil {
		t.Fatalf("unexpected error querying Player.Name: %v", err)
	}
	if result != "Hero" {
		t.Errorf("expected Player.Name=Hero, got %v", result)
	}

	// Query by Go field name for tagged field - should work
	result, err = navigatePath(state, "Player.Health")
	if err != nil {
		t.Fatalf("unexpected error querying Player.Health: %v", err)
	}
	if result != 100 {
		t.Errorf("expected Player.Health=100, got %v", result)
	}

	// Query by JSON tag name - should NOT work (path not found)
	_, err = navigatePath(state, "player_data.Name")
	if err != ErrPathNotFound {
		t.Errorf("expected ErrPathNotFound for JSON tag path 'player_data.Name', got %v", err)
	}

	// Query by JSON tag name for nested field - should NOT work
	_, err = navigatePath(state, "Player.player_name")
	if err != ErrPathNotFound {
		t.Errorf("expected ErrPathNotFound for JSON tag path 'Player.player_name', got %v", err)
	}
}

// TestStateExporterInterfaceField tests interface field behavior.
// You can query the interface field itself, but cannot traverse into it.
func TestStateExporterInterfaceField(t *testing.T) {
	// Interface holding a struct
	state := &gameWithInterface{
		Entity: taggedPlayer{
			Name:   "Hero",
			Health: 100,
			X:      10.5,
		},
	}

	// Query interface field directly - should return the interface value
	result, err := navigatePath(state, "Entity")
	if err != nil {
		t.Fatalf("unexpected error querying Entity: %v", err)
	}

	// The result should be the taggedPlayer value stored in the interface
	_, ok := result.(taggedPlayer)
	if !ok {
		t.Errorf("expected taggedPlayer, got %T", result)
	}

	// Query nested field inside interface - not supported
	_, err = navigatePath(state, "Entity.Name")
	if err != ErrPathNotFound {
		t.Errorf("expected ErrPathNotFound for nested interface query, got %v", err)
	}
}

// TestStateExporterInterfaceWithPointer tests interface holding pointer.
func TestStateExporterInterfaceWithPointer(t *testing.T) {
	state := &gameWithInterface{
		Entity: &taggedPlayer{
			Name:   "Hero",
			Health: 100,
			X:      10.5,
		},
	}

	// Query interface field directly
	result, err := navigatePath(state, "Entity")
	if err != nil {
		t.Fatalf("unexpected error querying Entity: %v", err)
	}

	// Should be a pointer to taggedPlayer
	_, ok := result.(*taggedPlayer)
	if !ok {
		t.Errorf("expected *taggedPlayer, got %T", result)
	}

	// Query nested field inside interface containing pointer - not supported
	_, err = navigatePath(state, "Entity.Name")
	if err != ErrPathNotFound {
		t.Errorf("expected ErrPathNotFound for nested interface query, got %v", err)
	}
}

// TestStateExporterNilInterface tests nil interface field.
func TestStateExporterNilInterface(t *testing.T) {
	state := &gameWithInterface{
		Entity: nil,
	}

	// Query nil interface - returns nil with no error
	result, err := navigatePath(state, "Entity")
	if err != nil {
		t.Fatalf("unexpected error querying nil Entity: %v", err)
	}
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}
