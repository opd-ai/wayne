//go:build windows || darwin || android || ios

package wayne

import "time"

// EventType identifies the category of an event.
type EventType int

const (
	// EventTypePointer identifies mouse/touchpad pointer events.
	EventTypePointer EventType = iota
	// EventTypeKey identifies keyboard events.
	EventTypeKey
	// EventTypeTouch identifies touch screen events.
	EventTypeTouch
	// EventTypeWindow identifies window state events.
	EventTypeWindow
	// EventTypeCustom identifies application-defined events.
	EventTypeCustom
)

// Event is the common interface for all events.
type Event interface {
	Type() EventType
	Timestamp() time.Time
	Consumed() bool
	Consume()
}

// baseEvent provides common fields for all events.
type baseEvent struct {
	timestamp time.Time
	consumed  bool
}

func (e *baseEvent) Timestamp() time.Time { return e.timestamp }
func (e *baseEvent) Consumed() bool       { return e.consumed }
func (e *baseEvent) Consume()             { e.consumed = true }

// PointerEventType specifies the type of pointer event.
type PointerEventType int

const (
	// PointerMove indicates the pointer has moved.
	PointerMove PointerEventType = iota
	// PointerButtonPress indicates a mouse button was pressed.
	PointerButtonPress
	// PointerButtonRelease indicates a mouse button was released.
	PointerButtonRelease
	// PointerScroll indicates a scroll wheel event.
	PointerScroll
	// PointerEnter indicates the pointer entered the window.
	PointerEnter
	// PointerLeave indicates the pointer left the window.
	PointerLeave
)

// PointerButton represents a mouse button.
type PointerButton uint32

const (
	// PointerButtonLeft is the left mouse button.
	PointerButtonLeft PointerButton = 0x110
	// PointerButtonRight is the right mouse button.
	PointerButtonRight PointerButton = 0x111
	// PointerButtonMiddle is the middle mouse button.
	PointerButtonMiddle PointerButton = 0x112
)

// ScrollAxis represents the scroll direction.
type ScrollAxis int

const (
	// ScrollAxisVertical indicates vertical scrolling.
	ScrollAxisVertical ScrollAxis = iota
	// ScrollAxisHorizontal indicates horizontal scrolling.
	ScrollAxisHorizontal
)

// PointerEvent represents mouse/touchpad pointer events.
type PointerEvent struct {
	baseEvent
	eventType PointerEventType
	x, y      float64
	button    PointerButton
	axis      ScrollAxis
	value     float64
}

func (e *PointerEvent) Type() EventType             { return EventTypePointer }
func (e *PointerEvent) EventType() PointerEventType { return e.eventType }
func (e *PointerEvent) X() float64                  { return e.x }
func (e *PointerEvent) Y() float64                  { return e.y }
func (e *PointerEvent) Button() PointerButton       { return e.button }
func (e *PointerEvent) Axis() ScrollAxis            { return e.axis }
func (e *PointerEvent) Value() float64              { return e.value }

// NewPointerEvent creates a new PointerEvent with the given parameters.
func NewPointerEvent(evtType PointerEventType, x, y float64, button PointerButton, axis ScrollAxis, value float64) *PointerEvent {
	return &PointerEvent{
		baseEvent: baseEvent{timestamp: time.Now()},
		eventType: evtType,
		x:         x,
		y:         y,
		button:    button,
		axis:      axis,
		value:     value,
	}
}

// KeyEventType specifies the type of keyboard event.
type KeyEventType int

const (
	// KeyPress indicates a key was pressed.
	KeyPress KeyEventType = iota
	// KeyRelease indicates a key was released.
	KeyRelease
	// KeyRepeat indicates a key is repeating (held down).
	KeyRepeat
)

// Key represents a keyboard key symbol (compatible with X11 keysyms).
type Key uint32

// Common key constants.
const (
	KeyEscape    Key = 0xFF1B
	KeyReturn    Key = 0xFF0D
	KeyTab       Key = 0xFF09
	KeyBackspace Key = 0xFF08
	KeyDelete    Key = 0xFFFF
	KeyLeft      Key = 0xFF51
	KeyUp        Key = 0xFF52
	KeyRight     Key = 0xFF53
	KeyDown      Key = 0xFF54
	KeyHome      Key = 0xFF50
	KeyEnd       Key = 0xFF57
	KeyPageUp    Key = 0xFF55
	KeyPageDown  Key = 0xFF56
	KeySpace     Key = 0x0020
	KeyShiftL    Key = 0xFFE1
	KeyShiftR    Key = 0xFFE2
	KeyControlL  Key = 0xFFE3
	KeyControlR  Key = 0xFFE4
	KeyAltL      Key = 0xFFE9
	KeyAltR      Key = 0xFFEA
	KeySuperL    Key = 0xFFEB
	KeySuperR    Key = 0xFFEC
)

// Modifier represents keyboard modifiers.
type Modifier uint32

const (
	// ModShift indicates the Shift key is held.
	ModShift Modifier = 1 << 0
	// ModControl indicates the Control key is held.
	ModControl Modifier = 1 << 1
	// ModAlt indicates the Alt key is held.
	ModAlt Modifier = 1 << 2
	// ModSuper indicates the Super (Windows/Command) key is held.
	ModSuper Modifier = 1 << 3
)

// KeyEvent represents keyboard events.
type KeyEvent struct {
	baseEvent
	eventType KeyEventType
	key       Key
	modifiers Modifier
	r         rune
}

func (e *KeyEvent) Type() EventType         { return EventTypeKey }
func (e *KeyEvent) EventType() KeyEventType { return e.eventType }
func (e *KeyEvent) Key() Key                { return e.key }
func (e *KeyEvent) Modifiers() Modifier     { return e.modifiers }
func (e *KeyEvent) Rune() rune              { return e.r }

// IsPress returns true if this is a key press or repeat event.
func (e *KeyEvent) IsPress() bool {
	return e.eventType == KeyPress || e.eventType == KeyRepeat
}

// NewKeyEvent creates a new KeyEvent.
func NewKeyEvent(evtType KeyEventType, key Key, mods Modifier, r rune) *KeyEvent {
	return &KeyEvent{
		baseEvent: baseEvent{timestamp: time.Now()},
		eventType: evtType,
		key:       key,
		modifiers: mods,
		r:         r,
	}
}

// TouchEventType specifies the type of touch event.
type TouchEventType int

const (
	// TouchDown indicates a touch point was pressed.
	TouchDown TouchEventType = iota
	// TouchUp indicates a touch point was released.
	TouchUp
	// TouchMotion indicates a touch point moved.
	TouchMotion
	// TouchCancel indicates a touch sequence was cancelled.
	TouchCancel
)

// TouchPhase is an alias for TouchEventType for API surface compatibility.
type TouchPhase = TouchEventType

// TouchEvent represents touch screen events.
type TouchEvent struct {
	baseEvent
	eventType TouchEventType
	id        int32
	x, y      float64
}

func (e *TouchEvent) Type() EventType           { return EventTypeTouch }
func (e *TouchEvent) EventType() TouchEventType { return e.eventType }

// TouchID returns the unique identifier for this touch point.
func (e *TouchEvent) TouchID() int32 { return e.id }

// ID returns the unique identifier for this touch point.
func (e *TouchEvent) ID() int32 { return e.id }

func (e *TouchEvent) X() float64 { return e.x }
func (e *TouchEvent) Y() float64 { return e.y }

// Phase returns the phase of this touch event.
func (e *TouchEvent) Phase() TouchPhase { return e.eventType }

// NewTouchEvent creates a new TouchEvent.
func NewTouchEvent(phase TouchEventType, id int32, x, y float64) *TouchEvent {
	return &TouchEvent{
		baseEvent: baseEvent{timestamp: time.Now()},
		eventType: phase,
		id:        id,
		x:         x,
		y:         y,
	}
}

// WindowEventType specifies the type of window event.
type WindowEventType int

const (
	// WindowResize indicates the window was resized.
	WindowResize WindowEventType = iota
	// WindowClose indicates the window close was requested.
	WindowClose
	// WindowFocus indicates the window gained focus.
	WindowFocus
	// WindowUnfocus indicates the window lost focus.
	WindowUnfocus
	// WindowScaleChange indicates the window's scale factor changed.
	WindowScaleChange
)

// WindowEvent represents window state events.
type WindowEvent struct {
	baseEvent
	eventType WindowEventType
	width     int
	height    int
	scale     float64
}

func (e *WindowEvent) Type() EventType            { return EventTypeWindow }
func (e *WindowEvent) EventType() WindowEventType { return e.eventType }
func (e *WindowEvent) Width() int                 { return e.width }
func (e *WindowEvent) Height() int                { return e.height }
func (e *WindowEvent) Scale() float64             { return e.scale }

// CustomEventPayload is an opaque payload for application-defined custom events.
type CustomEventPayload interface{}

// CustomEvent represents application-defined events.
type CustomEvent struct {
	baseEvent
	data CustomEventPayload
}

func (e *CustomEvent) Type() EventType          { return EventTypeCustom }
func (e *CustomEvent) Data() CustomEventPayload { return e.data }

// NewCustomEvent creates a new custom event with the given payload.
func NewCustomEvent(data CustomEventPayload) *CustomEvent {
	return &CustomEvent{
		baseEvent: baseEvent{timestamp: time.Now()},
		data:      data,
	}
}

// EventHandler is a function that processes events.
type EventHandler func(Event)
