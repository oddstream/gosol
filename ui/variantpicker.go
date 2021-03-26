package ui

import (
	"oddstream.games/gosol/input"
)

// Picker object (hamburger button, variant name, undo, help buttons)
type Picker struct {
	DrawerBase
}

// NewVariantPicker creates a new container
func NewVariantPicker(input *input.Input) *Picker {
	p := &Picker{DrawerBase: DrawerBase{input: input, x: -300, y: 48, width: 300}} // height will be set when drawn
	return p
}

// ShowVariantPicker makes the variant picker visible
func (u *UI) ShowVariantPicker(content []string) {
	con := u.VisibleDrawer()
	if con == u.variantPicker {
		return
	}
	if con != nil {
		con.Hide()
	}
	u.variantPicker.widgets = u.variantPicker.widgets[:0]
	for _, c := range content {
		u.variantPicker.widgets = append(u.variantPicker.widgets, NewLabel(u.variantPicker, u.input, 0, c, "Variant"))
	}
	u.variantPicker.LayoutWidgets()
	u.variantPicker.Show()
}
