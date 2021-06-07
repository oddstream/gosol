package sol

func (b *Baize) ShowSettingsDrawer() {
	b.ui.ShowSettingsDrawer(TheUserData.CardStyle, TheUserData.SingleTap, TheUserData.HighlightMovable, TheUserData.PowerMoves, TheUserData.MuteSounds)
}
