package sol

func (b *Baize) ShowSettingsDrawer() {
	b.ui.ShowSettingsDrawer(TheUserData.RetroCards, TheUserData.FixedCards, TheUserData.SingleTap, TheUserData.HighlightMovable, TheUserData.PowerMoves, TheUserData.MuteSounds)
}
