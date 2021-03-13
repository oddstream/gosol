package ui

import "github.com/hajimehoshi/ebiten/v2"

// Widget is an interface for widget objects
type Widget interface {
	Size() (int, int)
	Position() (int, int)
	Rect() (int, int, int, int)
	OffsetRect() (int, int, int, int)
	SetPosition(int, int)
	Align() int
	Activate()
	Deactivate()
	Update()
	Draw(*ebiten.Image)
	NotifyCallback(interface{})
}

// Container contains a list of widgets
type Container interface {
	Rect() (int, int, int, int)
	FindWidgetAt(int, int) Widget
	LayoutWidgets()
	StartDrag() bool
	DragBy(int, int)
	StopDrag()
	Update()
	Draw(*ebiten.Image)
}
