package ui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/util"
)

type CardBackWidget struct {
	WidgetBase
	name string
}

// NewCardBackWidget creates a new cardBack widget for the CardBackPicker
func NewCardBackWidget(parent Container, input *input.Input, name string, img *ebiten.Image) *CardBackWidget {
	w, h := img.Size()
	if img == nil {
		println("warning nil img")
	}
	cb := &CardBackWidget{WidgetBase: WidgetBase{parent: parent, input: input, x: -256, y: 0, width: w, height: h, img: img},
		name: name}
	cb.input.Add(cb)
	return cb
}

// Activate tells the input we need notifications
func (cb *CardBackWidget) Activate() {
	cb.disabled = false
	cb.input.Add(cb)
}

// Deactivate tells the input we no longer need notifications
func (cb *CardBackWidget) Deactivate() {
	cb.disabled = true
	cb.input.Remove(cb)
}

// NotifyCallback is called by the Subject (Input/Stroke) when something interesting happens
func (cb *CardBackWidget) NotifyCallback(event interface{}) {
	switch v := event.(type) { // Type switch https://tour.golang.org/methods/16
	case image.Point:
		// println("Label image.Point", v.X, v.Y)
		if util.InRect(v.X, v.Y, cb.OffsetRect) {
			println("card back notify", cb.name)
			cb.input.Notify(ChangeRequest{ChangeRequested: "cardback", Data: cb.name})
		}
	}
}
