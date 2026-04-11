package autoui

import (
	"sync"

	"github.com/ebitenui/ebitenui/widget"
)

// radioGroupRegistry stores registered RadioGroup instances.
// Access is protected by mutex for thread safety.
var radioGroupRegistry = struct {
	sync.RWMutex
	groups map[string]*widget.RadioGroup
}{
	groups: make(map[string]*widget.RadioGroup),
}

// RegisterRadioGroup registers a RadioGroup with a unique name.
// If a RadioGroup with the same name exists, it will be replaced.
// Use this to make RadioGroups discoverable via autoui.call and autoui.tree.
func RegisterRadioGroup(name string, rg *widget.RadioGroup) {
	radioGroupRegistry.Lock()
	radioGroupRegistry.groups[name] = rg
	radioGroupRegistry.Unlock()
}

// UnregisterRadioGroup removes a registered RadioGroup by name.
// Call this when the RadioGroup is no longer needed.
func UnregisterRadioGroup(name string) {
	radioGroupRegistry.Lock()
	delete(radioGroupRegistry.groups, name)
	radioGroupRegistry.Unlock()
}

// GetRadioGroup returns a registered RadioGroup by name, or nil if not found.
func GetRadioGroup(name string) *widget.RadioGroup {
	radioGroupRegistry.RLock()
	defer radioGroupRegistry.RUnlock()
	return radioGroupRegistry.groups[name]
}

// GetRegisteredRadioGroups returns all registered RadioGroup names.
// Used by tree traversal to inject RadioGroups into output.
func GetRegisteredRadioGroups() []string {
	radioGroupRegistry.RLock()
	defer radioGroupRegistry.RUnlock()

	names := make([]string, 0, len(radioGroupRegistry.groups))
	for name := range radioGroupRegistry.groups {
		names = append(names, name)
	}
	return names
}
