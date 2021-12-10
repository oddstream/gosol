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
	MuteSounds                bool
	MirrorBaize               bool
	WindowX, WindowY          int
	WindowWidth, WindowHeight int
}
