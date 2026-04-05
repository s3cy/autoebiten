package testkit

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/s3cy/autoebiten/internal/input"
)

// getStatefulTestGameBinary returns the path to the stateful test game binary.
func getStatefulTestGameBinary() string {
	candidates := []string{
		"./internal/testgames/stateful/stateful",
		"./internal/testgames/stateful/main",
	}
	if os.PathSeparator == '\\' {
		for i := range candidates {
			candidates[i] += ".exe"
		}
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return ""
}

// getCustomTestGameBinary returns the path to the custom test game binary.
func getCustomTestGameBinary() string {
	candidates := []string{
		"./internal/testgames/custom/custom",
		"./internal/testgames/custom/main",
	}
	if os.PathSeparator == '\\' {
		for i := range candidates {
			candidates[i] += ".exe"
		}
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return ""
}

// TestE2EPlayerMovement tests player movement in a stateful game.
func TestE2EPlayerMovement(t *testing.T) {
	binaryPath := getStatefulTestGameBinary()
	if binaryPath == "" {
		t.Skip("stateful test game binary not available - build with: go build -o testkit/internal/testgames/stateful/stateful ./testkit/internal/testgames/stateful")
	}

	game := Launch(t, binaryPath, WithTimeout(30*time.Second))
	defer game.Shutdown()

	// Wait for game to be ready
	ready := game.WaitFor(func() bool {
		return game.Ping() == nil
	}, 5*time.Second)
	if !ready {
		t.Fatal("game did not become ready in time")
	}

	// Get initial player position
	initialX, err := game.StateQuery("Player.X")
	if err != nil {
		t.Fatalf("failed to query Player.X: %v", err)
	}

	// Press right arrow key for 10 ticks
	if err := game.HoldKey(input.KeyArrowRight, 10); err != nil {
		t.Fatalf("failed to hold key: %v", err)
	}

	// Wait for movement to take effect
	time.Sleep(100 * time.Millisecond)

	// Get new player position
	newX, err := game.StateQuery("Player.X")
	if err != nil {
		t.Fatalf("failed to query Player.X after movement: %v", err)
	}

	// Verify player moved
	initialXFloat := toFloat64(initialX)
	newXFloat := toFloat64(newX)

	if newXFloat <= initialXFloat {
		t.Errorf("expected Player.X to increase, got %v -> %v", initialXFloat, newXFloat)
	}
}

// TestE2EScreenshot tests taking screenshots.
func TestE2EScreenshot(t *testing.T) {
	binaryPath := getSimpleTestGameBinary()
	if binaryPath == "" {
		t.Skip("simple test game binary not available - build with: go build -o testkit/internal/testgames/simple/simple ./testkit/internal/testgames/simple")
	}

	game := Launch(t, binaryPath, WithTimeout(30*time.Second))
	defer game.Shutdown()

	// Wait for game to be ready
	ready := game.WaitFor(func() bool {
		return game.Ping() == nil
	}, 5*time.Second)
	if !ready {
		t.Fatal("game did not become ready in time")
	}

	// Take screenshot
	img, err := game.Screenshot()
	if err != nil {
		t.Fatalf("failed to take screenshot: %v", err)
	}

	// Verify image dimensions
	if img.Bounds().Empty() {
		t.Error("screenshot is empty")
	}

	// Test ScreenshotBase64
	b64, err := game.ScreenshotBase64()
	if err != nil {
		t.Fatalf("failed to take base64 screenshot: %v", err)
	}
	if len(b64) == 0 {
		t.Error("base64 screenshot is empty")
	}
}

// TestE2ECustomCommands tests custom command execution.
func TestE2ECustomCommands(t *testing.T) {
	binaryPath := getCustomTestGameBinary()
	if binaryPath == "" {
		t.Skip("custom test game binary not available - build with: go build -o testkit/internal/testgames/custom/custom ./testkit/internal/testgames/custom")
	}

	game := Launch(t, binaryPath, WithTimeout(30*time.Second))
	defer game.Shutdown()

	// Wait for game to be ready
	ready := game.WaitFor(func() bool {
		return game.Ping() == nil
	}, 5*time.Second)
	if !ready {
		t.Fatal("game did not become ready in time")
	}

	// Test echo command
	response, err := game.RunCustom("echo", "hello world")
	if err != nil {
		t.Fatalf("failed to run echo command: %v", err)
	}
	if response != "hello world" {
		t.Errorf("echo returned %q, expected %q", response, "hello world")
	}

	// Test counter command
	for i := 1; i <= 3; i++ {
		response, err = game.RunCustom("counter", "")
		if err != nil {
			t.Fatalf("failed to run counter command: %v", err)
		}
		expected := string(rune('0' + i))
		if response != expected {
			t.Errorf("counter returned %q, expected %q", response, expected)
		}
	}
}

// TestE2EStateQuery tests state query functionality.
func TestE2EStateQuery(t *testing.T) {
	binaryPath := getStatefulTestGameBinary()
	if binaryPath == "" {
		t.Skip("stateful test game binary not available")
	}

	game := Launch(t, binaryPath, WithTimeout(30*time.Second))
	defer game.Shutdown()

	// Wait for game to be ready
	ready := game.WaitFor(func() bool {
		return game.Ping() == nil
	}, 5*time.Second)
	if !ready {
		t.Fatal("game did not become ready in time")
	}

	// Test simple field query
	health, err := game.StateQuery("Player.Health")
	if err != nil {
		t.Fatalf("failed to query Player.Health: %v", err)
	}
	if toFloat64(health) != 100 {
		t.Errorf("expected Player.Health=100, got %v", health)
	}

	// Test slice query
	itemName, err := game.StateQuery("Inventory.0.Name")
	if err != nil {
		t.Fatalf("failed to query Inventory.0.Name: %v", err)
	}
	if itemName != "Sword" {
		t.Errorf("expected Inventory.0.Name=Sword, got %v", itemName)
	}

	// Test map query
	skillLevel, err := game.StateQuery("Skills.Sword")
	if err != nil {
		t.Fatalf("failed to query Skills.Sword: %v", err)
	}
	if toFloat64(skillLevel) != 10 {
		t.Errorf("expected Skills.Sword=10, got %v", skillLevel)
	}
}

// toFloat64 converts a numeric value to float64.
func toFloat64(v any) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case int32:
		return float64(val)
	case string:
		// Try to parse as float
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}
	return 0
}
