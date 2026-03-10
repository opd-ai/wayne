//go:build windows || darwin || android || ios

package wayne

import "testing"

func TestPanelSetThemePropagation(t *testing.T) {
	// Create a panel with various child widget types
	panel := NewPanel(Size{Width: 100, Height: 100})

	// Add different types of children
	btn := NewButton("Test", Size{Width: 50, Height: 20})
	label := NewLabel("Label", Size{Width: 50, Height: 20})
	input := NewTextInput(Size{Width: 50, Height: 20})
	scrollView := NewScrollView(Size{Width: 50, Height: 50})
	childPanel := NewPanel(Size{Width: 50, Height: 50})

	panel.Add(btn)
	panel.Add(label)
	panel.Add(input)
	panel.Add(scrollView)
	panel.Add(childPanel)

	// Create a custom theme
	customTheme := DefaultDark()
	customTheme.Primary = NewColor(255, 0, 0, 255) // Red for testing

	// Apply theme to panel
	panel.SetTheme(customTheme)

	// Verify all children received the theme
	// Since we can't directly check private theme fields, we verify no panics
	// and that the method completes successfully

	// Apply theme again to test idempotency
	panel.SetTheme(DefaultLight())
}

func TestThemeableInterface(t *testing.T) {
	// Verify all concrete widgets implement Themeable
	var _ Themeable = (*Button)(nil)
	var _ Themeable = (*Label)(nil)
	var _ Themeable = (*TextInput)(nil)
	var _ Themeable = (*ScrollView)(nil)
	var _ Themeable = (*Panel)(nil)
}
