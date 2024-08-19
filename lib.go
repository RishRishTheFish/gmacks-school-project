package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

const (
	ColorNameBackground = "background"
	ColorNameButton     = "button"
	ColorNameText       = "text"
	ColorNameSidebar    = "sidebar"
	ColorNameTop        = "top"
	ColorNameContent    = "content"
	ColorNameLabel      = "label"
)

func colorToRGBA(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(a >> 8),
	}
}

type CustomTheme struct {
	backgroundColor color.Color
	sideBarColor    color.Color
	topColor        color.Color
	contentColor    color.Color
	themeLabelColor color.Color
	foregroundColor color.Color
	optionsColor    color.Color
	buttonColor     color.Color
}

func NewCustomTheme(
	bg color.Color,
	sideBar color.Color,
	top color.Color,
	content color.Color,
	themeLabel color.Color,
	options color.Color,
	button color.Color,
	foreground color.Color,
) *CustomTheme {
	if bg == nil {
		bg = color.Black
	}
	if sideBar == nil {
		sideBar = color.Gray{Y: 0x80}
	}
	if top == nil {
		top = color.Gray{Y: 0x80}
	}
	if content == nil {
		content = color.Gray{Y: 0x80}
	}
	if themeLabel == nil {
		themeLabel = color.Gray{Y: 0x88}
	}
	if options == nil {
		options = color.Gray{Y: 0x80}
	}
	if button == nil {
		button = color.Gray{Y: 0x80}
	}
	if foreground == nil {
		foreground = color.White
	}

	return &CustomTheme{
		backgroundColor: bg,
		sideBarColor:    sideBar,
		topColor:        top,
		contentColor:    content,
		themeLabelColor: themeLabel,
		optionsColor:    options,
		buttonColor:     button,
		foregroundColor: foreground,
	}
}

func (c *CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return c.backgroundColor

}

func (c *CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (c *CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (c *CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
