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

// func MinSize(top, objects []fyne.CanvasObject) fyne.Size {
// 	// return fyne.NewSize(10, 10)

// 	borders := fyne.NewSize(
// 		sideWidth*2,
// 		top.MinSize().Height,
// 	)
// 	return borders.AddWidthHeight(100, 100)
// }
// func setPosAndSize(top, bottom, left, right, textbox, content fyne.CanvasObject, dividers [3]fyne.CanvasObject) {
// 	topHeight := top.MinSize().Height
// 	bottomHeight := bottom.MinSize().Height

// 	leftEdgeRight := float32(sideWidth)
// 	// rightEdgeLeft := size.Width - float32(sideWidth)
// 	// middleX := (leftEdgeRight + rightEdgeLeft) / 2

// 	top.Resize(fyne.NewSize(size.Width, topHeight))

// 	left.Move(fyne.NewPos(0, topHeight))
// 	left.Resize(fyne.NewSize(sideWidth, size.Height-topHeight))

// 	// l.right.Move(fyne.NewPos(0, 0))
// 	// l.right.Resize(fyne.NewSize(0, 0))
// 	right.Move(fyne.NewPos(size.Width-sideWidth, topHeight))
// 	right.Resize(fyne.NewSize(sideWidth, size.Height-topHeight))
// 	// l.right.Visible()
// 	right.Hide()

// 	content.Move(fyne.NewPos(sideWidth, topHeight))
// 	content.Resize(fyne.NewSize(size.Width-2*sideWidth, size.Height-topHeight))

// 	dividerThickness := theme.SeparatorThicknessSize()
// 	dividers[0].Move(fyne.NewPos(0, topHeight))
// 	dividers[0].Resize(fyne.NewSize(size.Width, dividerThickness))

// 	dividers[1].Move(fyne.NewPos(sideWidth, topHeight))
// 	dividers[1].Resize(fyne.NewSize(dividerThickness, size.Height-topHeight))

// 	dividers[2].Move(fyne.NewPos(size.Width-sideWidth, topHeight))
// 	dividers[2].Resize(fyne.NewSize(dividerThickness, size.Height-topHeight))

// 	// Resize and position the bottom component
// 	bottom.Move(fyne.NewPos(0, size.Height-bottomHeight))
// 	bottom.Resize(fyne.NewSize(size.Width, bottomHeight))

// 	// Resize and position the textbox component within the bottom area
// 	// textboxHeight := l.textbox.MinSize().Height
// 	// l.textbox.Move(fyne.NewPos(middleX/3, size.Height-bottomHeight-40))
// 	// l.textbox.Resize(fyne.NewSize(size.Width, textboxHeight))
// 	textboxHeight := l.textbox.MinSize().Height
// 	padding := float32(35) // adjust this value as needed for padding
// 	textbox.Move(fyne.NewPos(leftEdgeRight+padding, size.Height-bottomHeight-40))
// 	textbox.Resize(fyne.NewSize(size.Width, textboxHeight))
// }

func makeGUI() fyne.CanvasObject {
	left := widget.NewLabel("left")
	right := widget.NewLabel("right")
	// test := widget.NewLabel("right")
	// Create a single-line text box
	singleLineEntry := widget.NewEntry()
	singleLineEntry.SetPlaceHolder("Enter text...")
	textbox := container.NewVBox(
		widget.NewLabel("Single-line Entry:"),
		singleLineEntry,
	)
	bottom := widget.NewLabel("")
	top := makeBanner()

	content := canvas.NewRectangle(color.Gray{Y: 0xee})

	dividers := [3]fyne.CanvasObject{
		widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(),
	}
	// setPosAndSize(top, bottom, left, right, textbox, content, dividers[0], dividers[1], dividers[2])
	// right.Hide()
	// appList := widget.NewList()
	objs := []fyne.CanvasObject{content, top, bottom, left, right, textbox, dividers[0], dividers[1], dividers[2]}
	//return container.NewBorder(top, bottom, left, right, textbox, content, dividers[0], dividers[1], dividers[2], container.Size())
	return container.New(newLayout(top, bottom, left, right, textbox, content, dividers), objs...)
}

// func main() {
// 	// Initialize the app and set the custom theme
// 	app := fyne.NewApp()
// 	defer app.Quit()

// 	w := app.NewWindow("Custom Theme")
// 	w.SetContent(makeGUI())
// 	w.ShowAndRun()
// }
