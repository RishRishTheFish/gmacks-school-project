package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"github.com/notnil/chess"
)

func resourceForPiece(c chess.Color, p chess.PieceType) fyne.Resource {
	// Get the single-letter string representation of the piece type
	pieceType := p.String()
	peiceColor := c.String()
	// fmt.Println(peiceColor)
	switch peiceColor {
	case "b":
		switch pieceType {
		case "p":
			return resourceBlackPawnSvg
		case "r":
			return resourceBlackRookSvg
		case "n":
			return resourceBlackKnightSvg
		case "b":
			return resourceBlackBishopSvg
		case "q":
			return resourceBlackQueenSvg
		case "k":
			return resourceBlackKingSvg
		}
	case "w":
		fmt.Println("white")
		switch pieceType {
		case "p":
			return resourceWhitePawnSvg
		case "r":
			return resourceWhiteRookSvg
		case "n":
			return resourceWhiteKnightSvg
		case "b":
			return resourceWhiteBishopSvg
		case "q":
			return resourceWhiteQueenSvg
		case "k":
			return resourceWhiteKingSvg
		}
	}

	return nil
}
