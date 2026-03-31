package script

import (
	"fmt"
	"time"
)

// Executor executes a script.
type Executor struct {
	script         *Script
	inputFunc      func(key, action string, durationTicks int64, async bool) error
	mouseFunc      func(action string, x, y int, button string, durationTicks int64, async bool) error
	wheelFunc      func(x, y float64, async bool) error
	screenshotFunc func(output string, async bool) error
	customFunc     func(name, request string) error
	commandCount   int
}

// NewExecutor creates a new script executor.
func NewExecutor(s *Script) *Executor {
	return &Executor{
		script: s,
	}
}

// SetInputFunc sets the function to call for input commands.
func (e *Executor) SetInputFunc(f func(key, action string, durationTicks int64, async bool) error) {
	e.inputFunc = f
}

// SetMouseFunc sets the function to call for mouse commands.
func (e *Executor) SetMouseFunc(f func(action string, x, y int, button string, durationTicks int64, async bool) error) {
	e.mouseFunc = f
}

// SetWheelFunc sets the function to call for wheel commands.
func (e *Executor) SetWheelFunc(f func(x, y float64, async bool) error) {
	e.wheelFunc = f
}

// SetScreenshotFunc sets the function to call for screenshot commands.
func (e *Executor) SetScreenshotFunc(f func(output string, async bool) error) {
	e.screenshotFunc = f
}

// SetCustomFunc sets the function to call for custom commands.
func (e *Executor) SetCustomFunc(f func(name, request string) error) {
	e.customFunc = f
}

// Execute runs the script.
func (e *Executor) Execute() (int, error) {
	e.commandCount = 0
	return e.commandCount, e.executeCommands(e.script.Commands)
}

func (e *Executor) executeCommands(commands []CommandWrapper) error {
	for _, cmd := range commands {
		if err := e.executeCommand(cmd); err != nil {
			return err
		}
	}
	return nil
}

func (e *Executor) executeCommand(cmd CommandWrapper) error {
	e.commandCount++

	switch c := cmd.(type) {
	case *InputCmd:
		if e.inputFunc == nil {
			return fmt.Errorf("input function not set")
		}
		if err := e.inputFunc(c.Key, c.Action, c.DurationTicks, c.Async); err != nil {
			return fmt.Errorf("%s command failed: %w", formatInputCmd(c), err)
		}
		return nil

	case *MouseCmd:
		if e.mouseFunc == nil {
			return fmt.Errorf("mouse function not set")
		}
		if err := e.mouseFunc(c.Action, c.X, c.Y, c.Button, c.DurationTicks, c.Async); err != nil {
			return fmt.Errorf("%s command failed: %w", formatMouseCmd(c), err)
		}
		return nil

	case *WheelCmd:
		if e.wheelFunc == nil {
			return fmt.Errorf("wheel function not set")
		}
		if err := e.wheelFunc(c.X, c.Y, c.Async); err != nil {
			return fmt.Errorf("%s command failed: %w", formatWheelCmd(c), err)
		}
		return nil

	case *ScreenshotCmd:
		if e.screenshotFunc == nil {
			return fmt.Errorf("screenshot function not set")
		}
		if err := e.screenshotFunc(c.Output, c.Async); err != nil {
			return fmt.Errorf("%s command failed: %w", formatScreenshotCmd(c), err)
		}
		return nil

	case *DelayCmd:
		time.Sleep(time.Duration(c.Ms) * time.Millisecond)
		return nil

	case *CustomCmd:
		if e.customFunc == nil {
			return fmt.Errorf("custom function not set")
		}
		if err := e.customFunc(c.Name, c.Request); err != nil {
			return fmt.Errorf("%s command failed: %w", formatCustomCmd(c), err)
		}
		return nil

	case *RepeatCmd:
		for i := 0; i < c.Times; i++ {
			if err := e.executeCommands(c.Commands); err != nil {
				return err
			}
		}
		return nil

	default:
		return fmt.Errorf("unknown command type: %T", cmd)
	}
}

func formatInputCmd(cmd *InputCmd) string {
	return fmt.Sprintf(`{"key": %q, "action": %q}`, cmd.Key, cmd.Action)
}

func formatMouseCmd(cmd *MouseCmd) string {
	return fmt.Sprintf(`{"action": %q, "x": %d, "y": %d, "button": %q}`, cmd.Action, cmd.X, cmd.Y, cmd.Button)
}

func formatWheelCmd(cmd *WheelCmd) string {
	return fmt.Sprintf(`{"wheel": %.2f, %.2f}`, cmd.X, cmd.Y)
}

func formatScreenshotCmd(cmd *ScreenshotCmd) string {
	return fmt.Sprintf(`{"output": %q, "async": %v}`, cmd.Output, cmd.Async)
}

func formatCustomCmd(cmd *CustomCmd) string {
	return fmt.Sprintf(`{"custom": %q, "request": %q}`, cmd.Name, cmd.Request)
}
