package main

import (
	"image/color"

	//	chess "onslow.collage/chess"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	// "fyne.io/fyne/v2/key"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	// "./chess"
)

var tetrisScores []int
var chessScores []int
var snakeScores []int

func makeBanner(color color.Color) fyne.CanvasObject {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.HomeIcon(), func() {}),
	)

	logo := canvas.NewImageFromResource(resourcePfpJpg)
	logo.FillMode = canvas.ImageFillContain

	background := canvas.NewRectangle(color)

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
	sideWidth := float32(100)

	top.Resize(fyne.NewSize(size.Width, topHeight))

	// Position and resize left
	left.Move(fyne.NewPos(0, topHeight))
	left.Resize(fyne.NewSize(sideWidth, size.Height-topHeight-bottomHeight))

	rightWidth := float32(0)
	if showRight {
		rightWidth = sideWidth
		right.Show()
	} else {
		right.Hide()
	}
	right.Move(fyne.NewPos(size.Width-rightWidth, topHeight))
	right.Resize(fyne.NewSize(rightWidth, size.Height-topHeight-bottomHeight))

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

// func handleEnterKey(w fyne.Window, entry *widget.Entry, ev *fyne.KeyEvent) {
// 	if ev.Name == fyne.KeyEnter {
// 		// Get the text from the entry widget
// 		text := entry.Text
// 		fmt.Println("Text entered:", text)

// 		// You can also show the text in a dialog or use it as needed
// 		dialog.ShowInformation("Entered Text", text, w)
// 	}
// }

// Function to resize and refresh the layout
func resizeAndRefresh(top fyne.CanvasObject, bottom fyne.CanvasObject, left fyne.CanvasObject, right fyne.CanvasObject, center *fyne.Container, textbox *fyne.Container, content fyne.CanvasObject, dividers [3]fyne.CanvasObject, size fyne.Size, isVisible bool, enableOptions bool, options *widget.PopUp, root *fyne.Container) {
	setPosAndSize(top, bottom, left, right, center, textbox, content, dividers, size, isVisible, enableOptions, options)
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

var buttonColor *canvas.Rectangle

func makeGUI(w fyne.Window, CustomTheme *CustomTheme) fyne.CanvasObject {
	var enableOptions bool
	allowBeta := false
	//bgColor := color.RGBA{R: 0xFF, G: 0x69, B: 0xB4, A: 0xFF}
	buttonColor = canvas.NewRectangle(colorToRGBA(CustomTheme.buttonColor))
	// globalButtonColor = globalButtonColor
	// canvas.NewRectangle(&color.RGBA{R: 0xFF, G: 0x69, B: 0xB4, A: 0xFF})
	//optionsColor := canvas.NewRectangle(bgColor)
	optionsColor := canvas.NewRectangle(colorToRGBA(CustomTheme.optionsColor))
	//canvas.NewRectangle(&color.RGBA{R: 0x33, G: 0x99, B: 0xFF, A: 0xFF})
	// Create a colored background rectangle
	sideBarColor := canvas.NewRectangle(colorToRGBA(CustomTheme.sideBarColor))
	// canvas.NewRectangle(&color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF})
	topColor := colorToRGBA(CustomTheme.topColor)
	//&color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF}
	contentColor := colorToRGBA(CustomTheme.contentColor)
	//color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF}
	themeLabelColor := color.Gray{Y: 0x88}

	right := container.NewMax(widget.NewLabel("right"), canvas.NewRectangle(sideBarColor.FillColor))
	bottom := widget.NewLabel("")
	top := makeBanner(topColor)

	scoreBoardRect := *buttonColor
	scoreBoardRect.FillColor = buttonColor.FillColor
	scoreBoardRect.SetMinSize(fyne.NewSize(200, 300))

	centerRect := canvas.NewRectangle(contentColor)

	commandRect := *buttonColor
	commandRect.FillColor = buttonColor.FillColor
	commandRect.SetMinSize(fyne.NewSize(centerRect.MinSize().Width, 50))
	// Create the entry widget
	entry := widget.NewEntry()

	// Define a function to handle the Enter key press
	// handleEnterKey := func(ev *fyne.KeyEvent) {
	// 	if ev.Name == fyne.KeyEnter {
	// 		// Get the text from the entry widget
	// 		text := entry.Text
	// 		fmt.Println("Text entered:", text)

	// 		// You can also show the text in a dialog or use it as needed
	// 		dialog.ShowInformation("Entered Text", text, w)
	// 	}
	// }
	// w.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
	// 	//handleEnterKey(w, entry, ev)
	// 	handleEnterKey(ev)
	// })

	center := container.NewMax(centerRect, container.NewVBox(
		// container.NewHBox(
		// 	layout.NewSpacer(),
		// 	container.NewMax(
		// 		&scoreBoardRect,
		// 		container.NewVBox(
		// 			widget.NewLabel("test"),
		// 		),
		// 	),
		// 	layout.NewSpacer(),
		// 	container.NewMax(
		// 		&scoreBoardRect,
		// 		container.NewVBox(
		// 			widget.NewLabel("test"),
		// 		),
		// 	),
		// 	layout.NewSpacer(),
		// 	container.NewMax(
		// 		&scoreBoardRect,
		// 		container.NewVBox(
		// 			widget.NewLabel("test"),
		// 		),
		// 	),
		// 	layout.NewSpacer(),
		// ),
		layout.NewSpacer(),
		container.NewMax(&commandRect, container.NewVBox(
			widget.NewLabel("Type a message here!"),
			entry,
		)),
	),
	// layout.NewSpacer(),
	// container.NewHBox(
	// 	layout.NewSpacer(), // Center the command rectangle
	// 	container.NewVBox(
	// 		container.NewMax(&commandRect, widget.NewLabel("command, rect")),
	// 	),
	// 	layout.NewSpacer(),
	// ),
	)

	// overlay := container.NewStac
	// Make sure the rectangle expands to fill the available space
	centerRect.Resize(fyne.NewSize(800, 600)) // Adjust this size as needed, or leave it to auto-resize
	entry.OnSubmitted = func(content string) {
		// Create a new label with the submitted content
		newMessage := widget.NewLabel(content)

		// Add the new label to the VBox in the center container
		center.Objects[1].(*fyne.Container).Add(newMessage)

		// Optionally clear the entry after submission
		entry.SetText("")
	}
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

	cantRunBackroundSquare := canvas.NewRectangle(contentColor)
	cantRunBackroundSquare.SetMinSize(fyne.NewSize(300, 150))
	cantRunVboxBackround := *buttonColor
	cantRunVboxBackround.SetMinSize(fyne.NewSize(250, 125))
	var cantRun *widget.PopUp
	cantRunVbox := container.NewMax(
		&cantRunVboxBackround,
		// container.NewFixedSize
		// fyne.NewSize(300, 150),
		container.NewVBox(
			// widget.NewLabel("You cannot run this"),
			container.NewMax(buttonColor, returnText("Sorry! this functionality isnt avalible at the moment")),

			container.NewMax(widget.NewButton("", func() {
				cantRun.Hide()
			}), container.NewMax(buttonColor, returnText("Click here to exit"))),
		))

	cantRun = widget.NewModalPopUp(container.NewMax(
		cantRunBackroundSquare,
		container.NewCenter(cantRunVbox)),
		w.Canvas(),
	)
	// cantDoPvp = widget.NewModalPopUp(
	// 	container.NewMax(
	// 		optionsColor, container.NewVBox(
	// 			returnText("Sorry! the functionality is not avalible at the moment"),
	// 		),
	// 	),
	// 	w.Canvas(),
	// )
	var chessMode *widget.PopUp
	toggleButton8 := container.NewMax(widget.NewButton("", func() {
		options.Hide()
		chessMode.Hide()
		createChess(w, true, CustomTheme)
	}), container.NewMax(buttonColor, returnText("chess with pve")))
	toggleButton9 := container.NewMax(widget.NewButton("", func() {
		// options.Hide()
		chessMode.Hide()
		cantRun.Show()
		//createChess(w, false, CustomTheme)
	}), container.NewMax(buttonColor, returnText("chess with pvp")))
	chessMode = widget.NewModalPopUp(
		container.NewMax(
			optionsColor,
			container.NewVBox(
				toggleButton8,
				toggleButton9,
			),
		),
		w.Canvas(),
	)
	toggleButton2 := container.NewMax(widget.NewButton("", func() {
		enableOptions = true
		resizeAndRefresh(top, bottom, left, right, center, textbox, content, dividers, root.Size(), right.Visible(), enableOptions, options, root)
	}), container.NewMax(buttonColor, returnText("show options")))
	toggleButton3 := container.NewMax(widget.NewButton("", func() {
		options.Hide()
		chessMode.Show()
		// createChess(w, true, CustomTheme)
	}), container.NewMax(buttonColor, returnText("chess")))

	toggleButton4 := container.NewMax(widget.NewButton("", func() {
		if allowBeta {
			options.Hide()
			w.SetContent(createTetris(w, CustomTheme, true))
		} else {
			cantRun.Show()
		}
	}), container.NewMax(buttonColor, returnText("tetris-beta")))
	toggleButton5 := container.NewMax(widget.NewButton("", func() {
		options.Hide()
		cantRun.Show()
		//w.SetContent(createSnake(w, CustomTheme))
		// options.Hide()
		// w.SetContent(createTetris(w, false))
	}), container.NewMax(buttonColor, returnText("snake")))
	toggleButton6 := container.NewMax(widget.NewButton("", func() {
		// options.Hide()
		// w.SetContent(createSnake(w))
		options.Hide()
		w.SetContent(createTetris(w, CustomTheme, false))
	}), container.NewMax(buttonColor, returnText("tetris-latest")))
	toggleButton1 := container.NewMax(widget.NewButton("", func() {
		if right.Visible() {
			right.Hide()
		} else {
			right.Show()
		}
		//fyne.CanvasObject
		resizeAndRefresh(top, bottom, left, right, center, textbox, content, dividers, root.Size(), right.Visible(), enableOptions, options, root)
	}), container.NewMax(buttonColor, returnText("show/hide sidebar")))
	toggleButton7 := container.NewMax(widget.NewButton("", func() {
		enableOptions = false
		resizeAndRefresh(top, bottom, left, right, center, textbox, content, dividers, root.Size(), right.Visible(), enableOptions, options, root)
	}), container.NewMax(buttonColor, returnText("exit options")))
	// Define initial positions
	initialPositions := map[fyne.CanvasObject]float32{
		toggleButton3: 40,
		toggleButton5: 80,
		toggleButton6: 120,
		toggleButton4: 160,
		toggleButton7: 220, // Starting position of toggleButton4

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

	originalPosY := initialPositions[toggleButton4]
	optionsContent := container.NewMax(optionsColor, container.NewVBox(
		slider,
		toggleButton3,
		toggleButton5,
		toggleButton6,
		toggleButton4,
		spacer, // Add spacer to start with
		toggleButton7,
	))

	slider.OnChanged = func(value float64) {
		if value > 0 {
			optionsContent.Remove(spacer)
		}

		for btn, initialY := range initialPositions {
			if btn != toggleButton7 {
				newY := initialY - float32(value)
				btn.Move(fyne.NewPos(0, newY))
			}
		}

		if value >= 100 {

			newY := originalPosY - (float32(value) - 50)
			toggleButton7.Move(fyne.NewPos(0, newY))
		} else {
			toggleButton7.Move(fyne.NewPos(0, originalPosY))
		}
	}

	options = widget.NewModalPopUp(
		optionsContent,
		w.Canvas(),
	)
	options.Hide()

	var themeName string
	if fyne.CurrentApp().Settings().Theme() == theme.LightTheme() {
		themeName = "Light"
	} else {
		themeName = "Dark"
	}

	spacer2 := layout.NewSpacer()

	themeLabel := widget.NewLabel("theme:  " + themeName)
	square := canvas.NewRectangle(themeLabelColor)
	square.Resize(fyne.NewSize(50, 50))
	themeContainer := container.NewCenter(themeLabel)       // Center the label inside the square
	themeSquare := container.NewMax(square, themeContainer) // Overlay the label on the square

	leftContent := container.NewVBox(
		widget.NewLabel("Buttons:"),
		toggleButton1,
		toggleButton2,
		spacer2,
		themeSquare,
	)

	// Layer the background and the VBox using container.NewMax
	left = container.NewMax(
		sideBarColor,
		leftContent,
	)

	// Update the root container to use the rectangle in the center
	//dividers[0], dividers[1], dividers[2]
	// root = container.NewBorder(top, bottom, left, right, centerRect, options)
	root = container.NewBorder(top, bottom, left, right, options, dividers[0], dividers[1], dividers[2], center)

	enableOptions = false

	root.Resize(fyne.NewSize(800, 600))
	resizeAndRefresh(top, bottom, left, right, center, textbox, content, dividers, root.Size(), right.Visible(), enableOptions, options, root)

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

var darkMoistYellow = color.RGBA{R: 0xB8, G: 0x8F, B: 0x00, A: 0xFF} // Adjust the RGB values as needed

func getSidebar(w fyne.Window, customTheme *CustomTheme) *fyne.Container {
	// Create a rectangle with the dark moist yellow color
	containerRect := canvas.NewRectangle(darkMoistYellow)

	// Set a specific size for the container
	containerRect.SetMinSize(fyne.NewSize(100, 200)) // Set desired size
	returnToggleButton := container.NewMax(widget.NewButton("", func() {
		returnToMenu(w, customTheme)
	}), container.NewMax(buttonColor, returnText("exit options")))
	// Create a container with the rectangle and a label
	container := container.NewMax(containerRect, returnToggleButton)

	return container
}
func returnToMenu(w fyne.Window, customTheme *CustomTheme) {
	// fmt.Println("test")
	w.SetContent(makeGUI(w, customTheme))
}
func getFiller() *fyne.Container {
	// Define the black color
	blackColor := color.RGBA{0x10, 0x10, 0x10, 0xff}

	// Create a rectangle with the black color
	containerRect := canvas.NewRectangle(blackColor)
	containerRect.SetMinSize(fyne.NewSize(500, 200))
	// Return a container with the black rectangle
	return container.NewMax(containerRect)
}

type GlobalUI struct {
	Window    fyne.Window
	Button    *widget.Button
	Label     *widget.Label
	Image     *canvas.Image
	Container *fyne.Container
	Canvas    fyne.Canvas
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
