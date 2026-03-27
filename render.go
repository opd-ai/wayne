//go:build windows || darwin || android || ios

package wayne

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	textv2 "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

// defaultFontFace is the package-level default font for text rendering.
var defaultFontFace textv2.Face

func init() {
	defaultFontFace = textv2.NewGoXFace(basicfont.Face7x13)
}

// colorToRGBA converts a wayne Color to a standard image/color.RGBA.
func colorToRGBA(c Color) color.RGBA {
	return color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}
}

// ebitenCanvas implements the Canvas interface using an *ebiten.Image as the
// drawing target.
type ebitenCanvas struct {
	dst   *ebiten.Image
	theme Theme
}

// newEbitenCanvas creates a Canvas backed by an ebiten.Image with the given theme.
func newEbitenCanvas(dst *ebiten.Image, theme Theme) Canvas {
	return &ebitenCanvas{dst: dst, theme: theme}
}

// FillRect fills a solid rectangle.
func (c *ebitenCanvas) FillRect(x, y, width, height int, col Color) {
	if width <= 0 || height <= 0 {
		return
	}
	vector.FillRect(c.dst,
		float32(x), float32(y),
		float32(width), float32(height),
		colorToRGBA(col), false)
}

// FillRoundedRect fills a rounded rectangle with the specified corner radius.
func (c *ebitenCanvas) FillRoundedRect(x, y, width, height, radius int, col Color) {
	if width <= 0 || height <= 0 {
		return
	}
	if radius <= 0 {
		c.FillRect(x, y, width, height, col)
		return
	}
	// Clamp radius to half the smaller dimension.
	maxR := min(width, height) / 2
	if radius > maxR {
		radius = maxR
	}

	clr := colorToRGBA(col)
	r := float32(radius)
	fx, fy := float32(x), float32(y)
	fw, fh := float32(width), float32(height)

	// Draw three overlapping rectangles to cover the rounded rect body.
	vector.FillRect(c.dst, fx+r, fy, fw-2*r, fh, clr, false)
	vector.FillRect(c.dst, fx, fy+r, fw, fh-2*r, clr, false)

	// Draw four corner circles.
	vector.FillCircle(c.dst, fx+r, fy+r, r, clr, false)
	vector.FillCircle(c.dst, fx+fw-r, fy+r, r, clr, false)
	vector.FillCircle(c.dst, fx+r, fy+fh-r, r, clr, false)
	vector.FillCircle(c.dst, fx+fw-r, fy+fh-r, r, clr, false)
}

// DrawLine draws a line from (x1,y1) to (x2,y2).
func (c *ebitenCanvas) DrawLine(x1, y1, x2, y2 int, col Color, thickness int) {
	if thickness <= 0 {
		thickness = 1
	}
	vector.StrokeLine(c.dst,
		float32(x1), float32(y1),
		float32(x2), float32(y2),
		float32(thickness),
		colorToRGBA(col), false)
}

// DrawText renders text at (x, y) using the specified font and color.
// If font is nil, the package default font is used.
// The text is scaled according to the theme's HiDPI scale factor.
func (c *ebitenCanvas) DrawText(text string, x, y int, f *Font, col Color) {
	if text == "" {
		return
	}

	face := defaultFontFace
	scale := c.Scale() // Apply HiDPI scale factor
	if f != nil && f.face != nil {
		face = f.face
		// Scale text based on requested font size relative to default (13pt)
		if f.size > 0 {
			scale = (f.size / 13.0) * c.Scale()
		}
	}

	op := &textv2.DrawOptions{}
	if scale != 1.0 {
		op.GeoM.Scale(scale, scale)
	}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(colorToRGBA(col))
	textv2.Draw(c.dst, text, face, op)
}

// DrawImage renders an image at the given position and size.
func (c *ebitenCanvas) DrawImage(img *Image, x, y, width, height int) {
	if img == nil || img.eimg == nil {
		return
	}
	bw := img.eimg.Bounds().Dx()
	bh := img.eimg.Bounds().Dy()
	if bw == 0 || bh == 0 {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(width)/float64(bw), float64(height)/float64(bh))
	op.GeoM.Translate(float64(x), float64(y))
	c.dst.DrawImage(img.eimg, op)
}

// LinearGradient fills a rectangle with a linear gradient.
// angle is in degrees (0 = left-to-right, 90 = top-to-bottom).
//
// Implementation note: This method uses band rendering to approximate arbitrary-angle
// gradients. Diagonal gradients use scan-line rendering with color interpolation based
// on vector projection. For optimal performance and quality, prefer axis-aligned angles
// (0°, 90°, 180°, 270°).
func (c *ebitenCanvas) LinearGradient(x, y, width, height int, startColor, endColor Color, angle float64) {
	if width <= 0 || height <= 0 {
		return
	}

	params := computeGradientParams(x, y, width, height, angle)
	if params.gradientLength == 0 {
		return
	}

	c.renderGradient(x, y, width, height, startColor, endColor, params)
}

// gradientParams holds precomputed values for gradient rendering.
type gradientParams struct {
	angle          float64
	dx, dy         float64
	centerX        float64
	centerY        float64
	minProj        float64
	gradientLength float64
	sinA, cosA     float64
}

// computeGradientParams calculates gradient direction and projection bounds.
func computeGradientParams(x, y, width, height int, angle float64) gradientParams {
	// Normalize angle to [0, 360)
	angle = math.Mod(angle, 360)
	if angle < 0 {
		angle += 360
	}

	angleRad := angle * math.Pi / 180.0
	dx := math.Cos(angleRad)
	dy := math.Sin(angleRad)

	centerX := float64(x) + float64(width)/2.0
	centerY := float64(y) + float64(height)/2.0

	minProj, maxProj := computeCornerProjections(x, y, width, height, centerX, centerY, dx, dy)

	return gradientParams{
		angle:          angle,
		dx:             dx,
		dy:             dy,
		centerX:        centerX,
		centerY:        centerY,
		minProj:        minProj,
		gradientLength: maxProj - minProj,
		sinA:           math.Abs(math.Sin(angleRad)),
		cosA:           math.Abs(math.Cos(angleRad)),
	}
}

// computeCornerProjections projects rectangle corners onto gradient direction.
func computeCornerProjections(x, y, width, height int, centerX, centerY, dx, dy float64) (minProj, maxProj float64) {
	corners := [][2]float64{
		{float64(x), float64(y)},
		{float64(x + width), float64(y)},
		{float64(x), float64(y + height)},
		{float64(x + width), float64(y + height)},
	}

	minProj, maxProj = math.Inf(1), math.Inf(-1)
	for _, corner := range corners {
		proj := (corner[0]-centerX)*dx + (corner[1]-centerY)*dy
		if proj < minProj {
			minProj = proj
		}
		if proj > maxProj {
			maxProj = proj
		}
	}
	return minProj, maxProj
}

// renderGradient dispatches to the appropriate rendering strategy.
func (c *ebitenCanvas) renderGradient(x, y, width, height int, startColor, endColor Color, p gradientParams) {
	switch {
	case p.cosA > 0.99:
		c.renderHorizontalGradient(x, y, width, height, startColor, endColor, p.angle)
	case p.sinA > 0.99:
		c.renderVerticalGradient(x, y, width, height, startColor, endColor, p.angle)
	default:
		c.renderDiagonalGradient(x, y, width, height, startColor, endColor, p)
	}
}

// renderHorizontalGradient renders a nearly horizontal gradient (0° or 180°).
func (c *ebitenCanvas) renderHorizontalGradient(x, y, width, height int, startColor, endColor Color, angle float64) {
	for i := 0; i < width; i++ {
		t := float64(i) / float64(width)
		if angle > 90 && angle < 270 {
			t = 1 - t
		}
		clr := interpolateColor(startColor, endColor, t)
		vector.FillRect(c.dst, float32(x+i), float32(y), 1, float32(height), clr, false)
	}
}

// renderVerticalGradient renders a nearly vertical gradient (90° or 270°).
func (c *ebitenCanvas) renderVerticalGradient(x, y, width, height int, startColor, endColor Color, angle float64) {
	for i := 0; i < height; i++ {
		t := float64(i) / float64(height)
		if angle > 180 {
			t = 1 - t
		}
		clr := interpolateColor(startColor, endColor, t)
		vector.FillRect(c.dst, float32(x), float32(y+i), float32(width), 1, clr, false)
	}
}

// renderDiagonalGradient renders an arbitrary-angle gradient using scanlines.
func (c *ebitenCanvas) renderDiagonalGradient(x, y, width, height int, startColor, endColor Color, p gradientParams) {
	for py := 0; py < height; py++ {
		wy := float64(y + py)
		for px := 0; px < width; px++ {
			t := computePixelGradientT(float64(x+px), wy, p)
			clr := interpolateColor(startColor, endColor, t)
			vector.FillRect(c.dst, float32(x+px), float32(y+py), 1, 1, clr, false)
		}
	}
}

// computePixelGradientT computes the normalized gradient position for a pixel.
func computePixelGradientT(wx, wy float64, p gradientParams) float64 {
	proj := (wx-p.centerX)*p.dx + (wy-p.centerY)*p.dy
	t := (proj - p.minProj) / p.gradientLength
	if t < 0 {
		return 0
	}
	if t > 1 {
		return 1
	}
	return t
}

// interpolateColor interpolates between two colors using parameter t [0, 1]
func interpolateColor(start, end Color, t float64) color.RGBA {
	return color.RGBA{
		R: uint8(float64(start.R)*(1-t) + float64(end.R)*t),
		G: uint8(float64(start.G)*(1-t) + float64(end.G)*t),
		B: uint8(float64(start.B)*(1-t) + float64(end.B)*t),
		A: uint8(float64(start.A)*(1-t) + float64(end.A)*t),
	}
}

// RadialGradient fills a rectangle with a radial gradient.
//
// Implementation note: This method renders concentric circles to approximate
// a radial gradient. The number of rings is dynamically computed based on the
// gradient radius to balance quality and performance. Very large gradients
// (>1000px diagonal) may show subtle banding on high-DPI displays. The gradient
// is centered in the rectangle and extends to the corners.
func (c *ebitenCanvas) RadialGradient(x, y, width, height int, centerColor, edgeColor Color) {
	if width <= 0 || height <= 0 {
		return
	}
	cx := float64(x + width/2)
	cy := float64(y + height/2)
	maxR := math.Sqrt(float64(width*width+height*height)) / 2.0

	// Dynamic ring count: use at least 128 rings, or scale with radius for large gradients
	steps := int(math.Max(128, maxR))
	if steps > 512 {
		steps = 512 // Cap at 512 for performance on very large gradients
	}

	for i := steps - 1; i >= 0; i-- {
		t := float64(i) / float64(steps)
		r := uint8(float64(centerColor.R)*t + float64(edgeColor.R)*(1-t))
		g := uint8(float64(centerColor.G)*t + float64(edgeColor.G)*(1-t))
		b := uint8(float64(centerColor.B)*t + float64(edgeColor.B)*(1-t))
		a := uint8(float64(centerColor.A)*t + float64(edgeColor.A)*(1-t))
		clr := color.RGBA{R: r, G: g, B: b, A: a}
		ri := float32(maxR * float64(steps-i) / float64(steps))
		vector.FillCircle(c.dst, float32(cx), float32(cy), ri, clr, false)
	}
}

// BoxShadow renders a box shadow around the given rectangle.
//
// Implementation note: This method provides a simplified shadow approximation
// using a rounded rectangle with reduced opacity. It does NOT implement Gaussian
// blur or true soft shadows as found in CSS box-shadow. The 'blur' parameter
// controls the shadow rectangle's corner radius and expansion, not blur kernel size.
// For production-quality shadows, consider pre-rendering shadow images with proper
// blur and using Canvas.DrawImage instead. This method is suitable for simple
// drop-shadow effects.
func (c *ebitenCanvas) BoxShadow(x, y, width, height, offsetX, offsetY, blur int, col Color) {
	if width <= 0 || height <= 0 {
		return
	}
	// Draw a semi-transparent shadow rectangle offset from the main rect.
	shadowCol := col.WithAlpha(uint8(float64(col.A) * 0.5))
	shadowX := x + offsetX
	shadowY := y + offsetY
	c.FillRoundedRect(shadowX-blur, shadowY-blur, width+2*blur, height+2*blur, blur, shadowCol)
}

// Theme returns the application-wide theme for this rendering context.
func (c *ebitenCanvas) Theme() Theme {
	return c.theme
}

// Scale returns the current HiDPI scale factor from the theme.
// A value of 1.0 means standard resolution, 2.0 means retina/HiDPI.
func (c *ebitenCanvas) Scale() float64 {
	if c.theme.Scale <= 0 {
		return 1.0
	}
	return c.theme.Scale
}
