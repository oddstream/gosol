package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/schriftbank"
)

// Picker object (hamburger button, variant name, undo, help buttons)
type Picker struct {
	ContainerBase
}

// NewPicker creates a new toolbar
func NewPicker(input *input.Input, content []string) *Picker {
	p := &Picker{ContainerBase: ContainerBase{input: input}} // x,y,width,height will be set when drawn
	for _, c := range content {
		p.widgets = append(p.widgets, NewLabel(p, input, 0, 0, 0, 48, 0, c, schriftbank.RobotoRegular24))
	}
	return p
}

// LayoutWidgets that belong to this container
func (p *Picker) LayoutWidgets() {
	wpx0, wpy0, _, _ := p.Rect()

	var x, y int
	x = 24
	y = 24

	for _, w := range p.widgets {
		w.SetPosition(wpx0+x, wpy0+y+p.yOffset)
		_, widgetHeight := w.Size()
		y += widgetHeight + 14
	}
}

// Update the window
func (p *Picker) Update() {
	for _, w := range p.widgets {
		w.Update()
	}
}

// Draw the window
func (p *Picker) Draw(screen *ebiten.Image) {
	width, height := screen.Size()
	if p.img == nil {
		p.width = width / 2
		p.height = height / 2
		p.x = (width - p.width) / 2
		p.y = (height - p.height) / 2
		p.img = p.createImg()
		p.LayoutWidgets()
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.x), float64(p.y))
	screen.DrawImage(p.img, op)

	for _, w := range p.widgets {
		w.Draw(screen)
	}
}

// OpenPicker create window
func (u *UI) OpenPicker(content []string) {
	u.CloseActiveModal()
	u.modal = NewPicker(u.input, content)
}

// ClosePicker create window
func (u *UI) ClosePicker() {
	u.CloseActiveModal()
}
