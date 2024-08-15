package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/notnil/chess"
)

var moveStart chess.Square = chess.NoSquare

type peice struct {
	// *canvas.Image
	widget.Icon

	square chess.Square
}

func newPeice(board *chess.Board, square chess.Square) *peice {
	p := board.Piece(square)
	img := canvas.NewImageFromResource(resourceForPiece(p))
	img.FillMode = canvas.ImageFillContain
	return &peice{square: square}
}
func (p *peice) Tapped(ev *fyne.PointEvent) {
	//log.Println("tapped sq", p.square)
	// b.startMove()
	if moveStart == chess.NoSquare {
		moveStart = p.square
		return
	}
}
