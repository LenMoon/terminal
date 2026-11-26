package terminal

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	widget2 "github.com/fyne-io/terminal/internal/widget"
	"github.com/stretchr/testify/assert"
)

// inverted is the selection counterpart of a color, the way termgridhelper.go derives it:
// each 8-bit channel XOR mask, alpha kept.
func inverted(c color.Color, mask uint8) color.RGBA {
	r, g, b, _ := c.RGBA()
	return color.RGBA{R: uint8(r>>8) ^ mask, G: uint8(g>>8) ^ mask, B: uint8(b>>8) ^ mask, A: 0xff}
}

func TestOnScreenColorsFollowTheme(t *testing.T) {
	pinVariant(t, theme.VariantLight)
	term := New()
	term.Resize(fyne.NewSize(800, 400))
	term.handleOutput([]byte(esc("[1;34mdir") + esc("[0m ") + esc("[91;44mcoloured ") + esc("[0m") +
		esc("[5;31mblink") + esc("[0m ") + esc("[38;2;1;2;3mtruecolor ") + esc("[38;5;196mcube")))
	// 列号：dir=0–2  空格=3  coloured=4–11  空格=12  blink=13–17  空格=18  truecolor=19–27  空格=28  cube=29–32
	fg := func(col int) color.Color { return term.content.Rows[0].Cells[col].Style.TextColor() }
	bg := func(col int) color.Color { return term.content.Rows[0].Cells[col].Style.BackgroundColor() }

	assert.Equal(t, rgba(hexColor("#0000aa")), rgba(fg(0)))
	assert.True(t, term.content.Rows[0].Cells[0].Style.Style().Bold)
	assert.Equal(t, rgba(hexColor("#cc0000")), rgba(fg(4)))
	assert.Equal(t, rgba(hexColor("#0000aa")), rgba(bg(4)))
	assert.Equal(t, rgba(hexColor("#aa0000")), rgba(fg(13))) // blinking cell
	assert.Equal(t, rgba(&color.RGBA{1, 2, 3, 255}), rgba(fg(19)))
	assert.Equal(t, rgba(&color.RGBA{255, 0, 0, 255}), rgba(fg(29))) // 256-color cube entry

	lightBlue := fg(0)
	assert.True(t, lightBlue == fg(0)) // stable within one variant: same cache key

	pinVariant(t, theme.VariantDark)
	assert.Equal(t, rgba(hexColor("#7b7bff")), rgba(fg(0)))
	assert.Equal(t, rgba(hexColor("#ff8a80")), rgba(fg(4)))
	assert.Equal(t, rgba(hexColor("#7b7bff")), rgba(bg(4)))
	assert.Equal(t, rgba(hexColor("#ff5c57")), rgba(fg(13)))
	assert.Equal(t, rgba(&color.RGBA{1, 2, 3, 255}), rgba(fg(19))) // explicit colors do not follow
	assert.Equal(t, rgba(&color.RGBA{255, 0, 0, 255}), rgba(fg(29)))
	assert.True(t, lightBlue != fg(0)) // the two variants are different cache keys

	widget2.HighlightRange(term.content, false, 0, 4, 0, 4, highlightBitMask)
	// selected after the flip: both inverted colors derive from the active variant's palette
	assert.Equal(t, rgba(inverted(hexColor("#ff8a80"), highlightBitMask)), rgba(fg(4)))
	assert.Equal(t, rgba(inverted(hexColor("#7b7bff"), highlightBitMask)), rgba(bg(4)))
	widget2.ClearHighlightRange(term.content, false, 0, 4, 0, 4)
	assert.Equal(t, rgba(hexColor("#ff8a80")), rgba(fg(4))) // a cell that was selected keeps following
}

// ECH and ICH write blanks that carry the SGR background.
func TestBlanksFollowTheme(t *testing.T) {
	for _, tc := range []struct {
		name string
		seq  string
	}{
		{"erase chars", esc("[45m") + esc("[2X")},
		{"insert chars", esc("[45m") + esc("[2@")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pinVariant(t, theme.VariantLight)
			term := New()
			term.Resize(fyne.NewSize(800, 400))
			term.handleOutput([]byte(tc.seq))
			assert.Equal(t, rgba(hexColor("#aa00aa")), rgba(term.content.Rows[0].Cells[0].Style.BackgroundColor()))
			pinVariant(t, theme.VariantDark)
			assert.Equal(t, rgba(hexColor("#d670d6")), rgba(term.content.Rows[0].Cells[0].Style.BackgroundColor()))
		})
	}
}
