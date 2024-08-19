package main

import (
	"image/color"
	"math/rand"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	maxGridSize = 30
	cellSize    = 20
)

type SnakeGame struct {
	grid       [][]*canvas.Rectangle
	snake      []fyne.Position
	direction  fyne.Position
	food       fyne.Position
	gameOver   bool
	gridSize   int
	ticker     *time.Ticker
	stopChan   chan struct{}
	score      int
	scoreLabel *widget.Label
}

func NewSnakeGame(gridSize int) *SnakeGame {
	if gridSize > maxGridSize {
		gridSize = maxGridSize
	}

	g := &SnakeGame{
		gridSize: gridSize,
		grid:     make([][]*canvas.Rectangle, gridSize),
		snake:    []fyne.Position{{X: float32(gridSize / 2), Y: float32(gridSize / 2)}},
		direction: fyne.Position{
			X: 0,
			Y: -1, // Start moving up
		},
		gameOver: false,
		stopChan: make(chan struct{}),
		score:    0,
	}

	for i := range g.grid {
		g.grid[i] = make([]*canvas.Rectangle, gridSize)
		for j := range g.grid[i] {
			g.grid[i][j] = canvas.NewRectangle(g.GetGridCellColor(i, j))
		}
	}

	g.respawnFood() // Place initial food
	g.startTicker()
	return g
}

func (g *SnakeGame) startTicker() {
	if g.ticker != nil {
		g.ticker.Stop()
	}

	g.ticker = time.NewTicker(time.Millisecond * 200)

	go func() {
		for {
			select {
			case <-g.ticker.C:
				g.update()
				if g.gameOver {
					return
				}
			case <-g.stopChan:
				return
			}
		}
	}()
}

func (g *SnakeGame) stopTicker() {
	if g.ticker != nil {
		g.ticker.Stop()
	}
	close(g.stopChan)
	g.stopChan = make(chan struct{}) // Create a new channel for the next game
}

func (g *SnakeGame) respawnFood() {
	g.food = fyne.Position{
		X: float32(rand.Intn(g.gridSize)),
		Y: float32(rand.Intn(g.gridSize)),
	}
	g.updateGridCell(int(g.food.X), int(g.food.Y), theme.PrimaryColor())
}

func (g *SnakeGame) update() {
	if g.gameOver {
		g.stopTicker()
		return
	}

	// Move the snake
	head := g.snake[0]
	newHead := fyne.Position{X: head.X + g.direction.X, Y: head.Y + g.direction.Y}

	// Check for collision with walls
	if newHead.X < 0 || int(newHead.X) >= g.gridSize || newHead.Y < 0 || int(newHead.Y) >= g.gridSize {
		g.gameOver = true
		return
	}

	// Check for collision with itself
	for _, part := range g.snake {
		if part == newHead {
			g.gameOver = true
			return
		}
	}

	// Add new head to the snake
	g.snake = append([]fyne.Position{newHead}, g.snake...)

	if newHead == g.food {
		// Snake eats the food, place new food and increase score
		g.respawnFood()
		g.score++
		g.scoreLabel.SetText("Score: " + strconv.Itoa(g.score))
	} else {
		// Remove tail
		tail := g.snake[len(g.snake)-1]
		g.snake = g.snake[:len(g.snake)-1]
		g.updateGridCell(int(tail.X), int(tail.Y), g.GetGridCellColor(int(tail.X), int(tail.Y)))
	}

	// Update the grid with the snake's new position
	for _, part := range g.snake {
		g.updateGridCell(int(part.X), int(part.Y), color.RGBA{0x00, 0xFF, 0x00, 0xFF}) // Snake color
	}

	// Update food position
	g.updateGridCell(int(g.food.X), int(g.food.Y), color.RGBA{0xFF, 0x00, 0x00, 0xFF}) // Food color
}

func (g *SnakeGame) GetGridCellColor(x, y int) color.Color {
	lightGreen := color.RGBA{R: 0x8F, G: 0xBC, B: 0x8F, A: 0xFF} // Lighter shade of green for the grid
	darkGreen := color.RGBA{R: 0x2E, G: 0x8B, B: 0x57, A: 0xFF}  // Darker shade of green for the grid

	if (x+y)%2 == 0 {
		return lightGreen
	}
	return darkGreen
}

func (g *SnakeGame) updateGridCell(x, y int, color color.Color) {
	g.grid[x][y].FillColor = color
	g.grid[x][y].Refresh()
}

func (g *SnakeGame) reset() {
	// Stop the ticker to avoid multiple game loops
	g.stopTicker()

	// Clear the grid and reset the snake
	for _, row := range g.grid {
		for _, cell := range row {
			cell.FillColor = g.GetGridCellColor(int(cell.Position().X), int(cell.Position().Y))
			cell.Refresh()
		}
	}
	g.snake = []fyne.Position{{X: float32(g.gridSize / 2), Y: float32(g.gridSize / 2)}}
	g.direction = fyne.Position{X: 0, Y: -1}
	g.gameOver = false
	g.score = 0
	g.scoreLabel.SetText("Score: 0")
	g.respawnFood()

	// Restart the ticker and the game loop
	g.startTicker()
}

func (g *SnakeGame) GetGrid() *fyne.Container {
	grid := container.NewGridWithColumns(g.gridSize)
	for y := 0; y < g.gridSize; y++ {
		for x := 0; x < g.gridSize; x++ {
			cellColor := g.GetGridCellColor(x, y)
			cell := canvas.NewRectangle(cellColor)
			cell.SetMinSize(fyne.NewSize(cellSize, cellSize))
			g.grid[x][y] = cell
			grid.Add(cell)
		}
	}
	return grid
}
func createSnake(w fyne.Window, customTheme *CustomTheme) fyne.CanvasObject {
	// Define colors
	backgroundColor := color.RGBA{0x10, 0x10, 0x10, 0xff} // Dark background
	textColor := color.RGBA{0xcc, 0xcc, 0xcc, 0xff}       // Light gray for text

	// Define the size of the content area
	contentSize := fyne.NewSize(400, 400)

	// Create a dark background rectangle
	background := canvas.NewRectangle(backgroundColor)
	background.SetMinSize(contentSize)

	// Initialize the game
	game := NewSnakeGame(int(contentSize.Width / cellSize))

	// Initialize the score label
	scoreLabel := widget.NewLabel("Score: 0")
	scoreLabel.TextStyle = fyne.TextStyle{Bold: true}
	scoreLabel.Alignment = fyne.TextAlignCenter
	game.scoreLabel = scoreLabel

	// Create the grid for the Snake game with alternating green colors
	grid := game.GetGrid()

	// Create a footer with control instructions using canvas.Text
	footerLabel := canvas.NewText("Controls: Arrow keys to move, Space to restart", textColor)
	footerLabel.TextStyle = fyne.TextStyle{Bold: true}
	footerLabel.Alignment = fyne.TextAlignCenter

	// Create a header and footer for the game
	header := container.NewVBox(
		layout.NewSpacer(),
		scoreLabel,
		layout.NewSpacer(),
	)

	footer := container.NewVBox(
		layout.NewSpacer(),
		footerLabel,
		layout.NewSpacer(),
	)

	// Create a container for the grid
	gridContainer := container.NewVBox(
		layout.NewSpacer(),
		grid,
		layout.NewSpacer(),
	)

	// Create the main game container with a background
	gameContainer := container.NewVBox(
		header,
		gridContainer,
		footer,
	)

	// Overlay the game UI on top of the background
	content := container.NewMax(background, gameContainer)

	// Add a sidebar using container.NewBorder
	contentWithSidebar := container.NewBorder(nil, nil, nil, container.NewHBox(getFiller(), getSidebar(w, customTheme)), content)

	// Set the content for the window
	w.SetContent(contentWithSidebar)

	// Key event handling
	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		switch k.Name {
		case fyne.KeyUp:
			game.changeDirection(fyne.Position{X: 0, Y: -1})
		case fyne.KeyDown:
			game.changeDirection(fyne.Position{X: 0, Y: 1})
		case fyne.KeyLeft:
			game.changeDirection(fyne.Position{X: -1, Y: 0})
		case fyne.KeyRight:
			game.changeDirection(fyne.Position{X: 1, Y: 0})
		case fyne.KeySpace:
			w.SetContent(createSnake(w, customTheme))
		}
	})

	return contentWithSidebar
}

func (g *SnakeGame) changeDirection(newDirection fyne.Position) {
	if (newDirection.X == -g.direction.X && newDirection.Y == 0) || (newDirection.Y == -g.direction.Y && newDirection.X == 0) {
		// Prevent reversing direction
		return
	}
	g.direction = newDirection
}

// func main() {
// 	myApp := fyne.NewApp()
// 	myWindow := myApp.NewWindow("Snake Game")
// 	myWindow.Resize(fyne.NewSize(600, 600))

// 	snakeContent := createSnake(myWindow)
// 	myWindow.SetContent(snakeContent)

// 	myWindow.ShowAndRun()
// }
