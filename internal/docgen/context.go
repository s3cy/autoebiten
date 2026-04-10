package docgen

// Context holds template execution state.
type Context struct {
	GameSession *GameSession
	Config      *Config
	outputs     []string
}

// NewContext creates a new template context.
func NewContext() *Context {
	return &Context{}
}

// SetConfig sets the configuration for the template.
func (c *Context) SetConfig(cfg *Config) {
	c.Config = cfg
}

// AddOutput stores a captured output.
func (c *Context) AddOutput(output string) {
	c.outputs = append(c.outputs, output)
}

// GetOutputs returns a copy of all captured outputs.
func (c *Context) GetOutputs() []string {
	// Return a copy to prevent callers from modifying internal state
	result := make([]string, len(c.outputs))
	copy(result, c.outputs)
	return result
}
