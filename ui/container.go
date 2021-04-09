package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
)

// Container contains a list of widgets
type Container interface {
	Position() (int, int)
	Size() (int, int)
	Rect() (int, int, int, int)
	FindWidgetAt(int, int) Widget
	LayoutWidgets()
	StartDrag(*input.Stroke) bool
	DragBy(int, int)
	StopDrag()
	Notify(interface{})
	Visible() bool
	Show()
	Hide()
	Layout(int, int) (int, int)
	Update()
	Draw(*ebiten.Image)
}
