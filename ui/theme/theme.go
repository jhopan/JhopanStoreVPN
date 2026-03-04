package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// LightTheme is a clean white/light theme for JhopanStoreVPN.
type DarkTheme struct{}

var _ fyne.Theme = (*DarkTheme)(nil)

func (d *DarkTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 250, G: 250, B: 252, A: 255}
	case theme.ColorNameButton:
		return color.NRGBA{R: 230, G: 232, B: 238, A: 255}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 30, G: 30, B: 35, A: 255}
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 0, G: 120, B: 215, A: 255}
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 160, G: 160, B: 170, A: 255}
	case theme.ColorNameDisabled:
		return color.NRGBA{R: 180, G: 180, B: 185, A: 255}
	case theme.ColorNameOverlayBackground:
		return color.NRGBA{R: 245, G: 245, B: 248, A: 250}
	case theme.ColorNameSeparator:
		return color.NRGBA{R: 210, G: 212, B: 218, A: 255}
	case theme.ColorNameInputBorder:
		return color.NRGBA{R: 190, G: 192, B: 200, A: 255}
	case theme.ColorNameShadow:
		return color.NRGBA{R: 0, G: 0, B: 0, A: 20}
	}
	return theme.DefaultTheme().Color(name, theme.VariantLight)
}

func (d *DarkTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (d *DarkTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (d *DarkTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 13
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameInnerPadding:
		return 6
	}
	return theme.DefaultTheme().Size(name)
}
