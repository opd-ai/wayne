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

// buildScrollViewPanel creates the main panel with scrollable content.
func buildScrollViewPanel() *wayne.Panel {
	root := wayne.NewPanel(wayne.Size{Width: 100, Height: 100})
	root.SetFlowDirection(wayne.FlowColumn)
	root.SetPadding(10)
	root.SetGap(5)

	title := wayne.NewLabel("Scrollable Item List", wayne.Size{Width: 100, Height: 8})
	root.Add(title)

	scrollInfo := wayne.NewLabel("Scroll: 0px", wayne.Size{Width: 100, Height: 5})
	root.Add(scrollInfo)

	scrollView := buildScrollView(title, scrollInfo)
	root.Add(scrollView)

	buttonRow := buildScrollButtons(scrollView, scrollInfo)
	root.Add(buttonRow)

	return root
}

// buildScrollView creates the scrollable list with items.
func buildScrollView(title, scrollInfo *wayne.Label) *wayne.ScrollView {
	scrollView := wayne.NewScrollView(wayne.Size{Width: 100, Height: 75})
	scrollView.OnScroll(func(offset int) {
		scrollInfo.SetText(fmt.Sprintf("Scroll: %dpx", offset))
	})

	itemContainer := buildItemContainer(title)
	scrollView.Add(itemContainer)
	return scrollView
}

// buildItemContainer creates the panel containing all scrollable items.
func buildItemContainer(title *wayne.Label) *wayne.Panel {
	container := wayne.NewPanel(wayne.Size{Width: 100, Height: 500})
	container.SetFlowDirection(wayne.FlowColumn)
	container.SetGap(2)

	for i := 1; i <= 30; i++ {
		container.Add(buildItemRow(i, title))
	}
	return container
}

// buildItemRow creates a single item row with label and button.
func buildItemRow(itemNum int, title *wayne.Label) *wayne.Row {
	itemRow := wayne.NewRow()
	itemRow.Add(wayne.NewLabel(fmt.Sprintf("Item #%d", itemNum), wayne.Size{Width: 70, Height: 100}))

	btn := wayne.NewButton("Select", wayne.Size{Width: 30, Height: 100})
	btn.OnClick(func() {
		title.SetText(fmt.Sprintf("Selected: Item #%d", itemNum))
	})
	itemRow.Add(btn)
	return itemRow
}

// buildScrollButtons creates the scroll control buttons.
func buildScrollButtons(scrollView *wayne.ScrollView, scrollInfo *wayne.Label) *wayne.Row {
	buttonRow := wayne.NewRow()

	topBtn := wayne.NewButton("Scroll to Top", wayne.Size{Width: 50, Height: 100})
	topBtn.OnClick(func() {
		scrollView.SetScrollOffset(0)
		scrollInfo.SetText("Scroll: 0px")
	})
	buttonRow.Add(topBtn)

	bottomBtn := wayne.NewButton("Scroll Down", wayne.Size{Width: 50, Height: 100})
	bottomBtn.OnClick(func() {
		scrollView.SetScrollOffset(scrollView.ScrollOffset() + 100)
	})
	buttonRow.Add(bottomBtn)

	return buttonRow
}

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

	window.SetRoot(buildScrollViewPanel())

	if err := app.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
