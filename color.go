//go:build windows || darwin || android || ios

package wayne

// Color represents an RGBA color with 8-bit channels.
//
// Colors are specified in sRGB color space with separate red, green, blue,
// and alpha components. Alpha of 255 is fully opaque, 0 is fully transparent.
//
// Note: Fields use int32 for gomobile compatibility (uint8 generates invalid
// Objective-C code). Values should be in range 0-255.
//
// Example:
//
//	red := wayne.RGB(255, 0, 0)
//	transparentBlue := wayne.RGBA(0, 0, 255, 128)
type Color struct {
	R, G, B, A int32
}

// RGB creates an opaque color from red, green, and blue components.
// Values should be in range 0-255.
func RGB(r, g, b int32) Color {
	return Color{R: r, G: g, B: b, A: 255}
}

// RGBA creates a color from red, green, blue, and alpha components.
// Values should be in range 0-255.
func RGBA(r, g, b, a int32) Color {
	return Color{R: r, G: g, B: b, A: a}
}

// WithAlpha returns a new color with the specified alpha value.
// Alpha should be in range 0-255.
func (c Color) WithAlpha(a int32) Color {
	return Color{R: c.R, G: c.G, B: c.B, A: a}
}

// Common color constants for convenience.
var (
	// Transparent is fully transparent black.
	Transparent = RGBA(0, 0, 0, 0)

	// Black is opaque black.
	Black = RGB(0, 0, 0)

	// White is opaque white.
	White = RGB(255, 255, 255)

	// Red is opaque red.
	Red = RGB(255, 0, 0)

	// Green is opaque green.
	Green = RGB(0, 255, 0)

	// Blue is opaque blue.
	Blue = RGB(0, 0, 255)

	// Gray is opaque medium gray.
	Gray = RGB(128, 128, 128)

	// LightGray is opaque light gray.
	LightGray = RGB(192, 192, 192)

	// DarkGray is opaque dark gray.
	DarkGray = RGB(64, 64, 64)

	// Yellow is opaque yellow.
	Yellow = RGB(255, 255, 0)
)

// Theme color palette constants for dark theme (Catppuccin-inspired).
var (
	// DarkBase is the primary background color for dark themes.
	DarkBase = RGB(30, 30, 46)

	// DarkText is the primary text color for dark themes.
	DarkText = RGB(205, 214, 244)

	// DarkAccent is the accent/highlight color for dark themes.
	DarkAccent = RGB(137, 180, 250)

	// DarkBorder is the border color for dark themes.
	DarkBorder = RGB(88, 91, 112)
)

// Theme color palette constants for light theme.
var (
	// LightBase is the primary background color for light themes.
	LightBase = RGB(245, 245, 245)

	// LightText is the primary text color for light themes.
	LightText = RGB(30, 30, 30)

	// LightAccent is the accent/highlight color for light themes.
	LightAccent = RGB(74, 144, 226)

	// LightBorder is the border color for light themes.
	LightBorder = RGB(200, 200, 200)
)
