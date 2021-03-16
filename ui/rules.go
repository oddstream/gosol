package ui

import (
	"oddstream.games/gosol/input"
)

// Picker object (hamburger button, variant name, undo, help buttons)
type Rules struct {
	DrawerBase
}

// NewRules creates a new container
func NewRules(input *input.Input, content []string) *Rules {
	r := &Rules{DrawerBase: DrawerBase{input: input, x: -300, y: 48, width: 300}} // height will be set when drawn
	for _, c := range content {                                                   // content may be nil
		r.widgets = append(r.widgets, NewText(r, input, -300, 0, 300, 48, c))
	}
	return r
}

// ShowRules makes the variant picker visible
func (u *UI) ShowRules(content []string) {
	con := u.VisibleDrawer()
	if con == u.rules {
		return
	}
	if con != nil {
		con.Hide()
	}

	u.rules.widgets = nil
	for _, c := range content { // content may be nil
		u.rules.widgets = append(u.rules.widgets, NewText(u.rules, u.input, -300, 0, 300, 48, c))
	}

	u.rules.Show()
}
