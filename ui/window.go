package ui

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/schriftbank"
)

// Window object (hamburger button, variant name, undo, help buttons)
type Window struct {
	img           *ebiten.Image
	input         *input.Input
	title         *Label
	x, y          int
	width, height int
	widgets       []Widget
}

func (w *Window) createImg() *ebiten.Image {
	dc := gg.NewContext(w.width, w.height)
	dc.SetColor(color.RGBA{R: 0x32, G: 0x32, B: 0x32, A: 0xff})
	dc.DrawRectangle(0, 0, float64(w.width), float64(w.height))
	dc.Fill()
	dc.Stroke()
	return ebiten.NewImageFromImage(dc.Image())
}

// NewWindow creates a new toolbar
func NewWindow(input *input.Input, title string) *Window {
	w := &Window{input: input} // x,y,width,height will be set when drawn
	w.title = NewLabel(w, input, 0, 0, 0, 48, 0, title, schriftbank.RobotoMedium24)
	return w
}

// Rect gives the screen position
func (w *Window) Rect() (x0, y0, x1, y1 int) {
	x0 = w.x
	y0 = w.y
	x1 = w.x + w.width
	y1 = w.y + w.height
	return // using named parameters
}

// LayoutWidgets that belong to this container
func (w *Window) LayoutWidgets() {
	wpx, wpy, x1, _ := w.Rect()
	wsx := x1 - wpx
	x := wpx + (wsx / 2) // center of window
	lsx, lsy := w.title.Size()
	x -= lsx / 2
	w.title.SetPosition(x, wpy+lsy)

	y := lsy + 48
	for _, w := range w.widgets {
		w.SetPosition(wpx+48, wpy+y)
		y += 48
	}
}

// AppendText to this window
func (w *Window) AppendText(text string) {
	w.widgets = append(w.widgets,
		NewLabel(w, w.input, 0, 0, 640, 48, 0, text, schriftbank.RobotoRegular14),
	)
	w.LayoutWidgets()
}

// Update the window
func (w *Window) Update() {
	w.title.Update()
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

	w.title.Draw(screen)
	for _, w := range w.widgets {
		w.Draw(screen)
	}
}

// OpenWindow create window
func (u *UI) OpenWindow(input *input.Input, title string) {
	if u.window != nil {
		u.window = nil
	}
	u.window = NewWindow(input, title)
}

// AppendTextToWindow create window
func (u *UI) AppendTextToWindow(text string) {
	if u.window == nil {
		return
	}
	u.window.AppendText(text)
}

// CloseWindow create window
func (u *UI) CloseWindow() {
	if u.window != nil {
		u.window = nil
	}
}
