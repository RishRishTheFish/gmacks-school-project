package main

import (
	"fmt"
	"image/color"
	"math"
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

type cellsParams func(grid fyne.Container, cells [][]*canvas.Rectangle, randNum int, color color.Color, extraParams ExtraParams)

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

// var groupCells []*fyne.Position
func minArr(arr []int) int {
	if len(arr) == 0 {
		panic("Cannot find minimum of an empty slice")
	}
	minVal := arr[0]
	for _, val := range arr {
		if val < minVal {
			minVal = val
		}
	}
	return minVal
}

// Function to find the maximum value in a slice of integers
func maxArr(arr []int) int {
	if len(arr) == 0 {
		panic("Cannot find maximum of an empty slice")
	}
	maxVal := arr[0]
	for _, val := range arr {
		if val > maxVal {
			maxVal = val
		}
	}
	return maxVal
}

// var tempGroupCells []*fyne.Position
func removePos(positions []fyne.Position, pos fyne.Position) []fyne.Position {
	for i, p := range positions {
		if p == pos {
			return append(positions[:i], positions[i+1:]...)
		}
	}
	return positions
}

func findMinValue(m map[int]int) (int, error) {
	if len(m) == 0 {
		return 0, fmt.Errorf("map is empty")
	}

	maxValue := math.MaxInt // Start with the maximum possible integer value

	for _, value := range m {
		if value < maxValue {
			maxValue = value
		}
	}

	return maxValue, nil
}
func findMaxValue(m map[int]int) (int, error) {
	if len(m) == 0 {
		return 0, fmt.Errorf("map is empty")
	}

	// Start with the minimum possible integer value
	maxValue := math.MinInt

	for _, value := range m {
		if value > maxValue {
			maxValue = value
			// key.FillColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
			// key.Refresh
		}
	}

	return maxValue, nil
}

// ConcurrentMap is a thread-safe wrapper around a map[int]int
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

// var maxYForX = NewConcurrentMap()
var minYForX = NewConcurrentMap()
var currentY int
var currentMin int

// var //maxMap map[int]int

var rowCounters map[int]int
var totalGenerations int // Track the total number of generations
var dontColor bool

func init() {
	rowCounters = make(map[int]int)
}

var latestKeyEvent KeyEvent
var previousKeyEvent KeyEvent
var rowCountersMutex sync.Mutex

var yoffset int

// var clearLowestPoint int

func fall(
	grid fyne.Container,
	cells [][]*canvas.Rectangle,
	groupCells []fyne.Position,
	allCells []fyne.Position,
	color color.Color,
	isNormal bool,
	params ExtraParams,
	limit int,
) {
	defer func() {
		closeOnce.Do(func() {
			close(stopChan)
		})
	}()

	var doesContainPos bool
	if isNormal {
		// Initialize or reset limit
	}

	toBeCleared := make(map[fyne.Position]bool)

	// Initialize previous positions if not already done
	if isNormal && !previousPositionsHasBottom {
		previousPositionsHasBottom = true
		for i := 0; i < gridWidth; i++ {
			previousPositions = append(previousPositions, fyne.NewPos(float32(i), float32(limit)+1))
		}
	}

	for i := 0; i < limit; i++ {
		select {
		case <-stopChan:
			return // Exit the function if stop is triggered
		default:
			// Continue falling
		}

		if i == 3 && isNormal {
			for j := 0; j < gridWidth; j++ {
				var tempCellArr []fyne.Position
				tempCellArr = append(tempCellArr, fyne.NewPos(float32(j), 1))
				go fall(grid, cells, tempCellArr, tempCellArr, randomColor(), false, params, limit)
			}
		}

		xoffset := handleKeyEvents(isNormal)
		newGroupCells := []fyne.Position{}

		for _, pos := range groupCells {
			x, y := int(pos.X), int(pos.Y)
			if y < len(cells) && x < len(cells[y]) {
				toBeCleared[pos] = true
			}
		}

		for pos := range toBeCleared {
			x, y := int(pos.X), int(pos.Y)
			if y < len(cells) && x < len(cells[y]) {
				cells[y][x].FillColor = rgbaGrayColor
				cells[y][x].Refresh()
			}
		}

		for _, pos := range groupCells {
			x, y := int(pos.X), int(pos.Y)
			newX := x + xoffset
			newY := y + 1

			if newY < len(cells) && newX >= 0 && newX < len(cells[newY]) {
				newPos := fyne.NewPos(float32(newX), float32(newY))
				oldPos := fyne.NewPos(float32(x), float32(y))

				if !isNormal {
					allignmentPos = removePos(allignmentPos, oldPos)
					allignmentPos = append(allignmentPos, newPos)
				}

				if containsPos(previousPositions, newPos) {
					doesContainPos = true
					if isNormal {
						allignmentPos = []fyne.Position{}
					}
					limit = y

					if !isNormal {
						minYForX.Set(x, newY)
					}
				}

				if isNormal {
					if doesContainPos {
						// Notify other segments to stop
						closeOnce.Do(func() {
							close(stopChan)
						})

						for _, cell := range groupCells {
							if len(cells) >= int(cell.Y) && len(cells[int(cell.Y)]) >= int(cell.X) {
								cells[int(cell.Y)][int(cell.X)].FillColor = color
								cells[int(cell.Y)][int(cell.X)].Refresh()
							}
							previousPositions = append(previousPositions, cell)
						}

						_, maxY := MaxYPosition(minYForX.data, 1)

						for i := params.Length; i > 0; i-- {
							_, maxY = MaxYPosition(minYForX.data, i)
							if gridWidth-rowCounters[maxY] >= i {
								break
							}
						}

						if currentMin != maxY {
							for _, cell := range allignmentPos {
								if cell.Y < float32(maxY) {
									rowCounters[maxY]++
								}
							}
						}

						if newY < maxY {
							rowCounters[maxY]--
						}
						currentMin = maxY
						rowCounters[maxY]++

						rowCountersMutex.Lock()
						defer rowCountersMutex.Unlock()

						if rowCounters[maxY] >= gridWidth-1 {
							increaseScore(1)
							go func() {
								time.Sleep(1 * time.Millisecond)
								for x := 0; x < gridWidth; x++ {
									previousPositions = removePos(previousPositions, fyne.NewPos(float32(x), float32(maxY)-2))
									cells[maxY][x].FillColor = rgbaGrayColor
									cells[maxY][x].Refresh()
									if maxY-1 >= 0 {
										previousPositions = removePos(previousPositions, fyne.NewPos(float32(x), float32(maxY)-1))
										cells[maxY-1][x].FillColor = rgbaGrayColor
										cells[maxY-1][x].Refresh()
									}
								}
								time.Sleep(500 * time.Millisecond) // Delay to prevent early block spawning
							}()
							delete(rowCounters, maxY)
						}
					}
					cells[newY][newX].FillColor = color
					cells[newY][newX].Refresh()
				} else {
					// Additional logic for non-normal cells
				}

				newGroupCells = append(newGroupCells, newPos)
			}
		}

		groupCells = newGroupCells
		if isNormal {
			currentGroup = groupCells
		}
		toBeCleared = make(map[fyne.Position]bool)
		if latestKeyEvent.KeyName != fyne.KeySpace {
			time.Sleep(500 * time.Millisecond)
		}
	}

	if isNormal {
		latestKeyEvent = KeyEvent{}
	}
}
func splitAndFall(
	grid fyne.Container,
	cells [][]*canvas.Rectangle,
	groupCells []fyne.Position,
	color color.Color,
	params ExtraParams,
	limit int,
) {
	// Create a map to group cells by their X position
	cellGroups := make(map[int][]fyne.Position)
	for _, pos := range groupCells {
		x := int(pos.X)
		cellGroups[x] = append(cellGroups[x], pos)
	}

	// Run goroutines for each group of cells with the same X position
	for _, group := range cellGroups {
		go fall(grid, cells, group, groupCells, color, true, params, limit)
	}
}

func makeLine(
	grid fyne.Container,
	cells [][]*canvas.Rectangle,
	randNum int,
	color color.Color,
	params ExtraParams,
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
	splitAndFall(grid, cells, groupCells, color, params, bottomLimit)
}

func makeSquare(
	grid fyne.Container,
	cells [][]*canvas.Rectangle,
	randNum int,
	color color.Color,
	params ExtraParams,
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
	splitAndFall(grid, cells, groupCells, color, params, bottomLimit)
}
func makeCorner(
	grid fyne.Container,
	cells [][]*canvas.Rectangle,
	randNum int,
	color color.Color,
	params ExtraParams,
	//keyEventChannel KeyEvent
) {
	cells[yoffset+0][randNum].FillColor = color
	cells[yoffset+0][randNum].Refresh()
	pos1 := fyne.NewPos(float32(randNum), 0)
	groupCells := []fyne.Position{pos1}
	cells[yoffset+1][randNum].FillColor = color
	cells[yoffset+1][randNum].Refresh()
	pos2 := fyne.NewPos(float32(randNum), 1)
	groupCells = append(groupCells, pos2)
	time.Sleep(1 * time.Second)
	fall(
		grid,
		cells,
		groupCells,
		groupCells,
		color,
		true,
		params,
		bottomLimit,
		//keyEventChannel
	)
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

// func MaxYPosition(positions map[int]int) (int, int) {
// 	// Lock the mutex for reading
// 	// positionsRWMutex.RLock()
// 	// Ensure that the mutex is unlocked when the function returns
// 	// positions.cm.mu.Lock()
// 	minYForX.mu.Lock()
// 	defer minYForX.mu.Unlock()
// 	// copy()
// 	// positionsRWMutex.RLock()
// 	if len(positions) == 0 {
// 		// minYForX.mu.Unlock()
// 		return 0, 0 // Return a default value if the map is empty
// 	}

// 	var maxX int
// 	maxY := positions[maxX]
// 	// defer positionsRWMutex.RUnlock()
// 	for x, y := range positions {
// 		if y > maxY {
// 			maxY = y
// 			maxX = x
// 		}
// 	}
// 	// minYForX.mu.Unlock()

// 	return maxX, maxY
// }

// Example function that modifies the map
//
//	func updatePositions(positions map[int]int, x int, y int) {
//		// Lock the mutex before modifying the map
//		positionsMutex.Lock()
//		positions[x] = y
//		positionsMutex.Unlock()
//	}
// var positionsRWMutex sync.RWMutex

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
	//keyEventChannel KeyEvent
) {
	for {
		randNum := rand.Intn(gridWidth)
		colour := randomColor()

		// Define actions and corresponding types
		actions := map[string]func(grid fyne.Container, cells [][]*canvas.Rectangle, randNum int, colour color.Color, params ExtraParams){
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

// func createGrid() *fyne.Container {
// 	grid := container.NewGridWithColumns(8)

// 	for y := 0; y < 8; y++ {
// 		for x := 0; x < 8; x++ {
// 			bg := canvas.NewRectangle(color.Gray{0x30})
// 			if x%2 == y%2 {
// 				bg.FillColor = color.Gray{0xE0}
// 			}

//				img := canvas.NewImageFromResource(resourceForPiece())
//				img.FillMode = canvas.ImageFillContain
//				grid.Add(container.NewMax(bg, img))
//			}
//		}
//		return grid
//	}
var (
	score      int
	scoreLabel *canvas.Text
)

func createTetris(w fyne.Window) *fyne.Container {
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
	go applyRandomColors(grid, cells)
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

// func fall(grid fyne.Container, cells [][]*canvas.Rectangle, groupCells []fyne.Position, color color.Color, isNormal bool, params ExtraParams, limit int) {
// 	var doesContainPos bool

// 	if isNormal {
// 		globalLimit = 0
// 	}

// 	toBeCleared := make(map[fyne.Position]bool)

// 	// Initialize previous positions if not already done
// 	if isNormal && !previousPositionsHasBottom {
// 		previousPositionsHasBottom = true
// 		for i := 0; i < gridWidth; i++ {
// 			previousPositions = append(previousPositions, fyne.NewPos(float32(i), float32(bottomLimit)))
// 		}
// 	}

// 	for i := 0; i < limit; i++ {
// 		if i == 3 && isNormal {
// 			for j := 0; j < gridWidth; j++ {
// 				var tempCellArr []fyne.Position
// 				tempCellArr = append(tempCellArr, fyne.NewPos(float32(j), 1))
// 				// if params.Type == "line" {
// 				// 	tempCellArr = append(tempCellArr, fyne.NewPos(float32(j), 1))
// 				// } else if params.Type == "square" {
// 				// 	tempCellArr = append(tempCellArr, fyne.NewPos(float32(j), 1))
// 				// 	// Append positions for a 2x2 square
// 				// 	// if j < gridWidth-1 { // Ensure there's enough space horizontally
// 				// 	// 	tempCellArr = append(tempCellArr, fyne.NewPos(float32(j), 1))
// 				// 	// 	tempCellArr = append(tempCellArr, fyne.NewPos(float32(j+1), 1))
// 				// 	// 	tempCellArr = append(tempCellArr, fyne.NewPos(float32(j), 2))
// 				// 	// 	tempCellArr = append(tempCellArr, fyne.NewPos(float32(j+1), 2))
// 				// 	// }
// 				// }
// 				//fmt.Println("falling")
// 				go fall(grid, cells, tempCellArr, randomColor(), false, params, limit)
// 			}
// 		}

// 		// Update positions based on key events
// 		//if isNormal {
// 		xoffset := handleKeyEvents(isNormal)
// 		newGroupCells := []fyne.Position{}

// 		// Clear previous positions
// 		// for _, pos := range groupCells {
// 		// 	x, y := int(pos.X), int(pos.Y)
// 		// 	//fmt.Println(y)
// 		// 	if y < len(cells) && x < len(cells[max(yoffset, 0)+y]) {
// 		// 		toBeCleared[pos] = true
// 		// 	}
// 		// }
// 		for _, pos := range groupCells {
// 			x, y := int(pos.X), int(pos.Y)
// 			//fmt.Println(y)
// 			if y < len(cells) && x < len(cells[max(yoffset+y, y)]) {
// 				toBeCleared[pos] = true
// 			}
// 		}
// 		// for pos := range toBeCleared {
// 		// 	x, y := int(pos.X), int(pos.Y)
// 		// 	if y < len(cells) && x < len(cells[max(yoffset, 0)+y]) {
// 		// 		cells[max(yoffset, 0)+y][x].FillColor = rgbaGrayColor
// 		// 		cells[max(yoffset, 0)+y][x].Refresh()
// 		// 	}
// 		// }
// 		for pos := range toBeCleared {
// 			x, y := int(pos.X), int(pos.Y)
// 			if y < len(cells) && x < len(cells[max(yoffset+y, y)]) {
// 				cells[max(yoffset+y, y)][x].FillColor = rgbaGrayColor
// 				cells[max(yoffset+y, y)][x].Refresh()
// 			}
// 		}

// 		for _, pos := range groupCells {
// 			x, y := int(pos.X), int(pos.Y)
// 			newX := x + xoffset
// 			newY := y + 1

// 			if newY < len(cells) && newX >= 0 && newX < len(cells[yoffset+newY]) {
// 				newPos := fyne.NewPos(float32(newX), float32(newY))
// 				oldPos := fyne.NewPos(float32(x), float32(y))

// 				if !isNormal {
// 					allignmentPos = removePos(allignmentPos, oldPos)
// 					allignmentPos = append(allignmentPos, newPos)
// 				}

// 				if containsPos(previousPositions, newPos) {
// 					doesContainPos = true
// 					if isNormal {
// 						allignmentPos = []fyne.Position{}
// 					}
// 					limit = y

// 					if !isNormal {
// 						minYForX.Set(x, newY)
// 						// for i := newY - 1; i < limit; i++ {
// 						// 	previousPositions = removePos(previousPositions, fyne.NewPos(float32(x), float32(i)))
// 						// }
// 						// fmt.Println("prev pos", previousPositions)
// 					}
// 					// for i := newY; i < limit; i++ { // start of the execution block
// 					// 	previousPositions = removePos(previousPositions, fyne.NewPos(float32(x), float32(i)))
// 					// 	//fmt.Println("Hello") // prints "Hello" 3 times
// 					// }
// 				}

// 				if isNormal {
// 					if doesContainPos {
// 						for _, cell := range groupCells {
// 							//cell.FillColor = color
// 							//cell.Refresh()
// 							if len(cells) >= int(cell.Y) && len(cells[yoffset+int(cell.Y)]) >= int(cell.X) {
// 								cells[yoffset+int(cell.Y+1)][int(cell.X)].FillColor = color
// 								cells[yoffset+int(cell.Y+1)][int(cell.X)].Refresh()
// 							}
// 							previousPositions = append(previousPositions, cell)
// 							// previousPositions = removePos(previousPositions, fyne.NewPos(cell.X, cell.Y+1))
// 							// previousPositions = removePos(previousPositions, fyne.NewPos(cell.X, cell.Y+2))
// 						}
// 						// fmt.Println(previousPositions)
// 						//ensureMapInitialized(&//maxMap)
// 						maxY := 0
// 						_, maxY = MaxYPosition(minYForX.data, 1)
// 						// for i := params.Length; i > 0; i-- {
// 						// 	_, maxY = MaxYPosition(minYForX.data, i)
// 						// 	//fmt.Println(gridWidth-rowCounters[maxY], i)
// 						// 	if gridWidth-rowCounters[maxY] >= i {
// 						// 		break
// 						// 	}
// 						// }
// 						//rowCounters[maxY] -= 2
// 						// fmt.Println(y, maxY)
// 						//fmt.Println(currentMin, currentY)
// 						// if currentMin != maxY {
// 						////maxMap[maxX] = maxY

// 						if newY < maxY && rowCounters[maxY] != 0 {
// 							rowCounters[maxY] -= 1
// 						}
// 						if (rowCounters[maxY] == 0 || currentMin != maxY) && newPos == fyne.NewPos(newPos.X, float32(maxY)) {
// 							// fmt.Println("b")
// 							// rowCounters[maxY] = 0
// 							// fmt.Println("New generation since last clearing")
// 							// if newY < maxY {
// 							// 	rowCounters[maxY]++
// 							// }
// 							for _, cell := range allignmentPos {
// 								if cell.Y < float32(maxY) {
// 									// fmt.Println("cell size")
// 									// fmt.Println(cell.Y, maxY)
// 									rowCounters[maxY] += params.Length
// 								}
// 							}
// 						}

// 						currentMin = maxY

// 						rowCounters[maxY] += 1
// 						// /&& newPos == fyne.NewPos(newPos.X, float32(maxY))
// 						if rowCounters[maxY] >= gridWidth-1 {
// 							increaseScore(1)
// 							// for _, cell := range groupCells {
// 							// 	cells[int(cell.Y)][int(cell.X)].FillColor = color
// 							// 	cells[int(cell.Y+1)][int(cell.X)].Refresh()
// 							// 	cells[int(cell.Y+1)][int(cell.X)].FillColor = color
// 							// 	cells[int(cell.Y+1)][int(cell.X)].Refresh()
// 							// }
// 							// fmt.Printf("Row %d is full\n", maxY)
// 							// Clear or update the filled row
// 							// go func() {
// 							//yoffset += 2
// 							time.Sleep(1 * time.Millisecond)
// 							// fmt.Println(previousPositions)
// 							// fmt.Println(maxY)
// 							for x := 0; x < gridWidth; x++ {

// 								cells[yoffset+maxY-1][x].FillColor = rgbaGrayColor
// 								cells[yoffset+maxY-1][x].Refresh()
// 								cells[yoffset+maxY][x].FillColor = rgbaGrayColor
// 								cells[yoffset+maxY][x].Refresh()

// 							}
// 							//fmt.Println(previousPositions)
// 							// fall(cells, previousPositions, color, true, params, bottomLimit)
// 							// }()
// 							delete(rowCounters, maxY)

// 						}
// 					}

// 					fmt.Println(allignmentPos)
// 					if newY != bottomLimit {
// 						cells[yoffset+newY][newX].FillColor = color
// 						cells[yoffset+newY][newX].Refresh()
// 					}
// 				} else {
// 					// Additional logic for non-normal cells
// 					fmt.Println(pos)
// 				}

// 				newGroupCells = append(newGroupCells, newPos)
// 			}
// 		}

// 		groupCells = newGroupCells
// 		if isNormal {
// 			currentGroup = groupCells
// 		}
// 		toBeCleared = make(map[fyne.Position]bool)
// 		if latestKeyEvent.KeyName != fyne.KeySpace {
// 			time.Sleep(500 * time.Millisecond)
// 		}
// 		// } else {
// 		// 	// Handle non-normal cases, if necessary
// 		// }
// 		// fmt.Println(i)
// 	}

// 	if isNormal {
// 		// Reset state for next fall
// 		latestKeyEvent = KeyEvent{}
// 		//keyProcessed = false
// 		return
// 	}
// }

// previousPositions = removePos(previousPositions, fyne.NewPos(float32(yoffset+maxY-3), float32(x)))
// previousPositions = removePos(previousPositions, fyne.NewPos(float32(yoffset+maxY-4), float32(x)))
//fmt.Println(maxY)
//	if containsPos(previousPositions, fyne.NewPos(float32(newY), float32(x))) {
//fmt.Println("A")
// if containsPos(previousPositions, fyne.NewPos(float32(maxY), float32(x))) {
// 	fmt.Println(true)
// }
// if containsPos(previousPositions, fyne.NewPos(float32(maxY-1), float32(x))) {
// 	fmt.Println(true)
// }
// if containsPos(previousPositions, fyne.NewPos(float32(maxY-2), float32(x))) {
// 	fmt.Println(true)
// }
// previousPositions = removePos(previousPositions, fyne.NewPos(float32(maxY)+1, float32(x)))
// previousPositions = removePos(previousPositions, fyne.NewPos(float32(maxY), float32(x)))
// fmt.Println(previousPositions)
// previousPositions = removePos(previousPositions, fyne.NewPos(float32(maxY)+1, float32(x)))
//	}
// fmt.Println(maxY)
// limit = maxY + 1
// if clearLowestPoint != -1 {
// 	newY++
// 	// for _, cell := range previousPositions {
// 	// 	if cell.Y == float32(newY) {
// 	// 		previousPositions = removePos(previousPositions, cell)
// 	// 	}
// 	// }
// 	// newY--
// }
// if newY < len(cells)-1 && newX >= 0 && newX < len(cells[yoffset+newY])-1 && clearLowestPoint {
// 	for _, cell := range previousPositions {
// 		if cell.Y ==
// 	}
// }
// if newY+1 == clearLowestPoint {
// 	for _, cell := range previousPositions {
// 		if cell.Y == float32(newY) {
// 			previousPositions = removePos(previousPositions, cell)
// 		}
// 	}
// }
// fmt.Println(newY, clearLowestPoint)
// previousPositions = removePos(previousPositions, fyne.NewPos(float32(x), float32(maxY-4)))
// previousPositions = removePos(previousPositions, fyne.NewPos(float32(x), float32(maxY-1)))
// previousPositions = removePos(previousPositions, fyne.NewPos(float32(x), float32(maxY)-2))
// previousPositions = removePos(previousPositions, fyne.NewPos(float32(x), float32(maxY)-3))
// previousPositions = removePos(previousPositions, fyne.NewPos(float32(x), float32(newY)-3))
// previousPositions = removePos(previousPositions, fyne.NewPos(float32(x), float32(newY)))
// previousPositions = removePos(previousPositions, fyne.NewPos(float32(x), float32(newY)-2))
// previousPositions = removePos(previousPositions, fyne.NewPos(float32(x), float32(newY-1)))
// cells[yoffset+newY-1][x].FillColor = rgbaGrayColor
// cells[yoffset+newY-1][x].Refresh()
// cells[yoffset+newY][x].FillColor = rgbaGrayColor
// cells[yoffset+newY][x].Refresh()
// previousPositions = removePos(previousPositions, fyne.NewPos(float32(x), float32(maxY-3)))
// previousPositions = removePos(previousPositions, fyne.NewPos(float32(x), float32(maxY-2)))
// previousPositions = removePos(previousPositions, fyne.NewPos(float32(x), float32(maxY-1)))
// // if maxY != bottomLimit-1 {
// previousPositions = removePos(previousPositions, fyne.NewPos(float32(x), float32(maxY)))
// previousPositions = removePos(previousPositions, fyne.NewPos(float32(x), float32(maxY+1)))
//}
// Before clearing previousPositions, find the highest point for each x value
// highestPoints := make(map[int]float32)

// for _, pos := range previousPositions {
// 	x := int(pos.X)
// 	y := pos.Y
// 	if currentY, exists := highestPoints[x]; !exists || y > currentY {
// 		highestPoints[x] = y
// 		previousPositions = append(previousPositions, fyne.NewPos(float32(x), float32(y)))
// 	}
// }

// // Clear the row and update the cells based on the highest points
// previousPositions = []fyne.Position{}
// for i := 0; i < gridWidth; i++ {
// 	previousPositions = append(previousPositions, fyne.NewPos(float32(i), float32(bottomLimit)))
// }

// // Clearing the row and reapplying the highest positions
// for x := 0; x < gridWidth; x++ {
// 	// Clear the row
// 	cells[yoffset+bottomLimit][x].FillColor = rgbaGrayColor
// 	cells[yoffset+bottomLimit][x].Refresh()

// 	// Get the highest Y value for this X
// 	if maxY, exists := highestPoints[x]; exists && maxY < float32(bottomLimit) {
// 		// Move the cell down to the bottom limit
// 		cells[yoffset+int(maxY)][x].FillColor = color // Apply the color to the new position
// 		cells[yoffset+int(maxY)][x].Refresh()
// 		fmt.Printf("Moving cell from Y=%f to bottom limit X=%d\n", maxY, x) // Debug: Show cell movement
// 	}
// }

//	for _, cell := range allignmentPos {
//		fmt.Println(cell)
//		if cell.Y != float32(currentMin) {
//			previousPositions = append(previousPositions, cell)
//		}
//	}
//
// previousPositions = removePos(previousPositions, fyne.NewPos(float32(x), float32(maxY-1)))
// clearLowestPoint = maxY
