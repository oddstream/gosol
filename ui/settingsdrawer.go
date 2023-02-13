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
	d := &SettingsDrawer{DrawerBase: DrawerBase{WindowBase: WindowBase{width: 360, height: 0, x: -400, y: ToolbarHeight}}}
	return d
}

type BooleanSetting struct {
	Title  string
	Var    *bool
	Update func()
}

// ShowSettingsDrawer makes the card back picker visible
func (u *UI) ShowSettingsDrawer(booleanSettings *[]BooleanSetting) {
	con := u.VisibleDrawer()
	if con == u.settingsDrawer {
		return
	}
	if con != nil {
		con.Hide()
	}
	u.settingsDrawer.widgets = u.settingsDrawer.widgets[:0]
	u.settingsDrawer.widgets = []Widgety{
		// widget x, y will be set by LayoutWidgets()
		NewNavItem(u.settingsDrawer, "", "speed", "Card speed...", ebiten.KeyA),
	}
	for _, p := range *booleanSettings {
		u.settingsDrawer.widgets = append(u.settingsDrawer.widgets, NewCheckbox(u.settingsDrawer, "", p.Title, p.Var, p.Update))
	}
	u.settingsDrawer.LayoutWidgets()
	u.settingsDrawer.Show()
}

type FloatSetting struct {
	Title string
	Var   *float64
	Value float64
}

func (u *UI) ShowAniSpeedDrawer(floatSettings *[]FloatSetting) {
	con := u.VisibleDrawer()
	if con == u.aniSpeedDrawer {
		return
	}
	if con != nil {
		con.Hide()
	}
	u.aniSpeedDrawer.widgets = u.settingsDrawer.widgets[:0]
	u.aniSpeedDrawer.widgets = []Widgety{
		NewText(u.aniSpeedDrawer, "aniTitle", "Card Animation Speed"),
	}
	for _, p := range *floatSettings {
		u.aniSpeedDrawer.widgets = append(u.aniSpeedDrawer.widgets, NewRadioButton(u.aniSpeedDrawer, "", p.Title, p.Var, p.Value))
	}
	u.aniSpeedDrawer.LayoutWidgets()
	u.aniSpeedDrawer.Show()
}
