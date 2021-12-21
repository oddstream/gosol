package sol

// Preferences contains the settings and preferences for the user
type Preferences struct {
	// Capitals to emit to json
	Title                     string
	BaizeColor                string
	Variant                   string
	CardFaceColor             string
	CardBackColor             string
	BlackColor                string
	RedColor                  string
	ClubColor                 string
	DiamondColor              string
	HeartColor                string
	SpadeColor                string
	ExtraColors               bool
	FixedCards                bool
	PowerMoves                bool
	Mute                      bool
	Volume                    float64
	MirrorBaize               bool
	WindowX, WindowY          int
	WindowWidth, WindowHeight int
}
