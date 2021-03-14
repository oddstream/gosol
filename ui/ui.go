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
	modal        Container
}

// New creates a new UI object
func New(input *input.Input) *UI {
	ui := &UI{input: input}

	ui.toastManager = &ToastManager{}
	ui.toolbar = NewToolbar(input)
	ui.navdrawer = NewNavDrawer(input)

	return ui
}

// FindWidgetAt finds the widget at the screen coords
func (u *UI) FindWidgetAt(x, y int) Widget {
	if u.modal == nil {
		return nil
	}
	return u.modal.FindWidgetAt(x, y)
}

func (u *UI) FindContainerAt(x, y int) Container {

	if u.modal != nil && util.InRect(x, y, u.modal.Rect) {
		return u.modal
	}
	if util.InRect(x, y, u.toolbar.Rect) {
		return u.toolbar
	}
	if u.IsNavDrawerOpen() && util.InRect(x, y, u.navdrawer.Rect) {
		return u.navdrawer
	}
	return nil

}

// ActiveModal returns true if there is an active modal
func (u *UI) ActiveModal() bool {
	if u.navdrawer.Visible() {
		return true
	}
	if u.modal != nil {
		return true
	}
	return false
}

// ActiveRect returns the rect coords of any active UI object (dialog, drawer)
func (u *UI) ActiveRect() (int, int, int, int) {
	if u.navdrawer.Visible() {
		return u.navdrawer.Rect()
	}
	if u.modal != nil {
		return u.modal.Rect()
	}
	return 0, 0, 0, 0
}

// CloseActiveModal closes any open dialogs, drawers &c
func (u *UI) CloseActiveModal() {
	if u.navdrawer.Visible() {
		u.navdrawer.Hide()
	}
	if u.modal != nil {
		u.modal.DeactivateWidgets()
		u.modal = nil
	}
}

// Update is called once per tick and updates the UI's state
func (u *UI) Update() {
	u.toolbar.Update()
	u.navdrawer.Update()
	u.toastManager.Update()
	if u.modal != nil {
		u.modal.Update()
	}
}

// Draw is called once per tick and renders the UI to the screen
func (u *UI) Draw(screen *ebiten.Image) {
	u.toolbar.Draw(screen)
	u.navdrawer.Draw(screen)
	u.toastManager.Draw(screen)
	if u.modal != nil {
		u.modal.Draw(screen)
	}
}
