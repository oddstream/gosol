package ui

import (
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
)

// NavDrawer slide out modal menu
type NavDrawer struct {
	DrawerBase
}

// NewNavDrawer creates the NavDrawer object; it starts life off screen to the left
func NewNavDrawer() *NavDrawer {
	// according to https://material.io/components/navigation-drawer#specs, always 256 wide
	nd := &NavDrawer{
		DrawerBase: DrawerBase{
			WindowBase: WindowBase{width: 300, height: 0, x: -300, y: ToolbarHeight},
		},
	}
	nd.widgets = []Widgety{
		// widget x, y will be set by LayoutWidgets()
		NewNavItem(nd, "newDeal", "star", "New deal", ebiten.KeyN),
		NewNavItem(nd, "restartDeal", "restore", "Restart deal", ebiten.KeyR),
		NewNavItem(nd, "findGame", "search", "Find game...", ebiten.KeyF),
		NewNavItem(nd, "bookmark", "bookmark_add", "Set bookmark", ebiten.KeyS),
		NewNavItem(nd, "gotoBookmark", "bookmark", "Go to bookmark", ebiten.KeyL),
		NewNavItem(nd, "wikipedia", "wikipedia", "Wikipedia...", ebiten.KeyF1),
		NewNavItem(nd, "statistics", "poll", "Statistics...", ebiten.KeyF2),
		NewNavItem(nd, "settings", "settings", "Settings...", ebiten.KeyF3),
	}
	// don't know how to ask a browser window to close
	if runtime.GOARCH != "wasm" {
		nd.widgets = append(nd.widgets, NewNavItem(nd, "exit", "close", "Save and exit", ebiten.KeyX))
	}
	nd.LayoutWidgets()
	return nd
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
