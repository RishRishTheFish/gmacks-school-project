package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func makeBanner() fyne.CanvasObject {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.HomeIcon(), func() {}),
	)
	logo := canvas.NewImageFromResource(resourcePfpJpg)
	logo.FillMode = canvas.ImageFillContain

	return container.NewMax(toolbar, logo)
}

func makeGUI() fyne.CanvasObject {
	left := widget.NewLabel("left")
	right := widget.NewLabel("right")
	top := makeBanner()

	content := canvas.NewRectangle(color.Gray{Y: 0xee})

	dividers := [3]fyne.CanvasObject{
		widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(),
	}
	objs := []fyne.CanvasObject{content, top, left, right, dividers[0], dividers[1], dividers[2]}
	// return container.NewBorder(makeBanner(), nil, left, right, content)
	return container.New(newLayout(top, left, right, content, dividers), objs...)
}

// func main() {
// 	// Initialize the app and set the custom theme
// 	app := fyne.NewApp()
// 	defer app.Quit()

// 	w := app.NewWindow("Custom Theme")
// 	w.SetContent(makeGUI())
// 	w.ShowAndRun()
// }
