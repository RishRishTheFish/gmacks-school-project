package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/notnil/chess"
)

var moveStart chess.Square = chess.NoSquare

type peice struct {
	widget.Icon
	game   *chess.Game
	square chess.Square
	image  *canvas.Image
}

func newPeice(g *chess.Game, square chess.Square) *peice {
	p := g.Position().Board().Piece(square)
	pt := p.Type()
	pc := p.Color()
	resource := resourceForPiece(pc, pt)

	if resource == nil {
		fmt.Println("Warning: Resource for piece is nil")
	}

	img := canvas.NewImageFromResource(resource)
	img.FillMode = canvas.ImageFillContain

	ret := &peice{
		game:   g,
		square: square,
		image:  img,
	}

	ret.ExtendBaseWidget(ret)
	return ret
}

func (p *peice) Tapped(ev *fyne.PointEvent) {
	fmt.Println("Tapped on square:", p.square)
	if moveStart == chess.NoSquare {
		moveStart = p.square
		return
	}

	valid := p.game.ValidMoves()
	for _, m := range valid {
		if m.S1() == moveStart && m.S2() == p.square {
			move(m, p.game, grid, over)
			return
		}
	}

	pos := p.game.Position().Board().Piece(p.square)
	dialog.ShowInformation("Invalid move", "Cannot move piece "+pos.String()+" to square "+p.square.String(), win)
}
