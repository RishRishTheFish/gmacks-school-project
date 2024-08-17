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
//
// func setPosAndSize(top, bottom, left, right, textbox, content fyne.CanvasObject, options fyne.CanvasObject, dividers [3]fyne.CanvasObject, size fyne.Size, showRight bool, showOptions bool) {
func setPosAndSize(top, bottom, left, right, textbox, content fyne.CanvasObject, dividers [3]fyne.CanvasObject, size fyne.Size, showRight bool, showOptions bool, options *widget.PopUp) {
	topHeight := top.MinSize().Height
	bottomHeight := bottom.MinSize().Height

	leftEdgeRight := float32(sideWidth)

	top.Resize(fyne.NewSize(size.Width, topHeight))

	left.Move(fyne.NewPos(0, topHeight))
	left.Resize(fyne.NewSize(sideWidth, size.Height-topHeight))

	// right.Hide()

	rightWidth := float32(0)
	if showRight {
		rightWidth = sideWidth
		right.Show()
	} else {
		right.Hide()
	}

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

	dividerThickness := theme.SeparatorThicknessSize()
	dividers[0].Move(fyne.NewPos(0, topHeight))
	dividers[0].Resize(fyne.NewSize(size.Width, dividerThickness))

	dividers[1].Move(fyne.NewPos(sideWidth, topHeight))
	dividers[1].Resize(fyne.NewSize(dividerThickness, size.Height-topHeight))
	//1465
	//1237
	// dividers[2].Move(fyne.NewPos(size.Width-rightWidth, topHeight))
	// fmt.Println(size.Width - rightWidth)
	dividers[2].Move(fyne.NewPos(
		max(1465, size.Width-rightWidth),
		topHeight,
	))
	dividers[2].Resize(fyne.NewSize(dividerThickness, size.Height-topHeight))

	bottom.Move(fyne.NewPos(0, size.Height-bottomHeight))
	bottom.Resize(fyne.NewSize(size.Width, bottomHeight))

	textboxHeight := textbox.MinSize().Height
	padding := float32(35)
	textbox.Move(fyne.NewPos(leftEdgeRight+padding, size.Height-bottomHeight-40))
	textbox.Resize(fyne.NewSize(size.Width, textboxHeight))
}
func makeGUI(w fyne.Window) fyne.CanvasObject {
	// Create buttons
	toggleButton1 := widget.NewButton("Toggle Right", nil)
	toggleButton2 := widget.NewButton("Show options", nil)
	toggleButton3 := widget.NewButton("Chess", nil)
	toggleButton4 := widget.NewButton("Tetris-latest", nil)
	toggleButton5 := widget.NewButton("Tetris-BETA", nil)
	toggleButton6 := widget.NewButton("Snake", nil)

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
	optionsContent := container.NewVBox(
		slider,
		toggleButton3,
		toggleButton5,
		toggleButton6,
		spacer, // Add spacer to start with
		toggleButton4,
	)

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
	options := widget.NewModalPopUp(
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
	square := canvas.NewRectangle(color.Gray{Y: 0x88})
	square.Resize(fyne.NewSize(50, 50)) // Adjust the size of the square as needed

	themeContainer := container.NewCenter(themeLabel)       // Center the label inside the square
	themeSquare := container.NewMax(square, themeContainer) // Overlay the label on the square

	// Position the square at the bottom left by adding it after the spacer
	left := container.NewVBox(
		widget.NewLabel("Buttons:"),
		toggleButton1,
		toggleButton2,
		spacer2,     // This spacer will take up all the space, pushing the themeSquare to the bottom
		themeSquare, // The square with the theme label
	)

	right := widget.NewLabel("right") // Placeholder for the right section

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

	// Create the main container
	root := container.NewBorder(top, bottom, left, right, textbox, content, options, dividers[0], dividers[1], dividers[2])

	enableOptions := false
	// Function to resize and refresh the layout
	resizeAndRefresh := func() {
		setPosAndSize(top, bottom, left, right, textbox, content, dividers, root.Size(), right.Visible(), enableOptions, options)
		root.Refresh()
	}

	root.Resize(fyne.NewSize(800, 600))
	resizeAndRefresh()

	// Define button actions
	toggleButton2.OnTapped = func() {
		enableOptions = true
		resizeAndRefresh()
	}
	toggleButton5.OnTapped = func() {
		options.Hide()
		w.SetContent(createTetris(w, false))
	}
	toggleButton1.OnTapped = func() {
		if right.Visible() {
			right.Hide()
		} else {
			right.Show()
		}
		resizeAndRefresh()
	}

	toggleButton3.OnTapped = func() {
		options.Hide()
		createChess(w)
	}
	toggleButton4.OnTapped = func() {
		options.Hide()
		w.SetContent(createTetris(w, true))
	}
	toggleButton6.OnTapped = func() {
		options.Hide()
		w.SetContent(createSnake(w))

	}

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
