package sol

import (
	"oddstream.games/gosol/sound"
	"oddstream.games/gosol/ui"
)

// Settings contains the settings and preferences for the user
type Settings struct {
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

// TheSettings holds serialized game progress data
// Colors are named from the web extended colors at https://en.wikipedia.org/wiki/Web_colors
var TheSettings = &Settings{
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

var BooleanSettings = []ui.BooleanSetting{
	{Title: "Power moves", Var: &TheSettings.PowerMoves, Update: func() {}},
	{Title: "Auto collect", Var: &TheSettings.AutoCollect, Update: func() {}},
	{Title: "Safe collect", Var: &TheSettings.SafeCollect, Update: func() {}},
	{Title: "Show movable cards", Var: &TheSettings.ShowMovableCards, Update: func() {}},
	{Title: "Colorful cards", Var: &TheSettings.ColorfulCards, Update: func() { TheBaize.setFlag(dirtyCardImages) }},
	{Title: "Mute sounds", Var: &TheSettings.Mute, Update: func() {
		if TheSettings.Mute {
			sound.SetVolume(0.0)
		} else {
			sound.SetVolume(TheSettings.Volume)
		}
	}},
	{Title: "Mirror baize", Var: &TheSettings.MirrorBaize, Update: func() {
		savedUndoStack := TheBaize.undoStack
		TheBaize.StartFreshGame()
		TheBaize.SetUndoStack(savedUndoStack)
	}},
}

func ShowSettingsDrawer() {
	TheUI.ShowSettingsDrawer(&BooleanSettings)
}

var AniSpeedSettings = []ui.FloatSetting{
	{Title: "Fast", Var: &TheSettings.AniSpeed, Value: 0.3},
	{Title: "Normal", Var: &TheSettings.AniSpeed, Value: 0.6},
	{Title: "Slow", Var: &TheSettings.AniSpeed, Value: 0.9},
}

func ShowAniSpeedDrawer() {
	TheUI.ShowAniSpeedDrawer(&AniSpeedSettings)
}
