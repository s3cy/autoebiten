package script

import (
	"testing"
)

func TestExecutorBasic(t *testing.T) {
	script := &Script{
		Version: "1.0",
		Commands: []CommandWrapper{
			&DelayCmd{Ms: 1},
		},
	}

	executor := NewExecutor(script)
	count, err := executor.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}
}

func TestExecutorCustomCommand(t *testing.T) {
	script := &Script{
		Version: "1.0",
		Commands: []CommandWrapper{
			&CustomCmd{Name: "testCmd", Request: "test data"},
		},
	}

	var capturedName, capturedRequest string

	executor := NewExecutor(script)
	executor.SetCustomFunc(func(name, request string) error {
		capturedName = name
		capturedRequest = request
		return nil
	})

	count, err := executor.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}

	if capturedName != "testCmd" {
		t.Errorf("Expected name 'testCmd', got %s", capturedName)
	}
	if capturedRequest != "test data" {
		t.Errorf("Expected request 'test data', got %s", capturedRequest)
	}
}

func TestExecutorCustomCommandError(t *testing.T) {
	script := &Script{
		Version: "1.0",
		Commands: []CommandWrapper{
			&CustomCmd{Name: "failingCmd", Request: ""},
		},
	}

	executor := NewExecutor(script)
	executor.SetCustomFunc(func(name, request string) error {
		return errCustom
	})

	_, err := executor.Execute()
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestExecutorCustomFuncNotSet(t *testing.T) {
	script := &Script{
		Version: "1.0",
		Commands: []CommandWrapper{
			&CustomCmd{Name: "testCmd"},
		},
	}

	executor := NewExecutor(script)
	_, err := executor.Execute()
	if err == nil {
		t.Fatal("Expected error when custom func not set, got nil")
	}
}

func TestExecutorRepeatWithCustom(t *testing.T) {
	script := &Script{
		Version: "1.0",
		Commands: []CommandWrapper{
			&RepeatCmd{
				Times: 3,
				Commands: []CommandWrapper{
					&CustomCmd{Name: "repeatCmd"},
				},
			},
		},
	}

	callCount := 0
	executor := NewExecutor(script)
	executor.SetCustomFunc(func(name, request string) error {
		callCount++
		return nil
	})

	count, err := executor.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if count != 4 { // 1 repeat + 3 custom
		t.Errorf("Expected count 4, got %d", count)
	}
	if callCount != 3 {
		t.Errorf("Expected custom func called 3 times, got %d", callCount)
	}
}

func TestExecutorMixedCommands(t *testing.T) {
	script := &Script{
		Version: "1.0",
		Commands: []CommandWrapper{
			&InputCmd{Key: "KeyA", Action: "press"},
			&CustomCmd{Name: "testCmd"},
			&DelayCmd{Ms: 1},
		},
	}

	inputCalled := false
	customCalled := false

	executor := NewExecutor(script)
	executor.SetInputFunc(func(key, action string, durationTicks int64, async bool) error {
		inputCalled = true
		return nil
	})
	executor.SetCustomFunc(func(name, request string) error {
		customCalled = true
		return nil
	})

	_, err := executor.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !inputCalled {
		t.Error("Input func was not called")
	}
	if !customCalled {
		t.Error("Custom func was not called")
	}
}

func TestExecutorInputFuncNotSet(t *testing.T) {
	script := &Script{
		Version: "1.0",
		Commands: []CommandWrapper{
			&InputCmd{Key: "KeyA", Action: "press"},
		},
	}

	executor := NewExecutor(script)
	_, err := executor.Execute()
	if err == nil {
		t.Fatal("Expected error when input func not set, got nil")
	}
}

func TestExecutorMouseFuncNotSet(t *testing.T) {
	script := &Script{
		Version: "1.0",
		Commands: []CommandWrapper{
			&MouseCmd{Action: "click", X: 10, Y: 20},
		},
	}

	executor := NewExecutor(script)
	_, err := executor.Execute()
	if err == nil {
		t.Fatal("Expected error when mouse func not set, got nil")
	}
}

func TestExecutorWheelFuncNotSet(t *testing.T) {
	script := &Script{
		Version: "1.0",
		Commands: []CommandWrapper{
			&WheelCmd{X: 0, Y: 10},
		},
	}

	executor := NewExecutor(script)
	_, err := executor.Execute()
	if err == nil {
		t.Fatal("Expected error when wheel func not set, got nil")
	}
}

func TestExecutorScreenshotFuncNotSet(t *testing.T) {
	script := &Script{
		Version: "1.0",
		Commands: []CommandWrapper{
			&ScreenshotCmd{Output: "test.png"},
		},
	}

	executor := NewExecutor(script)
	_, err := executor.Execute()
	if err == nil {
		t.Fatal("Expected error when screenshot func not set, got nil")
	}
}

func TestExecutorStateCommand(t *testing.T) {
	script := &Script{
		Version: "1.0",
		Commands: []CommandWrapper{
			&StateCmd{Name: "gamestate", Path: "Player.X"},
		},
	}

	var capturedName, capturedPath string

	executor := NewExecutor(script)
	executor.SetStateFunc(func(name, path string) error {
		capturedName = name
		capturedPath = path
		return nil
	})

	count, err := executor.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}

	if capturedName != "gamestate" {
		t.Errorf("Expected name 'gamestate', got %s", capturedName)
	}
	if capturedPath != "Player.X" {
		t.Errorf("Expected path 'Player.X', got %s", capturedPath)
	}
}

func TestExecutorStateFuncNotSet(t *testing.T) {
	script := &Script{
		Version: "1.0",
		Commands: []CommandWrapper{
			&StateCmd{Name: "gamestate", Path: "Player.X"},
		},
	}

	executor := NewExecutor(script)
	_, err := executor.Execute()
	if err == nil {
		t.Fatal("Expected error when state func not set, got nil")
	}
}

func TestExecutorWaitCommand(t *testing.T) {
	script := &Script{
		Version: "1.0",
		Commands: []CommandWrapper{
			&WaitCmd{Condition: "state:gamestate:Player.X == 100", Timeout: "5s", Interval: "100ms"},
		},
	}

	var capturedCondition, capturedTimeout, capturedInterval string

	executor := NewExecutor(script)
	executor.SetWaitFunc(func(condition, timeout, interval string, verbose bool) error {
		capturedCondition = condition
		capturedTimeout = timeout
		capturedInterval = interval
		return nil
	})

	count, err := executor.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}

	if capturedCondition != "state:gamestate:Player.X == 100" {
		t.Errorf("Expected condition 'state:gamestate:Player.X == 100', got %s", capturedCondition)
	}
	if capturedTimeout != "5s" {
		t.Errorf("Expected timeout '5s', got %s", capturedTimeout)
	}
	if capturedInterval != "100ms" {
		t.Errorf("Expected interval '100ms', got %s", capturedInterval)
	}
}

func TestExecutorWaitFuncNotSet(t *testing.T) {
	script := &Script{
		Version: "1.0",
		Commands: []CommandWrapper{
			&WaitCmd{Condition: "state:gamestate:Player.X == 100", Timeout: "5s"},
		},
	}

	executor := NewExecutor(script)
	_, err := executor.Execute()
	if err == nil {
		t.Fatal("Expected error when wait func not set, got nil")
	}
}

var errCustom = &testError{msg: "custom error"}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
