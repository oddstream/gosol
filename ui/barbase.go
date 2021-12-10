package ui

import (
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gomps5/input"
	"oddstream.games/gomps5/util"
)

type BarBase struct {
	img           *ebiten.Image
	stroke        *input.Stroke
	widgets       []Widget
	x, y          int
	width, height int
}

func (bb *BarBase) createImg() *ebiten.Image {
	// println("BarBase createImg", bb.x, bb.y, bb.width, bb.height)
	dc := gg.NewContext(bb.width, bb.height)
	dc.SetColor(BackgroundColor)
	dc.DrawRectangle(0, 0, float64(bb.width), float64(bb.height))
	dc.Fill()
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
	nextLeft := 24
	nextRight := bb.width - 24
	for _, w := range bb.widgets {
		widgetWidth, widgetHeight := w.Size()
		_, parentHeight := w.Parent().Size()
		var y int = parentHeight/2 - widgetHeight/2
		switch w.Align() {
		case -1: // left align
			w.SetPosition(nextLeft, y)
			nextLeft += widgetWidth + 24 // add padding for big fingers
		case 0: // center
			w.SetPosition(bb.width/2-widgetWidth/2, y)
		case 1: // right align
			w.SetPosition(nextRight-widgetWidth, y)
			nextRight -= widgetWidth + 24 // add padding for big fingers
		}
	}
}

// ReplaceWidget replaces a widget
// func (bb *BarBase) ReplaceWidget(n int, w Widget) {
// 	bb.widgets[n].Deactivate()
// 	bb.widgets[n] = w
// }

// StartDrag this container, if it is allowed
func (bb *BarBase) StartDrag(stroke *input.Stroke) bool {
	// println("BarBase start drag, adding widgets")
	bb.stroke = stroke
	for _, w := range bb.widgets {
		if !w.Disabled() {
			stroke.Add(w)
		}
	}
	return true
}

// DragBy this widget
func (bb *BarBase) DragBy(dx, dy int) {
	// you can't drag a bar
}

// StopDrag this widget
func (bb *BarBase) StopDrag() {
	// println("BarBase stop drag, removing widgets")
	for _, w := range bb.widgets {
		if !w.Disabled() {
			bb.stroke.Remove(w)
		}
	}
	bb.stroke = nil
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

// Layout implements Ebiten's Layout
func (bb *BarBase) Layout(outsideWidth, outsideHeight int) (int, int) {
	if bb.img == nil || outsideWidth != bb.width {
		bb.width = outsideWidth
		bb.img = bb.createImg()
		bb.LayoutWidgets()
	}
	return outsideWidth, outsideHeight
}

// Update the bar
func (bb *BarBase) Update() {
	for _, w := range bb.widgets {
		w.Update()
	}
}

// Draw the bar
func (bb *BarBase) Draw(screen *ebiten.Image) {
	if bb.img == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(bb.x), float64(bb.y))
	screen.DrawImage(bb.img, op)

	for _, w := range bb.widgets {
		w.Draw(screen)
	}
}
