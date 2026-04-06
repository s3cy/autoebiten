package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/s3cy/autoebiten/internal/input"
	"github.com/s3cy/autoebiten/internal/recording"
	"github.com/s3cy/autoebiten/internal/rpc"
	"github.com/s3cy/autoebiten/internal/script"
	"github.com/s3cy/autoebiten/internal/version"
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
func (e *CommandExecutor) RunInputCommand(key, action string, durationTicks int64, async bool, shouldRecord bool) error {
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
		Async:         async,
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

	// Record after successful execution
	if shouldRecord {
		recorder := recording.NewRecorderFromSocket(rpc.SocketPath())
		cmd := &script.InputCmd{
			Action:        action,
			Key:           key,
			DurationTicks: durationTicks,
			Async:         async,
		}
		if err := recorder.Record(cmd); err != nil {
			fmt.Fprintf(os.Stderr, "autoebiten: recording failed: %v\n", err)
		}
	}

	e.writer.Success(fmt.Sprintf("input %s %s", action, key))
	return nil
}

// RunMouseCommand runs a mouse command.
func (e *CommandExecutor) RunMouseCommand(action string, x, y int, button string, durationTicks int64, async bool, shouldRecord bool) error {
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
		Async:         async,
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

	var result rpc.MouseResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return fmt.Errorf("invalid response format: %w", err)
	}

	// Record after successful execution
	if shouldRecord {
		recorder := recording.NewRecorderFromSocket(rpc.SocketPath())
		cmd := &script.MouseCmd{
			Action:        action,
			X:             x,
			Y:             y,
			Button:        button,
			DurationTicks: durationTicks,
			Async:         async,
		}
		if err := recorder.Record(cmd); err != nil {
			fmt.Fprintf(os.Stderr, "autoebiten: recording failed: %v\n", err)
		}
	}

	e.writer.Success(fmt.Sprintf("mouse %s at (%d, %d)", action, result.X, result.Y))
	return nil
}

// RunWheelCommand runs a wheel command.
func (e *CommandExecutor) RunWheelCommand(x, y float64, async bool, shouldRecord bool) error {
	params := &rpc.WheelParams{
		X:     x,
		Y:     y,
		Async: async,
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

	// Record after successful execution
	if shouldRecord {
		recorder := recording.NewRecorderFromSocket(rpc.SocketPath())
		cmd := &script.WheelCmd{
			X:     x,
			Y:     y,
			Async: async,
		}
		if err := recorder.Record(cmd); err != nil {
			fmt.Fprintf(os.Stderr, "autoebiten: recording failed: %v\n", err)
		}
	}

	e.writer.Success(fmt.Sprintf("wheel moved by (%.2f, %.2f)", x, y))
	return nil
}

// RunScreenshotCommand runs a screenshot command.
func (e *CommandExecutor) RunScreenshotCommand(output string, b64 bool, async bool, shouldRecord bool) error {
	params := &rpc.ScreenshotParams{
		Output: output,
		Base64: b64,
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

	var result rpc.ScreenshotResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return fmt.Errorf("invalid response format: %w", err)
	}

	if result.Path != "" && result.Data != "" {
		e.writer.Success(fmt.Sprintf("screenshot saved to %s and captured (base64)\n%s", result.Path, result.Data))
	} else if result.Path != "" {
		e.writer.Success(fmt.Sprintf("screenshot saved to %s", result.Path))
	} else if result.Data != "" {
		e.writer.Success(fmt.Sprintf("screenshot captured (base64)\n%s", result.Data))
	} else {
		e.writer.Success("screenshot captured")
	}

	// Record after successful execution (only if output path was provided)
	if shouldRecord && output != "" {
		recorder := recording.NewRecorderFromSocket(rpc.SocketPath())
		cmd := &script.ScreenshotCmd{
			Output: output,
			Base64: b64,
			Async:  async,
		}
		if err := recorder.Record(cmd); err != nil {
			fmt.Fprintf(os.Stderr, "autoebiten: recording failed: %v\n", err)
		}
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

	var result rpc.GetMousePositionResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return fmt.Errorf("invalid response format: %w", err)
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

	var result rpc.GetWheelPositionResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return fmt.Errorf("invalid response format: %w", err)
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
	executor.SetInputFunc(func(key, action string, durationTicks int64, async bool) error {
		return e.RunInputCommand(key, action, durationTicks, async, false)
	})

	executor.SetMouseFunc(func(action string, x, y int, button string, durationTicks int64, async bool) error {
		return e.RunMouseCommand(action, x, y, button, durationTicks, async, false)
	})

	executor.SetWheelFunc(func(x, y float64, async bool) error {
		return e.RunWheelCommand(x, y, async, false)
	})

	executor.SetScreenshotFunc(func(output string, b64 bool, async bool) error {
		return e.RunScreenshotCommand(output, b64, async, false)
	})

	executor.SetCustomFunc(func(name, request string) error {
		return e.RunCustomCommand(name, request, false)
	})

	executor.SetStateFunc(func(name, path string) error {
		return e.RunStateCommand(name, path, false)
	})

	executor.SetWaitFunc(func(condition, timeout, interval string, verbose bool) error {
		return e.RunWaitForCommand(condition, timeout, interval, verbose, false)
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

// ListCustomCommands lists all registered custom commands.
func (e *CommandExecutor) ListCustomCommands() error {
	req, err := rpc.BuildRequest("list_custom_commands", nil)
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

	var names []string
	if err := json.Unmarshal(resp.Result, &names); err != nil {
		return fmt.Errorf("invalid response format: %w", err)
	}

	if len(names) == 0 {
		e.writer.Success("no custom commands registered")
		return nil
	}

	return e.writer.PrintJSON(names)
}

// RunCustomCommand runs a custom command.
func (e *CommandExecutor) RunCustomCommand(name, request string, shouldRecord bool) error {
	params := &rpc.CustomParams{
		Name:    name,
		Request: request,
	}

	req, err := rpc.BuildRequest("custom", params)
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

	var result rpc.CustomResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return fmt.Errorf("invalid response format: %w", err)
	}

	// Record after successful execution
	if shouldRecord {
		recorder := recording.NewRecorderFromSocket(rpc.SocketPath())
		cmd := &script.CustomCmd{
			Name:    name,
			Request: request,
		}
		if err := recorder.Record(cmd); err != nil {
			fmt.Fprintf(os.Stderr, "autoebiten: recording failed: %v\n", err)
		}
	}

	e.writer.Success(result.Response)
	return nil
}

// RunVersionCommand runs the version command.
// Shows CLI version and attempts to get game version via ping.
func (e *CommandExecutor) RunVersionCommand() error {
	cliVersion := version.VersionForCLI()
	fmt.Printf("CLI version:    %s\n", cliVersion)

	// Try to get game version via ping (optional - don't fail if no game)
	req, err := rpc.BuildRequest("ping", nil)
	if err != nil {
		return err
	}

	resp, err := rpc.SendRequestSocket(req)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("Game version:   not connected\n")
		return nil
	}

	if resp.Error != nil {
		return err
	}

	var result rpc.PingResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return err
	}

	fmt.Printf("Game version:   %s\n", result.Version)
	return nil
}

// ClearRecording clears the recording file for the current game.
func (e *CommandExecutor) ClearRecording() error {
	if err := recording.Clear(rpc.SocketPath()); err != nil {
		return fmt.Errorf("failed to clear recording: %w", err)
	}
	e.writer.Success("recording cleared")
	return nil
}

// Replay replays the recorded session.
func (e *CommandExecutor) Replay(speed float64, dumpPath string) error {
	// Read recording
	reader := recording.NewReaderFromSocket(rpc.SocketPath())
	entries, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read recording: %w", err)
	}

	if len(entries) == 0 {
		return fmt.Errorf("no recording found for game")
	}

	// Generate script
	gen := recording.NewGenerator(speed)
	script, err := gen.Generate(entries)
	if err != nil {
		return fmt.Errorf("failed to generate script: %w", err)
	}

	// Either dump or execute
	if dumpPath != "" {
		data, err := json.MarshalIndent(script, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal script: %w", err)
		}
		if err := os.WriteFile(dumpPath, data, 0644); err != nil {
			return fmt.Errorf("failed to write script: %w", err)
		}
		e.writer.Success(fmt.Sprintf("script dumped to %s", dumpPath))
		return nil
	}

	// Execute via RunScriptCommand
	data, err := json.Marshal(script)
	if err != nil {
		return fmt.Errorf("failed to marshal script: %w", err)
	}
	return e.RunScriptCommand(string(data), false)
}
