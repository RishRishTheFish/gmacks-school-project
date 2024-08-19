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

type peice struct {
	widget.Icon

	window    fyne.Window
	container *fyne.Container
	game      *chess.Game
	square    chess.Square
	image     *canvas.Image
	theme     *CustomTheme
	pve       bool
}

func newPeice(g *chess.Game, square chess.Square, grid *fyne.Container, pve bool, theme *CustomTheme) *peice {
	p := g.Position().Board().Piece(square)
	pt := p.Type()
	pc := p.Color()
	resource := resourceForPiece(pc, pt)

	if resource == nil {

	}

	img := canvas.NewImageFromResource(resource)
	img.FillMode = canvas.ImageFillContain

	ret := &peice{
		game:      g,
		square:    square,
		image:     img,
		window:    win,
		container: grid,
		theme:     theme,
		pve:       pve,
	}

	ret.ExtendBaseWidget(ret)
	return ret
}

func showCustomModalPopup(w fyne.Window, title, message string) {
	bgColor := color.RGBA{R: 0x00, G: 0x00, B: 0xFF, A: 0xFF}
	mainRect := canvas.NewRectangle(bgColor)

	popupSize := fyne.NewSize(w.Canvas().Size().Width, w.Canvas().Size().Height)

	bg := canvas.NewRectangle(bgColor)
	bg.Resize(popupSize)

	msgLabel := canvas.NewText(message, color.White)
	msgLabel.TextStyle = fyne.TextStyle{Bold: true}
	msgLabel.TextSize = 16

	var overlay fyne.CanvasObject
	// toggleButton3 := container.NewMax(widget.NewButton("", func() {
	// 	options.Hide()
	// 	createChess(w, CustomTheme)
	// }), container.NewMax(buttonColor, returnText("chess")))
	content := container.NewMax(
		mainRect,
		container.NewVBox(
			container.NewMax(bg, msgLabel),
			container.NewHBox(
				layout.NewSpacer(),
				// container.NewMax(
				// 	mainRect,
				container.NewMax(widget.NewButton("", func() {
					w.Canvas().Overlays().Remove(overlay)
				}), container.NewMax(buttonColor, returnText("Return to game"))),
				//),
				layout.NewSpacer(),
			),
		))

	overlay = container.NewCenter(content)

	// overlay = container.NewMax(centeredContent)
	overlay.Resize(popupSize)

	w.Canvas().Overlays().Add(overlay)
}

// container.NewHBox(
// 	layout.NewSpacer(),
// 	// container.NewMax(
// 	// 	mainRect,
// 	container.NewMax(widget.NewButton("", func() {
// 		w.Canvas().Overlays().Remove(overlay)
// 	}), container.NewMax(buttonColor, returnText("OK"))),
// 	//),
// 	layout.NewSpacer(),
// ),
// func mirrorSquare(square chess.Square) chess.Square {
// 	file := square.File()
// 	rank := square.Rank()

// 	mirroredFile := chess.File(7 - int(file))
// 	mirroredRank := chess.Rank(7 - int(rank))

//		return chess.NewSquare(mirroredFile, mirroredRank)
//	}
func (p *peice) Tapped(ev *fyne.PointEvent) {
	w := p.window
	g := p.container

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
		move(m, p.game, g, over, p)

		go func() {
			time.Sleep(time.Second)

			if !gameEnded {
				if p.pve {
					randomResponse(p.game, g, over, p.theme)
					//w.SetContent(createFlippedGrid(p.game, p.theme))
				} else {
					w.SetContent(createFlippedGrid(p.game, p.theme))
					// randomResponse(p.game, g, over, p.theme)
				}
			}
		}()
		return
	}
	showCustomModalPopup(w, "Invalid move", fmt.Sprintf("Cannot move piece %d to square %v",
		p.game.Position().Board().Piece(moveStart), p.square))

	moveStart = chess.NoSquare
}

func isValidMove(s1, s2 chess.Square, g *chess.Game) *chess.Move {
	valid := g.ValidMoves()

	for _, m := range valid {

		if m.S1() == s1 && (s2 == chess.NoSquare || m.S2() == s2) {

			return m
		}
	}

	return nil
}

func randomResponse(game *chess.Game, g *fyne.Container, over *canvas.Image, theme *CustomTheme) {
	valid := game.ValidMoves()
	m := valid[rand.Intn(len(valid))]
	piece := newPeice(game, m.S1(), g, true, theme)
	// Get the piece from the starting square of the move
	// piece := game.Position().Board().Piece(m.S1())

	// Here you could add logic to work with the piece in Fyne
	// fmt.Printf("Moving piece: %s from %s to %s\n", piece, m.S1(), m.S2())

	// Move the piece visually in Fyne
	move(m, game, g, over, piece)
}
