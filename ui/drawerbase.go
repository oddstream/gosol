package ui

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/util"
)

const (
	aniLeft  = -1
	aniStop  = 0
	aniRight = 1
)

type DrawerBase struct {
	img              *ebiten.Image
	input            *input.Input
	widgets          []Widget
	x, y             int
	width, height    int
	aniState         int
	xOffset, yOffset int // used when dragging group of widgets
	xOffsetBase      int // used when dragging group of widgets more than once
	yOffsetBase      int // used when dragging group of widgets more than once
}

func (db *DrawerBase) createImg() *ebiten.Image {
	dc := gg.NewContext(db.width, db.height)
	dc.SetColor(color.RGBA{R: 0x32, G: 0x32, B: 0x32, A: 0xff})
	// dc.SetColor(color.RGBA{R: 0, G: 0, B: 0, A: 0xff})
	dc.DrawRectangle(0, 0, float64(db.width), float64(db.height))
	dc.Fill()
	dc.Stroke()
	return ebiten.NewImageFromImage(dc.Image())
}

// Position gives the screen position
func (db *DrawerBase) Position() (x, y int) {
	x = db.x
	y = db.y
	return // using named parameters
}

// Size gives the size of the container
func (db *DrawerBase) Size() (width, height int) {
	width = db.width
	height = db.height
	return // using named parameters
}

// Rect gives the screen position
func (db *DrawerBase) Rect() (x0, y0, x1, y1 int) {
	x0 = db.x
	y0 = db.y
	x1 = db.x + db.width
	y1 = db.y + db.height
	return // using named parameters
}

func (db *DrawerBase) FindWidgetAt(x, y int) Widget {
	for _, w := range db.widgets {
		if util.InRect(x, y, w.OffsetRect) {
			return w
		}
	}
	return nil
}

// LayoutWidgets that belong to this container
// by setting the x,y of each relative to their parent
func (db *DrawerBase) LayoutWidgets() {
	const padding = 8
	var x, y int
	y = padding
	for _, w := range db.widgets {
		w.SetPosition(x, y+db.yOffset)
		_, widgetHeight := w.Size()
		y += widgetHeight + padding
	}
}

// Show starts to animate the drawer on screen from the left
func (db *DrawerBase) Show() {
	for _, w := range db.widgets {
		w.Activate()
	}
	db.aniState = aniRight
}

// Hide starts to animate the drawer off screen to the left
func (db *DrawerBase) Hide() {
	for _, w := range db.widgets {
		w.Deactivate()
	}
	if db.x == -db.width {
		db.aniState = aniStop
	} else {
		db.aniState = aniLeft
	}
}

// Visible returns true if the NavDrawer is showing
func (db *DrawerBase) Visible() bool {
	return db.x == 0
}

// StartDrag this widget, if it is allowed
func (db *DrawerBase) StartDrag() bool {
	// println("start drag with offset base", db.yOffsetBase)
	return true
}

// DragBy this widget
func (db *DrawerBase) DragBy(dx, dy int) {
	db.xOffset = db.xOffsetBase + dx
	db.xOffset = util.ClampInt(db.xOffset, -db.width, 0)

	numWidgets := len(db.widgets)
	_, widgetHeight := db.widgets[0].Size()
	widgetHeight += 14 // see Picker.LayoutWidgets, which uses 24 (widget height) + 14 as vertical spacing
	_, pickerHeight := db.Size()
	// println("picker height", pickerHeight, "widget height", widgetHeight)
	visibleWidgets := pickerHeight / widgetHeight
	hiddenWidgets := numWidgets - visibleWidgets
	// println("total", numWidgets, "visible", visibleWidgets, "hidden", hiddenWidgets)
	heightOfHiddenWidgets := hiddenWidgets * widgetHeight
	db.yOffset = db.yOffsetBase + dy
	db.yOffset = util.ClampInt(db.yOffset, -heightOfHiddenWidgets, 0)
	db.LayoutWidgets()
}

// StopDrag this widget
func (db *DrawerBase) StopDrag() {
	// remember the amount of drag incase the widgets are dragged again
	db.xOffsetBase = db.xOffset
	db.yOffsetBase = db.yOffset
}

// ResetScroll state for this drawer
func (db *DrawerBase) ResetScroll() {
	db.xOffset = 0
	db.xOffsetBase = 0
	db.yOffset = 0
	db.yOffsetBase = 0
}

// DeactivateWidgets stops the widgets from receiving input
func (db *DrawerBase) DeactivateWidgets() {
	for _, w := range db.widgets {
		db.input.Remove(w)
	}
}

// Update the Drawer
func (db *DrawerBase) Update() {
	switch db.aniState {
	case aniLeft:
		if db.x <= -db.width {
			db.x = -db.width
			db.aniState = aniStop
		} else {
			db.x -= 16
		}
	case aniRight:
		if db.x >= 0 {
			db.x = 0
			db.aniState = aniStop
		} else {
			db.x += 16
		}
	}
	for _, w := range db.widgets {
		w.Update()
	}
}

// Draw the Drawer
func (db *DrawerBase) Draw(screen *ebiten.Image) {

	const toolbarHeight int = 48 // draw drawer below toolbar

	_, screenHeight := screen.Size()
	if db.img == nil || screenHeight != toolbarHeight+db.height {
		db.height = screenHeight - toolbarHeight
		db.img = db.createImg()
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(db.x), float64(db.y))
	// op.ColorM.Scale(1, 1, 1, 0.95)
	screen.DrawImage(db.img, op)

	for _, w := range db.widgets {
		w.Draw(screen)
	}

}
