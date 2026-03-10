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
func (c *ebitenCanvas) DrawText(text string, x, y int, f *Font, col Color) {
	if text == "" {
		return
	}

	face := defaultFontFace
	if f != nil && f.face != nil {
		face = f.face
	}

	op := &textv2.DrawOptions{}
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
func (c *ebitenCanvas) LinearGradient(x, y, width, height int, startColor, endColor Color, angle float64) {
	if width <= 0 || height <= 0 {
		return
	}
	// Approximate with horizontal or vertical bands of interpolated colors.
	steps := width
	isVertical := math.Abs(math.Sin(angle*math.Pi/180)) > 0.5
	if isVertical {
		steps = height
	}
	if steps <= 0 {
		return
	}

	for i := 0; i < steps; i++ {
		t := float64(i) / float64(steps)
		r := uint8(float64(startColor.R)*(1-t) + float64(endColor.R)*t)
		g := uint8(float64(startColor.G)*(1-t) + float64(endColor.G)*t)
		b := uint8(float64(startColor.B)*(1-t) + float64(endColor.B)*t)
		a := uint8(float64(startColor.A)*(1-t) + float64(endColor.A)*t)
		clr := color.RGBA{R: r, G: g, B: b, A: a}

		if isVertical {
			vector.FillRect(c.dst, float32(x), float32(y+i), float32(width), 1, clr, false)
		} else {
			vector.FillRect(c.dst, float32(x+i), float32(y), 1, float32(height), clr, false)
		}
	}
}

// RadialGradient fills a rectangle with a radial gradient.
func (c *ebitenCanvas) RadialGradient(x, y, width, height int, centerColor, edgeColor Color) {
	if width <= 0 || height <= 0 {
		return
	}
	cx := float64(x + width/2)
	cy := float64(y + height/2)
	maxR := math.Sqrt(float64(width*width+height*height)) / 2.0

	steps := 32 // number of rings
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
