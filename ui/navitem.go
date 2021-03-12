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
	WidgetBase
	r    rune
	text string
	key  ebiten.Key
}

func (n *NavItem) createImg() *ebiten.Image {

	dc := gg.NewContext(n.width, n.height)

	dc.SetRGBA(1, 1, 1, 1)

	// nota bene - text is drawn with y as a baseline

	dc.SetFontFace(schriftbank.Symbol24)
	dc.DrawString(string(n.r), 24, float64(n.height)*0.7)
	dc.SetFontFace(schriftbank.RobotoRegular24)
	dc.DrawString(n.text, float64(24+48), float64(n.height)*0.7)

	// uncomment this to show the area we expect the text to occupy
	// dc.DrawLine(0, float64(0), float64(n.width), float64(0))
	// dc.DrawLine(0, float64(n.height), float64(n.width), float64(n.height))
	// dc.DrawLine(0, float64(0), float64(n.width), float64(n.height))

	dc.Stroke()

	return ebiten.NewImageFromImage(dc.Image())
}

// NewNavItem creates a new NavItem
func NewNavItem(parent Container, input *input.Input, x, y, width, height, align int, r rune, text string, key ebiten.Key) *NavItem {
	n := &NavItem{WidgetBase: WidgetBase{parent: parent, input: input, img: nil, x: x, y: y, width: width, height: height, align: align}, r: r, text: text, key: key}
	n.img = n.createImg()
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

// Update the state of this widget
func (n *NavItem) Update() {
}
