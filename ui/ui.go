package ui

import (
	_ "embed" // go:embed only allowed in Go files that import "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
)

// UI encapsulates a complete user interface that can be rendered onto the screen.
type UI struct {
	input        *input.Input // place to receive clicks, taps and key presses from
	toastManager *ToastManager
	toolbar      *Toolbar
	navdrawer    *NavDrawer
	window       *Window
}

// New creates a new UI object
func New(input *input.Input) *UI {
	ui := &UI{input: input}

	ui.toastManager = &ToastManager{}
	ui.toolbar = NewToolbar(input)
	ui.navdrawer = NewNavDrawer(input)
	ui.window = nil

	return ui
}

// NotifyCallback is called by the Subject (Input) when something interesting happens
// func (u *UI) NotifyCallback(event interface{}) {
// 	switch v := event.(type) { // Type switch https://tour.golang.org/methods/16
// 	case image.Point:
// 		println("UI event", v.X, v.Y)
// 		if util.InRect(v.X, v.Y, u.navdrawer.Rect) {
// 			println("UI click over navdrawer")
// 			u.navdrawer.Tapped(v.X, v.Y)
// 		} else if util.InRect(v.X, v.Y, u.toolbar.Rect) {
// 			println("UI click over toolbar")
// 			u.toolbar.Tapped(v.X, v.Y)
// 		}
// 	}
// }

// ActiveModal returns true if there is an active modal
func (u *UI) ActiveModal() bool {
	if u.navdrawer.Visible() {
		return true
	}
	if u.window != nil {
		return true
	}
	return false
}

// ActiveRect returns the rect coords of any active UI object (dialog, drawer)
func (u *UI) ActiveRect() (int, int, int, int) {
	if u.navdrawer.Visible() {
		return u.navdrawer.Rect()
	}
	if u.window != nil {
		return u.window.Rect()
	}
	return 0, 0, 0, 0
}

// CloseActiveModal closes any open dialogs, drawers &c
func (u *UI) CloseActiveModal() {
	if u.navdrawer.Visible() {
		u.navdrawer.Hide()
	}
	if u.window != nil {
		u.window = nil
	}
}

// Update is called once per tick and updates the UI's state
func (u *UI) Update() {
	u.toolbar.Update()
	u.navdrawer.Update()
	u.toastManager.Update()
	if u.window != nil {
		u.window.Update()
	}
}

// Draw is called once per tick and renders the UI to the screen
func (u *UI) Draw(screen *ebiten.Image) {
	u.toolbar.Draw(screen)
	u.navdrawer.Draw(screen)
	u.toastManager.Draw(screen)
	if u.window != nil {
		u.window.Draw(screen)
	}
}
