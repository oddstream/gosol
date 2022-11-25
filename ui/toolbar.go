package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/schriftbank"
)

// Toolbar object (hamburger button, variant name, undo, help buttons)
type Toolbar struct {
	BarBase
}

// NewToolbar creates a new toolbar
func NewToolbar() *Toolbar {
	// img will created first time it's drawn if width == 0
	tb := &Toolbar{BarBase: BarBase{x: 0, y: 0, width: 0, height: 48}}

	tb.widgets = []Widget{
		// button's x will be set by LayoutWidgets() (y will always be 0 in a toolbar)
		NewIconButton(tb, 0, 0, 48, 48, -1, "menu", ebiten.KeyMenu),
		NewLabel(tb, 0, "title", schriftbank.RobotoMedium24, ""),
		NewIconButton(tb, 0, 0, 48, 48, 1, "undo", ebiten.KeyU),      // U for Undo
		NewIconButton(tb, 0, 0, 48, 48, 1, "done", ebiten.KeyC),      // C for Collect
		NewIconButton(tb, 0, 0, 48, 48, 1, "lightbulb", ebiten.KeyH), // H for Hint
	}
	return tb
}

// SetTitle of the toolbar
func (u *UI) SetTitle(title string) {
	var l *Label = u.toolbar.widgets[1].(*Label)
	l.UpdateText(title)
	// u.toolbar.LayoutWidgets()
}

// Layout implements Ebiten's Layout
func (tb *Toolbar) Layout(outsideWidth, outsideHeight int) (int, int) {
	// override BarBase.Layout to get screen height and position bar
	if tb.img == nil || outsideWidth != tb.width {
		tb.width = outsideWidth
		// tb.height is fixed (at 48)
		tb.img = tb.createImg()
		tb.LayoutWidgets()
	}
	// tb.x, tb.y = 0, 0
	return outsideWidth, outsideHeight
}
