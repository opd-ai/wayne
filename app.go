//go:build windows || darwin || android || ios

// Package wayne provides cross-platform UI widgets. See doc.go for details.
package wayne

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// ErrNotRunning is returned when calling methods that require Run() to be active.
	ErrNotRunning = errors.New("wayne: app not running")

	// ErrAlreadyRunning is returned when Run() is called more than once.
	ErrAlreadyRunning = errors.New("wayne: app already running")

	// ErrInvalidWindowConfig is returned when window configuration is invalid.
	ErrInvalidWindowConfig = errors.New("wayne: invalid window configuration")
)

// AppConfig contains configuration options for creating an App.
type AppConfig struct {
	// Width is the initial window width in pixels (default: 800).
	Width int

	// Height is the initial window height in pixels (default: 600).
	Height int

	// Verbose enables logging of backend selection decisions (default: false).
	Verbose bool
}

// DefaultConfig returns the default application configuration.
func DefaultConfig() AppConfig {
	return AppConfig{
		Width:   800,
		Height:  600,
		Verbose: false,
	}
}

// App represents a UI application backed by Ebitengine.
//
// App manages the main event loop and the primary window. It is the entry point
// for all wayne applications.
//
// Example:
//
//	app := wayne.NewApp()
//	defer app.Close()
//	win, _ := app.NewWindow(wayne.WindowConfig{Title: "Hello", Width: 800, Height: 600})
//	win.Show()
//	app.Run()
type App struct {
	mu sync.Mutex

	width   int
	height  int
	verbose bool

	theme     Theme
	resources *ResourceManager

	windows       []*Window
	primaryWindow *Window
	dispatcher    *EventDispatcher

	primaryConfig WindowConfig

	quitFlag   atomic.Bool
	notifyChan chan func()

	lastMX, lastMY int
}

// NewApp creates a new application with default configuration.
func NewApp() *App {
	return NewAppWithConfig(DefaultConfig())
}

// NewAppWithConfig creates a new application with the specified configuration.
func NewAppWithConfig(cfg AppConfig) *App {
	if cfg.Width <= 0 {
		cfg.Width = 800
	}
	if cfg.Height <= 0 {
		cfg.Height = 600
	}

	a := &App{
		width:      cfg.Width,
		height:     cfg.Height,
		verbose:    cfg.Verbose,
		theme:      DefaultDark(),
		resources:  newResourceManager(),
		dispatcher: NewEventDispatcher(),
		notifyChan: make(chan func(), 256),
	}
	return a
}

// SetTheme sets the application-wide theme.
func (a *App) SetTheme(theme Theme) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.theme = theme
}

// SetRoot sets the root widget on the primary window.
// This is a convenience method equivalent to primaryWindow.SetRoot(w).
func (a *App) SetRoot(w PublicWidget) {
	a.mu.Lock()
	win := a.primaryWindow
	a.mu.Unlock()
	if win != nil {
		win.SetRoot(w)
	}
}

// Notify schedules a function to be called from the main goroutine on the next tick.
// This is safe to call from any goroutine.
func (a *App) Notify(fn func()) {
	select {
	case a.notifyChan <- fn:
	default:
		// Channel full; drop notification to avoid blocking.
	}
}

// Run starts the main event loop. This blocks until the app is quit or an error occurs.
// Run must be called from the main goroutine.
func (a *App) Run() error {
	ebiten.SetWindowTitle(a.primaryWindowTitle())
	ebiten.SetWindowSize(a.width, a.height)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// Apply primary window configuration
	a.mu.Lock()
	cfg := a.primaryConfig
	a.mu.Unlock()

	// Apply window size limits if specified
	if cfg.MinWidth > 0 || cfg.MinHeight > 0 || cfg.MaxWidth > 0 || cfg.MaxHeight > 0 {
		minW := cfg.MinWidth
		minH := cfg.MinHeight
		maxW := cfg.MaxWidth
		maxH := cfg.MaxHeight

		// Use -1 for unspecified limits (Ebitengine convention)
		if minW <= 0 {
			minW = -1
		}
		if minH <= 0 {
			minH = -1
		}
		if maxW <= 0 {
			maxW = -1
		}
		if maxH <= 0 {
			maxH = -1
		}

		ebiten.SetWindowSizeLimits(minW, minH, maxW, maxH)
	}

	// Apply fullscreen setting
	if cfg.Fullscreen {
		ebiten.SetFullscreen(true)
	}

	// Apply window decoration setting
	if cfg.Decorations != nil {
		ebiten.SetWindowDecorated(*cfg.Decorations)
	}

	g := &ebitenGame{app: a}
	return ebiten.RunGame(g)
}

// Quit signals the application to exit on the next tick.
func (a *App) Quit() {
	a.quitFlag.Store(true)
}

// Close releases resources associated with the app.
// It does not stop a running event loop; call Quit() for that.
func (a *App) Close() {
	if a.resources != nil {
		a.resources.cleanup()
		a.resources = nil
	}
}

// NewWindow creates a new window.
//
// Since Ebitengine manages a single OS window, the first call to NewWindow
// configures the primary window. Subsequent calls create logical windows that
// are rendered as overlapping root widget subtrees.
func (a *App) NewWindow(cfg WindowConfig) (*Window, error) {
	if cfg.Width <= 0 {
		cfg.Width = a.width
	}
	if cfg.Height <= 0 {
		cfg.Height = a.height
	}

	win := &Window{
		app:        a,
		title:      cfg.Title,
		width:      cfg.Width,
		height:     cfg.Height,
		dispatcher: NewEventDispatcher(),
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	a.windows = append(a.windows, win)
	if a.primaryWindow == nil {
		a.primaryWindow = win
		a.primaryConfig = cfg
		// Resize the Ebitengine window for the first window.
		if cfg.Width > 0 && cfg.Width != a.width {
			a.width = cfg.Width
		}
		if cfg.Height > 0 && cfg.Height != a.height {
			a.height = cfg.Height
		}
	}

	return win, nil
}

// LoadFont loads a font from the specified path at the given size.
//
// Supported formats: TrueType (.ttf), OpenType (.otf) - currently not implemented,
// falls back to embedded font.
//
// Size is specified in points and must be positive.
//
// Returns an error if the app is not running or if the font parameters are invalid.
func (a *App) LoadFont(path string, size float64) (*Font, error) {
	if a.resources == nil {
		return nil, ErrNotRunning
	}
	return a.resources.LoadFont(path, size)
}

// LoadImage loads an image from the specified path.
//
// Supported formats: PNG, JPEG, GIF (via standard library image decoders).
//
// Returns an error if the app is not running, the file does not exist,
// or the image format is unsupported.
func (a *App) LoadImage(path string) (*Image, error) {
	if a.resources == nil {
		return nil, ErrNotRunning
	}
	return a.resources.LoadImage(path)
}

// DefaultFont returns the embedded default font.
func (a *App) DefaultFont() *Font {
	if a.resources == nil {
		return nil
	}
	return a.resources.DefaultFont()
}

// --- internal helpers ---

func (a *App) shouldQuit() bool {
	return a.quitFlag.Load()
}

func (a *App) drainNotify() {
	for {
		select {
		case fn := <-a.notifyChan:
			fn()
		default:
			return
		}
	}
}

func (a *App) primaryRoot() PublicWidget {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.primaryWindow != nil {
		return a.primaryWindow.root
	}
	return nil
}

func (a *App) dimensions() (int, int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.width, a.height
}

func (a *App) primaryWindowTitle() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.primaryWindow != nil {
		return a.primaryWindow.title
	}
	return "wayne"
}

func (a *App) dispatchEvent(evt Event) {
	// Dispatch to the primary window's dispatcher.
	a.mu.Lock()
	win := a.primaryWindow
	a.mu.Unlock()
	if win != nil && win.dispatcher != nil {
		win.dispatcher.Dispatch(evt)
	}
	// Also dispatch directly to the root widget.
	root := a.primaryRoot()
	if root != nil {
		root.HandleEvent(evt)
	}
}

func (a *App) lastMousePos() [2]int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return [2]int{a.lastMX, a.lastMY}
}

func (a *App) setLastMousePos(x, y int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.lastMX = x
	a.lastMY = y
}

// WindowConfig contains configuration options for creating a Window.
type WindowConfig struct {
	// Title is the window title (default: "").
	Title string

	// Width is the initial window width in pixels (default: 800).
	Width int

	// Height is the initial window height in pixels (default: 600).
	Height int

	// MinWidth is the minimum window width in pixels.
	MinWidth int

	// MinHeight is the minimum window height in pixels.
	MinHeight int

	// MaxWidth is the maximum window width in pixels.
	MaxWidth int

	// MaxHeight is the maximum window height in pixels.
	MaxHeight int

	// Fullscreen indicates whether the window should start fullscreen.
	Fullscreen bool

	// Decorations controls window decorations (title bar, borders).
	// nil (default) = decorated window (default platform behavior)
	// true = force decorated window
	// false = borderless window (no title bar, borders)
	//
	// Note: On mobile platforms (Android, iOS), this setting has no effect
	// as Ebitengine does not support window decoration control on mobile.
	Decorations *bool
}

// Window represents a UI window.
//
// Since Ebitengine manages a single OS window, the primary Window corresponds
// to the Ebitengine window. Methods like SetTitle and Show delegate to Ebitengine.
type Window struct {
	mu sync.Mutex

	app    *App
	title  string
	width  int
	height int
	root   PublicWidget
	closed bool

	dispatcher *EventDispatcher
}

// Show displays the window. For Ebitengine, this is a no-op because the window
// appears automatically when App.Run() is called.
func (w *Window) Show() {
	// Ebitengine shows the window when RunGame is called.
	// For the primary window, optionally set the title.
	w.mu.Lock()
	title := w.title
	w.mu.Unlock()
	ebiten.SetWindowTitle(title)
}

// SetTitle sets the window title.
func (w *Window) SetTitle(title string) {
	w.mu.Lock()
	w.title = title
	w.mu.Unlock()
	ebiten.SetWindowTitle(title)
}

// SetRoot sets the root widget for this window.
func (w *Window) SetRoot(root PublicWidget) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.root = root
	if w.dispatcher != nil {
		w.dispatcher.SetWidgetRoot(root)
	}
}

// Close closes the window. If this is the primary window, it quits the app.
func (w *Window) Close() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed {
		return
	}
	w.closed = true
	// Quit the Ebitengine loop if the primary window is closed.
	if w.app != nil && w.app.primaryWindow == w {
		w.app.Quit()
	}
}
