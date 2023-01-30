package ui

import "strings"

// TextDrawer provides a drawer for displaying rules of the current variant
type TextDrawer struct {
	DrawerBase
}

// NewTextDrawer creates a new container
func NewTextDrawer() *TextDrawer {
	r := &TextDrawer{DrawerBase: DrawerBase{WindowBase: WindowBase{x: -400, y: ToolbarHeight, width: 400}}} // height will be set when drawn
	return r
}

// ShowTextDrawer makes the text drawer container visible
func (u *UI) ShowTextDrawer(content []string) {
	con := u.VisibleDrawer()
	if con == u.textDrawer {
		return
	}
	if con != nil {
		con.Hide()
	}

	u.textDrawer.widgets = nil
	for _, c := range content { // content may be nil
		if strings.HasPrefix(c, "https://") {
			u.textDrawer.widgets = append(u.textDrawer.widgets, NewTextUrl(u.textDrawer, "", c))
		} else {
			u.textDrawer.widgets = append(u.textDrawer.widgets, NewText(u.textDrawer, "", c))
		}
	}
	u.textDrawer.LayoutWidgets()
	u.textDrawer.Show()
}
