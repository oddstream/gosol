package ui

import (
	"fmt"

	"oddstream.games/gosol/input"
	"oddstream.games/gosol/schriftbank"
)

// Statusbar object (hamburger button, variant name, undo, help buttons)
type Statusbar struct {
	BarBase
}

// NewStatusbar creates a new statusbar
func NewStatusbar(input *input.Input) *Statusbar {
	// img will created first time it's drawn if width == 0
	sb := &Statusbar{BarBase: BarBase{input: input, x: 0, y: 0, width: 0, height: 24}}

	sb.widgets = []Widget{
		// button's x will be set by LayoutWidgets()
		NewLabel(sb, input, -1, "Moves", schriftbank.RobotoRegular14, ""),
		NewLabel(sb, input, 1, "Complete", schriftbank.RobotoRegular14, ""),
	}
	return sb
}

// SetMoves of the statusbar
func (u *UI) SetMoves(moves int) {
	var l *Label = u.statusbar.widgets[0].(*Label)
	l.UpdateText(fmt.Sprintf("Moves %d", moves))
	// u.statusbar.LayoutWidgets()
}

// SetPercent of the statusbar
func (u *UI) SetPercent(percent int) {
	var l *Label = u.statusbar.widgets[1].(*Label)
	l.UpdateText(fmt.Sprintf("Complete %d%%", percent))
	// u.statusbar.LayoutWidgets()
}

// Layout implements Ebiten's Layout
func (sb *Statusbar) Layout(outsideWidth, outsideHeight int) (int, int) {
	// override BarBase.Layout to get screen height and position statusbar
	if sb.img == nil || outsideWidth != sb.width {
		sb.width = outsideWidth
		// sb.height is fixed (at 24)
		sb.img = sb.createImg()
		sb.LayoutWidgets()
	}
	sb.x, sb.y = 0, outsideHeight-sb.height
	return outsideWidth, outsideHeight
}
