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
