package ui

import (
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/schriftbank"
)

// Picker object (hamburger button, variant name, undo, help buttons)
type Picker struct {
	DrawerBase
}

// NewPicker creates a new container
func NewPicker(input *input.Input, content []string) *Picker {
	p := &Picker{DrawerBase: DrawerBase{input: input, x: -300, y: 48, width: 300}} // height will be set when drawn
	for _, c := range content {
		p.widgets = append(p.widgets, NewLabel(p, input, -300, 0, 300, 48, 0, c, schriftbank.RobotoRegular24, "Variant"))
	}
	p.LayoutWidgets()
	return p
}

// ShowVariantPicker makes the variant picker visible
func (u *UI) ShowVariantPicker() {
	con := u.VisibleDrawer()
	if con == u.variantPicker {
		return
	}
	if con != nil {
		con.Hide()
	}
	u.variantPicker.Show()
}
