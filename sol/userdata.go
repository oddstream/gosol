package sol

// UserData contains the settings and preferences for the user
type UserData struct {
	// Capitals to emit to json
	Game                      string
	Variant                   string
	CardBackPattern           string
	CardBackColor             string
	CardStyle                 string
	SingleTap                 bool
	HighlightMovable          bool
	PowerMoves                bool
	MuteSounds                bool
	WindowX, WindowY          int
	WindowWidth, WindowHeight int
	WindowMaximized           bool
}
