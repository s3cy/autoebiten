package autoui

import (
	"sync"

	"github.com/ebitenui/ebitenui"
	"github.com/s3cy/autoebiten"
	"github.com/s3cy/autoebiten/integrate"
)

// uiReference holds the registered UI instance.
var uiReference *ebitenui.UI

// uiMu protects uiReference for concurrent access.
var uiMu sync.RWMutex

// treeRWLock protects the widget tree during traversal.
// Set via RegisterOptions.RWLock for users who modify UI from goroutines.
var treeRWLock *sync.RWMutex

// RegisterOptions holds optional configuration for autoui registration.
type RegisterOptions struct {
	// RWLock protects concurrent UI tree access.
	// Provide this if you modify the UI tree from goroutines.
	// Acquire WriteLock when modifying, autoui acquires RLock during traversal.
	RWLock *sync.RWMutex
}

// Register registers all autoui commands with the default "autoui." prefix.
// The UI must be non-nil; otherwise Register panics.
func Register(ui *ebitenui.UI) {
	RegisterWithPrefix(ui, "autoui")
}

// RegisterWithPrefix registers all autoui commands with a custom prefix.
// The UI must be non-nil; otherwise RegisterWithPrefix panics.
func RegisterWithPrefix(ui *ebitenui.UI, prefix string) {
	RegisterWithOptions(ui, prefix, nil)
}

// RegisterWithOptions registers all autoui commands with a custom prefix and options.
// The UI must be non-nil; otherwise RegisterWithOptions panics.
// Use RegisterOptions.RWLock if you modify the UI tree from goroutines.
//
// Example:
//
//	var uiLock sync.RWMutex
//	ui := &ebitenui.UI{Container: root}
//	autoui.RegisterWithOptions(ui, "autoui", &autoui.RegisterOptions{
//	    RWLock: &uiLock,
//	})
//
//	// In a goroutine that modifies the UI:
//	uiLock.Lock()
//	container.AddChild(newWidget)
//	uiLock.Unlock()
func RegisterWithOptions(ui *ebitenui.UI, prefix string, opts *RegisterOptions) {
	if ui == nil {
		panic("autoui.RegisterWithOptions: UI cannot be nil")
	}

	// Store UI reference
	uiMu.Lock()
	uiReference = ui
	uiMu.Unlock()

	// Store optional tree RWLock
	if opts != nil && opts.RWLock != nil {
		treeRWLock = opts.RWLock
	} else {
		treeRWLock = nil
	}

	// Register highlight callback for patch method
	integrate.RegisterDrawHighlights(drawHighlightsCallback)

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