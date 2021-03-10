package ui

import (
	_ "embed" // go:embed only allowed in Go files that import "embed"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/util"
)

// Container holds Widgets
// type Container struct {
// 	widgets         []*Widget
// 	BackgroundImage *ebiten.Image
// }

// UI encapsulates a complete user interface that can be rendered onto the screen.
type UI struct {
	input        *input.Input   // place to receive clicks, taps and key presses from
	observer     input.Observer // place to send commands to
	toastManager *ToastManager
	toolbar      *Toolbar
	navdrawer    *NavDrawer
}

// New creates a new UI object
func New(i *input.Input, observer input.Observer) *UI {
	ui := &UI{input: i}

	i.Add(ui)

	ui.toastManager = &ToastManager{}
	ui.toolbar = NewToolbar(observer)

	ui.navdrawer = NewNavDrawer()

	return ui
}

// NotifyCallback is called by the Subject (Input) when something interesting happens
func (u *UI) NotifyCallback(event interface{}) {
	switch v := event.(type) { // Type switch https://tour.golang.org/methods/16
	case image.Point:
		println("UI event", v.X, v.Y)
		if util.InRect(v.X, v.Y, u.navdrawer.Rect) {
			println("UI click over navdrawer")
			u.navdrawer.Tapped(v.X, v.Y)
		} else if util.InRect(v.X, v.Y, u.toolbar.Rect) {
			println("UI click over toolbar")
			u.toolbar.Tapped(v.X, v.Y)
		}
	}
}

// ToggleNavDrawer animates the drawer on/off screen to the left
func (u *UI) ToggleNavDrawer() {
	if u.navdrawer.Visible() {
		u.navdrawer.Hide()
	} else {
		u.navdrawer.Show()
	}
}

// Update is called once per tick and updates the UI's state
func (u *UI) Update() {
	u.toolbar.Update()
	u.navdrawer.Update()
	u.toastManager.Update()
}

// Draw is called once per tick and renders the UI to the screen
func (u *UI) Draw(screen *ebiten.Image) {
	u.toolbar.Draw(screen)
	u.navdrawer.Draw(screen)
	u.toastManager.Draw(screen)
}
