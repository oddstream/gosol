package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/schriftbank"
	"oddstream.games/gosol/util"
)

// Picker object (hamburger button, variant name, undo, help buttons)
type Picker struct {
	ContainerBase
}

// NewPicker creates a new toolbar
func NewPicker(input *input.Input, content []string) *Picker {
	p := &Picker{ContainerBase: ContainerBase{input: input}} // x,y,width,height will be set when drawn
	for _, c := range content {
		p.widgets = append(p.widgets, NewLabel(p, input, 0, 0, 0, 48, 0, c, schriftbank.RobotoRegular24, "Variant"))
	}
	return p
}

// LayoutWidgets that belong to this container
func (p *Picker) LayoutWidgets() {
	var x, y int
	x = 24
	y = 24

	for _, w := range p.widgets {
		w.SetPosition(p.x+x, p.y+y+p.yOffset)
		_, widgetHeight := w.Size()
		y += widgetHeight + 14
	}
	// println("yOffset is", p.yOffset)
}

// StartDrag this widget, if it is allowed
func (p *Picker) StartDrag() bool {
	// println("start drag with offset base", p.yOffsetBase)
	return true
}

// DragBy this widget
func (p *Picker) DragBy(dx, dy int) {
	p.xOffset = p.xOffsetBase + dx
	p.xOffset = util.ClampInt(p.xOffset, -p.width, 0)

	numWidgets := len(p.widgets)
	_, widgetHeight := p.widgets[0].Size()
	widgetHeight += 14 // see Picker.LayoutWidgets, which uses 24 (widget height) + 14 as vertical spacing
	_, pickerHeight := p.Size()
	// println("picker height", pickerHeight, "widget height", widgetHeight)
	visibleWidgets := pickerHeight / widgetHeight
	hiddenWidgets := numWidgets - visibleWidgets
	// println("total", numWidgets, "visible", visibleWidgets, "hidden", hiddenWidgets)
	heightOfHiddenWidgets := hiddenWidgets * widgetHeight
	p.yOffset = p.yOffsetBase + dy
	p.yOffset = util.ClampInt(p.yOffset, -heightOfHiddenWidgets, 0)
	p.LayoutWidgets()
}

// StopDrag this widget
func (p *Picker) StopDrag() {
	p.xOffsetBase = p.xOffset
	p.yOffsetBase = p.yOffset
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
