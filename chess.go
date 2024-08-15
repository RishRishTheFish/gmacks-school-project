package main

import (
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/notnil/chess"
)

func createGrid(game *chess.Board) *fyne.Container {
	grid := container.NewGridWithColumns(8)

	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			bg := canvas.NewRectangle(color.Gray{0x30})
			if x%2 == y%2 {
				bg.FillColor = color.Gray{0xE0}
			}

			piece := game.Piece(chess.Square(x + (7-y)*8))
			img := canvas.NewImageFromResource(resourceForPiece(piece))
			img.FillMode = canvas.ImageFillContain
			grid.Add(container.NewMax(bg, img))
		}
	}
	return grid
}

func createChess(w fyne.Window) {
	game := chess.NewGame()
	w.Resize(fyne.NewSize(480, 480))

	grid := createGrid(game.Position().Board())
	over := canvas.NewImageFromResource(nil)
	over.Hide()
	w.SetContent(container.NewMax(grid, container.NewWithoutLayout(over)))
	go func() {
		rand.Seed(time.Now().UnixNano())
		for game.Outcome() == chess.NoOutcome {
			time.Sleep(500 * time.Millisecond)
			valid := game.ValidMoves()
			m := valid[rand.Intn(len(valid))]
			move(m, game, grid, over)
		}
	}()
	// w.ShowAndRun()
}

func move(m *chess.Move, game *chess.Game, grid *fyne.Container, over *canvas.Image) {
	// Get the initial position and image
	off := squareToOffset(m.S1())
	cell := grid.Objects[off].(*fyne.Container)
	img := cell.Objects[1].(*canvas.Image)

	// Calculate the absolute position for pos1
	pos1 := img.Position().Add(cell.Position())

	// Update over with the current image properties
	over.Resource = img.Resource
	over.Resize(img.Size())
	over.Move(pos1)
	over.Refresh()
	over.Show()

	// Clear the resource from the original image
	img.Resource = nil
	img.Refresh()

	// Calculate the final position by updating the cell for the destination square
	off = squareToOffset(m.S2())
	cell = grid.Objects[off].(*fyne.Container)

	// Calculate the absolute position for pos2
	pos2 := cell.Position().Add(grid.Position())

	// Animate the movement from pos1 to pos2
	a := canvas.NewPositionAnimation(pos1, pos2, time.Millisecond*500, func(p fyne.Position) {
		over.Move(p)
		over.Refresh()
	})
	a.Start()

	// Wait for the animation to complete
	time.Sleep(time.Millisecond * 500)

	// Perform the move in the game after the animation
	game.Move(m)

	// Hide the over image and refresh the grid
	over.Hide()
	refreshGrid(grid, game.Position().Board())
}

func refreshGrid(grid *fyne.Container, b *chess.Board) {
	y, x := 7, 0
	for _, cell := range grid.Objects {
		if container, ok := cell.(*fyne.Container); ok {
			if len(container.Objects) > 1 {
				if img, ok := container.Objects[1].(*canvas.Image); ok {
					p := b.Piece(chess.Square(x + y*8))
					img.Resource = resourceForPiece(p)
					img.Refresh()
				}
			}
		}

		x++
		if x == 8 {
			x = 0
			y--
		}
	}
}
func squareToOffset(sq chess.Square) int {
	x := sq % 8
	y := 7 - (sq / 8)
	return int(x + y*8)
}
