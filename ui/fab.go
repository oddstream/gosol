package ui

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/util"
)

type FAB struct {
	WidgetBase
	img           *ebiten.Image
	x, y          int // position relative to parent
	width, height int
	iconName      string
	key           ebiten.Key
}

func (f *FAB) createImg() *ebiten.Image {
	dc := gg.NewContext(f.width, f.height)
	dc.SetColor(color.RGBA{R: 0x64, G: 0x95, B: 0xed, A: 0xff}) // CornflowerBlue
	dc.DrawCircle(float64(f.width/2), float64(f.height/2), float64(f.height/2))
	dc.Fill()
	dc.Stroke()
	dc.SetRGBA(1, 1, 1, 1)
	dc.DrawImageAnchored(IconMap[f.iconName], f.width/2, f.height/2, 0.5, 0.5)
	return ebiten.NewImageFromImage(dc.Image())
}

func NewFAB(parent Container, iconName string, key ebiten.Key) *FAB {
	// conceptually, a FAB is a toolbar button widget
	f := &FAB{WidgetBase: WidgetBase{parent: parent}, width: 72, height: 72, iconName: iconName, key: key}
	// x, y will be set by Draw()
	f.img = f.createImg()
	return f
}

// NotifyCallback is called by the Subject (Input/Stroke) when something interesting happens
func (f *FAB) NotifyCallback(event interface{}) {
	switch v := event.(type) { // Type switch https://tour.golang.org/methods/16
	case image.Point:
		if util.InRect(v.X, v.Y, f.Rect) {
			// println("FAB notify", f.key)
			f.parent.Notify(f.key)
		}
	}
}

// Activate this widget
func (f *FAB) Activate() {
	f.disabled = false
	f.img = f.createImg()
}

// Deactivate this widget
func (f *FAB) Deactivate() {
	f.disabled = true
	f.img = f.createImg()
}

func (f *FAB) Update() {
}

func (f *FAB) Draw(screen *ebiten.Image) {
	w, h := screen.Size()
	f.x = w - f.width - (f.width / 2)
	f.y = h - f.height - (f.height / 2) - 24 // statusbar is 24 high
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(f.x), float64(f.y))
	screen.DrawImage(f.img, op)
}

func (u *UI) ShowFAB(iconName string, key ebiten.Key) {
	if u.fab == nil {
		u.fab = NewFAB(u.toolbar, iconName, key)
		u.fab.Activate()
	}
}

func (u *UI) HideFAB() {
	if u.fab == nil {
		return
	}
	u.fab.Deactivate()
	u.fab = nil
}
