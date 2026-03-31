// Package custom provides custom command registration and execution.
package custom

import (
	"fmt"
	"sync"
)

// Context provides context for custom command execution.
type Context interface {
	// Request returns the request string sent from the CLI.
	Request() string

	// Respond sends the response back to the CLI.
	// This can be called immediately or deferred to a later time.
	// Only the first call to Respond will send the response;
	// subsequent calls are ignored.
	Respond(response string)
}

// context implements Context.
type context struct {
	request  string
	respond  func(string)
	called   bool
}

func (c *context) Request() string {
	return c.request
}

func (c *context) Respond(response string) {
	if c.called {
		return
	}
	c.called = true
	c.respond(response)
}

// NewContext creates a new Context.
func NewContext(request string, respond func(string)) Context {
	return &context{
		request: request,
		respond: respond,
		called:  false,
	}
}

var (
	// commands stores registered custom command handlers.
	commands   = make(map[string]func(Context))
	commandsMu sync.RWMutex
)

// Register registers a custom command handler.
// The name must be unique; registering a duplicate name will panic.
func Register(name string, handler func(Context)) {
	if name == "" {
		panic("custom.Register: command name cannot be empty")
	}
	if handler == nil {
		panic("custom.Register: handler cannot be nil")
	}

	commandsMu.Lock()
	defer commandsMu.Unlock()

	if _, exists := commands[name]; exists {
		panic(fmt.Sprintf("custom.Register: command %q already registered", name))
	}

	commands[name] = handler
}

// Unregister removes a custom command handler.
// Returns true if the command was found and removed, false otherwise.
func Unregister(name string) bool {
	commandsMu.Lock()
	defer commandsMu.Unlock()

	if _, exists := commands[name]; !exists {
		return false
	}

	delete(commands, name)
	return true
}

// Get returns the handler for a custom command.
// Returns nil if the command is not registered.
func Get(name string) func(Context) {
	commandsMu.RLock()
	defer commandsMu.RUnlock()

	return commands[name]
}

// List returns a list of all registered custom command names.
func List() []string {
	commandsMu.RLock()
	defer commandsMu.RUnlock()

	names := make([]string, 0, len(commands))
	for name := range commands {
		names = append(names, name)
	}
	return names
}
