package autoebiten

import "github.com/s3cy/autoebiten/internal/custom"

// CommandContext provides context for custom command execution.
// It contains the request data and a method to send the response.
type CommandContext = custom.Context

// Register registers a custom command handler.
// The name must be unique; registering a duplicate name will override.
// The handler receives a CommandContext containing the request and a Respond method.
//
// Example:
//
//	autoebiten.Register("getPlayerInfo", func(ctx autoebiten.CommandContext) {
//		info := getPlayerInfo() // user-defined function
//		ctx.Respond(fmt.Sprintf("Health: %d, Mana: %d", info.Health, info.Mana))
//	})
func Register(name string, handler func(CommandContext)) {
	custom.Register(name, handler)
}

// Unregister removes a custom command handler.
// Returns true if the command was found and removed, false otherwise.
func Unregister(name string) bool {
	return custom.Unregister(name)
}

// GetCustomCommand returns the handler for a custom command.
// Returns nil if the command is not registered.
// This is primarily for internal use.
func GetCustomCommand(name string) func(CommandContext) {
	return custom.Get(name)
}

// ListCustomCommands returns a list of all registered custom command names.
func ListCustomCommands() []string {
	return custom.List()
}
