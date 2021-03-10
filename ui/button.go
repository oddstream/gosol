package ui

import (
	"github.com/fogleman/gg"
	"oddstream.games/gosol/schriftbank"
)

// RuneButton is a button that displays a single rune
type RuneButton struct {
	r             rune
	action        func()
	align         int
	x, y          int // screen position
	width, height int
}

// NewRuneButton creates a new RuneButton
func NewRuneButton(r rune, align int, action func()) *RuneButton {
	return &RuneButton{r: r, action: action, align: align, width: 48, height: 48}
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

// Action invokes the action func
func (rb *RuneButton) Action() {
	if rb.action != nil {
		rb.action()
	}
}
