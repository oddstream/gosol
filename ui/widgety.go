package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Widget is an interface for widget objects
type Widgety interface {
	Parent() Containery
	ID() string
	Size() (int, int)
	Position() (int, int)
	Rect() (int, int, int, int)
	OffsetRect() (int, int, int, int)
	SetPosition(int, int)
	Align() int
	Disabled() bool
	Activate()
	Deactivate()
	Tapped()
	Update()
	Draw(*ebiten.Image)
}
