package ui

import (
	"fmt"

	"oddstream.games/gosol/schriftbank"
)

// Statusbar object (hamburger button, variant name, undo, help buttons)
type Statusbar struct {
	BarBase
}

// NewStatusbar creates a new statusbar
func NewStatusbar() *Statusbar {
	// img will created first time it's drawn if width == 0
	sb := &Statusbar{BarBase: BarBase{WindowBase: WindowBase{x: 0, y: 0, width: 0, height: StatusbarHeight}}}

	sb.widgets = []Widgety{
		// button's x will be set by LayoutWidgets()
		NewLabel(sb, "statusbarStock", -1, "", schriftbank.RobotoRegular14, ""),  // 0 stock
		NewLabel(sb, "statusbarWaste", -1, "", schriftbank.RobotoRegular14, ""),  // 1 waste
		NewLabel(sb, "statusbarMiddle", 0, "", schriftbank.RobotoRegular14, ""),  // 2 middle (debug)
		NewLabel(sb, "statusbarPercent", 1, "", schriftbank.RobotoRegular14, ""), // 3 percent
	}
	return sb
}

// SetStock of the statusbar
func (u *UI) SetStock(cards int) {
	var l *Label = u.statusbar.widgets[0].(*Label)
	if cards == -1 {
		l.UpdateText("") // hide hidden stock
	} else {
		l.UpdateText(fmt.Sprintf("STOCK: %d", cards))
	}
	u.statusbar.LayoutWidgets()
}

// SetWaste of the statusbar
func (u *UI) SetWaste(cards int) {
	var l *Label = u.statusbar.widgets[1].(*Label)
	if cards == -1 {
		l.UpdateText("")
	} else {
		l.UpdateText(fmt.Sprintf("WASTE: %d", cards))
	}
	u.statusbar.LayoutWidgets()
}

// SetPercent of the statusbar
func (u *UI) SetMiddle(str string) {
	var l *Label = u.statusbar.widgets[2].(*Label)
	l.UpdateText(str)
	u.statusbar.LayoutWidgets()
}

// SetPercent of the statusbar
func (u *UI) SetPercent(percent int) {
	var l *Label = u.statusbar.widgets[3].(*Label)
	if percent == 100 {
		l.UpdateText("COMPLETE")
	} else {
		l.UpdateText(fmt.Sprintf("COMPLETE: %d%%", percent))
	}
	u.statusbar.LayoutWidgets()
}

// Layout implements Ebiten's Layout
func (sb *Statusbar) Layout(outsideWidth, outsideHeight int) (int, int) {
	// override BarBase.Layout to get screen height and position statusbar
	if sb.img == nil || outsideWidth != sb.width {
		sb.width = outsideWidth
		// sb.height is fixed (at 24)
		sb.img = sb.createImg(BackgroundColor)
		sb.LayoutWidgets()
	}
	sb.x, sb.y = 0, outsideHeight-sb.height
	return outsideWidth, outsideHeight
}
