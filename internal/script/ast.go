package script

import (
	"encoding/json"
	"fmt"
)

// CommandWrapper is the interface for all command types.
type CommandWrapper interface {
	commandType() string
}

// Script is the root node of a script.
type Script struct {
	Version  string           `json:"version"`
	Commands []CommandWrapper `json:"commands"`
}

// InputCmd represents an input command.
type InputCmd struct {
	Action        string `json:"action" jsonschema:"enum=press,enum=release,enum=hold,description=Action to perform"`
	Key           string `json:"key" jsonschema:"description=Key name (use 'autoebiten keys' to list all)"`
	DurationTicks int64  `json:"duration_ticks" jsonschema:"default=6,description=Duration in game ticks for hold action"`
	Async         bool   `json:"async" jsonschema:"default=false,description=Return immediately without waiting"`
}

func (InputCmd) commandType() string { return "input" }

// MouseCmd represents a mouse command.
type MouseCmd struct {
	Action        string `json:"action" jsonschema:"enum=position,enum=press,enum=release,enum=hold,description=Action to perform (default is position or hold when button is used)"`
	X             int    `json:"x" jsonschema:"default=0,description=X coordinate"`
	Y             int    `json:"y" jsonschema:"default=0,description=Y coordinate"`
	Button        string `json:"button" jsonschema:"description=Mouse button (use 'autoebiten mouse_buttons' to list all)"`
	DurationTicks int64  `json:"duration_ticks" jsonschema:"default=6,description=Duration in game ticks for hold action"`
	Async         bool   `json:"async" jsonschema:"default=false,description=Return immediately without waiting"`
}

func (MouseCmd) commandType() string { return "mouse" }

// WheelCmd represents a wheel command.
type WheelCmd struct {
	X     float64 `json:"x" jsonschema:"default=0,description=Horizontal scroll"`
	Y     float64 `json:"y" jsonschema:"default=0,description=Vertical scroll"`
	Async bool    `json:"async" jsonschema:"default=false,description=Return immediately without waiting"`
}

func (WheelCmd) commandType() string { return "wheel" }

// ScreenshotCmd represents a screenshot command.
type ScreenshotCmd struct {
	Output string `json:"output" jsonschema:"description=Output file path (optional, auto-generated if not provided)"`
	Base64 bool   `json:"base64" jsonschema:"default=false,description=Return screenshot as base64 string in the response instead of saving to a file"`
	Async  bool   `json:"async" jsonschema:"default=false,description=Return immediately without waiting"`
}

func (ScreenshotCmd) commandType() string { return "screenshot" }

// DelayCmd represents a delay command.
type DelayCmd struct {
	Ms int64 `json:"ms" jsonschema:"minimum=0,description=Milliseconds to wait"`
}

func (DelayCmd) commandType() string { return "delay" }

// CustomCmd represents a custom command.
type CustomCmd struct {
	Name    string `json:"name" jsonschema:"description=Name of the registered custom command"`
	Request string `json:"request,omitempty" jsonschema:"description=Optional request data to pass to the custom command"`
}

func (CustomCmd) commandType() string { return "custom" }

// RepeatCmd represents a repeat block.
type RepeatCmd struct {
	Times    int              `json:"times"`
	Commands []CommandWrapper `json:"commands"`
}

func (RepeatCmd) commandType() string { return "repeat" }

// internalWrapper is used for JSON unmarshaling.
type internalWrapper struct {
	Input      *InputCmd        `json:"input,omitempty"`
	Mouse      *MouseCmd        `json:"mouse,omitempty"`
	Wheel      *WheelCmd        `json:"wheel,omitempty"`
	Screenshot *ScreenshotCmd   `json:"screenshot,omitempty"`
	Delay      *DelayCmd        `json:"delay,omitempty"`
	Custom     *CustomCmd       `json:"custom,omitempty"`
	Repeat     *json.RawMessage `json:"repeat,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshaling for CommandWrapper.
func (s *Script) UnmarshalJSON(data []byte) error {
	var raw struct {
		Version  string            `json:"version"`
		Commands []json.RawMessage `json:"commands"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	s.Version = raw.Version
	s.Commands = make([]CommandWrapper, 0, len(raw.Commands))

	for _, cmdData := range raw.Commands {
		cmd, err := unmarshalCommand(cmdData)
		if err != nil {
			return err
		}
		s.Commands = append(s.Commands, cmd)
	}

	return nil
}

// unmarshalCommand unmarshals a single command from JSON.
func unmarshalCommand(data []byte) (CommandWrapper, error) {
	return UnmarshalCommand(data)
}

// UnmarshalCommand unmarshals a single command from JSON.
// This is the public version that can be used by external packages.
func UnmarshalCommand(data []byte) (CommandWrapper, error) {
	var w internalWrapper
	if err := json.Unmarshal(data, &w); err != nil {
		return nil, err
	}

	switch {
	case w.Input != nil:
		return w.Input, nil
	case w.Mouse != nil:
		return w.Mouse, nil
	case w.Wheel != nil:
		return w.Wheel, nil
	case w.Screenshot != nil:
		return w.Screenshot, nil
	case w.Delay != nil:
		return w.Delay, nil
	case w.Custom != nil:
		return w.Custom, nil
	case w.Repeat != nil:
		return unmarshalRepeat(*w.Repeat)
	default:
		return nil, fmt.Errorf("unknown command type")
	}
}

// unmarshalRepeat unmarshals a repeat command.
func unmarshalRepeat(data []byte) (*RepeatCmd, error) {
	var raw struct {
		Times    int               `json:"times"`
		Commands []json.RawMessage `json:"commands"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	commands := make([]CommandWrapper, 0, len(raw.Commands))
	for _, cmdData := range raw.Commands {
		cmd, err := unmarshalCommand(cmdData)
		if err != nil {
			return nil, err
		}
		commands = append(commands, cmd)
	}

	return &RepeatCmd{
		Times:    raw.Times,
		Commands: commands,
	}, nil
}

// CommandSchema represents a command in the script format for schema generation.
// Only one field should be set at a time.
type CommandSchema struct {
	// Input command - inject keyboard input
	Input *InputCmd `json:"input,omitempty" jsonschema:"oneof_required=input,description=Inject keyboard input"`

	// Mouse command - inject mouse input
	Mouse *MouseCmd `json:"mouse,omitempty" jsonschema:"oneof_required=mouse,description=Inject mouse input"`

	// Wheel command - inject wheel/scroll input
	Wheel *WheelCmd `json:"wheel,omitempty" jsonschema:"oneof_required=wheel,description=Inject wheel/scroll input"`

	// Screenshot command - capture game screenshot
	Screenshot *ScreenshotCmd `json:"screenshot,omitempty" jsonschema:"oneof_required=screenshot,description=Capture game screenshot"`

	// Delay command - pause execution
	Delay *DelayCmd `json:"delay,omitempty" jsonschema:"oneof_required=delay,description=Pause execution for a duration"`

	// Custom command - execute a registered custom command
	Custom *CustomCmd `json:"custom,omitempty" jsonschema:"oneof_required=custom,description=Execute a registered custom command"`

	// Repeat command - repeat a block of commands
	Repeat *RepeatSchema `json:"repeat,omitempty" jsonschema:"oneof_required=repeat,description=Repeat a block of commands"`
}

// RepeatSchema represents a repeat block for schema generation.
type RepeatSchema struct {
	Times    int             `json:"times" jsonschema:"minimum=1,description=Number of times to repeat"`
	Commands []CommandSchema `json:"commands" jsonschema:"description=Commands to repeat"`
}

// ScriptSchema represents the root script structure for schema generation.
type ScriptSchema struct {
	// Schema is the JSON Schema URI for this document
	Schema string `json:"$schema,omitempty" jsonschema:"type=string,format=uri-reference,description=JSON Schema URI for this document"`

	Version string `json:"version" jsonschema:"enum=1.0,description=Script format version"`

	// Commands to execute in order
	Commands []CommandSchema `json:"commands" jsonschema:"description=List of commands to execute in order"`
}
