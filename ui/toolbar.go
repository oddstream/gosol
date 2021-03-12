package ui

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/schriftbank"
)

// Toolbar object (hamburger button, variant name, undo, help buttons)
type Toolbar struct {
	input         *input.Input
	img           *ebiten.Image
	title         string
	x, y          int
	width, height int
	widgets       []Widget
}

func (tb *Toolbar) createImg() *ebiten.Image {
	dc := gg.NewContext(tb.width, 48)
	dc.SetColor(color.RGBA{R: 0x32, G: 0x32, B: 0x32, A: 0xff})
	dc.DrawRectangle(0, 0, float64(tb.width), 48)
	dc.Fill()
	dc.Stroke()
	return ebiten.NewImageFromImage(dc.Image())
}

// NewToolbar creates a new toolbar
func NewToolbar(input *input.Input) *Toolbar {
	tb := &Toolbar{input: input, x: 0, y: 0, width: 0, height: 48}

	tb.widgets = []Widget{
		NewRuneButton(tb, input, 0, 0, 48, 48, -1, rune(9776), ebiten.KeyMenu),
		NewLabel(tb, input, 0, 0, 0, 48, 0, "", schriftbank.RobotoMedium24),
		NewRuneButton(tb, input, 0, 0, 48, 48, 1, '?', ebiten.KeyH),
		NewRuneButton(tb, input, 0, 0, 48, 48, 1, rune(8592), ebiten.KeyU),
	}
	// img will created first time it's drawn if width == 0
	return tb
}

// LayoutWidgets that belong to this container
func (tb *Toolbar) LayoutWidgets() {
	nextLeft := 0
	nextRight := tb.width - 48
	for _, w := range tb.widgets {
		switch w.Align() {
		case -1:
			w.SetPosition(nextLeft, tb.y)
			nextLeft += 48
		case 0:
			widgetWidth, widgetHeight := w.Size()
			w.SetPosition(tb.width/2-widgetWidth/2, tb.y+widgetHeight/2)
		case 1:
			w.SetPosition(nextRight, tb.y)
			nextRight -= 48
		}
	}
}

// Rect returns the area this toolbar covers
func (tb *Toolbar) Rect() (x0, y0, x1, y1 int) {
	x0 = 0
	y0 = 0
	x1 = tb.width
	y1 = tb.height
	return // using named parameters
}

// SetTitle of the toolbar
func (u *UI) SetTitle(title string) {
	u.toolbar.ReplaceWidget(1, NewLabel(u.toolbar, u.input, 0, 0, 0, 48, 0, title, schriftbank.RobotoMedium24))
	u.toolbar.width = 0 // force img to be recreated
}

// ReplaceWidget replaces a widget
func (tb *Toolbar) ReplaceWidget(n int, w Widget) {
	tb.widgets[n].Deactivate()
	tb.widgets[n] = w
}

// Update the toolbar
func (tb *Toolbar) Update() {
}

// Draw the toolbar
func (tb *Toolbar) Draw(screen *ebiten.Image) {
	w, _ := screen.Size()
	if tb.img == nil || w != tb.width {
		tb.width = w
		tb.img = tb.createImg()
		tb.LayoutWidgets()
	}
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(tb.img, op)

	for _, w := range tb.widgets {
		w.Draw(screen)
	}
}
