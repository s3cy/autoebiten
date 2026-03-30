package input

import "github.com/hajimehoshi/ebiten/v2"

type MouseButton = ebiten.MouseButton

// MouseButtons
const (
	MouseButtonLeft   MouseButton = MouseButton0
	MouseButtonMiddle MouseButton = MouseButton1
	MouseButtonRight  MouseButton = MouseButton2

	MouseButton0   MouseButton = ebiten.MouseButton0
	MouseButton1   MouseButton = ebiten.MouseButton1
	MouseButton2   MouseButton = ebiten.MouseButton2
	MouseButton3   MouseButton = ebiten.MouseButton3
	MouseButton4   MouseButton = ebiten.MouseButton4
	MouseButtonMax MouseButton = ebiten.MouseButton4
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
