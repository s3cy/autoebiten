package script

import "encoding/json"

// Node represents a node in the script AST.
type Node interface {
	nodeType() string
}

// Script is the root node of a script.
type Script struct {
	Version  string `json:"version"`
	Commands []Node `json:"commands"`
}

// InputCmd represents an input command.
type InputCmd struct {
	Action        string `json:"action"`
	Key           string `json:"key"`
	DurationTicks int64  `json:"duration_ticks"`
	Async         bool   `json:"async"`
}

func (InputCmd) nodeType() string { return "input" }

// MouseCmd represents a mouse command.
type MouseCmd struct {
	Action        string `json:"action"`
	X             int    `json:"x"`
	Y             int    `json:"y"`
	Button        string `json:"button"`
	DurationTicks int64  `json:"duration_ticks"`
}

func (MouseCmd) nodeType() string { return "mouse" }

// WheelCmd represents a wheel command.
type WheelCmd struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (WheelCmd) nodeType() string { return "wheel" }

// ScreenshotCmd represents a screenshot command.
type ScreenshotCmd struct {
	Output string `json:"output"`
	Async  bool   `json:"async"`
}

func (ScreenshotCmd) nodeType() string { return "screenshot" }

// DelayCmd represents a delay command.
type DelayCmd struct {
	Ms int64 `json:"ms"`
}

func (DelayCmd) nodeType() string { return "delay" }

// RepeatCmdRaw is used during parsing to hold raw JSON commands.
type RepeatCmdRaw struct {
	Times    int               `json:"times"`
	Commands []json.RawMessage `json:"commands"`
}

// RepeatCmd represents a repeat block.
type RepeatCmd struct {
	Times    int    `json:"times"`
	Commands []Node `json:"commands"`
}

func (RepeatCmd) nodeType() string { return "repeat" }
