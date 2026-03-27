//go:build windows || darwin || android || ios

package wayne

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// TestCanvasScale verifies that Canvas.Scale() returns the theme scale factor.
func TestCanvasScale(t *testing.T) {
	tests := []struct {
		name     string
		scale    float64
		expected float64
	}{
		{"default scale", 1.0, 1.0},
		{"retina scale", 2.0, 2.0},
		{"custom scale", 1.5, 1.5},
		{"zero scale defaults to 1.0", 0, 1.0},
		{"negative scale defaults to 1.0", -1.0, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theme := DefaultDark()
			theme.Scale = tt.scale
			img := ebiten.NewImage(100, 100)
			defer img.Deallocate()
			canvas := newEbitenCanvas(img, theme)

			got := canvas.Scale()
			if got != tt.expected {
				t.Errorf("Canvas.Scale() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestScaledInt verifies the drawContext.scaledInt helper function.
func TestScaledInt(t *testing.T) {
	tests := []struct {
		name     string
		scale    float64
		input    int
		expected int
	}{
		{"scale 1.0", 1.0, 10, 10},
		{"scale 2.0", 2.0, 10, 20},
		{"scale 1.5", 1.5, 10, 15},
		{"scale 2.0 with odd value", 2.0, 7, 14},
		{"scale 1.25", 1.25, 8, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &drawContext{scale: tt.scale}
			got := ctx.scaledInt(tt.input)
			if got != tt.expected {
				t.Errorf("scaledInt(%d) with scale %v = %d, want %d", tt.input, tt.scale, got, tt.expected)
			}
		})
	}
}

// TestDefaultThemeScale verifies that default themes have Scale = 1.0.
func TestDefaultThemeScale(t *testing.T) {
	themes := []struct {
		name  string
		theme Theme
	}{
		{"DefaultDark", DefaultDark()},
		{"DefaultLight", DefaultLight()},
		{"HighContrast", HighContrast()},
	}

	for _, tt := range themes {
		t.Run(tt.name, func(t *testing.T) {
			if tt.theme.Scale != 1.0 {
				t.Errorf("%s.Scale = %v, want 1.0", tt.name, tt.theme.Scale)
			}
		})
	}
}

// TestPrepareDrawContextScale verifies that prepareDrawContext captures the scale.
func TestPrepareDrawContextScale(t *testing.T) {
	tests := []struct {
		name  string
		scale float64
	}{
		{"scale 1.0", 1.0},
		{"scale 2.0", 2.0},
		{"scale 0 defaults to 1.0", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theme := DefaultDark()
			theme.Scale = tt.scale
			img := ebiten.NewImage(100, 100)
			defer img.Deallocate()
			canvas := newEbitenCanvas(img, theme)

			widget := &BasePublicWidget{
				x: 0, y: 0, width: 50, height: 50,
			}

			ctx := prepareDrawContext(canvas, widget, nil)
			if ctx == nil {
				t.Fatal("prepareDrawContext returned nil")
			}

			expectedScale := tt.scale
			if expectedScale <= 0 {
				expectedScale = 1.0
			}

			if ctx.scale != expectedScale {
				t.Errorf("drawContext.scale = %v, want %v", ctx.scale, expectedScale)
			}
		})
	}
}

// TestScaleAffectsText verifies that DrawText respects the scale factor.
// This is a basic test that DrawText doesn't panic with various scales.
func TestScaleAffectsText(t *testing.T) {
	scales := []float64{0.5, 1.0, 1.5, 2.0, 3.0}

	for _, scale := range scales {
		theme := DefaultDark()
		theme.Scale = scale
		img := ebiten.NewImage(200, 100)
		defer img.Deallocate()
		canvas := newEbitenCanvas(img, theme)

		// Should not panic
		canvas.DrawText("Test", 10, 10, nil, White)
	}
}

// TestButtonDrawWithScale verifies Button.Draw doesn't panic at various scales.
func TestButtonDrawWithScale(t *testing.T) {
	scales := []float64{0.5, 1.0, 1.5, 2.0, 3.0}

	for _, scale := range scales {
		theme := DefaultDark()
		theme.Scale = scale
		img := ebiten.NewImage(200, 100)
		defer img.Deallocate()
		canvas := newEbitenCanvas(img, theme)

		btn := NewButton("Test", Size{Width: 50, Height: 10})
		btn.SetBounds(10, 10, 80, 30)

		// Should not panic
		btn.Draw(canvas)
	}
}

// TestTextInputDrawWithScale verifies TextInput.Draw doesn't panic at various scales.
func TestTextInputDrawWithScale(t *testing.T) {
	scales := []float64{0.5, 1.0, 1.5, 2.0, 3.0}

	for _, scale := range scales {
		theme := DefaultDark()
		theme.Scale = scale
		img := ebiten.NewImage(200, 100)
		defer img.Deallocate()
		canvas := newEbitenCanvas(img, theme)

		input := NewTextInput("Placeholder", Size{Width: 50, Height: 10})
		input.SetBounds(10, 10, 150, 30)
		input.SetText("Hello")
		input.SetFocused(true)

		// Should not panic
		input.Draw(canvas)
	}
}
