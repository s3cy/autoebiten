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
// Users can provide their own RWMutex via RegisterOptions if they modify
// the UI tree from goroutines. autoui acquires RLock during tree traversal
// in Update(). Users must acquire WriteLock when modifying the tree.
var treeRWLock *sync.RWMutex

// RegisterOptions holds optional configuration for autoui registration.
type RegisterOptions struct {
	// RWLock protects concurrent UI tree access.
	// autoui reads the widget tree during Update() processing (main thread).
	// If you modify the UI tree from goroutines (e.g., AddChild, RemoveChild),
	// you must provide this lock and acquire WriteLock when modifying.
	//
	// Example:
	//   var uiLock sync.RWMutex
	//   autoui.RegisterWithOptions(ui, "autoui", &autoui.RegisterOptions{
	//       RWLock: &uiLock,
	//   })
	//
	//   // In your goroutine that modifies UI:
	//   uiLock.Lock()
	//   container.AddChild(newWidget)
	//   uiLock.Unlock()
	RWLock *sync.RWMutex
}

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
//
// Note: autoui reads the widget tree during Update() processing. If you modify
// the UI tree from goroutines, use RegisterWithOptions with an RWLock.
func Register(ui *ebitenui.UI) {
	RegisterWithPrefix(ui, "autoui")
}

// RegisterWithPrefix registers all autoui commands with a custom prefix.
// The UI must be non-nil; otherwise RegisterWithPrefix panics.
// The prefix is prepended to all command names (e.g., prefix + ".tree").
//
// Note: autoui reads the widget tree during Update() processing. If you modify
// the UI tree from goroutines, use RegisterWithOptions with an RWLock.
func RegisterWithPrefix(ui *ebitenui.UI, prefix string) {
	RegisterWithOptions(ui, prefix, nil)
}

// RegisterWithOptions registers all autoui commands with a custom prefix and options.
// The UI must be non-nil; otherwise RegisterWithOptions panics.
//
// The prefix is prepended to all command names (e.g., prefix + ".tree").
//
// Thread Safety:
// autoui reads the widget tree during Update() processing (called from main thread).
// If you modify the UI tree from goroutines (AddChild, RemoveChild, etc.), you must:
//  1. Create a sync.RWMutex
//  2. Pass it via RegisterOptions.RWLock
//  3. Acquire WriteLock when modifying the UI tree from goroutines
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