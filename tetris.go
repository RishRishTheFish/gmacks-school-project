package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

const (
	gridWidth   = 10
	gridHeight  = 20
	bottomLimit = 15
	cellSize    = 30
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

// func createTetrisGrid(grid *fyne.Container, cells [][]*canvas.Rectangle) (*fyne.Container, [][]*canvas.Rectangle) {
// 	// grid := container.NewGridWithColumns(gridWidth)
// 	// cells := make([][]*canvas.Rectangle, gridHeight)

//		for y := 0; y < gridHeight; y++ {
//			cells[y] = make([]*canvas.Rectangle, gridWidth) // Initialize the inner slice
//			for x := 0; x < gridWidth; x++ {
//				// bg := canvas.NewRectangle(color.Gray{0x30})
//				// bg.SetMinSize(fyne.NewSize(cellSize, cellSize))
//				// cells[y][x] = bg
//				// grid.Add(bg)
//				cells[y][x]
//			}
//		}
//		return grid, cells
//	}
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

type cellsParams func(cells [][]*canvas.Rectangle, randNum int, color color.Color, extraParams ExtraParams)

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

// func contains(slice []int, value int) bool {
// 	if slice > 1 {
// 		for _, v := range slice {
// 			if v == value {
// 				return true
// 			}
// 		}
// 	}
// 	return false
// }

// func checkLimits() {

// }
// var limitsOnX = make([]int, 0, 10)
// var limitsWithCords []fyne.Position
//
//	window.Canvas().AddShortcut(&fyne.KeyEvent{
//		Name: key.NameUp, // Example: Handle Up Arrow key
//	}, func() {
//
//		label.SetText("Up Arrow key pressed")
//	})
var previousPositions []fyne.Position

// func clearRow(cells [][]*canvas.Rectangle, y int) {
// 	rowTally := 0
// 	fmt.Println("running clear")
// 	for _, existingCell := range previousPositions {
// 		fmt.Println("existing cells")
// 		fmt.Println(existingCell.Y)
// 		fmt.Println("y:")
// 		fmt.Println(y)
// 		if existingCell.Y == float32(y) {
// 			fmt.Println("matches y (cells)")
// 			rowTally++
// 		}
// 	}
// 	fmt.Println(rowTally)
// 	if rowTally == gridWidth {
// 		fmt.Println("rowTally and gridWidth match")
// 		for _, existingCell := range previousPositions {
// 			if existingCell.Y == float32(y) {
// 				fmt.Println("matched Y")
// 				cells[int(existingCell.Y)][int(existingCell.X)].FillColor = rgbaGrayColor
// 				// existingCell.FillColor = rgbaGrayColor
// 			}
// 		}
// 	}
// }

// func clearRow(cells [][]*canvas.Rectangle, y int) {
// 	rowTally := 0
// 	fmt.Println("running clear")
// 	for _, existingCell := range cells[y] {
// 		// fmt.Println("existing cells")
// 		// fmt.Println(existingCell.Y)
// 		// fmt.Println("y:")
// 		// fmt.Println(y)
// 		if existingCell.Position().Y == float32(y) {
// 			fmt.Println("matches y (cells)")
// 			rowTally++
// 		}
// 	}
// 	fmt.Println(rowTally)
// 	if rowTally == gridWidth {
// 		fmt.Println("rowTally and gridWidth match")
// 		for _, existingCell := range cells[y] {
// 			if existingCell.Position().Y == float32(y) {
// 				fmt.Println("matched Y")
// 				cells[int(existingCell.Position().Y)][int(existingCell.Position().X)].FillColor = rgbaGrayColor
// 				// existingCell.FillColor = rgbaGrayColor
// 			}
// 		}
// 	}

// }
var limit int
var previousPositionsHasBottom bool

//	func checkLevel(cells [][]*canvas.Rectangle) {
//		isNotLevel := true
//		// fmt.Println("checking level")
//		// fmt.Println(previousPositions)
//		// for x, _ := range cells[limit] {
//
//			pos := fyne.NewPos(float32(x), float32(limit))
//			// fmt.Println(pos)
//			if containsPos(previousPositions, pos) {
//				isNotLevel = false
//			}
//		}
//		// fmt.Println()
//		if !isNotLevel {
//			fmt.Println("Level at ", limit)
//		}
//	}

// Define the global map to store matching rows
// var matchingRows = make(map[int]bool)

// // checkLevel function to check for matching rows
// func checkLevel(y int) {
// 	var preBlockInstances int
// 	for _, pos := range previousPositions {
// 		if pos.Y == float32(y+1) {
// 			preBlockInstances++
// 		}
// 	}

// 	if preBlockInstances == gridWidth {
// 		fmt.Println("Matching row")

// 		// Check if y value already exists in the map
// 		if _, exists := matchingRows[y]; !exists {
// 			matchingRows[y] = true
// 		}

//			// Print something if y value is greater than 5
//			if y > 5 {
//				fmt.Println("Y value exceeds 5:", y)
//			}
//		}
//	}
// var fakeGroupCells []fyne.Position

// func fallInvis(cells [][]*canvas.Rectangle) {
// 	var doesNotMatch bool
// 	for i := 0; i < gridWidth; i++ {
// 		fakeGroupCells = append(fakeGroupCells, fyne.NewPos(float32(i), float32(limit-1)))
// 	}
// 	for i := limit - 1; i < limit; i++ {
// 		// fmt.Println(i)
// 		newGroupCells := []fyne.Position{}
// 		for _, pos := range fakeGroupCells {
// 			x, y := int(pos.X), int(pos.Y)
// 			newPos := fyne.NewPos(float32(x), float32(y+1))
// 			// Check if cell can move down
// 			if y+1 < len(cells) && x < len(cells[y+1]) {
// 				// fmt.Println("A")
// 				// if !containsPos(previousPositions, newPos) {
// 				// 	// fmt.Println("B")
// 				// 	doesNotMatch = true
// 				// 	// previousPositions = append(previousPositions, fyne.NewPos(float32(x), float32(y)-1))
// 				// 	// previousPositions = append(previousPositions, fyne.NewPos(float32(x), float32(y)-1))
// 				// 	// previousPositions = append(previousPositions, fyne.NewPos(float32(x), float32(y)-2))
// 				// 	// fmt.Println("reached end")
// 				// 	// limit = y
// 				// 	// clearRow(cells, y)
// 				// 	// break
// 				// } else {
// 				// 	// fmt.Println("C")
// 				// }

// 				// Update the cell color and add new position
// 				// cells[y+1][x].FillColor = color
// 				// cells[y+1][x].Refresh()
// 				// cells[y][x].FillColor = rgbaGrayColor
// 				// cells[y][x].Refresh()
// 				// if containsPos(previousPositions, newPos) {
// 				// 	// previousPositions = append(previousPositions, fyne.NewPos(float32(x), float32(y)-1))
// 				// 	// previousPositions = append(previousPositions, fyne.NewPos(float32(x), float32(y)-1))
// 				// 	// previousPositions = append(previousPositions, fyne.NewPos(float32(x), float32(y)-2))
// 				// 	// fmt.Println("reached end")
// 				// 	doesNotMatch = true
// 				// 	limit = y
// 				// 	// fmt.Println(x)
// 				// 	// fallInvis(cells)
// 				// 	// checkLevel(y)
// 				// 	// checkLevel(cells)
// 				// 	// clearRow(cells, y)
// 				// 	// break
// 				// }
// 				newGroupCells = append(newGroupCells, newPos)
// 			} else {
// 				if containsPos(previousPositions, newPos) {
// 					// previousPositions = append(previousPositions, fyne.NewPos(float32(x), float32(y)-1))
// 					// previousPositions = append(previousPositions, fyne.NewPos(float32(x), float32(y)-1))
// 					// previousPositions = append(previousPositions, fyne.NewPos(float32(x), float32(y)-2))
// 					// fmt.Println("reached end")
// 					doesNotMatch = true
// 					limit = y
// 					fmt.Println("test")
// 					// fmt.Println(x)
// 					// fallInvis(cells)
// 					// checkLevel(y)
// 					// checkLevel(cells)
// 					// clearRow(cells, y)
// 					// break
// 				}
// 				// cells[y][x].FillColor = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
// 				// cells[y][x].Refresh()
// 				// fmt.Println("Does not fit")
// 			}
// 		}

//			// Step 4: Update groupCells with new positions
//			fakeGroupCells = newGroupCells
//		}
//		if doesNotMatch == false {
//			fmt.Println("line matches")
//		}
//	}
//
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

// func abs(x int) int {
// 	if x < 0 {
// 		return -x
// 	}
// 	return x
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

// var heightMapping = make(map[int]int)

// func printMinValue() {
// 	// Check if the number of entries in the map is less than gridWidth
// 	if len(heightMapping) < gridWidth {
// 		fmt.Println(0)
// 		return
// 	}

// 	// Find the minimum value in the map
// 	minValue := math.MaxInt32
// 	for _, value := range heightMapping {
// 		if value < minValue {
// 			minValue = value
// 		}
// 	}

//		// Print the minimum value if valid
//		if minValue == math.MaxInt32 {
//			fmt.Println(0)
//		} else {
//			fmt.Println(minValue)
//		}
//	}
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

// }
// func clearRows(cells [][]*canvas.Rectangle) {
// 	// 	fmt.Println(heightMapping)
// 	minimum, _ := findMinValue(heightMapping)
// 	// maximum, _ := findMaxValue(heightMapping)
// 	// fmt.Println(maximum)

// 	fmt.Println(minimum)
// 	// // if minimum != 1 {
// 	// // fmt.Println("A")
// 	// // for i := minimum; i == 0; i-- {
// 	// for _, cell := range cells[maximum] {
// 	// 	cell.FillColor = rgbaRedColor
// 	// 	cell.Refresh()
// 	// }
// 	// for i := maximum; i < gridHeight; i++ {
// 	// 	// for _, cell := range cells[i] {
// 	// 	// }
// 	// 	//cells[i][i].FillColor = rgbaRedColor
// 	// 	// fmt.Println(20 - max(i, 1))
// 	// 	// for _, cell := range cells[i] {
// 	// 	// 	// if cell.Position().Y > float32(maximum) {
// 	// 	// 	// 	cell.FillColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
// 	// 	// 	// 	cell.Refresh()
// 	// 	// 	// }
// 	// 	// 	// cell.FillColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
// 	// 	// 	// cell.Refresh()
// 	// 	// }
// 	// }

// 	//}
// }

// func clearRows(cells [][]*canvas.Rectangle) {
// 	// for _, height := range heightMapping {
// 	// fmt.Println(min(heightMapping))
// 	// minimum, _ := findMinValue(heightMapping)
// 	// fmt.Println(minimum)
// 	// }
// 	// printMinValue()
// 	// Track filled rows
// 	// filledRows := make(map[int]bool)

//		// // Find the maximum height for each column
//		// for _, height := range heightMapping {
//		// 	for y := 0; y <= height; y++ {
//		// 		filledRows[y] = true
//		// 	}
//		//}
//		// fmt.Println(filledRows)
//		// // Check if all cells in each row are filled
//		// for y := range filledRows {
//		// 	rowFilled := true
//		// 	for x := 0; x < len(cells[y]); x++ {
//		// 		if cells[y][x] == nil || cells[y][x].FillColor != rgbaGrayColor {
//		// 			rowFilled = false
//		// 			break
//		// 		}
//		// 	}
//		// 	if rowFilled {
//		// 		// Clear the row
//		// 		for x := 0; x < len(cells[y]); x++ {
//		// 			cells[y][x].FillColor = rgbaGrayColor
//		// 			cells[y][x].Refresh()
//		// 		}
//		// 	}
//		// }
//	}
//
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

// var updateToNewKey bool

func fall(cells [][]*canvas.Rectangle,
	groupCells []fyne.Position,
	color color.Color,
	isNormal bool,
	params ExtraParams,
	limit int,
	//keyEventChannel KeyEvent
) {
	// go func() {
	var doesContainPos bool
	xoffset := 0
	if isNormal {
		globalLimit = 0
		if latestKeyEvent.KeyName != previousKeyEvent.KeyName {
			previousKeyEvent.KeyName = latestKeyEvent.KeyName
			switch previousKeyEvent.KeyName {
			case fyne.KeyLeft:
				xoffset = -1
			}

		}
	}

	toBeCleared := make(map[fyne.Position]bool)

	if !previousPositionsHasBottom {
		previousPositionsHasBottom = true
		for i := 0; i < gridWidth; i++ {
			previousPositions = append(previousPositions, fyne.NewPos(float32(i), float32(limit)+1))
		}
	}

	for i := 0; i < limit; i++ {
		if i == 3 && isNormal {
			for i := 0; i < gridWidth; i++ {
				var tempCellArr []fyne.Position
				tempCellArr = append(tempCellArr, fyne.NewPos(float32(i), 1))
				go fall(cells, tempCellArr, randomColor(), false, params, limit)
			}
		}

		newGroupCells := []fyne.Position{}
		if isNormal {
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
		}

		for _, pos := range groupCells {
			x, y := int(pos.X), int(pos.Y)
			if y+1 < len(cells) && x < len(cells[y+1]) {
				newPos := fyne.NewPos(float32(x+xoffset), float32(y+1))
				oldPos := fyne.NewPos(float32(x), float32(y))

				if xoffset != 0 {
					xoffset = 0
				}
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
						minYForX.Set(x, y+1)

					}
				}

				if isNormal {
					if doesContainPos {
						for _, cell := range groupCells {
							previousPositions = append(previousPositions, cell)
						}
						ensureMapInitialized(&maxMap)
						maxX, maxY := MaxYPosition(minYForX.data)
						//fmt.Println(y, maxY)
						if currentMin > maxY {
							//	fmt.Println("New generation since last clearing")

							dontColor = true
							for _, cell := range allignmentPos {
								if cell.Y < float32(maxY) {
									//fmt.Println("cell size")
									//fmt.Println(cell.Y, maxY)
									rowCounters[maxY]++
								}
							}
						}
						if y+1 < maxY {
							// fmt.Println("A")
							rowCounters[maxY]--
						}
						maxMap[maxX] = maxY

						currentMin = maxY

						cells[maxY][maxX].FillColor = rgbaRedColor
						cells[maxY][maxX].Refresh()

						rowCounters[maxY]++
						fmt.Println(rowCounters[maxY])
						// Check for filled rows
						if rowCounters[maxY] >= gridWidth-1 {
							//fmt.Printf("Row %d is full\n", maxY)
							// Clear or update the filled row
							for x := 0; x < gridWidth; x++ {
								// fmt.Println("clearing")
								cells[maxY-1][x].FillColor = rgbaGrayColor
								cells[maxY-1][x].Refresh()
								cells[maxY][x].FillColor = rgbaGrayColor
								cells[maxY][x].Refresh()
							}
							delete(rowCounters, maxY)
						}
					}
					cells[y+1][x].FillColor = color
					cells[y+1][x].Refresh()
				}

				newGroupCells = append(newGroupCells, newPos)
			}
		}

		groupCells = newGroupCells
		if isNormal {
			// totalGenerations++ // Increment totalGenerations each time fall is run
			currentGroup = groupCells
		}
		toBeCleared = make(map[fyne.Position]bool)

		time.Sleep(1 * time.Millisecond)
	}
}

//}

func ensureMapInitialized(m *map[int]int) {
	if *m == nil {
		*m = make(map[int]int)
	}
}

func MaxYPosition(positions map[int]int) (int, int) {
	if len(positions) == 0 {
		return 0, 0 // Return a default value if the map is empty
	}

	var maxX int
	maxY := positions[maxX]

	for x, y := range positions {
		if y > maxY {
			maxY = y
			maxX = x
		}
	}

	return maxX, maxY
}

// func CheckFilledRows(maxMap map[int]int, gridWidth int) []int {
// 	filledRows := []int{}

// 	minY, maxY := findMinAndMaxY(maxMap)

// 	for row := minY; row <= maxY; row++ {
// 		filled := true
// 		for x := 0; x < gridWidth; x++ {
// 			if y, exists := maxMap[x]; !exists || y != row {
// 				filled = false
// 				break
// 			}
// 		}
// 		if filled {
// 			filledRows = append(filledRows, row)
// 		}
// 	}

// 	return filledRows
// }

// func findMinAndMaxY(maxMap map[int]int) (int, int) {
// 	minY, maxY := int(^uint(0)>>1), -int(^uint(0)>>1) // Initialize min and max

// 	for _, y := range maxMap {
// 		if y < minY {
// 			minY = y
// 		}
// 		if y > maxY {
// 			maxY = y
// 		}
// 	}

// 	return minY, maxY
// }

// fmt.Println(CheckFilledRows(maxMap, gridWidth))
//
//	if currentMin == maxY+2 {
//		fmt.Println("Matches")
//		fmt.Println(currentMin, maxY)
//		// for _, cell := range cells[currentMin] {
//		// 	cell.FillColor = rgbaRedColor
//		// 	cell.Refresh()
//		// }
//		// for i := maxY; i < gridHeight; i++ {
//		// 	// for _, cell := range cells[i] {
//		// 	// 	cell.FillColor = rgbaRedColor
//		// 	// 	cell.Refresh()
//		// 	// }
//		// 	// for _, cell := range cells[i] {
//		// 	// 	cell.FillColor = color
//		// 	// 	cell.Refresh()
//		// 	// }
//		// }
//	}
func makeCorner(
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
	// fmt.Println(randNum)
	// cells[1][randNum+1].FillColor = color
	// cells[1][randNum+1].Refresh()
	// pos3 := fyne.NewPos(float32(randNum+1), 1)
	// groupCells = append(groupCells, pos3)
	time.Sleep(1 * time.Second)
	fall(cells,
		groupCells,
		color,
		true,
		params,
		bottomLimit,
		//keyEventChannel
	)
}

func makeSquare(
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
			// makeSquare,
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
				cells,
				randNum,
				color,
				params,
				//keyEventChannel
			)
		}
	}
	// randNum := rand.Intn(10)
	// for x := 0; x < len(cells[0]); x++ {
	// 	// Generate a random color
	// 	color := randomColor()

	// 	// Update the cell's background color
	// 	cells[0][x].FillColor = color

	// 	// Refresh the rectangle to apply the new color
	// 	cells[0][x].Refresh()
	// }
	// }
}

// func applyRandomColors(grid *fyne.Container, cells [][]*canvas.Rectangle) {
// 	// Define gray color
// 	grayColor := color.Gray{Y: 0x30}

// 	// Convert gray to RGBA
// 	rgbaGrayColor := color.RGBA{
// 		R: grayColor.Y,
// 		G: grayColor.Y,
// 		B: grayColor.Y,
// 		A: 255, // Fully opaque
// 	}

// 	// Keep track of the last modified time for each cell
// 	lastModified := make(map[*canvas.Rectangle]time.Time)
// 	resetDuration := 1 * time.Millisecond // Time to wait before resetting color

// 	for {
// 		// Randomly select a cell
// 		randNum := rand.Intn(len(cells[0]))
// 		currentCell := cells[0][randNum]

// 		// Generate a random color
// 		randomColor := randomColor()

// 		// Update the cell's color
// 		currentCell.FillColor = randomColor
// 		currentCell.Refresh()

// 		// Record the time when the cell was last modified
// 		lastModified[currentCell] = time.Now()

// 		// Wait for a specified time
// 		time.Sleep(1 * time.Second)

// 		// Reset the cell color to gray
// 		currentCell.FillColor = rgbaGrayColor
// 		currentCell.Refresh()

// 		// Check if any cell should be reset
// 		for cell, modTime := range lastModified {
// 			if time.Since(modTime) > resetDuration {
// 				cell.FillColor = rgbaGrayColor
// 				cell.Refresh()
// 				delete(lastModified, cell)
// 			}
// 		}

//			// Wait before applying the next color
//			time.Sleep(1)
//		}
//	}
type KeyEvent struct {
	KeyName fyne.KeyName
}

func listenKeyEvent(
	w fyne.Window,
	//keyEventChannel KeyEvent
) {
	w.Canvas().SetOnTypedKey(func(keyEvent *fyne.KeyEvent) {
		latestKeyEvent = KeyEvent{
			KeyName: keyEvent.Name,
		}
	})

}
func createTetris(w fyne.Window) *fyne.Container {
	//keyEventChannel := make(chan KeyEvent)

	//	w.Canvas().SetOnTypedKey(func(keyEvent *fyne.KeyEvent) {
	//		switch keyEvent.Name {

	//		}
	//	})
	cells := make([][]*canvas.Rectangle, gridHeight)
	for y := 0; y < gridHeight; y++ {
		cells[y] = make([]*canvas.Rectangle, gridWidth) // Initialize the inner slice
	}

	grid := container.NewGridWithColumns(gridWidth)

	for y := 0; y < gridHeight; y++ {
		for x := 0; x < gridWidth; x++ {
			bg := canvas.NewRectangle(color.Gray{0x30})
			bg.SetMinSize(fyne.NewSize(cellSize, cellSize*2))
			cells[y][x] = bg
			grid.Add(bg)
		}
	}
	go applyRandomColors(
		grid,
		cells,
		//keyEventChannel
	)
	go listenKeyEvent(
		w,
		//keyEventChannel
	)

	return grid
}
