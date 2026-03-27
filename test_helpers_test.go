//go:build windows || darwin || android || ios

package wayne

import "testing"

// createTestApp creates a basic App for testing with default dimensions.
func createTestApp(t *testing.T) *App {
	t.Helper()
	cfg := &AppConfig{
		Width:  800,
		Height: 600,
	}

	app := NewAppWithConfig(*cfg)
	if app == nil {
		t.Fatal("NewAppWithConfig returned nil")
	}
	return app
}

// createTestAppWithWindow creates an App and Window for testing.
func createTestAppWithWindow(t *testing.T, title string) (*App, *Window) {
	t.Helper()
	app := createTestApp(t)

	winCfg := &WindowConfig{
		Title:  title,
		Width:  400,
		Height: 300,
	}

	win, err := app.NewWindow(*winCfg)
	if err != nil {
		t.Fatalf("NewWindow failed: %v", err)
	}
	if win == nil {
		t.Fatal("Window should not be nil")
	}
	return app, win
}

// createTestWidgetPanel creates a Panel with basic widgets for testing.
// Note: Size values are percentages (0-100) of parent, not pixel values.
func createTestWidgetPanel(t *testing.T) *Panel {
	t.Helper()
	btn := NewButton("Test Button", Size{Width: 80, Height: 10})
	lbl := NewLabel("Test Label", Size{Width: 80, Height: 8})
	input := NewTextInput("", Size{Width: 80, Height: 10})
	input.SetPlaceholder("Test Input")

	panel := NewPanel(Size{Width: 100, Height: 100})
	panel.SetFlowDirection(FlowColumn)
	panel.Add(btn)
	panel.Add(lbl)
	panel.Add(input)

	return panel
}
