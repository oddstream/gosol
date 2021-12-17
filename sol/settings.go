package sol

func ShowSettingsDrawer() {
	// TODO this pattern is well ugly
	// consider using callbacks so UI can query each setting
	var booleanSettings = map[string]bool{
		"FixedCards":  ThePreferences.FixedCards,
		"PowerMoves":  ThePreferences.PowerMoves,
		"MirrorBaize": ThePreferences.MirrorBaize,
		"MuteSounds":  ThePreferences.Mute,
	}
	TheUI.ShowSettingsDrawer(booleanSettings)
}
