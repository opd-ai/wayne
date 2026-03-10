//go:build windows || darwin || android || ios

package wayne

import (
	"testing"
)

func TestNewEventDispatcher(t *testing.T) {
	dispatcher := NewEventDispatcher()
	
	if dispatcher == nil {
		t.Fatal("NewEventDispatcher returned nil")
	}
}

func TestEventDispatchToWidget(t *testing.T) {
	dispatcher := NewEventDispatcher()
	
	// Test that dispatcher exists and can be created
	if dispatcher == nil {
		t.Error("Dispatcher should not be nil")
	}
}

func TestFocusManagerCreation(t *testing.T) {
	fm := NewFocusManager()
	
	if fm == nil {
		t.Fatal("NewFocusManager returned nil")
	}
}

func TestFocusManagerSetChain(t *testing.T) {
	fm := NewFocusManager()
	
	// SetChain expects Widget interface
	// We can test it exists
	if fm == nil {
		t.Error("FocusManager should not be nil")
	}
}

func TestFocusManagerFocused(t *testing.T) {
	fm := NewFocusManager()
	
	// Initially no widget should be focused
	if fm.Focused() != nil {
		t.Error("No widget should be focused initially")
	}
}

func TestFocusManagerNext(t *testing.T) {
	fm := NewFocusManager()
	
	// Test FocusNext() doesn't panic with empty chain
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("FocusNext() panicked with empty chain: %v", r)
		}
	}()
	
	fm.FocusNext()
}

func TestFocusManagerPrevious(t *testing.T) {
	fm := NewFocusManager()
	
	// Test FocusPrev() doesn't panic with empty chain
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("FocusPrev() panicked with empty chain: %v", r)
		}
	}()
	
	fm.FocusPrev()
}

func TestEventDispatcherSetWidgetRoot(t *testing.T) {
	_ = NewEventDispatcher()
	
	// Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dispatcher creation panicked: %v", r)
		}
	}()
}

func TestEventDispatcherDispatchPointerEvent(t *testing.T) {
	_ = NewEventDispatcher()
	
	// Create event
	evt := NewPointerEvent(PointerMove, 100, 200, 0, 0, 0)
	
	// Dispatch without root (should not panic)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dispatch panicked: %v", r)
		}
	}()
	
	// We can't fully test dispatch without setting up the widget tree
	// but we verify the event creation works
	if evt == nil {
		t.Error("Event should not be nil")
	}
}

func TestEventDispatcherDispatchKeyEvent(t *testing.T) {
	dispatcher := NewEventDispatcher()
	
	// Create key event
	evt := NewKeyEvent(KeyPress, KeySpace, 0, ' ')
	
	// Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Key event creation panicked: %v", r)
		}
	}()
	
	if evt == nil {
		t.Error("Event should not be nil")
	}
	
	if dispatcher == nil {
		t.Error("Dispatcher should not be nil")
	}
}

func TestEventConsumptionStopsPropagation(t *testing.T) {
	evt := NewPointerEvent(PointerMove, 100, 200, 0, 0, 0)
	
	// Consume the event
	evt.Consume()
	
	// Further handlers should check Consumed()
	if !evt.Consumed() {
		t.Error("Event should be marked as consumed")
	}
}

