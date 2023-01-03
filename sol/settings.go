package sol

func ShowSettingsDrawer() {
	// TODO this pattern is well ugly
	// consider using callbacks so UI can query each setting
	var booleanSettings = map[string]bool{
		// "FixedCards":    ThePreferences.FixedCards,
		"PowerMoves":       ThePreferences.PowerMoves,
		"ColorfulCards":    ThePreferences.ColorfulCards,
		"MirrorBaize":      ThePreferences.MirrorBaize,
		"ShowMovableCards": ThePreferences.ShowMovableCards,
		"Mute":             ThePreferences.Mute,
	}
	TheUI.ShowSettingsDrawer(booleanSettings)
}
