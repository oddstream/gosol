package sol

func ShowSettingsDrawer() {
	// TODO this smells real bad
	TheUI.ShowSettingsDrawer(ThePreferences.RetroCards, ThePreferences.FixedCards, ThePreferences.SingleTap, ThePreferences.HighlightMovable, ThePreferences.PowerMoves, ThePreferences.MuteSounds)
}
