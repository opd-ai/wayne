//go:build windows || darwin || android || ios

package wayne

import (
	"testing"
)

// BenchmarkResolveTree_SmallTree benchmarks resolveTree with a small widget tree.
func BenchmarkResolveTree_SmallTree(b *testing.B) {
	// Create a simple tree: panel with 3 buttons
	panel := NewPanel(Size{Width: 100, Height: 100})
	panel.Add(NewButton("Button 1", Size{Width: 30, Height: 10}))
	panel.Add(NewButton("Button 2", Size{Width: 30, Height: 10}))
	panel.Add(NewButton("Button 3", Size{Width: 30, Height: 10}))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resolveTree(panel, 0, 0, 800, 600)
	}
}

// BenchmarkResolveTree_MediumTree benchmarks resolveTree with a medium widget tree.
func BenchmarkResolveTree_MediumTree(b *testing.B) {
	// Create a tree with nested panels: 3 panels, each with 5 widgets
	root := NewPanel(Size{Width: 100, Height: 100})
	for i := 0; i < 3; i++ {
		panel := NewPanel(Size{Width: 100, Height: 30})
		for j := 0; j < 5; j++ {
			panel.Add(NewButton("Btn", Size{Width: 18, Height: 100}))
		}
		root.Add(panel)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resolveTree(root, 0, 0, 800, 600)
	}
}

// BenchmarkResolveTree_LargeTree benchmarks resolveTree with a large widget tree.
func BenchmarkResolveTree_LargeTree(b *testing.B) {
	// Create a deep tree: 5 levels of nested panels with widgets
	root := NewPanel(Size{Width: 100, Height: 100})
	current := root
	for level := 0; level < 5; level++ {
		panel := NewPanel(Size{Width: 100, Height: 80})
		for j := 0; j < 4; j++ {
			panel.Add(NewLabel("Label", Size{Width: 20, Height: 10}))
		}
		current.Add(panel)
		current = panel
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resolveTree(root, 0, 0, 1920, 1080)
	}
}

// BenchmarkResolveTree_GridLayout benchmarks Grid layout specifically.
func BenchmarkResolveTree_GridLayout(b *testing.B) {
	// Create a 4-column grid with 16 items
	grid := NewGrid(4)
	for i := 0; i < 16; i++ {
		grid.Add(NewButton("Cell", Size{Width: 100, Height: 100}))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resolveTree(grid, 0, 0, 800, 600)
	}
}

// BenchmarkResolveTree_StackLayout benchmarks Stack layout specifically.
func BenchmarkResolveTree_StackLayout(b *testing.B) {
	// Create a stack with 5 overlapping panels
	stack := NewStack()
	for i := 0; i < 5; i++ {
		stack.Add(NewPanel(Size{Width: 80, Height: 80}))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resolveTree(stack, 0, 0, 800, 600)
	}
}

// BenchmarkHitTest_ShallowTree benchmarks hitTest with a flat widget structure.
func BenchmarkHitTest_ShallowTree(b *testing.B) {
	panel := NewPanel(Size{Width: 100, Height: 100})
	for i := 0; i < 10; i++ {
		panel.Add(NewButton("Btn", Size{Width: 10, Height: 100}))
	}
	// Resolve layout first
	resolveTree(panel, 0, 0, 800, 600)

	dispatcher := NewEventDispatcher()
	dispatcher.SetWidgetRoot(panel)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Hit test at various positions
		dispatcher.hitTest(panel, 400, 300)
	}
}

// BenchmarkHitTest_DeepTree benchmarks hitTest with nested widgets.
func BenchmarkHitTest_DeepTree(b *testing.B) {
	root := NewPanel(Size{Width: 100, Height: 100})
	current := root
	for level := 0; level < 5; level++ {
		panel := NewPanel(Size{Width: 90, Height: 90})
		panel.Add(NewButton("Nested", Size{Width: 50, Height: 20}))
		current.Add(panel)
		current = panel
	}
	// Resolve layout first
	resolveTree(root, 0, 0, 800, 600)

	dispatcher := NewEventDispatcher()
	dispatcher.SetWidgetRoot(root)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Hit test at center where nested widgets overlap
		dispatcher.hitTest(root, 100, 100)
	}
}

// BenchmarkHitTest_Miss benchmarks hitTest when no widget is hit.
func BenchmarkHitTest_Miss(b *testing.B) {
	panel := NewPanel(Size{Width: 50, Height: 50}) // Only covers top-left quadrant
	panel.Add(NewButton("Btn", Size{Width: 100, Height: 100}))
	// Resolve layout first
	resolveTree(panel, 0, 0, 800, 600)

	dispatcher := NewEventDispatcher()
	dispatcher.SetWidgetRoot(panel)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Hit test outside the panel
		dispatcher.hitTest(panel, 700, 500)
	}
}
