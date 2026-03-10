//go:build windows || darwin || android || ios

package wayne

import "testing"

// TestColorRGB verifies that RGB creates opaque colors correctly.
func TestColorRGB(t *testing.T) {
	c := RGB(100, 150, 200)

	if c.R != 100 {
		t.Errorf("RGB R incorrect: got %d, want 100", c.R)
	}
	if c.G != 150 {
		t.Errorf("RGB G incorrect: got %d, want 150", c.G)
	}
	if c.B != 200 {
		t.Errorf("RGB B incorrect: got %d, want 200", c.B)
	}
	if c.A != 255 {
		t.Errorf("RGB A incorrect: got %d, want 255 (opaque)", c.A)
	}
}

// TestColorRGBA verifies that RGBA creates colors with custom alpha correctly.
func TestColorRGBA(t *testing.T) {
	c := RGBA(50, 100, 150, 128)

	if c.R != 50 {
		t.Errorf("RGBA R incorrect: got %d, want 50", c.R)
	}
	if c.G != 100 {
		t.Errorf("RGBA G incorrect: got %d, want 100", c.G)
	}
	if c.B != 150 {
		t.Errorf("RGBA B incorrect: got %d, want 150", c.B)
	}
	if c.A != 128 {
		t.Errorf("RGBA A incorrect: got %d, want 128", c.A)
	}
}

// TestColorWithAlpha verifies that WithAlpha creates a new color with modified alpha.
func TestColorWithAlpha(t *testing.T) {
	original := RGB(255, 0, 0)
	modified := original.WithAlpha(128)

	// Original should be unchanged
	if original.A != 255 {
		t.Errorf("Original color modified: got A=%d, want 255", original.A)
	}

	// Modified should have new alpha but same RGB
	if modified.R != 255 || modified.G != 0 || modified.B != 0 {
		t.Errorf("Modified RGB incorrect: got (%d,%d,%d), want (255,0,0)",
			modified.R, modified.G, modified.B)
	}
	if modified.A != 128 {
		t.Errorf("Modified alpha incorrect: got %d, want 128", modified.A)
	}
}

// TestColorConstants verifies that predefined color constants have correct values.
func TestColorConstants(t *testing.T) {
	tests := []struct {
		name       string
		color      Color
		r, g, b, a uint8
	}{
		{"Transparent", Transparent, 0, 0, 0, 0},
		{"Black", Black, 0, 0, 0, 255},
		{"White", White, 255, 255, 255, 255},
		{"Red", Red, 255, 0, 0, 255},
		{"Green", Green, 0, 255, 0, 255},
		{"Blue", Blue, 0, 0, 255, 255},
		{"Gray", Gray, 128, 128, 128, 255},
		{"LightGray", LightGray, 192, 192, 192, 255},
		{"DarkGray", DarkGray, 64, 64, 64, 255},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.color.R != tt.r || tt.color.G != tt.g ||
				tt.color.B != tt.b || tt.color.A != tt.a {
				t.Errorf("%s color incorrect: got RGBA(%d,%d,%d,%d), want RGBA(%d,%d,%d,%d)",
					tt.name, tt.color.R, tt.color.G, tt.color.B, tt.color.A,
					tt.r, tt.g, tt.b, tt.a)
			}
		})
	}
}

// TestDefaultDarkTheme verifies that DefaultDark returns consistent theme values.
func TestDefaultDarkTheme(t *testing.T) {
	theme := DefaultDark()

	if theme.FontSize != 14.0 {
		t.Errorf("DefaultDark FontSize incorrect: got %f, want 14.0", theme.FontSize)
	}
	if theme.Padding != DefaultPadding {
		t.Errorf("DefaultDark Padding incorrect: got %d, want %d", theme.Padding, DefaultPadding)
	}
	if theme.Gap != DefaultGap {
		t.Errorf("DefaultDark Gap incorrect: got %d, want %d", theme.Gap, DefaultGap)
	}
	if theme.BorderWidth != 1 {
		t.Errorf("DefaultDark BorderWidth incorrect: got %d, want 1", theme.BorderWidth)
	}
	if theme.BorderRadius != DefaultBorderRadius {
		t.Errorf("DefaultDark BorderRadius incorrect: got %d, want %d",
			theme.BorderRadius, DefaultBorderRadius)
	}
	if theme.Scale != 1.0 {
		t.Errorf("DefaultDark Scale incorrect: got %f, want 1.0", theme.Scale)
	}

	// Verify color is dark
	if theme.Background.R > 50 || theme.Background.G > 50 || theme.Background.B > 50 {
		t.Errorf("DefaultDark Background not dark: RGB(%d,%d,%d)",
			theme.Background.R, theme.Background.G, theme.Background.B)
	}
}

// TestDefaultLightTheme verifies that DefaultLight returns consistent theme values.
func TestDefaultLightTheme(t *testing.T) {
	theme := DefaultLight()

	if theme.FontSize != 14.0 {
		t.Errorf("DefaultLight FontSize incorrect: got %f, want 14.0", theme.FontSize)
	}
	if theme.Padding != DefaultPadding {
		t.Errorf("DefaultLight Padding incorrect: got %d, want %d", theme.Padding, DefaultPadding)
	}
	if theme.Gap != DefaultGap {
		t.Errorf("DefaultLight Gap incorrect: got %d, want %d", theme.Gap, DefaultGap)
	}
	if theme.BorderWidth != 1 {
		t.Errorf("DefaultLight BorderWidth incorrect: got %d, want 1", theme.BorderWidth)
	}
	if theme.BorderRadius != DefaultBorderRadius {
		t.Errorf("DefaultLight BorderRadius incorrect: got %d, want %d",
			theme.BorderRadius, DefaultBorderRadius)
	}
	if theme.Scale != 1.0 {
		t.Errorf("DefaultLight Scale incorrect: got %f, want 1.0", theme.Scale)
	}

	// Verify color is light
	if theme.Background.R < 200 || theme.Background.G < 200 || theme.Background.B < 200 {
		t.Errorf("DefaultLight Background not light: RGB(%d,%d,%d)",
			theme.Background.R, theme.Background.G, theme.Background.B)
	}
}

// TestHighContrastTheme verifies that HighContrast returns accessible theme values.
func TestHighContrastTheme(t *testing.T) {
	theme := HighContrast()

	if theme.FontSize != 16.0 {
		t.Errorf("HighContrast FontSize incorrect: got %f, want 16.0", theme.FontSize)
	}
	if theme.Padding != 10 {
		t.Errorf("HighContrast Padding incorrect: got %d, want 10", theme.Padding)
	}
	if theme.Gap != 8 {
		t.Errorf("HighContrast Gap incorrect: got %d, want 8", theme.Gap)
	}
	if theme.BorderWidth != 2 {
		t.Errorf("HighContrast BorderWidth incorrect: got %d, want 2", theme.BorderWidth)
	}
	if theme.BorderRadius != 0 {
		t.Errorf("HighContrast BorderRadius incorrect: got %d, want 0 (sharp corners)",
			theme.BorderRadius)
	}
	if theme.Scale != 1.0 {
		t.Errorf("HighContrast Scale incorrect: got %f, want 1.0", theme.Scale)
	}

	// Verify maximum contrast (pure black background)
	if theme.Background != Black {
		t.Errorf("HighContrast Background not black: RGB(%d,%d,%d)",
			theme.Background.R, theme.Background.G, theme.Background.B)
	}
	// Verify maximum contrast (pure white foreground)
	if theme.Foreground != White {
		t.Errorf("HighContrast Foreground not white: RGB(%d,%d,%d)",
			theme.Foreground.R, theme.Foreground.G, theme.Foreground.B)
	}
}

// TestStyleOverrideApplyToTheme_NoOverrides verifies that an empty StyleOverride
// returns the base theme unchanged.
func TestStyleOverrideApplyToTheme_NoOverrides(t *testing.T) {
	base := DefaultDark()
	override := StyleOverride{}

	result := override.applyToTheme(base)

	// All fields should match the base theme
	if result.Background != base.Background {
		t.Error("Background changed when no override specified")
	}
	if result.Foreground != base.Foreground {
		t.Error("Foreground changed when no override specified")
	}
	if result.Accent != base.Accent {
		t.Error("Accent changed when no override specified")
	}
	if result.Border != base.Border {
		t.Error("Border changed when no override specified")
	}
	if result.FontSize != base.FontSize {
		t.Error("FontSize changed when no override specified")
	}
	if result.Padding != base.Padding {
		t.Error("Padding changed when no override specified")
	}
	if result.Gap != base.Gap {
		t.Error("Gap changed when no override specified")
	}
	if result.BorderWidth != base.BorderWidth {
		t.Error("BorderWidth changed when no override specified")
	}
	if result.BorderRadius != base.BorderRadius {
		t.Error("BorderRadius changed when no override specified")
	}
}

// TestStyleOverrideApplyToTheme_BackgroundOnly verifies that overriding only
// background color leaves other fields unchanged.
func TestStyleOverrideApplyToTheme_BackgroundOnly(t *testing.T) {
	base := DefaultDark()
	customBg := RGB(100, 50, 25)
	override := StyleOverride{
		Background: &customBg,
	}

	result := override.applyToTheme(base)

	// Background should be overridden
	if result.Background != customBg {
		t.Errorf("Background not overridden: got RGB(%d,%d,%d), want RGB(100,50,25)",
			result.Background.R, result.Background.G, result.Background.B)
	}

	// All other fields should match base
	if result.Foreground != base.Foreground {
		t.Error("Foreground changed when not overridden")
	}
	if result.Accent != base.Accent {
		t.Error("Accent changed when not overridden")
	}
	if result.Border != base.Border {
		t.Error("Border changed when not overridden")
	}
	if result.FontSize != base.FontSize {
		t.Error("FontSize changed when not overridden")
	}
	if result.Padding != base.Padding {
		t.Error("Padding changed when not overridden")
	}
	if result.Gap != base.Gap {
		t.Error("Gap changed when not overridden")
	}
	if result.BorderWidth != base.BorderWidth {
		t.Error("BorderWidth changed when not overridden")
	}
	if result.BorderRadius != base.BorderRadius {
		t.Error("BorderRadius changed when not overridden")
	}
}

// TestStyleOverrideApplyToTheme_AllColors verifies that all color fields
// can be overridden independently.
func TestStyleOverrideApplyToTheme_AllColors(t *testing.T) {
	base := DefaultDark()
	customBg := RGB(1, 2, 3)
	customFg := RGB(4, 5, 6)
	customAccent := RGB(7, 8, 9)
	customBorder := RGB(10, 11, 12)

	override := StyleOverride{
		Background: &customBg,
		Foreground: &customFg,
		Accent:     &customAccent,
		Border:     &customBorder,
	}

	result := override.applyToTheme(base)

	if result.Background != customBg {
		t.Error("Background not overridden correctly")
	}
	if result.Foreground != customFg {
		t.Error("Foreground not overridden correctly")
	}
	if result.Accent != customAccent {
		t.Error("Accent not overridden correctly")
	}
	if result.Border != customBorder {
		t.Error("Border not overridden correctly")
	}
}

// TestStyleOverrideApplyToTheme_AllMetrics verifies that all numeric fields
// can be overridden independently.
func TestStyleOverrideApplyToTheme_AllMetrics(t *testing.T) {
	base := DefaultDark()
	customFontSize := 20.0
	customPadding := 15
	customGap := 12
	customBorderWidth := 3
	customBorderRadius := 8

	override := StyleOverride{
		FontSize:     &customFontSize,
		Padding:      &customPadding,
		Gap:          &customGap,
		BorderWidth:  &customBorderWidth,
		BorderRadius: &customBorderRadius,
	}

	result := override.applyToTheme(base)

	if result.FontSize != customFontSize {
		t.Errorf("FontSize not overridden: got %f, want %f", result.FontSize, customFontSize)
	}
	if result.Padding != customPadding {
		t.Errorf("Padding not overridden: got %d, want %d", result.Padding, customPadding)
	}
	if result.Gap != customGap {
		t.Errorf("Gap not overridden: got %d, want %d", result.Gap, customGap)
	}
	if result.BorderWidth != customBorderWidth {
		t.Errorf("BorderWidth not overridden: got %d, want %d", result.BorderWidth, customBorderWidth)
	}
	if result.BorderRadius != customBorderRadius {
		t.Errorf("BorderRadius not overridden: got %d, want %d", result.BorderRadius, customBorderRadius)
	}
}

// TestStyleOverrideApplyToTheme_Mixed verifies that partial overrides work correctly.
func TestStyleOverrideApplyToTheme_Mixed(t *testing.T) {
	base := DefaultLight()
	customBg := RGB(255, 128, 64)
	customPadding := 20
	customBorderRadius := 0

	override := StyleOverride{
		Background:   &customBg,
		Padding:      &customPadding,
		BorderRadius: &customBorderRadius,
	}

	result := override.applyToTheme(base)

	// Overridden fields
	if result.Background != customBg {
		t.Error("Background not overridden correctly")
	}
	if result.Padding != customPadding {
		t.Errorf("Padding not overridden: got %d, want %d", result.Padding, customPadding)
	}
	if result.BorderRadius != customBorderRadius {
		t.Errorf("BorderRadius not overridden: got %d, want %d", result.BorderRadius, customBorderRadius)
	}

	// Non-overridden fields should match base
	if result.Foreground != base.Foreground {
		t.Error("Foreground changed when not overridden")
	}
	if result.Accent != base.Accent {
		t.Error("Accent changed when not overridden")
	}
	if result.Border != base.Border {
		t.Error("Border changed when not overridden")
	}
	if result.FontSize != base.FontSize {
		t.Error("FontSize changed when not overridden")
	}
	if result.Gap != base.Gap {
		t.Error("Gap changed when not overridden")
	}
	if result.BorderWidth != base.BorderWidth {
		t.Error("BorderWidth changed when not overridden")
	}
}

// TestStyleOverrideApplyToTheme_ChainedOverrides verifies that multiple
// StyleOverrides can be applied in sequence.
func TestStyleOverrideApplyToTheme_ChainedOverrides(t *testing.T) {
	base := DefaultDark()

	// First override: change background
	customBg := RGB(50, 50, 50)
	override1 := StyleOverride{
		Background: &customBg,
	}
	result1 := override1.applyToTheme(base)

	// Second override: change foreground and padding
	customFg := RGB(220, 220, 220)
	customPadding := 16
	override2 := StyleOverride{
		Foreground: &customFg,
		Padding:    &customPadding,
	}
	result2 := override2.applyToTheme(result1)

	// Final result should have all overrides applied
	if result2.Background != customBg {
		t.Error("First override (Background) lost in chain")
	}
	if result2.Foreground != customFg {
		t.Error("Second override (Foreground) not applied")
	}
	if result2.Padding != customPadding {
		t.Error("Second override (Padding) not applied")
	}

	// Non-overridden fields should still match base
	if result2.Accent != base.Accent {
		t.Error("Non-overridden field (Accent) changed")
	}
	if result2.Gap != base.Gap {
		t.Error("Non-overridden field (Gap) changed")
	}
}

// TestStyleOverrideApplyToTheme_ZeroValues verifies that zero values can be
// set as overrides (not confused with nil).
func TestStyleOverrideApplyToTheme_ZeroValues(t *testing.T) {
	base := DefaultDark()
	zeroPadding := 0
	zeroGap := 0
	zeroRadius := 0
	zeroFontSize := 0.0

	override := StyleOverride{
		Padding:      &zeroPadding,
		Gap:          &zeroGap,
		BorderRadius: &zeroRadius,
		FontSize:     &zeroFontSize,
	}

	result := override.applyToTheme(base)

	// All zero values should be applied (not ignored)
	if result.Padding != 0 {
		t.Errorf("Zero padding not applied: got %d, want 0", result.Padding)
	}
	if result.Gap != 0 {
		t.Errorf("Zero gap not applied: got %d, want 0", result.Gap)
	}
	if result.BorderRadius != 0 {
		t.Errorf("Zero border radius not applied: got %d, want 0", result.BorderRadius)
	}
	if result.FontSize != 0.0 {
		t.Errorf("Zero font size not applied: got %f, want 0.0", result.FontSize)
	}
}

// TestStyleOverrideApplyToTheme_TransparentColor verifies that transparent
// colors can be set as overrides.
func TestStyleOverrideApplyToTheme_TransparentColor(t *testing.T) {
	base := DefaultLight()
	transparentBg := Transparent

	override := StyleOverride{
		Background: &transparentBg,
	}

	result := override.applyToTheme(base)

	if result.Background != Transparent {
		t.Errorf("Transparent background not applied: got RGBA(%d,%d,%d,%d), want RGBA(0,0,0,0)",
			result.Background.R, result.Background.G, result.Background.B, result.Background.A)
	}
}
