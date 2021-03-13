package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
)

// WidgetBase is a button that displays a single rune
type WidgetBase struct {
	parent        Container
	input         *input.Input
	img           *ebiten.Image
	align         int
	disabled      bool
	x, y          int // screen position
	width, height int
}

// NewWidgetBase creates a new WidgetBase
func NewWidgetBase(parent Container, input *input.Input, width, height, x, y, align int) *WidgetBase {
	wb := &WidgetBase{parent: parent, input: input, width: width, height: height, x: x, y: y, align: align}
	return wb
}

// Size of the widget
func (wb *WidgetBase) Size() (int, int) {
	return wb.width, wb.height
}

// Position of the widget
func (wb *WidgetBase) Position() (int, int) {
	return wb.x, wb.y
}

// Rect gives the screen position
func (wb *WidgetBase) Rect() (x0, y0, x1, y1 int) {
	x0 = wb.x
	y0 = wb.y
	x1 = wb.x + wb.width
	y1 = wb.y + wb.height
	return // using named parameters
}

// OffsetRect gives the screen position in relation to parent's position
func (wb *WidgetBase) OffsetRect() (x0, y0, x1, y1 int) {
	px, py, _, _ := wb.parent.Rect()
	x0 = px + wb.x
	y0 = py + wb.y
	x1 = px + wb.x + wb.width
	y1 = py + wb.y + wb.height
	// println(x0, y0, x1, y1)
	return // using named parameters
}

// SetPosition of this widget
func (wb *WidgetBase) SetPosition(x, y int) {
	wb.x, wb.y = x, y
}

// Align returns the x axis alignment (-1, 0, 1)
func (wb *WidgetBase) Align() int {
	return wb.align
}

// StartDrag this widget, if it is allowed
func (wb *WidgetBase) StartDrag() bool {
	return false
}

// DragBy this widget
func (wb *WidgetBase) DragBy(dx, dy int) {
}

// StopDrag this widget
func (wb *WidgetBase) StopDrag() {
}

// Draw the widget
func (wb *WidgetBase) Draw(screen *ebiten.Image) {
	// don't draw a widget unless it is fully contained within it's parent
	_, py0, _, py1 := wb.parent.Rect()
	_, _, _, wy1 := wb.Rect()
	if wy1 > py1 || wy1 < py0 {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(wb.x), float64(wb.y))
	screen.DrawImage(wb.img, op)
}
