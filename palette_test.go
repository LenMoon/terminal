package terminal

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	widget2 "github.com/fyne-io/terminal/internal/widget"
	"github.com/stretchr/testify/assert"
)

// pinVariant keeps palette lookups on one theme variant for the rest of the test.
func pinVariant(t *testing.T, v fyne.ThemeVariant) {
	t.Helper()
	prev := widget2.ThemeVariant
	widget2.ThemeVariant = func() fyne.ThemeVariant { return v }
	t.Cleanup(func() { widget2.ThemeVariant = prev })
}

// hexColor parses "#rrggbb" into an opaque color.
func hexColor(hex string) color.Color {
	v, err := strconv.ParseUint(strings.TrimPrefix(hex, "#"), 16, 32)
	if err != nil {
		panic(err)
	}
	return &color.RGBA{R: uint8(v >> 16), G: uint8(v >> 8), B: uint8(v), A: 0xff}
}

// rgba is the 16-bit premultiplied channel tuple of a color.
func rgba(c color.Color) [4]uint32 {
	r, g, b, a := c.RGBA()
	return [4]uint32{r, g, b, a}
}

// SGR foreground, SGR background and 256-color forms all route to the slot of the active variant.
func TestHandleOutput_PaletteSlots(t *testing.T) {
	variants := []struct {
		name    string
		variant fyne.ThemeVariant
	}{{"light", theme.VariantLight}, {"dark", theme.VariantDark}}
	for _, v := range variants {
		t.Run(v.name, func(t *testing.T) {
			pinVariant(t, v.variant)
			for slot := 0; slot < 16; slot++ {
				fg, bg := 30+slot%8, 40+slot%8
				if slot >= 8 {
					fg, bg = 90+slot-8, 100+slot-8
				}
				for _, tc := range []struct {
					name string
					seq  string
					slot uint8
					fg   bool
				}{
					{"SGR foreground", fmt.Sprintf("[%dm", fg), uint8(slot), true},
					{"SGR background", fmt.Sprintf("[%dm", bg), uint8(slot), false},
					{"256-color foreground", fmt.Sprintf("[38;5;%dm", slot), uint8(slot), true},
					{"256-color background", fmt.Sprintf("[48;5;%dm", slot), uint8(slot), false},
				} {
					t.Run(fmt.Sprintf("slot %d %s", slot, tc.name), func(t *testing.T) {
						term := New()
						term.handleOutput([]byte(esc(tc.seq)))
						got := term.currentBG
						if tc.fg {
							got = term.currentFG
						}
						assert.Equal(t, rgba(widget2.ANSI(tc.slot)), rgba(got))
					})
				}
			}
		})
	}
}
