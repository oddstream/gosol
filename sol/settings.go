package sol

func (b *Baize) ShowSettingsDrawer() {
	retro := TheUserData.CardStyle == "retro"
	b.ui.ShowSettingsDrawer(retro, TheUserData.HighlightMovable)
}
