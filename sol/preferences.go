package sol

// Preferences contains the settings and preferences for the user
type Preferences struct {
	// Capitals to emit to json
	Title                string
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
	Mute                               bool
	Volume                             float64
	MirrorBaize                        bool
	PreferredWindow                    bool
	CardRatio                          float64
	LastVersionMajor, LastVersionMinor int
}

// ThePreferences holds serialized game progress data
// Colors are named from the web extended colors at https://en.wikipedia.org/wiki/Web_colors
var ThePreferences = &Preferences{
	Title:                "Solitaire",
	Variant:              "Klondike",
	BaizeColor:           "BaizeGreen",
	PowerMoves:           true,
	CardFaceColor:        "Ivory",
	CardBackColor:        "CornflowerBlue",
	MovableCardBackColor: "Gold",
	ColorfulCards:        true,
	RedColor:             "Crimson",
	BlackColor:           "Black",
	ClubColor:            "DarkGreen",
	DiamondColor:         "DarkBlue",
	HeartColor:           "Crimson",
	SpadeColor:           "Black",
	Mute:                 false,
	Volume:               0.75,
	// FixedCards:       false,
	// FixedCardWidth:   90,
	// FixedCardHeight:  122,
	CardRatio:        1.357,
	LastVersionMajor: 0,
	LastVersionMinor: 0,
}
