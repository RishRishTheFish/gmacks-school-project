package main

import (
	"fmt"
	"image/color"

	//	chess "onslow.collage/chess"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	// "./chess"
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
<<<<<<< HEAD
//
// func setPosAndSize(top, bottom, left, right, textbox, content fyne.CanvasObject, options fyne.CanvasObject, dividers [3]fyne.CanvasObject, size fyne.Size, showRight bool, showOptions bool) {
func setPosAndSize(top, bottom, left, right, textbox, content fyne.CanvasObject, dividers [3]fyne.CanvasObject, size fyne.Size, showRight bool, showOptions bool, options *widget.PopUp) {
=======
func setPosAndSize(top, bottom, left, right, textbox, content fyne.CanvasObject, dividers [3]fyne.CanvasObject, size fyne.Size, showRight bool) {
>>>>>>> 0238b6a (without layout i moved everything to a container)
	topHeight := top.MinSize().Height
	bottomHeight := bottom.MinSize().Height

	leftEdgeRight := float32(sideWidth)

	top.Resize(fyne.NewSize(size.Width, topHeight))

	left.Move(fyne.NewPos(0, topHeight))
	left.Resize(fyne.NewSize(sideWidth, size.Height-topHeight))

<<<<<<< HEAD
	// right.Hide()

=======
>>>>>>> 0238b6a (without layout i moved everything to a container)
	rightWidth := float32(0)
	if showRight {
		rightWidth = sideWidth
		right.Show()
	} else {
		right.Hide()
	}

<<<<<<< HEAD
	if showOptions {
		options.Show()
		optionsWidth := size.Width / 2
		optionsHeight := size.Height - topHeight - bottomHeight
		options.Resize(fyne.NewSize(optionsWidth/2, optionsHeight/2))

		// Calculate the center position
		centerX := (size.Width - optionsWidth) / 2
		centerY := (size.Height - optionsHeight) / 2
		options.Move(fyne.NewPos(centerX, topHeight+centerY))
	} else {
		options.Hide()
	}

	right.Move(fyne.NewPos(size.Width-rightWidth, topHeight))
	right.Resize(fyne.NewSize(rightWidth, size.Height-topHeight))

	contentMinWidth := max(1000, size.Width-sideWidth-rightWidth)
	// if
	// fmt.Println(size.Width - sideWidth - rightWidth)

	content.Move(fyne.NewPos(sideWidth, topHeight))
	// content.Resize(fyne.NewSize(size.Width-sideWidth-rightWidth, size.Height-topHeight))
	content.Resize(fyne.NewSize(
		float32(contentMinWidth),
		max(1000, size.Height-topHeight),
	))
=======
	right.Move(fyne.NewPos(size.Width-rightWidth, topHeight))
	right.Resize(fyne.NewSize(rightWidth, size.Height-topHeight))

	content.Move(fyne.NewPos(sideWidth, topHeight))
	content.Resize(fyne.NewSize(size.Width-sideWidth-rightWidth, size.Height-topHeight))
>>>>>>> 0238b6a (without layout i moved everything to a container)

	dividerThickness := theme.SeparatorThicknessSize()
	dividers[0].Move(fyne.NewPos(0, topHeight))
	dividers[0].Resize(fyne.NewSize(size.Width, dividerThickness))

	dividers[1].Move(fyne.NewPos(sideWidth, topHeight))
	dividers[1].Resize(fyne.NewSize(dividerThickness, size.Height-topHeight))
<<<<<<< HEAD
	//1465
	//1237
	// dividers[2].Move(fyne.NewPos(size.Width-rightWidth, topHeight))
	// fmt.Println(size.Width - rightWidth)
	dividers[2].Move(fyne.NewPos(
		max(1465, size.Width-rightWidth),
		topHeight,
	))
=======

	dividers[2].Move(fyne.NewPos(size.Width-rightWidth, topHeight))
>>>>>>> 0238b6a (without layout i moved everything to a container)
	dividers[2].Resize(fyne.NewSize(dividerThickness, size.Height-topHeight))

	bottom.Move(fyne.NewPos(0, size.Height-bottomHeight))
	bottom.Resize(fyne.NewSize(size.Width, bottomHeight))

	textboxHeight := textbox.MinSize().Height
	padding := float32(35)
	textbox.Move(fyne.NewPos(leftEdgeRight+padding, size.Height-bottomHeight-40))
	textbox.Resize(fyne.NewSize(size.Width, textboxHeight))
}

<<<<<<< HEAD
func makeGUI(w fyne.Window) fyne.CanvasObject {
	fmt.Println()
	toggleButton1 := widget.NewButton("Toggle Right", nil)
	toggleButton2 := widget.NewButton("Show options", nil)
	toggleButton3 := widget.NewButton("Chess", nil)
	toggleButton4 := widget.NewButton("Tetris", nil)

	left := container.NewVBox(
		widget.NewLabel("Buttons:"),
		toggleButton1,
		toggleButton2,
	)

	enableOptions := false

=======
func makeGUI() fyne.CanvasObject {
	leftButton := widget.NewButton("Toggle Right", nil)
>>>>>>> 0238b6a (without layout i moved everything to a container)
	right := widget.NewLabel("right")
	options := widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabel("First option"),
			toggleButton3,
			toggleButton4,
		),
		w.Canvas(),
	)
	options.Hide() // Ensure options is hidden initially

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

<<<<<<< HEAD
	root := container.NewWithoutLayout(top, bottom, left, right, textbox, content, options, dividers[0], dividers[1], dividers[2])
	// root.Refresh()
	resizeAndRefresh := func() {
		setPosAndSize(top, bottom, left, right, textbox, content, dividers, root.Size(), right.Visible(), enableOptions, options)
		root.Refresh()
	}
	// right.Hide()
	// resizeAndRefresh()
	// right.Show()
	// resizeAndRefresh()
	// right.Show()

	root.Resize(fyne.NewSize(800, 600))
	resizeAndRefresh()

	toggleButton2.OnTapped = func() {
		enableOptions = true
		resizeAndRefresh()
	}

	toggleButton1.OnTapped = func() {
=======
	root := container.NewWithoutLayout(top, bottom, leftButton, right, textbox, content, dividers[0], dividers[1], dividers[2])

	resizeAndRefresh := func() {
		setPosAndSize(top, bottom, leftButton, right, textbox, content, dividers, root.Size(), right.Visible())
		root.Refresh()
	}

	leftButton.OnTapped = func() {
>>>>>>> 0238b6a (without layout i moved everything to a container)
		if right.Visible() {
			right.Hide()
		} else {
			right.Show()
		}
		resizeAndRefresh()
	}

<<<<<<< HEAD
	toggleButton3.OnTapped = func() {
		options.Hide()
		w.Resize(fyne.NewSize(480, 480))
		w.SetContent(createGrid())
	}
	toggleButton4.OnTapped = func() {
		options.Hide()
		w.SetContent(createTetris())
	}
	return root
}

/*

	// Use custom layout
	return fyne.NewContainerWithLayout(
		&customLayout{
			top:     top,
			bottom:  bottom,
			left:    left,
			right:   right,
			textbox: textbox,
			content: content,
			dividers: [3]fyne.CanvasObject{
				dividers[0], dividers[1], dividers[2],
			},
		},
		top, bottom, left, right, textbox, content, dividers[0], dividers[1], dividers[2],
	)
*/
=======
	root.Resize(fyne.NewSize(800, 600))
	resizeAndRefresh()

	return root
}

>>>>>>> 0238b6a (without layout i moved everything to a container)
// func main() {
// 	// Initialize the app and set the custom theme
// 	app := fyne.NewApp()
// 	defer app.Quit()

// 	w := app.NewWindow("Custom Theme")
// 	w.SetContent(makeGUI())
// 	w.ShowAndRun()
// }
