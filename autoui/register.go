package autoui

import (
	"sync"

	"github.com/ebitenui/ebitenui"
	"github.com/s3cy/autoebiten"
)

// uiReference holds the registered UI instance.
var uiReference *ebitenui.UI

// uiMu protects uiReference for concurrent access.
var uiMu sync.RWMutex

// Register registers all autoui commands with the default "autoui." prefix.
// The UI must be non-nil; otherwise Register panics.
//
// Commands registered:
//   - autoui.tree - Returns full widget tree as XML
//   - autoui.at - Returns widget at coordinates (x,y or JSON)
//   - autoui.find - Returns widgets matching query (key=value or JSON)
//   - autoui.xpath - Returns widgets matching XPath expression
//   - autoui.call - Invokes method on widget (JSON request)
//   - autoui.highlight - Adds visual highlights (clear, coordinates, or query)
func Register(ui *ebitenui.UI) {
	RegisterWithPrefix(ui, "autoui")
}

// RegisterWithPrefix registers all autoui commands with a custom prefix.
// The UI must be non-nil; otherwise RegisterWithPrefix panics.
// The prefix is prepended to all command names (e.g., prefix + ".tree").
func RegisterWithPrefix(ui *ebitenui.UI, prefix string) {
	if ui == nil {
		panic("autoui.RegisterWithPrefix: UI cannot be nil")
	}

	// Store UI reference
	uiMu.Lock()
	uiReference = ui
	uiMu.Unlock()

	// Register all command handlers
	registerCommands(prefix)
}

// registerCommands registers all command handlers with the given prefix.
func registerCommands(prefix string) {
	// Get UI reference for handlers
	uiMu.RLock()
	ui := uiReference
	uiMu.RUnlock()

	// Register tree command
	autoebiten.Register(prefix+".tree", handleTreeCommand(ui))

	// Register at command
	autoebiten.Register(prefix+".at", handleAtCommand(ui))

	// Register find command
	autoebiten.Register(prefix+".find", handleFindCommand(ui))

	// Register xpath command
	autoebiten.Register(prefix+".xpath", handleXPathCommand(ui))

	// Register call command
	autoebiten.Register(prefix+".call", handleCallCommand(ui))

	// Register highlight command
	autoebiten.Register(prefix+".highlight", handleHighlightCommand(ui))
}

// GetUI returns the registered UI instance.
// Returns nil if no UI has been registered.
func GetUI() *ebitenui.UI {
	uiMu.RLock()
	defer uiMu.RUnlock()
	return uiReference
}