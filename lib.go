package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type CustomTheme struct {
	backgroundColor color.Color
	foregroundColor color.Color
}

// NewCustomTheme creates a new CustomTheme with the specified background color
// If the background color is nil, it defaults to black

// NewCustomTheme creates a new theme with specified background and foreground colors.
func NewCustomTheme(bg *color.RGBA, fg color.Color) *CustomTheme {
	// Define default colors as color.RGBA
	var defaultBackground color.RGBA = color.RGBA{R: 0, G: 0, B: 0, A: 255}       // Black
	var defaultForeground color.RGBA = color.RGBA{R: 255, G: 255, B: 255, A: 255} // White

	// If bg is nil, use default background color
	var background color.RGBA
	if bg == nil {
		background = defaultBackground
	} else {
		background = *bg
	}

	// If fg is nil, use default foreground color
	var foreground color.RGBA
	if fg == nil {
		foreground = defaultForeground
	} else {
		// Convert color.Color to color.RGBA
		rgba, ok := fg.(color.RGBA)
		if !ok {
			// If fg is not of type color.RGBA, default to white
			rgba = defaultForeground
		}
		foreground = rgba
	}

	return &CustomTheme{
		backgroundColor: background,
		foregroundColor: foreground,
	}
}

// Color returns the color for a specific element of the theme
func (c *CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return c.backgroundColor
	case theme.ColorNameForeground:
		return c.foregroundColor
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
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
