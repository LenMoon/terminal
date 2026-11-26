package widget

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/stretchr/testify/assert"
)

// Fyne's theme.ColorNameBackground per theme variant.
const (
	lightBackground = "#ffffff"
	darkBackground  = "#171718"
)

// paletteWant pins every slot of both palettes; changing a color means changing this too.
var paletteWant = []struct {
	slot  uint8
	light string
	dark  string
}{
	{0, "#000000", "#000000"},
	{1, "#aa0000", "#ff5c57"},
	{2, "#008000", "#3fb950"},
	{3, "#707000", "#c0c000"},
	{4, "#0000aa", "#7b7bff"},
	{5, "#aa00aa", "#d670d6"},
	{6, "#007070", "#29b8db"},
	{7, "#707070", "#d0d0d0"},
	{8, "#555555", "#8a8f98"},
	{9, "#cc0000", "#ff8a80"},
	{10, "#006400", "#5af78e"},
	{11, "#5c5c00", "#ffff55"},
	{12, "#0000cc", "#a5b6ff"},
	{13, "#c400c4", "#ff92df"},
	{14, "#005c5c", "#9aedfe"},
	{15, "#585858", "#ffffff"},
}

func TestPalette(t *testing.T) {
	for _, want := range paletteWant {
		t.Run(fmt.Sprint(want.slot), func(t *testing.T) {
			pinVariant(t, theme.VariantLight)
			assert.Equal(t, rgba(hexColor(want.light)), rgba(ANSI(want.slot)))
			pinVariant(t, theme.VariantDark)
			assert.Equal(t, rgba(hexColor(want.dark)), rgba(ANSI(want.slot)))
		})
	}
}

// Any slot used as foreground has to stay readable on its own background.
func TestPaletteContrast(t *testing.T) {
	backgrounds := []struct {
		name       string
		variant    fyne.ThemeVariant
		background string
	}{{"light", theme.VariantLight, lightBackground}, {"dark", theme.VariantDark, darkBackground}}
	for _, tc := range backgrounds {
		t.Run(tc.name, func(t *testing.T) {
			pinVariant(t, tc.variant)
			for slot := uint8(1); slot < 16; slot++ {
				assert.GreaterOrEqual(t, contrast(ANSI(slot), hexColor(tc.background)), 4.5, "slot %d", slot)
			}
		})
	}
}

func TestSystemThemeVariant(t *testing.T) {
	prev := fyne.CurrentApp()
	t.Cleanup(func() { fyne.SetCurrentApp(prev) })

	fyne.SetCurrentApp(nil) // no app at all
	assert.Equal(t, theme.VariantLight, systemThemeVariant())
	fyne.SetCurrentApp(stubApp{variant: theme.VariantDark})
	assert.Equal(t, theme.VariantDark, systemThemeVariant())
	fyne.SetCurrentApp(stubApp{variant: fyne.ThemeVariant(2)}) // what the Fyne test driver reports
	assert.Equal(t, theme.VariantLight, systemThemeVariant())
}

// pinVariant keeps palette lookups on one theme variant for the rest of the test.
func pinVariant(t *testing.T, v fyne.ThemeVariant) {
	t.Helper()
	prev := ThemeVariant
	ThemeVariant = func() fyne.ThemeVariant { return v }
	t.Cleanup(func() { ThemeVariant = prev })
}

// hexColor parses "#rrggbb" into an opaque color.
func hexColor(hex string) color.RGBA {
	v, err := strconv.ParseUint(strings.TrimPrefix(hex, "#"), 16, 32)
	if err != nil {
		panic(err)
	}
	return color.RGBA{R: uint8(v >> 16), G: uint8(v >> 8), B: uint8(v), A: 0xff}
}

// rgba is the 16-bit premultiplied channel tuple of a color.
func rgba(c color.Color) [4]uint32 {
	r, g, b, a := c.RGBA()
	return [4]uint32{r, g, b, a}
}

// contrast is the WCAG contrast ratio of fg on bg.
func contrast(fg, bg color.Color) float64 {
	a, b := luminance(fg), luminance(bg)
	if a < b {
		a, b = b, a
	}
	return (a + 0.05) / (b + 0.05)
}

// luminance is the WCAG relative luminance of a color.
func luminance(c color.Color) float64 {
	r, g, b, _ := c.RGBA()
	channel := func(v float64) float64 {
		v /= 0xffff
		if v <= 0.03928 {
			return v / 12.92
		}
		return math.Pow((v+0.055)/1.055, 2.4)
	}
	return 0.2126*channel(float64(r)) + 0.7152*channel(float64(g)) + 0.0722*channel(float64(b))
}

type stubApp struct {
	fyne.App
	variant fyne.ThemeVariant
}

func (a stubApp) Settings() fyne.Settings { return stubSettings{variant: a.variant} }

type stubSettings struct {
	fyne.Settings
	variant fyne.ThemeVariant
}

func (s stubSettings) ThemeVariant() fyne.ThemeVariant { return s.variant }
