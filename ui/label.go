package ui

import (
	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

// Label is a button that displays a single rune
type Label struct {
	text          string
	face          font.Face
	align         int // -1 left, 0=center, 1=right
	x, y          int // screen position
	width, height int
}

// NewLabel creates a new Label
func NewLabel(text string, face font.Face, align int) *Label {
	return &Label{text: text, face: face, align: align, width: 48, height: 48}
}

// Size of the Label
func (l *Label) Size() (int, int) {
	return l.width, l.height
}

// Rect gives the screen position
func (l *Label) Rect() (x0, y0, x1, y1 int) {
	x0 = l.x
	y0 = l.y
	x1 = l.x + l.width
	y1 = l.y + l.height
	return // using named parameters
}

// Align returns the x axis alignment (-1, 0, 1)
func (l *Label) Align() int {
	return l.align
}

// Draw into a gg context, not to the screen
func (l *Label) Draw(dc *gg.Context, x, y int) {
	dc.SetFontFace(l.face)
	dc.SetRGBA(1, 1, 1, 1)
	dc.DrawStringAnchored(l.text, float64(x), float64(y), 0.5, 0.5)
	dc.Stroke()

	l.x, l.y = x, y
	w, h := dc.MeasureString(l.text)
	l.width, l.height = int(w), int(h)
}

// Action invokes the action func
func (l *Label) Action() {
}

// NotifyCallback is called by the Subject (Input) when something interesting happens
func (l *Label) NotifyCallback(event interface{}) {
}
