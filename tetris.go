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
var limitsWithCords []fyne.Position

func fall(cells [][]*canvas.Rectangle, groupCells []fyne.Position, color color.Color, params ExtraParams) {
	toBeCleared := make(map[fyne.Position]bool)
	occupied := make(map[fyne.Position]bool)

	for _, pos := range limitsWithCords {
		occupied[pos] = true
	}

	limit := 15
	for i := 0; i < limit; i++ {
		newGroupCells := []fyne.Position{}
		collisionDetected := false

		fmt.Printf("Iteration %d\n", i)
		fmt.Printf("Current groupCells: %v\n", groupCells)
		fmt.Printf("Limits with Cords: %v\n", limitsWithCords)
		fmt.Printf("Occupied: %v\n", occupied)

		// Step 1: Mark cells to be cleared
		for _, pos := range groupCells {
			x, y := int(pos.X), int(pos.Y)
			if y < len(cells) && x < len(cells[y]) {
				toBeCleared[pos] = true
				fmt.Printf("Checking cell at (%d, %d) to be cleared\n", x, y)
			}
		}

		// Step 2: Clear the marked cells
		for pos := range toBeCleared {
			x, y := int(pos.X), int(pos.Y)
			if y < len(cells) && x < len(cells[y]) {
				fmt.Printf("Clearing cell at (%d, %d)\n", x, y)
				cells[y][x].FillColor = rgbaGrayColor
				cells[y][x].Refresh()
			}
		}

		// Step 3: Move cells down and track new positions
		for _, pos := range groupCells {
			x, y := int(pos.X), int(pos.Y)
			newY := y + 1
			fmt.Printf("Moving cell at (%d, %d) down\n", x, y)

			if newY < len(cells) {
				if x < len(cells[newY]) {
					newPos := fyne.NewPos(float32(x), float32(newY))
					if _, exists := occupied[newPos]; exists {
						fmt.Printf("Collision detected at (%d, %d)\n", x, newY)
						collisionDetected = true
						limitsWithCords = append(limitsWithCords, pos)
						occupied[pos] = true
						continue
					}

					cells[newY][x].FillColor = color
					cells[newY][x].Refresh()
					newGroupCells = append(newGroupCells, newPos)
					occupied[newPos] = true
				} else {
					fmt.Printf("Final position at (%d, %d)\n", x, y)
					limitsWithCords = append(limitsWithCords, pos)
					occupied[pos] = true
					collisionDetected = true
				}
			} else {
				fmt.Printf("Final position at (%d, %d)\n", x, y)
				limitsWithCords = append(limitsWithCords, pos)
				occupied[pos] = true
				collisionDetected = true
			}
		}

		groupCells = newGroupCells
		toBeCleared = make(map[fyne.Position]bool)

		if collisionDetected {
			break
		}

		time.Sleep(1 * time.Second) // Delay to visualize the falling effect
	}
}

var groupCells []*fyne.Position

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
	fall(cells, groupCells, color, params)
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
	fall(cells, groupCells, color, params)
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
	fall(cells, groupCells, color, params)
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
			makeLine,
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
