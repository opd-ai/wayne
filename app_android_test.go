//go:build android

package wayne

import "testing"

// TestAndroidSmokeTest verifies basic app and window creation on Android.
//
// Note: Running tests on Android requires gomobile and an emulator or physical device.
// See: https://pkg.go.dev/golang.org/x/mobile/cmd/gomobile
func TestAndroidSmokeTest(t *testing.T) {
	_, _ = createTestAppWithWindow(t, "Test Window")
	// Note: Cannot call app.Run() in a test as it blocks indefinitely
}

// TestAndroidBasicWidgetCreation verifies widgets can be instantiated on Android.
func TestAndroidBasicWidgetCreation(t *testing.T) {
	_, win := createTestAppWithWindow(t, "Widget Test")
	panel := createTestWidgetPanel(t)
	win.SetRoot(panel)
}
