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

//		borders := fyne.NewSize(
//			sideWidth*2,
//			top.MinSize().Height,
//		)
//		return borders.AddWidthHeight(100, 100)
//	}
func setPosAndSize(top, bottom, left, right, textbox, content fyne.CanvasObject, dividers [3]fyne.CanvasObject, size fyne.Size, showRight bool) {
	topHeight := top.MinSize().Height
	bottomHeight := bottom.MinSize().Height

	leftEdgeRight := float32(sideWidth)

	top.Resize(fyne.NewSize(size.Width, topHeight))

	left.Move(fyne.NewPos(0, topHeight))
	left.Resize(fyne.NewSize(sideWidth, size.Height-topHeight))

	rightWidth := float32(0)
	if showRight {
		rightWidth = sideWidth
		right.Show()
	} else {
		right.Hide()
	}

	right.Move(fyne.NewPos(size.Width-rightWidth, topHeight))
	right.Resize(fyne.NewSize(rightWidth, size.Height-topHeight))

	content.Move(fyne.NewPos(sideWidth, topHeight))
	content.Resize(fyne.NewSize(size.Width-sideWidth-rightWidth, size.Height-topHeight))

	dividerThickness := theme.SeparatorThicknessSize()
	dividers[0].Move(fyne.NewPos(0, topHeight))
	dividers[0].Resize(fyne.NewSize(size.Width, dividerThickness))

	dividers[1].Move(fyne.NewPos(sideWidth, topHeight))
	dividers[1].Resize(fyne.NewSize(dividerThickness, size.Height-topHeight))

	dividers[2].Move(fyne.NewPos(size.Width-rightWidth, topHeight))
	dividers[2].Resize(fyne.NewSize(dividerThickness, size.Height-topHeight))

	bottom.Move(fyne.NewPos(0, size.Height-bottomHeight))
	bottom.Resize(fyne.NewSize(size.Width, bottomHeight))

	textboxHeight := textbox.MinSize().Height
	padding := float32(35)
	textbox.Move(fyne.NewPos(leftEdgeRight+padding, size.Height-bottomHeight-40))
	textbox.Resize(fyne.NewSize(size.Width, textboxHeight))
}

func makeGUI() fyne.CanvasObject {
	leftButton := widget.NewButton("Toggle Right", nil)
	right := widget.NewLabel("right")

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

	root := container.NewWithoutLayout(top, bottom, leftButton, right, textbox, content, dividers[0], dividers[1], dividers[2])

	resizeAndRefresh := func() {
		setPosAndSize(top, bottom, leftButton, right, textbox, content, dividers, root.Size(), right.Visible())
		root.Refresh()
	}

	leftButton.OnTapped = func() {
		if right.Visible() {
			right.Hide()
		} else {
			right.Show()
		}
		resizeAndRefresh()
	}

	root.Resize(fyne.NewSize(800, 600))
	resizeAndRefresh()

	return root
}

// func main() {
// 	// Initialize the app and set the custom theme
// 	app := fyne.NewApp()
// 	defer app.Quit()

// 	w := app.NewWindow("Custom Theme")
// 	w.SetContent(makeGUI())
// 	w.ShowAndRun()
// }
