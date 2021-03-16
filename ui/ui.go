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
	toolbar      *Toolbar
	navdrawer    *NavDrawer
	picker       *Picker
	rules        *Rules
	containers   []Container
	bars         []Container
	drawers      []Container
	toastManager *ToastManager
}

// New creates a new UI object
func New(input *input.Input, pickerContents []string) *UI {
	ui := &UI{input: input}

	ui.toastManager = &ToastManager{}
	ui.toolbar = NewToolbar(input)
	ui.navdrawer = NewNavDrawer(input)
	ui.picker = NewPicker(input, pickerContents)
	ui.rules = NewRules(input, nil)

	ui.bars = []Container{ui.toolbar}
	ui.drawers = []Container{ui.navdrawer, ui.picker, ui.rules}
	ui.containers = []Container{ui.toolbar, ui.navdrawer, ui.picker, ui.rules}

	return ui
}

// FindWidgetAt finds the widget at the screen coords
func (u *UI) FindWidgetAt(x, y int) Widget {
	for _, con := range u.containers {
		if w := con.FindWidgetAt(x, y); con != nil {
			return w
		}
	}
	return nil
}

func (u *UI) FindContainerAt(x, y int) Container {
	for _, con := range u.containers {
		if util.InRect(x, y, con.Rect) {
			return con
		}
	}
	return nil
}

// VisibleDrawer returns the drawer that is currently open
func (u *UI) VisibleDrawer() Container {
	for _, con := range u.drawers {
		if con.Visible() {
			return con
		}
	}
	return nil
}

// InVisibleDrawRect
// func (u *UI) InVisibleDrawerRect(x, y int) bool {
// 	for _, con := range u.drawers {
// 		if con.Visible() {
// 			return util.InRect(x, y, con.Rect)
// 		}
// 	}
// 	return false
// }

// ActiveDrawerRect returns the rect coords of the active drawer
// func (u *UI) ActiveDrawerRect() (int, int, int, int) {
// 	d := u.VisibleDrawer()
// 	if d != nil {
// 		return d.Rect()
// 	}
// 	return 0, 0, 0, 0
// }

// HideActiveDrawer closes the active/open drawer
func (u *UI) HideActiveDrawer() {
	if con := u.VisibleDrawer(); con != nil {
		con.Hide()
	}
}

// Update is called once per tick and updates the UI's state
func (u *UI) Update() {
	for _, con := range u.containers {
		con.Update()
	}
	u.toastManager.Update()
}

// Draw is called once per tick and renders the UI to the screen
func (u *UI) Draw(screen *ebiten.Image) {
	for _, con := range u.containers {
		con.Draw(screen)
	}
	u.toastManager.Draw(screen)
}
