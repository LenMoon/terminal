package widget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var (
	// lightPalette is the xterm palette, for text on a light background. The slots
	// that fell short of a 4.5:1 contrast ratio on white are darkened.
	lightPalette = [16]color.Color{
		color.Black,
		&color.RGBA{R: 0xaa, G: 0x00, B: 0x00, A: 0xff},
		&color.RGBA{R: 0x00, G: 0x80, B: 0x00, A: 0xff},
		&color.RGBA{R: 0x70, G: 0x70, B: 0x00, A: 0xff},
		&color.RGBA{R: 0x00, G: 0x00, B: 0xaa, A: 0xff},
		&color.RGBA{R: 0xaa, G: 0x00, B: 0xaa, A: 0xff},
		&color.RGBA{R: 0x00, G: 0x70, B: 0x70, A: 0xff},
		&color.RGBA{R: 0x70, G: 0x70, B: 0x70, A: 0xff},
		&color.RGBA{R: 0x55, G: 0x55, B: 0x55, A: 0xff},
		&color.RGBA{R: 0xcc, G: 0x00, B: 0x00, A: 0xff},
		&color.RGBA{R: 0x00, G: 0x64, B: 0x00, A: 0xff},
		&color.RGBA{R: 0x5c, G: 0x5c, B: 0x00, A: 0xff},
		&color.RGBA{R: 0x00, G: 0x00, B: 0xcc, A: 0xff},
		&color.RGBA{R: 0xc4, G: 0x00, B: 0xc4, A: 0xff},
		&color.RGBA{R: 0x00, G: 0x5c, B: 0x5c, A: 0xff},
		&color.RGBA{R: 0x58, G: 0x58, B: 0x58, A: 0xff},
	}
	// darkPalette is the same slots for text on a dark background.
	darkPalette = [16]color.Color{
		color.Black,
		&color.RGBA{R: 0xff, G: 0x5c, B: 0x57, A: 0xff},
		&color.RGBA{R: 0x3f, G: 0xb9, B: 0x50, A: 0xff},
		&color.RGBA{R: 0xc0, G: 0xc0, B: 0x00, A: 0xff},
		&color.RGBA{R: 0x7b, G: 0x7b, B: 0xff, A: 0xff},
		&color.RGBA{R: 0xd6, G: 0x70, B: 0xd6, A: 0xff},
		&color.RGBA{R: 0x29, G: 0xb8, B: 0xdb, A: 0xff},
		&color.RGBA{R: 0xd0, G: 0xd0, B: 0xd0, A: 0xff},
		&color.RGBA{R: 0x8a, G: 0x8f, B: 0x98, A: 0xff},
		&color.RGBA{R: 0xff, G: 0x8a, B: 0x80, A: 0xff},
		&color.RGBA{R: 0x5a, G: 0xf7, B: 0x8e, A: 0xff},
		&color.RGBA{R: 0xff, G: 0xff, B: 0x55, A: 0xff},
		&color.RGBA{R: 0xa5, G: 0xb6, B: 0xff, A: 0xff},
		&color.RGBA{R: 0xff, G: 0x92, B: 0xdf, A: 0xff},
		&color.RGBA{R: 0x9a, G: 0xed, B: 0xfe, A: 0xff},
		&color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
	}
)

// ThemeVariant reports the theme variant that palette colors have to match. It is
// a variable so that tests can pin the variant.
var ThemeVariant = systemThemeVariant

func systemThemeVariant() fyne.ThemeVariant {
	app := fyne.CurrentApp()
	if app == nil || app.Settings().ThemeVariant() != theme.VariantDark {
		return theme.VariantLight
	}
	return theme.VariantDark
}

// PaletteColor is one of the 16 ANSI colors, tied to the theme variant it was
// looked up under. Being a comparable value, Fyne's text texture cache can tell
// the two variants of a slot apart.
type PaletteColor struct {
	Slot    uint8
	Variant fyne.ThemeVariant
}

// RGBA implements color.Color.
func (c PaletteColor) RGBA() (uint32, uint32, uint32, uint32) {
	pal := lightPalette
	if c.Variant == theme.VariantDark {
		pal = darkPalette
	}
	return pal[c.Slot].RGBA()
}

// ANSI returns the color of ANSI slot (0–15) for the active theme variant.
func ANSI(slot uint8) PaletteColor {
	return PaletteColor{Slot: slot, Variant: ThemeVariant()}
}

// Resolve returns c for the active theme variant: a palette color is looked up
// again, any other color is returned unchanged.
func Resolve(c color.Color) color.Color {
	if p, ok := c.(PaletteColor); ok {
		return ANSI(p.Slot)
	}
	return c
}
