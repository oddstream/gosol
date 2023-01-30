package ui

type BarBase struct {
	WindowBase
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

// StartDrag notifies this container that dragging has started
func (bb *BarBase) StartDrag() {
	// println("BarBase start drag, adding widgets")
}

// DragBy this container
func (bb *BarBase) DragBy(dx, dy int) {
	// you can't drag a bar
}

// StopDrag notifies the container that dragging has been stopped
func (bb *BarBase) StopDrag() {
}

// CancelDrag notifies the container that dragging has been cancelled
func (bb *BarBase) CancelDrag() {
}

func (bb *BarBase) Tapped() {
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
		bb.img = bb.createImg(BackgroundColor)
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
