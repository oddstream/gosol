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
		NewNavItem(n, input, -256, 0, 256, 48, 0, rune(0x2605), "New deal", ebiten.KeyN),
		NewNavItem(n, input, -256, 0, 256, 48, 0, rune(0x267b), "Restart deal", ebiten.KeyR),
		NewNavItem(n, input, -256, 0, 256, 48, 0, rune(0x2618), "Find game...", ebiten.KeyF),
		NewNavItem(n, input, -256, 0, 256, 48, 0, rune(0x2696), "Rules...", ebiten.KeyF1),
		NewNavItem(n, input, -256, 0, 256, 48, 0, rune(0x2611), "Settings...", ebiten.KeyHome),
		NewNavItem(n, input, -256, 0, 256, 48, 0, rune('x'), "Save and exit", ebiten.KeyX),
	}
	// n.widgets[2].Deactivate()
	// n.widgets[3].Deactivate()
	// n.widgets[4].Deactivate()
	return n
}

// LayoutWidgets belonging to this container
// func (n *NavDrawer) LayoutWidgets() {

// 	var y int = 64
// 	for _, w := range n.widgets {
// 		w.SetPosition(n.x, n.y+y)
// 		y += 48
// 	}

// }

// ShowNavDrawer animates the drawer on/off screen to the left
// func (u *UI) ShowNavDrawer() {

// 	con := u.VisibleDrawer()
// 	if con == u.navdrawer {
// 		return
// 	}
// 	if con != nil {
// 		con.Hide()
// 	}
// 	u.navdrawer.Show()

// }

// ToggleNavDrawer animates the drawer on/off screen to the left
func (u *UI) ToggleNavDrawer() {

	con := u.VisibleDrawer()
	if con == u.navdrawer {
		con.Hide()
		return
	}
	if con == nil {
		u.navdrawer.Show()
		return
	}
	con.Hide()
	u.navdrawer.Show()
}
