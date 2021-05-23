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
func (u *UI) ShowSettingsDrawer(cardStyle string, highlight, powerMoves, muteSounds bool) {
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
		NewRadioButton(u.settingsDrawer, "Scaled cards", cardStyle == "scaled"),
		NewRadioButton(u.settingsDrawer, "Fixed cards", cardStyle == "fixed"),
		NewRadioButton(u.settingsDrawer, "Retro cards", cardStyle == "retro"),
		NewNavItem(u.settingsDrawer, "settings", "Card back ...", ebiten.KeyF2),
		NewCheckbox(u.settingsDrawer, "Highlights", highlight),
		NewCheckbox(u.settingsDrawer, "Power moves", powerMoves),
		NewCheckbox(u.settingsDrawer, "Mute sounds", muteSounds),
	}
	u.settingsDrawer.LayoutWidgets()
	u.settingsDrawer.Show()
}
