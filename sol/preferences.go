package sol

// Preferences contains the settings and preferences for the user
type Preferences struct {
	// Capitals to emit to json
	Game                      string
	Variant                   string
	CardBackPattern           string
	CardBackColor             string
	RetroCards                bool
	FixedCards                bool
	SingleTap                 bool
	HighlightMovable          bool
	PowerMoves                bool
	MuteSounds                bool
	WindowX, WindowY          int
	WindowWidth, WindowHeight int
	WindowMaximized           bool
}
