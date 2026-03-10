//go:build windows || darwin || android || ios

// Package wayne provides a 100% API-compatible alternative to the opd-ai/wain
// widget system, targeting Windows, macOS, Android, and iOS.
//
// Wayne uses Ebitengine (github.com/hajimehoshi/ebiten/v2) as its rendering
// and windowing backend, replacing wain's Linux-specific Wayland/X11 backend.
//
// # Quick Start
//
//	app := wayne.NewApp()
//	defer app.Close()
//
//	win, _ := app.NewWindow(wayne.WindowConfig{Title: "Hello", Width: 800, Height: 600})
//	win.Show()
//
//	btn := wayne.NewButton("Click me", wayne.Size{Width: 30, Height: 10})
//	btn.OnClick(func() { fmt.Println("clicked!") })
//
//	col := wayne.NewColumn()
//	col.Add(btn)
//	win.SetRoot(col)
//
//	app.Run()
//
// # Widget Sizing Convention
//
// Widgets with text content (Button, Label, TextInput) require explicit Size
// parameters. Spacers and containers (Panel, ScrollView, ImageWidget) also
// require explicit Size. Container widgets (Row, Column, Stack, Grid) use
// percentage-based defaults of 100x100 (full parent dimensions).
//
// All Size values are percentages (0-100) of the parent container. For example,
// Size{Width: 50, Height: 100} means "50% of parent width, 100% of parent height".
//
// # Supported Platforms
//
// Wayne compiles on Windows, macOS, Android, and iOS.
// It explicitly does NOT support Linux or BSD.
//
// # Architecture
//
// Wayne's App type wraps an Ebitengine game loop (ebiten.RunGame).
// The widget tree is resolved and rendered each frame via Ebitengine's
// Draw callback, with layout computed from percentage-based Size values.
package wayne
