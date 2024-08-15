package main

import (
	"fmt"
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// var maxYForX = NewConcurrentMap()
var minYForX = NewConcurrentMap()
var currentY int
var currentMin int

// var //maxMap map[int]int

var rowCounters map[int]int
var squaresInRows map[int]int
var totalGenerations int // Track the total number of generations
var dontColor bool

func init() {
	rowCounters = make(map[int]int)
	squaresInRows = make(map[int]int)
}

var latestKeyEvent KeyEvent
var previousKeyEvent KeyEvent
var rowCountersMutex sync.Mutex

var yoffset int

// var clearLowestPoint int
func fallBeta(grid fyne.Container, cells [][]*canvas.Rectangle, groupCells []fyne.Position, color color.Color, isNormal bool, params ExtraParams, limit int) {
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
				//fmt.Println("fallBetaing")
				go fallBeta(grid, cells, tempCellArr, randomColor(), false, params, limit)
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
						// ensureMapInitialized(&maxMap)
						_, maxY := MaxYPosition(minYForX.data, 1)

						for i := params.Length; i > 0; i-- {
							_, maxY = MaxYPosition(minYForX.data, i)
							if gridWidth-rowCounters[maxY] >= i {
								break
							}
						}
						if params.Type == "square" {
							squaresInRows[maxY]++
						}
						// fmt.Println(y, maxY)
						if currentMin != maxY || rowCounters[maxY] <= 1 {
							// rowCounters[maxY] = 0
							// fmt.Println("New generation since last clearing")
							// if newY < maxY {
							// 	rowCounters[maxY]++
							// }
							for _, cell := range allignmentPos {
								// fmt.Println(cell.Y, maxY)
								if cell.Y < float32(maxY) {
									// fmt.Println("cell size")
									// fmt.Println(cell.Y, maxY)
									rowCounters[maxY]++
								}
							}
						}
						//maxMap[maxX] = maxY
						if newY < maxY {
							rowCounters[maxY]--
						}
						currentMin = maxY
						fmt.Println(maxY, rowCounters[maxY])
						// fmt.Println(rowCounters[maxY], maxY)
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
						// fmt.Println("a")
						if rowCounters[maxY] >= gridWidth-squaresInRows[maxY] {
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
								// fallBeta(cells, previousPositions, color, true, params, bottomLimit)
							}()
							delete(rowCounters, maxY)
							delete(squaresInRows, maxY)
						}
						break
						//go makeLine(grid, cells, rand.Intn(gridWidth-1), randomColor(), ExtraParams{})
						//go applyRandomColors(&grid, cells)
						// fmt.Println("b")
					} else {
						//break
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
		// Reset state for next fallBeta
		latestKeyEvent = KeyEvent{}
		//keyProcessed = false
	}
}

var totalSquares int
