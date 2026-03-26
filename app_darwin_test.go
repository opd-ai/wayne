//go:build darwin

package wayne

import "testing"

// TestMacOSSmokeTest verifies basic app and window creation on macOS.
func TestMacOSSmokeTest(t *testing.T) {
	_, _ = createTestAppWithWindow(t, "Test Window")
	// Note: Cannot call app.Run() in a test as it blocks indefinitely
}

// TestMacOSBasicWidgetCreation verifies widgets can be instantiated on macOS.
func TestMacOSBasicWidgetCreation(t *testing.T) {
	_, win := createTestAppWithWindow(t, "Widget Test")
	panel := createTestWidgetPanel(t)
	win.SetRoot(panel)
}
