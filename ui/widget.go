package ui

import (
	"github.com/fogleman/gg"
)

// Widget is an interface for widget objects
type Widget interface {
	Size() (int, int)
	Align() int
	Draw(*gg.Context, int, int)
	Action()
	NotifyCallback(interface{})
}
