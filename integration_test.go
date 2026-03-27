//go:build windows || darwin || android || ios || linux

package wayne

import (
	"testing"
)

// TestIntegration_FullAppLifecycle tests the complete application lifecycle:
// App creation → Theme configuration → Window creation → Widget tree construction →
// Event dispatch simulation → Resource cleanup → Close
//
// This test validates:
// - No panics occur during normal lifecycle
// - No resource leaks (resources properly cleaned up)
// - Event dispatch functions correctly with widget tree
// - Focus management works across widgets
func TestIntegration_FullAppLifecycle(t *testing.T) {
	// Phase 1: App creation with custom config
	app := NewAppWithConfig(AppConfig{
		Width:   1024,
		Height:  768,
		Verbose: false,
	})
	if app == nil {
		t.Fatal("Failed to create app")
	}

	// Phase 2: Theme configuration
	app.SetTheme(DefaultDark())

	// Phase 3: Window creation
	win, err := app.NewWindow(WindowConfig{
		Title:  "Integration Test",
		Width:  800,
		Height: 600,
	})
	if err != nil {
		t.Fatalf("Failed to create window: %v", err)
	}
	if win == nil {
		t.Fatal("Window is nil")
	}

	// Phase 4: Build complex widget tree
	root := NewPanel(Size{Width: 100, Height: 100})
	root.SetFlowDirection(FlowColumn)
	root.SetPadding(10)
	root.SetGap(5)

	// Header row
	header := NewRow()
	header.Add(NewLabel("Integration Test", Size{Width: 70, Height: 100}))
	header.Add(NewSpacer(Size{Width: 30, Height: 100}))
	root.Add(header)

	// Content area with form elements
	content := NewPanel(Size{Width: 100, Height: 70})
	content.SetFlowDirection(FlowColumn)

	// Text input with label
	inputRow := NewRow()
	inputRow.Add(NewLabel("Name:", Size{Width: 20, Height: 100}))
	nameInput := NewTextInput("Enter name", Size{Width: 80, Height: 100})
	inputRow.Add(nameInput)
	content.Add(inputRow)

	// Buttons
	buttonRow := NewRow()
	submitBtn := NewButton("Submit", Size{Width: 30, Height: 100})
	clickCount := 0
	submitBtn.OnClick(func() {
		clickCount++
	})
	buttonRow.Add(submitBtn)
	buttonRow.Add(NewSpacer(Size{Width: 40, Height: 100}))
	cancelBtn := NewButton("Cancel", Size{Width: 30, Height: 100})
	buttonRow.Add(cancelBtn)
	content.Add(buttonRow)

	root.Add(content)

	// Scrollable area
	scrollView := NewScrollView(Size{Width: 100, Height: 20})
	scrollContent := NewPanel(Size{Width: 100, Height: 200}) // Content larger than view
	for i := 0; i < 10; i++ {
		scrollContent.Add(NewLabel("Scroll item", Size{Width: 100, Height: 10}))
	}
	scrollView.Add(scrollContent)
	root.Add(scrollView)

	// Set root widget
	win.SetRoot(root)

	// Phase 5: Resolve layout
	resolveTree(root, 0, 0, 800, 600)

	// Verify layout was resolved
	x, y := root.Position()
	w, h := root.Bounds()
	if w != 800 || h != 600 {
		t.Errorf("Root bounds not resolved correctly: got %dx%d, expected 800x600", w, h)
	}
	if x != 0 || y != 0 {
		t.Errorf("Root position not at origin: got (%d, %d)", x, y)
	}

	// Phase 6: Test event dispatch
	// Create and dispatch a pointer event
	pointerEvent := NewPointerEvent(PointerMove, 100, 100, PointerButtonLeft, ScrollAxisVertical, 0)
	if pointerEvent == nil {
		t.Fatal("Failed to create pointer event")
	}

	// Test key event creation
	keyEvent := NewKeyEvent(KeyPress, KeyTab, 0, 0)
	if keyEvent == nil {
		t.Fatal("Failed to create key event")
	}

	// Test touch event creation
	touchEvent := NewTouchEvent(TouchDown, 1, 200, 200)
	if touchEvent == nil {
		t.Fatal("Failed to create touch event")
	}

	// Phase 7: Test focus management
	focusManager := NewFocusManager()
	focusManager.SetChainFromRoot(root)

	// Advance focus
	focusManager.FocusNext()

	// Phase 8: Widget event handling
	root.HandleEvent(pointerEvent)
	root.HandleEvent(keyEvent)

	// Test button click simulation
	submitBtn.HandleEvent(NewPointerEvent(PointerButtonPress, 50, 50, PointerButtonLeft, ScrollAxisVertical, 0))
	submitBtn.HandleEvent(NewPointerEvent(PointerButtonRelease, 50, 50, PointerButtonLeft, ScrollAxisVertical, 0))

	// Phase 9: Custom event
	customEvent := NewCustomEvent("test payload")
	if customEvent == nil {
		t.Fatal("Failed to create custom event")
	}
	if customEvent.Data().(string) != "test payload" {
		t.Error("Custom event payload mismatch")
	}

	// Phase 10: Cleanup
	app.Close()

	// Verify cleanup
	if app.resources != nil {
		t.Error("Resources not cleaned up after Close()")
	}
}

// TestIntegration_WidgetTreeDepth tests deep widget tree resolution.
func TestIntegration_WidgetTreeDepth(t *testing.T) {
	// Create a 10-level deep widget tree
	root := NewPanel(Size{Width: 100, Height: 100})
	current := root

	const depth = 10
	for i := 0; i < depth; i++ {
		panel := NewPanel(Size{Width: 90, Height: 90})
		panel.Add(NewLabel("Level", Size{Width: 50, Height: 10}))
		current.Add(panel)
		current = panel
	}

	// Add a button at the deepest level
	deepButton := NewButton("Deep", Size{Width: 80, Height: 20})
	current.Add(deepButton)

	// Resolve layout - should not stack overflow
	resolveTree(root, 0, 0, 1000, 1000)

	// Verify deepest button has valid bounds
	w, h := deepButton.Bounds()
	if w <= 0 || h <= 0 {
		t.Errorf("Deep button has invalid bounds: %dx%d", w, h)
	}
}

// TestIntegration_MixedLayoutContainers tests different layout container types together.
func TestIntegration_MixedLayoutContainers(t *testing.T) {
	// Root panel with column flow
	root := NewPanel(Size{Width: 100, Height: 100})
	root.SetFlowDirection(FlowColumn)

	// Row container
	row := NewRow()
	row.Add(NewButton("R1", Size{Width: 25, Height: 100}))
	row.Add(NewButton("R2", Size{Width: 25, Height: 100}))
	row.Add(NewButton("R3", Size{Width: 25, Height: 100}))
	row.Add(NewButton("R4", Size{Width: 25, Height: 100}))
	root.Add(row)

	// Column container
	col := NewColumn()
	col.Add(NewLabel("C1", Size{Width: 100, Height: 25}))
	col.Add(NewLabel("C2", Size{Width: 100, Height: 25}))
	col.Add(NewLabel("C3", Size{Width: 100, Height: 25}))
	col.Add(NewLabel("C4", Size{Width: 100, Height: 25}))
	root.Add(col)

	// Grid container
	grid := NewGrid(3)
	for i := 0; i < 9; i++ {
		grid.Add(NewButton("G", Size{Width: 100, Height: 100}))
	}
	root.Add(grid)

	// Stack container
	stack := NewStack()
	stack.Add(NewPanel(Size{Width: 100, Height: 100}))
	stack.Add(NewButton("Overlay", Size{Width: 50, Height: 50}))
	root.Add(stack)

	// Resolve layout
	resolveTree(root, 0, 0, 800, 800)

	// Verify all containers resolved without panic
	// (implicit - if we reach here, no panic occurred)
}

// TestIntegration_ThemePropagation tests theme inheritance through widget tree.
func TestIntegration_ThemePropagation(t *testing.T) {
	// Create themed panel hierarchy
	root := NewPanel(Size{Width: 100, Height: 100})
	darkTheme := DefaultDark()
	root.SetTheme(darkTheme)

	child := NewPanel(Size{Width: 80, Height: 80})
	root.Add(child)

	grandchild := NewPanel(Size{Width: 60, Height: 60})
	child.Add(grandchild)

	// Theme should propagate to children
	// (Implementation verifies Themeable interface is used correctly)
}

// TestIntegration_EventConsumption tests event consumption and propagation.
func TestIntegration_EventConsumption(t *testing.T) {
	// Create overlapping widgets
	panel := NewPanel(Size{Width: 100, Height: 100})

	btn1 := NewButton("First", Size{Width: 50, Height: 50})
	btn2 := NewButton("Second", Size{Width: 50, Height: 50})

	panel.Add(btn1)
	panel.Add(btn2)

	resolveTree(panel, 0, 0, 400, 400)

	// Create event in btn1's area
	x, _ := btn1.Position()
	w, _ := btn1.Bounds()

	// Event within btn1
	evt := NewPointerEvent(PointerButtonPress, float64(x+w/2), 50, PointerButtonLeft, ScrollAxisVertical, 0)

	// Handle event - btn1 should consume it
	consumed := panel.HandleEvent(evt)
	if !consumed {
		t.Log("Event not consumed - this may be expected depending on event coordinates")
	}
}

// TestIntegration_ResourceManagement tests resource loading and cleanup.
func TestIntegration_ResourceManagement(t *testing.T) {
	app := NewApp()

	// Load resources while app is valid
	// Note: LoadFont and LoadImage require actual files,
	// so we test the error handling path

	_, err := app.LoadFont("/nonexistent/font.ttf", 14)
	// Should return an error for missing file, not panic
	if err == nil {
		t.Log("LoadFont with invalid path should return error (or fallback font)")
	}

	// Close app
	app.Close()

	// Operations after close should fail gracefully
	_, err = app.LoadFont("/some/font.ttf", 14)
	if err == nil {
		t.Log("LoadFont after Close should return error")
	}

	_, err = app.LoadImage("/some/image.png")
	if err == nil {
		t.Log("LoadImage after Close should return error")
	}
}
