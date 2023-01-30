package ui

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/util"
)

type WindowBase struct {
	img           *ebiten.Image
	widgets       []Widgety
	x, y          int
	width, height int
}

// createImg creates the background grey image for this window
func (wb *WindowBase) createImg(bkgrCol color.Color) *ebiten.Image {
	if wb.width == 0 || wb.height == 0 {
		return nil
	}
	dc := gg.NewContext(wb.width, wb.height)
	dc.SetColor(bkgrCol)
	dc.DrawRectangle(0, 0, float64(wb.width), float64(wb.height))
	dc.Fill()
	return ebiten.NewImageFromImage(dc.Image())
}

// Position gives the screen position of the window
func (wb *WindowBase) Position() (x, y int) {
	x = wb.x
	y = wb.y
	return // using named parameters
}

// Size gives the size of the window
func (wb *WindowBase) Size() (width, height int) {
	width = wb.width
	height = wb.height
	return // using named parameters
}

// Rect gives the screen position of the window
func (wb *WindowBase) Rect() (x0, y0, x1, y1 int) {
	x0 = wb.x
	y0 = wb.y
	x1 = wb.x + wb.width
	y1 = wb.y + wb.height
	return // using named parameters
}

func (wb WindowBase) Widgets() []Widgety {
	return wb.widgets
}

func (wb *WindowBase) FindWidgetAt(x, y int) Widgety {
	for _, w := range wb.widgets {
		if util.InRect(x, y, w.OffsetRect) {
			return w
		}
	}
	return nil
}

// Draw the window background and widgets
func (wb *WindowBase) Draw(screen *ebiten.Image) {
	if wb.img == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(wb.x), float64(wb.y))
	screen.DrawImage(wb.img, op)

	for _, w := range wb.widgets {
		w.Draw(screen)
	}
}
