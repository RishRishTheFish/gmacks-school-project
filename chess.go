package main

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	//"fyne.io/fyne/v2/internal/widget"
	"github.com/notnil/chess"
)

var (
	grid *fyne.Container
	over *canvas.Image
	win  fyne.Window // Changed from *fyne.Window to fyne.Window
)

// func createGrid(game *chess.Game) *fyne.Container {
// 	var cells []fyne.CanvasObject

// 	for y := 0; y < 8; y++ {
// 		for x := 0; x < 8; x++ {
// 			bg := canvas.NewRectangle(color.Gray{0x30})
// 			if x%2 == y%2 {
// 				bg.FillColor = color.Gray{0xE0}
// 			}

//				piece := newPeice(game, chess.Square(x+y*8))
//				cells = append(cells, container.NewMax(bg, piece.image))
//			}
//		}
//		return container.New(&boardLayout{}, cells...)
//	}
func createGrid(g *chess.Game) *fyne.Container {
	var cells []fyne.CanvasObject

	// Create a placeholder for the grid pointer
	grid := new(fyne.Container)

	for y := 7; y >= 0; y-- {
		for x := 0; x < 8; x++ {
			bg := canvas.NewRectangle(color.NRGBA{0xF4, 0xE2, 0xB6, 0xFF})
			if x%2 == y%2 {
				bg.FillColor = color.RGBA{0x73, 0x50, 0x32, 0xFF}
			}

			// Create the piece and pass the grid pointer to it
			p := newPeice(g, chess.Square(x+y*8), grid)
			cells = append(cells, container.NewMax(bg, p))
		}
	}

	// Now that the grid is fully constructed, reassign grid
	*grid = *container.New(&boardLayout{}, cells...)

	return grid
}

var gameEnded bool

// func prepareToMove(m *chess.Move, game *chess.Game, grid *fyne.Container, over *canvas.Image) {

// }
// func move(m *chess.Move, game *chess.Game, grid *fyne.Container, over *canvas.Image) {
// 	// Mirror the start and end squares
// 	mirroredStart := mirrorSquare(m.S1())
// 	mirroredEnd := mirrorSquare(m.S2())

// 	// Safeguard checks
// 	off := squareToOffset(mirroredStart)
// 	var cell *fyne.Container

// 	if grid == nil {
// 		fmt.Println("Warning: Grid is nil, cannot get cell")
// 		return
// 	}
// 	if grid.Objects == nil {
// 		fmt.Println("Warning: Grid objects are nil, cannot get cell")
// 		return
// 	}
// 	if off < 0 || off >= len(grid.Objects) {
// 		fmt.Println("Warning: Offset is out of bounds, cannot get cell")
// 		return
// 	}

// 	// If all checks pass, proceed to get the cell
// 	cell = grid.Objects[off].(*fyne.Container)
// 	img := cell.Objects[1].(*peice)
// 	pos1 := cell.Position()

// 	// Handle overlay
// 	over.Resource = img.Resource
// 	over.Move(pos1)
// 	over.Resize(img.Size())
// 	over.Refresh() // clear old resource before showing

// 	over.Show()
// 	img.Resource = nil
// 	img.Refresh()

// 	// Move to the destination cell
// 	off = squareToOffset(mirroredEnd)
// 	if off < 0 || off >= len(grid.Objects) {
// 		fmt.Println("Warning: Offset is out of bounds for target square")
// 		return
// 	}
// 	cell = grid.Objects[off].(*fyne.Container)
// 	pos2 := cell.Position()

// 	// Animate the move
// 	a := canvas.NewPositionAnimation(pos1, pos2, time.Millisecond*500, func(p fyne.Position) {
// 		over.Move(p)
// 		over.Refresh()
// 	})
// 	a.Start()
// 	time.Sleep(time.Millisecond * 550)

// 	// Perform the mirrored move on the game board
// 	mirroredMove := chess.Move{
// 		// Assuming Move has methods to create or set moves
// 		// You may need to use appropriate methods to apply mirrored move
// 		S1: mirroredStart,
// 		S2: mirroredEnd,
// 	}
// 	game.Move(&mirroredMove)
// 	refreshGrid(grid, game.Position().Board())
// 	over.Hide()

// 	// Check for game outcome
// 	if game.Outcome() != chess.NoOutcome {
// 		result := "draw"
// 		switch game.Outcome().String() {
// 		case "1-0":
// 			result = "won"
// 		case "0-1":
// 			result = "lost"
// 		}
// 		dialog.ShowInformation("Game ended",
// 			"Game "+result+" because "+game.Method().String(), win)
// 	}
// }

func move(m *chess.Move, game *chess.Game, grid *fyne.Container, over *canvas.Image) {
	off := squareToOffset(m.S1())

	// Attempt to get the cell, even if there might be issues
	var cell *fyne.Container
	if grid == nil {
		fmt.Println("Warning: Grid is nil, cannot get cell")
		return
	}
	if grid.Objects == nil {
		fmt.Println("Warning: Grid objects are nil, cannot get cell")
		return
	}
	if off < 0 || off >= len(grid.Objects) {
		fmt.Println("Warning: Offset is out of bounds, cannot get cell")
		return
	}

	// If all checks pass, proceed to get the cell
	cell = grid.Objects[off].(*fyne.Container)
	img := cell.Objects[1].(*peice)
	pos1 := cell.Position()

	over.Resource = img.Resource
	over.Move(pos1)
	over.Resize(img.Size())
	over.Refresh() // clear old resource before showing

	over.Show()
	img.Resource = nil
	img.Refresh()

	off = squareToOffset(m.S2())
	if off < 0 || off >= len(grid.Objects) {
		fmt.Println("Warning: Offset is out of bounds for target square")
		return
	}
	cell = grid.Objects[off].(*fyne.Container)
	pos2 := cell.Position()

	a := canvas.NewPositionAnimation(pos1, pos2, time.Millisecond*500, func(p fyne.Position) {
		over.Move(p)
		over.Refresh()
	})
	a.Start()
	time.Sleep(time.Millisecond * 550)

	game.Move(m)
	refreshGrid(grid, game.Position().Board())
	over.Hide()

	if game.Outcome() != chess.NoOutcome {
		result := "draw"
		switch game.Outcome().String() {
		case "1-0":
			result = "won"
		case "0-1":
			result = "lost"
		}
		showCustomModalPopup(win, "Game ended", "Game "+result+" because "+game.Method().String())
		gameEnded = true
		// dialog.ShowInformation("Game ended",
		// 	"Game "+result+" because "+game.Method().String(), win)
	}
}

// func move(m *chess.Move, game *chess.Game, grid *fyne.Container, over *canvas.Image) {
// 	// m := mirrorSquare(b.S1())
// 	//flipGrid(grid, game.Position().Board())
// 	off := squareToOffset(mirrorSquare(m.S1()))

// 	// Attempt to get the cell, even if there might be issues
// 	var cell *fyne.Container
// 	if grid == nil {
// 		fmt.Println("Warning: Grid is nil, cannot get cell")
// 		return
// 	}
// 	if grid.Objects == nil {
// 		fmt.Println("Warning: Grid objects are nil, cannot get cell")
// 		return
// 	}
// 	if off < 0 || off >= len(grid.Objects) {
// 		fmt.Println("Warning: Offset is out of bounds, cannot get cell")
// 		return
// 	}

// 	// If all checks pass, proceed to get the cell
// 	cell = grid.Objects[off].(*fyne.Container)
// 	img := cell.Objects[1].(*peice)
// 	pos1 := cell.Position()

// 	over.Resource = img.Resource
// 	over.Move(pos1)
// 	over.Resize(img.Size())
// 	over.Refresh() // clear old resource before showing

// 	over.Show()
// 	img.Resource = nil
// 	img.Refresh()

// 	off = squareToOffset(mirrorSquare(m.S2()))
// 	if off < 0 || off >= len(grid.Objects) {
// 		fmt.Println("Warning: Offset is out of bounds for target square")
// 		return
// 	}
// 	cell = grid.Objects[off].(*fyne.Container)
// 	pos2 := cell.Position()

// 	a := canvas.NewPositionAnimation(pos1, pos2, time.Millisecond*500, func(p fyne.Position) {
// 		over.Move(p)
// 		over.Refresh()
// 	})
// 	a.Start()
// 	time.Sleep(time.Millisecond * 550)

// 	game.Move(m)
// 	flipGrid(grid, game.Position().Board())
// 	//refreshGrid(grid, game.Position().Board())
// 	over.Hide()

// 	if game.Outcome() != chess.NoOutcome {
// 		result := "draw"
// 		switch game.Outcome().String() {
// 		case "1-0":
// 			result = "won"
// 		case "0-1":
// 			result = "lost"
// 		}
// 		dialog.ShowInformation("Game ended",
// 			"Game "+result+" because "+game.Method().String(), win)
// 	}
// }

func createChess(w fyne.Window, customTheme *CustomTheme) {
	game := chess.NewGame()
	win = w // Assign the window to the global variable

	win.Resize(fyne.NewSize(480, 480))

	grid := createGrid(game)
	over = canvas.NewImageFromResource(nil)
	over.Hide()

	// Create a container with the chess widgets and place it in a specific area
	chessWidgets := getSidebar(w, customTheme)

	content := container.NewBorder(
		nil,          // No top content
		nil,          // No bottom content
		nil,          // No left content
		chessWidgets, // No right content
		container.NewMax(grid, container.NewWithoutLayout(over)), // Main content (center)
		// Additional content (could be placed in top, bottom, left, or right
	)

	win.SetContent(content)
	refreshGrid(grid, game.Position().Board())
}

func refreshGrid(grid *fyne.Container, b *chess.Board) {
	y, x := 7, 0
	for _, cell := range grid.Objects {
		p := b.Piece(chess.Square(x + y*8))

		img := cell.(*fyne.Container).Objects[1].(*peice)
		img.Resource = resourceForPiece(p.Color(), p.Type())
		img.Refresh()

		x++
		if x == 8 {
			x = 0
			y--
		}
	}
}
func flipGrid(grid *fyne.Container, b *chess.Board) {
	// Retrieve the objects in the grid
	objects := grid.Objects

	// Prepare a list to hold the pieces
	pieces := make([]*peice, len(objects))

	// Extract the piece widgets from the grid objects
	for i, obj := range objects {
		pieceContainer := obj.(*fyne.Container)
		piece := pieceContainer.Objects[1].(*peice)
		pieces[i] = piece
	}

	// Flip each row of the grid horizontally
	for y := 0; y < 8; y++ {
		for x := 0; x < 4; x++ { // Only need to go halfway because we'll swap pairs
			// Calculate the index for the current and opposite cells in the row
			index1 := y*8 + x
			index2 := y*8 + (7 - x)

			// Swap the pieces at index1 and index2
			piece1 := pieces[index1]
			piece2 := pieces[index2]

			// Swap the resources
			p1 := b.Piece(chess.Square(index1))
			p2 := b.Piece(chess.Square(index2))

			piece1.Resource = resourceForPiece(p2.Color(), p2.Type())
			piece2.Resource = resourceForPiece(p1.Color(), p1.Type())

			piece1.Refresh()
			piece2.Refresh()
		}
	}

	// Optionally, refresh the grid container to reflect changes
	grid.Refresh()
}

// func flipGrid(grid *fyne.Container, b *chess.Board) {
// 	// Retrieve the objects in the grid
// 	objects := grid.Objects

// 	// Prepare a list to hold the pieces
// 	pieces := make([]*peice, len(objects))

// 	// Extract the piece widgets from the grid objects
// 	for i, obj := range objects {
// 		pieceContainer := obj.(*fyne.Container)
// 		piece := pieceContainer.Objects[1].(*peice)
// 		pieces[i] = piece
// 	}

// 	// Reverse the pieces along the vertical axis
// 	y, x := 7, 0
// 	for i := range pieces {
// 		// Calculate the new index after flipping
// 		newX := 7 - x
// 		newY := y

// 		// Update the piece's position and resource
// 		//p := b.Piece(chess.Square(x + y*8))
// 		newPiece := b.Piece(chess.Square(newX + newY*8))

// 		pieces[i].Resource = resourceForPiece(newPiece.Color(), newPiece.Type())
// 		pieces[i].Refresh()

// 		// Update x and y for the next piece
// 		x++
// 		if x == 8 {
// 			x = 0
// 			y--
// 		}
// 	}

// 	// Optionally, refresh the grid container to reflect changes
// 	grid.Refresh()
// }

// returnToggleButton := container.NewMax(widget.NewButton("", func() {
// 	enableOptions := false
// 	resizeAndRefresh(top, bottom, left, right, center, textbox, content, dividers, root.Size(), right.Visible(), enableOptions, options, root)
// }), container.NewMax(buttonColor, returnText("exit options")))
// Use container.NewBorder to manage the layout more effectively

// func move(m *chess.Move, game *chess.Game, grid *fyne.Container, over *canvas.Image) {
// 	off := squareToOffset(m.S1())

// 	// Safeguard to ensure grid.Objects[off] is a *fyne.Container
// 	var cell *fyne.Container
// 	// Check if grid.Objects is nil or if the offset is out of bounds
// 	if grid == nil || grid.Objects == nil || off < 0 || off >= len(grid.Objects) {
// 		fmt.Println("Error: grid.Objects is nil or offset is out of bounds:", off)
// 		return
// 	}

// 	if obj := grid.Objects[off]; obj != nil {
// 		var ok bool
// 		cell, ok = obj.(*fyne.Container)
// 		if !ok {
// 			fmt.Println("Error: Object at offset is not a *fyne.Container")
// 			cell = fyne.NewContainer()
// 			grid.Objects[off] = cell
// 		}
// 	} else {
// 		fmt.Println("Warning: Object at offset is nil, creating new container")
// 		cell = fyne.NewContainer()
// 		grid.Objects[off] = cell
// 	}

// 	// Find the piece in the cell
// 	var movedPiece *peice
// 	for _, obj := range cell.Objects {
// 		if p, ok := obj.(*peice); ok {
// 			movedPiece = p
// 			break
// 		}
// 	}

// 	// If the piece is not found or doesn't have an image, create one
// 	if movedPiece == nil {
// 		fmt.Println("Error: piece not found, creating new piece")
// 		square := m.S1()
// 		boardPiece := game.Position().Board().Piece(square)

// 		// Get the resource for the piece
// 		resource := resourceForPiece(boardPiece.Color(), boardPiece.Type())
// 		if resource == nil {
// 			fmt.Println("Error: No resource found for piece type")
// 			return
// 		}

// 		// Create a new image
// 		img := canvas.NewImageFromResource(resource)
// 		img.FillMode = canvas.ImageFillContain

// 		// Create the piece and assign the image and container
// 		movedPiece = &peice{
// 			container: cell,
// 			image:     img,
// 			square:    square,
// 		}
// 		cell.Objects = append(cell.Objects, movedPiece)
// 	}

// 	if movedPiece.image == nil {
// 		fmt.Println("Error: piece image not found, creating new image")
// 		boardPiece := game.Position().Board().Piece(m.S1())

// 		// Get the resource for the piece
// 		resource := resourceForPiece(boardPiece.Color(), boardPiece.Type())
// 		if resource == nil {
// 			fmt.Println("Error: No resource found for piece type")
// 			return
// 		}

// 		// Create a new image
// 		img := canvas.NewImageFromResource(resource)
// 		img.FillMode = canvas.ImageFillContain
// 		movedPiece.image = img
// 		cell.Objects = append(cell.Objects, img)
// 	}

// 	pos1 := movedPiece.image.Position().Add(movedPiece.container.Position())
// 	over.Resource = movedPiece.image.Resource
// 	over.Resize(movedPiece.image.Size())
// 	over.Move(pos1)
// 	over.Refresh()
// 	over.Show()

// 	// Clear resource of the original image
// 	movedPiece.image.Resource = nil
// 	movedPiece.image.Refresh()

// 	off = squareToOffset(m.S2())

// 	// Safeguard for destination square
// 	if off < 0 || off >= len(grid.Objects) {
// 		fmt.Println("Error: Offset out of bounds:", off)
// 		return
// 	}

// 	if obj := grid.Objects[off]; obj != nil {
// 		var ok bool
// 		cell, ok = obj.(*fyne.Container)
// 		if !ok {
// 			fmt.Println("Error: Object at offset is not a *fyne.Container")
// 			cell = fyne.NewContainer()
// 			grid.Objects[off] = cell
// 		}
// 	} else {
// 		fmt.Println("Warning: Object at offset is nil, creating new container")
// 		cell = fyne.NewContainer()
// 		grid.Objects[off] = cell
// 	}

// 	pos2 := cell.Position().Add(grid.Position())

// 	a := canvas.NewPositionAnimation(pos1, pos2, time.Millisecond*500, func(p fyne.Position) {
// 		over.Move(p)
// 		over.Refresh()
// 	})
// 	a.Start()

// 	time.Sleep(time.Millisecond * 500)

// 	game.Move(m)

//		over.Hide()
//		refreshGrid(grid, game.Position().Board(), game)
//	}

// Uncomment the following code if you need automated game moves
// go func() {
// 	rand.Seed(time.Now().UnixNano())
// 	for game.Outcome() == chess.NoOutcome {
// 		time.Sleep(500 * time.Millisecond)
// 		valid := game.ValidMoves()
// 		m := valid[rand.Intn(len(valid))]
// 		move(m, game, grid, over)
// 	}
// }()

// func refreshGrid(grid *fyne.Container, b *chess.Board, game *chess.Game) {
// 	y, x := 7, 0

// 	for _, cell := range grid.Objects {
// 		if container, ok := cell.(*fyne.Container); ok {
// 			// Remove existing images from the container
// 			for i := len(container.Objects) - 1; i >= 0; i-- {
// 				if _, ok := container.Objects[i].(*canvas.Image); ok {
// 					container.Objects = append(container.Objects[:i], container.Objects[i+1:]...)
// 				}
// 			}

// 			// Get the square on the chessboard
// 			square := chess.Square(x + y*8)
// 			boardPiece := b.Piece(square)

// 			// Create a new piece widget for this square
// 			pieceWidget := newPeice(game, square)
// 			pieceWidget.container = container

// 			// If there is a chess piece on this square, create an image for it
// 			//p r n b q k
// 			pieceType := boardPiece.Type().String()
// 			if pieceType == "p" || pieceType == "r" || pieceType == "n" || pieceType == "b" || pieceType == "q" || pieceType == "k" {
// 				resource := resourceForPiece(boardPiece.Color(), boardPiece.Type())
// 				if resource != nil {
// 					img := canvas.NewImageFromResource(resource)
// 					img.FillMode = canvas.ImageFillContain
// 					container.Objects = append(container.Objects, img)

// 					// Assign the image to the piece widget
// 					pieceWidget.image = img
// 				}
// 			} else {
// 				// If there's no piece on this square, set the image to nil
// 				pieceWidget.image = nil
// 			}

// 			// Add the piece widget to the container
// 			container.Objects = append(container.Objects, pieceWidget)
// 			container.Refresh()
// 		}

// 		x++
// 		if x == 8 {
// 			x = 0
// 			y--
// 		}
// 	}
// }

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
