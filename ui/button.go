package ui

import (
	"image"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/schriftbank"
	"oddstream.games/gosol/util"
)

// RuneButton is a button that displays a single rune
type RuneButton struct {
	parent        Container
	img           *ebiten.Image
	r             rune
	align         int
	x, y          int // screen position
	width, height int // always 48x48
	input         *input.Input
	key           ebiten.Key
}

func (rb *RuneButton) createImg() *ebiten.Image {
	dc := gg.NewContext(rb.width, rb.height)
	dc.SetRGBA(1, 1, 1, 1)
	dc.SetFontFace(schriftbank.Symbol24)
	dc.DrawStringAnchored(string(rb.r), float64(rb.width/2), float64(rb.height/2), 0.5, 0.5)
	dc.Stroke()
	return ebiten.NewImageFromImage(dc.Image())
}

// NewRuneButton creates a new RuneButton
func NewRuneButton(parent Container, input *input.Input, r rune, align int, key ebiten.Key) *RuneButton {
	rb := &RuneButton{parent: parent, r: r, align: align, width: 48, height: 48, input: input, key: key}
	rb.img = rb.createImg()
	rb.Activate()
	return rb
}

// Activate tells the input we need notifications
func (rb *RuneButton) Activate() {
	rb.input.Add(rb)
}

// Deactivate tells the input we no longer need notifications
func (rb *RuneButton) Deactivate() {
	rb.input.Remove(rb)
}

// Position of the widget
func (rb *RuneButton) Position() (int, int) {
	return rb.x, rb.y
}

// Size of the widget
func (rb *RuneButton) Size() (int, int) {
	return rb.width, rb.height
}

// Rect gives the screen position
func (rb *RuneButton) Rect() (x0, y0, x1, y1 int) {
	x0 = rb.x
	y0 = rb.y
	x1 = rb.x + rb.width
	y1 = rb.y + rb.height
	return // using named parameters
}

// SetPosition of this widget
func (rb *RuneButton) SetPosition(x, y int) {
	rb.x, rb.y = x, y
}

// NotifyCallback is called by the Subject (Input/Stroke) when something interesting happens
func (rb *RuneButton) NotifyCallback(event interface{}) {
	switch v := event.(type) { // Type switch https://tour.golang.org/methods/16
	case image.Point:
		// println("RuneButton image.Point", v.X, v.Y)
		if util.InRect(v.X, v.Y, rb.Rect) {
			rb.input.Notify(rb.key)
		}
	}
}

// Align returns the x axis alignment (-1, 0, 1)
func (rb *RuneButton) Align() int {
	return rb.align
}

// Update the state of this widget
func (rb *RuneButton) Update() {

}

// Draw the widget
func (rb *RuneButton) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(rb.x), float64(rb.y))
	screen.DrawImage(rb.img, op)
}
