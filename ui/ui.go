package ui

import (
	_ "embed" // go:embed only allowed in Go files that import "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/util"
)

// UI encapsulates a complete user interface that can be rendered onto the screen.
type UI struct {
	input        *input.Input // place to receive clicks, taps and key presses from
	toastManager *ToastManager
	toolbar      *Toolbar
	navdrawer    *NavDrawer
	picker       *Picker
}

// New creates a new UI object
func New(input *input.Input, pickerContents []string) *UI {
	ui := &UI{input: input}

	ui.toastManager = &ToastManager{}
	ui.toolbar = NewToolbar(input)
	ui.navdrawer = NewNavDrawer(input)
	ui.picker = NewPicker(input, pickerContents)

	return ui
}

// FindWidgetAt finds the widget at the screen coords
func (u *UI) FindWidgetAt(x, y int) Widget {

	var w Widget
	w = u.toolbar.FindWidgetAt(x, y)
	if w != nil {
		return w
	}
	w = u.navdrawer.FindWidgetAt(x, y)
	if w != nil {
		return w
	}
	w = u.picker.FindWidgetAt(x, y)
	if w != nil {
		return w
	}
	return nil
}

func (u *UI) FindContainerAt(x, y int) Container {

	if util.InRect(x, y, u.toolbar.Rect) {
		return u.toolbar
	}
	if util.InRect(x, y, u.navdrawer.Rect) {
		return u.navdrawer
	}
	if util.InRect(x, y, u.picker.Rect) {
		return u.picker
	}
	return nil

}

// VisibleDrawer returns the drawer that is currently open
func (u *UI) VisibleDrawer() Container {
	if u.navdrawer.Visible() {
		return u.navdrawer
	}
	if u.picker.Visible() {
		return u.picker
	}
	return nil
}

// ActiveRect returns the rect coords of the active drawer
func (u *UI) ActiveDrawerRect() (int, int, int, int) {
	d := u.VisibleDrawer()
	if d != nil {
		return d.Rect()
	}
	return 0, 0, 0, 0
}

// HideActiveDrawer closes the active/open drawer
func (u *UI) HideActiveDrawer() {
	if con := u.VisibleDrawer(); con != nil {
		con.Hide()
	}
}

// Update is called once per tick and updates the UI's state
func (u *UI) Update() {
	u.toolbar.Update()
	u.navdrawer.Update()
	u.picker.Update()
	u.toastManager.Update()
}

// Draw is called once per tick and renders the UI to the screen
func (u *UI) Draw(screen *ebiten.Image) {
	u.toolbar.Draw(screen)
	u.navdrawer.Draw(screen)
	u.picker.Draw(screen)
	u.toastManager.Draw(screen)
}
