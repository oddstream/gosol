package ui

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/util"
)

type BarBase struct {
	img           *ebiten.Image
	input         *input.Input
	widgets       []Widget
	x, y          int
	width, height int
}

func (bb *BarBase) createImg() *ebiten.Image {
	dc := gg.NewContext(bb.width, bb.height)
	dc.SetColor(color.RGBA{R: 0x32, G: 0x32, B: 0x32, A: 0xff})
	dc.DrawRectangle(0, 0, float64(bb.width), float64(bb.height))
	dc.Fill()
	dc.Stroke()
	return ebiten.NewImageFromImage(dc.Image())
}

// Position gives the screen position of this container
func (bb *BarBase) Position() (x, y int) {
	x = bb.x
	y = bb.y
	return // using named parameters
}

// Size gives the size of the container
func (bb *BarBase) Size() (width, height int) {
	width = bb.width
	height = bb.height
	return // using named parameters
}

// Rect gives the screen position and extent of this container
func (bb *BarBase) Rect() (x0, y0, x1, y1 int) {
	x0 = bb.x
	y0 = bb.y
	x1 = bb.x + bb.width
	y1 = bb.y + bb.height
	return // using named parameters
}

// FindWidgetAt given screen coordinates
func (bb *BarBase) FindWidgetAt(x, y int) Widget {
	for _, w := range bb.widgets {
		if util.InRect(x, y, w.OffsetRect) {
			return w
		}
	}
	return nil
}

// LayoutWidgets that belong to this container
// by setting the x,y of each relative to their parent
func (bb *BarBase) LayoutWidgets() {
	nextLeft := 0
	nextRight := bb.width - 48
	for _, w := range bb.widgets {
		widgetWidth, widgetHeight := w.Size()
		switch w.Align() {
		case -1: // left align
			w.SetPosition(nextLeft, bb.y)
			nextLeft += widgetWidth + 24 // add padding for big fingers
		case 0: // center
			w.SetPosition(bb.width/2-widgetWidth/2, bb.y+widgetHeight/2)
		case 1: // right align
			w.SetPosition(nextRight, bb.y)
			nextRight -= widgetWidth + 24 // add padding for big fingers
		}
	}
}

// ReplaceWidget replaces a widget
// func (bb *BarBase) ReplaceWidget(n int, w Widget) {
// 	bb.widgets[n].Deactivate()
// 	bb.widgets[n] = w
// }

// StartDrag this widget, if it is allowed
func (bb *BarBase) StartDrag() bool {
	return false // no widget dragging in bars
}

// DragBy this widget
func (bb *BarBase) DragBy(dx, dy int) {
}

// StopDrag this widget
func (bb *BarBase) StopDrag() {
}

// DeactivateWidgets stops the widgets from receiving input
func (bb *BarBase) DeactivateWidgets() {
	for _, w := range bb.widgets {
		bb.input.Remove(w)
	}
}

// Show the bar
func (bb *BarBase) Show() {
}

// Hide the bar
func (bb *BarBase) Hide() {
}

// Visible is the bar
func (bb *BarBase) Visible() bool {
	return true
}

// Update the toolbar
func (bb *BarBase) Update() {
	for _, w := range bb.widgets {
		w.Update()
	}
}

// Draw the bar
func (bb *BarBase) Draw(screen *ebiten.Image) {
	w, _ := screen.Size()
	if bb.img == nil || w != bb.width {
		bb.width = w
		bb.img = bb.createImg()
		bb.LayoutWidgets()
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(bb.x), float64(bb.y))
	screen.DrawImage(bb.img, op)

	for _, w := range bb.widgets {
		w.Draw(screen)
	}
}
