package ui

import (
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
		dc.DrawImage(img, 0, n.height/4)
	}
	dc.SetRGBA(1, 1, 1, 1)
	dc.SetFontFace(schriftbank.RobotoMedium24)
	dc.DrawString(n.text, float64(48), float64(n.height)*0.8)

	// uncomment this to show the area we expect the text to occupy
	// dc.DrawLine(0, float64(0), float64(n.width), float64(0))
	// dc.DrawLine(0, float64(n.height), float64(n.width), float64(n.height))
	// dc.DrawLine(0, float64(0), float64(n.width), float64(n.height))
	// dc.Stroke()

	return ebiten.NewImageFromImage(dc.Image())
}

// NewNavItem creates a new NavItem
func NewNavItem(parent Container, iconName string, text string, key ebiten.Key) *NavItem {
	w, _ := parent.Size()
	n := &NavItem{WidgetBase: WidgetBase{parent: parent, img: nil, x: -w, y: 0, width: w, height: 48, align: 0},
		iconName: iconName, text: text, key: key}
	return n
}

// Activate tells the input we need notifications
func (n *NavItem) Activate() {
	n.disabled = false
	n.img = n.createImg() // incase disabled flag has changed
}

// Deactivate tells the input we no longer need notifications
func (n *NavItem) Deactivate() {
	n.disabled = true
	n.img = n.createImg() // incase disabled flag has changed
}

// NotifyCallback is called by the Subject (Input/Stroke) when something interesting happens
func (n *NavItem) NotifyCallback(v input.StrokeEvent) {
	if n.disabled {
		return
	}
	switch v.Event {
	case input.Tap:
		if util.InRect(v.X, v.Y, n.OffsetRect) {
			cmdFn(n.key)
		}
	}
}
