package cli

import (
	"fmt"
	"sort"

	"github.com/s3cy/autoebiten/internal/input"
	"github.com/s3cy/autoebiten/internal/rpc"
	"github.com/s3cy/autoebiten/internal/script"
)

// EnsureTargetPID auto-detects the game PID if not already set.
func EnsureTargetPID() error {
	// Only auto-detect if not already set via SetTargetPID
	game, err := rpc.AutoSelectGame()
	if err != nil {
		return err
	}
	rpc.SetTargetPID(game.PID)
	return nil
}

// CommandExecutor executes CLI commands via RPC.
type CommandExecutor struct {
	writer *Writer
}

// NewCommandExecutor creates a new CommandExecutor.
func NewCommandExecutor() *CommandExecutor {
	return &CommandExecutor{
		writer: NewWriter(),
	}
}

// RunInputCommand runs an input command.
func (e *CommandExecutor) RunInputCommand(key, action string, durationTicks int64) error {
	if action == "" {
		action = "hold"
	}
	if durationTicks == 0 {
		durationTicks = 6
	}
	params := &rpc.InputParams{
		Action:        action,
		Key:           key,
		DurationTicks: durationTicks,
	}

	req, err := rpc.BuildRequest("input", params)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := rpc.SendRequestSocket(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("rpc error: %s", resp.Error.Message)
	}

	e.writer.Success(fmt.Sprintf("input %s %s", action, key))
	return nil
}

// RunMouseCommand runs a mouse command.
func (e *CommandExecutor) RunMouseCommand(action string, x, y int, button string, durationTicks int64) error {
	if action == "" {
		if button != "" {
			action = "hold"
		} else {
			action = "position"
		}
	}
	if durationTicks == 0 {
		durationTicks = 6
	}
	params := &rpc.MouseParams{
		Action:        action,
		X:             x,
		Y:             y,
		Button:        button,
		DurationTicks: durationTicks,
	}

	req, err := rpc.BuildRequest("mouse", params)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := rpc.SendRequestSocket(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("rpc error: %s", resp.Error.Message)
	}

	e.writer.Success(fmt.Sprintf("mouse %s at (%d, %d)", action, x, y))
	return nil
}

// RunWheelCommand runs a wheel command.
func (e *CommandExecutor) RunWheelCommand(x, y float64) error {
	params := &rpc.WheelParams{
		X: x,
		Y: y,
	}

	req, err := rpc.BuildRequest("wheel", params)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := rpc.SendRequestSocket(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("rpc error: %s", resp.Error.Message)
	}

	e.writer.Success(fmt.Sprintf("wheel moved by (%.2f, %.2f)", x, y))
	return nil
}

// RunScreenshotCommand runs a screenshot command.
func (e *CommandExecutor) RunScreenshotCommand(output string, async bool) error {
	params := &rpc.ScreenshotParams{
		Output: output,
		Async:  async,
	}

	req, err := rpc.BuildRequest("screenshot", params)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := rpc.SendRequestSocket(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("rpc error: %s", resp.Error.Message)
	}

	result, ok := resp.Result.(map[string]any)
	if !ok {
		return fmt.Errorf("invalid response format")
	}

	if path, ok := result["path"].(string); ok && path != "" {
		e.writer.Success(fmt.Sprintf("screenshot saved to %s", path))
	} else if _, ok := result["data"]; ok {
		e.writer.Success(fmt.Sprintf("screenshot captured (base64)\n%s", result["data"]))
	} else {
		e.writer.Success("screenshot captured")
	}

	return nil
}

// RunPingCommand runs a ping command to check connection.
func (e *CommandExecutor) RunPingCommand() error {
	req, err := rpc.BuildRequest("ping", nil)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := rpc.SendRequestSocket(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("rpc error: %s", resp.Error.Message)
	}

	e.writer.Success("game is running")
	return nil
}

// RunGetMousePositionCommand runs a get_mouse_position command.
func (e *CommandExecutor) RunGetMousePositionCommand() error {
	req, err := rpc.BuildRequest("get_mouse_position", nil)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := rpc.SendRequestSocket(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("rpc error: %s", resp.Error.Message)
	}

	result, ok := resp.Result.(*rpc.GetMousePositionResult)
	if !ok {
		return fmt.Errorf("invalid response format")
	}

	e.writer.Success(fmt.Sprintf("mouse position: (%d, %d)", result.X, result.Y))
	return nil
}

// RunGetWheelPositionCommand runs a get_wheel_position command.
func (e *CommandExecutor) RunGetWheelPositionCommand() error {
	req, err := rpc.BuildRequest("get_wheel_position", nil)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := rpc.SendRequestSocket(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("rpc error: %s", resp.Error.Message)
	}

	result, ok := resp.Result.(*rpc.GetWheelPositionResult)
	if !ok {
		return fmt.Errorf("invalid response format")
	}

	e.writer.Success(fmt.Sprintf("wheel position: (%.2f, %.2f)", result.X, result.Y))
	return nil
}

// RunScriptCommand runs a script from a file path or inline JSON string.
func (e *CommandExecutor) RunScriptCommand(input string, isFile bool) error {
	// Parse the script
	var s *script.Script
	var err error

	if isFile {
		s, err = script.Parse(input)
	} else {
		s, err = script.ParseString(input)
	}
	if err != nil {
		return fmt.Errorf("failed to parse script: %w", err)
	}

	// Create executor
	executor := script.NewExecutor(s)

	// Set up handlers that send RPC commands
	executor.SetInputFunc(func(key, action string, durationTicks int64) error {
		return e.RunInputCommand(key, action, durationTicks)
	})

	executor.SetMouseFunc(func(action string, x, y int, button string, durationTicks int64) error {
		return e.RunMouseCommand(action, x, y, button, durationTicks)
	})

	executor.SetWheelFunc(func(x, y float64) error {
		return e.RunWheelCommand(x, y)
	})

	executor.SetScreenshotFunc(func(output string, async bool) error {
		return e.RunScreenshotCommand(output, async)
	})

	// Execute
	count, err := executor.Execute()
	if err != nil {
		return fmt.Errorf("script execution failed: %w", err)
	}

	e.writer.Success(fmt.Sprintf("executed %d commands", count))
	return nil
}

// ListKeysCommand lists all available key names.
func (e *CommandExecutor) ListKeysCommand() error {
	var keys []string
	for name := range input.StringKeyMap {
		keys = append(keys, name)
	}
	sort.Strings(keys)
	return e.writer.PrintJSON(keys)
}

// ListMouseButtonsCommand lists all available mouse button names.
func (e *CommandExecutor) ListMouseButtonsCommand() error {
	var buttons []string
	for name := range input.StringMouseButtonMap {
		buttons = append(buttons, name)
	}
	sort.Strings(buttons)
	return e.writer.PrintJSON(buttons)
}
