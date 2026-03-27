//go:build windows || darwin || android || ios

package wayne_test

import (
	"fmt"

	"github.com/opd-ai/wayne"
)

// Example demonstrates creating a simple application with a button.
func Example() {
	app := wayne.NewApp()
	panel := wayne.NewPanel(wayne.Size{Width: 100, Height: 100})
	btn := wayne.NewButton("Click Me", wayne.Size{Width: 50, Height: 10})
	btn.OnClick(func() {
		fmt.Println("Button clicked!")
	})
	panel.Add(btn)
	app.SetRoot(panel)
	// In a real application: app.Run(wayne.WindowConfig{Title: "Demo"})
	fmt.Println("App created with button")
	// Output: App created with button
}

// ExampleNewButton demonstrates creating and configuring a Button widget.
func ExampleNewButton() {
	btn := wayne.NewButton("Submit", wayne.Size{Width: 30, Height: 8})
	btn.OnClick(func() {
		fmt.Println("Submitted!")
	})
	btn.SetEnabled(true)
	fmt.Printf("Button label can be set: %T\n", btn)
	// Output: Button label can be set: *wayne.Button
}

// ExampleNewLabel demonstrates creating a text label widget.
func ExampleNewLabel() {
	label := wayne.NewLabel("Hello, World!", wayne.Size{Width: 50, Height: 5})
	label.SetText("Updated text")
	fmt.Printf("Label created: %T\n", label)
	// Output: Label created: *wayne.Label
}

// ExampleNewTextInput demonstrates creating a text input field.
func ExampleNewTextInput() {
	input := wayne.NewTextInput("Enter name...", wayne.Size{Width: 60, Height: 8})
	input.SetText("Default value")
	text := input.Text()
	fmt.Printf("Input text: %s\n", text)
	// Output: Input text: Default value
}

// ExampleNewPanel demonstrates creating a container panel.
func ExampleNewPanel() {
	panel := wayne.NewPanel(wayne.Size{Width: 100, Height: 100})
	panel.Add(wayne.NewLabel("Item 1", wayne.Size{Width: 50, Height: 10}))
	panel.Add(wayne.NewLabel("Item 2", wayne.Size{Width: 50, Height: 10}))
	fmt.Printf("Panel children: %d\n", len(panel.Children()))
	// Output: Panel children: 2
}

// ExampleNewRow demonstrates creating a horizontal layout container.
func ExampleNewRow() {
	row := wayne.NewRow()
	row.SetGap(5)
	row.SetAlign(wayne.AlignCenter)
	row.Add(wayne.NewButton("A", wayne.Size{Width: 30, Height: 10}))
	row.Add(wayne.NewButton("B", wayne.Size{Width: 30, Height: 10}))
	fmt.Printf("Row children: %d\n", len(row.Children()))
	// Output: Row children: 2
}

// ExampleNewColumn demonstrates creating a vertical layout container.
func ExampleNewColumn() {
	col := wayne.NewColumn()
	col.SetGap(10)
	col.Add(wayne.NewLabel("Top", wayne.Size{Width: 100, Height: 10}))
	col.Add(wayne.NewLabel("Bottom", wayne.Size{Width: 100, Height: 10}))
	fmt.Printf("Column children: %d\n", len(col.Children()))
	// Output: Column children: 2
}

// ExampleNewGrid demonstrates creating a grid layout container.
func ExampleNewGrid() {
	grid := wayne.NewGrid(3) // 3 columns
	for i := 0; i < 6; i++ {
		grid.Add(wayne.NewLabel(fmt.Sprintf("Cell %d", i), wayne.Size{Width: 30, Height: 10}))
	}
	fmt.Printf("Grid children: %d\n", len(grid.Children()))
	// Output: Grid children: 6
}

// ExampleNewStack demonstrates creating a stacked (overlay) layout.
func ExampleNewStack() {
	stack := wayne.NewStack()
	stack.Add(wayne.NewLabel("Background", wayne.Size{Width: 100, Height: 100}))
	stack.Add(wayne.NewButton("Overlay", wayne.Size{Width: 50, Height: 10}))
	fmt.Printf("Stack children: %d\n", len(stack.Children()))
	// Output: Stack children: 2
}

// ExampleNewScrollView demonstrates creating a scrollable container.
func ExampleNewScrollView() {
	scroll := wayne.NewScrollView(wayne.Size{Width: 100, Height: 50})
	content := wayne.NewColumn()
	for i := 0; i < 20; i++ {
		content.Add(wayne.NewLabel(fmt.Sprintf("Item %d", i), wayne.Size{Width: 100, Height: 5}))
	}
	scroll.Add(content) // Add content to scroll view
	fmt.Printf("ScrollView created: %T\n", scroll)
	// Output: ScrollView created: *wayne.ScrollView
}

// ExampleNewSpacer demonstrates creating a flexible spacer widget.
func ExampleNewSpacer() {
	row := wayne.NewRow()
	row.Add(wayne.NewLabel("Left", wayne.Size{Width: 20, Height: 10}))
	row.Add(wayne.NewSpacer(wayne.Size{Width: 60, Height: 10})) // Pushes content apart
	row.Add(wayne.NewLabel("Right", wayne.Size{Width: 20, Height: 10}))
	fmt.Printf("Row with spacer: %d children\n", len(row.Children()))
	// Output: Row with spacer: 3 children
}

// ExampleNewImageWidget demonstrates creating an image display widget.
func ExampleNewImageWidget() {
	img := wayne.NewImageWidget(wayne.Size{Width: 50, Height: 50})
	// Load image: img.SetImage(loadedImage)
	fmt.Printf("ImageWidget created: %T\n", img)
	// Output: ImageWidget created: *wayne.ImageWidget
}

// ExampleNewApp demonstrates creating a basic application.
func ExampleNewApp() {
	app := wayne.NewApp()
	app.SetTheme(wayne.DefaultDark())
	fmt.Printf("App created: %T\n", app)
	// Output: App created: *wayne.App
}

// ExampleNewAppWithConfig demonstrates creating an app with custom configuration.
func ExampleNewAppWithConfig() {
	cfg := wayne.AppConfig{
		Width:  800,
		Height: 600,
	}
	app := wayne.NewAppWithConfig(cfg)
	fmt.Printf("App with config: %T\n", app)
	// Output: App with config: *wayne.App
}

// Example_themes demonstrates using the built-in themes.
func Example_themes() {
	dark := wayne.DefaultDark()
	light := wayne.DefaultLight()
	contrast := wayne.HighContrast()
	fmt.Printf("Themes: dark=%T, light=%T, contrast=%T\n", dark, light, contrast)
	// Output: Themes: dark=wayne.Theme, light=wayne.Theme, contrast=wayne.Theme
}

// ExampleSize demonstrates the percentage-based size specification.
func ExampleSize() {
	// Sizes are specified as percentages of parent container
	full := wayne.Size{Width: 100, Height: 100}   // Full width and height
	half := wayne.Size{Width: 50, Height: 50}     // Half width and height
	narrow := wayne.Size{Width: 20, Height: 100}  // 20% width, full height
	fmt.Printf("Sizes: full=%v, half=%v, narrow=%v\n", full, half, narrow)
	// Output: Sizes: full={100 100}, half={50 50}, narrow={20 100}
}
