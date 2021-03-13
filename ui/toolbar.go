package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/schriftbank"
)

// Toolbar object (hamburger button, variant name, undo, help buttons)
type Toolbar struct {
	ContainerBase
}

// NewToolbar creates a new toolbar
func NewToolbar(input *input.Input) *Toolbar {
	tb := &Toolbar{ContainerBase: ContainerBase{input: input, x: 0, y: 0, width: 0, height: 48}}

	tb.widgets = []Widget{
		NewRuneButton(tb, input, 0, 0, 48, 48, -1, rune(9776), ebiten.KeyMenu),
		NewLabel(tb, input, 0, 0, 0, 48, 0, "", schriftbank.RobotoMedium24),
		NewRuneButton(tb, input, 0, 0, 48, 48, 1, '?', ebiten.KeyH),
		// NewRuneButton(tb, input, 0, 0, 48, 48, 1, rune(0x238c), ebiten.KeyU),	// does not display unicode undo glyph
		NewRuneButton(tb, input, 0, 0, 48, 48, 1, rune(8592), ebiten.KeyU),
	}
	tb.widgets[2].Deactivate()
	// img will created first time it's drawn if width == 0
	return tb
}

// LayoutWidgets that belong to this container
func (tb *Toolbar) LayoutWidgets() {
	nextLeft := 0
	nextRight := tb.width - 48
	for _, w := range tb.widgets {
		switch w.Align() {
		case -1:
			w.SetPosition(nextLeft, tb.y)
			nextLeft += 48
		case 0:
			widgetWidth, widgetHeight := w.Size()
			w.SetPosition(tb.width/2-widgetWidth/2, tb.y+widgetHeight/2)
		case 1:
			w.SetPosition(nextRight, tb.y)
			nextRight -= 48
		}
	}
}

// SetTitle of the toolbar
func (u *UI) SetTitle(title string) {
	u.toolbar.ReplaceWidget(1, NewLabel(u.toolbar, u.input, 0, 0, 0, 48, 0, title, schriftbank.RobotoMedium24))
	u.toolbar.width = 0 // force img to be recreated
}

// ReplaceWidget replaces a widget
func (tb *Toolbar) ReplaceWidget(n int, w Widget) {
	tb.widgets[n].Deactivate()
	tb.widgets[n] = w
}

// Update the toolbar
func (tb *Toolbar) Update() {
}

// Draw the toolbar
func (tb *Toolbar) Draw(screen *ebiten.Image) {
	w, _ := screen.Size()
	if tb.img == nil || w != tb.width {
		tb.width = w
		tb.img = tb.createImg()
		tb.LayoutWidgets()
	}
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(tb.img, op)

	for _, w := range tb.widgets {
		w.Draw(screen)
	}
}
