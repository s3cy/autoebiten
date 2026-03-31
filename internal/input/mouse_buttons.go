package input

type MouseButton int

// MouseButton constants matching ebiten.MouseButton.
const (
	MouseButtonLeft   MouseButton = MouseButton0
	MouseButtonMiddle MouseButton = MouseButton1
	MouseButtonRight  MouseButton = MouseButton2
)

const (
	MouseButton0   MouseButton = iota // The 'left' button
	MouseButton1                      // The 'right' button
	MouseButton2                      // The 'middle' button
	MouseButton3                      // The additional button (usually browser-back)
	MouseButton4                      // The additional button (usually browser-forward)
	MouseButtonMax = MouseButton4
)

var StringMouseButtonMap = map[string]MouseButton{
	"MouseButtonLeft":   MouseButtonLeft,
	"MouseButtonMiddle": MouseButtonMiddle,
	"MouseButtonRight":  MouseButtonRight,
	"MouseButton0":      MouseButton0,
	"MouseButton1":      MouseButton1,
	"MouseButton2":      MouseButton2,
	"MouseButton3":      MouseButton3,
	"MouseButton4":      MouseButton4,
}

func LookupMouseButton(name string) (MouseButton, bool) {
	btn, ok := StringMouseButtonMap[name]
	return btn, ok
}
