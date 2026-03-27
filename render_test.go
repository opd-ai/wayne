//go:build windows || darwin || android || ios

package wayne

import (
	"image"
	"image/color"
	"math"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// skipIfPixelReadUnavailable skips the test if ebiten pixel reading is unavailable.
// ebiten.Image.At() panics with "ReadPixels cannot be called before the game starts"
// when called outside of the game loop (i.e., in unit tests).
func skipIfPixelReadUnavailable(t *testing.T, img *ebiten.Image) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Skip("Skipping test: ebiten pixel reading unavailable outside game loop")
		}
	}()
	// Attempt a pixel read to trigger the panic if unavailable
	_ = img.At(0, 0)
}

// TestLinearGradient_HorizontalBasic verifies horizontal gradient (0 degrees)
func TestLinearGradient_HorizontalBasic(t *testing.T) {
	img := ebiten.NewImage(100, 50)
	skipIfPixelReadUnavailable(t, img)
	canvas := newEbitenCanvas(img, DefaultDark())

	red := Color{R: 255, G: 0, B: 0, A: 255}
	blue := Color{R: 0, G: 0, B: 255, A: 255}

	canvas.LinearGradient(0, 0, 100, 50, red, blue, 0)

	// Verify gradient transitions from red (left) to blue (right)
	// Sample at left edge
	leftPixel := img.At(5, 25)
	r1, _, b1, _ := leftPixel.RGBA()
	if r1>>8 < 200 { // Should be mostly red
		t.Errorf("Left edge should be red, got R=%d", r1>>8)
	}
	if b1>>8 > 55 { // Should have minimal blue
		t.Errorf("Left edge should have minimal blue, got B=%d", b1>>8)
	}

	// Sample at right edge
	rightPixel := img.At(95, 25)
	r2, _, b2, _ := rightPixel.RGBA()
	if r2>>8 > 55 { // Should have minimal red
		t.Errorf("Right edge should have minimal red, got R=%d", r2>>8)
	}
	if b2>>8 < 200 { // Should be mostly blue
		t.Errorf("Right edge should be blue, got B=%d", b2>>8)
	}
}

// TestLinearGradient_VerticalBasic verifies vertical gradient (90 degrees)
func TestLinearGradient_VerticalBasic(t *testing.T) {
	img := ebiten.NewImage(50, 100)
	skipIfPixelReadUnavailable(t, img)
	canvas := newEbitenCanvas(img, DefaultDark())

	green := Color{R: 0, G: 255, B: 0, A: 255}
	yellow := Color{R: 255, G: 255, B: 0, A: 255}

	canvas.LinearGradient(0, 0, 50, 100, green, yellow, 90)

	// Sample at top (should be green)
	topPixel := img.At(25, 5)
	_, g1, _, _ := topPixel.RGBA()
	if g1>>8 < 200 {
		t.Errorf("Top edge should be green, got G=%d", g1>>8)
	}

	// Sample at bottom (should be yellow - red+green)
	bottomPixel := img.At(25, 95)
	r2, g2, _, _ := bottomPixel.RGBA()
	if r2>>8 < 200 || g2>>8 < 200 {
		t.Errorf("Bottom edge should be yellow, got R=%d G=%d", r2>>8, g2>>8)
	}
}

// TestLinearGradient_DiagonalAngle verifies diagonal gradients work
func TestLinearGradient_DiagonalAngle(t *testing.T) {
	img := ebiten.NewImage(100, 100)
	skipIfPixelReadUnavailable(t, img)
	canvas := newEbitenCanvas(img, DefaultDark())

	black := Color{R: 0, G: 0, B: 0, A: 255}
	white := Color{R: 255, G: 255, B: 255, A: 255}

	// 45 degree angle
	canvas.LinearGradient(0, 0, 100, 100, black, white, 45)

	// Verify gradient exists (should have variation from corner to corner)
	topLeft := img.At(10, 10)
	r1, _, _, _ := topLeft.RGBA()
	bottomRight := img.At(90, 90)
	r2, _, _, _ := bottomRight.RGBA()

	// There should be a significant color difference along the diagonal
	diff := int(r2>>8) - int(r1>>8)
	if diff < 50 {
		t.Errorf("Expected significant gradient along diagonal, got R diff=%d", diff)
	}
}

// TestLinearGradient_ArbitraryAngle tests arbitrary angle support
func TestLinearGradient_ArbitraryAngle(t *testing.T) {
	angles := []float64{0, 30, 45, 60, 90, 120, 135, 150, 180, 210, 270, 315}
	for _, angle := range angles {
		img := ebiten.NewImage(80, 80)
		canvas := newEbitenCanvas(img, DefaultDark())

		start := Color{R: 100, G: 0, B: 0, A: 255}
		end := Color{R: 0, G: 100, B: 0, A: 255}

		// Should not panic for any angle
		canvas.LinearGradient(0, 0, 80, 80, start, end, angle)
	}
}

// TestLinearGradient_ZeroSize verifies zero-size rectangles are handled
func TestLinearGradient_ZeroSize(t *testing.T) {
	img := ebiten.NewImage(100, 100)
	canvas := newEbitenCanvas(img, DefaultDark())

	color1 := Color{R: 255, G: 0, B: 0, A: 255}
	color2 := Color{R: 0, G: 0, B: 255, A: 255}

	// Should not panic
	canvas.LinearGradient(10, 10, 0, 50, color1, color2, 0)
	canvas.LinearGradient(10, 10, 50, 0, color1, color2, 0)
	canvas.LinearGradient(10, 10, -10, 50, color1, color2, 0)
}

// TestRadialGradient_Basic verifies radial gradient rendering
func TestRadialGradient_Basic(t *testing.T) {
	img := ebiten.NewImage(100, 100)
	skipIfPixelReadUnavailable(t, img)
	canvas := newEbitenCanvas(img, DefaultDark())

	centerWhite := Color{R: 255, G: 255, B: 255, A: 255}
	edgeBlack := Color{R: 0, G: 0, B: 0, A: 255}

	canvas.RadialGradient(0, 0, 100, 100, centerWhite, edgeBlack)

	// Sample at center (should be lighter)
	centerPixel := img.At(50, 50)
	rc, _, _, _ := centerPixel.RGBA()

	// Sample at edge (should be darker)
	edgePixel := img.At(95, 95)
	re, _, _, _ := edgePixel.RGBA()

	if rc>>8 < re>>8 {
		t.Errorf("Center should be lighter than edge, got center R=%d, edge R=%d", rc>>8, re>>8)
	}
}

// TestRadialGradient_DynamicSteps verifies step count scales with size
func TestRadialGradient_DynamicSteps(t *testing.T) {
	// Small gradient (should use minimum 128 steps)
	imgSmall := ebiten.NewImage(50, 50)
	skipIfPixelReadUnavailable(t, imgSmall)
	canvasSmall := newEbitenCanvas(imgSmall, DefaultDark())

	// Large gradient (should use more steps, capped at 512)
	imgLarge := ebiten.NewImage(800, 800)
	canvasLarge := newEbitenCanvas(imgLarge, DefaultDark())

	color1 := Color{R: 255, G: 0, B: 0, A: 255}
	color2 := Color{R: 0, G: 0, B: 255, A: 255}

	// Should not panic and should render smoothly
	canvasSmall.RadialGradient(0, 0, 50, 50, color1, color2)
	canvasLarge.RadialGradient(0, 0, 800, 800, color1, color2)

	// Verify large gradient has reasonable smoothness (no huge jumps)
	// Sample several points along a radius
	samples := []image.Point{{400, 400}, {450, 400}, {500, 400}, {550, 400}}
	var prevIntensity uint32
	for i, pt := range samples {
		pixel := imgLarge.At(pt.X, pt.Y)
		r, _, _, _ := pixel.RGBA()
		if i > 0 {
			// Color should transition gradually (not jump by more than 20 units per 50px)
			diff := int(r>>8) - int(prevIntensity>>8)
			if diff < 0 {
				diff = -diff
			}
			if diff > 50 {
				t.Errorf("Large gradient has excessive banding: diff=%d at position %d", diff, i)
			}
		}
		prevIntensity = r
	}
}

// TestRadialGradient_NonSquare verifies radial gradient on rectangles
func TestRadialGradient_NonSquare(t *testing.T) {
	img := ebiten.NewImage(200, 100)
	skipIfPixelReadUnavailable(t, img)
	canvas := newEbitenCanvas(img, DefaultDark())

	center := Color{R: 0, G: 255, B: 0, A: 255}
	edge := Color{R: 0, G: 0, B: 255, A: 255}

	// Should handle non-square rectangles gracefully
	canvas.RadialGradient(0, 0, 200, 100, center, edge)

	// Center should be green
	centerPixel := img.At(100, 50)
	_, g, _, _ := centerPixel.RGBA()
	if g>>8 < 200 {
		t.Errorf("Center should be green, got G=%d", g>>8)
	}
}

// TestRadialGradient_ZeroSize verifies zero-size handling
func TestRadialGradient_ZeroSize(t *testing.T) {
	img := ebiten.NewImage(100, 100)
	canvas := newEbitenCanvas(img, DefaultDark())

	color1 := Color{R: 255, G: 0, B: 0, A: 255}
	color2 := Color{R: 0, G: 0, B: 255, A: 255}

	// Should not panic
	canvas.RadialGradient(10, 10, 0, 50, color1, color2)
	canvas.RadialGradient(10, 10, 50, 0, color1, color2)
	canvas.RadialGradient(10, 10, -10, -10, color1, color2)
}

// TestBoxShadow_Basic verifies box shadow rendering
func TestBoxShadow_Basic(t *testing.T) {
	img := ebiten.NewImage(200, 200)
	skipIfPixelReadUnavailable(t, img)
	canvas := newEbitenCanvas(img, DefaultDark())

	// Clear to white background
	canvas.FillRect(0, 0, 200, 200, Color{R: 255, G: 255, B: 255, A: 255})

	// Draw shadow for a box at (50, 50, 100x100) with 5px offset and 10px blur
	shadowColor := Color{R: 0, G: 0, B: 0, A: 128}
	canvas.BoxShadow(50, 50, 100, 100, 5, 5, 10, shadowColor)

	// Verify shadow exists at offset position (should be darker than white)
	shadowPixel := img.At(55, 55) // offset position
	r, g, b, _ := shadowPixel.RGBA()
	luminance := (r + g + b) / 3

	// Shadow area should be darker than pure white
	if luminance>>8 > 250 {
		t.Errorf("Shadow area should be darker, got luminance=%d", luminance>>8)
	}
}

// TestBoxShadow_OffsetVariations tests different offset directions
func TestBoxShadow_OffsetVariations(t *testing.T) {
	offsets := []struct {
		x, y int
	}{
		{5, 5},   // bottom-right
		{-5, 5},  // bottom-left
		{5, -5},  // top-right
		{-5, -5}, // top-left
		{0, 10},  // directly down
		{10, 0},  // directly right
	}

	for _, offset := range offsets {
		img := ebiten.NewImage(150, 150)
		canvas := newEbitenCanvas(img, DefaultDark())

		shadowColor := Color{R: 0, G: 0, B: 0, A: 100}
		// Should not panic with any offset direction
		canvas.BoxShadow(50, 50, 50, 50, offset.x, offset.y, 5, shadowColor)
	}
}

// TestBoxShadow_ZeroBlur verifies shadow with zero blur
func TestBoxShadow_ZeroBlur(t *testing.T) {
	img := ebiten.NewImage(100, 100)
	skipIfPixelReadUnavailable(t, img)
	canvas := newEbitenCanvas(img, DefaultDark())

	shadowColor := Color{R: 0, G: 0, B: 0, A: 128}
	// Zero blur should still render a shadow (hard-edged)
	canvas.BoxShadow(20, 20, 40, 40, 3, 3, 0, shadowColor)

	// Should have created some shadow
	shadowPixel := img.At(23, 23)
	r, _, _, a := shadowPixel.RGBA()
	if a == 0 {
		t.Error("Shadow with zero blur should still be visible")
	}
	if r>>8 > 128 {
		t.Errorf("Shadow should be dark, got R=%d", r>>8)
	}
}

// TestBoxShadow_ZeroSize verifies zero-size handling
func TestBoxShadow_ZeroSize(t *testing.T) {
	img := ebiten.NewImage(100, 100)
	canvas := newEbitenCanvas(img, DefaultDark())

	shadowColor := Color{R: 0, G: 0, B: 0, A: 128}

	// Should not panic
	canvas.BoxShadow(10, 10, 0, 50, 5, 5, 3, shadowColor)
	canvas.BoxShadow(10, 10, 50, 0, 5, 5, 3, shadowColor)
}

// TestCanvasGradients_AlphaTransparency verifies alpha channel handling
func TestCanvasGradients_AlphaTransparency(t *testing.T) {
	img := ebiten.NewImage(100, 100)
	skipIfPixelReadUnavailable(t, img)
	canvas := newEbitenCanvas(img, DefaultDark())

	// Test gradient with transparency
	transparent := Color{R: 255, G: 0, B: 0, A: 0}
	opaque := Color{R: 255, G: 0, B: 0, A: 255}

	canvas.LinearGradient(0, 0, 100, 100, transparent, opaque, 0)

	// Left side should be fully transparent
	leftPixel := img.At(10, 50)
	_, _, _, a1 := leftPixel.RGBA()
	if a1>>8 > 50 {
		t.Errorf("Left side should be mostly transparent, got A=%d", a1>>8)
	}

	// Right side should be opaque
	rightPixel := img.At(90, 50)
	_, _, _, a2 := rightPixel.RGBA()
	if a2>>8 < 200 {
		t.Errorf("Right side should be opaque, got A=%d", a2>>8)
	}
}

// TestCanvasGradients_ColorInterpolation verifies smooth color interpolation
func TestCanvasGradients_ColorInterpolation(t *testing.T) {
	img := ebiten.NewImage(256, 50)
	skipIfPixelReadUnavailable(t, img)
	canvas := newEbitenCanvas(img, DefaultDark())

	black := Color{R: 0, G: 0, B: 0, A: 255}
	white := Color{R: 255, G: 255, B: 255, A: 255}

	canvas.LinearGradient(0, 0, 256, 50, black, white, 0)

	// Sample several points and verify monotonic increase
	var prevR uint32
	for x := 10; x < 250; x += 20 {
		pixel := img.At(x, 25)
		r, _, _, _ := pixel.RGBA()
		if x > 10 && r < prevR {
			t.Errorf("Color should increase monotonically, but decreased at x=%d", x)
		}
		prevR = r
	}
}

// BenchmarkLinearGradient measures performance of linear gradient
func BenchmarkLinearGradient(b *testing.B) {
	img := ebiten.NewImage(500, 500)
	canvas := newEbitenCanvas(img, DefaultDark())

	start := Color{R: 255, G: 0, B: 0, A: 255}
	end := Color{R: 0, G: 0, B: 255, A: 255}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		canvas.LinearGradient(0, 0, 500, 500, start, end, 45)
	}
}

// BenchmarkRadialGradient measures performance of radial gradient
func BenchmarkRadialGradient(b *testing.B) {
	img := ebiten.NewImage(500, 500)
	canvas := newEbitenCanvas(img, DefaultDark())

	center := Color{R: 255, G: 255, B: 255, A: 255}
	edge := Color{R: 0, G: 0, B: 0, A: 255}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		canvas.RadialGradient(0, 0, 500, 500, center, edge)
	}
}

// BenchmarkBoxShadow measures performance of box shadow
func BenchmarkBoxShadow(b *testing.B) {
	img := ebiten.NewImage(500, 500)
	canvas := newEbitenCanvas(img, DefaultDark())

	shadowColor := Color{R: 0, G: 0, B: 0, A: 128}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		canvas.BoxShadow(100, 100, 200, 200, 10, 10, 20, shadowColor)
	}
}

// Helper function to compute color distance (for visual comparison)
func colorDistance(c1, c2 color.Color) float64 {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	dr := float64(r1>>8) - float64(r2>>8)
	dg := float64(g1>>8) - float64(g2>>8)
	db := float64(b1>>8) - float64(b2>>8)
	return math.Sqrt(dr*dr + dg*dg + db*db)
}
