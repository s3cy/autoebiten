package script

import (
	"encoding/json"
	"fmt"
	"os"
)

// Parse reads and parses a script file.
func Parse(path string) (*Script, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return ParseBytes(data)
}

// ParseString parses a script from a JSON string.
func ParseString(s string) (*Script, error) {
	return ParseBytes([]byte(s))
}

// ParseBytes parses a script from JSON bytes.
func ParseBytes(data []byte) (*Script, error) {
	var raw struct {
		Version  string          `json:"version"`
		Commands []json.RawMessage `json:"commands"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if raw.Version == "" {
		return nil, fmt.Errorf("missing version field")
	}

	commands := make([]Node, 0, len(raw.Commands))
	for i, cmdData := range raw.Commands {
		cmd, err := parseCommand(cmdData)
		if err != nil {
			return nil, fmt.Errorf("command %d: %w", i, err)
		}
		commands = append(commands, cmd)
	}

	return &Script{
		Version:  raw.Version,
		Commands: commands,
	}, nil
}

func parseCommand(data json.RawMessage) (Node, error) {
	// First, determine the command type
	var typeCheck struct {
		Input      *InputCmd      `json:"input"`
		Mouse      *MouseCmd      `json:"mouse"`
		Wheel      *WheelCmd      `json:"wheel"`
		Screenshot *ScreenshotCmd `json:"screenshot"`
		Delay      *DelayCmd      `json:"delay"`
		Repeat     *RepeatCmdRaw  `json:"repeat"`
	}

	if err := json.Unmarshal(data, &typeCheck); err != nil {
		return nil, fmt.Errorf("failed to parse command: %w", err)
	}

	switch {
	case typeCheck.Input != nil:
		return typeCheck.Input, nil
	case typeCheck.Mouse != nil:
		return typeCheck.Mouse, nil
	case typeCheck.Wheel != nil:
		return typeCheck.Wheel, nil
	case typeCheck.Screenshot != nil:
		return typeCheck.Screenshot, nil
	case typeCheck.Delay != nil:
		return typeCheck.Delay, nil
	case typeCheck.Repeat != nil:
		// Parse nested commands
		repeatCommands := make([]Node, 0, len(typeCheck.Repeat.Commands))
		for i, cmdData := range typeCheck.Repeat.Commands {
			cmd, err := parseCommand(cmdData)
			if err != nil {
				return nil, fmt.Errorf("repeat command %d: %w", i, err)
			}
			repeatCommands = append(repeatCommands, cmd)
		}
		return &RepeatCmd{
			Times:    typeCheck.Repeat.Times,
			Commands: repeatCommands,
		}, nil
	default:
		return nil, fmt.Errorf("unknown command type")
	}
}
