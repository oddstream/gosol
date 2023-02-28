package sol

import (
	"oddstream.games/gosol/sound"
	"oddstream.games/gosol/ui"
)

// Settings holds user preferences.
// Colors are named from the web extended colors at https://en.wikipedia.org/wiki/Web_colors
type Settings struct {
	// Capitals to emit to json
	Variant                            string
	BaizeColor                         string
	CardFaceColor                      string
	CardBackColor                      string
	MovableCardBackColor               string
	BlackColor                         string
	RedColor                           string
	ClubColor                          string
	DiamondColor                       string
	HeartColor                         string
	SpadeColor                         string
	ColorfulCards                      bool
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
	// FixedCards                         bool
	// FixedCardWidth, FixedCardHeight    int
}

func NewSettings() *Settings {
	s := &Settings{
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
	s.Load()
	return s
}

func ShowSettingsDrawer() {
	var BooleanSettings = []ui.BooleanSetting{
		{Title: "Power moves", Var: &TheGame.Settings.PowerMoves},
		{Title: "Auto collect", Var: &TheGame.Settings.AutoCollect},
		{Title: "Safe collect", Var: &TheGame.Settings.SafeCollect},
		{Title: "Show movable cards", Var: &TheGame.Settings.ShowMovableCards},
		{Title: "Colorful cards", Var: &TheGame.Settings.ColorfulCards, Update: func() { TheGame.Baize.setFlag(dirtyCardImages) }},
		{Title: "Mute sounds", Var: &TheGame.Settings.Mute, Update: func() {
			if TheGame.Settings.Mute {
				sound.SetVolume(0.0)
			} else {
				sound.SetVolume(TheGame.Settings.Volume)
			}
		}},
		{Title: "Mirror baize", Var: &TheGame.Settings.MirrorBaize, Update: func() {
			savedUndoStack := TheGame.Baize.undoStack
			TheGame.Baize.StartFreshGame()
			TheGame.Baize.SetUndoStack(savedUndoStack)
		}},
	}

	TheGame.UI.ShowSettingsDrawer(&BooleanSettings)
}

func ShowAniSpeedDrawer() {
	var AniSpeedSettings = []ui.FloatSetting{
		{Title: "Fast", Var: &TheGame.Settings.AniSpeed, Value: 0.3},
		{Title: "Normal", Var: &TheGame.Settings.AniSpeed, Value: 0.6},
		{Title: "Slow", Var: &TheGame.Settings.AniSpeed, Value: 0.9},
	}

	TheGame.UI.ShowAniSpeedDrawer(&AniSpeedSettings)
}
