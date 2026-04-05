package autoebiten

import (
	"testing"
)

// Test player struct for state tests
type testPlayer struct {
	X      float64
	Y      float64
	Health int
	Name   string
}

type testInventoryItem struct {
	Name string
	Qty  int
}

type testGameState struct {
	Player    testPlayer
	Inventory []testInventoryItem
	Skills    map[string]int
	Level     int
}

func newTestGameState() *testGameState {
	return &testGameState{
		Player: testPlayer{
			X:      100.5,
			Y:      200.5,
			Health: 100,
			Name:   "Hero",
		},
		Inventory: []testInventoryItem{
			{Name: "Sword", Qty: 1},
			{Name: "Shield", Qty: 2},
		},
		Skills: map[string]int{
			"Sword":  10,
			"Shield": 5,
		},
		Level: 5,
	}
}

// TestStateExporterFieldAccess tests accessing simple struct fields.
func TestStateExporterFieldAccess(t *testing.T) {
	state := newTestGameState()

	// Test Level field
	result, err := navigatePath(state, "Level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 5 {
		t.Errorf("expected Level=5, got %v", result)
	}

	// Test Player.X
	result, err = navigatePath(state, "Player.X")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 100.5 {
		t.Errorf("expected Player.X=100.5, got %v", result)
	}

	// Test Player.Name
	result, err = navigatePath(state, "Player.Name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Hero" {
		t.Errorf("expected Player.Name=Hero, got %v", result)
	}
}

// TestStateExporterSliceAccess tests accessing slice elements.
func TestStateExporterSliceAccess(t *testing.T) {
	state := newTestGameState()

	// Test first inventory item
	result, err := navigatePath(state, "Inventory.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	item, ok := result.(testInventoryItem)
	if !ok {
		t.Fatalf("expected testInventoryItem, got %T", result)
	}
	if item.Name != "Sword" {
		t.Errorf("expected item.Name=Sword, got %v", item.Name)
	}

	// Test second inventory item
	result, err = navigatePath(state, "Inventory.1.Name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Shield" {
		t.Errorf("expected Inventory.1.Name=Shield, got %v", result)
	}
}

// TestStateExporterMapAccess tests accessing map elements.
func TestStateExporterMapAccess(t *testing.T) {
	state := newTestGameState()

	// Test skill access
	result, err := navigatePath(state, "Skills.Sword")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 10 {
		t.Errorf("expected Skills.Sword=10, got %v", result)
	}

	// Test another skill
	result, err = navigatePath(state, "Skills.Shield")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 5 {
		t.Errorf("expected Skills.Shield=5, got %v", result)
	}
}

// TestStateExporterPathNotFound tests error handling for invalid paths.
func TestStateExporterPathNotFound(t *testing.T) {
	state := newTestGameState()

	// Test non-existent field
	_, err := navigatePath(state, "NonExistent")
	if err != ErrPathNotFound {
		t.Errorf("expected ErrPathNotFound, got %v", err)
	}

	// Test non-existent nested field
	_, err = navigatePath(state, "Player.NonExistent")
	if err != ErrPathNotFound {
		t.Errorf("expected ErrPathNotFound, got %v", err)
	}

	// Test out of bounds slice index
	_, err = navigatePath(state, "Inventory.100")
	if err != ErrPathNotFound {
		t.Errorf("expected ErrPathNotFound, got %v", err)
	}

	// Test non-existent map key
	_, err = navigatePath(state, "Skills.NonExistent")
	if err != ErrPathNotFound {
		t.Errorf("expected ErrPathNotFound, got %v", err)
	}
}

// TestStateExporterNestedAccess tests deeply nested paths.
func TestStateExporterNestedAccess(t *testing.T) {
	state := newTestGameState()

	// Test Inventory.0.Name
	result, err := navigatePath(state, "Inventory.0.Name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Sword" {
		t.Errorf("expected Inventory.0.Name=Sword, got %v", result)
	}

	// Test Inventory.1.Qty
	result, err = navigatePath(state, "Inventory.1.Qty")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 2 {
		t.Errorf("expected Inventory.1.Qty=2, got %v", result)
	}
}

// TestParsePath tests the path parsing function.
func TestParsePath(t *testing.T) {
	tests := []struct {
		path     string
		expected []string
	}{
		{"Player.X", []string{"Player", "X"}},
		{"Inventory.0.Name", []string{"Inventory", "0", "Name"}},
		{"Skills.Sword", []string{"Skills", "Sword"}},
		{"Level", []string{"Level"}},
	}

	for _, tc := range tests {
		parts := parsePath(tc.path)
		if len(parts) != len(tc.expected) {
			t.Errorf("parsePath(%s): expected %v, got %v", tc.path, tc.expected, parts)
			continue
		}
		for i := 0; i < len(parts); i++ {
			if parts[i] != tc.expected[i] {
				t.Errorf("parsePath(%s)[%d]: expected %s, got %s", tc.path, i, tc.expected[i], parts[i])
			}
		}
	}
}

// TestStateQueryPath tests the reusable path parsing.
func TestStateQueryPath(t *testing.T) {
	state := newTestGameState()

	// Parse path once
	path, err := ParseStateQueryPath("Player.X")
	if err != nil {
		t.Fatalf("failed to parse path: %v", err)
	}

	// Query multiple times
	for i := 0; i < 3; i++ {
		result, err := path.Query(state)
		if err != nil {
			t.Fatalf("query %d failed: %v", i, err)
		}
		if result != 100.5 {
			t.Errorf("query %d: expected 100.5, got %v", i, result)
		}
	}
}
