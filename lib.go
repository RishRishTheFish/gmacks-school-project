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

// type CustomTheme struct {
// 	backgroundColor color.Color
// 	sideBarColor    color.Color
// 	topColor        color.Color
// 	contentColor    color.Color
// 	themeLabelColor color.Color
// }

// NewCustomTheme creates a new CustomTheme with the specified background color
// If the background color is nil, it defaults to black
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
		themeLabel = color.Gray{Y: 0x88} // Example default gray color
	}
	if options == nil {
		options = color.Gray{Y: 0x80} // Example default gray color
	}
	if button == nil {
		button = color.Gray{Y: 0x80} // Example default gray color
	}
	if foreground == nil {
		foreground = color.White // Example default foreground color
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
	// switch name {
	// case ColorNameBackground:
	// 	return c.backgroundColor
	// case ColorNameButton:
	// 	return c.contentColor
	// case ColorNameText:
	// 	return c.themeLabelColor
	// case ColorNameSidebar:
	// 	return c.sideBarColor
	// case ColorNameTop:
	// 	return c.topColor
	// case ColorNameContent:
	// 	return c.contentColor
	// default:
	// 	return c.foregroundColor
	// }
}

// Font returns the font for a specific element of the theme
func (c *CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Icon returns the icon for a specific element of the theme
func (c *CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size returns the size for a specific element of the theme
func (c *CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
