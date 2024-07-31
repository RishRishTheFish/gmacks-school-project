//go:generate ../bin/fyne bundle -o chess.go assets

package main

import (
	"fyne.io/fyne/v2"
)

func resourceForPiece() fyne.Resource {
	return resourceBlackBishopSvg
}
