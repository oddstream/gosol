package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
)

// SettingsDrawer slide out modal menu
type SettingsDrawer struct {
	DrawerBase
}

// NewSettingsDrawer creates the SettingsDrawer object; it starts life off screen to the left
func NewSettingsDrawer(input *input.Input) *SettingsDrawer {
	// according to https://material.io/components/navigation-drawer#specs, always 256 wide
	d := &SettingsDrawer{DrawerBase: DrawerBase{input: input, width: 256, height: 0, x: -256, y: 48}}
	return d
}

// ShowSettingsDrawer makes the card back picker visible
func (u *UI) ShowSettingsDrawer(retro, highlight, powerMoves bool) {
	con := u.VisibleDrawer()
	if con == u.settingsDrawer {
		return
	}
	if con != nil {
		con.Hide()
	}
	u.settingsDrawer.widgets = u.settingsDrawer.widgets[:0]
	u.settingsDrawer.widgets = []Widget{
		// widget x, y will be set by LayoutWidgets()
		NewCheckbox(u.settingsDrawer, u.input, "Highlight", highlight),
		NewCheckbox(u.settingsDrawer, u.input, "Power moves", powerMoves),
		NewCheckbox(u.settingsDrawer, u.input, "Retro cards", retro),
		NewNavItem(u.settingsDrawer, u.input, "settings", "Card back ...", ebiten.KeyF2),
	}
	u.settingsDrawer.LayoutWidgets()
	u.settingsDrawer.Show()
}
