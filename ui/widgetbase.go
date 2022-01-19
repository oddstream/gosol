package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/util"
)

// WidgetBase is a button that displays a single rune
type WidgetBase struct {
	parent        Container
	img           *ebiten.Image
	align         int
	disabled      bool
	x, y          int // position relative to parent
	width, height int
}

// NewWidgetBase creates a new WidgetBase
func NewWidgetBase(parent Container, width, height, x, y, align int) *WidgetBase {
	wb := &WidgetBase{parent: parent, width: width, height: height, x: x, y: y, align: align}
	return wb
}

// Parent of this widget
func (wb *WidgetBase) Parent() Container {
	return wb.parent
}

// Size of the widget
func (wb *WidgetBase) Size() (int, int) {
	return wb.width, wb.height
}

// Position of the widget, relative to parent
func (wb *WidgetBase) Position() (int, int) {
	return wb.x, wb.y
}

// Rect gives the position and extent of widget, relative to parent
func (wb *WidgetBase) Rect() (x0, y0, x1, y1 int) {
	x0 = wb.x
	y0 = wb.y
	x1 = x0 + wb.width
	y1 = y0 + wb.height
	return // using named parameters
}

// OffsetRect gives the screen position in relation to parent's position
func (wb *WidgetBase) OffsetRect() (x0, y0, x1, y1 int) {
	px, py := wb.parent.Position()
	x0 = px + wb.x
	y0 = py + wb.y
	x1 = x0 + wb.width
	y1 = y0 + wb.height
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

// Disabled returns the disabled status
func (wb *WidgetBase) Disabled() bool {
	return wb.disabled
}

// Update the state of this widget
func (wb *WidgetBase) Update() {
}

// Draw the widget
func (wb *WidgetBase) Draw(screen *ebiten.Image) {

	if wb.img == nil {
		return
	}

	// don't draw a widget unless it is fully contained within it's parent
	parentLeft, parentTop, _, parentBottom := wb.parent.Rect()
	_, _, _, widgetBottom := wb.OffsetRect()
	_, widgetHeight := wb.Size()
	if widgetBottom > parentBottom || widgetBottom-widgetHeight < parentTop {
		return
	}
	// parentLeft, parentTop, _, _ := wb.parent.Rect()

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(parentLeft+wb.x), float64(parentTop+wb.y))
	if wb.disabled {
		op.ColorM.Scale(1, 1, 1, 0.1)
	} else if x, y := ebiten.CursorPosition(); util.InRect(x, y, wb.OffsetRect) {
		// op.ColorM.Scale(100.0/255.0, 149.0/255.0, 237.0/255.0, 1) // CornflowerBlue
		op.ColorM.Scale(1, 1, 1, 0.5)
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			op.GeoM.Translate(2, 2)
		}
	}
	screen.DrawImage(wb.img, op)
}
