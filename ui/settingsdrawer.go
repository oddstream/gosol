package ui

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
func (u *UI) ShowSettingsDrawer(booleanSettings map[string]bool) {
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
		NewCheckbox(u.settingsDrawer, "Fixed cards", booleanSettings["FixedCards"]),
		NewCheckbox(u.settingsDrawer, "Power moves", booleanSettings["PowerMoves"]),
		NewCheckbox(u.settingsDrawer, "Four colors", booleanSettings["FourColors"]),
		NewCheckbox(u.settingsDrawer, "Mirror baize", booleanSettings["MirrorBaize"]),
		NewCheckbox(u.settingsDrawer, "Mute sounds", booleanSettings["Mute"]),
	}
	u.settingsDrawer.LayoutWidgets()
	u.settingsDrawer.Show()
}
