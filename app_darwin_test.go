//go:build darwin

package wayne

import (
	"testing"
)

// TestMacOSSmokeTest verifies basic app and window creation on macOS.
// This test ensures the platform-specific initialization works correctly.
func TestMacOSSmokeTest(t *testing.T) {
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

	if win == nil {
		t.Fatal("Window should not be nil")
	}

	// Note: Cannot verify rendering without Run()
	// Widgets were created successfully (no panic)
	}

	// Note: Cannot call app.Run() in a test as it blocks indefinitely
	// Production apps would call Run() in main()
}

// TestMacOSBasicWidgetCreation verifies widgets can be instantiated on macOS.
func TestMacOSBasicWidgetCreation(t *testing.T) {
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

	// Create basic widgets to verify platform compatibility
	btn := NewButton("Test Button", Size{Width: 100, Height: 30})

	lbl := NewLabel("Test Label", Size{Width: 100, Height: 20})

	input := NewTextInput("", Size{Width: 200, Height: 30})
	input.SetPlaceholder("Test Input")

	panel := NewPanel(Size{Width: 400, Height: 300})
panel.SetFlowDirection(FlowColumn)
	panel.Add(btn)
	panel.Add(lbl)
	panel.Add(input)

	win.SetRoot(panel)

	// Verify widget tree was set
}
