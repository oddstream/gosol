package ui

import (
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/schriftbank"
	"oddstream.games/gosol/util"
)

type CardBackWidget struct {
	WidgetBase
	name    string
	backImg *ebiten.Image
}

func (cb *CardBackWidget) createImg() *ebiten.Image {

	w, _ := cb.backImg.Size()

	dc := gg.NewContext(cb.width, cb.height)

	dc.DrawImage(cb.backImg, 24, 0)

	// nota bene - text is drawn with y as a baseline

	dc.SetRGBA(1, 1, 1, 1)
	dc.SetFontFace(schriftbank.RobotoMedium24)
	dc.DrawString(cb.name, float64(24+w+24), float64(cb.height)*0.6)

	// uncomment this to show the area we expect the text to occupy
	// dc.DrawLine(0, float64(0), float64(cb.width), float64(0))
	// dc.DrawLine(0, float64(cb.height), float64(cb.width), float64(cb.height))
	// dc.DrawLine(0, float64(0), float64(cb.width), float64(cb.height))
	// dc.Stroke()

	return ebiten.NewImageFromImage(dc.Image())
}

// NewCardBackWidget creates a new cardBack widget for the CardBackPicker
func NewCardBackWidget(parent Containery, id string, name string, backImg *ebiten.Image) *CardBackWidget {
	_, h := backImg.Size()
	w, _ := parent.Size()
	// widget x, y will be set by LayoutWidgets
	cb := &CardBackWidget{WidgetBase: WidgetBase{parent: parent, id: id, width: w, height: h},
		name: name, backImg: backImg}
	return cb
}

// Activate tells the input we need notifications
func (cb *CardBackWidget) Activate() {
	cb.disabled = false
	cb.img = cb.createImg()
	// cb.input.Add(cb)
}

// Deactivate tells the input we no longer need notifications
func (cb *CardBackWidget) Deactivate() {
	cb.disabled = true
	cb.img = cb.createImg()
	// cb.input.Remove(cb)
}

// NotifyCallback is called by the Subject (Input/Stroke) when something interesting happens
func (cb *CardBackWidget) NotifyCallback(v input.StrokeEvent) {
	switch v.Event {
	case input.Tap:
		if util.InRect(v.X, v.Y, cb.OffsetRect) {
			cmdFn(ChangeRequest{ChangeRequested: "CardBack", Data: cb.name})
		}
	}
}
