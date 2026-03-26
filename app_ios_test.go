//go:build ios

package wayne

import "testing"

// TestIOSSmokeTest verifies basic app and window creation on iOS.
//
// Note: Running tests on iOS requires gomobile and an iOS simulator or physical device.
// See: https://pkg.go.dev/golang.org/x/mobile/cmd/gomobile
func TestIOSSmokeTest(t *testing.T) {
	_, _ = createTestAppWithWindow(t, "Test Window")
	// Note: Cannot call app.Run() in a test as it blocks indefinitely
}

// TestIOSBasicWidgetCreation verifies widgets can be instantiated on iOS.
func TestIOSBasicWidgetCreation(t *testing.T) {
	_, win := createTestAppWithWindow(t, "Widget Test")
	panel := createTestWidgetPanel(t)
	win.SetRoot(panel)
}
