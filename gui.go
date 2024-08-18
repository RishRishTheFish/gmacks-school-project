package main

import (
	"image/color"

	//	chess "onslow.collage/chess"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	// "./chess"
)

func makeBanner(color color.Color) fyne.CanvasObject {
	// Create the toolbar
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.HomeIcon(), func() {}),
	)

	// Create the logo
	logo := canvas.NewImageFromResource(resourcePfpJpg)
	logo.FillMode = canvas.ImageFillContain

	// Create the background rectangle with the specified color
	background := canvas.NewRectangle(color)

	// Layer the background, toolbar, and logo using container.NewMax
	return container.NewMax(background, toolbar, logo)
}

// func MinSize(top, objects []fyne.CanvasObject) fyne.Size {
// 	// return fyne.NewSize(10, 10)

//		borders := fyne.NewSize(
//			sideWidth*2,
//			top.MinSize().Height,
//		)
//		return borders.AddWidthHeight(100, 100)
//	}
//
// func setPosAndSize(top, bottom, left, right, textbox, content fyne.CanvasObject, options fyne.CanvasObject, dividers [3]fyne.CanvasObject, size fyne.Size, showRight bool, showOptions bool) {
func setPosAndSize(top, bottom, left, right, centerRect fyne.CanvasObject, textbox, content fyne.CanvasObject, dividers [3]fyne.CanvasObject, size fyne.Size, showRight bool, showOptions bool, options *widget.PopUp) {
	topHeight := top.MinSize().Height
	bottomHeight := bottom.MinSize().Height
	sideWidth := float32(100) // Assuming a fixed width for left and right sidebars

	// Resize top
	top.Resize(fyne.NewSize(size.Width, topHeight))

	// Position and resize left
	left.Move(fyne.NewPos(0, topHeight))
	left.Resize(fyne.NewSize(sideWidth, size.Height-topHeight-bottomHeight))

	// Handle right sidebar visibility and positioning
	rightWidth := float32(0)
	if showRight {
		rightWidth = sideWidth
		right.Show()
	} else {
		right.Hide()
	}
	right.Move(fyne.NewPos(size.Width-rightWidth, topHeight))
	right.Resize(fyne.NewSize(rightWidth, size.Height-topHeight-bottomHeight))

	// Resize content
	contentMinWidth := max(1000, size.Width-sideWidth-rightWidth)
	content.Resize(fyne.NewSize(
		float32(contentMinWidth),
		max(1000, size.Height-topHeight-bottomHeight),
	))
	content.Move(fyne.NewPos(sideWidth, topHeight))

	// Position and resize dividers
	dividerThickness := theme.SeparatorThicknessSize()
	dividers[0].Move(fyne.NewPos(0, topHeight))
	dividers[0].Resize(fyne.NewSize(size.Width, dividerThickness))

	dividers[1].Move(fyne.NewPos(sideWidth, topHeight))
	dividers[1].Resize(fyne.NewSize(dividerThickness, size.Height-topHeight-bottomHeight))

	dividers[2].Move(fyne.NewPos(
		size.Width-rightWidth-dividerThickness,
		topHeight,
	))
	dividers[2].Resize(fyne.NewSize(dividerThickness, size.Height-topHeight-bottomHeight))

	// Position and resize bottom
	bottom.Move(fyne.NewPos(0, size.Height-bottomHeight))
	bottom.Resize(fyne.NewSize(size.Width, bottomHeight))

	// Position and resize textbox
	textboxHeight := textbox.MinSize().Height
	padding := float32(35)
	textbox.Move(fyne.NewPos(sideWidth+padding, size.Height-bottomHeight-40))
	textbox.Resize(fyne.NewSize(size.Width-sideWidth-rightWidth-2*padding, textboxHeight))

	// Position and resize centerRect
	// centerRect.Move(fyne.NewPos(sideWidth, topHeight))
	// centerRect.Resize(fyne.NewSize(
	// 	size.Width-sideWidth-rightWidth,
	// 	size.Height-topHeight-bottomHeight,
	// ))

	// Debug: Output the position and size of centerRect
	//fmt.Printf("centerRect Position: %v\n", centerRect.Position())
	//fmt.Printf("centerRect Size: %v\n", centerRect.Size())

	// Handle options pop-up
	if showOptions {
		options.Show()
		optionsWidth := size.Width / 2
		optionsHeight := size.Height - topHeight - bottomHeight
		options.Resize(fyne.NewSize(optionsWidth/2, optionsHeight/2))

		// Calculate the center position for options
		centerX := (size.Width - optionsWidth) / 2
		centerY := (size.Height - optionsHeight) / 2
		options.Move(fyne.NewPos(centerX, topHeight+centerY))
	} else {
		options.Hide()
	}
}

// Function to resize and refresh the layout
func resizeAndRefresh(top fyne.CanvasObject, bottom fyne.CanvasObject, left fyne.CanvasObject, right fyne.CanvasObject, centerRect *canvas.Rectangle, textbox *fyne.Container, content fyne.CanvasObject, dividers [3]fyne.CanvasObject, size fyne.Size, isVisible bool, enableOptions bool, options *widget.PopUp, root *fyne.Container) {
	setPosAndSize(top, bottom, left, right, centerRect, textbox, content, dividers, size, isVisible, enableOptions, options)
	root.Refresh()
}

//	func returnLabel(text string) *widget.Label {
//		label := widget.NewLabel(text)
//		label.TextStyle = fyne.TextStyle{Bold: true} // Make the text bold
//		// label.TextColor = color.Black                // Set the text color to black
//		return label
//	}
func returnText(text string) *canvas.Text {
	textColor := color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
	finalText := canvas.NewText(text, textColor)
	finalText.TextStyle = fyne.TextStyle{Bold: true} // Apply bold text style
	return finalText
}

func makeGUI(w fyne.Window, CustomTheme *CustomTheme) fyne.CanvasObject {
	var enableOptions bool
	//bgColor := color.RGBA{R: 0xFF, G: 0x69, B: 0xB4, A: 0xFF}
	buttonColor := canvas.NewRectangle(colorToRGBA(CustomTheme.buttonColor))
	// canvas.NewRectangle(&color.RGBA{R: 0xFF, G: 0x69, B: 0xB4, A: 0xFF})
	//optionsColor := canvas.NewRectangle(bgColor)
	optionsColor := canvas.NewRectangle(colorToRGBA(CustomTheme.optionsColor))
	//canvas.NewRectangle(&color.RGBA{R: 0x33, G: 0x99, B: 0xFF, A: 0xFF})
	// Create a colored background rectangle
	sideBarColor := canvas.NewRectangle(colorToRGBA(CustomTheme.sideBarColor))
	// canvas.NewRectangle(&color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF}) // Example gray background color
	topColor := colorToRGBA(CustomTheme.topColor)
	//&color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF} // Example gray background color
	contentColor := colorToRGBA(CustomTheme.contentColor)
	//color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF}
	themeLabelColor := color.Gray{Y: 0x88}

	right := container.NewMax(widget.NewLabel("right"), canvas.NewRectangle(sideBarColor.FillColor)) // Dereference the pointer
	bottom := widget.NewLabel("")
	top := makeBanner(topColor)
	// Create a rectangle that fills the entire center area
	centerRect := canvas.NewRectangle(contentColor) // Green rectangle

	// Make sure the rectangle expands to fill the available space
	centerRect.Resize(fyne.NewSize(800, 600)) // Adjust this size as needed, or leave it to auto-resize

	// Create buttons
	var options *widget.PopUp
	var root *fyne.Container
	var left *fyne.Container

	singleLineEntry := widget.NewEntry()
	singleLineEntry.SetPlaceHolder("Enter text...")
	textbox := container.NewVBox(
		widget.NewLabel("Single-line Entry:"),
		singleLineEntry,
	)

	dividers := [3]fyne.CanvasObject{
		widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(),
	}

	content := canvas.NewRectangle(contentColor)

	toggleButton2 := container.NewMax(widget.NewButton("", func() {
		enableOptions = true
		resizeAndRefresh(top, bottom, left, right, centerRect, textbox, content, dividers, root.Size(), right.Visible(), enableOptions, options, root)
	}), container.NewMax(buttonColor, returnText("show options")))

	toggleButton3 := container.NewMax(widget.NewButton("", func() {
		options.Hide()
		createChess(w)
	}), container.NewMax(buttonColor, returnText("chess")))
	toggleButton4 := container.NewMax(widget.NewButton("", func() {
		options.Hide()
		w.SetContent(createTetris(w, true))
	}), container.NewMax(buttonColor, returnText("tetris-beta")))
	toggleButton5 := container.NewMax(widget.NewButton("", func() {
		options.Hide()
		w.SetContent(createSnake(w))
		// options.Hide()
		// w.SetContent(createTetris(w, false))
	}), container.NewMax(buttonColor, returnText("snake")))
	toggleButton6 := container.NewMax(widget.NewButton("", func() {
		// options.Hide()
		// w.SetContent(createSnake(w))
		options.Hide()
		w.SetContent(createTetris(w, false))
	}), container.NewMax(buttonColor, returnText("tetris-latest")))
	toggleButton1 := container.NewMax(widget.NewButton("", func() {
		if right.Visible() {
			right.Hide()
		} else {
			right.Show()
		}
		//fyne.CanvasObject
		resizeAndRefresh(top, bottom, left, right, centerRect, textbox, content, dividers, root.Size(), right.Visible(), enableOptions, options, root)
	}), container.NewMax(buttonColor, returnText("show/hide sidebar")))
	toggleButton7 := container.NewMax(widget.NewButton("", func() {
		enableOptions = false
		resizeAndRefresh(top, bottom, left, right, centerRect, textbox, content, dividers, root.Size(), right.Visible(), enableOptions, options, root)
	}), container.NewMax(buttonColor, returnText("exit options")))
	// Define initial positions
	initialPositions := map[fyne.CanvasObject]float32{
		toggleButton3: 40,
		toggleButton5: 80,
		toggleButton6: 120,
		toggleButton4: 180, // Starting position of toggleButton4
	}

	// Apply initial positions to buttons
	for btn, pos := range initialPositions {
		btn.Move(fyne.NewPos(0, pos))
	}

	// Create a slider
	slider := widget.NewSlider(0, 200)

	// Create a large spacer
	spacer := widget.NewLabel("")       // Spacer with empty content, large enough to push buttons
	spacer.Resize(fyne.NewSize(0, 200)) // Adjust size as needed

	// Store the original Y position of toggleButton4
	originalPosY := initialPositions[toggleButton4]
	optionsContent := container.NewMax(optionsColor, container.NewVBox(
		slider,
		toggleButton3,
		toggleButton5,
		toggleButton6,
		toggleButton7,
		spacer, // Add spacer to start with
		toggleButton4,
	))

	// Update positions based on slider value
	slider.OnChanged = func(value float64) {
		if value > 0 {
			// Remove spacer when slider value changes
			optionsContent.Remove(spacer)
		}

		for btn, initialY := range initialPositions {
			if btn != toggleButton4 {
				// Move buttons except toggleButton4
				newY := initialY - float32(value)
				btn.Move(fyne.NewPos(0, newY))
			}
		}

		// Move toggleButton4 separately
		if value >= 100 {
			// Move toggleButton4 only after halfway
			newY := originalPosY - (float32(value) - 50) // Adjust position relative to halfway
			toggleButton4.Move(fyne.NewPos(0, newY))
		} else {
			// Return toggleButton4 to original position
			toggleButton4.Move(fyne.NewPos(0, originalPosY))
		}
	}

	// Create content for the options
	options = widget.NewModalPopUp(
		optionsContent,
		w.Canvas(),
	)
	options.Hide() // Ensure options is hidden initially

	// Determine the current system theme (Light or Dark)
	var themeName string
	if fyne.CurrentApp().Settings().Theme() == theme.LightTheme() {
		themeName = "Light"
	} else {
		themeName = "Dark"
	}

	// Create an expanding spacer
	spacer2 := layout.NewSpacer()

	// Create the square with a label to display the system theme
	themeLabel := widget.NewLabel(themeName)
	square := canvas.NewRectangle(themeLabelColor)
	square.Resize(fyne.NewSize(50, 50)) // Adjust the size of the square as needed

	themeContainer := container.NewCenter(themeLabel)       // Center the label inside the square
	themeSquare := container.NewMax(square, themeContainer) // Overlay the label on the square

	// Position the square at the bottom left by adding it after the spacer
	// left = container.NewMax(container.NewVBox(
	// 	widget.NewLabel("Buttons:"),
	// 	toggleButton1,
	// 	toggleButton2,
	// 	spacer2,     // This spacer will take up all the space, pushing the themeSquare to the bottom
	// 	themeSquare, // The square with the theme label
	// ), bgColor)
	leftContent := container.NewVBox(
		widget.NewLabel("Buttons:"),
		toggleButton1,
		toggleButton2,
		spacer2,     // This spacer will take up all the space, pushing the themeSquare to the bottom
		themeSquare, // The square with the theme label
	)

	// Layer the background and the VBox using container.NewMax
	left = container.NewMax(
		sideBarColor,
		leftContent,
	)

	// Update the root container to use the rectangle in the center
	//dividers[0], dividers[1], dividers[2]
	// root = container.NewBorder(top, bottom, left, right, centerRect, options)
	root = container.NewBorder(top, bottom, left, right, options, dividers[0], dividers[1], dividers[2], centerRect)

	enableOptions = false

	root.Resize(fyne.NewSize(800, 600))
	resizeAndRefresh(top, bottom, left, right, centerRect, textbox, content, dividers, root.Size(), right.Visible(), enableOptions, options, root)

	// Declare the welcomeModal variable
	var welcomeModal *widget.PopUp

	// Create the Welcome modal pop-up content
	welcomeLabel := widget.NewLabel("Welcome to the App!")
	closeButton := widget.NewButton("Close", func() {
		welcomeModal.Hide()
	})

	// Create a vertical box to hold the label and close button
	welcomeContent := container.NewVBox(
		welcomeLabel,
		closeButton,
	)

	// Create the modal pop-up
	welcomeModal = widget.NewModalPopUp(
		welcomeContent,
		w.Canvas(),
	)

	// Show the welcome modal when the app starts
	// welcomeModal.Show()

	return root
}

// Define button actions
// toggleButton2.OnTapped = func() {
// 	// enableOptions = true
// 	// resizeAndRefresh(top, bottom, left, right, textbox, content, dividers, root.Size(), right.Visible(), enableOptions, options, root)
// }
// toggleButton5.OnTapped = func() {
// 	options.Hide()
// 	w.SetContent(createTetris(w, false))
// }
// toggleButton1.OnTapped = func() {
// 	if right.Visible() {
// 		right.Hide()
// 	} else {
// 		right.Show()
// 	}
// 	resizeAndRefresh()
// }

// toggleButton3.OnTapped = func() {
// 	options.Hide()
// 	createChess(w)
// }
// toggleButton4.OnTapped = func() {
// 	options.Hide()
// 	w.SetContent(createTetris(w, true))
// }
// toggleButton6.OnTapped = func() {
// 	options.Hide()
// 	w.SetContent(createSnake(w))

// }

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
// func main() {
// 	// Initialize the app and set the custom theme
// 	app := fyne.NewApp()
// 	defer app.Quit()

// 	w := app.NewWindow("Custom Theme")
// 	w.SetContent(makeGUI())
// 	w.ShowAndRun()
// }
