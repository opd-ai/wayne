//go:build windows || darwin || android || ios

package wayne

import (
	"testing"
	"time"
)

func TestPointerEventCreation(t *testing.T) {
	tests := []struct {
		name      string
		evtType   PointerEventType
		x, y      float64
		button    PointerButton
		axis      ScrollAxis
		value     float64
	}{
		{
			name:    "pointer move",
			evtType: PointerMove,
			x:       100.5,
			y:       200.5,
		},
		{
			name:    "left button press",
			evtType: PointerButtonPress,
			x:       150.0,
			y:       250.0,
			button:  PointerButtonLeft,
		},
		{
			name:    "right button press",
			evtType: PointerButtonPress,
			x:       150.0,
			y:       250.0,
			button:  PointerButtonRight,
		},
		{
			name:    "middle button press",
			evtType: PointerButtonPress,
			x:       150.0,
			y:       250.0,
			button:  PointerButtonMiddle,
		},
		{
			name:    "button release",
			evtType: PointerButtonRelease,
			x:       155.0,
			y:       255.0,
			button:  PointerButtonLeft,
		},
		{
			name:    "vertical scroll",
			evtType: PointerScroll,
			x:       200.0,
			y:       300.0,
			axis:    ScrollAxisVertical,
			value:   10.0,
		},
		{
			name:    "horizontal scroll",
			evtType: PointerScroll,
			x:       200.0,
			y:       300.0,
			axis:    ScrollAxisHorizontal,
			value:   -5.0,
		},
		{
			name:    "pointer enter",
			evtType: PointerEnter,
			x:       0.0,
			y:       0.0,
		},
		{
			name:    "pointer leave",
			evtType: PointerLeave,
			x:       800.0,
			y:       600.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evt := NewPointerEvent(tt.evtType, tt.x, tt.y, tt.button, tt.axis, tt.value)
			
			if evt == nil {
				t.Fatal("NewPointerEvent returned nil")
			}
			
			if evt.Type() != EventTypePointer {
				t.Errorf("Expected Type() = EventTypePointer, got %v", evt.Type())
			}
			
			if evt.EventType() != tt.evtType {
				t.Errorf("Expected EventType() = %v, got %v", tt.evtType, evt.EventType())
			}
			
			if evt.X() != tt.x {
				t.Errorf("Expected X() = %f, got %f", tt.x, evt.X())
			}
			
			if evt.Y() != tt.y {
				t.Errorf("Expected Y() = %f, got %f", tt.y, evt.Y())
			}
			
			if evt.Button() != tt.button {
				t.Errorf("Expected Button() = %v, got %v", tt.button, evt.Button())
			}
			
			if evt.Axis() != tt.axis {
				t.Errorf("Expected Axis() = %v, got %v", tt.axis, evt.Axis())
			}
			
			if evt.Value() != tt.value {
				t.Errorf("Expected Value() = %f, got %f", tt.value, evt.Value())
			}
			
			// Check timestamp is recent
			if time.Since(evt.Timestamp()) > time.Second {
				t.Error("Timestamp is too old")
			}
		})
	}
}

func TestKeyEventCreation(t *testing.T) {
	tests := []struct {
		name     string
		evtType  KeyEventType
		key      Key
		modifier Modifier
		char     rune
	}{
		{
			name:     "key press backspace",
			evtType:  KeyPress,
			key:      KeyBackspace,
			modifier: 0,
			char:     0,
		},
		{
			name:     "key release backspace",
			evtType:  KeyRelease,
			key:      KeyBackspace,
			modifier: 0,
		},
		{
			name:     "shift+key",
			evtType:  KeyPress,
			key:      KeySpace,
			modifier: ModShift,
			char:     ' ',
		},
		{
			name:     "ctrl+key",
			evtType:  KeyPress,
			key:      KeyReturn,
			modifier: ModControl,
			char:     '\n',
		},
		{
			name:     "alt+key",
			evtType:  KeyPress,
			key:      KeyTab,
			modifier: ModAlt,
			char:     '\t',
		},
		{
			name:     "ctrl+shift+key",
			evtType:  KeyPress,
			key:      KeySpace,
			modifier: ModControl | ModShift,
			char:     ' ',
		},
		{
			name:     "enter key",
			evtType:  KeyPress,
			key:      KeyReturn,
			modifier: 0,
		},
		{
			name:     "backspace key",
			evtType:  KeyPress,
			key:      KeyBackspace,
			modifier: 0,
		},
		{
			name:     "delete key",
			evtType:  KeyPress,
			key:      KeyDelete,
			modifier: 0,
		},
		{
			name:     "arrow left",
			evtType:  KeyPress,
			key:      KeyLeft,
			modifier: 0,
		},
		{
			name:     "arrow right",
			evtType:  KeyPress,
			key:      KeyRight,
			modifier: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evt := NewKeyEvent(tt.evtType, tt.key, tt.modifier, tt.char)
			
			if evt == nil {
				t.Fatal("NewKeyEvent returned nil")
			}
			
			if evt.Type() != EventTypeKey {
				t.Errorf("Expected Type() = EventTypeKey, got %v", evt.Type())
			}
			
			if evt.EventType() != tt.evtType {
				t.Errorf("Expected EventType() = %v, got %v", tt.evtType, evt.EventType())
			}
			
			if evt.Key() != tt.key {
				t.Errorf("Expected Key() = %v, got %v", tt.key, evt.Key())
			}
			
			if evt.Modifiers() != tt.modifier {
				t.Errorf("Expected Modifier() = %v, got %v", tt.modifier, evt.Modifiers())
			}
			
			if evt.Rune() != tt.char {
				t.Errorf("Expected Rune() = %v, got %v", tt.char, evt.Rune())
			}
			
			// Check timestamp is recent
			if time.Since(evt.Timestamp()) > time.Second {
				t.Error("Timestamp is too old")
			}
		})
	}
}

func TestTouchEventCreation(t *testing.T) {
	tests := []struct {
		name    string
		evtType TouchEventType
		id      int32
		x, y    float64
	}{
		{
			name:    "touch down",
			evtType: TouchDown,
			id:      1,
			x:       100.0,
			y:       200.0,
		},
		{
			name:    "touch motion",
			evtType: TouchMotion,
			id:      1,
			x:       105.0,
			y:       205.0,
		},
		{
			name:    "touch up",
			evtType: TouchUp,
			id:      1,
			x:       110.0,
			y:       210.0,
		},
		{
			name:    "multi-touch first finger",
			evtType: TouchDown,
			id:      0,
			x:       50.0,
			y:       50.0,
		},
		{
			name:    "multi-touch second finger",
			evtType: TouchDown,
			id:      1,
			x:       100.0,
			y:       100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evt := NewTouchEvent(tt.evtType, tt.id, tt.x, tt.y)
			
			if evt == nil {
				t.Fatal("NewTouchEvent returned nil")
			}
			
			if evt.Type() != EventTypeTouch {
				t.Errorf("Expected Type() = EventTypeTouch, got %v", evt.Type())
			}
			
			if evt.EventType() != tt.evtType {
				t.Errorf("Expected EventType() = %v, got %v", tt.evtType, evt.EventType())
			}
			
			if evt.ID() != tt.id {
				t.Errorf("Expected ID() = %d, got %d", tt.id, evt.ID())
			}
			
			if evt.X() != tt.x {
				t.Errorf("Expected X() = %f, got %f", tt.x, evt.X())
			}
			
			if evt.Y() != tt.y {
				t.Errorf("Expected Y() = %f, got %f", tt.y, evt.Y())
			}
			
			// Check timestamp is recent
			if time.Since(evt.Timestamp()) > time.Second {
				t.Error("Timestamp is too old")
			}
		})
	}
}

func TestEventConsumption(t *testing.T) {
	evt := NewPointerEvent(PointerMove, 100, 200, 0, 0, 0)
	
	if evt.Consumed() {
		t.Error("Event should not be consumed initially")
	}
	
	evt.Consume()
	
	if !evt.Consumed() {
		t.Error("Event should be consumed after Consume() call")
	}
}

func TestEventTimestamp(t *testing.T) {
	before := time.Now()
	evt := NewPointerEvent(PointerMove, 0, 0, 0, 0, 0)
	after := time.Now()
	
	ts := evt.Timestamp()
	
	if ts.Before(before) || ts.After(after) {
		t.Error("Event timestamp should be between before and after times")
	}
}

func TestModifierCombinations(t *testing.T) {
	tests := []struct {
		name     string
		modifier Modifier
		hasShift bool
		hasCtrl  bool
		hasAlt   bool
	}{
		{
			name:     "no modifiers",
			modifier: 0,
			hasShift: false,
			hasCtrl:  false,
			hasAlt:   false,
		},
		{
			name:     "shift only",
			modifier: ModShift,
			hasShift: true,
			hasCtrl:  false,
			hasAlt:   false,
		},
		{
			name:     "ctrl only",
			modifier: ModControl,
			hasShift: false,
			hasCtrl:  true,
			hasAlt:   false,
		},
		{
			name:     "alt only",
			modifier: ModAlt,
			hasShift: false,
			hasCtrl:  false,
			hasAlt:   true,
		},
		{
			name:     "ctrl+shift",
			modifier: ModControl | ModShift,
			hasShift: true,
			hasCtrl:  true,
			hasAlt:   false,
		},
		{
			name:     "ctrl+alt",
			modifier: ModControl | ModAlt,
			hasShift: false,
			hasCtrl:  true,
			hasAlt:   true,
		},
		{
			name:     "shift+alt",
			modifier: ModShift | ModAlt,
			hasShift: true,
			hasCtrl:  false,
			hasAlt:   true,
		},
		{
			name:     "ctrl+shift+alt",
			modifier: ModControl | ModShift | ModAlt,
			hasShift: true,
			hasCtrl:  true,
			hasAlt:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evt := NewKeyEvent(KeyPress, KeySpace, tt.modifier, ' ')
			
			hasShift := (evt.Modifiers() & ModShift) != 0
			hasCtrl := (evt.Modifiers() & ModControl) != 0
			hasAlt := (evt.Modifiers() & ModAlt) != 0
			
			if hasShift != tt.hasShift {
				t.Errorf("Expected hasShift = %v, got %v", tt.hasShift, hasShift)
			}
			if hasCtrl != tt.hasCtrl {
				t.Errorf("Expected hasCtrl = %v, got %v", tt.hasCtrl, hasCtrl)
			}
			if hasAlt != tt.hasAlt {
				t.Errorf("Expected hasAlt = %v, got %v", tt.hasAlt, hasAlt)
			}
		})
	}
}

func TestPointerButtonConstants(t *testing.T) {
	if PointerButtonLeft != 0x110 {
		t.Errorf("Expected PointerButtonLeft = 0x110, got %x", PointerButtonLeft)
	}
	if PointerButtonRight != 0x111 {
		t.Errorf("Expected PointerButtonRight = 0x111, got %x", PointerButtonRight)
	}
	if PointerButtonMiddle != 0x112 {
		t.Errorf("Expected PointerButtonMiddle = 0x112, got %x", PointerButtonMiddle)
	}
}

func TestScrollAxisConstants(t *testing.T) {
	if ScrollAxisVertical != 0 {
		t.Errorf("Expected ScrollAxisVertical = 0, got %d", ScrollAxisVertical)
	}
	if ScrollAxisHorizontal != 1 {
		t.Errorf("Expected ScrollAxisHorizontal = 1, got %d", ScrollAxisHorizontal)
	}
}
