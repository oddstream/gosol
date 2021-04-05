package ui

import (
	"oddstream.games/gosol/input"
)

// TextDrawer provides a drawer for displaying rules of the current variant
type TextDrawer struct {
	DrawerBase
}

// NewTextDrawer creates a new container
func NewTextDrawer(input *input.Input) *TextDrawer {
	r := &TextDrawer{DrawerBase: DrawerBase{input: input, x: -300, y: 48, width: 300}} // height will be set when drawn
	return r
}

// ShowTextDrawer makes the text drawer container visible
func (u *UI) ShowTextDrawer(content []string) {
	con := u.VisibleDrawer()
	if con == u.textDrawer {
		return
	}
	if con != nil {
		con.Hide()
	}

	u.textDrawer.widgets = nil
	for _, c := range content { // content may be nil
		u.textDrawer.widgets = append(u.textDrawer.widgets, NewText(u.textDrawer, u.input, c))
	}
	u.textDrawer.LayoutWidgets()
	u.textDrawer.Show()
}
