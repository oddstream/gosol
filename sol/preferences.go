package sol

// Preferences contains the settings and preferences for the user
type Preferences struct {
	// Capitals to emit to json
	Game                      string
	BaizeColor                string
	Variant                   string
	CardFaceColor             string
	CardBackColor             string
	FixedCards                bool
	PowerMoves                bool
	Mute                      bool
	Volume                    float64
	MirrorBaize               bool
	WindowX, WindowY          int
	WindowWidth, WindowHeight int
}
