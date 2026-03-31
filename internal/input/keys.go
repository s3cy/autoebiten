package input

// Key represents a keyboard key.
type Key int

// Key constants matching ebiten.Key.
const (
	KeyA Key = iota
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ
	KeyAltLeft
	KeyAltRight
	KeyArrowDown
	KeyArrowLeft
	KeyArrowRight
	KeyArrowUp
	KeyBackquote
	KeyBackslash
	KeyBackspace
	KeyBracketLeft
	KeyBracketRight
	KeyCapsLock
	KeyComma
	KeyContextMenu
	KeyControlLeft
	KeyControlRight
	KeyDelete
	KeyDigit0
	KeyDigit1
	KeyDigit2
	KeyDigit3
	KeyDigit4
	KeyDigit5
	KeyDigit6
	KeyDigit7
	KeyDigit8
	KeyDigit9
	KeyEnd
	KeyEnter
	KeyEqual
	KeyEscape
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyF13
	KeyF14
	KeyF15
	KeyF16
	KeyF17
	KeyF18
	KeyF19
	KeyF20
	KeyF21
	KeyF22
	KeyF23
	KeyF24
	KeyHome
	KeyInsert
	KeyIntlBackslash
	KeyMetaLeft
	KeyMetaRight
	KeyMinus
	KeyNumLock
	KeyNumpad0
	KeyNumpad1
	KeyNumpad2
	KeyNumpad3
	KeyNumpad4
	KeyNumpad5
	KeyNumpad6
	KeyNumpad7
	KeyNumpad8
	KeyNumpad9
	KeyNumpadAdd
	KeyNumpadDecimal
	KeyNumpadDivide
	KeyNumpadEnter
	KeyNumpadEqual
	KeyNumpadMultiply
	KeyNumpadSubtract
	KeyPageDown
	KeyPageUp
	KeyPause
	KeyPeriod
	KeyPrintScreen
	KeyQuote
	KeyScrollLock
	KeySemicolon
	KeyShiftLeft
	KeyShiftRight
	KeySlash
	KeySpace
	KeyTab
	KeyAlt
	KeyControl
	KeyShift
	KeyMeta
	KeyMax = KeyMeta
)

// StringKeyMap maps string key names to ebiten.Key values.
var StringKeyMap = map[string]Key{
	"KeyA":              KeyA,
	"KeyB":              KeyB,
	"KeyC":              KeyC,
	"KeyD":              KeyD,
	"KeyE":              KeyE,
	"KeyF":              KeyF,
	"KeyG":              KeyG,
	"KeyH":              KeyH,
	"KeyI":              KeyI,
	"KeyJ":              KeyJ,
	"KeyK":              KeyK,
	"KeyL":              KeyL,
	"KeyM":              KeyM,
	"KeyN":              KeyN,
	"KeyO":              KeyO,
	"KeyP":              KeyP,
	"KeyQ":              KeyQ,
	"KeyR":              KeyR,
	"KeyS":              KeyS,
	"KeyT":              KeyT,
	"KeyU":              KeyU,
	"KeyV":              KeyV,
	"KeyW":              KeyW,
	"KeyX":              KeyX,
	"KeyY":              KeyY,
	"KeyZ":              KeyZ,
	"KeyAltLeft":        KeyAltLeft,
	"KeyAltRight":       KeyAltRight,
	"KeyArrowDown":      KeyArrowDown,
	"KeyArrowLeft":      KeyArrowLeft,
	"KeyArrowRight":     KeyArrowRight,
	"KeyArrowUp":        KeyArrowUp,
	"KeyBackquote":      KeyBackquote,
	"KeyBackslash":      KeyBackslash,
	"KeyBackspace":      KeyBackspace,
	"KeyBracketLeft":    KeyBracketLeft,
	"KeyBracketRight":   KeyBracketRight,
	"KeyCapsLock":       KeyCapsLock,
	"KeyComma":          KeyComma,
	"KeyContextMenu":    KeyContextMenu,
	"KeyControlLeft":    KeyControlLeft,
	"KeyControlRight":   KeyControlRight,
	"KeyDelete":         KeyDelete,
	"KeyDigit0":         KeyDigit0,
	"KeyDigit1":         KeyDigit1,
	"KeyDigit2":         KeyDigit2,
	"KeyDigit3":         KeyDigit3,
	"KeyDigit4":         KeyDigit4,
	"KeyDigit5":         KeyDigit5,
	"KeyDigit6":         KeyDigit6,
	"KeyDigit7":         KeyDigit7,
	"KeyDigit8":         KeyDigit8,
	"KeyDigit9":         KeyDigit9,
	"KeyEnd":            KeyEnd,
	"KeyEnter":          KeyEnter,
	"KeyEqual":          KeyEqual,
	"KeyEscape":         KeyEscape,
	"KeyF1":             KeyF1,
	"KeyF2":             KeyF2,
	"KeyF3":             KeyF3,
	"KeyF4":             KeyF4,
	"KeyF5":             KeyF5,
	"KeyF6":             KeyF6,
	"KeyF7":             KeyF7,
	"KeyF8":             KeyF8,
	"KeyF9":             KeyF9,
	"KeyF10":            KeyF10,
	"KeyF11":            KeyF11,
	"KeyF12":            KeyF12,
	"KeyF13":            KeyF13,
	"KeyF14":            KeyF14,
	"KeyF15":            KeyF15,
	"KeyF16":            KeyF16,
	"KeyF17":            KeyF17,
	"KeyF18":            KeyF18,
	"KeyF19":            KeyF19,
	"KeyF20":            KeyF20,
	"KeyF21":            KeyF21,
	"KeyF22":            KeyF22,
	"KeyF23":            KeyF23,
	"KeyF24":            KeyF24,
	"KeyHome":           KeyHome,
	"KeyInsert":         KeyInsert,
	"KeyIntlBackslash":  KeyIntlBackslash,
	"KeyMetaLeft":       KeyMetaLeft,
	"KeyMetaRight":      KeyMetaRight,
	"KeyMinus":          KeyMinus,
	"KeyNumLock":        KeyNumLock,
	"KeyNumpad0":        KeyNumpad0,
	"KeyNumpad1":        KeyNumpad1,
	"KeyNumpad2":        KeyNumpad2,
	"KeyNumpad3":        KeyNumpad3,
	"KeyNumpad4":        KeyNumpad4,
	"KeyNumpad5":        KeyNumpad5,
	"KeyNumpad6":        KeyNumpad6,
	"KeyNumpad7":        KeyNumpad7,
	"KeyNumpad8":        KeyNumpad8,
	"KeyNumpad9":        KeyNumpad9,
	"KeyNumpadAdd":      KeyNumpadAdd,
	"KeyNumpadDecimal":  KeyNumpadDecimal,
	"KeyNumpadDivide":   KeyNumpadDivide,
	"KeyNumpadEnter":    KeyNumpadEnter,
	"KeyNumpadEqual":    KeyNumpadEqual,
	"KeyNumpadMultiply": KeyNumpadMultiply,
	"KeyNumpadSubtract": KeyNumpadSubtract,
	"KeyPageDown":       KeyPageDown,
	"KeyPageUp":         KeyPageUp,
	"KeyPause":          KeyPause,
	"KeyPeriod":         KeyPeriod,
	"KeyPrintScreen":    KeyPrintScreen,
	"KeyQuote":          KeyQuote,
	"KeyScrollLock":     KeyScrollLock,
	"KeySemicolon":      KeySemicolon,
	"KeyShiftLeft":      KeyShiftLeft,
	"KeyShiftRight":     KeyShiftRight,
	"KeySlash":          KeySlash,
	"KeySpace":          KeySpace,
	"KeyTab":            KeyTab,
	"KeyAlt":            KeyAlt,
	"KeyControl":        KeyControl,
	"KeyShift":          KeyShift,
	"KeyMeta":           KeyMeta,
}

// LookupKey looks up a key by its string name.
func LookupKey(name string) (Key, bool) {
	k, ok := StringKeyMap[name]
	return k, ok
}
