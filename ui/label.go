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
	text string
	face font.Face
}

func (l *Label) createImg() *ebiten.Image {
	dc := gg.NewContext(8, 8)
	dc.SetFontFace(l.face)
	w, h := dc.MeasureString(l.text)

	l.width = int(w)
	l.height = int(h)

	dc = gg.NewContext(int(w), int(h))
	dc.SetFontFace(l.face)
	if l.disabled {
		dc.SetRGBA(0.5, 0.5, 0.5, 1)
	} else {
		dc.SetRGBA(1, 1, 1, 1)
	}
	dc.DrawStringAnchored(l.text, w/2, h/2, 0.5, 0.4)
	dc.Stroke()

	return ebiten.NewImageFromImage(dc.Image())
}

// NewLabel creates a new Label
func NewLabel(parent Container, input *input.Input, x, y, width, height, align int, text string, face font.Face) *Label {
	l := &Label{WidgetBase: WidgetBase{parent: parent, input: input, img: nil, x: x, y: y, width: width, height: height, align: align}, text: text, face: face}
	l.Activate()
	return l
}

// Activate tells the input we need notifications
func (l *Label) Activate() {
	l.disabled = false
	l.img = l.createImg() // also sets width, height
	l.input.Add(l)
}

// Deactivate tells the input we no longer need notofications
func (l *Label) Deactivate() {
	l.disabled = true
	l.img = l.createImg() // also sets width, height
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
		if util.InRect(v.X, v.Y, l.Rect) {
			l.input.Notify(l.text)
		}
	}
}

// Update the state of this widget
func (l *Label) Update() {

}
