package sol

// Preferences contains the settings and preferences for the user
type Preferences struct {
	// Capitals to emit to json
	Title                           string
	Variant                         string
	BaizeColor                      string
	CardFaceColor                   string
	CardBackColor                   string
	BlackColor                      string
	RedColor                        string
	ClubColor                       string
	DiamondColor                    string
	HeartColor                      string
	SpadeColor                      string
	FourColors                      bool
	FixedCards                      bool
	PowerMoves                      bool
	Relaxed                         bool
	Mute                            bool
	Volume                          float64
	MirrorBaize                     bool
	PreferredWindow                 bool
	CardRatio                       float64
	FixedCardWidth, FixedCardHeight int
}

// ThePreferences holds serialized game progress data
// Colors are named from the web extended colors at https://en.wikipedia.org/wiki/Web_colors
var ThePreferences = &Preferences{
	Title:           "Solitaire",
	Variant:         "Klondike",
	BaizeColor:      "BaizeGreen",
	PowerMoves:      true,
	CardFaceColor:   "Ivory",
	CardBackColor:   "CornflowerBlue",
	FourColors:      false,
	RedColor:        "Crimson",
	BlackColor:      "Black",
	ClubColor:       "Indigo",
	DiamondColor:    "OrangeRed",
	HeartColor:      "Crimson",
	SpadeColor:      "Black",
	FixedCards:      true,
	Mute:            false,
	Volume:          1.0,
	FixedCardWidth:  90,
	FixedCardHeight: 122,
	CardRatio:       1.357,
}
