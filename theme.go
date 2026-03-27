//go:build windows || darwin || android || ios

package wayne

const (
	// DefaultPadding is the default inner padding in pixels.
	DefaultPadding = 8

	// DefaultGap is the default gap between sibling widgets in pixels.
	DefaultGap = 6

	// DefaultBorderRadius is the default border radius for rounded corners in pixels.
	DefaultBorderRadius = 4
)

// Theme defines the global visual appearance of the application.
//
// A Theme specifies colors, fonts, spacing, and scale applied application-wide
// or per-widget. Widgets inherit theme from their parent container unless
// overridden with a StyleOverride.
//
// Example:
//
//	app.SetTheme(wayne.DefaultDark())
//	app.SetTheme(wayne.HighContrast())
type Theme struct {
	// Background is the primary background color.
	Background Color

	// Foreground is the primary text/foreground color.
	Foreground Color

	// Accent is the accent/highlight color for interactive elements.
	Accent Color

	// Border is the default border color.
	Border Color

	// FontSize is the base font size in pixels.
	FontSize float64

	// Padding is the default inner padding in pixels.
	Padding int

	// Gap is the default gap between sibling widgets in pixels.
	Gap int

	// BorderWidth is the default border width in pixels.
	BorderWidth int

	// BorderRadius is the default border radius for rounded corners in pixels.
	BorderRadius int

	// Scale is the HiDPI scale factor (1.0 = standard, 2.0 = retina).
	Scale float64
}

// DefaultDark returns the built-in dark theme.
func DefaultDark() Theme {
	return Theme{
		Background:   DarkBase,
		Foreground:   DarkText,
		Accent:       DarkAccent,
		Border:       DarkBorder,
		FontSize:     14.0,
		Padding:      DefaultPadding,
		Gap:          DefaultGap,
		BorderWidth:  1,
		BorderRadius: DefaultBorderRadius,
		Scale:        1.0,
	}
}

// DefaultLight returns the built-in light theme.
func DefaultLight() Theme {
	return Theme{
		Background:   LightBase,
		Foreground:   LightText,
		Accent:       LightAccent,
		Border:       LightBorder,
		FontSize:     14.0,
		Padding:      DefaultPadding,
		Gap:          DefaultGap,
		BorderWidth:  1,
		BorderRadius: DefaultBorderRadius,
		Scale:        1.0,
	}
}

// HighContrast returns a high-contrast theme for accessibility.
func HighContrast() Theme {
	return Theme{
		Background:   Black,
		Foreground:   White,
		Accent:       Yellow,
		Border:       White,
		FontSize:     16.0,
		Padding:      10,
		Gap:          8,
		BorderWidth:  2,
		BorderRadius: 0,
		Scale:        1.0,
	}
}

// StyleOverride provides per-widget visual customization.
//
// Any field left as nil will inherit from the parent container's theme.
//
// Example:
//
//	bg := wayne.RGB(40, 40, 60)
//	panel.SetStyle(wayne.StyleOverride{Background: &bg})
type StyleOverride struct {
	// Background overrides the background color if non-nil.
	Background *Color

	// Foreground overrides the foreground color if non-nil.
	Foreground *Color

	// Accent overrides the accent color if non-nil.
	Accent *Color

	// Border overrides the border color if non-nil.
	Border *Color

	// FontSize overrides the font size if non-nil.
	FontSize *float64

	// Padding overrides the padding if non-nil.
	Padding *int

	// Gap overrides the gap if non-nil.
	Gap *int

	// BorderWidth overrides the border width if non-nil.
	BorderWidth *int

	// BorderRadius overrides the border radius if non-nil.
	BorderRadius *int
}

// applyToTheme applies the style override to a theme, returning a new merged theme.
func (s StyleOverride) applyToTheme(base Theme) Theme {
	result := base
	applyColorOverride(&result.Background, s.Background)
	applyColorOverride(&result.Foreground, s.Foreground)
	applyColorOverride(&result.Accent, s.Accent)
	applyColorOverride(&result.Border, s.Border)
	applyFloat64Override(&result.FontSize, s.FontSize)
	applyIntOverride(&result.Padding, s.Padding)
	applyIntOverride(&result.Gap, s.Gap)
	applyIntOverride(&result.BorderWidth, s.BorderWidth)
	applyIntOverride(&result.BorderRadius, s.BorderRadius)
	return result
}

// applyColorOverride applies a color override if the source is non-nil.
func applyColorOverride(target *Color, source *Color) {
	if source != nil {
		*target = *source
	}
}

// applyFloat64Override applies a float64 override if the source is non-nil.
func applyFloat64Override(target *float64, source *float64) {
	if source != nil {
		*target = *source
	}
}

// applyIntOverride applies an int override if the source is non-nil.
func applyIntOverride(target *int, source *int) {
	if source != nil {
		*target = *source
	}
}
