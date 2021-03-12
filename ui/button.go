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
	WidgetBase
	r   rune
	key ebiten.Key
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
func NewRuneButton(parent Container, input *input.Input, x, y, width, height, align int, r rune, key ebiten.Key) *RuneButton {
	rb := &RuneButton{WidgetBase: WidgetBase{parent: parent, input: input, img: nil, x: x, y: y, width: width, height: height, align: align}, r: r, key: key}
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

// Update the state of this widget
func (rb *RuneButton) Update() {
}
