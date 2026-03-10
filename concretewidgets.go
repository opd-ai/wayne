//go:build windows || darwin || android || ios

package wayne

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Button is a clickable button widget with text and onClick callback.
//
// Example:
//
//	btn := wayne.NewButton("Submit", wayne.Size{Width: 30, Height: 8})
//	btn.OnClick(func() {
//	    fmt.Println("Button clicked!")
//	})
//	panel.Add(btn)
type Button struct {
	BasePublicWidget

	label   string
	size    Size
	onClick func()
	enabled bool
	hovered bool
	pressed bool
	theme   *Theme
}

// NewButton creates a new button with the specified label and percentage-based size.
func NewButton(label string, size Size) *Button {
	validateSize(size)
	return &Button{
		BasePublicWidget: NewBasePublicWidget(0, 0),
		label:            label,
		size:             size,
		enabled:          true,
	}
}

func (b *Button) sizeHint() Size { return b.size }

// OnClick registers a callback to be invoked when the button is clicked.
func (b *Button) OnClick(handler func()) {
	b.onClick = handler
}

// SetLabel changes the button's displayed text.
func (b *Button) SetLabel(text string) {
	b.label = text
}

// SetText changes the button's displayed text (alias for SetLabel).
//
// Deprecated: Use SetLabel instead. SetText is retained for compatibility but will be removed in v2.0.0.
func (b *Button) SetText(text string) {
	b.label = text
}

// Text returns the button's current label.
func (b *Button) Text() string {
	return b.label
}

// SetEnabled enables or disables the button.
func (b *Button) SetEnabled(enabled bool) {
	b.enabled = enabled
}

// SetTheme applies a theme to this button.
func (b *Button) SetTheme(theme Theme) {
	b.theme = &theme
}

// CanTakeFocus returns true if the button can currently receive keyboard focus.
func (b *Button) CanTakeFocus() bool {
	return b.enabled
}

// HandleEvent processes pointer and keyboard events for button interaction.
func (b *Button) HandleEvent(evt Event) bool {
	// Handle keyboard events for focused buttons
	if ke, ok := evt.(*KeyEvent); ok {
		if b.IsFocused() && b.enabled && ke.EventType() == KeyPress {
			if ke.Key() == KeyReturn || ke.Key() == KeySpace {
				if b.onClick != nil {
					b.onClick()
				}
				return true
			}
		}
		return false
	}

	// Handle pointer events
	pe, ok := evt.(*PointerEvent)
	if !ok {
		return false
	}

	switch pe.EventType() {
	case PointerEnter:
		b.hovered = true
		return true
	case PointerLeave:
		b.hovered = false
		b.pressed = false
		return true
	case PointerButtonPress:
		if pe.Button() == PointerButtonLeft && b.enabled {
			x, y := b.Position()
			w, h := b.Bounds()
			if !contains(x, y, w, h, pe.X(), pe.Y()) {
				return false
			}
			b.pressed = true
			return true
		}
	case PointerButtonRelease:
		if pe.Button() == PointerButtonLeft && b.pressed && b.enabled {
			x, y := b.Position()
			w, h := b.Bounds()
			if !contains(x, y, w, h, pe.X(), pe.Y()) {
				return false
			}
			b.pressed = false
			if b.onClick != nil {
				b.onClick()
			}
			return true
		}
		b.pressed = false
	}
	return false
}

// Draw renders the button to the canvas.
func (b *Button) Draw(c Canvas) {
	x, y := b.Position()
	w, h := b.Bounds()
	if w <= 0 || h <= 0 {
		return
	}

	theme := c.Theme()
	if b.theme != nil {
		theme = *b.theme
	}

	bgColor := theme.Accent
	if !b.enabled {
		bgColor = theme.Border
	} else if b.pressed {
		bgColor = RGB(
			uint8(max(0, int(theme.Accent.R)-20)),
			uint8(max(0, int(theme.Accent.G)-20)),
			uint8(max(0, int(theme.Accent.B)-20)),
		)
	} else if b.hovered {
		bgColor = RGB(
			uint8(min(255, int(theme.Accent.R)+20)),
			uint8(min(255, int(theme.Accent.G)+20)),
			uint8(min(255, int(theme.Accent.B)+20)),
		)
	}

	c.FillRoundedRect(x, y, w, h, theme.BorderRadius, bgColor)

	// Draw focus indicator (border)
	if b.IsFocused() {
		focusColor := theme.Accent
		// Draw a 2-pixel border around the button
		c.DrawLine(x, y, x+w, y, focusColor, 2)     // top
		c.DrawLine(x, y+h, x+w, y+h, focusColor, 2) // bottom
		c.DrawLine(x, y, x, y+h, focusColor, 2)     // left
		c.DrawLine(x+w, y, x+w, y+h, focusColor, 2) // right
	}

	if b.label != "" {
		// Center text (approximate: basicfont is 7px wide per char, 13px tall)
		textW := len(b.label) * 7
		textX := x + (w-textW)/2
		textY := y + (h-13)/2
		if textY < y {
			textY = y + 2
		}
		c.DrawText(b.label, textX, textY, nil, theme.Foreground)
	}
}

// Label is a static text display widget.
//
// Example:
//
//	label := wayne.NewLabel("Welcome", wayne.Size{Width: 100, Height: 5})
//	panel.Add(label)
type Label struct {
	BasePublicWidget

	text      string
	size      Size
	textColor *Color
	fontSize  int
	theme     *Theme
}

// NewLabel creates a new label with the specified text and percentage-based size.
func NewLabel(text string, size Size) *Label {
	validateSize(size)
	return &Label{
		BasePublicWidget: NewBasePublicWidget(0, 0),
		text:             text,
		size:             size,
	}
}

func (l *Label) sizeHint() Size { return l.size }

// SetText changes the label's displayed text.
func (l *Label) SetText(text string) {
	l.text = text
}

// Text returns the label's current text.
func (l *Label) Text() string {
	return l.text
}

// SetTextColor sets the color of the label's text.
func (l *Label) SetTextColor(color Color) {
	l.textColor = &color
}

// SetFontSize sets the font size in pixels.
func (l *Label) SetFontSize(size int) {
	l.fontSize = size
}

// SetTheme applies a theme to this label.
func (l *Label) SetTheme(theme Theme) {
	l.theme = &theme
}

// HandleEvent does nothing for labels (they don't respond to input).
func (l *Label) HandleEvent(_ Event) bool { return false }

// Draw renders the label text to the canvas.
func (l *Label) Draw(c Canvas) {
	if l.text == "" {
		return
	}

	theme := c.Theme()
	if l.theme != nil {
		theme = *l.theme
	}

	col := theme.Foreground
	if l.textColor != nil {
		col = *l.textColor
	}

	x, y := l.Position()
	w, h := l.Bounds()
	_ = h
	if w <= 0 {
		return
	}

	c.DrawText(l.text, x, y, nil, col)
}

// TextInput is a single-line editable text field.
//
// Example:
//
//	input := wayne.NewTextInput("Enter name...", wayne.Size{Width: 50, Height: 6})
//	input.OnChange(func(text string) {
//	    fmt.Println("Input:", text)
//	})
//	panel.Add(input)
type TextInput struct {
	BasePublicWidget

	text        string
	placeholder string
	size        Size
	onChange    func(string)
	cursorPos   int
	theme       *Theme
}

// NewTextInput creates a new text input field with placeholder text and size.
func NewTextInput(placeholder string, size Size) *TextInput {
	validateSize(size)
	return &TextInput{
		BasePublicWidget: NewBasePublicWidget(0, 0),
		placeholder:      placeholder,
		size:             size,
	}
}

func (t *TextInput) sizeHint() Size { return t.size }

// OnChange registers a callback invoked when the text changes.
func (t *TextInput) OnChange(handler func(string)) {
	t.onChange = handler
}

// SetText changes the input's text content.
func (t *TextInput) SetText(text string) {
	t.text = text
	t.cursorPos = len([]rune(text))
	if t.onChange != nil {
		t.onChange(text)
	}
}

// Text returns the current text content.
func (t *TextInput) Text() string {
	return t.text
}

// SetPlaceholder sets the placeholder text shown when the input is empty.
func (t *TextInput) SetPlaceholder(placeholder string) {
	t.placeholder = placeholder
}

// SetTheme applies a theme to this text input.
func (t *TextInput) SetTheme(theme Theme) {
	t.theme = &theme
}

// CanTakeFocus returns true (text inputs can always receive focus).
func (t *TextInput) CanTakeFocus() bool {
	return true
}

// HandleEvent processes keyboard and pointer events for text input.
func (t *TextInput) HandleEvent(evt Event) bool {
	switch e := evt.(type) {
	case *PointerEvent:
		return t.handlePointerEvent(e)
	case *KeyEvent:
		if !t.IsFocused() {
			return false
		}
		if e.EventType() == KeyPress || e.EventType() == KeyRepeat {
			return t.handleKeyEvent(e)
		}
	}
	return false
}

func (t *TextInput) handlePointerEvent(pe *PointerEvent) bool {
	if pe.EventType() == PointerButtonPress {
		x, y := t.Position()
		w, h := t.Bounds()
		if contains(x, y, w, h, pe.X(), pe.Y()) {
			t.SetFocused(true)
			return true
		}
		t.SetFocused(false)
	}
	return false
}

func (t *TextInput) handleKeyEvent(ke *KeyEvent) bool {
	switch ke.Key() {
	case KeyBackspace:
		return t.handleBackspace()
	case KeyLeft, KeyRight, KeyHome, KeyEnd:
		return t.handleCursorMovement(ke.Key())
	default:
		if r := ke.Rune(); r != 0 && r >= 0x20 {
			return t.handleTextInput(r)
		}
	}
	return false
}

func (t *TextInput) handleBackspace() bool {
	runes := []rune(t.text)
	if t.cursorPos > 0 {
		runes = append(runes[:t.cursorPos-1], runes[t.cursorPos:]...)
		t.cursorPos--
		t.text = string(runes)
		if t.onChange != nil {
			t.onChange(t.text)
		}
	}
	return true
}

func (t *TextInput) handleCursorMovement(key Key) bool {
	switch key {
	case KeyLeft:
		if t.cursorPos > 0 {
			t.cursorPos--
		}
	case KeyRight:
		if t.cursorPos < len([]rune(t.text)) {
			t.cursorPos++
		}
	case KeyHome:
		t.cursorPos = 0
	case KeyEnd:
		t.cursorPos = len([]rune(t.text))
	}
	return true
}

func (t *TextInput) handleTextInput(r rune) bool {
	runes := []rune(t.text)
	runes = append(runes[:t.cursorPos], append([]rune{r}, runes[t.cursorPos:]...)...)
	t.text = string(runes)
	t.cursorPos++
	if t.onChange != nil {
		t.onChange(t.text)
	}
	return true
}

// Draw renders the text input to the canvas.
func (t *TextInput) Draw(c Canvas) {
	x, y := t.Position()
	w, h := t.Bounds()
	if w <= 0 || h <= 0 {
		return
	}

	theme := c.Theme()
	if t.theme != nil {
		theme = *t.theme
	}

	// Background
	bg := theme.Background
	if t.IsFocused() {
		bg = RGB(
			uint8(min(255, int(bg.R)+15)),
			uint8(min(255, int(bg.G)+15)),
			uint8(min(255, int(bg.B)+15)),
		)
	}
	c.FillRoundedRect(x, y, w, h, theme.BorderRadius, bg)

	// Border (brighter when focused)
	borderColor := theme.Border
	if t.IsFocused() {
		borderColor = theme.Accent
	}
	c.DrawLine(x, y, x+w, y, borderColor, theme.BorderWidth)
	c.DrawLine(x+w, y, x+w, y+h, borderColor, theme.BorderWidth)
	c.DrawLine(x, y+h, x+w, y+h, borderColor, theme.BorderWidth)
	c.DrawLine(x, y, x, y+h, borderColor, theme.BorderWidth)

	// Text
	displayText := t.text
	textColor := theme.Foreground
	if displayText == "" && t.placeholder != "" {
		displayText = t.placeholder
		textColor = theme.Border
	}

	if displayText != "" {
		c.DrawText(displayText, x+4, y+2, nil, textColor)
	}

	// Cursor
	if t.IsFocused() {
		cursorX := x + 4 + t.cursorPos*7
		c.DrawLine(cursorX, y+2, cursorX, y+h-2, theme.Foreground, 1)
	}
}

// ScrollView is a scrollable container for overflow content.
//
// Example:
//
//	scroll := wayne.NewScrollView(wayne.Size{Width: 100, Height: 80})
//	for i := 0; i < 50; i++ {
//	    scroll.Add(wayne.NewLabel(fmt.Sprintf("Item %d", i), wayne.Size{Width: 100, Height: 5}))
//	}
type ScrollView struct {
	BasePublicWidget

	size     Size
	scrollY  int
	onScroll func(offset int)
	theme    *Theme
}

// NewScrollView creates a new scrollable container.
func NewScrollView(size Size) *ScrollView {
	validateSize(size)
	return &ScrollView{
		BasePublicWidget: NewBasePublicWidget(0, 0),
		size:             size,
	}
}

func (s *ScrollView) sizeHint() Size { return s.size }

// resolveChildren lays out the scroll view's children.
func (s *ScrollView) resolveChildren(parentX, parentY, parentW, parentH int) {
	cursorY := parentY - s.scrollY
	for _, child := range s.children {
		childW, childH := computeChildPixelSize(child, parentW, parentH)
		resolveTree(child, parentX, cursorY, childW, childH)
		cursorY += childH
	}
}

// OnScroll registers a callback invoked when the scroll offset changes.
func (s *ScrollView) OnScroll(handler func(offset int)) {
	s.onScroll = handler
}

// SetScrollOffset sets the current scroll position in pixels.
func (s *ScrollView) SetScrollOffset(offset int) {
	if offset < 0 {
		offset = 0
	}
	s.scrollY = offset
	if s.onScroll != nil {
		s.onScroll(offset)
	}
}

// ScrollOffset returns the current scroll position in pixels.
func (s *ScrollView) ScrollOffset() int {
	return s.scrollY
}

// Add appends a child widget to the scroll view.
func (s *ScrollView) Add(child PublicWidget) {
	s.children = append(s.children, child)
}

// Children returns the scroll view's children.
func (s *ScrollView) Children() []PublicWidget {
	return s.children
}

// HandleEvent processes scroll events.
func (s *ScrollView) HandleEvent(evt Event) bool {
	if pe, ok := evt.(*PointerEvent); ok && pe.EventType() == PointerScroll {
		x, y := s.Position()
		w, h := s.Bounds()
		if contains(x, y, w, h, pe.X(), pe.Y()) {
			s.SetScrollOffset(s.scrollY + int(pe.Value()*20))
			return true
		}
	}
	// Propagate to children
	for _, child := range s.children {
		if child.HandleEvent(evt) {
			return true
		}
	}
	return false
}

// Draw renders the scroll view and its visible children with clipping.
func (s *ScrollView) Draw(c Canvas) {
	x, y := s.Position()
	w, h := s.Bounds()
	if w <= 0 || h <= 0 {
		return
	}

	theme := c.Theme()
	if s.theme != nil {
		theme = *s.theme
	}

	c.FillRect(x, y, w, h, theme.Background)

	// Create an offscreen buffer for the viewport to implement clipping
	viewport := ebiten.NewImage(w, h)
	defer viewport.Deallocate()

	// Render children to the offscreen viewport
	viewportCanvas := newEbitenCanvas(viewport, theme)
	for _, child := range s.children {
		cx, cy := 0, 0
		if bp, ok := child.(interface{ Position() (int, int) }); ok {
			cx, cy = bp.Position()
		}
		cw, ch := child.Bounds()

		// Translate child position by scroll offset
		adjustedY := cy - s.scrollY

		// Only draw children that are visible within the viewport
		if adjustedY+ch >= 0 && adjustedY < h {
			// Temporarily adjust child position for viewport-relative rendering
			if bp, ok := child.(interface{ SetBounds(int, int, int, int) }); ok {
				bp.SetBounds(cx, adjustedY, cw, ch)
				child.Draw(viewportCanvas)
				// Restore original position
				bp.SetBounds(cx, cy, cw, ch)
			}
		}
	}

	// Draw the clipped viewport to the main canvas
	if ec, ok := c.(*ebitenCanvas); ok {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		ec.dst.DrawImage(viewport, op)
	}
}

// SetTheme applies a theme to this scroll view.
func (s *ScrollView) SetTheme(theme Theme) {
	s.theme = &theme
	for _, child := range s.children {
		if themeable, ok := child.(Themeable); ok {
			themeable.SetTheme(theme)
		}
	}
}

// ImageWidget displays an image resource.
//
// Example:
//
//	imageWidget := wayne.NewImageWidget(wayne.Size{Width: 20, Height: 20})
//	imageWidget.SetImage(img)
//	panel.Add(imageWidget)
type ImageWidget struct {
	BasePublicWidget

	img  *Image
	size Size
}

// NewImageWidget creates a new image display widget.
func NewImageWidget(size Size) *ImageWidget {
	validateSize(size)
	return &ImageWidget{
		BasePublicWidget: NewBasePublicWidget(0, 0),
		size:             size,
	}
}

func (iw *ImageWidget) sizeHint() Size { return iw.size }

// SetImage changes the displayed image.
func (iw *ImageWidget) SetImage(img *Image) {
	iw.img = img
}

// GetImage returns the currently displayed image.
func (iw *ImageWidget) GetImage() *Image {
	return iw.img
}

// HandleEvent does nothing for image widgets.
func (iw *ImageWidget) HandleEvent(_ Event) bool { return false }

// Draw renders the image to the canvas.
func (iw *ImageWidget) Draw(c Canvas) {
	if iw.img == nil {
		return
	}
	x, y := iw.Position()
	w, h := iw.Bounds()
	if w > 0 && h > 0 {
		c.DrawImage(iw.img, x, y, w, h)
	}
}

// Spacer is an invisible widget that consumes percentage space for layout.
//
// Example:
//
//	row := wayne.NewRow()
//	row.Add(wayne.NewButton("Left", wayne.Size{Width: 20, Height: 10}))
//	row.Add(wayne.NewSpacer(wayne.Size{Width: 60, Height: 10}))
//	row.Add(wayne.NewButton("Right", wayne.Size{Width: 20, Height: 10}))
type Spacer struct {
	BasePublicWidget
	size Size
}

// NewSpacer creates a new invisible spacer widget.
func NewSpacer(size Size) *Spacer {
	validateSize(size)
	return &Spacer{
		BasePublicWidget: NewBasePublicWidget(0, 0),
		size:             size,
	}
}

func (s *Spacer) sizeHint() Size { return s.size }

// HandleEvent does nothing for spacers.
func (s *Spacer) HandleEvent(_ Event) bool { return false }

// Draw does nothing (spacers are invisible).
func (s *Spacer) Draw(_ Canvas) {}

// contains checks if a point (px, py) is within a widget's bounds.
func contains(x, y, w, h int, px, py float64) bool {
	return px >= float64(x) && px < float64(x+w) &&
		py >= float64(y) && py < float64(y+h)
}
