package ui

import (
	"oddstream.games/gomps5/schriftbank"
)

// Picker object (hamburger button, variant name, undo, help buttons)
type Picker struct {
	DrawerBase
}

// NewVariantPicker creates a new container
func NewVariantPicker() *Picker {
	p := &Picker{DrawerBase: DrawerBase{x: -300, y: 48, width: 300}} // height will be set when drawn
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
		u.variantPicker.widgets = append(u.variantPicker.widgets, NewLabel(u.variantPicker, 0, c, schriftbank.RobotoMedium24, "Variant"))
	}
	u.variantPicker.LayoutWidgets()
	u.variantPicker.Show()
}
