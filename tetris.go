package main

import (
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
)

type BlockGroup struct {
	Shape  [][]int     // 2D array to represent the shape
	Color  color.NRGBA // Color of the block
	Width  int         // Width of the block
	Height int         // Height of the block
}

//	var blockGroups = []BlockGroup{
//		{
//			Shape: [][]int{
//				{1, 1, 1, 1}, // I shape
//			},
//			Color:  color.NRGBA{0, 255, 255, 255}, // Cyan
//			Width:  4,
//			Height: 1,
//		},
//		{
//			Shape: [][]int{
//				{1, 1},
//				{1, 1}, // O shape
//			},
//			Color:  color.NRGBA{255, 255, 0, 255}, // Yellow
//			Width:  2,
//			Height: 2,
//		},
//		// Add other shapes like T, S, Z, J, L...
//	}
func makeSquare(x float32) []fyne.Position {
	return []fyne.Position{
		fyne.NewPos(x, 1),
		fyne.NewPos(x, 2),
		fyne.NewPos(x+1, 1),
		fyne.NewPos(x+1, 2),
	}
}
func makeLine(x float32) []fyne.Position {
	return []fyne.Position{
		fyne.NewPos(x, 1),
		fyne.NewPos(x, 2),
	}
}

func createTetris(w fyne.Window, customTheme *CustomTheme, isBeta bool) *fyne.Container {
	const gridWidth, gridHeight = 10, 20
	const cellSize = 20

	// Define retro colors
	backgroundColor := color.RGBA{0x10, 0x10, 0x10, 0xff} // Dark background
	cellColor := color.RGBA{0x22, 0x22, 0x22, 0xff}       // Darker cell color
	textColor := color.RGBA{0xcc, 0xcc, 0xcc, 0xff}       // Retro light gray for text

	// Create a background rectangle
	background := canvas.NewRectangle(backgroundColor)
	background.SetMinSize(fyne.NewSize(gridWidth*cellSize, gridHeight*cellSize))

	// Create the grid for Tetris blocks
	cells := make([]*canvas.Rectangle, gridWidth*gridHeight)
	// var groupCells []*canvas.Rectangle
	grid := container.NewGridWithColumns(gridWidth)

	for i := range cells {
		cells[i] = canvas.NewRectangle(cellColor)
		cells[i].SetMinSize(fyne.NewSize(cellSize, cellSize))
		grid.Add(cells[i])
	}

	keys := make(map[fyne.KeyName]int)
	var keysMutex sync.Mutex
	if can, ok := w.Canvas().(desktop.Canvas); ok {
		can.SetOnKeyDown(func(ev *fyne.KeyEvent) {
			keysMutex.Lock()
			keys[ev.Name] = 0
			keysMutex.Unlock()
		})
		can.SetOnKeyUp(func(ev *fyne.KeyEvent) {
			keysMutex.Lock()
			delete(keys, ev.Name)
			keysMutex.Unlock()
		})
	}

	isKeyPressed := func(key fyne.KeyName) bool {
		keysMutex.Lock()
		defer keysMutex.Unlock()
		return keys[key] == 1 || keys[key] > 15
	}

	processInput := func() {
		keysMutex.Lock()
		defer keysMutex.Unlock()
		for key := range keys {
			keys[key]++
		}
	}

	bufferedGrid := make([]color.Color, gridWidth*gridHeight)
	renderedGrid := make([]color.Color, gridWidth*gridHeight)

	pieceX, pieceY := 4, 0
	lockedCells := make([]*color.NRGBA, gridWidth*gridHeight)

	canFall := func() bool {
		if pieceY == gridHeight-1 {
			return false
		}
		index := (pieceY+1)*gridWidth + pieceX
		return lockedCells[index] == nil
	}
	// actions := map[string]func(x float32) []fyne.Position{
	// 	"square": makeSquare,
	// 	"line":   makeLine,
	// 	// "corner": makeCorner, // Uncomment if makeCorner is available
	// }
	// actionNames := make([]string, 0, len(actions))
	// for name := range actions {

	// 	actionNames = append(actionNames, name)
	// }
	// randomIndex := rand.Intn(max(1, len(actionNames)))
	// selectedActionName := actionNames[randomIndex]

	renderPiece := func() {
		// for _, cell := range actions[selectedActionName](float32(pieceX)) {
		// Calculate the index based on the x and y positions of the piece
		// index := int(cell.Y)*gridWidth + int(cell.X)
		index := pieceY*gridWidth + pieceX

		// Update the buffered grid with the desired color
		bufferedGrid[index] = color.NRGBA{255, 0, 0, 255}
		//}
	}

	deleteRow := func(y int) {
		for x := 0; x < gridWidth; x++ {
			lockedCells[y*gridWidth+x] = nil
		}
		for y2 := y; y2 > 0; y2-- {
			for x := 0; x < gridWidth; x++ {
				lockedCells[y2*gridWidth+x] = lockedCells[(y2-1)*gridWidth+x]
			}
		}
	}

	clearRows := func() {
		for y := 0; y < gridHeight; y++ {
			full := true
			for x := 0; x < gridWidth; x++ {
				if lockedCells[y*gridWidth+x] == nil {
					full = false
					break
				}
			}
			if full {
				deleteRow(y)
			}
		}
	}

	renderLockedCells := func() {
		for i, c := range lockedCells {
			if c != nil {
				bufferedGrid[i] = *c
			}
		}
	}

	clearBuffer := func() {
		for i := range bufferedGrid {
			bufferedGrid[i] = color.NRGBA{127, 127, 127, 255}
		}
	}

	blit := func() {
		for i := range renderedGrid {
			if bufferedGrid[i] != renderedGrid[i] {
				renderedGrid[i] = bufferedGrid[i]
				cells[i].FillColor = bufferedGrid[i]
				cells[i].Refresh()
			}
		}
	}

	var fallTick int

	go loop(10, func(tick int) {
		processInput()

		if isKeyPressed(fyne.KeyLeft) && pieceX > 0 {
			pieceX--
		}
		if isKeyPressed(fyne.KeyRight) && pieceX < gridWidth-1 {
			pieceX++
		}
		if isKeyPressed(fyne.KeySpace) {
			for canFall() {
				pieceY++
			}
			fallTick = tick
		}

		if tick-fallTick >= 5 {
			if canFall() {
				pieceY++
				fallTick = tick
			} else {
				index := pieceY*gridWidth + pieceX
				lockedCells[index] = &color.NRGBA{255, 0, 0, 255}
				pieceX, pieceY = 4, 0
				clearRows()
			}
		}

		renderLockedCells()
		renderPiece()
		blit()
		clearBuffer()
	})

	// Create labels for score and level
	scoreLabel := canvas.NewText("Score: 0", textColor)
	scoreLabel.TextStyle = fyne.TextStyle{Bold: true}
	scoreLabel.Alignment = fyne.TextAlignCenter

	levelLabel := canvas.NewText("Level: 1", textColor)
	levelLabel.TextStyle = fyne.TextStyle{Bold: true}
	levelLabel.Alignment = fyne.TextAlignCenter

	header := container.NewHBox(
		layout.NewSpacer(),
		scoreLabel,
		layout.NewSpacer(),
		levelLabel,
		layout.NewSpacer(),
	)

	// Create a footer with control instructions
	footerLabel := canvas.NewText("Controls: Arrow keys to move, Space to drop", textColor)
	footerLabel.TextStyle = fyne.TextStyle{Bold: true}
	footerLabel.Alignment = fyne.TextAlignCenter

	footer := container.NewHBox(
		layout.NewSpacer(),
		footerLabel,
		layout.NewSpacer(),
	)

	// Create a centered container with padding
	gridWrapper := container.NewVBox(
		layout.NewSpacer(),
		container.NewHBox(
			layout.NewSpacer(),
			grid,
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
	)

	// Combine all elements into the main game container
	gameContainer := container.NewVBox(
		header,
		gridWrapper,
		footer,
	)
	// returnToggleButton := container.NewMax(widget.NewButton("", func() {
	// 	returnToMenu(w, customTheme)
	// }), container.NewMax(buttonColor, returnText("exit options")))
	// Overlay the game UI on top of the background
	content := container.NewMax(background, gameContainer)
	content = container.NewBorder(nil, nil, nil, getSidebar(w, customTheme), content)

	return content
}

// func isKeyPressed(key fyne.KeyName) bool {
// 	keysMutex.Lock()
// 	defer keysMutex.Unlock()
// 	return keys[key] == 0 || keys[key] > 0
// }

// loop runs a fixed timestep game loop. It calls fn a fixed number of times per second.
func loop(tps int, fn func(tick int)) {
	// if isKeyPressed(fyne.KeySpace) {
	// 	// Drop the piece all the way down
	// 	for canFall() {
	// 		pieceY++
	// 	}
	// 	fallTick = tick // Reset the fall tick to prevent immediate locking
	// }
	lastTick := time.Now().UnixNano()
	var tick int
	for {
		now := time.Now().UnixNano()
		for now-lastTick >= 1e9/int64(tps) {
			fn(tick)
			lastTick += 1e9 / int64(tps)
			tick++
		}
	}
}
