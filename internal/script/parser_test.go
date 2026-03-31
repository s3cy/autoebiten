package script

import (
	"testing"
)

func TestParseSimple(t *testing.T) {
	input := `{
		"version": "1.0",
		"commands": [
			{"input": {"action": "press", "key": "KeyA"}},
			{"delay": {"ms": 100}}
		]
	}`

	script, err := ParseBytes([]byte(input))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	if script.Version != "1.0" {
		t.Errorf("Expected version 1.0, got %s", script.Version)
	}

	if len(script.Commands) != 2 {
		t.Errorf("Expected 2 commands, got %d", len(script.Commands))
	}
}

func TestParseRepeat(t *testing.T) {
	input := `{
		"version": "1.0",
		"commands": [
			{"repeat": {"times": 3, "commands": [
				{"input": {"action": "press", "key": "KeyA"}}
			]}}
		]
	}`

	script, err := ParseBytes([]byte(input))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	if len(script.Commands) != 1 {
		t.Fatalf("Expected 1 command, got %d", len(script.Commands))
	}

	repeatCmd, ok := script.Commands[0].(*RepeatCmd)
	if !ok {
		t.Fatal("Expected RepeatCmd")
	}

	if repeatCmd.Times != 3 {
		t.Errorf("Expected times 3, got %d", repeatCmd.Times)
	}

	if len(repeatCmd.Commands) != 1 {
		t.Errorf("Expected 1 nested command, got %d", len(repeatCmd.Commands))
	}
}

func TestParseCustom(t *testing.T) {
	input := `{
		"version": "1.0",
		"commands": [
			{"custom": {"name": "getPlayerInfo"}},
			{"custom": {"name": "heal", "request": "+20"}}
		]
	}`

	script, err := ParseBytes([]byte(input))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	if len(script.Commands) != 2 {
		t.Fatalf("Expected 2 commands, got %d", len(script.Commands))
	}

	cmd1, ok := script.Commands[0].(*CustomCmd)
	if !ok {
		t.Fatal("Expected CustomCmd for first command")
	}
	if cmd1.Name != "getPlayerInfo" {
		t.Errorf("Expected name 'getPlayerInfo', got %s", cmd1.Name)
	}
	if cmd1.Request != "" {
		t.Errorf("Expected empty request, got %s", cmd1.Request)
	}

	cmd2, ok := script.Commands[1].(*CustomCmd)
	if !ok {
		t.Fatal("Expected CustomCmd for second command")
	}
	if cmd2.Name != "heal" {
		t.Errorf("Expected name 'heal', got %s", cmd2.Name)
	}
	if cmd2.Request != "+20" {
		t.Errorf("Expected request '+20', got %s", cmd2.Request)
	}
}

func TestParseInvalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"invalid JSON", `{invalid}`},
		{"unknown command type", `{"version": "1.0", "commands": [{"unknown": {}}]}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseBytes([]byte(tt.input))
			if err == nil {
				t.Error("Expected error, got nil")
			}
		})
	}
}
