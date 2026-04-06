package recording

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sync"
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

func TestRecord(t *testing.T) {
	// Use a temp directory for test isolation
	tmpDir := t.TempDir()

	// Override RecordingDir for this test
	oldDir := RecordingDir

	tests := []struct {
		name    string
		pid     int
		cmd     script.CommandWrapper
	}{
		{"input command", 12345, &script.InputCmd{Action: "press", Key: "KeyA"}},
		{"mouse command", 12346, &script.MouseCmd{Action: "position", X: 100, Y: 200}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set test dir
			RecordingDir = tmpDir
			defer func() { RecordingDir = oldDir }()

			recorder := NewRecorder(tt.pid)
			if err := recorder.Record(tt.cmd); err != nil {
				t.Fatalf("Record failed: %v", err)
			}

			// Verify file exists and has content
			path := Path(tt.pid)
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read file: %v", err)
			}

			if len(data) == 0 {
				t.Error("Recording file is empty")
			}

			// Verify valid JSON line
			lines := bytes.Split(data, []byte("\n"))
			if len(lines) < 1 || len(lines[0]) == 0 {
				t.Error("No JSON line in recording file")
			}
		})
	}
}

func TestRecordConcurrent(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir := RecordingDir
	RecordingDir = tmpDir
	defer func() { RecordingDir = oldDir }()

	pid := 54321
	recorder := NewRecorder(pid)

	// Write 10 entries concurrently
	var wg sync.WaitGroup
	errCh := make(chan error, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cmd := &script.InputCmd{Action: "press", Key: fmt.Sprintf("Key%d", i)}
			if err := recorder.Record(cmd); err != nil {
				errCh <- err
			}
		}(i)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			t.Fatalf("Concurrent record error: %v", err)
		}
	}

	// Verify 10 lines written
	path := Path(pid)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	lines := bytes.Split(bytes.TrimSuffix(data, []byte("\n")), []byte("\n"))
	if len(lines) != 10 {
		t.Errorf("Expected 10 lines, got %d", len(lines))
	}

	// Verify each line is valid JSON
	for i, line := range lines {
		var entry Entry
		if err := json.Unmarshal(line, &entry); err != nil {
			t.Errorf("Line %d is invalid JSON: %v", i, err)
		}
	}
}

func TestReadAll(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir := RecordingDir
	RecordingDir = tmpDir
	defer func() { RecordingDir = oldDir }()

	pid := 99999

	// Write some entries first
	recorder := NewRecorder(pid)
	recorder.Record(&script.InputCmd{Action: "press", Key: "KeyA"})
	recorder.Record(&script.MouseCmd{Action: "position", X: 100, Y: 200})
	recorder.Record(&script.ScreenshotCmd{Output: "test.png"})

	// Read them back
	reader := NewReader(pid)
	entries, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll failed: %v", err)
	}

	if len(entries) != 3 {
		t.Fatalf("Expected 3 entries, got %d", len(entries))
	}

	// Verify first entry is InputCmd
	inputCmd, ok := entries[0].Command.(*script.InputCmd)
	if !ok {
		t.Fatal("First entry should be InputCmd")
	}
	if inputCmd.Key != "KeyA" {
		t.Errorf("Expected KeyA, got %s", inputCmd.Key)
	}

	// Verify second entry is MouseCmd
	mouseCmd, ok := entries[1].Command.(*script.MouseCmd)
	if !ok {
		t.Fatal("Second entry should be MouseCmd")
	}
	if mouseCmd.X != 100 || mouseCmd.Y != 200 {
		t.Errorf("Expected (100, 200), got (%d, %d)", mouseCmd.X, mouseCmd.Y)
	}
}

func TestReadAllEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir := RecordingDir
	RecordingDir = tmpDir
	defer func() { RecordingDir = oldDir }()

	pid := 77777
	reader := NewReader(pid)
	entries, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll should not fail on missing file: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("Expected 0 entries for missing file, got %d", len(entries))
	}
}

func TestGenerate(t *testing.T) {
	// Create entries with known timestamps
	baseTime := time.Now()
	entries := []Entry{
		{Timestamp: baseTime, Command: &script.InputCmd{Action: "press", Key: "KeyA"}},
		{Timestamp: baseTime.Add(500 * time.Millisecond), Command: &script.MouseCmd{Action: "position", X: 100, Y: 200}},
		{Timestamp: baseTime.Add(800 * time.Millisecond), Command: &script.ScreenshotCmd{Output: "shot.png"}},
	}

	gen := NewGenerator(1.0)
	s, err := gen.Generate(entries)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if s.Version != "1.0" {
		t.Errorf("Expected version 1.0, got %s", s.Version)
	}

	// Expected: input, delay(500ms), mouse, delay(300ms), screenshot
	if len(s.Commands) != 5 {
		t.Fatalf("Expected 5 commands (input + delay + mouse + delay + screenshot), got %d", len(s.Commands))
	}

	// Verify first command is input
	_, ok := s.Commands[0].(*script.InputCmd)
	if !ok {
		t.Fatal("First command should be InputCmd")
	}

	// Verify second command is delay 500ms
	delay1, ok := s.Commands[1].(*script.DelayCmd)
	if !ok {
		t.Fatal("Second command should be DelayCmd")
	}
	if delay1.Ms != 500 {
		t.Errorf("Expected delay 500ms, got %d", delay1.Ms)
	}

	// Verify third command is mouse
	_, ok = s.Commands[2].(*script.MouseCmd)
	if !ok {
		t.Fatal("Third command should be MouseCmd")
	}

	// Verify fourth command is delay 300ms
	delay2, ok := s.Commands[3].(*script.DelayCmd)
	if !ok {
		t.Fatal("Fourth command should be DelayCmd")
	}
	if delay2.Ms != 300 {
		t.Errorf("Expected delay 300ms, got %d", delay2.Ms)
	}

	// Verify fifth command is screenshot
	_, ok = s.Commands[4].(*script.ScreenshotCmd)
	if !ok {
		t.Fatal("Fifth command should be ScreenshotCmd")
	}
}

func TestGenerateWithSpeed(t *testing.T) {
	baseTime := time.Now()
	entries := []Entry{
		{Timestamp: baseTime, Command: &script.InputCmd{Action: "press", Key: "KeyA"}},
		{Timestamp: baseTime.Add(500 * time.Millisecond), Command: &script.MouseCmd{Action: "position", X: 100, Y: 200}},
	}

	gen := NewGenerator(2.0) // 2x speed
	s, err := gen.Generate(entries)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Delay should be 500/2 = 250ms
	delay, ok := s.Commands[1].(*script.DelayCmd)
	if !ok {
		t.Fatal("Second command should be DelayCmd")
	}
	if delay.Ms != 250 {
		t.Errorf("Expected delay 250ms (500/2), got %d", delay.Ms)
	}
}

func TestGenerateEmpty(t *testing.T) {
	gen := NewGenerator(1.0)
	_, err := gen.Generate([]Entry{})
	if err == nil {
		t.Fatal("Expected error for empty entries")
	}
}

func TestGenerateZeroSpeed(t *testing.T) {
	gen := NewGenerator(0.0)
	entries := []Entry{
		{Timestamp: time.Now(), Command: &script.InputCmd{Action: "press", Key: "KeyA"}},
	}
	_, err := gen.Generate(entries)
	if err == nil {
		t.Fatal("Expected error for speed 0")
	}
}

func TestRecordingFlow(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir := RecordingDir
	RecordingDir = tmpDir
	defer func() { RecordingDir = oldDir }()

	pid := 11111

	// Record several commands
	recorder := NewRecorder(pid)
	recorder.Record(&script.InputCmd{Action: "press", Key: "KeyA"})
	recorder.Record(&script.MouseCmd{Action: "position", X: 100, Y: 200})
	recorder.Record(&script.DelayCmd{Ms: 500})
	recorder.Record(&script.ScreenshotCmd{Output: "shot.png"})

	// Read back
	reader := NewReader(pid)
	entries, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll failed: %v", err)
	}
	if len(entries) != 4 {
		t.Fatalf("Expected 4 entries, got %d", len(entries))
	}

	// Generate script
	gen := NewGenerator(1.0)
	script, err := gen.Generate(entries)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Verify script structure
	if script.Version != "1.0" {
		t.Errorf("Expected version 1.0, got %s", script.Version)
	}

	// Should have delays inserted between commands
	if len(script.Commands) < 4 {
		t.Errorf("Expected at least 4 commands, got %d", len(script.Commands))
	}

	// Clear recording
	if err := Clear(pid); err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	// Verify file removed
	path := Path(pid)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("Recording file should be removed after clear")
	}
}
