package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
)

// NavDrawer slide out modal menu
type NavDrawer struct {
	DrawerBase
}

// NewNavDrawer creates the NavDrawer object; it starts life off screen to the left
func NewNavDrawer(input *input.Input) *NavDrawer {
	// according to https://material.io/components/navigation-drawer#specs, always 256 wide
	n := &NavDrawer{DrawerBase: DrawerBase{input: input, width: 256, height: 0, x: -256, y: 48}}
	n.widgets = []Widget{
		// NewLabel(n, input, 0, -100, 256, 48, 0, "Title", schriftbank.RobotoMedium24, ""),
		// give -ve x to make sure item is initially drawn off screen
		// y will be set by LayoutWidgets()
		NewNavItem(n, input, "star", "New deal", ebiten.KeyN),
		NewNavItem(n, input, "restore", "Restart deal", ebiten.KeyR),
		NewNavItem(n, input, "search", "Find game...", ebiten.KeyF),
		NewNavItem(n, input, "bookmark_add", "Bookmark", ebiten.KeyS),
		NewNavItem(n, input, "bookmark", "Goto bookmark", ebiten.KeyL),
		NewNavItem(n, input, "list", "Rules...", ebiten.KeyF1),
		NewNavItem(n, input, "settings", "Settings...", ebiten.KeyHome),
		NewNavItem(n, input, "close", "Save and exit", ebiten.KeyX),
	}
	n.LayoutWidgets()
	n.widgets[6].Deactivate()
	return n
}

// ToggleNavDrawer animates the drawer on/off screen to the left
func (u *UI) ToggleNavDrawer() {

	con := u.VisibleDrawer()
	if con == u.navDrawer {
		con.Hide()
		return
	}
	if con == nil {
		u.navDrawer.Show()
		return
	}
	con.Hide()
	u.navDrawer.Show()
}
