package recording

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/s3cy/autoebiten/internal/script"
)

func TestPath(t *testing.T) {
	tests := []struct {
		name     string
		pid      int
		expected string
	}{
		{"pid 1234", 1234, "/tmp/autoebiten/recording-1234.jsonl"},
		{"pid 1", 1, "/tmp/autoebiten/recording-1.jsonl"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Path(tt.pid); got != tt.expected {
				t.Errorf("Path(%d) = %s, want %s", tt.pid, got, tt.expected)
			}
		})
	}
}

func TestEntry_MarshalJSON_InputCmd(t *testing.T) {
	entry := Entry{
		Timestamp: time.Date(2025, 4, 6, 14, 30, 22, 123456789, time.UTC),
		Command: &script.InputCmd{
			Action: "press",
			Key:    "KeyA",
		},
	}

	data, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to parse output: %v", err)
	}

	// Verify timestamp
	if ts, ok := result["timestamp"].(string); !ok || ts != "2025-04-06T14:30:22.123456789Z" {
		t.Errorf("timestamp = %v, want '2025-04-06T14:30:22.123456789Z'", result["timestamp"])
	}

	// Verify command is wrapped in "input" discriminator
	cmd, ok := result["command"].(map[string]interface{})
	if !ok {
		t.Fatal("command is not an object")
	}
	if _, ok := cmd["input"]; !ok {
		t.Error("command should have 'input' key")
	}
	if _, ok := cmd["mouse"]; ok {
		t.Error("command should not have 'mouse' key")
	}

	// Verify the input content
	input, _ := cmd["input"].(map[string]interface{})
	if input["action"] != "press" || input["key"] != "KeyA" {
		t.Errorf("input = %v, want {action: press, key: KeyA}", input)
	}
}

func TestEntry_MarshalJSON_MouseCmd(t *testing.T) {
	entry := Entry{
		Timestamp: time.Date(2025, 4, 6, 14, 30, 22, 0, time.UTC),
		Command: &script.MouseCmd{
			Action: "press",
			X:      100,
			Y:      200,
			Button: "ButtonLeft",
		},
	}

	data, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to parse output: %v", err)
	}

	// Verify command is wrapped in "mouse" discriminator
	cmd := result["command"].(map[string]interface{})
	if _, ok := cmd["mouse"]; !ok {
		t.Error("command should have 'mouse' key")
	}

	mouse, _ := cmd["mouse"].(map[string]interface{})
	if mouse["action"] != "press" || mouse["x"].(float64) != 100 || mouse["y"].(float64) != 200 {
		t.Errorf("mouse = %v", mouse)
	}
}

func TestEntry_MarshalJSON_WheelCmd(t *testing.T) {
	entry := Entry{
		Timestamp: time.Now().UTC().Truncate(time.Nanosecond),
		Command: &script.WheelCmd{
			X: 0.5,
			Y: -1.0,
		},
	}

	data, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to parse output: %v", err)
	}

	// Verify command is wrapped in "wheel" discriminator
	cmd := result["command"].(map[string]interface{})
	if _, ok := cmd["wheel"]; !ok {
		t.Error("command should have 'wheel' key")
	}
}

func TestEntry_UnmarshalJSON(t *testing.T) {
	// Test unmarshaling the expected format
	jsonData := `{"timestamp":"2025-04-06T14:30:22.123456789Z","command":{"input":{"action":"press","key":"KeyA"}}}`

	var entry Entry
	err := json.Unmarshal([]byte(jsonData), &entry)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	// Verify timestamp
	if !entry.Timestamp.Equal(time.Date(2025, 4, 6, 14, 30, 22, 123456789, time.UTC)) {
		t.Errorf("timestamp = %v, want 2025-04-06T14:30:22.123456789Z", entry.Timestamp)
	}

	// Verify command type
	inputCmd, ok := entry.Command.(*script.InputCmd)
	if !ok {
		t.Fatalf("command should be *InputCmd, got %T", entry.Command)
	}
	if inputCmd.Action != "press" || inputCmd.Key != "KeyA" {
		t.Errorf("inputCmd = %+v", inputCmd)
	}
}

func TestEntry_UnmarshalJSON_Mouse(t *testing.T) {
	jsonData := `{"timestamp":"2025-04-06T14:30:22Z","command":{"mouse":{"action":"press","x":100,"y":200,"button":"ButtonLeft"}}}`

	var entry Entry
	err := json.Unmarshal([]byte(jsonData), &entry)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	mouseCmd, ok := entry.Command.(*script.MouseCmd)
	if !ok {
		t.Fatalf("command should be *MouseCmd, got %T", entry.Command)
	}
	if mouseCmd.Action != "press" || mouseCmd.X != 100 || mouseCmd.Y != 200 {
		t.Errorf("mouseCmd = %+v", mouseCmd)
	}
}

func TestEntry_RoundTrip(t *testing.T) {
	// Test round-trip marshal/unmarshal preserves data
	original := Entry{
		Timestamp: time.Date(2025, 4, 6, 14, 30, 22, 123456789, time.UTC),
		Command: &script.InputCmd{
			Action:        "hold",
			Key:           "KeySpace",
			DurationTicks: 10,
			Async:         true,
		},
	}

	// Marshal
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	// Unmarshal
	var result Entry
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	// Verify
	if !result.Timestamp.Equal(original.Timestamp) {
		t.Errorf("timestamp = %v, want %v", result.Timestamp, original.Timestamp)
	}

	inputCmd, ok := result.Command.(*script.InputCmd)
	if !ok {
		t.Fatalf("command should be *InputCmd, got %T", result.Command)
	}
	if inputCmd.Action != original.Command.(*script.InputCmd).Action ||
		inputCmd.Key != original.Command.(*script.InputCmd).Key ||
		inputCmd.DurationTicks != original.Command.(*script.InputCmd).DurationTicks ||
		inputCmd.Async != original.Command.(*script.InputCmd).Async {
		t.Errorf("command = %+v, want %+v", inputCmd, original.Command)
	}
}

func TestEntry_RoundTrip_AllCommandTypes(t *testing.T) {
	commands := []script.CommandWrapper{
		&script.InputCmd{Action: "press", Key: "KeyA"},
		&script.MouseCmd{Action: "press", X: 100, Y: 200},
		&script.WheelCmd{X: 1.0, Y: -1.0},
		&script.ScreenshotCmd{Output: "test.png"},
		&script.DelayCmd{Ms: 100},
		&script.CustomCmd{Name: "test"},
	}

	for _, cmd := range commands {
		entry := Entry{
			Timestamp: time.Now().UTC().Truncate(time.Nanosecond),
			Command:   cmd,
		}

		// Marshal
		data, err := json.Marshal(entry)
		if err != nil {
			t.Fatalf("MarshalJSON failed for %T: %v", cmd, err)
		}

		// Unmarshal
		var result Entry
		if err := json.Unmarshal(data, &result); err != nil {
			t.Fatalf("UnmarshalJSON failed for %T: %v", cmd, err)
		}

		// Verify command type is preserved
		if fmt.Sprintf("%T", result.Command) != fmt.Sprintf("%T", cmd) {
			t.Errorf("command type = %T, want %T", result.Command, cmd)
		}
	}
}
