//go:generate ./bin/fyne bundle -o data.go ./assets/pfp.png

package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	//&color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}, // Background color
	// theme := NewCustomTheme(
	// 	// nil,
	// 	// nil,
	// 	// nil,
	// 	// nil,
	// 	// nil,
	// 	// nil,
	// 	// nil,
	// 	// nil,
	// 	theme := NewCustomTheme(
	// 		&color.RGBA{R: 0x2A, G: 0x2A, B: 0x2A, A: 0xCC}, // Background color: soft dark gray with transparency
	// 		&color.RGBA{R: 0xB5, G: 0x00, B: 0x00, A: 0xFF}, // Sidebar color: muted red
	// 		&color.RGBA{R: 0x00, G: 0xB5, B: 0x00, A: 0xFF}, // Top color: muted green
	// 		&color.RGBA{R: 0x00, G: 0x00, B: 0xB5, A: 0xFF}, // Content color: softer blue
	// 		&color.RGBA{R: 0xE0, G: 0xE0, B: 0x70, A: 0xFF}, // Label color: soft yellow
	// 		&color.RGBA{R: 0xF0, G: 0xB5, B: 0x70, A: 0xFF}, // Options color: softer orange
	// 		&color.RGBA{R: 0xFF, G: 0xA8, B: 0xB5, A: 0xFF}, // Button color: soft pink
	// 		&color.RGBA{R: 0xD0, G: 0xD0, B: 0xD0, A: 0xFF}, // Foreground color: light gray
	// 	)

	// )
	theme := NewCustomTheme(
		&color.RGBA{R: 0x2A, G: 0x2A, B: 0x2A, A: 0xCC}, // Background color: dark gray with transparency
		&color.RGBA{R: 0xD0, G: 0xD0, B: 0xD0, A: 0xFF}, // Sidebar color: light gray
		&color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}, // Top bar color: black
		&color.RGBA{R: 0x00, G: 0x00, B: 0x80, A: 0xFF}, // Content color: muted blue
		&color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}, // Label color: bright white for high contrast
		&color.RGBA{R: 0xF5, G: 0x9B, B: 0x00, A: 0xFF}, // Options color: muted orange
		&color.RGBA{R: 0x00, G: 0xFF, B: 0xFF, A: 0xFF}, // Button color: bright cyan
		&color.RGBA{R: 0xD0, G: 0xD0, B: 0xD0, A: 0xFF}, // Foreground color: light gray
	)

	// theme.backgroundColor = color.RGBA{0, 128, 0, 255}
	//theme := &CustomTheme{}
	a := app.New()
	a.Settings().SetTheme(theme)
	// a.Settings().SetTheme(newTheme())
	w := a.NewWindow("app")
	w.Resize(fyne.NewSize(1034, 768))

	w.SetContent(makeGUI(w, theme))
	// w.SetContent(widget.NewLabel("app"))
	w.ShowAndRun()

}
