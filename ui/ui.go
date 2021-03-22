// Package ui provides a minimal user interface for package sol
package ui

import (
	_ "embed" // go:embed only allowed in Go files that import "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
)

var (
	GenerateIcons bool = false
)

// UI encapsulates a complete user interface that can be rendered onto the screen.
type UI struct {
	input          *input.Input // place to receive clicks, taps and key presses from
	toolbar        *Toolbar
	navDrawer      *NavDrawer
	variantPicker  *Picker
	cardBackPicker *CardBackPicker
	rules          *Rules
	fab            *FAB
	containers     []Container
	bars           []Container
	drawers        []Container
	toastManager   *ToastManager
}

// New creates a new UI object
func New(input *input.Input, pickerContents []string, cardBackPickerContents map[string]*ebiten.Image) *UI {
	ui := &UI{input: input}

	LoadIconMap()

	ui.toastManager = &ToastManager{}
	ui.toolbar = NewToolbar(input)
	ui.navDrawer = NewNavDrawer(input)
	ui.variantPicker = NewPicker(input, pickerContents)
	ui.cardBackPicker = NewCardBackPicker(input, cardBackPickerContents)
	ui.rules = NewRules(input) // contents are added when shown

	ui.bars = []Container{ui.toolbar}
	ui.drawers = []Container{ui.navDrawer, ui.variantPicker, ui.cardBackPicker, ui.rules}
	ui.containers = []Container{ui.toolbar, ui.navDrawer, ui.variantPicker, ui.cardBackPicker, ui.rules}

	return ui
}

// FindWidgetAt finds the widget at the screen coords
// func (u *UI) FindWidgetAt(x, y int) Widget {
// 	for _, con := range u.containers {
// 		if w := con.FindWidgetAt(x, y); con != nil {
// 			return w
// 		}
// 	}
// 	// if u.fab != nil {
// 	// 	if util.InRect(x, y, u.fab.Rect) {
// 	// 		return u.fab
// 	// 	}
// 	// }
// 	return nil
// }

// func (u *UI) FindContainerAt(x, y int) Container {
// 	for _, con := range u.containers {
// 		if util.InRect(x, y, con.Rect) {
// 			return con
// 		}
// 	}
// 	return nil
// }

// VisibleDrawer returns the drawer that is currently open
func (u *UI) VisibleDrawer() Container {
	for _, con := range u.drawers {
		if con.Visible() {
			return con
		}
	}
	return nil
}

// HideActiveDrawer closes the active/open drawer
func (u *UI) HideActiveDrawer() {
	if con := u.VisibleDrawer(); con != nil {
		con.Hide()
	}
}

// Layout implements Ebiten's Layout
func (u *UI) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

// Update is called once per tick and updates the UI's state
func (u *UI) Update() {
	for _, con := range u.containers {
		con.Update()
	}
	if u.fab != nil {
		u.fab.Update()
	}
	u.toastManager.Update()
}

// Draw is called once per tick and renders the UI to the screen
func (u *UI) Draw(screen *ebiten.Image) {
	for _, con := range u.containers {
		con.Draw(screen)
	}
	if u.fab != nil {
		u.fab.Draw(screen)
	}
	u.toastManager.Draw(screen)
}
