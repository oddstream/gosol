package ui

import (
	"image"
	"log"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/schriftbank"
	"oddstream.games/gosol/util"
)

// NavItem is a button that displays a single rune
type NavItem struct {
	WidgetBase
	iconName string
	text     string
	key      ebiten.Key
}

func (n *NavItem) createImg() *ebiten.Image {

	dc := gg.NewContext(n.width, n.height)

	// nota bene - text is drawn with y as a baseline

	if n.iconName != "" {
		img, ok := IconMap[n.iconName]
		if !ok || img == nil {
			log.Fatal(n.iconName, " not in icon map")
		}
		dc.DrawImageAnchored(img, 18, n.height/2, 0, 0.5)
	}
	dc.SetRGBA(1, 1, 1, 1)
	dc.SetFontFace(schriftbank.RobotRegular24)
	dc.DrawString(n.text, float64(18+48), float64(n.height)*0.7)

	// uncomment this to show the area we expect the text to occupy
	// dc.DrawLine(0, float64(0), float64(n.width), float64(0))
	// dc.DrawLine(0, float64(n.height), float64(n.width), float64(n.height))
	// dc.DrawLine(0, float64(0), float64(n.width), float64(n.height))
	// dc.Stroke()

	return ebiten.NewImageFromImage(dc.Image())
}

// NewNavItem creates a new NavItem
func NewNavItem(parent Container, input *input.Input, iconName string, text string, key ebiten.Key) *NavItem {
	w, _ := parent.Size()
	n := &NavItem{WidgetBase: WidgetBase{parent: parent, input: input, img: nil, x: -w, y: 0, width: w, height: 48, align: 0},
		iconName: iconName, text: text, key: key}
	n.Activate()
	return n
}

// Activate tells the input we need notifications
func (n *NavItem) Activate() {
	n.disabled = false
	n.img = n.createImg() // incase disabled flag has changed
	n.input.Add(n)
}

// Deactivate tells the input we no longer need notifications
func (n *NavItem) Deactivate() {
	n.disabled = true
	n.img = n.createImg() // incase disabled flag has changed
	n.input.Remove(n)
}

// NotifyCallback is called by the Subject (Input/Stroke) when something interesting happens
func (n *NavItem) NotifyCallback(event interface{}) {
	if n.disabled {
		return
	}
	switch v := event.(type) { // Type switch https://tour.golang.org/methods/16
	case image.Point:
		// println("NavItem image.Point", v.X, v.Y)
		if util.InRect(v.X, v.Y, n.OffsetRect) {
			println("NavItem notify", n.key)
			n.input.Notify(n.key)
		}
	}
}
