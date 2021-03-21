package ui

import (
	"image"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/util"
)

// Label is a button that displays a single rune
type Label struct {
	WidgetBase
	text        string
	face        font.Face
	requestType string
}

func (l *Label) createImg() *ebiten.Image {
	dc := gg.NewContext(l.width, l.height)
	dc.SetRGBA(1, 1, 1, 1)
	// nota bene - text is drawn with y as a baseline
	dc.SetFontFace(l.face)
	dc.DrawString(l.text, 24, float64(l.height)*0.7)
	return ebiten.NewImageFromImage(dc.Image())
}

// NewLabel creates a new Label
func NewLabel(parent Container, input *input.Input, x, y, width, height, align int, text string, face font.Face, requestType string) *Label {
	l := &Label{
		WidgetBase: WidgetBase{parent: parent, input: input, img: nil, x: x, y: y, width: width, height: height, align: align},
		text:       text, face: face, requestType: requestType}
	l.Activate()
	return l
}

// Activate tells the input we need notifications
func (l *Label) Activate() {
	l.disabled = false
	l.img = l.createImg()
	l.input.Add(l)
}

// Deactivate tells the input we no longer need notofications
func (l *Label) Deactivate() {
	l.disabled = true
	l.img = l.createImg()
	l.input.Remove(l)
}

// NotifyCallback is called by the Subject (Input/Stroke) when something interesting happens
func (l *Label) NotifyCallback(event interface{}) {
	if l.disabled {
		return
	}
	switch v := event.(type) { // Type switch https://tour.golang.org/methods/16
	case image.Point:
		// println("Label image.Point", v.X, v.Y)
		if l.requestType != "" && util.InRect(v.X, v.Y, l.OffsetRect) {
			println("label notify", l.requestType, ":=", l.text)
			l.input.Notify(ChangeRequest{ChangeRequested: l.requestType, Data: l.text})
		}
	}
}
