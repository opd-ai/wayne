//go:build windows

package wayne

import "testing"

// TestWindowsSmokeTest verifies basic app and window creation on Windows.
func TestWindowsSmokeTest(t *testing.T) {
	_, _ = createTestAppWithWindow(t, "Test Window")
	// Note: Cannot call app.Run() in a test as it blocks indefinitely
}

// TestWindowsBasicWidgetCreation verifies widgets can be instantiated on Windows.
func TestWindowsBasicWidgetCreation(t *testing.T) {
	_, win := createTestAppWithWindow(t, "Widget Test")
	panel := createTestWidgetPanel(t)
	win.SetRoot(panel)
}
