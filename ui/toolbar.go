package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/schriftbank"
)

// Toolbar object (hamburger button, variant name, undo, help buttons)
type Toolbar struct {
	BarBase
}

// NewToolbar creates a new toolbar
func NewToolbar(input *input.Input) *Toolbar {
	tb := &Toolbar{BarBase: BarBase{input: input, x: 0, y: 0, width: 0, height: 48}}

	tb.widgets = []Widget{
		NewRuneButton(tb, input, 0, 0, 48, 48, -1, rune(9776), ebiten.KeyMenu),
		NewLabel(tb, input, 0, 0, 0, 48, 0, "", schriftbank.RobotoMedium24, ""),
		NewRuneButton(tb, input, 0, 0, 48, 48, 1, '?', ebiten.KeyH),
		// NewRuneButton(tb, input, 0, 0, 48, 48, 1, rune(0x238c), ebiten.KeyU),	// does not display unicode undo glyph
		NewRuneButton(tb, input, 0, 0, 48, 48, 1, rune(8592), ebiten.KeyU),
	}
	tb.widgets[2].Deactivate() // deactivate the help rune for now
	// img will created first time it's drawn if width == 0
	return tb
}

// SetTitle of the toolbar
func (u *UI) SetTitle(title string) {
	u.toolbar.ReplaceWidget(1, NewLabel(u.toolbar, u.input, 0, 0, 0, 48, 0, title, schriftbank.RobotoMedium24, ""))
	u.toolbar.width = 0 // force img to be recreated
}
