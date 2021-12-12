// Package ui provides a minimal user interface for package sol
package ui

import (
	_ "embed" // go:embed only allowed in Go files that import "embed"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gomps5/util"
)

var (
	GenerateIcons   bool        = false
	BackgroundColor color.Color = color.RGBA{R: 0x24, G: 0x24, B: 0x24, A: 0xee}
)

// UI encapsulates a complete user interface that can be rendered onto the screen.
type UI struct {
	toolbar        *Toolbar
	statusbar      *Statusbar
	fabbar         *FABBar
	navDrawer      *NavDrawer
	settingsDrawer *SettingsDrawer
	variantPicker  *Picker
	textDrawer     *TextDrawer
	containers     []Container
	bars           []Container
	drawers        []Container
	toastManager   *ToastManager
}

var cmdFn func(interface{})

// New creates a new UI object
func New(fn func(interface{})) *UI {
	cmdFn = fn
	ui := &UI{}

	LoadIconMap()

	ui.toastManager = &ToastManager{}
	ui.toolbar = NewToolbar()
	ui.statusbar = NewStatusbar()
	ui.fabbar = NewFABBar()
	ui.navDrawer = NewNavDrawer()
	ui.settingsDrawer = NewSettingsDrawer()
	ui.variantPicker = NewVariantPicker()
	ui.textDrawer = NewTextDrawer() // contents are added when shown

	ui.bars = []Container{ui.toolbar, ui.statusbar, ui.fabbar}
	ui.drawers = []Container{ui.navDrawer, ui.settingsDrawer, ui.variantPicker, ui.textDrawer}
	ui.containers = []Container{ui.toolbar, ui.statusbar, ui.fabbar, ui.navDrawer, ui.settingsDrawer, ui.variantPicker, ui.textDrawer}

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

// VisibleContainer returns the drawer that is currently open
func (u *UI) VisibleContainer() Container {
	for _, con := range u.containers {
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
	for _, con := range u.containers {
		con.Layout(outsideWidth, outsideHeight)
	}
	return outsideWidth, outsideHeight
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
