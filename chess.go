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
	fmt.Printf("imgObj: %+v\n", imgObj) // Print the details of imgObj

	//img, ok := imgObj.
	img, ok := imgObj.(*canvas.Image)
	if !ok {
		fmt.Println("Error: Expected *canvas.Image but got different type")
		return
	}

	fmt.Printf("Image: %+v\n", img)                   // Print the details of img
	fmt.Printf("Image Resource: %+v\n", img.Resource) // Print the image resource

	pos1 := img.Position().Add(cell.Position())
	over.Resource = img.Resource
	over.Resize(img.Size())
	over.Move(pos1)
	over.Refresh()
	over.Show()

	// Clear resource of the original image
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
	refreshGrid(grid, game.Position().Board(), game)
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
func refreshGrid(grid *fyne.Container, b *chess.Board, game *chess.Game) {
	y, x := 7, 0

	for _, cell := range grid.Objects {
		// Check if the cell is a *fyne.Container
		if container, ok := cell.(*fyne.Container); ok {
			// Find and remove any existing image object
			var img *canvas.Image
			for i := len(container.Objects) - 1; i >= 0; i-- {
				if _, ok := container.Objects[i].(*canvas.Image); ok {
					// Remove the image object
					container.Objects = append(container.Objects[:i], container.Objects[i+1:]...)
				}
			}

			// Get the piece on the board for the current square
			square := chess.Square(x + y*8)
			piece := b.Piece(square)
			// if piece == nil {
			// 	// If there's no piece, continue to the next cell
			// 	x++
			// 	if x == 8 {
			// 		x = 0
			// 		y--
			// 	}
			// 	continue
			// }

			// Get the resource (image) for this piece
			resource := resourceForPiece(piece.Color(), piece.Type())
			if resource == nil {
				fmt.Println("Warning: Resource for piece is nil")
				x++
				if x == 8 {
					x = 0
					y--
				}
				continue
			}

			// Create or update the image for the piece
			img = canvas.NewImageFromResource(resource)
			img.FillMode = canvas.ImageFillContain

			// Create a new piece widget for the current square
			pieceWidget := newPeice(game, square)
			pieceWidget.image = img

			// Add the new image and piece widget to the container
			container.Objects = append(container.Objects, img)
			container.Objects = append(container.Objects, pieceWidget)

			// Refresh the container
			container.Refresh()
		}

		// Move to the next square on the board
		x++
		if x == 8 {
			x = 0
			y--
		}
	}
	fmt.Println(grid.Objects)
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
