package sol

func ShowSettingsDrawer() {
	// TODO this pattern is well ugly
	// consider using callbacks so UI can query each setting
	var booleanSettings = map[string]bool{
		// "FixedCards":    ThePreferences.FixedCards,
		"PowerMoves":       ThePreferences.PowerMoves,
		"AutoCollect":      ThePreferences.AutoCollect,
		"SafeCollect":      ThePreferences.SafeCollect,
		"ColorfulCards":    ThePreferences.ColorfulCards,
		"MirrorBaize":      ThePreferences.MirrorBaize,
		"ShowMovableCards": ThePreferences.ShowMovableCards,
		"Mute":             ThePreferences.Mute,
	}
	TheUI.ShowSettingsDrawer(booleanSettings)
}

func ShowAniSpeedDrawer() {
	// ThePreferences.AniSpeed is a float64 0.75, 0.5, 0.25
	TheUI.ShowAniSpeedDrawer(ThePreferences.AniSpeed)
}
