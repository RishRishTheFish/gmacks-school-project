package main

import (
	"image/color"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

const (
	gridWidth   = 10
	gridHeight  = 20
	bottomLimit = 13
	cellSize    = 20
)

var randSource = rand.New(rand.NewSource(time.Now().UnixNano()))

// randomColor generates a random color with improved randomness
func randomColor() color.Color {
	return color.RGBA{
		R: uint8(randSource.Intn(256)),
		G: uint8(randSource.Intn(256)),
		B: uint8(randSource.Intn(256)),
		A: uint8(randSource.Intn(256)), // Random alpha value for more variation
	}
}

var grayColor = color.Gray{Y: 0x30}

// Convert to color.RGBA
var rgbaGrayColor = color.RGBA{
	R: grayColor.Y,
	G: grayColor.Y,
	B: grayColor.Y,
	A: 255, // Fully opaque
}
var rgbaRedColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}

type ExtraParams struct {
	Length int
	Type   string
	//	stopChan  chan bool
	//	closeOnce sync.Once
	// Add other fields if needed
}

var (
	stopChan  = make(chan bool)
	closeOnce sync.Once
)

type cellsParams func(grid fyne.Container, cells [][]*canvas.Rectangle, randNum int, color color.Color, doBeta bool, extraParams ExtraParams)

func containsPos(slice []fyne.Position, value fyne.Position) bool {
	if len(slice) > 1 {
		for _, v := range slice {
			if v == value {
				return true
			}
		}
	}
	return false
}

var previousPositions []fyne.Position
var limit int
var previousPositionsHasBottom bool

// var currentPos fyne.Position
var currentGroup []fyne.Position

var globalLimit int

// var allignment []int
var allignmentPos []fyne.Position

// // var groupCells []*fyne.Position
// func minArr(arr []int) int {
// 	if len(arr) == 0 {
// 		panic("Cannot find minimum of an empty slice")
// 	}
// 	minVal := arr[0]
// 	for _, val := range arr {
// 		if val < minVal {
// 			minVal = val
// 		}
// 	}
// 	return minVal
// }

// // Function to find the maximum value in a slice of integers
// func maxArr(arr []int) int {
// 	if len(arr) == 0 {
// 		panic("Cannot find maximum of an empty slice")
// 	}
// 	maxVal := arr[0]
// 	for _, val := range arr {
// 		if val > maxVal {
// 			maxVal = val
// 		}
// 	}
// 	return maxVal
// }

// var tempGroupCells []*fyne.Position
func removePos(positions []fyne.Position, pos fyne.Position) []fyne.Position {
	for i, p := range positions {
		if p == pos {
			return append(positions[:i], positions[i+1:]...)
		}
	}
	return positions
}

// func findMinValue(m map[int]int) (int, error) {
// 	if len(m) == 0 {
// 		return 0, fmt.Errorf("map is empty")
// 	}

// 	maxValue := math.MaxInt // Start with the maximum possible integer value

// 	for _, value := range m {
// 		if value < maxValue {
// 			maxValue = value
// 		}
// 	}

type ConcurrentMap struct {
	mu   sync.Mutex
	data map[int]int
}

// NewConcurrentMap creates a new ConcurrentMap
func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{
		data: make(map[int]int),
	}
}

// Set updates the map with a new value for the given key
func (cm *ConcurrentMap) Set(key, value int) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.data[key] = value
}
func colorToRGBA(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(a >> 8),
	}
}

func fall(
	grid fyne.Container,
	cells [][]*canvas.Rectangle,
	groupCells []fyne.Position,
	clr color.Color,
	params ExtraParams,
	bottomLimit int,
	//keyEventChannel KeyEvent
) {
	rgbaGrayColor := color.RGBA{128, 128, 128, 255} // Example gray color

	for {
		canMoveDown := true

		// Check if the group of cells can move down
		for _, pos := range groupCells {
			x, y := int(pos.X), int(pos.Y)
			if y+1 >= len(cells) || colorToRGBA(cells[y+1][x].FillColor) != colorToRGBA(rgbaGrayColor) {
				canMoveDown = false
				break
			}
		}

		if !canMoveDown {
			break
		}

		// Move the group of cells down
		newGroupCells := []fyne.Position{}
		for _, pos := range groupCells {
			x, y := int(pos.X), int(pos.Y)
			cells[y][x].FillColor = rgbaGrayColor
			cells[y][x].Refresh()

			newY := y + 1
			cells[newY][x].FillColor = clr
			cells[newY][x].Refresh()

			newGroupCells = append(newGroupCells, fyne.NewPos(float32(x), float32(newY)))
		}

		groupCells = newGroupCells

		// Delay for visual effect
		time.Sleep(200 * time.Millisecond)
	}

	// Once the cells can no longer move down, fix them in place
	for _, pos := range groupCells {
		x, y := int(pos.X), int(pos.Y)
		cells[y][x].FillColor = clr
		cells[y][x].Refresh()
	}

	// Check if any rows are filled and clear them
	clearFullRows(grid, cells, params)
}

func clearFullRows(
	grid fyne.Container,
	cells [][]*canvas.Rectangle,
	params ExtraParams,
) {
	for y := len(cells) - 1; y >= 0; y-- {
		isRowFull := true
		for x := 0; x < len(cells[y]); x++ {
			if cells[y][x].FillColor == rgbaGrayColor {
				isRowFull = false
				break
			}
		}

		if isRowFull {
			// Clear the row
			for x := 0; x < len(cells[y]); x++ {
				cells[y][x].FillColor = rgbaGrayColor
				cells[y][x].Refresh()
			}

			// Move all rows above down by one
			for yy := y; yy > 0; yy-- {
				for xx := 0; xx < len(cells[yy]); xx++ {
					cells[yy][xx].FillColor = cells[yy-1][xx].FillColor
					cells[yy][xx].Refresh()
				}
			}

			// Reset the top row to empty
			for xx := 0; xx < len(cells[0]); xx++ {
				cells[0][xx].FillColor = rgbaGrayColor
				cells[0][xx].Refresh()
			}

			// Check the same row again after clearing
			y++
		}
	}
}

func makeCorner(
	grid fyne.Container,
	cells [][]*canvas.Rectangle,
	randNum int,
	color color.Color,
	doBeta bool,
	params ExtraParams,
	//keyEventChannel KeyEvent
) {
	cells[0][randNum].FillColor = color
	cells[0][randNum].Refresh()
	pos1 := fyne.NewPos(float32(randNum), 0)
	groupCells := []fyne.Position{pos1}
	cells[1][randNum].FillColor = color
	cells[1][randNum].Refresh()
	pos2 := fyne.NewPos(float32(randNum), 1)
	groupCells = append(groupCells, pos2)
	time.Sleep(1 * time.Second)
	if doBeta {
		fallBeta(
			grid,
			cells,
			groupCells,
			color,
			true,
			params,
			bottomLimit,
			//keyEventChannel
		)
	} else {
		fall(grid, cells, groupCells, color, params, bottomLimit)
	}
}

func makeLine(
	grid fyne.Container,
	cells [][]*canvas.Rectangle,
	randNum int,
	color color.Color,
	doBeta bool,
	params ExtraParams,
	//keyEventChannel KeyEvent
) {
	cells[0][randNum].FillColor = color
	cells[0][randNum].Refresh()
	pos1 := fyne.NewPos(float32(randNum), 0)
	groupCells := []fyne.Position{pos1}
	cells[1][randNum].FillColor = color
	cells[1][randNum].Refresh()
	pos2 := fyne.NewPos(float32(randNum), 1)
	groupCells = append(groupCells, pos2)
	time.Sleep(1 * time.Second)
	if doBeta {
		fallBeta(
			grid,
			cells,
			groupCells,
			color,
			true,
			params,
			bottomLimit,
			//keyEventChannel
		)
	} else {
		fall(grid, cells, groupCells, color, params, bottomLimit)
	}
}

func makeSquare(
	grid fyne.Container,
	cells [][]*canvas.Rectangle,
	randNum int,
	color color.Color,
	doBeta bool,
	params ExtraParams,
	//keyEventChannel KeyEvent
) {
	cells[0][randNum].FillColor = color
	cells[0][randNum].Refresh()
	pos1 := fyne.NewPos(float32(randNum), 0)
	groupCells := []fyne.Position{pos1}
	cells[1][randNum].FillColor = color
	cells[1][randNum].Refresh()
	pos2 := fyne.NewPos(float32(randNum), 1)
	groupCells = append(groupCells, pos2)
	if randNum+1 != 10 {
		cells[1][randNum+1].FillColor = color
		cells[1][randNum+1].Refresh()
		pos3 := fyne.NewPos(float32(randNum+1), 1)
		groupCells = append(groupCells, pos3)
		cells[0][randNum+1].FillColor = color
		cells[0][randNum+1].Refresh()
		pos4 := fyne.NewPos(float32(randNum+1), 0)
		groupCells = append(groupCells, pos4)
	}
	time.Sleep(1 * time.Second)
	if doBeta {
		fallBeta(
			grid,
			cells,
			groupCells,
			color,
			true,
			params,
			bottomLimit,
			//keyEventChannel
		)
	}
}
func removeAction(actions []cellsParams, index int) []cellsParams {
	if index < 0 || index >= len(actions) {
		return actions
	}
	return append(actions[:index], actions[index+1:]...)
}
func moveBlocksDown(cells [][]*canvas.Rectangle, clearedRow int) {
	// Iterate from the row above the cleared row up to the top of the grid
	for y := clearedRow - 1; y >= 0; y-- {
		for x := 0; x < gridWidth; x++ {
			// Check if there is a block in the current cell
			if cells[yoffset+y][x].FillColor != rgbaGrayColor { // Assuming rgbaGrayColor is the color for empty cells
				// Move block to the row below
				newY := y + 1
				if newY < len(cells) {
					cells[yoffset+newY][x].FillColor = cells[yoffset+y][x].FillColor
					cells[yoffset+newY][x].Refresh()
					cells[yoffset+y][x].FillColor = rgbaGrayColor
					cells[yoffset+y][x].Refresh()
				}
			}
		}
	}
}

func handleKeyEvents(isNormal bool) int {
	xoffset := 0

	if isNormal && latestKeyEvent != previousKeyEvent {
		switch latestKeyEvent.KeyName {
		case fyne.KeyLeft:
			xoffset = -1
		case fyne.KeyRight:
			xoffset = 1
		}
		previousKeyEvent = latestKeyEvent
	}

	return xoffset
}

//}

func ensureMapInitialized(m *map[int]int) {
	if *m == nil {
		*m = make(map[int]int)
	}
}

// var positions sync.Map
func MaxYPosition(positions map[int]int, length int) (int, int) {
	minYForX.mu.Lock()
	defer minYForX.mu.Unlock()

	if len(positions) == 0 || length <= 0 {
		return 0, 0 // Return a default value if the map is empty or length is invalid
	}

	// Create a slice to hold the positions
	type position struct {
		x, y int
	}
	var posSlice []position

	for x, y := range positions {
		posSlice = append(posSlice, position{x, y})
	}

	// Sort the slice by the y value in descending order
	sort.Slice(posSlice, func(i, j int) bool {
		return posSlice[i].y > posSlice[j].y
	})

	// Ensure length does not exceed the number of positions
	if length > len(posSlice) {
		length = len(posSlice)
	}

	// Select the desired maximum position based on the length
	selectedPos := posSlice[length-1]

	return selectedPos.x, selectedPos.y
}
func applyRandomColors(
	grid *fyne.Container,
	cells [][]*canvas.Rectangle,
	doBeta bool,
	//keyEventChannel KeyEvent
) {
	for {
		randNum := rand.Intn(gridWidth)
		colour := randomColor()

		// Define actions and corresponding types
		actions := map[string]func(grid fyne.Container, cells [][]*canvas.Rectangle, randNum int, colour color.Color, doBeta bool, params ExtraParams){
			"square": makeSquare,
			"line":   makeLine,
			// "corner": makeCorner, // Uncomment if makeCorner is available
		}

		// Pick a random key (action name)
		actionNames := make([]string, 0, len(actions))
		for name := range actions {

			actionNames = append(actionNames, name)
		}
		randomIndex := rand.Intn(max(1, len(actionNames)))
		selectedActionName := actionNames[randomIndex]

		var length int
		switch selectedActionName {
		case "square":
			length = 2
		case "line":
			length = 1
		default:
			length = 1 // Default length if not specified
		}

		// Prepare ExtraParams with the appropriate type
		params := ExtraParams{
			Length: length,
			Type:   selectedActionName, // Set the type based on the selected action
			// Initialize other fields if needed
		}

		// Execute the selected function
		totalGenerations++
		actions[selectedActionName](
			*grid,
			cells,
			randNum,
			colour,
			doBeta,
			params,
			//keyEventChannel
		)
	}
}

type KeyEvent struct {
	KeyName   fyne.KeyName
	Increment int
}

func listenKeyEvent(
	w fyne.Window,
	//keyEventChannel KeyEvent
) {
	w.Canvas().SetOnTypedKey(func(keyEvent *fyne.KeyEvent) {
		latestKeyEvent = KeyEvent{
			KeyName:   keyEvent.Name,
			Increment: latestKeyEvent.Increment + 1,
		}
	})

}

var (
	score      int
	scoreLabel *canvas.Text
)

func createTetris(w fyne.Window, doBeta bool) *fyne.Container {
	// var a sync.Mutex
	// a.Lock()
	// a.Lock()
	// fmt.Println("hello world")
	const cellSize = 20 // Set the desired size for both width and height

	// Define retro colors
	backgroundColor := color.RGBA{0x10, 0x10, 0x10, 0xff} // Dark background
	cellColor := color.RGBA{0x22, 0x22, 0x22, 0xff}       // Darker cell color
	textColor := color.RGBA{0xcc, 0xcc, 0xcc, 0xff}       // Retro light gray for text

	// Create a dark background rectangle
	background := canvas.NewRectangle(backgroundColor)
	background.SetMinSize(fyne.NewSize(gridWidth*cellSize, gridHeight*cellSize))

	// Create a container for the Tetris game
	gameContainer := container.NewVBox()

	// Initialize the score label
	scoreLabel = canvas.NewText("Score: "+strconv.Itoa(score), textColor)
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

	// Create the grid for the Tetris blocks
	cells := make([][]*canvas.Rectangle, gridHeight)
	for y := 0; y < gridHeight; y++ {
		cells[yoffset+y] = make([]*canvas.Rectangle, gridWidth) // Initialize the inner slice
	}

	grid := container.NewGridWithColumns(gridWidth)

	for y := 0; y < gridHeight; y++ {
		for x := 0; x < gridWidth; x++ {
			bg := canvas.NewRectangle(cellColor)
			bg.SetMinSize(fyne.NewSize(cellSize, cellSize)) // Use the same size for both width and height
			cells[yoffset+y][x] = bg
			grid.Add(bg)
		}
	}

	// Create a footer with control instructions using canvas.Text
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
		layout.NewSpacer(), // Add vertical space at the top
		container.NewHBox(
			layout.NewSpacer(), // Add horizontal space to the left
			grid,               // The grid in the center
			layout.NewSpacer(), // Add horizontal space to the right
		),
		layout.NewSpacer(), // Add vertical space at the bottom
	)

	// Combine all elements into the main game container
	gameContainer.Add(header)
	gameContainer.Add(gridWrapper)
	gameContainer.Add(footer)

	// Overlay the game UI on top of the background
	content := container.NewMax(background, gameContainer)

	// Run the game logic in separate goroutines
	go applyRandomColors(grid, cells, doBeta)
	go listenKeyEvent(w)

	return content
}

// UpdateScore updates the score label with the current score
func UpdateScore() {
	scoreLabel.Text = "Score: " + strconv.Itoa(score)
	canvas.Refresh(scoreLabel)
}

// Example function to increase the score
func increaseScore(points int) {
	score += points
	UpdateScore() // Update the score label after changing the score
}
