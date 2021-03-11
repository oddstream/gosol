package ui

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/schriftbank"
)

// Toolbar object (hamburger button, variant name, undo button)
type Toolbar struct {
	img     *ebiten.Image
	title   string
	width   int
	widgets []Widget
}

// NewToolbar creates a new toolbar
func NewToolbar(input *input.Input) *Toolbar {
	tb := &Toolbar{}

	tb.widgets = []Widget{
		NewRuneButton(rune(9776), -1, input, ebiten.KeyMenu),
		NewLabel("", 0, schriftbank.RobotoMedium24, input),
		NewRuneButton('?', 1, input, ebiten.KeyH),
		NewRuneButton(rune(8592), 1, input, ebiten.KeyU),
	}

	return tb
}

func (tb *Toolbar) createImg() {

	dc := gg.NewContext(tb.width, 48)
	dc.SetColor(color.RGBA{R: 0x32, G: 0x32, B: 0x32, A: 0xff})
	dc.DrawRectangle(0, 0, float64(tb.width), 48)
	dc.Fill()
	dc.Stroke()

	width := dc.Width()
	height := dc.Height()
	nextLeft := 48 / 2
	nextRight := width - 48/2
	for _, w := range tb.widgets {
		switch w.Align() {
		case -1:
			w.Draw(dc, nextLeft, height/2)
			nextLeft += 48
		case 0:
			w.Draw(dc, width/2, height/2)
		case 1:
			w.Draw(dc, nextRight, height/2)
			nextRight -= 48
		}
	}
	tb.img = ebiten.NewImageFromImage(dc.Image())
}

// Rect returns the area this toolbar covers
func (tb *Toolbar) Rect() (x0, y0, x1, y1 int) {
	x0 = 0
	y0 = 0
	x1 = tb.width
	y1 = 48
	return // using named parameters
}

// SetTitle of the toolbar
func (u *UI) SetTitle(title string) {
	u.toolbar.ReplaceWidget(1, NewLabel(title, 0, schriftbank.RobotoMedium24, u.input))
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
		tb.createImg()
	}
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(tb.img, op)
}
