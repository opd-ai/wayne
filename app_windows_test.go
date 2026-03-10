//go:build windows

package wayne

import (
	"testing"
)

// TestWindowsSmokeTest verifies basic app and window creation on Windows.
// This test ensures the platform-specific initialization works correctly.
func TestWindowsSmokeTest(t *testing.T) {
	cfg := &AppConfig{
		Width:  800,
		Height: 600,
	}

	app := NewAppWithConfig(*cfg)
	if app == nil {
		t.Fatal("NewAppWithConfig returned nil")
	}

	// Create a window to verify platform initialization
	winCfg := &WindowConfig{
		Title:  "Test Window",
		Width:  400,
		Height: 300,
	}

	win, err := app.NewWindow(*winCfg)
	if err != nil {
		t.Fatalf("Failed to create window on Windows: %v", err)
	}

	if win == nil {
		t.Fatal("Window should not be nil")
	}

	// Note: Cannot call app.Run() in a test as it blocks indefinitely
	// Production apps would call Run() in main()
}

// TestWindowsBasicWidgetCreation verifies widgets can be instantiated on Windows.
func TestWindowsBasicWidgetCreation(t *testing.T) {
	cfg := &AppConfig{
		Width:  800,
		Height: 600,
	}

	app := NewAppWithConfig(*cfg)
	if app == nil {
		t.Fatal("NewAppWithConfig returned nil")
	}

	winCfg := &WindowConfig{
		Title:  "Widget Test",
		Width:  400,
		Height: 300,
	}

	win, err := app.NewWindow(*winCfg)
	if err != nil {
		t.Fatalf("Failed to create window: %v", err)
	}

	// Create basic widgets to verify platform compatibility
	// Size uses percentages (0-100) of parent container
	btn := NewButton("Test Button", Size{Width: 50, Height: 10})

	lbl := NewLabel("Test Label", Size{Width: 50, Height: 8})

	input := NewTextInput("", Size{Width: 80, Height: 10})
	input.SetPlaceholder("Test Input")

	panel := NewPanel(Size{Width: 100, Height: 100})
	panel.SetFlowDirection(FlowColumn)
	panel.Add(btn)
	panel.Add(lbl)
	panel.Add(input)

	win.SetRoot(panel)

	// Widgets were created successfully (no panic)
	// Note: Cannot verify rendering without Run()
}
