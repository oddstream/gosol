package ui

import (
	"github.com/fogleman/gg"
)

// Widget is an interface for widget objects
type Widget interface {
	// Size() (int, int)
	Rect() (int, int, int, int)
	Align() int
	Draw(*gg.Context, int, int)
	NotifyCallback(interface{})
	Activate()
	Deactivate()
}

// Container contains a list of widgets
type Container interface {
	Rect() (int, int, int, int)
}
