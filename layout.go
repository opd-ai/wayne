//go:build windows || darwin || android || ios

package wayne

// Size represents percentage-based dimensions for a widget.
//
// Width and Height are specified as percentages (0-100) of the parent container.
//
// Example:
//
//	sidebar := wayne.NewPanel(wayne.Size{Width: 25, Height: 100}) // 25% wide, full height
//	content := wayne.NewPanel(wayne.Size{Width: 75, Height: 100}) // 75% wide, full height
type Size struct {
	Width  float64 // Width as percentage of parent (0-100)
	Height float64 // Height as percentage of parent (0-100)
}

// FlowDirection controls how a container arranges its child widgets.
type FlowDirection int

const (
	// FlowRow arranges children horizontally, left to right.
	FlowRow FlowDirection = iota

	// FlowColumn arranges children vertically, top to bottom.
	FlowColumn
)

// Align specifies alignment on the cross axis of a container.
type Align int

const (
	// AlignStart aligns children to the start of the cross axis.
	AlignStart Align = iota
	// AlignCenter centers children on the cross axis.
	AlignCenter
	// AlignEnd aligns children to the end of the cross axis.
	AlignEnd
	// AlignStretch stretches children to fill the cross axis.
	AlignStretch
)

// sizeHinter is an internal interface for querying a widget's preferred size.
type sizeHinter interface {
	sizeHint() Size
}

// resolveTree resolves the layout tree recursively.
// It sets the pixel bounds of w and all its descendants.
func resolveTree(w PublicWidget, x, y, pixelW, pixelH int) {
	if bw, ok := w.(interface{ SetBounds(int, int, int, int) }); ok {
		bw.SetBounds(x, y, pixelW, pixelH)
	}

	switch c := w.(type) {
	case *Panel:
		c.resolveChildren(x, y, pixelW, pixelH)
	case *Row:
		c.resolveChildren(x, y, pixelW, pixelH)
	case *Column:
		c.resolveChildren(x, y, pixelW, pixelH)
	case *ScrollView:
		c.resolveChildren(x, y, pixelW, pixelH)
	}
}

// computeChildPixelSize computes the pixel size of a child widget relative to the parent.
func computeChildPixelSize(child PublicWidget, parentW, parentH int) (int, int) {
	if s, ok := child.(sizeHinter); ok {
		hint := s.sizeHint()
		w := clampPercent(hint.Width)
		h := clampPercent(hint.Height)
		return max(0, int(w*float64(parentW)/100.0)),
			max(0, int(h*float64(parentH)/100.0))
	}
	return parentW, parentH
}

func clampPercent(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return v
}

// Panel is a styled rectangular container that holds child widgets.
//
// Panel supports percentage-based sizing and automatic layout. It can be used
// as a building block for complex UIs by nesting panels.
//
// Example:
//
//	panel := wayne.NewPanel(wayne.Size{Width: 50, Height: 100})
//	panel.SetFlowDirection(wayne.FlowColumn)
//	panel.SetPadding(10)
//	panel.SetGap(5)
//	panel.Add(header)
//	panel.Add(content)
type Panel struct {
	BasePublicWidget

	size          Size
	flowDir       FlowDirection
	padding       int
	gap           int
	align         Align
	styleOverride *StyleOverride
	theme         *Theme
}

// NewPanel creates a new Panel with percentage-based dimensions.
func NewPanel(size Size) *Panel {
	return &Panel{
		BasePublicWidget: NewBasePublicWidget(0, 0),
		size:             size,
		flowDir:          FlowColumn,
		padding:          0,
		gap:              0,
	}
}

func (p *Panel) sizeHint() Size { return p.size }

// resolveChildren lays out this panel's children given the panel's pixel bounds.
func (p *Panel) resolveChildren(parentX, parentY, parentW, parentH int) {
	padding := p.padding
	gap := p.gap

	contentX := parentX + padding
	contentY := parentY + padding
	contentW := parentW - 2*padding
	contentH := parentH - 2*padding
	if contentW < 0 {
		contentW = 0
	}
	if contentH < 0 {
		contentH = 0
	}

	cursorX := contentX
	cursorY := contentY

	for _, child := range p.children {
		childW, childH := computeChildPixelSize(child, contentW, contentH)

		resolveTree(child, cursorX, cursorY, childW, childH)

		if p.flowDir == FlowRow {
			cursorX += childW + gap
		} else {
			cursorY += childH + gap
		}
	}
}

// Add appends a child widget to this panel.
func (p *Panel) Add(child PublicWidget) {
	p.children = append(p.children, child)
}

// Children returns this panel's child widgets.
func (p *Panel) Children() []PublicWidget {
	return p.children
}

// Bounds returns the resolved pixel dimensions.
func (p *Panel) Bounds() (width, height int) {
	return p.width, p.height
}

// HandleEvent passes events to children and returns true if any child consumed it.
func (p *Panel) HandleEvent(evt Event) bool {
	for _, child := range p.children {
		if child.HandleEvent(evt) {
			return true
		}
	}
	return false
}

// Draw renders the panel background and all children.
func (p *Panel) Draw(c Canvas) {
	theme := DefaultDark()
	if p.theme != nil {
		theme = *p.theme
	}
	effective := theme
	if p.styleOverride != nil {
		effective = p.styleOverride.applyToTheme(theme)
	}

	x, y := p.Position()
	w, h := p.Bounds()
	if w > 0 && h > 0 {
		c.FillRoundedRect(x, y, w, h, effective.BorderRadius, effective.Background)
		if effective.BorderWidth > 0 {
			c.DrawLine(x, y, x+w, y, effective.Border, effective.BorderWidth)
			c.DrawLine(x+w, y, x+w, y+h, effective.Border, effective.BorderWidth)
			c.DrawLine(x, y+h, x+w, y+h, effective.Border, effective.BorderWidth)
			c.DrawLine(x, y, x, y+h, effective.Border, effective.BorderWidth)
		}
	}

	for _, child := range p.children {
		child.Draw(c)
	}
}

// SetFlowDirection sets how this panel arranges its children.
func (p *Panel) SetFlowDirection(dir FlowDirection) {
	p.flowDir = dir
}

// FlowDirection returns the current flow direction.
func (p *Panel) FlowDirection() FlowDirection {
	return p.flowDir
}

// SetPadding sets the padding (in pixels) around the content area.
func (p *Panel) SetPadding(pixels int) {
	p.padding = pixels
}

// SetGap sets the spacing (in pixels) between child widgets.
func (p *Panel) SetGap(pixels int) {
	p.gap = pixels
}

// SetAlign sets the cross-axis alignment for children.
func (p *Panel) SetAlign(align Align) {
	p.align = align
}

// SetVisible controls whether this panel is drawn.
func (p *Panel) SetVisible(visible bool) {
	p.BasePublicWidget.SetVisible(visible)
}

// Visible reports whether the panel is currently visible.
func (p *Panel) Visible() bool {
	return p.BasePublicWidget.IsVisible()
}

// SetStyle applies a style override to this panel.
func (p *Panel) SetStyle(override StyleOverride) {
	p.styleOverride = &override
}

// SetTheme applies a theme to this panel and all its children.
func (p *Panel) SetTheme(theme Theme) {
	p.theme = &theme
	for _, child := range p.children {
		if panel, ok := child.(*Panel); ok {
			panel.SetTheme(theme)
		}
	}
}

// Row is a convenience container that arranges children horizontally.
//
// Row is equivalent to a Panel with FlowDirection set to FlowRow.
type Row struct {
	*Panel
}

// NewRow creates a new horizontal container.
func NewRow() *Row {
	panel := NewPanel(Size{Width: 100, Height: 100})
	panel.SetFlowDirection(FlowRow)
	return &Row{Panel: panel}
}

// Column is a convenience container that arranges children vertically.
//
// Column is equivalent to a Panel with FlowDirection set to FlowColumn.
type Column struct {
	*Panel
}

// NewColumn creates a new vertical container.
func NewColumn() *Column {
	panel := NewPanel(Size{Width: 100, Height: 100})
	panel.SetFlowDirection(FlowColumn)
	return &Column{Panel: panel}
}

// Stack is a layering container that places children on top of each other.
type Stack struct {
	*Panel
}

// NewStack creates a new layering container.
func NewStack() *Stack {
	panel := NewPanel(Size{Width: 100, Height: 100})
	return &Stack{Panel: panel}
}

// Grid is a fixed-column grid container.
type Grid struct {
	*Panel
	columns int
}

// NewGrid creates a new grid container with the specified number of columns.
func NewGrid(columns int) *Grid {
	if columns < 1 {
		columns = 1
	}
	panel := NewPanel(Size{Width: 100, Height: 100})
	return &Grid{Panel: panel, columns: columns}
}

// Columns returns the number of columns in the grid.
func (g *Grid) Columns() int { return g.columns }

// SetColumns changes the number of columns in the grid.
func (g *Grid) SetColumns(columns int) {
	if columns < 1 {
		columns = 1
	}
	g.columns = columns
}
