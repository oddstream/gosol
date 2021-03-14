package ui

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/util"
)

type ContainerBase struct {
	img              *ebiten.Image
	input            *input.Input
	widgets          []Widget
	x, y             int
	width, height    int
	xOffset, yOffset int // used when dragging group of widgets
	xOffsetBase      int // used when dragging group of widgets more than once
	yOffsetBase      int // used when dragging group of widgets more than once
}

func (cb *ContainerBase) createImg() *ebiten.Image {
	dc := gg.NewContext(cb.width, cb.height)
	dc.SetColor(color.RGBA{R: 0x32, G: 0x32, B: 0x32, A: 0xff})
	dc.DrawRectangle(0, 0, float64(cb.width), float64(cb.height))
	dc.Fill()
	dc.Stroke()
	return ebiten.NewImageFromImage(dc.Image())
}

// Position gives the screen position
func (cb *ContainerBase) Position() (x, y int) {
	x = cb.x
	y = cb.y
	return // using named parameters
}

// Size gives the size of the container
func (cb *ContainerBase) Size() (width, height int) {
	width = cb.width
	height = cb.height
	return // using named parameters
}

// Rect gives the screen position
func (cb *ContainerBase) Rect() (x0, y0, x1, y1 int) {
	x0 = cb.x
	y0 = cb.y
	x1 = cb.x + cb.width
	y1 = cb.y + cb.height
	return // using named parameters
}

func (cb *ContainerBase) FindWidgetAt(x, y int) Widget {
	for _, w := range cb.widgets {
		if util.InRect(x, y, w.Rect) {
			return w
		}
	}
	return nil
}

// StartDrag this widget, if it is allowed
func (cb *ContainerBase) StartDrag() bool {
	return false
}

// DragBy this widget
func (cb *ContainerBase) DragBy(dx, dy int) {
}

// StopDrag this widget
func (cb *ContainerBase) StopDrag() {
}

// DeactivateWidgets stops the widgets from receiving input
func (cb *ContainerBase) DeactivateWidgets() {
	for _, w := range cb.widgets {
		cb.input.Remove(w)
	}
}
