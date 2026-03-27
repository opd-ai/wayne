//go:build windows || darwin || android || ios

// Example: hello - Minimal wayne application with a button and label.
//
// This example demonstrates the most basic wayne application:
// - Creating an App with default configuration
// - Opening a Window
// - Building a simple widget tree with Panel, Label, and Button
// - Handling button clicks
// - Running the application event loop
//
// Build and run on Windows:
//
//	cd examples/hello
//	go build -o hello.exe .
//	./hello.exe
//
// Build and run on macOS:
//
//	cd examples/hello
//	go build -o hello .
//	./hello
package main

import (
	"fmt"
	"log"

	"github.com/opd-ai/wayne"
)

func main() {
	// Create the application with default configuration.
	app := wayne.NewApp()

	// Create the main window.
	window, err := app.NewWindow(wayne.WindowConfig{
		Title:  "Hello Wayne",
		Width:  400,
		Height: 300,
	})
	if err != nil {
		log.Fatalf("Failed to create window: %v", err)
	}

	// Build the widget tree.
	// Root panel takes 100% of the window area.
	root := wayne.NewPanel(wayne.Size{Width: 100, Height: 100})
	root.SetFlowDirection(wayne.FlowColumn)
	root.SetPadding(20)
	root.SetGap(10)

	// Add a greeting label (70% width, 20% height of parent).
	label := wayne.NewLabel("Hello, Wayne!", wayne.Size{Width: 70, Height: 20})
	root.Add(label)

	// Add a spacer to push button down.
	root.Add(wayne.NewSpacer(wayne.Size{Width: 100, Height: 40}))

	// Add a button (50% width, 20% height of parent).
	clickCount := 0
	button := wayne.NewButton("Click Me!", wayne.Size{Width: 50, Height: 20})
	button.OnClick(func() {
		clickCount++
		label.SetText(fmt.Sprintf("Clicked %d times!", clickCount))
	})
	root.Add(button)

	// Set the widget tree as the window's root.
	window.SetRoot(root)

	// Run the application (blocks until window is closed).
	if err := app.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
