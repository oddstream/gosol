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
func NewCardBackPicker(input *input.Input, content map[string]*ebiten.Image) *CardBackPicker {
	p := &CardBackPicker{DrawerBase: DrawerBase{input: input, x: -71, y: 48, width: 71}} // height will be set when drawn
	for name, img := range content {
		p.widgets = append(p.widgets, NewCardBackWidget(p, input, name, img))
	}
	p.LayoutWidgets()
	// for _, w := range p.widgets {
	// 	println(w.Rect())
	// }
	return p
}

// ShowCardBackPicker makes the card back picker visible
func (u *UI) ShowCardBackPicker() {
	con := u.VisibleDrawer()
	if con == u.cardBackPicker {
		return
	}
	if con != nil {
		con.Hide()
	}
	u.cardBackPicker.Show()
}
