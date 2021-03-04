package ui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Widget represents a display/interactable object
type Widget struct {
	Rect     image.Rectangle
	Disabled bool
}

// Container holds Widgets
type Container struct {
	widgets         []*Widget
	BackgroundImage *ebiten.Image
}

// UI encapsulates a complete user interface that can be rendered onto the screen.
// There should only be exactly one UI per application.
type UI struct {
	root *Container
}

// Update is called once per tick and updates the UI's state
func (u *UI) Update() {

}

// Draw is called once per tick and renders the UI to the screen
func (u *UI) Draw(screen *ebiten.Image) {

}
