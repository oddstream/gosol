package ui

import (
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
	d.widgets = []Widget{
		// NewLabel(n, input, 0, -100, 256, 48, 0, "Title", schriftbank.RobotRegular24, ""),
		// give -ve x to make sure item is initially drawn off screen
		// y will be set by LayoutWidgets()
		NewCheckbox(d, input, "Retro cards", false),
		NewCheckbox(d, input, "Highlight", true),
	}
	d.LayoutWidgets()
	return d
}

// ShowSettingsDrawer makes the card back picker visible
func (u *UI) ShowSettingsDrawer() {
	con := u.VisibleDrawer()
	if con == u.settingsDrawer {
		return
	}
	if con != nil {
		con.Hide()
	}
	u.settingsDrawer.Show()
}
