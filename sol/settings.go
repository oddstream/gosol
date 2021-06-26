package sol

func (b *Baize) ShowSettingsDrawer() {
	b.ui.ShowSettingsDrawer(ThePreferences.RetroCards, ThePreferences.FixedCards, ThePreferences.SingleTap, ThePreferences.HighlightMovable, ThePreferences.PowerMoves, ThePreferences.MuteSounds)
}
