//go:build windows || darwin || android || ios

package wayne

import (
	"testing"
)

func TestPanelLayout(t *testing.T) {
	panel := NewPanel(Size{Width: 100, Height: 100})

	if panel == nil {
		t.Fatal("NewPanel returned nil")
	}

	// Test adding children
	btn1 := NewButton("Button 1", Size{Width: 50, Height: 30})
	btn2 := NewButton("Button 2", Size{Width: 50, Height: 30})

	panel.Add(btn1)
	panel.Add(btn2)

	if len(panel.children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(panel.children))
	}
}

func TestRowLayoutPercentages(t *testing.T) {
	row := NewRow()

	if row == nil {
		t.Fatal("NewRow returned nil")
	}

	// Add widgets with percentage sizes
	btn1 := NewButton("50%", Size{Width: 50, Height: 100})
	btn2 := NewButton("50%", Size{Width: 50, Height: 100})

	row.Add(btn1)
	row.Add(btn2)

	if len(row.children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(row.children))
	}
}

func TestColumnLayoutPercentages(t *testing.T) {
	col := NewColumn()

	if col == nil {
		t.Fatal("NewColumn returned nil")
	}

	// Add widgets with percentage sizes
	btn1 := NewButton("50%", Size{Width: 100, Height: 50})
	btn2 := NewButton("50%", Size{Width: 100, Height: 50})

	col.Add(btn1)
	col.Add(btn2)

	if len(col.children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(col.children))
	}
}

func TestGridLayoutColumns(t *testing.T) {
	tests := []struct {
		name       string
		columns    int
		childCount int
	}{
		{
			name:       "2x2 grid",
			columns:    2,
			childCount: 4,
		},
		{
			name:       "3x3 grid",
			columns:    3,
			childCount: 9,
		},
		{
			name:       "4x2 grid",
			columns:    4,
			childCount: 8,
		},
		{
			name:       "incomplete grid",
			columns:    3,
			childCount: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grid := NewGrid(tt.columns)

			if grid == nil {
				t.Fatal("NewGrid returned nil")
			}

			if grid.columns != tt.columns {
				t.Errorf("Expected %d columns, got %d", tt.columns, grid.columns)
			}

			// Add children
			for i := 0; i < tt.childCount; i++ {
				btn := NewButton("Btn", Size{Width: 50, Height: 30})
				grid.Add(btn)
			}

			if len(grid.children) != tt.childCount {
				t.Errorf("Expected %d children, got %d", tt.childCount, len(grid.children))
			}

			// Set bounds and test
			grid.SetBounds(0, 0, 400, 400)
		})
	}
}

func TestPanelSetPadding(t *testing.T) {
	panel := NewPanel(Size{Width: 100, Height: 100})

	panel.SetPadding(10)

	if panel.padding != 10 {
		t.Errorf("Expected padding 10, got %d", panel.padding)
	}
}

func TestPanelSetGap(t *testing.T) {
	panel := NewPanel(Size{Width: 100, Height: 100})

	panel.SetGap(5)

	if panel.gap != 5 {
		t.Errorf("Expected gap 5, got %d", panel.gap)
	}
}

func TestPanelFlowRow(t *testing.T) {
	panel := NewPanel(Size{Width: 100, Height: 100})
	panel.SetFlowDirection(FlowRow)

	if panel.flowDir != FlowRow {
		t.Errorf("Expected flow FlowRow, got %v", panel.flowDir)
	}
}

func TestPanelFlowColumn(t *testing.T) {
	panel := NewPanel(Size{Width: 100, Height: 100})
	panel.SetFlowDirection(FlowColumn)

	if panel.flowDir != FlowColumn {
		t.Errorf("Expected flow FlowColumn, got %v", panel.flowDir)
	}
}

func TestPanelSetAlign(t *testing.T) {
	panel := NewPanel(Size{Width: 100, Height: 100})

	tests := []struct {
		name  string
		align Align
	}{
		{"align start", AlignStart},
		{"align center", AlignCenter},
		{"align end", AlignEnd},
		{"align stretch", AlignStretch},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			panel.SetAlign(tt.align)

			if panel.align != tt.align {
				t.Errorf("Expected align %v, got %v", tt.align, panel.align)
			}
		})
	}
}

func TestPanelHandleEvent(t *testing.T) {
	panel := NewPanel(Size{Width: 400, Height: 400})

	btn := NewButton("Click", Size{Width: 100, Height: 30})
	clicked := false
	btn.OnClick(func() {
		clicked = true
	})

	panel.Add(btn)
	panel.SetBounds(0, 0, 400, 400)
	btn.SetBounds(10, 10, 100, 30)

	// Send click event to panel
	pressEvt := NewPointerEvent(PointerButtonPress, 50, 20, PointerButtonLeft, 0, 0)
	panel.HandleEvent(pressEvt)

	releaseEvt := NewPointerEvent(PointerButtonRelease, 50, 20, PointerButtonLeft, 0, 0)
	panel.HandleEvent(releaseEvt)

	if !clicked {
		t.Error("Button click should be propagated through panel")
	}
}

func TestRowAdd(t *testing.T) {
	row := NewRow()

	btn := NewButton("Test", Size{Width: 100, Height: 30})
	row.Add(btn)

	if len(row.children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(row.children))
	}
}

func TestColumnAdd(t *testing.T) {
	col := NewColumn()

	label := NewLabel("Test", Size{Width: 200, Height: 40})
	col.Add(label)

	if len(col.children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(col.children))
	}
}

func TestGridAdd(t *testing.T) {
	grid := NewGrid(3)

	for i := 0; i < 6; i++ {
		btn := NewButton("Btn", Size{Width: 50, Height: 30})
		grid.Add(btn)
	}

	if len(grid.children) != 6 {
		t.Errorf("Expected 6 children, got %d", len(grid.children))
	}
}

func TestPanelSetTheme(t *testing.T) {
	panel := NewPanel(Size{Width: 100, Height: 100})
	theme := DefaultLight()

	// Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("SetTheme panicked: %v", r)
		}
	}()

	panel.SetTheme(theme)
}

func TestRowLayoutWithGap(t *testing.T) {
	row := NewRow()
	row.SetGap(10)

	btn1 := NewButton("1", Size{Width: 50, Height: 100})
	btn2 := NewButton("2", Size{Width: 50, Height: 100})

	row.Add(btn1)
	row.Add(btn2)

	if len(row.children) != 2 {
		t.Error("Row should have 2 children")
	}
}

func TestColumnLayoutWithGap(t *testing.T) {
	col := NewColumn()
	col.SetGap(10)

	btn1 := NewButton("1", Size{Width: 100, Height: 50})
	btn2 := NewButton("2", Size{Width: 100, Height: 50})

	col.Add(btn1)
	col.Add(btn2)

	if len(col.children) != 2 {
		t.Error("Column should have 2 children")
	}
}

func TestPanelLayoutWithPadding(t *testing.T) {
	panel := NewPanel(Size{Width: 100, Height: 100})
	panel.SetPadding(20)

	btn := NewButton("Test", Size{Width: 100, Height: 30})
	panel.Add(btn)

	if panel.padding != 20 {
		t.Error("Padding not set correctly")
	}
}

func TestGridZeroColumns(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("NewGrid(0) should panic")
		}
	}()

	NewGrid(0)
}

// TestPanelHandleEventBoundsCheck verifies that Panel only forwards pointer events
// to children whose bounds contain the event coordinates.
func TestPanelHandleEventBoundsCheck(t *testing.T) {
	panel := NewPanel(Size{Width: 100, Height: 100})
	panel.SetBounds(0, 0, 400, 400)

	// Create two buttons at different positions
	btn1 := NewButton("Button 1", Size{Width: 100, Height: 30})
	btn1.SetBounds(10, 10, 100, 30)
	btn1Clicked := false
	btn1.OnClick(func() {
		btn1Clicked = true
	})

	btn2 := NewButton("Button 2", Size{Width: 100, Height: 30})
	btn2.SetBounds(150, 10, 100, 30)
	btn2Clicked := false
	btn2.OnClick(func() {
		btn2Clicked = true
	})

	panel.Add(btn1)
	panel.Add(btn2)

	// Click on btn1's area (50, 20)
	pressEvt := NewPointerEvent(PointerButtonPress, 50, 20, PointerButtonLeft, 0, 0)
	panel.HandleEvent(pressEvt)

	releaseEvt := NewPointerEvent(PointerButtonRelease, 50, 20, PointerButtonLeft, 0, 0)
	panel.HandleEvent(releaseEvt)

	if !btn1Clicked {
		t.Error("Button 1 should be clicked when event is within its bounds")
	}
	if btn2Clicked {
		t.Error("Button 2 should NOT be clicked when event is outside its bounds")
	}

	// Reset
	btn1Clicked = false
	btn2Clicked = false

	// Click on btn2's area (180, 20)
	pressEvt2 := NewPointerEvent(PointerButtonPress, 180, 20, PointerButtonLeft, 0, 0)
	panel.HandleEvent(pressEvt2)

	releaseEvt2 := NewPointerEvent(PointerButtonRelease, 180, 20, PointerButtonLeft, 0, 0)
	panel.HandleEvent(releaseEvt2)

	if btn1Clicked {
		t.Error("Button 1 should NOT be clicked when event is outside its bounds")
	}
	if !btn2Clicked {
		t.Error("Button 2 should be clicked when event is within its bounds")
	}

	// Reset
	btn1Clicked = false
	btn2Clicked = false

	// Click outside both buttons (300, 300)
	pressEvt3 := NewPointerEvent(PointerButtonPress, 300, 300, PointerButtonLeft, 0, 0)
	panel.HandleEvent(pressEvt3)

	releaseEvt3 := NewPointerEvent(PointerButtonRelease, 300, 300, PointerButtonLeft, 0, 0)
	panel.HandleEvent(releaseEvt3)

	if btn1Clicked || btn2Clicked {
		t.Error("No button should be clicked when event is outside all child bounds")
	}
}

// TestPanelHandleEventTouchBoundsCheck verifies that Panel only forwards touch events
// to children whose bounds contain the event coordinates.
func TestPanelHandleEventTouchBoundsCheck(t *testing.T) {
	panel := NewPanel(Size{Width: 100, Height: 100})
	panel.SetBounds(0, 0, 400, 400)

	// Create a button
	btn := NewButton("Button", Size{Width: 100, Height: 30})
	btn.SetBounds(10, 10, 100, 30)

	eventReceived := false
	btn.OnEvent(func(evt Event) bool {
		if evt.Type() == EventTypeTouch {
			eventReceived = true
		}
		return true
	})

	panel.Add(btn)

	// Touch within button bounds
	touchEvt := NewTouchEvent(TouchDown, 1, 50, 20)
	panel.HandleEvent(touchEvt)

	if !eventReceived {
		t.Error("Button should receive touch event when touch is within its bounds")
	}

	// Reset
	eventReceived = false

	// Touch outside button bounds
	touchEvt2 := NewTouchEvent(TouchDown, 2, 200, 200)
	panel.HandleEvent(touchEvt2)

	if eventReceived {
		t.Error("Button should NOT receive touch event when touch is outside its bounds")
	}
}

// TestPanelHandleEventKeyBroadcast verifies that Panel broadcasts keyboard events
// to all children regardless of position (non-spatial events).
func TestPanelHandleEventKeyBroadcast(t *testing.T) {
	panel := NewPanel(Size{Width: 100, Height: 100})
	panel.SetBounds(0, 0, 400, 400)

	// Create a widget that consumes key events
	btn := NewButton("Button", Size{Width: 100, Height: 30})
	btn.SetBounds(10, 10, 100, 30)

	keyReceived := false
	btn.OnEvent(func(evt Event) bool {
		if evt.Type() == EventTypeKey {
			keyReceived = true
			return true
		}
		return false
	})

	panel.Add(btn)

	// Send key event (non-spatial)
	keyEvt := NewKeyEvent(KeyPress, KeyReturn, 0, '\n')
	panel.HandleEvent(keyEvt)

	if !keyReceived {
		t.Error("Button should receive keyboard event (broadcast to all children)")
	}
}

// TestResolveTreeFlowRow verifies that resolveTree correctly positions children
// in a horizontal flow layout.
func TestResolveTreeFlowRow(t *testing.T) {
	panel := NewPanel(Size{Width: 100, Height: 100})
	panel.SetFlowDirection(FlowRow)
	panel.SetPadding(0)
	panel.SetGap(0)

	// Add three children with 33.33% width each
	btn1 := NewButton("1", Size{Width: 33.33, Height: 100})
	btn2 := NewButton("2", Size{Width: 33.33, Height: 100})
	btn3 := NewButton("3", Size{Width: 33.33, Height: 100})

	panel.Add(btn1)
	panel.Add(btn2)
	panel.Add(btn3)

	// Resolve with 300x100 pixel bounds
	resolveTree(panel, 0, 0, 300, 100)

	// Check panel bounds
	x, y := panel.Position()
	w, h := panel.Bounds()
	if x != 0 || y != 0 || w != 300 || h != 100 {
		t.Errorf("Panel bounds incorrect: got (%d,%d,%d,%d), want (0,0,300,100)", x, y, w, h)
	}

	// Check first button position (should be at 0,0 with ~100px width)
	x1, y1 := btn1.Position()
	w1, h1 := btn1.Bounds()
	if x1 != 0 || y1 != 0 {
		t.Errorf("Button 1 position incorrect: got (%d,%d), want (0,0)", x1, y1)
	}
	if w1 < 99 || w1 > 100 {
		t.Errorf("Button 1 width incorrect: got %d, want ~100", w1)
	}
	if h1 != 100 {
		t.Errorf("Button 1 height incorrect: got %d, want 100", h1)
	}

	// Check second button position (should be at ~100,0)
	x2, y2 := btn2.Position()
	w2, h2 := btn2.Bounds()
	if x2 < 99 || x2 > 100 {
		t.Errorf("Button 2 x position incorrect: got %d, want ~100", x2)
	}
	if y2 != 0 {
		t.Errorf("Button 2 y position incorrect: got %d, want 0", y2)
	}
	if w2 < 99 || w2 > 100 {
		t.Errorf("Button 2 width incorrect: got %d, want ~100", w2)
	}
	if h2 != 100 {
		t.Errorf("Button 2 height incorrect: got %d, want 100", h2)
	}

	// Check third button position (should be at ~200,0)
	x3, y3 := btn3.Position()
	w3, h3 := btn3.Bounds()
	if x3 < 199 || x3 > 200 {
		t.Errorf("Button 3 x position incorrect: got %d, want ~200", x3)
	}
	if y3 != 0 {
		t.Errorf("Button 3 y position incorrect: got %d, want 0", y3)
	}
	if w3 < 99 || w3 > 100 {
		t.Errorf("Button 3 width incorrect: got %d, want ~100", w3)
	}
	if h3 != 100 {
		t.Errorf("Button 3 height incorrect: got %d, want 100", h3)
	}
}

// TestResolveTreeFlowColumn verifies that resolveTree correctly positions children
// in a vertical flow layout.
func TestResolveTreeFlowColumn(t *testing.T) {
	panel := NewPanel(Size{Width: 100, Height: 100})
	panel.SetFlowDirection(FlowColumn)
	panel.SetPadding(0)
	panel.SetGap(0)

	// Add three children with 33.33% height each
	btn1 := NewButton("1", Size{Width: 100, Height: 33.33})
	btn2 := NewButton("2", Size{Width: 100, Height: 33.33})
	btn3 := NewButton("3", Size{Width: 100, Height: 33.33})

	panel.Add(btn1)
	panel.Add(btn2)
	panel.Add(btn3)

	// Resolve with 200x300 pixel bounds
	resolveTree(panel, 0, 0, 200, 300)

	// Check panel bounds
	x, y := panel.Position()
	w, h := panel.Bounds()
	if x != 0 || y != 0 || w != 200 || h != 300 {
		t.Errorf("Panel bounds incorrect: got (%d,%d,%d,%d), want (0,0,200,300)", x, y, w, h)
	}

	// Check first button position (should be at 0,0 with ~100px height)
	x1, y1 := btn1.Position()
	w1, h1 := btn1.Bounds()
	if x1 != 0 || y1 != 0 {
		t.Errorf("Button 1 position incorrect: got (%d,%d), want (0,0)", x1, y1)
	}
	if w1 != 200 {
		t.Errorf("Button 1 width incorrect: got %d, want 200", w1)
	}
	if h1 < 99 || h1 > 100 {
		t.Errorf("Button 1 height incorrect: got %d, want ~100", h1)
	}

	// Check second button position (should be at 0,~100)
	x2, y2 := btn2.Position()
	w2, h2 := btn2.Bounds()
	if x2 != 0 {
		t.Errorf("Button 2 x position incorrect: got %d, want 0", x2)
	}
	if y2 < 99 || y2 > 100 {
		t.Errorf("Button 2 y position incorrect: got %d, want ~100", y2)
	}
	if w2 != 200 {
		t.Errorf("Button 2 width incorrect: got %d, want 200", w2)
	}
	if h2 < 99 || h2 > 100 {
		t.Errorf("Button 2 height incorrect: got %d, want ~100", h2)
	}

	// Check third button position (should be at 0,~200)
	x3, y3 := btn3.Position()
	w3, h3 := btn3.Bounds()
	if x3 != 0 {
		t.Errorf("Button 3 x position incorrect: got %d, want 0", x3)
	}
	if y3 < 199 || y3 > 200 {
		t.Errorf("Button 3 y position incorrect: got %d, want ~200", y3)
	}
	if w3 != 200 {
		t.Errorf("Button 3 width incorrect: got %d, want 200", w3)
	}
	if h3 < 99 || h3 > 100 {
		t.Errorf("Button 3 height incorrect: got %d, want ~100", h3)
	}
}

// TestResolveTreeWithPadding verifies that resolveTree correctly applies padding.
func TestResolveTreeWithPadding(t *testing.T) {
	panel := NewPanel(Size{Width: 100, Height: 100})
	panel.SetFlowDirection(FlowColumn)
	panel.SetPadding(20)
	panel.SetGap(0)

	// Add a single child with 100% size
	btn := NewButton("Button", Size{Width: 100, Height: 100})
	panel.Add(btn)

	// Resolve with 200x200 pixel bounds
	resolveTree(panel, 0, 0, 200, 200)

	// Panel should be 200x200 at (0,0)
	x, y := panel.Position()
	w, h := panel.Bounds()
	if x != 0 || y != 0 || w != 200 || h != 200 {
		t.Errorf("Panel bounds incorrect: got (%d,%d,%d,%d), want (0,0,200,200)", x, y, w, h)
	}

	// Button should be inset by padding (20px on all sides)
	// Content area: 160x160 at (20,20)
	btnX, btnY := btn.Position()
	btnW, btnH := btn.Bounds()
	if btnX != 20 || btnY != 20 {
		t.Errorf("Button position incorrect: got (%d,%d), want (20,20)", btnX, btnY)
	}
	if btnW != 160 || btnH != 160 {
		t.Errorf("Button size incorrect: got (%d,%d), want (160,160)", btnW, btnH)
	}
}

// TestResolveTreeWithGap verifies that resolveTree correctly applies gap between children.
func TestResolveTreeWithGap(t *testing.T) {
	panel := NewPanel(Size{Width: 100, Height: 100})
	panel.SetFlowDirection(FlowRow)
	panel.SetPadding(0)
	panel.SetGap(10)

	// Add two children with 50% width each
	btn1 := NewButton("1", Size{Width: 50, Height: 100})
	btn2 := NewButton("2", Size{Width: 50, Height: 100})

	panel.Add(btn1)
	panel.Add(btn2)

	// Resolve with 200x100 pixel bounds
	resolveTree(panel, 0, 0, 200, 100)

	// First button should be at (0,0) with 100px width
	x1, y1 := btn1.Position()
	w1, h1 := btn1.Bounds()
	if x1 != 0 || y1 != 0 {
		t.Errorf("Button 1 position incorrect: got (%d,%d), want (0,0)", x1, y1)
	}
	if w1 != 100 {
		t.Errorf("Button 1 width incorrect: got %d, want 100", w1)
	}
	if h1 != 100 {
		t.Errorf("Button 1 height incorrect: got %d, want 100", h1)
	}

	// Second button should be at (110,0) with 100px width (100 + 10 gap)
	x2, y2 := btn2.Position()
	w2, h2 := btn2.Bounds()
	if x2 != 110 || y2 != 0 {
		t.Errorf("Button 2 position incorrect: got (%d,%d), want (110,0)", x2, y2)
	}
	if w2 != 100 {
		t.Errorf("Button 2 width incorrect: got %d, want 100", w2)
	}
	if h2 != 100 {
		t.Errorf("Button 2 height incorrect: got %d, want 100", h2)
	}
}

// TestResolveTreeWithPaddingAndGap verifies that resolveTree correctly applies
// both padding and gap.
func TestResolveTreeWithPaddingAndGap(t *testing.T) {
	panel := NewPanel(Size{Width: 100, Height: 100})
	panel.SetFlowDirection(FlowColumn)
	panel.SetPadding(10)
	panel.SetGap(5)

	// Add two children with 50% height each
	btn1 := NewButton("1", Size{Width: 100, Height: 50})
	btn2 := NewButton("2", Size{Width: 100, Height: 50})

	panel.Add(btn1)
	panel.Add(btn2)

	// Resolve with 120x120 pixel bounds
	resolveTree(panel, 0, 0, 120, 120)

	// Content area after padding: 100x100 at (10,10)
	// First button: 50% of 100 = 50px height, at (10,10)
	x1, y1 := btn1.Position()
	w1, h1 := btn1.Bounds()
	if x1 != 10 || y1 != 10 {
		t.Errorf("Button 1 position incorrect: got (%d,%d), want (10,10)", x1, y1)
	}
	if w1 != 100 {
		t.Errorf("Button 1 width incorrect: got %d, want 100", w1)
	}
	if h1 != 50 {
		t.Errorf("Button 1 height incorrect: got %d, want 50", h1)
	}

	// Second button: at (10, 10+50+5) = (10,65)
	x2, y2 := btn2.Position()
	w2, h2 := btn2.Bounds()
	if x2 != 10 || y2 != 65 {
		t.Errorf("Button 2 position incorrect: got (%d,%d), want (10,65)", x2, y2)
	}
	if w2 != 100 {
		t.Errorf("Button 2 width incorrect: got %d, want 100", w2)
	}
	if h2 != 50 {
		t.Errorf("Button 2 height incorrect: got %d, want 50", h2)
	}
}

// TestResolveTreeAlignCenter verifies that resolveTree correctly centers children
// on the cross axis.
func TestResolveTreeAlignCenter(t *testing.T) {
	panel := NewPanel(Size{Width: 100, Height: 100})
	panel.SetFlowDirection(FlowRow)
	panel.SetAlign(AlignCenter)
	panel.SetPadding(0)
	panel.SetGap(0)

	// Add a child with 50% width and 50% height
	btn := NewButton("Button", Size{Width: 50, Height: 50})
	panel.Add(btn)

	// Resolve with 200x200 pixel bounds
	resolveTree(panel, 0, 0, 200, 200)

	// For FlowRow, cross axis is vertical, so button should be vertically centered
	// Button height: 50% of 200 = 100px
	// Center position: (200 - 100) / 2 = 50px from top
	x, y := btn.Position()
	w, h := btn.Bounds()
	if x != 0 {
		t.Errorf("Button x incorrect: got %d, want 0", x)
	}
	if y != 50 {
		t.Errorf("Button y incorrect: got %d, want 50 (centered)", y)
	}
	if w != 100 {
		t.Errorf("Button width incorrect: got %d, want 100", w)
	}
	if h != 100 {
		t.Errorf("Button height incorrect: got %d, want 100", h)
	}
}

// TestResolveTreeAlignEnd verifies that resolveTree correctly aligns children
// to the end of the cross axis.
func TestResolveTreeAlignEnd(t *testing.T) {
	panel := NewPanel(Size{Width: 100, Height: 100})
	panel.SetFlowDirection(FlowColumn)
	panel.SetAlign(AlignEnd)
	panel.SetPadding(0)
	panel.SetGap(0)

	// Add a child with 50% width and 50% height
	btn := NewButton("Button", Size{Width: 50, Height: 50})
	panel.Add(btn)

	// Resolve with 200x200 pixel bounds
	resolveTree(panel, 0, 0, 200, 200)

	// For FlowColumn, cross axis is horizontal, so button should be right-aligned
	// Button width: 50% of 200 = 100px
	// End position: 200 - 100 = 100px from left
	x, y := btn.Position()
	w, h := btn.Bounds()
	if x != 100 {
		t.Errorf("Button x incorrect: got %d, want 100 (end-aligned)", x)
	}
	if y != 0 {
		t.Errorf("Button y incorrect: got %d, want 0", y)
	}
	if w != 100 {
		t.Errorf("Button width incorrect: got %d, want 100", w)
	}
	if h != 100 {
		t.Errorf("Button height incorrect: got %d, want 100", h)
	}
}

// TestResolveTreeAlignStretch verifies that resolveTree correctly stretches children
// to fill the cross axis.
func TestResolveTreeAlignStretch(t *testing.T) {
	panel := NewPanel(Size{Width: 100, Height: 100})
	panel.SetFlowDirection(FlowRow)
	panel.SetAlign(AlignStretch)
	panel.SetPadding(0)
	panel.SetGap(0)

	// Add a child with 50% width and 50% height
	btn := NewButton("Button", Size{Width: 50, Height: 50})
	panel.Add(btn)

	// Resolve with 200x200 pixel bounds
	resolveTree(panel, 0, 0, 200, 200)

	// For FlowRow with AlignStretch, button height should be stretched to fill
	// the full content height (200px), ignoring the 50% size hint
	x, y := btn.Position()
	w, h := btn.Bounds()
	if x != 0 {
		t.Errorf("Button x incorrect: got %d, want 0", x)
	}
	if y != 0 {
		t.Errorf("Button y incorrect: got %d, want 0", y)
	}
	if w != 100 {
		t.Errorf("Button width incorrect: got %d, want 100", w)
	}
	if h != 200 {
		t.Errorf("Button height incorrect: got %d, want 200 (stretched)", h)
	}
}

// TestResolveTreeNestedPanels verifies that resolveTree correctly handles
// nested panel layouts.
func TestResolveTreeNestedPanels(t *testing.T) {
	// Create outer panel (vertical layout)
	outer := NewPanel(Size{Width: 100, Height: 100})
	outer.SetFlowDirection(FlowColumn)
	outer.SetPadding(10)
	outer.SetGap(0)

	// Create inner panel (horizontal layout) with 100% width, 50% height
	inner := NewPanel(Size{Width: 100, Height: 50})
	inner.SetFlowDirection(FlowRow)
	inner.SetPadding(5)
	inner.SetGap(0)

	// Add buttons to inner panel
	btn1 := NewButton("1", Size{Width: 50, Height: 100})
	btn2 := NewButton("2", Size{Width: 50, Height: 100})
	inner.Add(btn1)
	inner.Add(btn2)

	// Add inner panel to outer panel
	outer.Add(inner)

	// Resolve with 200x200 pixel bounds
	resolveTree(outer, 0, 0, 200, 200)

	// Outer panel: 200x200 at (0,0)
	outerX, outerY := outer.Position()
	outerW, outerH := outer.Bounds()
	if outerX != 0 || outerY != 0 || outerW != 200 || outerH != 200 {
		t.Errorf("Outer panel bounds incorrect: got (%d,%d,%d,%d), want (0,0,200,200)",
			outerX, outerY, outerW, outerH)
	}

	// Inner panel: content area of outer is 180x180 at (10,10)
	// Inner takes 50% height = 90px, positioned at (10,10)
	innerX, innerY := inner.Position()
	innerW, innerH := inner.Bounds()
	if innerX != 10 || innerY != 10 {
		t.Errorf("Inner panel position incorrect: got (%d,%d), want (10,10)", innerX, innerY)
	}
	if innerW != 180 || innerH != 90 {
		t.Errorf("Inner panel size incorrect: got (%d,%d), want (180,90)", innerW, innerH)
	}

	// Buttons: inner content area is 170x80 at (15,15) [inner pos + 5px padding]
	// Each button gets 50% width = 85px
	btn1X, btn1Y := btn1.Position()
	btn1W, btn1H := btn1.Bounds()
	if btn1X != 15 || btn1Y != 15 {
		t.Errorf("Button 1 position incorrect: got (%d,%d), want (15,15)", btn1X, btn1Y)
	}
	if btn1W != 85 || btn1H != 80 {
		t.Errorf("Button 1 size incorrect: got (%d,%d), want (85,80)", btn1W, btn1H)
	}

	btn2X, btn2Y := btn2.Position()
	btn2W, btn2H := btn2.Bounds()
	if btn2X != 100 || btn2Y != 15 {
		t.Errorf("Button 2 position incorrect: got (%d,%d), want (100,15)", btn2X, btn2Y)
	}
	if btn2W != 85 || btn2H != 80 {
		t.Errorf("Button 2 size incorrect: got (%d,%d), want (85,80)", btn2W, btn2H)
	}
}

// TestResolveTreeGrid verifies that Grid layout works correctly.
// Note: Grid currently inherits Panel's linear flow, not true grid layout.
func TestResolveTreeGrid(t *testing.T) {
	grid := NewGrid(2)
	grid.SetPadding(0)
	grid.SetGap(0)

	// Add 4 buttons (2x2 grid)
	for i := 0; i < 4; i++ {
		btn := NewButton("Btn", Size{Width: 100, Height: 100})
		grid.Add(btn)
	}

	// Resolve with 200x200 pixel bounds
	resolveTree(grid, 0, 0, 200, 200)

	// Grid should be positioned at (0,0) with 200x200 size
	x, y := grid.Position()
	w, h := grid.Bounds()
	if x != 0 || y != 0 || w != 200 || h != 200 {
		t.Errorf("Grid bounds incorrect: got (%d,%d,%d,%d), want (0,0,200,200)", x, y, w, h)
	}

	// Note: Grid currently uses Panel's FlowColumn layout (default)
	// so children stack vertically. A proper grid implementation would
	// need custom resolveChildren logic.
}
