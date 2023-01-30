package ui

import (
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

type FAB struct {
	WidgetBase
	iconName string
	key      ebiten.Key
}

func (f *FAB) createImg() *ebiten.Image {
	// WidgetBase doesn't have a default createImg
	if f.width == 0 || f.height == 0 {
		return nil
	}
	dc := gg.NewContext(f.width, f.height)
	dc.SetColor(BackgroundColor)
	dc.DrawCircle(float64(f.width/2), float64(f.height/2), float64(f.height/2))
	dc.Fill()
	dc.Stroke()
	dc.SetColor(ForegroundColor)
	dc.DrawImageAnchored(IconMap[f.iconName], f.width/2, f.height/2, 0.5, 0.5)
	return ebiten.NewImageFromImage(dc.Image())
}

func NewFAB(parent Containery, id string, iconName string, key ebiten.Key) *FAB {
	f := &FAB{WidgetBase: WidgetBase{parent: parent, id: id, x: 0, y: 0, width: FABWidth, height: FABWidth}, iconName: iconName, key: key}
	f.Activate()
	return f
}

func (f *FAB) Tapped() {
	if f.disabled {
		return
	}
	cmdFn(f.key)
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

//

type FABBar struct {
	BarBase
}

func (fb *FABBar) createImg() *ebiten.Image {
	// override BarBase createImg to create transparent image
	dc := gg.NewContext(fb.width, fb.height)
	dc.SetRGBA(0, 0, 0, 0)
	dc.DrawRectangle(0, 0, float64(fb.width), float64(fb.height))
	dc.Fill()
	return ebiten.NewImageFromImage(dc.Image())
}

func NewFABBar() *FABBar {
	fb := &FABBar{BarBase: BarBase{WindowBase: WindowBase{x: 0, y: 0, width: FABWidth, height: FABWidth}}}
	fb.img = fb.createImg()
	// no widgets yet
	return fb
}

// Layout implements Ebiten's Layout
func (fb *FABBar) Layout(outsideWidth, outsideHeight int) (int, int) {
	// override BarBase.Layout to get position near bottom right of screen
	fb.x = outsideWidth - fb.width - (fb.width / 2)
	fb.y = outsideHeight - fb.height - (fb.height / 2) - StatusbarHeight // statusbar is 24 high
	// println("FABBar.Layout() Window=", outsideWidth, outsideHeight, "FAB=", fb.x, fb.y)
	return outsideWidth, outsideHeight
}

//

func (u *UI) ShowFAB(iconName string, key ebiten.Key) {
	u.fabbar.widgets = nil
	u.fabbar.widgets = append(u.fabbar.widgets, NewFAB(u.fabbar, "", iconName, key))
}

func (u *UI) HideFAB() {
	u.fabbar.widgets = nil
}
