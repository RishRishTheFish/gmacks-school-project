package main

import (
	"image/color"

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

			img := canvas.NewImageFromResource(resourceForPiece())
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
	w.SetContent(grid)
	// w.ShowAndRun()
}
