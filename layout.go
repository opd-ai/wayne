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

// validateSize checks that a Size has valid percentage values (0-100).
// It panics if Width or Height are negative or exceed 100.
func validateSize(size Size) {
	if size.Width < 0 || size.Width > 100 {
		panic("wayne: Size.Width must be between 0 and 100 (percentage of parent)")
	}
	if size.Height < 0 || size.Height > 100 {
		panic("wayne: Size.Height must be between 0 and 100 (percentage of parent)")
	}
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
	case *Grid:
		c.resolveChildren(x, y, pixelW, pixelH)
	case *Stack:
		c.resolveChildren(x, y, pixelW, pixelH)
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
	validateSize(size)
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
	contentX, contentY, contentW, contentH := p.computeContentArea(parentX, parentY, parentW, parentH)
	cursorX, cursorY := contentX, contentY

	for _, child := range p.children {
		childW, childH := computeChildPixelSize(child, contentW, contentH)
		alignedX, alignedY, alignedW, alignedH := p.alignChild(cursorX, cursorY, contentX, contentY, contentW, contentH, childW, childH)
		resolveTree(child, alignedX, alignedY, alignedW, alignedH)
		cursorX, cursorY = p.advanceCursor(cursorX, cursorY, alignedW, alignedH)
	}
}

// computeContentArea calculates the content area after applying padding.
func (p *Panel) computeContentArea(parentX, parentY, parentW, parentH int) (x, y, w, h int) {
	x = parentX + p.padding
	y = parentY + p.padding
	w = max(0, parentW-2*p.padding)
	h = max(0, parentH-2*p.padding)
	return x, y, w, h
}

// alignChild applies cross-axis alignment to a child widget.
func (p *Panel) alignChild(cursorX, cursorY, contentX, contentY, contentW, contentH, childW, childH int) (x, y, w, h int) {
	x, y, w, h = cursorX, cursorY, childW, childH

	if p.flowDir == FlowRow {
		x, y, h = p.alignCrossAxisVertical(cursorY, contentY, contentH, childH, p.align)
		x = cursorX
	} else {
		x, y, w = p.alignCrossAxisHorizontal(cursorX, contentX, contentW, childW, p.align)
		y = cursorY
	}
	return x, y, w, h
}

// alignCrossAxisVertical computes vertical alignment for horizontal flow.
func (p *Panel) alignCrossAxisVertical(cursorY, contentY, contentH, childH int, align Align) (x, y, h int) {
	switch align {
	case AlignCenter:
		return 0, contentY + (contentH-childH)/2, childH
	case AlignEnd:
		return 0, contentY + contentH - childH, childH
	case AlignStretch:
		return 0, contentY, contentH
	default:
		return 0, cursorY, childH
	}
}

// alignCrossAxisHorizontal computes horizontal alignment for vertical flow.
func (p *Panel) alignCrossAxisHorizontal(cursorX, contentX, contentW, childW int, align Align) (x, y, w int) {
	switch align {
	case AlignCenter:
		return contentX + (contentW-childW)/2, 0, childW
	case AlignEnd:
		return contentX + contentW - childW, 0, childW
	case AlignStretch:
		return contentX, 0, contentW
	default:
		return cursorX, 0, childW
	}
}

// advanceCursor moves the layout cursor after placing a child.
func (p *Panel) advanceCursor(cursorX, cursorY, childW, childH int) (newX, newY int) {
	if p.flowDir == FlowRow {
		return cursorX + childW + p.gap, cursorY
	}
	return cursorX, cursorY + childH + p.gap
}

// Add appends a child widget to this panel.
func (p *Panel) Add(child PublicWidget) {
	p.children = append(p.children, child)
}

// Children returns this panel's child widgets.
func (p *Panel) Children() []PublicWidget {
	return p.children
}

// bounds returns the resolved pixel dimensions.
// This method is unexported for gomobile compatibility.
// Use Width() and Height() for interface-level access.
func (p *Panel) bounds() (width, height int) {
	return p.width, p.height
}

// Width returns the resolved width in pixels.
func (p *Panel) Width() int {
	return p.width
}

// Height returns the resolved height in pixels.
func (p *Panel) Height() int {
	return p.height
}

// HandleEvent passes events to children and returns true if any child consumed it.
// For pointer and touch events, only forwards to children whose bounds contain the event coordinates.
func (p *Panel) HandleEvent(evt Event) bool {
	switch evt.Type() {
	case EventTypePointer:
		return p.handleSpatialEvent(evt, evt.(*PointerEvent).X(), evt.(*PointerEvent).Y())
	case EventTypeTouch:
		return p.handleSpatialEvent(evt, evt.(*TouchEvent).X(), evt.(*TouchEvent).Y())
	default:
		return p.broadcastEvent(evt)
	}
}

// handleSpatialEvent forwards spatial events to the topmost child containing the coordinates.
func (p *Panel) handleSpatialEvent(evt Event, evtX, evtY float64) bool {
	for i := len(p.children) - 1; i >= 0; i-- {
		child := p.children[i]
		if p.childContainsPoint(child, evtX, evtY) {
			if child.HandleEvent(evt) {
				return true
			}
			break // Only forward to first matching child (z-order topmost)
		}
	}
	return false
}

// childContainsPoint checks if a child widget contains the given point.
func (p *Panel) childContainsPoint(child PublicWidget, px, py float64) bool {
	positioner, ok := child.(interface{ position() (int, int) })
	if !ok {
		return false
	}
	x, y := positioner.position()
	w, h := child.Width(), child.Height()
	return contains(x, y, w, h, px, py)
}

// broadcastEvent forwards non-spatial events to all children.
func (p *Panel) broadcastEvent(evt Event) bool {
	for _, child := range p.children {
		if child.HandleEvent(evt) {
			return true
		}
	}
	return false
}

// Draw renders the panel background and all children.
func (p *Panel) Draw(c Canvas) {
	effective := p.resolveEffectiveTheme(c)
	x, y := p.position()
	w, h := p.bounds()

	if w > 0 && h > 0 {
		p.drawBackground(c, x, y, w, h, effective)
	}

	for _, child := range p.children {
		child.Draw(c)
	}
}

// resolveEffectiveTheme returns the theme to use for rendering.
func (p *Panel) resolveEffectiveTheme(c Canvas) Theme {
	theme := c.Theme()
	if p.theme != nil {
		theme = *p.theme
	}
	if p.styleOverride != nil {
		return p.styleOverride.applyToTheme(theme)
	}
	return theme
}

// drawBackground renders the panel's background and border.
func (p *Panel) drawBackground(c Canvas, x, y, w, h int, theme Theme) {
	scale := c.Scale()
	borderRadius := int(float64(theme.BorderRadius) * scale)
	borderWidth := int(float64(theme.BorderWidth) * scale)

	c.FillRoundedRect(x, y, w, h, borderRadius, theme.Background)
	if borderWidth > 0 {
		p.drawBorder(c, x, y, w, h, theme.Border, borderWidth)
	}
}

// drawBorder renders the panel's border lines.
func (p *Panel) drawBorder(c Canvas, x, y, w, h int, color Color, width int) {
	c.DrawLine(x, y, x+w, y, color, width)
	c.DrawLine(x+w, y, x+w, y+h, color, width)
	c.DrawLine(x, y+h, x+w, y+h, color, width)
	c.DrawLine(x, y, x, y+h, color, width)
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
		if themeable, ok := child.(Themeable); ok {
			themeable.SetTheme(theme)
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

// resolveChildren lays out stack children at the same position (z-order stacking).
func (s *Stack) resolveChildren(parentX, parentY, parentW, parentH int) {
	padding := s.padding

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

	// All children are placed at the same position with full parent dimensions
	for _, child := range s.children {
		childW, childH := computeChildPixelSize(child, contentW, contentH)

		// Apply alignment within the stack area
		alignedX := contentX
		alignedY := contentY

		switch s.align {
		case AlignCenter:
			alignedX = contentX + (contentW-childW)/2
			alignedY = contentY + (contentH-childH)/2
		case AlignEnd:
			alignedX = contentX + contentW - childW
			alignedY = contentY + contentH - childH
		case AlignStretch:
			childW = contentW
			childH = contentH
		}

		resolveTree(child, alignedX, alignedY, childW, childH)
	}
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

// resolveChildren lays out grid children in rows and columns.
func (g *Grid) resolveChildren(parentX, parentY, parentW, parentH int) {
	numChildren := len(g.children)
	if numChildren == 0 {
		return
	}

	contentX, contentY, contentW, contentH := g.Panel.computeContentArea(parentX, parentY, parentW, parentH)
	cellW, cellH := g.computeCellDimensions(contentW, contentH, numChildren)

	for i, child := range g.children {
		g.layoutGridChild(child, i, contentX, contentY, cellW, cellH)
	}
}

// computeCellDimensions calculates the cell size based on grid dimensions.
func (g *Grid) computeCellDimensions(contentW, contentH, numChildren int) (cellW, cellH int) {
	rows := (numChildren + g.columns - 1) / g.columns
	totalGapW := g.gap * (g.columns - 1)
	totalGapH := g.gap * (rows - 1)
	cellW = max(0, (contentW-totalGapW)/g.columns)
	cellH = max(0, (contentH-totalGapH)/rows)
	return cellW, cellH
}

// layoutGridChild positions and aligns a single grid cell.
func (g *Grid) layoutGridChild(child PublicWidget, index, contentX, contentY, cellW, cellH int) {
	col := index % g.columns
	row := index / g.columns

	childX := contentX + col*(cellW+g.gap)
	childY := contentY + row*(cellH+g.gap)
	childW, childH := computeChildPixelSize(child, cellW, cellH)

	alignedX, alignedY, alignedW, alignedH := g.alignGridChild(childX, childY, cellW, cellH, childW, childH)
	resolveTree(child, alignedX, alignedY, alignedW, alignedH)
}

// alignGridChild applies alignment within a grid cell.
func (g *Grid) alignGridChild(cellX, cellY, cellW, cellH, childW, childH int) (x, y, w, h int) {
	switch g.align {
	case AlignCenter:
		return cellX + (cellW-childW)/2, cellY + (cellH-childH)/2, childW, childH
	case AlignEnd:
		return cellX + cellW - childW, cellY + cellH - childH, childW, childH
	case AlignStretch:
		return cellX, cellY, cellW, cellH
	default:
		return cellX, cellY, childW, childH
	}
}
