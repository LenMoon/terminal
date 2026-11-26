package terminal

import (
	"image/color"

	"fyne.io/fyne/v2/widget"
	widget2 "github.com/fyne-io/terminal/internal/widget"
)

// cellStyle is the style of a terminal cell. It answers palette colors for the
// theme variant that is active when the cell is painted, so text already on
// screen follows a theme change: Fyne refreshes every widget when the theme
// changes and asks each cell style for its colors again.
type cellStyle struct {
	widget.CustomTextGridStyle
}

func (s *cellStyle) TextColor() color.Color       { return widget2.Resolve(s.FGColor) }
func (s *cellStyle) BackgroundColor() color.Color { return widget2.Resolve(s.BGColor) }
