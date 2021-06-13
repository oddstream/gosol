package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// SettingsDrawer slide out modal menu
type SettingsDrawer struct {
	DrawerBase
}

// NewSettingsDrawer creates the SettingsDrawer object; it starts life off screen to the left
func NewSettingsDrawer() *SettingsDrawer {
	// according to https://material.io/components/navigation-drawer#specs, always 256 wide
	d := &SettingsDrawer{DrawerBase: DrawerBase{width: 256, height: 0, x: -256, y: 48}}
	return d
}

// ShowSettingsDrawer makes the card back picker visible
// TODO pass parameters better than this
func (u *UI) ShowSettingsDrawer(retroCards, fixedCards, singletap, highlight, powerMoves, muteSounds bool) {
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
		NewCheckbox(u.settingsDrawer, "Retro cards", retroCards),
		NewCheckbox(u.settingsDrawer, "Fixed cards", fixedCards),
		NewNavItem(u.settingsDrawer, "settings", "Card back ...", ebiten.KeyF2),
		NewCheckbox(u.settingsDrawer, "Single tap", singletap),
		NewCheckbox(u.settingsDrawer, "Highlights", highlight),
		NewCheckbox(u.settingsDrawer, "Power moves", powerMoves),
		NewCheckbox(u.settingsDrawer, "Mute sounds", muteSounds),
	}
	u.settingsDrawer.LayoutWidgets()
	u.settingsDrawer.Show()
}
