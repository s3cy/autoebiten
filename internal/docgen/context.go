package docgen

// GameSession represents a running game instance for command execution.
// This type will be fully implemented in a later task.
type GameSession struct {
	// Placeholder - will be expanded in Task 4
}

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

// GetOutputs returns all captured outputs.
func (c *Context) GetOutputs() []string {
	return c.outputs
}