//go:build windows || darwin || android || ios

// Example: scrollview - A scrollable list demonstration.
//
// This example demonstrates:
// - Creating a scrollable container with many items
// - Using ScrollView to handle content larger than the viewport
// - Dynamically adding items to a container
// - Scroll offset handling
//
// Build and run on Windows:
//
//	cd examples/scrollview
//	go build -o scrollview.exe .
//	./scrollview.exe
//
// Build and run on macOS:
//
//	cd examples/scrollview
//	go build -o scrollview .
//	./scrollview
package main

import (
	"fmt"
	"log"

	"github.com/opd-ai/wayne"
)

func main() {
	app := wayne.NewApp()

	window, err := app.NewWindow(wayne.WindowConfig{
		Title:  "Wayne ScrollView Example",
		Width:  400,
		Height: 500,
	})
	if err != nil {
		log.Fatalf("Failed to create window: %v", err)
	}

	// Root panel.
	root := wayne.NewPanel(wayne.Size{Width: 100, Height: 100})
	root.SetFlowDirection(wayne.FlowColumn)
	root.SetPadding(10)
	root.SetGap(5)

	// Title.
	title := wayne.NewLabel("Scrollable Item List", wayne.Size{Width: 100, Height: 8})
	root.Add(title)

	// Scroll position indicator.
	scrollInfo := wayne.NewLabel("Scroll: 0px", wayne.Size{Width: 100, Height: 5})
	root.Add(scrollInfo)

	// ScrollView takes up most of the window.
	scrollView := wayne.NewScrollView(wayne.Size{Width: 100, Height: 75})
	scrollView.OnScroll(func(offset int) {
		scrollInfo.SetText(fmt.Sprintf("Scroll: %dpx", offset))
	})

	// Create a container panel with many items.
	// Each item is 5% of parent height, but parent is virtual (larger than view).
	itemContainer := wayne.NewPanel(wayne.Size{Width: 100, Height: 500}) // 500% = 5x viewport
	itemContainer.SetFlowDirection(wayne.FlowColumn)
	itemContainer.SetGap(2)

	// Add 30 items to the scrollable list.
	for i := 1; i <= 30; i++ {
		itemRow := wayne.NewRow()

		// Item label.
		label := wayne.NewLabel(fmt.Sprintf("Item #%d", i), wayne.Size{Width: 70, Height: 100})
		itemRow.Add(label)

		// Action button for each item.
		btn := wayne.NewButton("Select", wayne.Size{Width: 30, Height: 100})
		itemNum := i // Capture for closure.
		btn.OnClick(func() {
			title.SetText(fmt.Sprintf("Selected: Item #%d", itemNum))
		})
		itemRow.Add(btn)

		itemContainer.Add(itemRow)
	}

	scrollView.Add(itemContainer)
	root.Add(scrollView)

	// Button row at bottom.
	buttonRow := wayne.NewRow()

	topBtn := wayne.NewButton("Scroll to Top", wayne.Size{Width: 50, Height: 100})
	topBtn.OnClick(func() {
		scrollView.SetScrollOffset(0)
		scrollInfo.SetText("Scroll: 0px")
	})
	buttonRow.Add(topBtn)

	bottomBtn := wayne.NewButton("Scroll Down", wayne.Size{Width: 50, Height: 100})
	bottomBtn.OnClick(func() {
		// Scroll down by 100 pixels.
		current := scrollView.ScrollOffset()
		scrollView.SetScrollOffset(current + 100)
	})
	buttonRow.Add(bottomBtn)

	root.Add(buttonRow)

	window.SetRoot(root)

	if err := app.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
