package ui

import (
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/schriftbank"
)

// Text is a button that displays a single rune
type Text struct {
	WidgetBase
	text string
}

func (t *Text) createImg() *ebiten.Image {
	dc := gg.NewContext(t.width, t.height)
	dc.SetFontFace(schriftbank.RobotoRegular14)
	// MeasureString says this text, requested to be 48 high, is 14 high
	lines := dc.WordWrap(t.text, float64(t.width-48)) // 24 padding left and right
	lineHeight := 24

	t.height = lineHeight*len(lines) + lineHeight
	dc = gg.NewContext(t.width, t.height)

	// nota bene - text is drawn with y as a baseline
	dc.SetFontFace(schriftbank.RobotoRegular14)
	y := lineHeight
	for _, str := range lines {
		dc.DrawString(str, 24, float64(y))
		y += lineHeight
	}
	// uncomment this line to visualize text box
	// dc.DrawLine(0, 0, float64(t.width), float64(t.height))
	dc.Stroke()

	return ebiten.NewImageFromImage(dc.Image())
}

// NewText creates a new Text
func NewText(parent Container, input *input.Input, x, y, width, height int, text string) *Text {
	l := &Text{
		WidgetBase: WidgetBase{parent: parent, input: input, img: nil, x: x, y: y, width: width, height: height},
		text:       text}
	l.Activate()
	return l
}

// Activate tells the input we need notifications
func (t *Text) Activate() {
	t.disabled = false
	t.img = t.createImg()
	t.input.Add(t)
}

// Deactivate tells the input we no longer need notofications
func (t *Text) Deactivate() {
	t.disabled = true
	t.img = t.createImg()
	t.input.Remove(t)
}

// NotifyCallback is called by the Subject (Input/Stroke) when something interesting happens
func (t *Text) NotifyCallback(event interface{}) {
}

// Update the state of this widget
func (t *Text) Update() {
}
