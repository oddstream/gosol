package ui

import "github.com/hajimehoshi/ebiten/v2"

// Container contains a list of widgets
type Container interface {
	Position() (int, int)
	Size() (int, int)
	Rect() (int, int, int, int)
	FindWidgetAt(int, int) Widget
	LayoutWidgets()
	DeactivateWidgets()
	StartDrag() bool
	DragBy(int, int)
	StopDrag()
	Visible() bool
	Show()
	Hide()
	Update()
	Draw(*ebiten.Image)
}
