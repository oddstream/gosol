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
	r             rune
	align         int
	x, y          int // screen position
	width, height int
	input         *input.Input
	key           ebiten.Key
}

// NewRuneButton creates a new RuneButton
func NewRuneButton(parent Container, r rune, align int, input *input.Input, key ebiten.Key) *RuneButton {
	rb := &RuneButton{parent: parent, r: r, align: align, width: 48, height: 48, input: input, key: key}
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

// Size of the RuneButton
// func (rb *RuneButton) Size() (int, int) {
// 	return rb.width, rb.height
// }

// Rect gives the screen position
func (rb *RuneButton) Rect() (x0, y0, x1, y1 int) {
	x0 = rb.x
	y0 = rb.y
	x1 = rb.x + rb.width
	y1 = rb.y + rb.height
	return // using named parameters
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

// Draw into a gg context, not to the screen; x,y is the center of the rune
func (rb *RuneButton) Draw(dc *gg.Context, x, y int) {
	rb.x, rb.y = x-(rb.width/2), y-(rb.height/2)
	dc.SetFontFace(schriftbank.Symbol24)
	dc.SetRGBA(1, 1, 1, 1)
	dc.DrawStringAnchored(string(rb.r), float64(x), float64(y), 0.5, 0.5)
	dc.Stroke()
}
