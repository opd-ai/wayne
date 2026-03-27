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

	// Name field row.
	nameRow := wayne.NewRow()
	nameRow.Add(wayne.NewLabel("Name:", wayne.Size{Width: 25, Height: 100}))
	nameInput := wayne.NewTextInput("Enter your name", wayne.Size{Width: 75, Height: 100})
	nameRow.Add(nameInput)
	root.Add(nameRow)

	// Email field row.
	emailRow := wayne.NewRow()
	emailRow.Add(wayne.NewLabel("Email:", wayne.Size{Width: 25, Height: 100}))
	emailInput := wayne.NewTextInput("Enter your email", wayne.Size{Width: 75, Height: 100})
	emailRow.Add(emailInput)
	root.Add(emailRow)

	// Message field row.
	msgRow := wayne.NewRow()
	msgRow.Add(wayne.NewLabel("Message:", wayne.Size{Width: 25, Height: 100}))
	msgInput := wayne.NewTextInput("Enter a message", wayne.Size{Width: 75, Height: 100})
	msgRow.Add(msgInput)
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
