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
}

func (cb *ContainerBase) createImg() *ebiten.Image {
	dc := gg.NewContext(cb.width, cb.height)
	dc.SetColor(color.RGBA{R: 0x32, G: 0x32, B: 0x32, A: 0xff})
	dc.DrawRectangle(0, 0, float64(cb.width), float64(cb.height))
	dc.Fill()
	dc.Stroke()
	return ebiten.NewImageFromImage(dc.Image())
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

func (cb *ContainerBase) OffsetWidgets(dx, dy int) {
	// TODO limit based on height of all widgets
	x0, y0, x1, y1 := cb.Rect()
	width := x1 - x0
	height := y1 - y0
	cb.xOffset += dx
	if cb.xOffset > width {
		cb.xOffset = width
	}
	cb.yOffset += dy
	if cb.yOffset > height {
		cb.yOffset = height
	}
}
