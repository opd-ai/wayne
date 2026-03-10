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
		Background:   RGB(30, 30, 46),
		Foreground:   RGB(205, 214, 244),
		Accent:       RGB(137, 180, 250),
		Border:       RGB(88, 91, 112),
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
		Background:   RGB(245, 245, 245),
		Foreground:   RGB(30, 30, 30),
		Accent:       RGB(74, 144, 226),
		Border:       RGB(200, 200, 200),
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
		Background:   RGB(0, 0, 0),
		Foreground:   RGB(255, 255, 255),
		Accent:       RGB(255, 255, 0),
		Border:       RGB(255, 255, 255),
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
	if s.Background != nil {
		result.Background = *s.Background
	}
	if s.Foreground != nil {
		result.Foreground = *s.Foreground
	}
	if s.Accent != nil {
		result.Accent = *s.Accent
	}
	if s.Border != nil {
		result.Border = *s.Border
	}
	if s.FontSize != nil {
		result.FontSize = *s.FontSize
	}
	if s.Padding != nil {
		result.Padding = *s.Padding
	}
	if s.Gap != nil {
		result.Gap = *s.Gap
	}
	if s.BorderWidth != nil {
		result.BorderWidth = *s.BorderWidth
	}
	if s.BorderRadius != nil {
		result.BorderRadius = *s.BorderRadius
	}
	return result
}
