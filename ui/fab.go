package ui

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/util"
)

type FAB struct {
	input         *input.Input
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
	dc.Stroke()
	return ebiten.NewImageFromImage(dc.Image())

}

func NewFAB(input *input.Input, iconName string, key ebiten.Key) *FAB {
	f := &FAB{input: input, width: 72, height: 72, iconName: iconName, key: key}
	// x, y will be set by Draw()
	f.img = f.createImg()
	return f
}

// Size of the widget
func (f *FAB) Size() (int, int) {
	return f.width, f.height
}

// Position of the widget, relative to parent (normally, in this case in screen coords)
func (f *FAB) Position() (int, int) {
	return f.x, f.y
}

// Rect gives the position and extent of widget, relative to parent
func (f *FAB) Rect() (x0, y0, x1, y1 int) {
	x0 = f.x
	y0 = f.y
	x1 = x0 + f.width
	y1 = y0 + f.height
	return // using named parameters
}

// OffsetRect gives the position and extent of widget, relative to parent
func (f *FAB) OffsetRect() (x0, y0, x1, y1 int) {
	x0 = f.x
	y0 = f.y
	x1 = x0 + f.width
	y1 = y0 + f.height
	return // using named parameters
}

// SetPosition of this widget
func (f *FAB) SetPosition(x, y int) {
	f.x, f.y = x, y
}

// Align returns the x axis alignment (-1, 0, 1)
func (f *FAB) Align() int {
	return 0
}

// NotifyCallback is called by the Subject (Input/Stroke) when something interesting happens
func (f *FAB) NotifyCallback(event interface{}) {
	switch v := event.(type) { // Type switch https://tour.golang.org/methods/16
	case image.Point:
		if util.InRect(v.X, v.Y, f.Rect) {
			// println("FAB notify", f.key)
			f.input.Notify(f.key)
		}
	}
}

// Activate tells the input we need notifications
func (f *FAB) Activate() {
	f.input.Add(f)
}

// Deactivate tells the input we no longer need notifications
func (f *FAB) Deactivate() {
	f.input.Remove(f)
}

func (f *FAB) Update() {
}

func (f *FAB) Draw(screen *ebiten.Image) {
	w, h := screen.Size()
	f.x = w - f.width - (f.width / 2)
	f.y = h - f.height - (f.height / 2)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(f.x), float64(f.y))
	screen.DrawImage(f.img, op)
}

func (u *UI) ShowFAB(iconName string, key ebiten.Key) {
	if u.fab == nil {
		u.fab = NewFAB(u.input, iconName, key)
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
