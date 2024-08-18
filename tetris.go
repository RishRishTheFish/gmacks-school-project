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

func createTetris(w fyne.Window, isBeta bool) *fyne.Container {
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

	renderPiece := func() {
		index := pieceY*gridWidth + pieceX
		bufferedGrid[index] = color.NRGBA{255, 0, 0, 255}
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

	// Overlay the game UI on top of the background
	content := container.NewMax(background, gameContainer)

	return content
}

// loop runs a fixed timestep game loop. It calls fn a fixed number of times per second.
func loop(tps int, fn func(tick int)) {
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
