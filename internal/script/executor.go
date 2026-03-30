package script

import (
	"fmt"
	"time"
)

// Executor executes a script.
type Executor struct {
	script         *Script
	inputFunc      func(key, action string, durationTicks int64) error
	mouseFunc      func(action string, x, y int, button string, durationTicks int64) error
	wheelFunc      func(x, y float64) error
	screenshotFunc func(output string, async bool) error
	commandCount   int
}

// NewExecutor creates a new script executor.
func NewExecutor(s *Script) *Executor {
	return &Executor{
		script: s,
	}
}

// SetInputFunc sets the function to call for input commands.
func (e *Executor) SetInputFunc(f func(key, action string, durationTicks int64) error) {
	e.inputFunc = f
}

// SetMouseFunc sets the function to call for mouse commands.
func (e *Executor) SetMouseFunc(f func(action string, x, y int, button string, durationTicks int64) error) {
	e.mouseFunc = f
}

// SetWheelFunc sets the function to call for wheel commands.
func (e *Executor) SetWheelFunc(f func(x, y float64) error) {
	e.wheelFunc = f
}

// SetScreenshotFunc sets the function to call for screenshot commands.
func (e *Executor) SetScreenshotFunc(f func(output string, async bool) error) {
	e.screenshotFunc = f
}

// Execute runs the script.
func (e *Executor) Execute() (int, error) {
	e.commandCount = 0
	return e.commandCount, e.executeNodes(e.script.Commands)
}

func (e *Executor) executeNodes(nodes []Node) error {
	for _, node := range nodes {
		if err := e.executeNode(node); err != nil {
			return err
		}
	}
	return nil
}

func (e *Executor) executeNode(node Node) error {
	e.commandCount++

	switch cmd := node.(type) {
	case *InputCmd:
		if e.inputFunc == nil {
			return fmt.Errorf("input function not set")
		}
		if err := e.inputFunc(cmd.Key, cmd.Action, cmd.DurationTicks); err != nil {
			return fmt.Errorf("%s command failed: %w", formatInputCmd(cmd), err)
		}
		return nil

	case *MouseCmd:
		if e.mouseFunc == nil {
			return fmt.Errorf("mouse function not set")
		}
		if err := e.mouseFunc(cmd.Action, cmd.X, cmd.Y, cmd.Button, cmd.DurationTicks); err != nil {
			return fmt.Errorf("%s command failed: %w", formatMouseCmd(cmd), err)
		}
		return nil

	case *WheelCmd:
		if e.wheelFunc == nil {
			return fmt.Errorf("wheel function not set")
		}
		if err := e.wheelFunc(cmd.X, cmd.Y); err != nil {
			return fmt.Errorf("%s command failed: %w", formatWheelCmd(cmd), err)
		}
		return nil

	case *ScreenshotCmd:
		if e.screenshotFunc == nil {
			return fmt.Errorf("screenshot function not set")
		}
		if err := e.screenshotFunc(cmd.Output, cmd.Async); err != nil {
			return fmt.Errorf("%s command failed: %w", formatScreenshotCmd(cmd), err)
		}
		return nil

	case *DelayCmd:
		time.Sleep(time.Duration(cmd.Ms) * time.Millisecond)
		return nil

	case *RepeatCmd:
		for i := 0; i < cmd.Times; i++ {
			if err := e.executeNodes(cmd.Commands); err != nil {
				return err
			}
		}
		return nil

	default:
		return fmt.Errorf("unknown node type: %T", node)
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
