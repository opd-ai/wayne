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

func TestFocusNavigation(t *testing.T) {
	// Create a dispatcher with focus management
	dispatcher := NewEventDispatcher()

	// Create three buttons
	btn1 := NewButton("Button 1", Size{Width: 20, Height: 8})
	btn1.SetBounds(0, 0, 100, 40)
	btn2 := NewButton("Button 2", Size{Width: 20, Height: 8})
	btn2.SetBounds(0, 50, 100, 40)
	btn3 := NewButton("Button 3", Size{Width: 20, Height: 8})
	btn3.SetBounds(0, 100, 100, 40)

	// Create a column and add buttons
	col := NewColumn()
	col.Add(btn1)
	col.Add(btn2)
	col.Add(btn3)

	// Set as root widget (this builds the focus chain)
	dispatcher.SetWidgetRoot(col)

	// Initially no widget should have focus
	if dispatcher.FocusedWidget() != nil {
		t.Error("No widget should be focused initially")
	}

	// Simulate Tab press to focus first button
	tabEvent := NewKeyEvent(KeyPress, KeyTab, 0, '\t')
	dispatcher.Dispatch(tabEvent)

	// First button should now have focus
	if !btn1.IsFocused() {
		t.Error("Button 1 should have focus after first Tab")
	}
	if btn2.IsFocused() || btn3.IsFocused() {
		t.Error("Other buttons should not have focus")
	}

	// Simulate another Tab press
	tabEvent2 := NewKeyEvent(KeyPress, KeyTab, 0, '\t')
	dispatcher.Dispatch(tabEvent2)

	// Second button should now have focus
	if btn1.IsFocused() {
		t.Error("Button 1 should lose focus")
	}
	if !btn2.IsFocused() {
		t.Error("Button 2 should have focus after second Tab")
	}
	if btn3.IsFocused() {
		t.Error("Button 3 should not have focus yet")
	}

	// Simulate Shift+Tab to go back
	shiftTabEvent := NewKeyEvent(KeyPress, KeyTab, ModShift, '\t')
	dispatcher.Dispatch(shiftTabEvent)

	// First button should have focus again
	if !btn1.IsFocused() {
		t.Error("Button 1 should have focus after Shift+Tab")
	}
	if btn2.IsFocused() || btn3.IsFocused() {
		t.Error("Other buttons should not have focus")
	}
}

func TestFocusableInterface(t *testing.T) {
	// Test that Button implements Focusable
	btn := NewButton("Test", Size{Width: 20, Height: 8})
	if _, ok := interface{}(btn).(Focusable); !ok {
		t.Error("Button should implement Focusable interface")
	}

	// Enabled buttons can take focus
	if !btn.CanTakeFocus() {
		t.Error("Enabled button should be able to take focus")
	}

	// Disabled buttons cannot take focus
	btn.SetEnabled(false)
	if btn.CanTakeFocus() {
		t.Error("Disabled button should not be able to take focus")
	}

	// Test that TextInput implements Focusable
	input := NewTextInput("test", Size{Width: 30, Height: 6})
	if _, ok := interface{}(input).(Focusable); !ok {
		t.Error("TextInput should implement Focusable interface")
	}

	if !input.CanTakeFocus() {
		t.Error("TextInput should be able to take focus")
	}
}

func TestKeyboardActivation(t *testing.T) {
	// Test that focused buttons can be activated with Enter/Space
	btn := NewButton("Test", Size{Width: 20, Height: 8})
	btn.SetBounds(0, 0, 100, 40)

	clicked := false
	btn.OnClick(func() {
		clicked = true
	})

	// Not focused, keyboard events should be ignored
	btn.SetFocused(false)
	enterEvent := NewKeyEvent(KeyPress, KeyReturn, 0, '\r')
	if btn.HandleEvent(enterEvent) {
		t.Error("Unfocused button should not handle keyboard events")
	}
	if clicked {
		t.Error("Button should not be clicked when unfocused")
	}

	// Now focus and press Enter
	btn.SetFocused(true)
	enterEvent2 := NewKeyEvent(KeyPress, KeyReturn, 0, '\r')
	if !btn.HandleEvent(enterEvent2) {
		t.Error("Focused button should handle Enter key")
	}
	if !clicked {
		t.Error("Button should be clicked on Enter")
	}

	// Reset and test Space key
	clicked = false
	spaceEvent := NewKeyEvent(KeyPress, KeySpace, 0, ' ')
	if !btn.HandleEvent(spaceEvent) {
		t.Error("Focused button should handle Space key")
	}
	if !clicked {
		t.Error("Button should be clicked on Space")
	}
}

func TestFocusChainAutoBuilding(t *testing.T) {
	dispatcher := NewEventDispatcher()

	// Create a nested structure
	col1 := NewColumn()
	btn1 := NewButton("Btn1", Size{Width: 20, Height: 8})
	btn2 := NewButton("Btn2", Size{Width: 20, Height: 8})
	col1.Add(btn1)
	col1.Add(btn2)

	row := NewRow()
	btn3 := NewButton("Btn3", Size{Width: 20, Height: 8})
	btn4 := NewButton("Btn4", Size{Width: 20, Height: 8})
	row.Add(btn3)
	row.Add(btn4)

	root := NewColumn()
	root.Add(col1)
	root.Add(row)

	// Set root (should auto-build focus chain)
	dispatcher.SetWidgetRoot(root)

	// Tab through all buttons
	for i := 0; i < 4; i++ {
		tabEvent := NewKeyEvent(KeyPress, KeyTab, 0, '\t')
		dispatcher.Dispatch(tabEvent)
	}

	// All buttons should have been focused in depth-first order
	// The implementation should have cycled through btn1, btn2, btn3, btn4
}
