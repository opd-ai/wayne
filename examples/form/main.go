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

	// Root panel with column layout.
	root := wayne.NewPanel(wayne.Size{Width: 100, Height: 100})
	root.SetFlowDirection(wayne.FlowColumn)
	root.SetPadding(15)
	root.SetGap(10)

	// Title label.
	title := wayne.NewLabel("Registration Form", wayne.Size{Width: 100, Height: 10})
	root.Add(title)

	// Form fields using helper function.
	nameRow, nameInput := createInputRow("Name:", "Enter your name")
	root.Add(nameRow)

	emailRow, emailInput := createInputRow("Email:", "Enter your email")
	root.Add(emailRow)

	msgRow, msgInput := createInputRow("Message:", "Enter a message")
	root.Add(msgRow)

	// Spacer.
	root.Add(wayne.NewSpacer(wayne.Size{Width: 100, Height: 10}))

	// Status label for validation feedback.
	statusLabel := wayne.NewLabel("", wayne.Size{Width: 100, Height: 10})
	root.Add(statusLabel)

	// Button row.
	buttonRow := wayne.NewRow()

	submitBtn := wayne.NewButton("Submit", wayne.Size{Width: 30, Height: 100})
	submitBtn.OnClick(func() {
		name := nameInput.Text()
		email := emailInput.Text()
		msg := msgInput.Text()

		// Simple validation.
		if name == "" {
			statusLabel.SetText("Error: Name is required")
			return
		}
		if email == "" {
			statusLabel.SetText("Error: Email is required")
			return
		}

		// Success.
		statusLabel.SetText(fmt.Sprintf("Submitted: %s <%s>", name, email))
		fmt.Printf("Form submitted:\n  Name: %s\n  Email: %s\n  Message: %s\n", name, email, msg)
	})
	buttonRow.Add(submitBtn)

	buttonRow.Add(wayne.NewSpacer(wayne.Size{Width: 40, Height: 100}))

	clearBtn := wayne.NewButton("Clear", wayne.Size{Width: 30, Height: 100})
	clearBtn.OnClick(func() {
		nameInput.SetText("")
		emailInput.SetText("")
		msgInput.SetText("")
		statusLabel.SetText("Form cleared")
	})
	buttonRow.Add(clearBtn)

	root.Add(buttonRow)

	window.SetRoot(root)

	if err := app.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
