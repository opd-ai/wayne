//go:build windows || darwin || android || ios

package wayne

import (
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/basicfont"

	textv2 "github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	// ErrUnsupportedImageFormat is returned for unsupported image formats.
	ErrUnsupportedImageFormat = errors.New("wayne: unsupported image format")

	// ErrInvalidFontData is returned when font data is malformed.
	ErrInvalidFontData = errors.New("wayne: invalid font data")
)

// Font represents a loaded font resource.
type Font struct {
	face textv2.Face
	size float64
}

// Size returns the font size in points.
func (f *Font) Size() float64 {
	return f.size
}

// Image represents a loaded image resource.
type Image struct {
	eimg   *ebiten.Image
	width  int
	height int
}

// Size returns the image dimensions in pixels.
func (img *Image) Size() (width, height int) {
	return img.width, img.height
}

// ResourceManager manages fonts and images for an application.
type ResourceManager struct {
	mu sync.RWMutex

	defaultFont *Font

	fonts  map[int]*Font
	images map[int]*Image

	nextFontID  int
	nextImageID int
}

// newResourceManager creates a new resource manager.
func newResourceManager() *ResourceManager {
	rm := &ResourceManager{
		fonts:       make(map[int]*Font),
		images:      make(map[int]*Image),
		nextFontID:  1,
		nextImageID: 1,
	}
	rm.defaultFont = &Font{
		face: textv2.NewGoXFace(basicfont.Face7x13),
		size: 13.0,
	}
	return rm
}

// DefaultFont returns the embedded default font.
func (rm *ResourceManager) DefaultFont() *Font {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.defaultFont
}

// LoadFont loads a font from the specified path at the given size.
//
// Currently returns the default embedded font with the requested size.
// Custom TTF loading may be added in a future version.
func (rm *ResourceManager) LoadFont(path string, size float64) (*Font, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if rm.defaultFont == nil {
		return nil, ErrInvalidFontData
	}

	font := &Font{
		face: rm.defaultFont.face,
		size: size,
	}
	rm.fonts[rm.nextFontID] = font
	rm.nextFontID++
	return font, nil
}

// LoadImage loads an image from the specified path.
//
// Supported formats: PNG, JPEG.
func (rm *ResourceManager) LoadImage(path string) (*Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return rm.LoadImageFromReader(f, path)
}

// LoadImageFromReader loads an image from an io.Reader.
func (rm *ResourceManager) LoadImageFromReader(r io.Reader, filenameHint string) (*Image, error) {
	img, format, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	if format != "png" && format != "jpeg" {
		return nil, ErrUnsupportedImageFormat
	}

	eimg := ebiten.NewImageFromImage(img)
	bounds := eimg.Bounds()

	rm.mu.Lock()
	defer rm.mu.Unlock()

	resource := &Image{
		eimg:   eimg,
		width:  bounds.Dx(),
		height: bounds.Dy(),
	}
	rm.images[rm.nextImageID] = resource
	rm.nextImageID++
	return resource, nil
}

// cleanup releases all loaded resources.
func (rm *ResourceManager) cleanup() {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	for id := range rm.fonts {
		delete(rm.fonts, id)
	}
	for id := range rm.images {
		delete(rm.images, id)
	}
	rm.defaultFont = nil
}
