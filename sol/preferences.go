package sol

import (
	"oddstream.games/gosol/sound"
	"oddstream.games/gosol/ui"
)

// Preferences contains the settings and preferences for the user
type Preferences struct {
	// Capitals to emit to json
	Variant              string
	BaizeColor           string
	CardFaceColor        string
	CardBackColor        string
	MovableCardBackColor string
	BlackColor           string
	RedColor             string
	ClubColor            string
	DiamondColor         string
	HeartColor           string
	SpadeColor           string
	ColorfulCards        bool
	// FixedCards                         bool
	// FixedCardWidth, FixedCardHeight    int
	PowerMoves                         bool
	SafeCollect, AutoCollect           bool
	Mute                               bool
	Volume                             float64
	MirrorBaize                        bool
	ShowMovableCards                   bool
	AlwaysShowMovableCards             bool
	CardRatio                          float64
	AniSpeed                           float64
	LastVersionMajor, LastVersionMinor int
}

// ThePreferences holds serialized game progress data
// Colors are named from the web extended colors at https://en.wikipedia.org/wiki/Web_colors
var ThePreferences = &Preferences{
	Variant:                "Klondike",
	BaizeColor:             "BaizeGreen",
	PowerMoves:             true,
	SafeCollect:            false,
	AutoCollect:            false,
	CardFaceColor:          "Ivory",
	CardBackColor:          "CornflowerBlue",
	MovableCardBackColor:   "Gold",
	ColorfulCards:          true,
	RedColor:               "Crimson",
	BlackColor:             "Black",
	ClubColor:              "DarkGreen",
	DiamondColor:           "DarkBlue",
	HeartColor:             "Crimson",
	SpadeColor:             "Black",
	Mute:                   false,
	Volume:                 0.75,
	ShowMovableCards:       false,
	AlwaysShowMovableCards: false,
	// FixedCards:       false,
	// FixedCardWidth:   90,
	// FixedCardHeight:  122,
	CardRatio:        1.39, // official poker size
	AniSpeed:         0.6,  // Normal
	LastVersionMajor: 0,
	LastVersionMinor: 0,
}

var BooleanPreferences = []ui.BooleanPreference{
	{Title: "Power moves", Var: &ThePreferences.PowerMoves, Update: func() {}},
	{Title: "Auto collect", Var: &ThePreferences.AutoCollect, Update: func() {}},
	{Title: "Safe collect", Var: &ThePreferences.SafeCollect, Update: func() {}},
	{Title: "Show movable cards", Var: &ThePreferences.ShowMovableCards, Update: func() {}},
	{Title: "Colorful cards", Var: &ThePreferences.ColorfulCards, Update: func() { TheBaize.setFlag(dirtyCardImages) }},
	{Title: "Mute sounds", Var: &ThePreferences.Mute, Update: func() {
		if ThePreferences.Mute {
			sound.SetVolume(0.0)
		} else {
			sound.SetVolume(ThePreferences.Volume)
		}
	}},
	{Title: "Mirror baize", Var: &ThePreferences.MirrorBaize, Update: func() {
		savedUndoStack := TheBaize.undoStack
		TheBaize.StartFreshGame()
		TheBaize.SetUndoStack(savedUndoStack)
	}},
}

func ShowSettingsDrawer() {
	TheUI.ShowSettingsDrawer(&BooleanPreferences)
}

var AniSpeedPreferences = []ui.FloatPreference{
	{Title: "Fast", Var: &ThePreferences.AniSpeed, Value: 0.3},
	{Title: "Normal", Var: &ThePreferences.AniSpeed, Value: 0.6},
	{Title: "Slow", Var: &ThePreferences.AniSpeed, Value: 0.9},
}

func ShowAniSpeedDrawer() {
	TheUI.ShowAniSpeedDrawer(&AniSpeedPreferences)
}
