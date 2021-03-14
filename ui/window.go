package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/schriftbank"
)

// Window object (hamburger button, variant name, undo, help buttons)
type Window struct {
	ContainerBase
	title *Label
}

// NewWindow creates a new toolbar
func NewWindow(input *input.Input, title string, content []string) *Window {
	w := &Window{ContainerBase: ContainerBase{input: input}} // x,y,width,height will be set when drawn
	if title != "" {
		w.title = NewLabel(w, input, 0, 0, 0, 48, 0, title, schriftbank.RobotoMedium24, "")
	}
	for _, c := range content {
		w.widgets = append(w.widgets, NewLabel(w, input, 0, 0, 0, 48, 0, c, schriftbank.RobotoRegular14, ""))
	}
	return w
}

// LayoutWidgets that belong to this container
func (w *Window) LayoutWidgets() {
	wpx0, wpy0, wpx1, _ := w.Rect()
	windowWidth := wpx1 - wpx0

	var titleWidth, titleHeight, x, y int
	if w.title != nil {
		x = wpx0 + (windowWidth / 2) // center of window
		titleWidth, titleHeight = w.title.Size()
		x -= titleWidth / 2
		w.title.SetPosition(x, wpy0+titleHeight)
		y = titleHeight + 48
	} else {
		y = 24
	}

	for _, w := range w.widgets {
		w.SetPosition(wpx0+48, wpy0+y)
		_, widgetHeight := w.Size()
		y += widgetHeight + 14
	}
}

// Update the window
func (w *Window) Update() {
	if w.title != nil {
		w.title.Update()
	}
	for _, w := range w.widgets {
		w.Update()
	}
}

// Draw the window
func (w *Window) Draw(screen *ebiten.Image) {
	width, height := screen.Size()
	if w.img == nil {
		w.width = width / 2
		w.height = height / 2
		w.x = (width - w.width) / 2
		w.y = (height - w.height) / 2
		w.img = w.createImg()
		w.LayoutWidgets()
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(w.x), float64(w.y))
	screen.DrawImage(w.img, op)

	if w.title != nil {
		w.title.Draw(screen)
	}
	for _, w := range w.widgets {
		w.Draw(screen)
	}
}

// OpenWindow create window
func (u *UI) OpenWindow(title string, content []string) {
	u.CloseActiveModal()
	u.modal = NewWindow(u.input, title, content)
}

// CloseWindow create window
func (u *UI) CloseWindow() {
	u.CloseActiveModal()
}
