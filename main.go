//go:generate ./bin/fyne bundle -o data.go Icon.png

package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	var bgColor *color.RGBA
	theme := NewCustomTheme(bgColor, nil)
	// theme.backgroundColor = color.RGBA{0, 128, 0, 255}
	//theme := &CustomTheme{}
	a := app.New()
	a.Settings().SetTheme(theme)
	// a.Settings().SetTheme(newTheme())
	w := a.NewWindow("app")
	w.Resize(fyne.NewSize(1034, 768))

	w.SetContent(makeGUI(w, bgColor))
	// w.SetContent(widget.NewLabel("app"))
	w.ShowAndRun()

}
