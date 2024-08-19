package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {

	theme := NewCustomTheme(
		&color.RGBA{R: 0x01, G: 0x01, B: 0x02, A: 0xCC},
		&color.RGBA{R: 0x0C, G: 0x06, B: 0x2E, A: 0xFF},
		&color.RGBA{R: 0x23, G: 0x14, B: 0x3C, A: 0xFF},
		&color.RGBA{R: 0x67, G: 0xB3, B: 0xB5, A: 0xFF},
		&color.RGBA{R: 0xF3, G: 0x19, B: 0x19, A: 0xFF},
		&color.RGBA{R: 0xF5, G: 0x9B, B: 0x00, A: 0xFF},

		&color.RGBA{R: 0x40, G: 0xE0, B: 0xD0, A: 0xFF},
		&color.RGBA{R: 0xD0, G: 0xD0, B: 0xD0, A: 0xFF},
	)

	a := app.New()
	a.Settings().SetTheme(theme)

	w := a.NewWindow("app")
	w.Resize(fyne.NewSize(1034, 768))

	w.SetContent(makeGUI(w, theme))

	w.ShowAndRun()

}
