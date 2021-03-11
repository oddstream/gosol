package ui

import (
	"image"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/schriftbank"
	"oddstream.games/gosol/util"
)

// NavItem is a button that displays a single rune
type NavItem struct {
	parent        Container
	r             rune
	text          string
	x, y          int // screen position
	width, height int // always 256 wide (same as navdrawer width) 48 high (same as text)
	input         *input.Input
	key           ebiten.Key
}

// NewNavItem creates a new NavItem
func NewNavItem(parent Container, r rune, text string, input *input.Input, key ebiten.Key) *NavItem {
	n := &NavItem{parent: parent, r: r, text: text, width: 256, height: 48, input: input, key: key}
	n.Activate()
	return n
}

// Activate tells the input we need notifications
func (n *NavItem) Activate() {
	n.input.Add(n)
}

// Deactivate tells the input we no longer need notifications
func (n *NavItem) Deactivate() {
	n.input.Remove(n)
}

// Rect gives the screen position
func (n *NavItem) Rect() (x0, y0, x1, y1 int) {
	x0 = n.x
	y0 = n.y
	x1 = n.x + n.width
	y1 = n.y + n.height
	return // using named parameters
}

// OffsetRect gives the screen position in relation to parent's position
func (n *NavItem) OffsetRect() (x0, y0, x1, y1 int) {
	px, py, _, _ := n.parent.Rect()
	x0 = px + n.x
	y0 = py + n.y
	x1 = px + n.x + n.width
	y1 = py + n.y + n.height
	// println(x0, y0, x1, y1)
	return // using named parameters
}

// NotifyCallback is called by the Subject (Input/Stroke) when something interesting happens
func (n *NavItem) NotifyCallback(event interface{}) {
	switch v := event.(type) { // Type switch https://tour.golang.org/methods/16
	case image.Point:
		// println("NavItem image.Point", v.X, v.Y)
		if util.InRect(v.X, v.Y, n.OffsetRect) {
			n.input.Notify(n.key)
		}
	}
}

// Align returns the x axis alignment
func (n *NavItem) Align() int {
	return 0 // not implemented
}

// Draw into a gg context, not to the screen
func (n *NavItem) Draw(dc *gg.Context, x, y int) {

	dc.SetRGBA(1, 1, 1, 1)

	// nota bene - text is drawn with y as a baseline

	dc.SetFontFace(schriftbank.Symbol24)
	dc.DrawString(string(n.r), 24, float64(y))
	dc.SetFontFace(schriftbank.RobotoRegular24)
	dc.DrawString(n.text, float64(24+48), float64(y))

	// uncomment this to show the area we expect the text to occupy
	// dc.DrawLine(float64(x), float64(y-32), float64(n.width), float64(y-32))
	// dc.DrawLine(float64(x), float64(y+16), float64(n.width), float64(y+16))
	// dc.DrawLine(float64(x), float64(y+16), float64(n.width), float64(y-32))

	dc.Stroke()

	n.x, n.y = x, y-32
}
