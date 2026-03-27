//go:build windows || darwin || android || ios

// Example: form - A form application with text inputs, labels, and validation.
//
// This example demonstrates:
// - Creating a form layout with labels and text inputs
// - Using Row and Panel containers for layout
// - Reading text input values
// - Simple form validation and feedback
// - Button actions with state
//
// Build and run on Windows:
//
//	cd examples/form
//	go build -o form.exe .
//	./form.exe
//
// Build and run on macOS:
//
//	cd examples/form
//	go build -o form .
//	./form
package main

import (
	"fmt"
	"log"

	"github.com/opd-ai/wayne"
)

// createInputRow creates a row with a label and text input for form fields.
func createInputRow(label, placeholder string) (*wayne.Row, *wayne.TextInput) {
	row := wayne.NewRow()
	row.Add(wayne.NewLabel(label, wayne.Size{Width: 25, Height: 100}))
	input := wayne.NewTextInput(placeholder, wayne.Size{Width: 75, Height: 100})
	row.Add(input)
	return row, input
}

// formInputs holds references to all form input fields.
type formInputs struct {
	name  *wayne.TextInput
	email *wayne.TextInput
	msg   *wayne.TextInput
}

// buildFormPanel creates the main form panel with all widgets.
func buildFormPanel() (*wayne.Panel, *formInputs, *wayne.Label) {
	root := wayne.NewPanel(wayne.Size{Width: 100, Height: 100})
	root.SetFlowDirection(wayne.FlowColumn)
	root.SetPadding(15)
	root.SetGap(10)

	root.Add(wayne.NewLabel("Registration Form", wayne.Size{Width: 100, Height: 10}))

	inputs := &formInputs{}
	var nameRow, emailRow, msgRow *wayne.Row
	nameRow, inputs.name = createInputRow("Name:", "Enter your name")
	emailRow, inputs.email = createInputRow("Email:", "Enter your email")
	msgRow, inputs.msg = createInputRow("Message:", "Enter a message")
	root.Add(nameRow)
	root.Add(emailRow)
	root.Add(msgRow)

	root.Add(wayne.NewSpacer(wayne.Size{Width: 100, Height: 10}))

	statusLabel := wayne.NewLabel("", wayne.Size{Width: 100, Height: 10})
	root.Add(statusLabel)

	buttonRow := buildButtonRow(inputs, statusLabel)
	root.Add(buttonRow)

	return root, inputs, statusLabel
}

// buildButtonRow creates the submit and clear buttons.
func buildButtonRow(inputs *formInputs, statusLabel *wayne.Label) *wayne.Row {
	buttonRow := wayne.NewRow()

	submitBtn := wayne.NewButton("Submit", wayne.Size{Width: 30, Height: 100})
	submitBtn.OnClick(func() { handleSubmit(inputs, statusLabel) })
	buttonRow.Add(submitBtn)

	buttonRow.Add(wayne.NewSpacer(wayne.Size{Width: 40, Height: 100}))

	clearBtn := wayne.NewButton("Clear", wayne.Size{Width: 30, Height: 100})
	clearBtn.OnClick(func() { handleClear(inputs, statusLabel) })
	buttonRow.Add(clearBtn)

	return buttonRow
}

// handleSubmit validates and submits the form.
func handleSubmit(inputs *formInputs, statusLabel *wayne.Label) {
	name := inputs.name.Text()
	email := inputs.email.Text()
	msg := inputs.msg.Text()

	if name == "" {
		statusLabel.SetText("Error: Name is required")
		return
	}
	if email == "" {
		statusLabel.SetText("Error: Email is required")
		return
	}

	statusLabel.SetText(fmt.Sprintf("Submitted: %s <%s>", name, email))
	fmt.Printf("Form submitted:\n  Name: %s\n  Email: %s\n  Message: %s\n", name, email, msg)
}

// handleClear resets all form fields.
func handleClear(inputs *formInputs, statusLabel *wayne.Label) {
	inputs.name.SetText("")
	inputs.email.SetText("")
	inputs.msg.SetText("")
	statusLabel.SetText("Form cleared")
}

func main() {
	app := wayne.NewApp()

	window, err := app.NewWindow(wayne.WindowConfig{
		Title:  "Wayne Form Example",
		Width:  500,
		Height: 400,
	})
	if err != nil {
		log.Fatalf("Failed to create window: %v", err)
	}

	root, _, _ := buildFormPanel()
	window.SetRoot(root)

	if err := app.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
