//go:build windows || darwin || android || ios

package wayne

import (
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"sync"
	"sync/atomic"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/basicfont"

	textv2 "github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	// ErrUnsupportedImageFormat is returned for unsupported image formats.
	ErrUnsupportedImageFormat = errors.New("wayne: unsupported image format")

	// ErrInvalidFontData is returned when font data is malformed.
	ErrInvalidFontData = errors.New("wayne: invalid font data")

	// ErrResourceManagerClosed is returned when attempting to load resources after cleanup.
	ErrResourceManagerClosed = errors.New("wayne: resource manager is closed")
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
	mu     sync.RWMutex
	closed atomic.Bool

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
//
// Supported formats: TrueType (.ttf), OpenType (.otf) - currently not implemented,
// falls back to embedded font.
//
// The path parameter is currently ignored but will be used in future versions.
// Size is specified in points.
func (rm *ResourceManager) LoadFont(path string, size float64) (*Font, error) {
	if rm.closed.Load() {
		return nil, ErrResourceManagerClosed
	}

	if path == "" {
		return nil, fmt.Errorf("font path cannot be empty")
	}

	if size <= 0 {
		return nil, fmt.Errorf("font size must be positive, got %f", size)
	}

	rm.mu.Lock()
	defer rm.mu.Unlock()

	if rm.defaultFont == nil {
		return nil, fmt.Errorf("default font not initialized: %w", ErrInvalidFontData)
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
// Supported formats: PNG, JPEG, GIF (via standard library image decoders).
//
// Returns an error if the file does not exist, is not accessible, or contains
// unsupported image data. Error messages include the file path for debugging.
func (rm *ResourceManager) LoadImage(path string) (*Image, error) {
	if path == "" {
		return nil, fmt.Errorf("image path cannot be empty")
	}

	// Validate file accessibility before attempting to load
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("resource path %q not accessible: %w", path, err)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open image %q: %w", path, err)
	}
	defer f.Close()

	img, err := rm.LoadImageFromReader(f, path)
	if err != nil {
		return nil, fmt.Errorf("failed to load image %q: %w", path, err)
	}
	return img, nil
}

// LoadImageFromReader loads an image from an io.Reader.
//
// Supported formats: PNG, JPEG, GIF. The filenameHint is used in error messages
// for better debugging context.
//
// Returns ErrUnsupportedImageFormat for unsupported formats, or a descriptive
// error for malformed image data.
func (rm *ResourceManager) LoadImageFromReader(r io.Reader, filenameHint string) (*Image, error) {
	if rm.closed.Load() {
		return nil, ErrResourceManagerClosed
	}

	img, format, err := image.Decode(r)
	if err != nil {
		if filenameHint != "" {
			return nil, fmt.Errorf("failed to decode image from %q: %w", filenameHint, err)
		}
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Check for supported formats
	if format != "png" && format != "jpeg" && format != "gif" {
		if filenameHint != "" {
			return nil, fmt.Errorf("unsupported image format %q in %q: %w", format, filenameHint, ErrUnsupportedImageFormat)
		}
		return nil, fmt.Errorf("unsupported image format %q: %w", format, ErrUnsupportedImageFormat)
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
	defer func() {
		if r := recover(); r != nil {
			// Ignore panics during cleanup to ensure resources are still marked closed
		}
	}()

	rm.closed.Store(true)

	rm.mu.Lock()
	defer rm.mu.Unlock()

	for id := range rm.fonts {
		delete(rm.fonts, id)
	}
	for id, img := range rm.images {
		if img.eimg != nil {
			img.eimg.Deallocate()
		}
		delete(rm.images, id)
	}
	rm.defaultFont = nil
}
