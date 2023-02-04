package ui

import (
	"oddstream.games/gosol/sound"
	"oddstream.games/gosol/util"
)

const (
	aniLeft  = -1
	aniStop  = 0
	aniRight = 1
)

type DrawerBase struct {
	WindowBase
	aniState         int
	xOffset, yOffset int // used when dragging group of widgets
	xOffsetBase      int // used when dragging group of widgets more than once
	yOffsetBase      int // used when dragging group of widgets more than once
}

// LayoutWidgets that belong to this container
// by setting the x,y of each relative to their parent
func (db *DrawerBase) LayoutWidgets() {
	const padding = 24
	var x, y int
	x = padding
	y = padding
	for _, w := range db.widgets {
		w.SetPosition(x, y+db.yOffset)
		_, widgetHeight := w.Size()
		y += widgetHeight + padding
	}
}

// Show starts to animate the drawer on screen from the left
func (db *DrawerBase) Show() {
	// for _, w := range db.widgets {
	// 	w.Activate()
	// }
	db.aniState = aniRight
	sound.Play("Click")
}

// Hide starts to animate the drawer off screen to the left
func (db *DrawerBase) Hide() {
	// for _, w := range db.widgets {
	// 	w.Deactivate()
	// }
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

// StartDrag this container, if it is allowed
func (db *DrawerBase) StartDrag() {
}

// DragBy this container
func (db *DrawerBase) DragBy(dx, dy int) {
	db.xOffset = db.xOffsetBase + dx
	db.xOffset = util.ClampInt(db.xOffset, -db.width, 0)

	numWidgets := len(db.widgets)
	_, widgetHeight := db.widgets[0].Size()
	widgetHeight += 24 // see Picker.LayoutWidgets, which uses (widget height) + 24 as vertical padding
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

// StopDrag this container
func (db *DrawerBase) StopDrag() {
	// remember the amount of drag incase the widgets are dragged again
	db.xOffsetBase = db.xOffset
	db.yOffsetBase = db.yOffset
}

func (db *DrawerBase) CancelDrag() {
}

func (db DrawerBase) Tapped() {
}

// ResetScroll state for this drawer
func (db *DrawerBase) ResetScroll() {
	db.xOffset = 0
	db.xOffsetBase = 0
	db.yOffset = 0
	db.yOffsetBase = 0
}

// Layout implements Ebiten's Layout
func (db *DrawerBase) Layout(outsideWidth, outsideHeight int) (int, int) {
	if db.img == nil || db.height != outsideHeight-ToolbarHeight-StatusbarHeight {
		db.height = outsideHeight - ToolbarHeight - StatusbarHeight
		db.img = db.createImg(BackgroundColor)
		db.LayoutWidgets()
	}
	return outsideWidth, outsideHeight
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
