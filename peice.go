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
	// *canvas.Image
	widget.Icon

	game   *chess.Game
	square chess.Square
	image  *canvas.Image
}

func newPeice(g *chess.Game, square chess.Square) *peice {
	p := g.Position().Board().Piece(square)
	pt := p.Type()
	pc := p.Color()
	fmt.Println(pt, pc)

	resource := resourceForPiece(pc, pt)
	// fmt.Println(resource)
	if resource == nil {
		fmt.Println("Warning 2: Resource for piece is nil")
	} else {
		// fmt.Println("A")
	}
	img := canvas.NewImageFromResource(resource)
	img.FillMode = canvas.ImageFillContain
	return &peice{
		game:   g,
		square: square,
		image:  img,
	}
}

//	func newPeice(g *chess.Game, square chess.Square) *peice {
//		p := g.Position().Board().Piece(square)
//		img := canvas.NewImageFromResource(resourceForPiece(p))
//		img.FillMode = canvas.ImageFillContain
//		return &peice{game: g, square: square}
//	}
func (p *peice) Tapped(ev *fyne.PointEvent) {
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
