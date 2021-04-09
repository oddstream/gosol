package ui

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

// CardBackPicker object
type CardBackPicker struct {
	DrawerBase
}

// NewCardBackPicker creates a new container
func NewCardBackPicker() *CardBackPicker {
	p := &CardBackPicker{DrawerBase: DrawerBase{x: -400, y: 48, width: 400}} // height will be set when drawn
	return p
}

// ShowCardBackPicker makes the card back picker visible
func (u *UI) ShowCardBackPicker(content map[string]*ebiten.Image) {
	con := u.VisibleDrawer()
	if con == u.cardBackPicker {
		return
	}
	if con != nil {
		con.Hide()
	}
	u.cardBackPicker.widgets = u.cardBackPicker.widgets[:0]
	strings := []string{}
	for name := range content {
		strings = append(strings, name)
	}
	sort.Slice(strings, func(i, j int) bool { return strings[i] < strings[j] })
	for _, name := range strings {
		u.cardBackPicker.widgets = append(u.cardBackPicker.widgets, NewCardBackWidget(u.cardBackPicker, name, content[name]))
	}
	u.cardBackPicker.ResetScroll()
	u.cardBackPicker.LayoutWidgets()
	u.cardBackPicker.Show()
}
