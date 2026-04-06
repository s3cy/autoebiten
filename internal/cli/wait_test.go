package cli

import (
	"testing"
)

func TestParseCondition(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      *Condition
		wantError bool
	}{
		{
			name:  "state query with number equality",
			input: "state:gamestate:Player.Health == 100",
			want: &Condition{
				Type:     "state",
				Name:     "gamestate",
				Path:     "Player.Health",
				Operator: "==",
				Value:    float64(100),
			},
		},
		{
			name:  "state query with greater than",
			input: "state:gamestate:Player.X > 50",
			want: &Condition{
				Type:     "state",
				Name:     "gamestate",
				Path:     "Player.X",
				Operator: ">",
				Value:    float64(50),
			},
		},
		{
			name:  "custom command with string inequality",
			input: `custom:getStatus:ready != "error"`,
			want: &Condition{
				Type:     "custom",
				Name:     "getStatus",
				Path:     "ready",
				Operator: "!=",
				Value:    "error",
			},
		},
		{
			name:  "state query with boolean",
			input: "state:gamestate:Player.Alive == true",
			want: &Condition{
				Type:     "state",
				Name:     "gamestate",
				Path:     "Player.Alive",
				Operator: "==",
				Value:    true,
			},
		},
		{
			name:  "state query with less than or equal",
			input: "state:gamestate:Enemies.Count <= 10",
			want: &Condition{
				Type:     "state",
				Name:     "gamestate",
				Path:     "Enemies.Count",
				Operator: "<=",
				Value:    float64(10),
			},
		},
		{
			name:      "missing operator",
			input:     "state:gamestate:Player.Health",
			wantError: true,
		},
		{
			name:      "invalid format - no colons",
			input:     "stategamestatePlayer.Health == 100",
			wantError: true,
		},
		{
			name:      "invalid operator",
			input:     "state:gamestate:Player.Health ~= 100",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCondition(tt.input)
			if tt.wantError {
				if err == nil {
					t.Errorf("ParseCondition(%q) expected error, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Errorf("ParseCondition(%q) unexpected error: %v", tt.input, err)
				return
			}
			if got.Type != tt.want.Type {
				t.Errorf("Type: got %q, want %q", got.Type, tt.want.Type)
			}
			if got.Name != tt.want.Name {
				t.Errorf("Name: got %q, want %q", got.Name, tt.want.Name)
			}
			if got.Path != tt.want.Path {
				t.Errorf("Path: got %q, want %q", got.Path, tt.want.Path)
			}
			if got.Operator != tt.want.Operator {
				t.Errorf("Operator: got %q, want %q", got.Operator, tt.want.Operator)
			}
		})
	}
}