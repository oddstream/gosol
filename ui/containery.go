package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Container is an interface for a UI widget
type Containery interface {
	Position() (int, int)
	Size() (int, int)
	Rect() (int, int, int, int)
	Widgets() []Widgety
	FindWidgetAt(int, int) Widgety
	LayoutWidgets()
	StartDrag()
	DragBy(int, int)
	StopDrag()
	CancelDrag()
	Tapped()
	Visible() bool
	Show()
	Hide()
	Layout(int, int) (int, int)
	Update()
	Draw(*ebiten.Image)
}
