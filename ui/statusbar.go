package ui

import (
	"fmt"

	"oddstream.games/gomps5/schriftbank"
)

// Statusbar object (hamburger button, variant name, undo, help buttons)
type Statusbar struct {
	BarBase
}

// NewStatusbar creates a new statusbar
func NewStatusbar() *Statusbar {
	// img will created first time it's drawn if width == 0
	sb := &Statusbar{BarBase: BarBase{x: 0, y: 0, width: 0, height: 24}}

	sb.widgets = []Widget{
		// button's x will be set by LayoutWidgets()
		NewLabel(sb, -1, "", schriftbank.RobotoRegular14, ""),
		NewLabel(sb, -1, "", schriftbank.RobotoRegular14, ""),
		NewLabel(sb, 0, "", schriftbank.RobotoRegular14, ""),
		NewLabel(sb, 1, "", schriftbank.RobotoRegular14, ""),
	}
	return sb
}

// SetStock of the statusbar
func (u *UI) SetStock(cards int) {
	var l *Label = u.statusbar.widgets[0].(*Label)
	switch cards {
	case 0:
		l.UpdateText("") // hide hidden stock
	// case 1:
	// 	l.UpdateText("1 STOCK CARD")
	// default:
	// 	l.UpdateText(fmt.Sprintf("%d STOCK CARDS", cards))
	default:
		l.UpdateText(fmt.Sprintf("STOCK: %d", cards))
	}
	u.statusbar.LayoutWidgets()
}

// SetWaste of the statusbar
func (u *UI) SetWaste(cards int) {
	var l *Label = u.statusbar.widgets[1].(*Label)
	switch cards {
	case 0:
		l.UpdateText("") // hide hidden stock
		// case 1:
		// 	l.UpdateText("1 WASTE CARD")
		// default:
		// 	l.UpdateText(fmt.Sprintf("%d WASTE CARDS", cards))
	default:
		l.UpdateText(fmt.Sprintf("WASTE: %d", cards))
	}
	u.statusbar.LayoutWidgets()
}

// SetMoves of the statusbar
func (u *UI) SetMoves(moves int) {
	var l *Label = u.statusbar.widgets[2].(*Label)
	switch moves {
	case 0:
		l.UpdateText("NO MOVES MADE")
	// case 1:
	// 	l.UpdateText("1 MOVE")
	// default:
	// 	l.UpdateText(fmt.Sprintf("%d MOVES", moves))
	default:
		l.UpdateText(fmt.Sprintf("MOVES: %d", moves))
	}
	u.statusbar.LayoutWidgets()
}

// SetPercent of the statusbar
func (u *UI) SetPercent(percent int) {
	var l *Label = u.statusbar.widgets[3].(*Label)
	l.UpdateText(fmt.Sprintf("COMPLETE: %d%%", percent))
	u.statusbar.LayoutWidgets()
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
