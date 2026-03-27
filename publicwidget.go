//go:build windows || darwin || android || ios || linux

package wayne

// PublicWidget is the stable public interface for all UI widgets in wayne.
//
// PublicWidget provides a simplified, stable contract for application developers.
// All concrete widget types implement this interface.
type PublicWidget interface {
	// Bounds returns the current pixel dimensions of the widget.
	Bounds() (width, height int)

	// HandleEvent processes a user interaction event. Returns true if consumed.
	HandleEvent(Event) bool

	// Draw renders the widget to the provided canvas.
	Draw(Canvas)

	// SetFocused sets the widget's keyboard focus state.
	SetFocused(focused bool)

	// IsFocused returns true if the widget currently has keyboard focus.
	IsFocused() bool
}

// Container extends PublicWidget for widgets that can contain child widgets.
type Container interface {
	PublicWidget

	// Add appends a child widget to this container.
	Add(child PublicWidget)

	// Children returns a slice of the container's child widgets.
	Children() []PublicWidget
}

// Themeable is the interface for widgets that support theme customization.
type Themeable interface {
	// SetTheme applies a custom theme to this widget.
	SetTheme(Theme)
}

// Focusable is the interface for widgets that can receive keyboard focus.
type Focusable interface {
	PublicWidget

	// CanTakeFocus returns true if the widget can currently receive focus.
	// Disabled or hidden widgets should return false.
	CanTakeFocus() bool
}

// Canvas provides a high-level drawing API for widget rendering.
//
// Canvas abstracts over the internal rendering backend (Ebitengine). Methods
// accept pixel coordinates and handle GPU rendering automatically.
//
// Canvas instances are provided by the framework during widget rendering;
// application code does not create Canvas instances directly.
type Canvas interface {
	// FillRect fills a solid rectangle at the given position and size.
	FillRect(x, y, width, height int, color Color)

	// FillRoundedRect fills a rounded rectangle with the specified corner radius.
	FillRoundedRect(x, y, width, height, radius int, color Color)

	// DrawLine draws a line segment from (x1, y1) to (x2, y2).
	DrawLine(x1, y1, x2, y2 int, color Color, thickness int)

	// DrawText renders text at the given position using the specified font and color.
	DrawText(text string, x, y int, font *Font, color Color)

	// DrawImage renders an image at the given position and size.
	DrawImage(img *Image, x, y, width, height int)

	// LinearGradient fills a rectangle with a linear gradient from startColor to endColor.
	// The angle parameter specifies the gradient direction in degrees (0 = left-to-right,
	// 90 = top-to-bottom, 180 = right-to-left, 270 = bottom-to-top).
	//
	// Note: Uses approximation rendering that may show banding on large gradients (>500px).
	LinearGradient(x, y, width, height int, startColor, endColor Color, angle float64)

	// RadialGradient fills a rectangle with a radial gradient from centerColor to edgeColor.
	// The gradient radiates from the rectangle's center to its corners.
	//
	// Note: Uses concentric circle approximation. Very large gradients (>1000px diagonal)
	// may show subtle banding on high-DPI displays.
	RadialGradient(x, y, width, height int, centerColor, edgeColor Color)

	// BoxShadow renders a simplified shadow around the given rectangle.
	// offsetX and offsetY control shadow position, blur controls spread and corner radius.
	//
	// Note: This is a simplified approximation without Gaussian blur. Does not match CSS
	// box-shadow semantics. For production-quality shadows, consider pre-rendered images.
	BoxShadow(x, y, width, height, offsetX, offsetY, blur int, color Color)

	// Theme returns the application-wide theme for this rendering context.
	Theme() Theme
}

// BasePublicWidget provides default implementations for the PublicWidget interface.
//
// Embed BasePublicWidget in custom widget types to get default event handling
// and bounds management. Override Draw() to provide custom rendering.
//
// Example:
//
//	type ColoredPanel struct {
//	    wayne.BasePublicWidget
//	    Color wayne.Color
//	}
//
//	func (p *ColoredPanel) Draw(c wayne.Canvas) {
//	    x, y := p.Position()
//	    w, h := p.Bounds()
//	    c.FillRect(x, y, w, h, p.Color)
//	}
type BasePublicWidget struct {
	x, y          int
	width, height int
	children      []PublicWidget
	visible       bool
	focused       bool

	onEvent func(Event) bool
}

// NewBasePublicWidget creates a BasePublicWidget with the given pixel dimensions.
func NewBasePublicWidget(width, height int) BasePublicWidget {
	return BasePublicWidget{
		width:   width,
		height:  height,
		visible: true,
	}
}

// Bounds returns the current pixel dimensions of the widget.
func (w *BasePublicWidget) Bounds() (width, height int) {
	return w.width, w.height
}

// Position returns the current position of the widget in pixels.
func (w *BasePublicWidget) Position() (x, y int) {
	return w.x, w.y
}

// HandleEvent processes an event. The default implementation invokes the
// registered event handler if one is set, otherwise returns false.
func (w *BasePublicWidget) HandleEvent(evt Event) bool {
	if w.onEvent != nil {
		return w.onEvent(evt)
	}
	return false
}

// Draw is a no-op. Override this in concrete widget types.
func (w *BasePublicWidget) Draw(c Canvas) {}

// Add appends a child widget.
func (w *BasePublicWidget) Add(child PublicWidget) {
	w.children = append(w.children, child)
}

// Children returns the list of child widgets.
func (w *BasePublicWidget) Children() []PublicWidget {
	return w.children
}

// SetBounds updates the widget's pixel dimensions and position.
func (w *BasePublicWidget) SetBounds(x, y, width, height int) {
	w.x = x
	w.y = y
	w.width = width
	w.height = height
}

// SetVisible controls whether the widget participates in layout and rendering.
func (w *BasePublicWidget) SetVisible(visible bool) {
	w.visible = visible
}

// IsVisible returns true if the widget is visible.
func (w *BasePublicWidget) IsVisible() bool {
	return w.visible
}

// OnEvent registers a callback to handle events for this widget.
func (w *BasePublicWidget) OnEvent(handler func(Event) bool) {
	w.onEvent = handler
}

// SetFocused sets the widget's keyboard focus state.
func (w *BasePublicWidget) SetFocused(focused bool) {
	w.focused = focused
}

// IsFocused returns true if the widget currently has keyboard focus.
func (w *BasePublicWidget) IsFocused() bool {
	return w.focused
}
