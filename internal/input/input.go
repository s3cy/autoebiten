package input

import (
	"sync"
)

// VirtualInput manages virtual keyboard input state.
type VirtualInput struct {
	mu sync.RWMutex

	keyPressedTimes          [KeyMax + 1]InputTime
	keyReleasedTimes         [KeyMax + 1]InputTime
	mouseButtonPressedTimes  [MouseButtonMax + 1]InputTime
	mouseButtonReleasedTimes [MouseButtonMax + 1]InputTime

	cursorX int
	cursorY int
	wheelX  float64
	wheelY  float64
}

func NewVirtualInput() *VirtualInput {
	return &VirtualInput{}
}

func (v *VirtualInput) InjectKeyPress(key Key, inputTime InputTime) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.keyPressedTimes[key] = inputTime
}

func (v *VirtualInput) InjectKeyRelease(key Key, inputTime InputTime) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.keyReleasedTimes[key] = inputTime
}

func (v *VirtualInput) InjectKeyHold(key Key, inputTime InputTime, durationTicks int64) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.keyPressedTimes[key] = inputTime
	v.keyReleasedTimes[key] = NewInputTimeFromTick(inputTime.Tick()+durationTicks, 0)
}

func (v *VirtualInput) InjectMouseButtonPress(button MouseButton, inputTime InputTime) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.mouseButtonPressedTimes[button] = inputTime
}

func (v *VirtualInput) InjectMouseButtonRelease(button MouseButton, inputTime InputTime) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.mouseButtonReleasedTimes[button] = inputTime
}

func (v *VirtualInput) InjectMouseButtonHold(button MouseButton, inputTime InputTime, durationTicks int64) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.mouseButtonPressedTimes[button] = inputTime
	v.mouseButtonReleasedTimes[button] = NewInputTimeFromTick(inputTime.Tick()+durationTicks, 0)
}

func (v *VirtualInput) InjectCursorMove(x, y int) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.cursorX = x
	v.cursorY = y
}

func (v *VirtualInput) InjectWheelMove(x, y float64) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.wheelX = x
	v.wheelY = y
}

func (v *VirtualInput) IsKeyPressed(key Key, tick int64) bool {
	v.mu.Lock()
	defer v.mu.Unlock()

	switch key {
	case KeyAlt:
		return v.isKeyPressed(KeyAlt, tick) ||
			v.isKeyPressed(KeyAltLeft, tick) ||
			v.isKeyPressed(KeyAltRight, tick)
	case KeyControl:
		return v.isKeyPressed(KeyControl, tick) ||
			v.isKeyPressed(KeyControlLeft, tick) ||
			v.isKeyPressed(KeyControlRight, tick)
	case KeyShift:
		return v.isKeyPressed(KeyShift, tick) ||
			v.isKeyPressed(KeyShiftLeft, tick) ||
			v.isKeyPressed(KeyShiftRight, tick)
	case KeyMeta:
		return v.isKeyPressed(KeyMeta, tick) ||
			v.isKeyPressed(KeyMetaLeft, tick) ||
			v.isKeyPressed(KeyMetaRight, tick)
	}

	return v.isKeyPressed(key, tick)
}

func (v *VirtualInput) isKeyPressed(key Key, tick int64) bool {
	if key < 0 || KeyMax < key {
		return false
	}
	p := v.keyPressedTimes[key]
	r := v.keyReleasedTimes[key]
	return inputStatePressed(p, r, tick)
}

func (v *VirtualInput) IsKeyJustPressed(key Key, tick int64) bool {
	v.mu.Lock()
	defer v.mu.Unlock()

	if key < 0 || KeyMax < key {
		return false
	}
	p := v.keyPressedTimes[key]
	return inputStateJustPressed(p, tick)
}

func (v *VirtualInput) IsKeyJustReleased(key Key, tick int64) bool {
	v.mu.Lock()
	defer v.mu.Unlock()

	if key < 0 || KeyMax < key {
		return false
	}
	r := v.keyReleasedTimes[key]
	return inputStateJustReleased(r, tick)
}

func (v *VirtualInput) KeyPressDuration(key Key, tick int64) int64 {
	v.mu.Lock()
	defer v.mu.Unlock()

	if key < 0 || KeyMax < key {
		return 0
	}
	p := v.keyPressedTimes[key]
	r := v.keyReleasedTimes[key]
	return inputStateDuration(p, r, tick)
}

func (v *VirtualInput) IsMouseButtonPressed(button MouseButton, tick int64) bool {
	v.mu.Lock()
	defer v.mu.Unlock()

	if button < 0 || MouseButtonMax < button {
		return false
	}
	p := v.mouseButtonPressedTimes[button]
	r := v.mouseButtonReleasedTimes[button]
	return inputStatePressed(p, r, tick)
}

func (v *VirtualInput) IsMouseButtonJustPressed(button MouseButton, tick int64) bool {
	v.mu.Lock()
	defer v.mu.Unlock()

	if button < 0 || MouseButtonMax < button {
		return false
	}
	p := v.mouseButtonPressedTimes[button]
	return inputStateJustPressed(p, tick)
}

func (v *VirtualInput) IsMouseButtonJustReleased(button MouseButton, tick int64) bool {
	v.mu.Lock()
	defer v.mu.Unlock()

	if button < 0 || MouseButtonMax < button {
		return false
	}
	r := v.mouseButtonReleasedTimes[button]
	return inputStateJustReleased(r, tick)
}

func (v *VirtualInput) MouseButtonPressDuration(button MouseButton, tick int64) int64 {
	v.mu.Lock()
	defer v.mu.Unlock()

	if button < 0 || MouseButtonMax < button {
		return 0
	}
	p := v.mouseButtonPressedTimes[button]
	r := v.mouseButtonReleasedTimes[button]
	return inputStateDuration(p, r, tick)
}

func (v *VirtualInput) CursorPosition() (x, y int) {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.cursorX, v.cursorY
}

func (v *VirtualInput) Wheel() (x, y float64) {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.wheelX, v.wheelY
}

func inputStatePressed(pressed, released InputTime, tick int64) bool {
	return released < pressed || tick < released.Tick() || inputStateJustPressed(pressed, tick)
}

func inputStateJustPressed(pressed InputTime, tick int64) bool {
	return pressed > 0 && pressed.Tick() == tick
}

func inputStateJustReleased(released InputTime, tick int64) bool {
	return released > 0 && released.Tick() == tick
}

func inputStateDuration(pressed, released InputTime, tick int64) int64 {
	if pressed == 0 {
		return 0
	}
	if pressed < released && released.Tick() <= tick {
		return 0
	}
	return tick - pressed.Tick() + 1
}

// Get returns the global VirtualInput instance.
func Get() *VirtualInput {
	return globalInput
}

var globalInput = NewVirtualInput()
