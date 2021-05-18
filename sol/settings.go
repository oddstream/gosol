package sol

func (b *Baize) ShowSettingsDrawer() {
	b.ui.ShowSettingsDrawer(TheUserData.CardStyle, TheUserData.HighlightMovable, TheUserData.PowerMoves, TheUserData.MuteSounds)
}
