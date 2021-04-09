package ui

import (
	"image"
	"log"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/util"
)

// IconButton is a button that displays a single rune
type IconButton struct {
	WidgetBase
	iconName string
	key      ebiten.Key
}

func (b *IconButton) createImg() *ebiten.Image {
	// println("button createImg", b.iconName)
	dc := gg.NewContext(b.width, b.height)
	img, ok := IconMap[b.iconName]
	if !ok || img == nil {
		log.Println(b.iconName, " not in icon map")
	}
	dc.DrawImageAnchored(img, b.width/2, b.height/2, 0.5, 0.5)
	return ebiten.NewImageFromImage(dc.Image())
}

// NewIconButton creates a new IconButton
func NewIconButton(parent Container, x, y, width, height, align int, iconName string, key ebiten.Key) *IconButton {
	b := &IconButton{WidgetBase: WidgetBase{parent: parent, img: nil, x: x, y: y, width: width, height: height, align: align},
		iconName: iconName, key: key}
	b.Activate()
	return b
}

// Activate this widget
func (b *IconButton) Activate() {
	b.disabled = false
	b.img = b.createImg()
}

// Deactivate this widget
func (b *IconButton) Deactivate() {
	b.disabled = true
	b.img = b.createImg()
}

// NotifyCallback is called by the Subject (Input/Stroke) when something interesting happens
func (b *IconButton) NotifyCallback(event interface{}) {
	println("IconButton NotifyCallback, disabled", b.disabled)
	if b.disabled {
		return
	}
	switch v := event.(type) { // Type switch https://tour.golang.org/methods/16
	case image.Point:
		println("IconButton image.Point", v.X, v.Y)
		if util.InRect(v.X, v.Y, b.OffsetRect) {
			println("IconButton sending notify to parent", b.key)
			b.parent.Notify(b.key)
		}
	case input.StrokeEvent:
		println("IconButton stroke event", v.Event)
	case ebiten.Key:
		println("IconButton key event", v)
	default:
		println("IconButton unknown event type")
	}
}
