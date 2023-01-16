package ui

import (
	"log"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
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
func NewIconButton(parent Containery, id string, x, y, width, height, align int, iconName string, key ebiten.Key) *IconButton {
	b := &IconButton{WidgetBase: WidgetBase{parent: parent, id: id, img: nil, x: x, y: y, width: width, height: height, align: align},
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

func (b *IconButton) Tapped() {
	if b.disabled {
		return
	}
	cmdFn(b.key)
}
