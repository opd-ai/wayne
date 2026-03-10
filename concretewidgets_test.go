//go:build windows || darwin || android || ios

package wayne

import (
	"testing"
)

func TestButtonClick(t *testing.T) {
	btn := NewButton("Test", Size{Width: 100, Height: 30})

	if btn == nil {
		t.Fatal("NewButton returned nil")
	}

	clicked := false
	btn.OnClick(func() {
		clicked = true
	})

	// Simulate button press
	pressEvt := NewPointerEvent(PointerButtonPress, 50, 15, PointerButtonLeft, 0, 0)
	btn.HandleEvent(pressEvt)

	if clicked {
		t.Error("Button should not fire onClick on press")
	}

	// Simulate button release
	releaseEvt := NewPointerEvent(PointerButtonRelease, 50, 15, PointerButtonLeft, 0, 0)
	btn.HandleEvent(releaseEvt)

	if !clicked {
		t.Error("Button should fire onClick on release after press")
	}
}

func TestButtonHoverState(t *testing.T) {
	btn := NewButton("Test", Size{Width: 100, Height: 30})

	if btn.hovered {
		t.Error("Button should not be hovered initially")
	}

	enterEvt := NewPointerEvent(PointerEnter, 50, 15, 0, 0, 0)
	btn.HandleEvent(enterEvt)

	if !btn.hovered {
		t.Error("Button should be hovered after PointerEnter")
	}

	leaveEvt := NewPointerEvent(PointerLeave, 150, 50, 0, 0, 0)
	btn.HandleEvent(leaveEvt)

	if btn.hovered {
		t.Error("Button should not be hovered after PointerLeave")
	}
}

func TestButtonDisabled(t *testing.T) {
	btn := NewButton("Test", Size{Width: 100, Height: 30})
	clicked := false
	btn.OnClick(func() {
		clicked = true
	})

	btn.SetEnabled(false)

	pressEvt := NewPointerEvent(PointerButtonPress, 50, 15, PointerButtonLeft, 0, 0)
	btn.HandleEvent(pressEvt)

	releaseEvt := NewPointerEvent(PointerButtonRelease, 50, 15, PointerButtonLeft, 0, 0)
	btn.HandleEvent(releaseEvt)

	if clicked {
		t.Error("Disabled button should not fire onClick")
	}
}

func TestButtonSetLabel(t *testing.T) {
	btn := NewButton("Initial", Size{Width: 100, Height: 30})

	if btn.Text() != "Initial" {
		t.Errorf("Expected initial text 'Initial', got '%s'", btn.Text())
	}

	btn.SetLabel("Updated")

	if btn.Text() != "Updated" {
		t.Errorf("Expected updated text 'Updated', got '%s'", btn.Text())
	}
}

func TestButtonSetText(t *testing.T) {
	btn := NewButton("Initial", Size{Width: 100, Height: 30})

	btn.SetText("Changed")

	if btn.Text() != "Changed" {
		t.Errorf("Expected text 'Changed', got '%s'", btn.Text())
	}
}

func TestLabelRendering(t *testing.T) {
	label := NewLabel("Test Label", Size{Width: 200, Height: 40})

	if label == nil {
		t.Fatal("NewLabel returned nil")
	}

	if label.Text() != "Test Label" {
		t.Errorf("Expected text 'Test Label', got '%s'", label.Text())
	}
}

func TestLabelSetText(t *testing.T) {
	label := NewLabel("Initial", Size{Width: 200, Height: 40})

	label.SetText("Updated")

	if label.Text() != "Updated" {
		t.Errorf("Expected text 'Updated', got '%s'", label.Text())
	}
}

func TestLabelSetTheme(t *testing.T) {
	label := NewLabel("Test", Size{Width: 200, Height: 40})
	customTheme := DefaultLight()

	// Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("SetTheme panicked: %v", r)
		}
	}()

	label.SetTheme(customTheme)
}

func TestTextInputKeyboardEntry(t *testing.T) {
	input := NewTextInput("placeholder", Size{Width: 200, Height: 30})

	if input == nil {
		t.Fatal("NewTextInput returned nil")
	}

	if input.Text() != "" {
		t.Errorf("Expected empty text initially, got '%s'", input.Text())
	}

	// Focus the input first
	input.focused = true

	// Simulate typing 'a'
	keyEvt := NewKeyEvent(KeyPress, KeySpace, 0, 'a')
	input.HandleEvent(keyEvt)

	if input.Text() != "a" {
		t.Errorf("Expected text 'a', got '%s'", input.Text())
	}

	// Simulate typing 'b'
	keyEvt = NewKeyEvent(KeyPress, KeySpace, 0, 'b')
	input.HandleEvent(keyEvt)

	if input.Text() != "ab" {
		t.Errorf("Expected text 'ab', got '%s'", input.Text())
	}
}

func TestTextInputBackspace(t *testing.T) {
	input := NewTextInput("", Size{Width: 200, Height: 30})
	input.focused = true
	input.text = "test"
	input.cursorPos = 4

	backspaceEvt := NewKeyEvent(KeyPress, KeyBackspace, 0, 0)
	input.HandleEvent(backspaceEvt)

	if input.Text() != "tes" {
		t.Errorf("Expected text 'tes' after backspace, got '%s'", input.Text())
	}
}

func TestTextInputDelete(t *testing.T) {
	input := NewTextInput("", Size{Width: 200, Height: 30})
	input.focused = true
	input.text = "test"
	input.cursorPos = 0

	deleteEvt := NewKeyEvent(KeyPress, KeyDelete, 0, 0)
	input.HandleEvent(deleteEvt)

	if input.Text() != "est" {
		t.Errorf("Expected text 'est' after delete, got '%s'", input.Text())
	}
}

func TestTextInputSetText(t *testing.T) {
	input := NewTextInput("", Size{Width: 200, Height: 30})

	input.SetText("new value")

	if input.Text() != "new value" {
		t.Errorf("Expected text 'new value', got '%s'", input.Text())
	}
}

func TestTextInputPlaceholder(t *testing.T) {
	input := NewTextInput("Enter text here", Size{Width: 200, Height: 30})

	if input.placeholder != "Enter text here" {
		t.Errorf("Expected placeholder 'Enter text here', got '%s'", input.placeholder)
	}
}

func TestImageWidget(t *testing.T) {
	img := NewImageWidget(Size{Width: 100, Height: 100})

	if img == nil {
		t.Fatal("NewImageWidget returned nil")
	}

	if img.sizeHint().Width != 100 || img.sizeHint().Height != 100 {
		t.Error("ImageWidget size mismatch")
	}
}

func TestSpacer(t *testing.T) {
	spacer := NewSpacer(Size{Width: 50, Height: 20})

	if spacer == nil {
		t.Fatal("NewSpacer returned nil")
	}

	if spacer.sizeHint().Width != 50 || spacer.sizeHint().Height != 20 {
		t.Error("Spacer size mismatch")
	}
}

func TestScrollView(t *testing.T) {
	sv := NewScrollView(Size{Width: 300, Height: 400})

	if sv == nil {
		t.Fatal("NewScrollView returned nil")
	}

	// Just test that it was created
	if sv.sizeHint().Width != 300 || sv.sizeHint().Height != 400 {
		t.Error("ScrollView size mismatch")
	}
}

func TestScrollViewScrolling(t *testing.T) {
	sv := NewScrollView(Size{Width: 300, Height: 400})

	initialScrollY := sv.scrollY

	// Simulate scroll down (positive value)
	scrollEvt := NewPointerEvent(PointerScroll, 150, 200, 0, ScrollAxisVertical, 10.0)
	sv.HandleEvent(scrollEvt)

	if sv.scrollY == initialScrollY {
		t.Error("ScrollView should update scrollY on scroll event")
	}
}
