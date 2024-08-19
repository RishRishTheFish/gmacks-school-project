package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/notnil/chess"
)

var moveStart chess.Square = chess.NoSquare

//var win *fyne.Window

type peice struct {
	widget.Icon
	//grid
	window    fyne.Window
	container *fyne.Container
	game      *chess.Game
	square    chess.Square
	image     *canvas.Image
}

func newPeice(g *chess.Game, square chess.Square, grid *fyne.Container) *peice {
	p := g.Position().Board().Piece(square)
	pt := p.Type()
	pc := p.Color()
	resource := resourceForPiece(pc, pt)

	if resource == nil {
		//fmt.Println("Warning: Resource for piece is nil")
	}

	img := canvas.NewImageFromResource(resource)
	img.FillMode = canvas.ImageFillContain

	ret := &peice{
		game:      g,
		square:    square,
		image:     img,
		window:    win,
		container: grid,
	}

	ret.ExtendBaseWidget(ret)
	return ret
}

// func (p *peice) Tapped(ev *fyne.PointEvent) {
// 	fmt.Println("Tapped on square:", p.square)
// 	if moveStart == chess.NoSquare {
// 		moveStart = p.square
// 		fmt.Println("Move start set to:", moveStart)
// 		return
// 	}

// 	valid := p.game.ValidMoves()
// 	fmt.Println("Valid moves:", valid)
// 	for _, m := range valid { // Renamed to `m`
// 		fmt.Println("Checking move:", m)
// 		if m.S1() == moveStart && m.S2() == p.square {
// 			fmt.Println("Move is valid:", m)

// 			// Ensure that the offset is within bounds and points to a valid container
// 			off := squareToOffset(m.S1())
// 			if off < 0 || off >= len(grid.Objects) || grid.Objects[off] == nil {
// 				fmt.Println("Error: grid.Objects is nil or offset is out of bounds:", off)
// 				return
// 			}

// 			// Perform the move, passing `m` as a pointer
// 			move(m, p.game, grid, over)
// 			moveStart = chess.NoSquare

// 			// Make the AI move after a delay
// 			go func() {
// 				time.Sleep(time.Second)
// 				randomResponse(p.game)
// 			}()

// 			return
// 		}
// 	}

//		pos := p.game.Position().Board().Piece(p.square)
//		dialog.ShowInformation("Invalid move", "Cannot move piece "+pos.String()+" to square "+p.square.String(), win)
//	}
func showCustomModalPopup(w fyne.Window, title, message string) {
	bgColor := color.RGBA{R: 0x00, G: 0x00, B: 0xFF, A: 0xFF} // Blue

	popupSize := fyne.NewSize(400, 150)

	bg := canvas.NewRectangle(bgColor)
	bg.Resize(popupSize)

	// Create a label for the message
	msgLabel := canvas.NewText(message, color.White)
	msgLabel.TextStyle = fyne.TextStyle{Bold: true}
	msgLabel.TextSize = 16

	var overlay fyne.CanvasObject

	// Create the content container with the background and message
	content := container.NewVBox(
		container.NewMax(bg, msgLabel),
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButton("OK", func() {
				// Close the popup when OK is pressed
				w.Canvas().Overlays().Remove(overlay)
			}),
			layout.NewSpacer(),
		),
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButton("Return to Game", func() {
				// Hide the custom dialog and return to the game
				w.Canvas().Overlays().Remove(overlay)
				// Implement any additional logic to return to the game
			}),
			layout.NewSpacer(),
		),
	)

	// Create a container to center the content
	centeredContent := container.NewCenter(content)

	// Create a custom overlay container
	overlay = container.NewMax(centeredContent)
	overlay.Resize(popupSize) // Set the size of the overlay to match the popup

	// Show the custom dialog as an overlay
	w.Canvas().Overlays().Add(overlay)
}
func mirrorSquare(square chess.Square) chess.Square {
	file := square.File() // Get the file (column)
	rank := square.Rank() // Get the rank (row)

	// Calculate the mirrored file and rank
	mirroredFile := chess.File(7 - int(file)) // 7 is the last index in an 8x8 board (0-indexed)
	mirroredRank := chess.Rank(7 - int(rank))

	// Return the mirrored square
	return chess.NewSquare(mirroredFile, mirroredRank)
}
func (p *peice) Tapped(ev *fyne.PointEvent) {
	w := p.window
	g := p.container
	// fmt.Println(g)
	if moveStart == chess.NoSquare {
		if m := isValidMove(p.square, chess.NoSquare, p.game); m != nil {
			moveStart = p.square
		} else {
			showCustomModalPopup(w, "Invalid move", fmt.Sprintf("Cannot move piece %d",
				p.game.Position().Board().Piece(p.square)))
		}
		return
	}

	if m := isValidMove(moveStart, p.square, p.game); m != nil {
		moveStart = chess.NoSquare
		move(m, p.game, g, over)

		go func() {
			time.Sleep(time.Second)
			// fmt.Println("a")
			if !gameEnded {
				randomResponse(p.game, g)
			}
		}()
		return
	}
	showCustomModalPopup(w, "Invalid move", fmt.Sprintf("Cannot move piece %d to square %v",
		p.game.Position().Board().Piece(moveStart), p.square))

	moveStart = chess.NoSquare
}

// func (p *peice) Tapped(ev *fyne.PointEvent) {
// 	w := p.window
// 	g := p.container

// 	// Check if a move is in progress
// 	if moveStart == chess.NoSquare {
// 		// Check if the move is valid from the current square
// 		if m := isValidMove(p.square, chess.NoSquare, p.game); m != nil {
// 			moveStart = p.square
// 			//flipGrid(p.container, p.game.Position().Board())
// 		} else {
// 			showCustomModalPopup(w, "Invalid move", fmt.Sprintf("Cannot move piece %d",
// 				p.game.Position().Board().Piece(p.square)))
// 		}
// 		return
// 	}

// 	// Check if the move from `moveStart` to the current square is valid
// 	if m := isValidMove(moveStart, p.square, p.game); m != nil {
// 		// Perform the original move
// 		//flipGrid(p.container, p.game.Position().Board())
// 		// flipGrid()
// 		move(m, p.game, g, over)
// 		//	flipGrid(p.container, p.game.Position().Board())
// 		// Now mirror the move to the opposite side of the board
// 		// mirroredStart := mirrorSquare(moveStart)
// 		// mirroredEnd := mirrorSquare(p.square)

// 		// fmt.Printf("Original Move: %v -> %v\n", moveStart, p.square)
// 		// fmt.Printf("Mirrored Move: %v -> %v\n", mirroredStart, mirroredEnd)

// 		// // Check if the mirrored move is valid
// 		// if mirroredMove := isValidMove(mirroredStart, mirroredEnd, p.game); mirroredMove != nil {
// 		// 	// Perform the mirrored move
// 		// 	move(mirroredMove, p.game, g, over)
// 		// }

// 		// Reset the move start
// 		moveStart = chess.NoSquare

// 		// Optionally flip the grid (if needed)
// 		// flipGrid(p.container, p.game.Position().Board())

// 		return
// 	}

// 	// If the move is invalid, show an error message
// 	showCustomModalPopup(w, "Invalid move", fmt.Sprintf("Cannot move piece %d to square %v",
// 		p.game.Position().Board().Piece(moveStart), p.square))

// 	// Reset the move start
// 	moveStart = chess.NoSquare
// }

// func (p *peice) Tapped(ev *fyne.PointEvent) {
// 	w := p.window
// 	g := p.container

// 	// Check if a move is in progress
// 	if moveStart == chess.NoSquare {
// 		// Check if the move is valid from the current square
// 		if m := isValidMove(p.square, chess.NoSquare, p.game); m != nil {
// 			moveStart = p.square
// 		} else {
// 			showCustomModalPopup(w, "Invalid move", fmt.Sprintf("Cannot move piece %d",
// 				p.game.Position().Board().Piece(p.square)))
// 		}
// 		return
// 	}

// 	// Check if the move from `moveStart` to the current square is valid
// 	if m := isValidMove(moveStart, p.square, p.game); m != nil {
// 		moveStart = chess.NoSquare
// 		move(m, p.game, g, over)

// 		// Mirror the move to the opposite side of the board
// 		mirroredStart := mirrorSquare(moveStart)
// 		mirroredEnd := mirrorSquare(p.square)

// 		// Create a mirrored move and apply it
// 		if mirroredMove := isValidMove(mirroredStart, mirroredEnd, p.game); mirroredMove != nil {
// 			move(mirroredMove, p.game, g, over)
// 		}

// 		flipGrid(p.container, p.game.Position().Board())
// 		return
// 	}

// 	showCustomModalPopup(w, "Invalid move", fmt.Sprintf("Cannot move piece %d to square %v",
// 		p.game.Position().Board().Piece(moveStart), p.square))

// 	moveStart = chess.NoSquare
// }

//	func randomResponse(game *chess.Game) {
//		rand.Seed(time.Now().Unix())
//		valid := game.ValidMoves()
//		m := valid[rand.Intn(len(valid))]
//		move(m, game, grid, over)
//	}
func isValidMove(s1, s2 chess.Square, g *chess.Game) *chess.Move {
	valid := g.ValidMoves()
	// fmt.Println(valid)
	for _, m := range valid {
		// fmt.Println(m.S1() == s1)
		// fmt.Println(m.S1())
		// fmt.Println(s1)
		if m.S1() == s1 && (s2 == chess.NoSquare || m.S2() == s2) {
			// fmt.Println("b")
			return m
		}
	}

	return nil
}

func randomResponse(game *chess.Game, g *fyne.Container) {
	valid := game.ValidMoves()
	m := valid[rand.Intn(len(valid))]
	// fmt.Println("b")
	move(m, game, g, over)
}
