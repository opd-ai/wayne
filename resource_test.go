//go:build windows || darwin || android || ios

package wayne

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"strings"
	"sync"
	"testing"
)

// TestResourceManagerConcurrentCleanup tests that LoadFont and LoadImage
// properly return errors when called concurrently with cleanup.
func TestResourceManagerConcurrentCleanup(t *testing.T) {
	rm := newResourceManager()

	var wg sync.WaitGroup
	errors := make(chan error, 100)

	// Start goroutines that attempt to load resources
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := rm.LoadFont("test.ttf", 12.0)
			if err != nil {
				errors <- err
			}
		}()
	}

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			img := createTestImage()
			_, err := rm.LoadImageFromReader(img, "test.png")
			if err != nil {
				errors <- err
			}
		}()
	}

	// Call cleanup concurrently
	go rm.cleanup()

	wg.Wait()
	close(errors)

	// After cleanup, all new calls should return ErrResourceManagerClosed
	_, err := rm.LoadFont("test.ttf", 12.0)
	if err != ErrResourceManagerClosed {
		t.Errorf("LoadFont after cleanup: expected ErrResourceManagerClosed, got %v", err)
	}

	img := createTestImage()
	_, err = rm.LoadImageFromReader(img, "test.png")
	if err != ErrResourceManagerClosed {
		t.Errorf("LoadImageFromReader after cleanup: expected ErrResourceManagerClosed, got %v", err)
	}
}

// TestResourceManagerCleanupPanic tests that cleanup handles dispose failures gracefully.
func TestResourceManagerCleanupPanic(t *testing.T) {
	rm := newResourceManager()

	// Load some resources
	_, err := rm.LoadFont("test.ttf", 12.0)
	if err != nil {
		t.Fatalf("Failed to load font: %v", err)
	}

	img := createTestImage()
	_, err = rm.LoadImageFromReader(img, "test.png")
	if err != nil {
		t.Fatalf("Failed to load image: %v", err)
	}

	// Cleanup should not panic even if Dispose fails
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("cleanup panicked: %v", r)
			}
		}()
		rm.cleanup()
	}()

	// Verify closed state
	if !rm.closed.Load() {
		t.Error("ResourceManager should be marked as closed after cleanup")
	}
}

// TestResourceManagerLoadAfterClose tests that resource loading after Close returns error.
func TestResourceManagerLoadAfterClose(t *testing.T) {
	rm := newResourceManager()

	// Load some resources successfully
	font, err := rm.LoadFont("test.ttf", 12.0)
	if err != nil {
		t.Fatalf("Failed to load font before cleanup: %v", err)
	}
	if font == nil {
		t.Fatal("LoadFont returned nil font before cleanup")
	}

	img := createTestImage()
	imgRes, err := rm.LoadImageFromReader(img, "test.png")
	if err != nil {
		t.Fatalf("Failed to load image before cleanup: %v", err)
	}
	if imgRes == nil {
		t.Fatal("LoadImageFromReader returned nil image before cleanup")
	}

	// Clean up resources
	rm.cleanup()

	// Attempt to load font after cleanup
	_, err = rm.LoadFont("test2.ttf", 14.0)
	if err != ErrResourceManagerClosed {
		t.Errorf("LoadFont after cleanup: expected ErrResourceManagerClosed, got %v", err)
	}

	// Attempt to load image after cleanup
	img2 := createTestImage()
	_, err = rm.LoadImageFromReader(img2, "test2.png")
	if err != ErrResourceManagerClosed {
		t.Errorf("LoadImageFromReader after cleanup: expected ErrResourceManagerClosed, got %v", err)
	}
}

// TestResourceManagerDefaultFont tests the DefaultFont method.
func TestResourceManagerDefaultFont(t *testing.T) {
	rm := newResourceManager()

	font := rm.DefaultFont()
	if font == nil {
		t.Fatal("DefaultFont returned nil")
	}
	if font.size != 13.0 {
		t.Errorf("DefaultFont size: expected 13.0, got %f", font.size)
	}
	if font.face == nil {
		t.Error("DefaultFont face is nil")
	}
}

// TestResourceManagerLoadFont tests the LoadFont method.
func TestResourceManagerLoadFont(t *testing.T) {
	rm := newResourceManager()

	tests := []struct {
		name string
		path string
		size float64
	}{
		{"Small font", "test.ttf", 10.0},
		{"Medium font", "test.ttf", 14.0},
		{"Large font", "test.ttf", 24.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			font, err := rm.LoadFont(tt.path, tt.size)
			if err != nil {
				t.Fatalf("LoadFont failed: %v", err)
			}
			if font == nil {
				t.Fatal("LoadFont returned nil font")
			}
			if font.size != tt.size {
				t.Errorf("Font size: expected %f, got %f", tt.size, font.size)
			}
			if font.face == nil {
				t.Error("Font face is nil")
			}
		})
	}
}

// TestResourceManagerLoadImage tests the LoadImage and LoadImageFromReader methods.
func TestResourceManagerLoadImage(t *testing.T) {
	rm := newResourceManager()

	img := createTestImage()
	resource, err := rm.LoadImageFromReader(img, "test.png")
	if err != nil {
		t.Fatalf("LoadImageFromReader failed: %v", err)
	}
	if resource == nil {
		t.Fatal("LoadImageFromReader returned nil resource")
	}
	if resource.eimg == nil {
		t.Error("Image eimg is nil")
	}

	width, height := resource.Size()
	if width != 100 || height != 100 {
		t.Errorf("Image size: expected (100, 100), got (%d, %d)", width, height)
	}
}

// TestResourceManagerUnsupportedImageFormat tests error handling for unsupported image formats.
func TestResourceManagerUnsupportedImageFormat(t *testing.T) {
	rm := newResourceManager()

	// Create a BMP image (unsupported format)
	invalidData := []byte("BM")
	reader := bytes.NewReader(invalidData)

	_, err := rm.LoadImageFromReader(reader, "test.bmp")
	if err == nil {
		t.Error("LoadImageFromReader should fail for invalid image data")
	}
}

// TestLoadFontErrors tests error handling for LoadFont.
func TestLoadFontErrors(t *testing.T) {
	tests := []struct {
		name        string
		setupRM     func() *ResourceManager
		path        string
		size        float64
		expectError string
	}{
		{
			name:        "empty path",
			setupRM:     newResourceManager,
			path:        "",
			size:        12.0,
			expectError: "font path cannot be empty",
		},
		{
			name:        "zero size",
			setupRM:     newResourceManager,
			path:        "test.ttf",
			size:        0,
			expectError: "font size must be positive",
		},
		{
			name:        "negative size",
			setupRM:     newResourceManager,
			path:        "test.ttf",
			size:        -10.0,
			expectError: "font size must be positive",
		},
		{
			name: "closed resource manager",
			setupRM: func() *ResourceManager {
				rm := newResourceManager()
				rm.cleanup()
				return rm
			},
			path:        "test.ttf",
			size:        12.0,
			expectError: "resource manager is closed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm := tt.setupRM()
			_, err := rm.LoadFont(tt.path, tt.size)
			if err == nil {
				t.Fatalf("LoadFont should fail for %s", tt.name)
			}
			if tt.expectError != "" && !strings.Contains(err.Error(), tt.expectError) {
				t.Errorf("LoadFont error: expected to contain %q, got %q", tt.expectError, err.Error())
			}
		})
	}
}

// TestLoadImageErrors tests error handling for LoadImage.
func TestLoadImageErrors(t *testing.T) {
	rm := newResourceManager()

	tests := []struct {
		name        string
		path        string
		expectError string
	}{
		{
			name:        "empty path",
			path:        "",
			expectError: "image path cannot be empty",
		},
		{
			name:        "non-existent file",
			path:        "/nonexistent/path/to/image.png",
			expectError: "not accessible",
		},
		{
			name:        "directory instead of file",
			path:        ".",
			expectError: "failed to load image",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := rm.LoadImage(tt.path)
			if err == nil {
				t.Fatalf("LoadImage should fail for %s", tt.name)
			}
			if tt.expectError != "" && !strings.Contains(err.Error(), tt.expectError) {
				t.Errorf("LoadImage error: expected to contain %q, got %q", tt.expectError, err.Error())
			}
		})
	}
}

// TestLoadImageAfterClose tests that LoadImage returns ErrResourceManagerClosed after cleanup.
func TestLoadImageAfterClose(t *testing.T) {
	rm := newResourceManager()
	rm.cleanup()

	_, err := rm.LoadImage("test.png")
	if err != ErrResourceManagerClosed {
		t.Errorf("LoadImage after cleanup: expected ErrResourceManagerClosed, got %v", err)
	}
}

// TestLoadImageFromReaderErrors tests error handling for LoadImageFromReader.
func TestLoadImageFromReaderErrors(t *testing.T) {
	rm := newResourceManager()

	tests := []struct {
		name        string
		data        []byte
		filename    string
		expectError string
	}{
		{
			name:        "invalid image data",
			data:        []byte("not an image"),
			filename:    "test.png",
			expectError: "failed to decode image",
		},
		{
			name:        "corrupt PNG header",
			data:        []byte("\x89PNG\r\n\x1a\nGARBAGE"),
			filename:    "corrupt.png",
			expectError: "failed to decode image",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader(tt.data)
			_, err := rm.LoadImageFromReader(reader, tt.filename)
			if err == nil {
				t.Fatalf("LoadImageFromReader should fail for %s", tt.name)
			}
			if tt.expectError != "" && !strings.Contains(err.Error(), tt.expectError) {
				t.Errorf("LoadImageFromReader error: expected to contain %q, got %q", tt.expectError, err.Error())
			}
			// Verify filename appears in error message
			if tt.filename != "" && !strings.Contains(err.Error(), tt.filename) {
				t.Errorf("Error message should contain filename %q, got %q", tt.filename, err.Error())
			}
		})
	}
}

// TestLoadImageErrorContext tests that error messages include helpful context.
func TestLoadImageErrorContext(t *testing.T) {
	rm := newResourceManager()

	// Test that file path appears in error message
	_, err := rm.LoadImage("/nonexistent/test.png")
	if err == nil {
		t.Fatal("LoadImage should fail for nonexistent file")
	}
	if !strings.Contains(err.Error(), "/nonexistent/test.png") {
		t.Errorf("Error message should contain file path, got: %v", err)
	}

	// Test that format appears in error message for unsupported formats
	invalidData := []byte("not an image")
	reader := bytes.NewReader(invalidData)
	_, err = rm.LoadImageFromReader(reader, "test.xyz")
	if err == nil {
		t.Fatal("LoadImageFromReader should fail for invalid data")
	}
	if !strings.Contains(err.Error(), "test.xyz") {
		t.Errorf("Error message should contain filename hint, got: %v", err)
	}
}

// TestLoadImageGIFSupport tests that GIF images are supported.
func TestLoadImageGIFSupport(t *testing.T) {
	rm := newResourceManager()

	// Create a minimal valid GIF (1x1 pixel, white)
	gifData := []byte{
		0x47, 0x49, 0x46, 0x38, 0x39, 0x61, // GIF89a
		0x01, 0x00, 0x01, 0x00, // width=1, height=1
		0x80, 0x00, 0x00, // global color table flag, background, aspect
		0xFF, 0xFF, 0xFF, // white color
		0x00, 0x00, 0x00, // black color
		0x2C, 0x00, 0x00, 0x00, 0x00, // image descriptor
		0x01, 0x00, 0x01, 0x00, 0x00, // image width=1, height=1
		0x02, 0x02, 0x44, 0x01, 0x00, // image data
		0x3B, // trailer
	}

	reader := bytes.NewReader(gifData)
	img, err := rm.LoadImageFromReader(reader, "test.gif")
	if err != nil {
		t.Fatalf("LoadImageFromReader should support GIF format: %v", err)
	}
	if img == nil {
		t.Fatal("LoadImageFromReader returned nil for valid GIF")
	}
}

// createTestImage creates a simple test image for testing.
func createTestImage() *bytes.Buffer {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}

	var buf bytes.Buffer
	png.Encode(&buf, img)
	return &buf
}
