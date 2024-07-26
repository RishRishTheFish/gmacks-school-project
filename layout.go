package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

const sideWidth = 228

type mainLayout struct {
	top      fyne.CanvasObject
	left     fyne.CanvasObject
	right    fyne.CanvasObject
	content  fyne.CanvasObject
	dividers [3]fyne.CanvasObject
}

func newLayout(top, left, right, content fyne.CanvasObject, dividers [3]fyne.CanvasObject) fyne.Layout {
	return &mainLayout{top: top, left: left, right: right, content: content, dividers: dividers}
}

func (l *mainLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	topHeight := l.top.MinSize().Height
	l.top.Resize(fyne.NewSize(size.Width, topHeight))

	l.left.Move(fyne.NewPos(0, topHeight))
	l.left.Resize(fyne.NewSize(sideWidth, size.Height-topHeight))

	l.right.Move(fyne.NewPos(size.Width-sideWidth, topHeight))
	l.right.Resize(fyne.NewSize(sideWidth, size.Height-topHeight))

	l.content.Move(fyne.NewPos(sideWidth, topHeight))
	l.content.Resize(fyne.NewSize(size.Width-2*sideWidth, size.Height-topHeight))

	dividerThickness := theme.SeparatorThicknessSize()
	l.dividers[0].Move(fyne.NewPos(0, topHeight))
	l.dividers[0].Resize(fyne.NewSize(size.Width, dividerThickness))

	l.dividers[1].Move(fyne.NewPos(sideWidth, topHeight))
	l.dividers[1].Resize(fyne.NewSize(dividerThickness, size.Height-topHeight))

	l.dividers[2].Move(fyne.NewPos(size.Width-sideWidth, topHeight))
	l.dividers[2].Resize(fyne.NewSize(dividerThickness, size.Height-topHeight))
}

func (l *mainLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	// return fyne.NewSize(10, 10)

	borders := fyne.NewSize(
		sideWidth*2,
		l.top.MinSize().Height,
	)
	return borders.AddWidthHeight(100, 100)
}
