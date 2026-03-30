package input

import "github.com/hajimehoshi/ebiten/v2"

// Key represents a keyboard key.
type Key = ebiten.Key

// Key constants matching ebiten.Key.
const (
	KeyA              = ebiten.KeyA
	KeyB              = ebiten.KeyB
	KeyC              = ebiten.KeyC
	KeyD              = ebiten.KeyD
	KeyE              = ebiten.KeyE
	KeyF              = ebiten.KeyF
	KeyG              = ebiten.KeyG
	KeyH              = ebiten.KeyH
	KeyI              = ebiten.KeyI
	KeyJ              = ebiten.KeyJ
	KeyK              = ebiten.KeyK
	KeyL              = ebiten.KeyL
	KeyM              = ebiten.KeyM
	KeyN              = ebiten.KeyN
	KeyO              = ebiten.KeyO
	KeyP              = ebiten.KeyP
	KeyQ              = ebiten.KeyQ
	KeyR              = ebiten.KeyR
	KeyS              = ebiten.KeyS
	KeyT              = ebiten.KeyT
	KeyU              = ebiten.KeyU
	KeyV              = ebiten.KeyV
	KeyW              = ebiten.KeyW
	KeyX              = ebiten.KeyX
	KeyY              = ebiten.KeyY
	KeyZ              = ebiten.KeyZ
	KeyAltLeft        = ebiten.KeyAltLeft
	KeyAltRight       = ebiten.KeyAltRight
	KeyArrowDown      = ebiten.KeyArrowDown
	KeyArrowLeft      = ebiten.KeyArrowLeft
	KeyArrowRight     = ebiten.KeyArrowRight
	KeyArrowUp        = ebiten.KeyArrowUp
	KeyBackquote      = ebiten.KeyBackquote
	KeyBackslash      = ebiten.KeyBackslash
	KeyBackspace      = ebiten.KeyBackspace
	KeyBracketLeft    = ebiten.KeyBracketLeft
	KeyBracketRight   = ebiten.KeyBracketRight
	KeyCapsLock       = ebiten.KeyCapsLock
	KeyComma          = ebiten.KeyComma
	KeyContextMenu    = ebiten.KeyContextMenu
	KeyControlLeft    = ebiten.KeyControlLeft
	KeyControlRight   = ebiten.KeyControlRight
	KeyDelete         = ebiten.KeyDelete
	KeyDigit0         = ebiten.KeyDigit0
	KeyDigit1         = ebiten.KeyDigit1
	KeyDigit2         = ebiten.KeyDigit2
	KeyDigit3         = ebiten.KeyDigit3
	KeyDigit4         = ebiten.KeyDigit4
	KeyDigit5         = ebiten.KeyDigit5
	KeyDigit6         = ebiten.KeyDigit6
	KeyDigit7         = ebiten.KeyDigit7
	KeyDigit8         = ebiten.KeyDigit8
	KeyDigit9         = ebiten.KeyDigit9
	KeyEnd            = ebiten.KeyEnd
	KeyEnter          = ebiten.KeyEnter
	KeyEqual          = ebiten.KeyEqual
	KeyEscape         = ebiten.KeyEscape
	KeyF1             = ebiten.KeyF1
	KeyF2             = ebiten.KeyF2
	KeyF3             = ebiten.KeyF3
	KeyF4             = ebiten.KeyF4
	KeyF5             = ebiten.KeyF5
	KeyF6             = ebiten.KeyF6
	KeyF7             = ebiten.KeyF7
	KeyF8             = ebiten.KeyF8
	KeyF9             = ebiten.KeyF9
	KeyF10            = ebiten.KeyF10
	KeyF11            = ebiten.KeyF11
	KeyF12            = ebiten.KeyF12
	KeyF13            = ebiten.KeyF13
	KeyF14            = ebiten.KeyF14
	KeyF15            = ebiten.KeyF15
	KeyF16            = ebiten.KeyF16
	KeyF17            = ebiten.KeyF17
	KeyF18            = ebiten.KeyF18
	KeyF19            = ebiten.KeyF19
	KeyF20            = ebiten.KeyF20
	KeyF21            = ebiten.KeyF21
	KeyF22            = ebiten.KeyF22
	KeyF23            = ebiten.KeyF23
	KeyF24            = ebiten.KeyF24
	KeyHome           = ebiten.KeyHome
	KeyInsert         = ebiten.KeyInsert
	KeyIntlBackslash  = ebiten.KeyIntlBackslash
	KeyMetaLeft       = ebiten.KeyMetaLeft
	KeyMetaRight      = ebiten.KeyMetaRight
	KeyMinus          = ebiten.KeyMinus
	KeyNumLock        = ebiten.KeyNumLock
	KeyNumpad0        = ebiten.KeyNumpad0
	KeyNumpad1        = ebiten.KeyNumpad1
	KeyNumpad2        = ebiten.KeyNumpad2
	KeyNumpad3        = ebiten.KeyNumpad3
	KeyNumpad4        = ebiten.KeyNumpad4
	KeyNumpad5        = ebiten.KeyNumpad5
	KeyNumpad6        = ebiten.KeyNumpad6
	KeyNumpad7        = ebiten.KeyNumpad7
	KeyNumpad8        = ebiten.KeyNumpad8
	KeyNumpad9        = ebiten.KeyNumpad9
	KeyNumpadAdd      = ebiten.KeyNumpadAdd
	KeyNumpadDecimal  = ebiten.KeyNumpadDecimal
	KeyNumpadDivide   = ebiten.KeyNumpadDivide
	KeyNumpadEnter    = ebiten.KeyNumpadEnter
	KeyNumpadEqual    = ebiten.KeyNumpadEqual
	KeyNumpadMultiply = ebiten.KeyNumpadMultiply
	KeyNumpadSubtract = ebiten.KeyNumpadSubtract
	KeyPageDown       = ebiten.KeyPageDown
	KeyPageUp         = ebiten.KeyPageUp
	KeyPause          = ebiten.KeyPause
	KeyPeriod         = ebiten.KeyPeriod
	KeyPrintScreen    = ebiten.KeyPrintScreen
	KeyQuote          = ebiten.KeyQuote
	KeyScrollLock     = ebiten.KeyScrollLock
	KeySemicolon      = ebiten.KeySemicolon
	KeyShiftLeft      = ebiten.KeyShiftLeft
	KeyShiftRight     = ebiten.KeyShiftRight
	KeySlash          = ebiten.KeySlash
	KeySpace          = ebiten.KeySpace
	KeyTab            = ebiten.KeyTab
	KeyAlt            = ebiten.KeyAlt
	KeyControl        = ebiten.KeyControl
	KeyShift          = ebiten.KeyShift
	KeyMeta           = ebiten.KeyMeta
	KeyMax            = ebiten.KeyMeta
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
