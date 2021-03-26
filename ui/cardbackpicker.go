package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
)

// CardBackPicker object
type CardBackPicker struct {
	DrawerBase
}

// NewCardBackPicker creates a new container
func NewCardBackPicker(input *input.Input) *CardBackPicker {
	p := &CardBackPicker{DrawerBase: DrawerBase{input: input, x: -400, y: 48, width: 400}} // height will be set when drawn
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
	for name, img := range content {
		u.cardBackPicker.widgets = append(u.cardBackPicker.widgets, NewCardBackWidget(u.cardBackPicker, u.input, name, img))
	}
	u.cardBackPicker.LayoutWidgets()
	u.cardBackPicker.Show()
}
