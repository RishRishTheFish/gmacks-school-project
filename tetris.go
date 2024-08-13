package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
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
	// Add other fields if needed
}

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
var maxMap map[int]int

var rowCounters map[int]int
var totalGenerations int // Track the total number of generations
var dontColor bool

func init() {
	rowCounters = make(map[int]int)
}

var latestKeyEvent KeyEvent
var previousKeyEvent KeyEvent
var rowCountersMutex sync.Mutex

func fall(grid fyne.Container, cells [][]*canvas.Rectangle, groupCells []fyne.Position, color color.Color, isNormal bool, params ExtraParams, limit int) {
	var doesContainPos bool

	if isNormal {
		globalLimit = 0
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
		if i == 3 && isNormal {
			for j := 0; j < gridWidth; j++ {
				var tempCellArr []fyne.Position
				tempCellArr = append(tempCellArr, fyne.NewPos(float32(j), 1))
				//fmt.Println("falling")
				go fall(grid, cells, tempCellArr, randomColor(), false, params, limit)
			}
		}

		// Update positions based on key events
		//if isNormal {
		xoffset := handleKeyEvents(isNormal)
		newGroupCells := []fyne.Position{}

		// Clear previous positions
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
						for _, cell := range groupCells {
							//cell.FillColor = color
							//cell.Refresh()
							if len(cells) >= int(cell.Y) && len(cells[int(cell.Y)]) >= int(cell.X) {
								cells[int(cell.Y)][int(cell.X)].FillColor = color
								cells[int(cell.Y)][int(cell.X)].Refresh()
							}
							previousPositions = append(previousPositions, cell)
						}
						ensureMapInitialized(&maxMap)
						maxX, maxY := MaxYPosition(minYForX.data)
						// fmt.Println(y, maxY)
						if currentMin != maxY {
							// rowCounters[maxY] = 0
							// fmt.Println("New generation since last clearing")
							// if newY < maxY {
							// 	rowCounters[maxY]++
							// }
							for _, cell := range allignmentPos {
								if cell.Y < float32(maxY) {
									// fmt.Println("cell size")
									// fmt.Println(cell.Y, maxY)
									rowCounters[maxY]++
								}
							}
						}
						maxMap[maxX] = maxY
						if newY < maxY {
							rowCounters[maxY]--
						}
						currentMin = maxY
						fmt.Println(rowCounters[maxY], maxY)
						// cells[maxY][maxX].FillColor = rgbaRedColor
						// cells[maxY][maxX].Refresh()
						// fmt.Println("a")
						// Update rowCounters
						rowCounters[maxY]++
						//fmt.Println(rowCounters[maxY])
						// Use mutex to synchronize access to rowCounters
						rowCountersMutex.Lock()
						defer rowCountersMutex.Unlock()
						// Check for filled rows
						fmt.Println("a")
						if rowCounters[maxY] >= gridWidth-1 {
							increaseScore(1)
							// fmt.Printf("Row %d is full\n", maxY)
							// Clear or update the filled row
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
								// fall(cells, previousPositions, color, true, params, bottomLimit)
							}()
							delete(rowCounters, maxY)

						}
						//go makeLine(grid, cells, rand.Intn(gridWidth-1), randomColor(), ExtraParams{})
						//go applyRandomColors(&grid, cells)
						// fmt.Println("b")
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
		// } else {
		// 	// Handle non-normal cases, if necessary
		// }
	}

	if isNormal {
		// Reset state for next fall
		latestKeyEvent = KeyEvent{}
		//keyProcessed = false
	}
}
func moveBlocksDown(cells [][]*canvas.Rectangle, clearedRow int) {
	// Iterate from the row above the cleared row up to the top of the grid
	for y := clearedRow - 1; y >= 0; y-- {
		for x := 0; x < gridWidth; x++ {
			// Check if there is a block in the current cell
			if cells[y][x].FillColor != rgbaGrayColor { // Assuming rgbaGrayColor is the color for empty cells
				// Move block to the row below
				newY := y + 1
				if newY < len(cells) {
					cells[newY][x].FillColor = cells[y][x].FillColor
					cells[newY][x].Refresh()
					cells[y][x].FillColor = rgbaGrayColor
					cells[y][x].Refresh()
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

var positionsRWMutex sync.RWMutex

// var positions sync.Map

func MaxYPosition(positions map[int]int) (int, int) {
	// Lock the mutex for reading
	// positionsRWMutex.RLock()
	// Ensure that the mutex is unlocked when the function returns

	// copy()

	if len(positions) == 0 {
		return 0, 0 // Return a default value if the map is empty
	}

	var maxX int
	maxY := positions[maxX]
	positionsRWMutex.RLock()
	defer positionsRWMutex.RUnlock()
	for x, y := range positions {
		if y > maxY {
			maxY = y
			maxX = x
		}
	}

	return maxX, maxY
}

// Example function that modifies the map
//
//	func updatePositions(positions map[int]int, x int, y int) {
//		// Lock the mutex before modifying the map
//		positionsMutex.Lock()
//		positions[x] = y
//		positionsMutex.Unlock()
//	}
func makeCorner(
	grid fyne.Container,
	cells [][]*canvas.Rectangle,
	randNum int,
	color color.Color,
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
	fall(
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

func makeLine(
	grid fyne.Container,
	cells [][]*canvas.Rectangle,
	randNum int,
	color color.Color,
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
	fall(
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

func makeSquare(
	grid fyne.Container,
	cells [][]*canvas.Rectangle,
	randNum int,
	color color.Color,
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
	cells[1][randNum+1].FillColor = color
	cells[1][randNum+1].Refresh()
	pos3 := fyne.NewPos(float32(randNum+1), 1)
	groupCells = append(groupCells, pos3)
	cells[0][randNum+1].FillColor = color
	cells[0][randNum+1].Refresh()
	pos4 := fyne.NewPos(float32(randNum+1), 0)
	groupCells = append(groupCells, pos4)
	time.Sleep(1 * time.Second)
	fall(
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

func removeAction(actions []cellsParams, index int) []cellsParams {
	if index < 0 || index >= len(actions) {
		return actions
	}
	return append(actions[:index], actions[index+1:]...)
}

func applyRandomColors(
	grid *fyne.Container,
	cells [][]*canvas.Rectangle,
	//keyEventChannel KeyEvent
) {
	// for y := 0; y < len(cells); y++ {

	// var currentCell *canvas.Rectangle
	// var groupCells [][]*canvas.Rectangle

	for {
		params := ExtraParams{
			Length: 10,
			// Initialize other fields if needed
		}
		randNum := rand.Intn(gridWidth)
		color := randomColor()
		// currentCell = cells[0][randNum]
		// makeSquare(cells, randNum, color)
		actions := []cellsParams{
			makeSquare,
			makeLine,
			//makeCorner,
		}
		// if randNum >= 9 || randNum <= 1 {
		// 	actions = removeAction(actions, 0)
		// }

		// Pick a random index
		randomIndex := rand.Intn(max(1, len(actions)))

		// Execute the function at the random index
		if len(actions) > 0 {
			totalGenerations++
			actions[randomIndex](
				*grid,
				cells,
				randNum,
				color,
				params,
				//keyEventChannel
			)
		}
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
		cells[y] = make([]*canvas.Rectangle, gridWidth) // Initialize the inner slice
	}

	grid := container.NewGridWithColumns(gridWidth)

	for y := 0; y < gridHeight; y++ {
		for x := 0; x < gridWidth; x++ {
			bg := canvas.NewRectangle(cellColor)
			bg.SetMinSize(fyne.NewSize(cellSize, cellSize)) // Use the same size for both width and height
			cells[y][x] = bg
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
