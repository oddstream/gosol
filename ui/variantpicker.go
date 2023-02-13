package ui

import (
	"oddstream.games/gosol/schriftbank"
)

type Picker struct {
	DrawerBase
}

// NewVariantPicker creates a new container
func NewVariantPicker() *Picker {
	p := &Picker{DrawerBase: DrawerBase{WindowBase: WindowBase{x: -300, y: ToolbarHeight, width: 300}}} // height will be set when drawn
	return p
}

// ShowVariantPicker makes the variant picker visible
func (u *UI) ShowVariantPickerEx(content []string, widgetCommand string) {
	u.variantPicker.widgets = u.variantPicker.widgets[:0]
	for _, c := range content {
		u.variantPicker.widgets = append(u.variantPicker.widgets, NewLabel(u.variantPicker, "", 0, c, schriftbank.RobotoMedium24, widgetCommand))
	}
	u.variantPicker.ResetScroll()
	u.variantPicker.LayoutWidgets()
	u.variantPicker.Show()
}
