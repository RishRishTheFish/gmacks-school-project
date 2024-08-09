package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

const (
	gridWidth  = 10
	gridHeight = 20
	cellSize   = 30
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

func clearRow(cells [][]*canvas.Rectangle, y int) {
	rowTally := 0
	fmt.Println("running clear")
	for _, existingCell := range previousPositions {
		fmt.Println("existing cells")
		fmt.Println(existingCell.Y)
		fmt.Println("y:")
		fmt.Println(y)
		if existingCell.Y == float32(y) {
			fmt.Println("matches y (cells)")
			rowTally++
		}
	}
	fmt.Println(rowTally)
	if rowTally == gridWidth {
		fmt.Println("rowTally and gridWidth match")
		for _, existingCell := range previousPositions {
			if existingCell.Y == float32(y) {
				fmt.Println("matched Y")
				cells[int(existingCell.Y)][int(existingCell.X)].FillColor = rgbaGrayColor
				// existingCell.FillColor = rgbaGrayColor
			}
		}
	}
}

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
var currentPos fyne.Position

// var allignment []int
var allignmentPos []fyne.Position

// var groupCells []*fyne.Position

// var tempGroupCells []*fyne.Position
func removePos(positions []fyne.Position, pos fyne.Position) []fyne.Position {
	for i, p := range positions {
		if p == pos {
			return append(positions[:i], positions[i+1:]...)
		}
	}
	return positions
}
func fall(cells [][]*canvas.Rectangle, groupCells []fyne.Position, color color.Color, isNormal bool, params ExtraParams) {

	// Maintain a set of cells to be cleared
	toBeCleared := make(map[fyne.Position]bool)

	limit = 15
	if previousPositionsHasBottom == false {
		previousPositionsHasBottom = true
		for i := 0; i < gridWidth; i++ {
			previousPositions = append(previousPositions, fyne.NewPos(float32(i), float32(limit)+1))
		}
	}
	for i := 0; i < limit; i++ {
		// if isNormal == false {
		// 	for i := 0; i < gridWidth; i++ {

		// 	}
		// }
		if i == 3 && isNormal == true {
			for i := 0; i < gridWidth; i++ {
				var tempCellArr []fyne.Position
				tempCellArr = append(tempCellArr, fyne.NewPos(float32(i), 2))
				// fmt.Println(tempCellArr)
				// fmt.Println("End of tempcellarr")
				go fall(cells, tempCellArr, randomColor(), false, params)
			}
		}
		newGroupCells := []fyne.Position{}

		// Step 1: Mark cells to be cleared
		if isNormal {
			for _, pos := range groupCells {
				x, y := int(pos.X), int(pos.Y)

				// Mark cell for clearing if it's within bounds
				if y < len(cells) && x < len(cells[y]) {
					toBeCleared[pos] = true
				}
			}

			// Step 2: Clear the marked cells
			for pos := range toBeCleared {
				x, y := int(pos.X), int(pos.Y)
				if y < len(cells) && x < len(cells[y]) {
					cells[y][x].FillColor = rgbaGrayColor
					cells[y][x].Refresh()
				}
			}
		}

		// Step 3: Move cells down and track new positions
		for _, pos := range groupCells {
			x, y := int(pos.X), int(pos.Y)
			// Check if cell can move down
			// if !(y+2 < len(cells) && x < len(cells[y+2])) {
			// for i := 0; i < gridWidth; i++ {
			// }
			if y+1 < len(cells) && x < len(cells[y+1]) {
				newPos := fyne.NewPos(float32(x), float32(y+1))
				oldPos := fyne.NewPos(float32(x), float32(y))
				if !isNormal {
					allignmentPos = removePos(allignmentPos, oldPos)
					allignmentPos = append(allignmentPos, newPos)
					// allignment = append(allignment, y)
					// fmt.Println(allignment)
				}
				if containsPos(previousPositions, newPos) {
					// previousPositions = append(previousPositions, fyne.NewPos(float32(x), float32(y)-1))
					if isNormal {
						previousPositions = append(previousPositions, fyne.NewPos(float32(x), float32(y)-1))
						fmt.Println(allignmentPos)
						// allignment = []int{}
						yAllignsArr := []int{}
						for _, cell := range allignmentPos {
							yAllignsArr = append(yAllignsArr, int(cell.Y))
							// cells[int(cell.Y)][int(cell.X)].FillColor = color
							// cells[int(cell.Y)][int(cell.X)].Refresh()
							// cells.FillColor = color
							// cells.Refresh()
						}
						if min(yAllignsArr) == max(yAllignsArr) {

						}
						// for _, y := range yAllignsArr {

						// }
						allignmentPos = []fyne.Position{}
						limit = y
					}
					// previousPositions = append(previousPositions, fyne.NewPos(float32(x), float32(y)-2))
					// fmt.Println("reached end")

					// fallInvis(cells)
					// checkLevel(y)
					// checkLevel(cells)
					// clearRow(cells, y)
					// break
				}

				// Update the cell color and add new position
				//if isNormal && !containsPos(previousPositions, fyne.NewPos(float32(x), float32(y)+3)) {
				if isNormal {
					currentPos = newPos
					cells[y+1][x].FillColor = color
					cells[y+1][x].Refresh()
				}
				// for _, pos := range allignmentPos {
				// 	cells[int(pos.Y)][int(pos.X)].FillColor = randomColor()
				// 	cells[int(pos.Y)][int(pos.X)].Refresh()
				// }
				// cells[y+1][x].FillColor = color
				// cells[y+1][x].Refresh()
				newGroupCells = append(newGroupCells, newPos)
			} else {
			}
		}

		// Step 4: Update groupCells with new positions
		groupCells = newGroupCells

		// Clear `toBeCleared` for the next iteration
		toBeCleared = make(map[fyne.Position]bool)
		// if isNormal {
		time.Sleep(100 * time.Millisecond)
		// }
		// Delay to visualize the falling effect
	}
	// fmt.Println("Stopped")
}

func makeCorner(cells [][]*canvas.Rectangle, randNum int, color color.Color, params ExtraParams) {
	cells[0][randNum].FillColor = color
	cells[0][randNum].Refresh()
	pos1 := fyne.NewPos(float32(randNum), 0)
	groupCells := []fyne.Position{pos1}
	cells[1][randNum].FillColor = color
	cells[1][randNum].Refresh()
	pos2 := fyne.NewPos(float32(randNum), 1)
	groupCells = append(groupCells, pos2)
	time.Sleep(1 * time.Second)
	fall(cells, groupCells, color, true, params)
}

func makeLine(cells [][]*canvas.Rectangle, randNum int, color color.Color, params ExtraParams) {
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
	time.Sleep(1 * time.Second)
	fall(cells, groupCells, color, true, params)
}

func makeSquare(cells [][]*canvas.Rectangle, randNum int, color color.Color, params ExtraParams) {
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
	fall(cells, groupCells, color, true, params)
}

func removeAction(actions []cellsParams, index int) []cellsParams {
	if index < 0 || index >= len(actions) {
		return actions
	}
	return append(actions[:index], actions[index+1:]...)
}

func applyRandomColors(grid *fyne.Container, cells [][]*canvas.Rectangle) {
	// for y := 0; y < len(cells); y++ {

	// var currentCell *canvas.Rectangle
	// var groupCells [][]*canvas.Rectangle

	for {
		params := ExtraParams{
			Length: 10,
			// Initialize other fields if needed
		}
		randNum := rand.Intn(10)
		color := randomColor()
		// currentCell = cells[0][randNum]
		// makeSquare(cells, randNum, color)
		actions := []cellsParams{
			makeSquare,
			// makeLine,
			makeCorner,
		}
		if randNum >= 9 || randNum <= 1 {
			actions = removeAction(actions, 0)
		}

		// Pick a random index
		randomIndex := rand.Intn(max(1, len(actions)))

		// Execute the function at the random index
		if len(actions) > 0 {
			actions[randomIndex](cells, randNum, color, params)
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

// 		// Wait before applying the next color
// 		time.Sleep(1)
// 	}
// }

func createTetris() *fyne.Container {
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
	go applyRandomColors(grid, cells)

	return grid
}
