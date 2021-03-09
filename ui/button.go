package ui

import (
	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

// RuneButton is a button that displays a single rune
type RuneButton struct {
	r             rune
	face          font.Face
	action        func()
	align         int
	width, height int
}

// NewRuneButton creates a new RuneButton
func NewRuneButton(r rune, face font.Face, action func(), align int) *RuneButton {
	return &RuneButton{r: r, face: face, action: action, align: align, width: 48, height: 48}
}

// Size of the RuneButton
func (rb *RuneButton) Size() (int, int) {
	return rb.width, rb.height
}

// Align returns the x axis alignment (-1, 0, 1)
func (rb *RuneButton) Align() int {
	return rb.align
}

// Draw into a gg context, not to the screen
func (rb *RuneButton) Draw(dc *gg.Context, x, y int) {
	dc.SetFontFace(rb.face)
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

// NotifyCallback is called by the Subject (Input/Stroke) when something interesting happens
func (rb *RuneButton) NotifyCallback(event interface{}) {
	// switch v := event.(type) { // Type switch https://tour.golang.org/methods/16
	// case image.Point:
	// case ebiten.Key:
	// }
}
