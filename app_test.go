//go:build windows || darwin || android || ios

package wayne

import (
	"testing"
)

func TestNewApp(t *testing.T) {
	app := NewApp()
	if app == nil {
		t.Fatal("NewApp() returned nil")
	}
	
	if app.width != 800 {
		t.Errorf("Expected default width 800, got %d", app.width)
	}
	if app.height != 600 {
		t.Errorf("Expected default height 600, got %d", app.height)
	}
	if app.resources == nil {
		t.Error("Expected resources to be initialized")
	}
	if app.dispatcher == nil {
		t.Error("Expected dispatcher to be initialized")
	}
	if app.notifyChan == nil {
		t.Error("Expected notifyChan to be initialized")
	}
}

func TestNewAppWithConfig(t *testing.T) {
	tests := []struct {
		name           string
		config         AppConfig
		expectedWidth  int
		expectedHeight int
		expectedVerbose bool
	}{
		{
			name:           "custom dimensions",
			config:         AppConfig{Width: 1024, Height: 768, Verbose: true},
			expectedWidth:  1024,
			expectedHeight: 768,
			expectedVerbose: true,
		},
		{
			name:           "zero width defaults to 800",
			config:         AppConfig{Width: 0, Height: 768},
			expectedWidth:  800,
			expectedHeight: 768,
			expectedVerbose: false,
		},
		{
			name:           "zero height defaults to 600",
			config:         AppConfig{Width: 1024, Height: 0},
			expectedWidth:  1024,
			expectedHeight: 600,
			expectedVerbose: false,
		},
		{
			name:           "negative dimensions default",
			config:         AppConfig{Width: -100, Height: -100},
			expectedWidth:  800,
			expectedHeight: 600,
			expectedVerbose: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewAppWithConfig(tt.config)
			if app == nil {
				t.Fatal("NewAppWithConfig() returned nil")
			}
			if app.width != tt.expectedWidth {
				t.Errorf("Expected width %d, got %d", tt.expectedWidth, app.width)
			}
			if app.height != tt.expectedHeight {
				t.Errorf("Expected height %d, got %d", tt.expectedHeight, app.height)
			}
			if app.verbose != tt.expectedVerbose {
				t.Errorf("Expected verbose %v, got %v", tt.expectedVerbose, app.verbose)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Width != 800 {
		t.Errorf("Expected default width 800, got %d", cfg.Width)
	}
	if cfg.Height != 600 {
		t.Errorf("Expected default height 600, got %d", cfg.Height)
	}
	if cfg.Verbose != false {
		t.Errorf("Expected default verbose false, got %v", cfg.Verbose)
	}
}

func TestAppSetTheme(t *testing.T) {
	app := NewApp()
	customTheme := DefaultLight()
	
	app.SetTheme(customTheme)
	
	// Theme should be stored
	if app.theme.Background.R != customTheme.Background.R {
		t.Error("Theme not properly set")
	}
}

func TestAppQuit(t *testing.T) {
	app := NewApp()
	
	// Initially should not be quitting
	if app.quitFlag.Load() {
		t.Error("App should not be quitting initially")
	}
	
	app.Quit()
	
	// After Quit() should be set
	if !app.quitFlag.Load() {
		t.Error("App should be quitting after Quit() call")
	}
}

func TestAppClose(t *testing.T) {
	app := NewApp()
	
	// Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Close() panicked: %v", r)
		}
	}()
	
	app.Close()
	
	// After close, resources should be nil
	if app.resources != nil {
		t.Error("Resources should be nil after Close()")
	}
}

func TestAppRunAfterRun(t *testing.T) {
	// This test verifies ErrAlreadyRunning is defined and accessible
	// We cannot actually test Run() in unit tests as it blocks indefinitely
	if ErrAlreadyRunning == nil {
		t.Error("ErrAlreadyRunning should be defined")
	}
	
	if ErrAlreadyRunning.Error() != "wayne: app already running" {
		t.Errorf("Unexpected error message: %s", ErrAlreadyRunning.Error())
	}
}

func TestAppNotRunningError(t *testing.T) {
	if ErrNotRunning == nil {
		t.Error("ErrNotRunning should be defined")
	}
	
	if ErrNotRunning.Error() != "wayne: app not running" {
		t.Errorf("Unexpected error message: %s", ErrNotRunning.Error())
	}
}

func TestAppInvalidWindowConfigError(t *testing.T) {
	if ErrInvalidWindowConfig == nil {
		t.Error("ErrInvalidWindowConfig should be defined")
	}
	
	if ErrInvalidWindowConfig.Error() != "wayne: invalid window configuration" {
		t.Errorf("Unexpected error message: %s", ErrInvalidWindowConfig.Error())
	}
}

func TestNewWindow(t *testing.T) {
	app := NewApp()
	
	tests := []struct {
		name   string
		config WindowConfig
		valid  bool
	}{
		{
			name:   "valid config",
			config: WindowConfig{Title: "Test", Width: 800, Height: 600},
			valid:  true,
		},
		{
			name:   "zero width defaults",
			config: WindowConfig{Title: "Test", Width: 0, Height: 600},
			valid:  true,
		},
		{
			name:   "zero height defaults",
			config: WindowConfig{Title: "Test", Width: 800, Height: 0},
			valid:  true,
		},
		{
			name:   "empty title",
			config: WindowConfig{Title: "", Width: 800, Height: 600},
			valid:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			win, err := app.NewWindow(tt.config)
			if tt.valid {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if win == nil {
					t.Error("Expected window to be created")
				}
			}
		})
	}
}

func TestLoadFont(t *testing.T) {
	app := NewApp()
	
	// Test with empty path
	font, err := app.LoadFont("", 12)
	if font != nil {
		t.Error("Expected nil font for empty path")
	}
	if err == nil {
		t.Error("Expected error for empty font path")
	}
	
	// Test with non-existent path
	font, err = app.LoadFont("/nonexistent/font.ttf", 12)
	if font != nil {
		t.Error("Expected nil font for non-existent path")
	}
	if err == nil {
		t.Error("Expected error for non-existent font path")
	}
}

func TestLoadImage(t *testing.T) {
	app := NewApp()
	
	// Test with empty path
	img, err := app.LoadImage("")
	if img != nil {
		t.Error("Expected nil image for empty path")
	}
	if err == nil {
		t.Error("Expected error for empty image path")
	}
	
	// Test with non-existent path
	img, err = app.LoadImage("/nonexistent/image.png")
	if img != nil {
		t.Error("Expected nil image for non-existent path")
	}
	if err == nil {
		t.Error("Expected error for non-existent image path")
	}
}

func TestDefaultFont(t *testing.T) {
	app := NewApp()
	
	font := app.DefaultFont()
	
	// Should return a font (embedded default)
	if font == nil {
		t.Error("DefaultFont should not return nil")
	}
}
