package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/notnil/chess"
)

var (
	grid *fyne.Container
	over *canvas.Image
	win  fyne.Window // Changed from *fyne.Window to fyne.Window
)

func createGrid(game *chess.Game) *fyne.Container {
	var cells []fyne.CanvasObject

	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			bg := canvas.NewRectangle(color.Gray{0x30})
			if x%2 == y%2 {
				bg.FillColor = color.Gray{0xE0}
			}

			piece := newPeice(game, chess.Square(x+y*8))
			cells = append(cells, container.NewMax(bg, piece.image))
		}
	}
	return container.New(&boardLayout{}, cells...)
}
func move(m *chess.Move, game *chess.Game, grid *fyne.Container, over *canvas.Image) {
	off := squareToOffset(m.S1())
	cell := grid.Objects[off].(*fyne.Container)

	imgObj := cell.Objects[1]
	img, ok := imgObj.(*canvas.Image)
	if !ok {
		fmt.Println("Error: Expected *canvas.Image but got different type")
		return
	}

	pos1 := img.Position().Add(cell.Position())
	// fmt.Println(img.Resource)
	over.Resource = img.Resource
	over.Resize(img.Size())
	over.Move(pos1)
	over.Refresh()
	over.Show()

	img.Resource = nil
	img.Refresh()

	off = squareToOffset(m.S2())
	cell = grid.Objects[off].(*fyne.Container)

	pos2 := cell.Position().Add(grid.Position())

	a := canvas.NewPositionAnimation(pos1, pos2, time.Millisecond*500, func(p fyne.Position) {
		over.Move(p)
		over.Refresh()
	})
	a.Start()

	time.Sleep(time.Millisecond * 500)

	game.Move(m)

	over.Hide()
	refreshGrid(grid, game.Position().Board())
}

func createChess(w fyne.Window) {
	game := chess.NewGame()
	win = w // Assign the window to the global variable

	win.Resize(fyne.NewSize(480, 480))

	grid := createGrid(game)
	over = canvas.NewImageFromResource(nil)
	over.Hide()
	win.SetContent(container.NewMax(grid, container.NewWithoutLayout(over)))
	go func() {
		rand.Seed(time.Now().UnixNano())
		for game.Outcome() == chess.NoOutcome {
			time.Sleep(500 * time.Millisecond)
			valid := game.ValidMoves()
			m := valid[rand.Intn(len(valid))]
			move(m, game, grid, over)
		}
	}()
}
func refreshGrid(grid *fyne.Container, b *chess.Board) {
	y, x := 7, 0
	var img *canvas.Image // Variable to hold the image temporarily

	for _, cell := range grid.Objects {
		fmt.Println(cell)

		// If cell is a *fyne.Container, handle objects
		if container, ok := cell.(*fyne.Container); ok {
			// Check if the container has more than one object (background + image)
			if len(container.Objects) > 1 {
				// Store the image object before removing it
				if imgObj, ok := container.Objects[1].(*canvas.Image); ok {
					img = imgObj
				}

				// Remove the custom object (e.g., img) from the container
				container.Objects = container.Objects[:1]
			}

			// Now do the type check after removal
			if len(container.Objects) == 1 {
				// Re-add the image object if it was removed
				if img != nil {
					container.Objects = append(container.Objects, img)
				}

				// Get the piece on the board for the current square
				piece := b.Piece(chess.Square(x + y*8))
				pieceType := piece.Type()
				pieceColor := piece.Color()

				// Get the resource (image) for this piece
				resource := resourceForPiece(pieceColor, pieceType)
				if resource == nil {
					fmt.Println("Warning: Resource for piece is nil")
				}

				// Set the image resource for the piece
				if img != nil {
					img.Resource = resource
					img.Refresh()
					// fmt.Println(img.Resource)
				}
			}
		}

		// Move to the next square on the board
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

// func createGrid(game *chess.Game) *fyne.Container {
// 	var cells []fyne.CanvasObject

// 	for y := 0; y < 8; y++ {
// 		for x := 0; x < 8; x++ {
// 			bg := canvas.NewRectangle(color.Gray{0x30})
// 			if x%2 == y%2 {
// 				bg.FillColor = color.Gray{0xE0}
// 			}

// 			piece := newPeice(game, chess.Square(x+y*8))
// 			cells = append(cells, container.NewMax(bg, piece))
// 		}
// 	}
// 	return container.New(&boardLayout{}, cells...)
// }
